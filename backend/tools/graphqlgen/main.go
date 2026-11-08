// graphqlgen —— gqlgen 生成恢复工具。
//
// 为什么存在：api/gen/go/graphql/{generated,models_gen}.go 是提交入库的生成产物，
// 但 `gow api` 会整目录清掉它们；且 app/app/service/internal/graphql/resolver.go
// 是手写的单文件实现（含 mapPost/optStrVal 等 helper 与注入字段），与 gqlgen
// single-file 布局的"合并校验"不兼容——resolver.go 在场时裸跑 gqlgen generate
// 必然 validation failed（stub 的空 Resolver 结构对不上手写方法体/包级 helper），
// resolver.go 缺席时生成成功但会产出丢实现的 stub。
// 因此恢复/再生成必须走"挪开 resolver.go → gqlgen generate → 原样放回"的舞步，
// 本工具把它自动化并附编译验证。`gow api` 之后请务必运行一次：
//
//	go run ./tools/graphqlgen        # backend/ 下；等价 make graphql
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	gqlDir       = "api/graphql"
	genDir       = "api/gen/go/graphql"
	resolverPath = "app/app/service/internal/graphql/resolver.go"
	stashSuffix  = ".gqlgen-stash"
	// go.mod 的 tool 指令已记全 gqlgen CLI 传递依赖的 go.sum 条目，
	// 走 go tool 不弄脏 go.mod/go.sum（勿改回 go run -mod=mod …）。
	gqlgenRunArgs = "tool gqlgen generate"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "graphqlgen: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	resolver, err := abs(resolverPath)
	if err != nil {
		return err
	}
	stash := resolver + stashSuffix
	generated, err := abs(filepath.Join(genDir, "generated.go"))
	if err != nil {
		return err
	}

	// 上次运行中断会留下 stash——先自愈放回，避免旧舞步残留盖掉本次产物。
	if _, err := os.Stat(stash); err == nil {
		fmt.Println("graphqlgen: 发现上次中断遗留的 resolver 暂存，先恢复")
		if err := restoreResolver(resolver, stash); err != nil {
			return err
		}
	}

	// 舞步第一步：挪开手写 resolver.go（.stash 后缀不会被 Go 视为包内文件）。
	stashed := false
	if _, err := os.Stat(resolver); err == nil {
		if err := os.Rename(resolver, stash); err != nil {
			return fmt.Errorf("挪开 resolver.go 失败: %w", err)
		}
		stashed = true
		fmt.Println("graphqlgen: 已暂存手写 resolver.go（gqlgen single-file 合并校验与其不兼容）")
	} else {
		fmt.Println("graphqlgen: 警告——手写 resolver.go 不存在，本次生成产出的是 stub；")
		fmt.Println("  若仓库历史中有手写版，生成后请 git checkout -- " + resolverPath + " 恢复")
	}

	fmt.Printf("graphqlgen: go %s\n", gqlgenRunArgs)
	if out, err := runGo(gqlDir, gqlgenRunArgs); err != nil {
		// 失败也要把手写文件放回去，别把仓库留残局。
		if stashed {
			_ = restoreResolver(resolver, stash)
			fmt.Println("graphqlgen: 生成失败，已恢复手写 resolver.go")
		}
		return fmt.Errorf("gqlgen 生成失败: %v\n%s", err, out)
	}

	// 舞步第二步：放回手写实现，覆盖 gqlgen 产出的 stub。
	if stashed {
		if err := restoreResolver(resolver, stash); err != nil {
			return err
		}
		fmt.Println("graphqlgen: 已放回手写 resolver.go（覆盖 stub）")
	}

	if _, err := os.Stat(generated); err != nil {
		return fmt.Errorf("生成后未找到 %s，检查 gqlgen.yml exec.filename 配置", generated)
	}

	// 编译验证：解析器包 + 产物包一起过一遍。
	if out, err := runGo(".", "build ./"+genDir+" ./app/app/service/internal/graphql"); err != nil {
		return fmt.Errorf("编译验证失败: %v\n%s", err, out)
	}
	fmt.Println("graphqlgen: 生成完成，编译验证通过")
	return nil
}

// restoreResolver 用暂存覆盖 gqlgen 产出的 stub（Windows 下 rename 不能覆盖已存在目标，先删再移）。
func restoreResolver(resolver, stash string) error {
	if err := os.Remove(resolver); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("清理 gqlgen stub 失败: %w", err)
	}
	if err := os.Rename(stash, resolver); err != nil {
		return fmt.Errorf("恢复手写 resolver.go 失败: %w", err)
	}
	return nil
}

func abs(rel string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(wd, rel))
}

// runGo 在 dir 下执行 go <args>（args 空格切分，首段为子命令），回传合并输出。
func runGo(dir, args string) (string, error) {
	fields := strings.Fields(args)
	cmd := exec.Command("go", fields...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}
