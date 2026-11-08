# rpc 包使用文档

本文档说明 `rpc` 包中 REST/gRPC 服务相关的构建方式与内置中间件（包含白名单、登录保护等）。目标读者为项目开发者/运维人员，文档偏实践、示例丰富。

---

## 概要

`rpc` 包提供对 HTTP (kratos rest) 与 gRPC 服务的封装、常用中间件的注册点和若干通用中间件实现（例如：白名单、登录保护、限流等）。主要出口函数：

- `CreateRestServer(cfg *conf.Bootstrap, mds ...middleware.Middleware) (*kratosRest.Server, error)`
  - 创建并配置 REST Server，返回 server 与可能的构建错误。
- `NewRestWhiteListMatcher() selector.MatchFunc`
  - 返回默认白名单 matcher（供 selector 使用），白名单实现位于 `rpc/whitelist.go`。

本包的设计原则：不在库内 `panic`，遇到配置/初始化错误应返回 error；中间件可配置并支持运行时白名单跳过。

---

## 快速开始（REST）

示例：在 `main` 中创建 REST server 并优雅关闭

```go
cfg := &conf.Bootstrap{ /* 从配置中心或文件加载 */ }
 srv, err := rpc.CreateRestServer(cfg)
 if err != nil {
     log.Fatalf("create rest server failed: %v", err)
 }
 // 启动 srv（kratos 调用或自行 Listen/Serve）
```

注意：`CreateRestServer` 已返回可能的初始化错误（例如 TLS 证书解析失败），调用方应处理该错误并决定是否退出或降级。

---

## initRestConfig 行为要点

`CreateRestServer` 内部通过 `initRestConfig` 构建 `kratosRest.ServerOption` 列表：

- CORS：使用 `github.com/gorilla/handlers` 的 CORS 中间件，来源自 `cfg.Server.Rest.Cors`。
- Middlewares：基于配置开启 recovery、tracing、validate、限流器、metadata 注入等；部分中间件会被 `selector.Server(..., NewRestWhiteListMatcher())` 包装，以便白名单跳过（例如 tracing、validate、ratelimit、metadata）。
- TLS：如果配置了 TLS，将调用 `loadServerTlsConfig`；若加载失败会返回 error。

若你想调整中间件顺序或注入自定义中间件，可以通过传入 `mds ...middleware.Middleware` 参数追加到中间件链尾部。

---

## 内置中间件（推荐用法）

下面列出项目中常用并已或建议内置的中间件，按常见部署场景排序并给出简短使用建议：

- Recovery：捕获 panic 并生成统一错误码，应位于最外层。默认已支持。
- RequestID：注入请求 ID（X-Request-Id），方便日志/trace 关联（建议在最外层）。
- Access Log（结构化日志）：记录请求、响应状态、耗时、trace id 等。建议与 RequestID 联用。
- Metrics（Prometheus）：记录请求数、失败数、latency 等；暴露 `/metrics`（建议受限访问）。
- Tracing（OpenTelemetry）：可选，支持 exporter 配置，并可被白名单跳过。
- Authentication / Authorization：JWT/ApiKey 等鉴权与权限检查（应在 validation 之前）。
- Rate Limit：请求速率限制（支持按 IP、按 key）；登录保护可基于此或用专用逻辑。已支持 selector 包装以跳过白名单。
- Validation（Proto/JSON）：请求语义校验。
- WhiteList（见下节）：供 selector 使用以跳过特定 operation 的中间件。
- LoginProtect：专用的登录保护中间件（按 IP/按账号失败计数与锁定）。
- Security Headers / HSTS：统一设置安全相关响应头。

中间件的具体默认实现/位置见代码 `rpc/rest.go` 中 `initRestConfig` 的组装逻辑。

---

## 白名单（Whitelist）

项目提供并发安全的白名单实现（`rpc/whitelist.go`），语义为“跳过某些中间件的判定集合”。核心点：

- 存储：`map[string]struct{}` + `sync.RWMutex`；写时 normalize（去除前导 `/`）以兼容 gRPC 完整方法名与 HTTP 路径。默认有 `DefaultWhiteList` 单例。
- API：
  - `AddWhiteList(ops ...string)`：追加项
  - `SetWhiteList(ops []string)`：替换列表
  - `ClearWhiteList()`：清空
  - `NewWhiteListMatcher()`：返回 `selector.MatchFunc`，在匹配时返回 `false` 表示“跳过中间件”。
- 匹配行为：默认 Exact 模式（先尝试完整名，如 `package.Service/Method`，再尝试 method-only `Method`），也支持 Prefix 模式（通过 `NewWhiteList( Prefix, ... )` 创建）。

示例：跳过验证中间件与限流

```go
ms = append(ms, selector.Server(validate.ProtoValidate(), rpc.NewRestWhiteListMatcher()))
ms = append(ms, selector.Server(ratelimitMiddleware, rpc.NewRestWhiteListMatcher()))
```

白名单适合用于健康检查、公开接口或内部监控回路等场景。

---

## 配置建议（示例片段）

下面给出简化的配置建议（YAML 片段），供参考：

```yaml
server:
  rest:
    addr: ":8080"
    cors:
      origins: ["https://example.com"]
      methods: ["GET","POST","OPTIONS"]
      headers: ["Content-Type","Authorization"]
    middleware:
      enableRecovery: true
      enableTracing: true
      enableValidate: true
      enableMetadata: true
      limiter:
        name: "bbr"
    tls: # 可选
      cert_file: "/path/to/cert.pem"
      key_file: "/path/to/key.pem"
```

注意：字段名与实际 proto 定义或 conf 结构应以仓库中的 `api/gen/go/conf/v1` 为准，此处为示意。

---

## 单元测试与验证

- 白名单：已添加并发与规范化测试（见 `rpc/whitelist_test.go`），应使用 `go test -race` 检查并发安全。

运行示例（在仓库根目录）：

```shell
# 仅 rpc 包测试（包含 race 检测）
go test ./rpc -race -v

# 全量测试（可能较慢）
go test ./... -race
```

---

## 其他建议与最佳实践

- pprof 与 metrics endpoint 应该限制访问（仅内网或需认证），避免在生产环境直接暴露到公网。可将 pprof 注册到单独的 admin server。
- 在高并发场景下，使用 Redis 作为限流/登录保护后端可获得更好的分布式一致性。
- 对于白名单与登录保护的运行时变更，建议提供简易 admin API（受 auth 保护）并记录变更审计日志。

---

## 联系与贡献

如需扩展中间件（例如接入 CASBIN、外部限流器、或自定义登录策略），欢迎提 PR 或在仓库中创建 issue 讨论实现细节。

---

感谢阅读本使用文档，祝你使用愉快！