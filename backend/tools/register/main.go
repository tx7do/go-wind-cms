// register 将新 CRUD 模块的登记代码一次性注入手写装配文件,替代原 Wire 时代的
// "改 wire set + make wire" 自动化。注入点为各文件中的 register:* 锚点注释:
//
//	core:  cmd/server/wiring.go                仓储行 / 服务行 / NewGrpcServer 实参
//	       internal/server/grpc_server.go      服务形参 / 路由注册调用
//	admin: cmd/server/wiring.go                服务客户端行 / 服务行 / NewRestServer 实参
//	       internal/server/rest_server.go      服务形参 / 路由注册调用
//	app:   同 admin(路由前缀默认 appV1)
//
// 仅覆盖 CRUD 生成器的标准形态(New<Ent>Repo(ctx, entClient) /
// New<Ent>Service(ctx, <ent>Repo) / New<Ent>Service(ctx, <ent>ServiceClient));
// 依赖更多的模块请手工调整注入行。
//
// 用法:
//	go run ./tools/register -entity product                    # core 仓储+服务+gRPC 注册
//	go run ./tools/register -entity product -svc admin         # admin BFF 客户端+服务+REST 注册
//	go run ./tools/register -entity collection -svc core -pkg contentV1   # 指定路由包别名
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type patch struct {
	path    string   // 目标文件
	anchors []string // 锚点行(按去空白后的全文匹配),插入到首个命中的锚点之后
	lines   []string // 插入的行(不含行尾符)
	skipIf  string   // 文件已包含该片段时跳过(幂等)
}

func main() {
	entity := flag.String("entity", "", "模块名:product / DictEntry / dict_entry 均可")
	svc := flag.String("svc", "core", "目标服务:core / admin / app")
	pkg := flag.String("pkg", "", "路由包别名;admin 默认 adminV1,app 默认 appV1,core 必填(如 contentV1)")
	flag.Parse()

	if *entity == "" {
		flag.Usage()
		os.Exit(2)
	}

	typ := pascal(*entity)
	name := lowerFirst(typ)

	var dir, defaultPkg string
	switch *svc {
	case "core":
		dir = "app/core/service"
	case "admin":
		dir, defaultPkg = "app/admin/service", "adminV1"
	case "app":
		dir, defaultPkg = "app/app/service", "appV1"
	default:
		fmt.Fprintf(os.Stderr, "register: 未知服务 %q(可选 core / admin / app)\n", *svc)
		os.Exit(2)
	}
	routePkg := *pkg
	if routePkg == "" {
		if defaultPkg == "" {
			fmt.Fprintln(os.Stderr, "register: core 服务跨多个 proto 域,请用 -pkg 指定路由包别名(如 -pkg contentV1)")
			os.Exit(2)
		}
		routePkg = defaultPkg
	}

	wiringPath := dir + "/cmd/server/wiring.go"
	serverPath := dir + "/internal/server/grpc_server.go"
	registerFn := "Register" + typ + "ServiceServer"
	if *svc != "core" {
		serverPath = dir + "/internal/server/rest_server.go"
		registerFn = "Register" + typ + "ServiceHTTPServer"
	}
	routeLine := "\t" + routePkg + "." + registerFn + "(srv, " + name + "Service)"

	// 仓储/客户端行与服务行按服务形态区分,其余四类锚点三服务同构。
	repoLine := "\t" + name + "Repo := data.New" + typ + "Repo(ctx, entClient)"
	repoSkip := "data.New" + typ + "Repo("
	serviceLine := "\t" + name + "Service := service.New" + typ + "Service(ctx, " + name + "Repo)"
	serviceSkip := "service.New" + typ + "Service("
	if *svc != "core" {
		repoLine = "\t" + name + "ServiceClient := data.New" + typ + "ServiceClient(ctx, discovery)"
		repoSkip = "data.New" + typ + "ServiceClient("
		serviceLine = "\t" + name + "Service := service.New" + typ + "Service(ctx, " + name + "ServiceClient)"
		serviceSkip = "service.New" + typ + "Service("
	}

	if err := apply(
		// 仓储/客户端构造行:注入到仓储(客户端)小节末尾
		patch{
			path:   wiringPath,
			skipIf: repoSkip,
			anchors: []string{
				"// ── register:repo ── 新模块仓储在此行后注册(make register 工具锚点,勿删)",
				"// ── register:client ── 新模块服务客户端在此行后注册(make register 工具锚点,勿删)",
			},
			lines: []string{repoLine},
		},
		// 服务构造行:注入到服务层小节末尾
		patch{
			path:   wiringPath,
			skipIf: serviceSkip,
			anchors: []string{
				"// ── register:service ── 新模块服务在此行后注册(make register 工具锚点,勿删)",
			},
			lines: []string{serviceLine},
		},
		// NewGrpcServer / NewRestServer 实参
		patch{
			path:   wiringPath,
			skipIf: "\t\t" + name + "Service,",
			anchors: []string{
				"// register:rest-arg ── 新模块服务实参在此行后追加(make register 工具锚点,勿删)",
				"// register:grpc-arg ── 新模块服务实参在此行后追加(make register 工具锚点,勿删)",
			},
			lines: []string{"\t\t" + name + "Service,"},
		},
		// server 文件形参
		patch{
			path:   serverPath,
			skipIf: name + "Service *service." + typ + "Service,",
			anchors: []string{
				"// register:param ── 新模块服务形参在此行后注册(make register 工具锚点,勿删)",
			},
			lines: []string{"\t" + name + "Service *service." + typ + "Service,"},
		},
		// server 文件路由注册
		patch{
			path:   serverPath,
			skipIf: registerFn + "(srv, " + name + "Service)",
			anchors: []string{
				"// register:route ── 新模块路由在此行后注册(make register 工具锚点,勿删)",
			},
			lines: []string{routeLine},
		},
	); err != nil {
		fmt.Fprintln(os.Stderr, "register:", err)
		os.Exit(1)
	}

	fmt.Printf("已登记 %s(%s):\n"+
		"  %s (repo|client / service / rest|grpc-arg)\n"+
		"  %s (param / route)\n"+
		"下一步: 实现 data/%s_repo.go 与 internal/service/%s_service.go,然后 go build。\n",
		typ, *svc, wiringPath, serverPath, lowerFirst(typ), lowerFirst(typ))
}

// apply 依次应用各 patch;同文件多个 patch 顺序生效,幂等由 skipIf 保证。
func apply(patches ...patch) error {
	for _, p := range patches {
		if err := applyOne(p); err != nil {
			return err
		}
	}
	return nil
}

func applyOne(p patch) error {
	raw, err := os.ReadFile(p.path)
	if err != nil {
		return err
	}
	content := string(raw)
	if strings.Contains(content, p.skipIf) {
		return nil // 已登记,幂等跳过
	}

	eol := "\n"
	if strings.Contains(content, "\r\n") {
		eol = "\r\n"
	}
	lines := strings.Split(content, eol)

	for _, anchor := range p.anchors {
		at := indexOfAnchor(lines, anchor)
		if at < 0 {
			continue
		}
		out := make([]string, 0, len(lines)+len(p.lines))
		out = append(out, lines[:at+1]...)
		out = append(out, p.lines...)
		out = append(out, lines[at+1:]...)
		return os.WriteFile(p.path, []byte(strings.Join(out, eol)), 0o644)
	}
	return fmt.Errorf("%s: 未找到注册锚点 %q,请检查锚点注释是否被移动或删除", p.path, p.anchors[0])
}

func indexOfAnchor(lines []string, anchor string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == anchor {
			return i
		}
	}
	return -1
}

// pascal 把 product / dict_entry / dictEntry 统一为 PascalCase 的实体名。
func pascal(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	var b strings.Builder
	for _, p := range parts {
		r := []rune(p)
		b.WriteRune(unicode.ToUpper(r[0]))
		b.WriteString(string(r[1:]))
	}
	out := b.String()
	if out == "" {
		return s
	}
	return out
}

func lowerFirst(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
