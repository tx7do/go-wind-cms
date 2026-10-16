# AGENTS.md — go-wind-cms Monorepo 开发指南（AI Agent 入口）

本文件是 monorepo 总入口：定位布局、指路各端文档、声明全仓铁律。**深入开发前先读本文件与对应端的相关文档**。

## 仓库布局

GoWind Content Hub（HCH，无头内容中枢）：一个 Go 后端 + 一个管理后台前端 + 四套 C 端前端的 Headless CMS。

```
backend/                        Go 1.26 + Kratos + Ent（DI 手写装配 wiring.go，已弃用 Wire）
├── api/
│   ├── protos/                 .proto 源文件（admin/app/core 各面 + content/site/identity 等域）
│   └── gen/go/                 生成代码（HTTP/gRPC/校验），不手改
├── app/
│   ├── admin/service/          管理后台服务（REST :6600，SSE :6601 /events）
│   ├── app/service/            C 端应用服务（REST :6700，SSE :6701）
│   └── core/service/           核心业务与数据层（Ent schema、仓储、业务逻辑，被 admin/app 引用）
├── pkg/                        公共包（鉴权/中间件/OSS/viewer/OpenSearch/...）
└── scripts/                    部署与依赖脚本（Docker/PM2/环境安装）
frontend/
├── admin/                      管理后台：React 19 + antd 6 + ProComponents + TanStack Query + TipTap（dev :15999）
└── app/
    ├── react/                  C 端：Next.js + shadcn/ui（dev :15001）
    ├── vue/                    C 端：Nuxt + shadcn-vue（dev :15011）
    ├── taro/                   C 端：Taro 微信小程序/H5（dev :15021）
    └── flutter_app/            C 端：Flutter + BLoC
docs/                           专题文档（设计语言、内容建模设计等）
```

管理后台以 **React 版 admin 为唯一基线**（2026-09-17 由 Vben 切换为 React，2026-10-07 目录定名 `frontend/admin`，勿再引入 vue-admin 代码）。C 端四套并列、面向不同技术栈的采用者，一次部署取一套。

## 门禁（提交前必须全绿）

| 范围 | 命令（在各自目录下） |
|---|---|
| backend | `make lint` + `go test ./...` |
| frontend/admin | `npm run typecheck` + `npm run lint` |
| C 端各应用（改动落在哪端跑哪端） | 各自的 lint / typecheck（见其 package.json） |

**门禁出现任何新报错，一律当作自己引入的 bug 修复**，不存在"可忽略的既有错误"。改完代码先跑门禁再声称完成。

> 别照抄 `make vet`：那条是裸 `go vet`，而 `backend/` 根目录没有 `.go` 文件，跑完只吐一句 `no Go files`——什么都没检查，却不像失败。要真的过一遍，写 `go vet ./...`。

## 全仓铁律

1. **不吞错**：任何 catch 至少二选一——`console.error/warn` 带出**原始错误对象**，或重新抛出。用户可见的通知 ≠ 日志（只有翻译文案）。合法裸 catch 仅限纯本地 best-effort 兜底且注释写明原因。历史教训：认证链路静默吞错曾让 bug 排查耗时数日。
2. **生成产物不手改**：`api/gen/go/`、各端 `src/api/generated/`、`assets/openapi.yaml`、ent 生成代码都是产物。改 `.proto` 或 `ent/schema` 后必须重新生成并**连同产物一并提交**（命令见 [CONTRIBUTING.md](./CONTRIBUTING.md)）。
3. **CRUD 写请求体必须包 `{ data: {...} }`**（请求消息的 `data` 字段），更新配合 `update_mask` 做部分更新。列表搜索条件一律 contains（模糊匹配），ID 类字段不进模糊搜索。
4. **新端点必须走完整生成链**：`make openapi` 的产物 `openapi.yaml` 经 `assets.go` **打进二进制**，管理页「接口同步」读的是进程内嵌的这份——所以新 proto 要 `make api && make openapi` 后**重启进程**再触发同步，否则同步"成功"而新端点依旧不在接口管理表里，路由看起来正常但权限链路缺表项。详见 [CONTRIBUTING.md](./CONTRIBUTING.md)。
5. **C 端公开页遇 401 静默处理**：清 token + `loginExpired`，不得强制跳登录页；公开资源本就允许游客访问，弹登录是体验事故。

## 后端任务：gow 优先，make 兜底

后端的运行与代码生成统一走 [gow CLI](https://github.com/tx7do/go-wind-toolkit/tree/main/gowind)（安装：`go install github.com/tx7do/go-wind-toolkit/gowind/cmd/gow@latest`），在 `backend/` 下执行。Windows 无原生 make、嵌套 Makefile 需要跨目录 cd，而 gow 自动发现 `app/*/service`，行为一致：

| 任务 | gow 命令 |
|---|---|
| 运行服务 | `gow run admin`（C 端服务：`gow run app`） |
| Ent 生成 | `gow ent`（全部服务）/ `gow ent admin` |
| Proto / API Go 代码 | `gow api` |
| 从数据库表生成 CRUD | `gow generate`（DSN 驱动；`--proto-only` 仅出 proto） |

gow 未覆盖的任务才退回 Makefile（`backend/` 根目录执行）：`make ts`（五份 TS/OpenAPI 产物一条命令）、`make openapi`、`make gen`（ent + api + openapi）、`make register`、`make test/lint`。

- 项目已弃用 Wire（DI 为手写 `cmd/server/wiring.go`），不要运行 wire 相关命令，也不要在新代码里引入 wire 依赖。
- 新 CRUD 模块接入装配：`make register ENTITY=xxx [SVC=core] [PKG=contentV1]`（实为 `go run ./tools/register`，靠锚点行注入 `wiring.go` 与 `rest_server.go`，幂等可重跑；依赖更多的模块请手工补装配）。

## 本地验证要点

- **依赖服务**：`backend/scripts/docker/libs_only.ps1`（Windows）或 `libs_only.sh` 起 Postgres / Redis / MinIO 等；启动方式详见 `backend/README.md` 与 `backend/scripts/WORKFLOWS_AND_BEST_PRACTICES.md`。
- **端口**：admin REST `:6600`（SSE `:6601`，协同 ws `:6602`）、C 端 app REST `:6700`（SSE `:6701`，GraphQL `:6702`，`GRAPHQL_ADDR=off` 停用）；前端 dev 见布局表。admin 后端的 CORS 默认只放行 `localhost:15999`（`app/admin/service/configs/server.yaml`），换端口联调要同步改。
- **登录账号**：全新环境由服务启动时自动播种——管理员 `admin / admin`、普通用户默认口令 `12345678`（`pkg/constants/default_data.go`），**不要**找 SQL 种子脚本，`backend/sql/` 下只有演示数据。登录带图形验证码（Redis 存储）。
- **refresh-token 勿手动调**：刷新接口会轮换令牌对，手动调用会使内存中的旧 access token 立即失效，表现为"假 401"，别误判成鉴权 bug。
- **admin 前端坑**：go-crud List 的 `total` 是字符串（uint64 JSON），展示/比较前须 `Number()`；antd Tabs 走 `items` prop 之外的旧写法会静默渲染空壳。
- **ListTable 必传 `request`**：不传时内部 wrappedRequest 拿到 undefined 结果即判定失败态，页面恒定渲染「加载失败/重试」且无任何日志。新列表页一律用 redirects 页模式（`loadItems` + `request` prop + 受控 dataSource，增删改后 `actionRef.reload()` 重跑）。
- **新 List 接口禁传空 `sorting=[]`**：生成客户端会序列化成空 `sorting=` 参数，kratos 解析 `pagination.Sorting` 报 400（CODEC）。排序走 `orderBy`（JSON 字符串参数）。
- **AI 内容生成助手**已落地（`api/protos/ai/`）；本地调试可配任意 OpenAI 兼容 mock 替代真实模型，生产在管理页配真实 provider。
- **Node 版本**：admin dev（vite 8 / rolldown）要求 Node ≥ 20.19，否则报 `styleText` 不存在——系统默认 20.10 太老，须前置 nvm 的 20.20 再 `npm run dev`。
- **本地 exe 部署模式**（backend/build/{core,admin,app}.exe --conf ...）有守护进程在进程死后数秒内用同路径 exe 重启：换代码必须按「kill → build → start」顺序，或构建后再 kill（守护会用旧镜像拉起并占住端口，新进程静默退出）。服务日志若重定向到 build/*.log 可直接读。
- **`go test ./...` 偶发** `fork/exec %TEMP%\go-build...\x.test.exe: Access is denied.`（Windows 对刚链接好的测试二进制执行被拦，疑似杀软实时扫描）：属环境问题、**不是代码失败**。复验办法是绕开 Temp——`go test -c -o <工作区内路径>/x.test.exe ./pkg/x` 后直接跑该 exe；判成"测试挂了"之前先这样确认一次。

- **新增环境变量（2026-10 批次）**：`SITE_BASE_URL`+`SEO_POST_PATH`/`SEO_PAGE_PATH`（app 服务，sitemap URL 生成，未配置返回空 sitemap）；`SMTP_HOST/PORT/USERNAME/PASSWORD/FROM/TLS`（邮件出站回退配置，email_providers 表可视化配置优先）；`MCP_ADDR`/`MCP_TOKEN` 已有；`GRAPHQL_ADDR`（app 服务 GraphQL 监听地址，默认 `:6702`，off/disabled 停用）；`COLLAB_ADDR`（admin 协同 ws，默认 `:6602`，off/disabled 停用）+ `COLLAB_ALLOWED_ORIGINS`（允许来源，默认 `http://localhost:15999`）+ `VITE_COLLAB_WS_URL`（前端协同 ws 地址，默认 `ws://localhost:6602`）。协同引擎依赖固定 fork：`replace github.com/reearth/ygo => github.com/tx7do/ygo v1.51.11`（上游更新按需同步）。协同 token 走 query 参数（与 SSE 同款先例），会进反代访问日志——生产日志留存策略注意脱敏。预览令牌签名复用 `ENCRYPTION_KEY`。
- **凭据加密（ENCRYPTION_KEY）**：core/admin 启动读取环境变量 `ENCRYPTION_KEY`
  （两服务必须同值），社交账号凭据与 AI provider key 据此 AES-GCM 加密落库；
  未设置时降级明文入库并打告警——生产必须配置，换 key 会使既有密文不可解。
- **接口同步（已升级为对账式 upsert，非空表安全）**：core `ApiRepo.SyncApis`
  以 (path, method) 为业务键对账——已存在则原位更新（**不换 ID**，保护角色/
  权限绑定）、缺失则 maxID 续号创建并自动 `setval` 对齐序列；admin 侧错误
  已传播到管理页。新增端点后管理页点一次「接口同步」即可，无需手工补录。
- **asynq 任务类型命名**：ServeMux 精确匹配失败后按**前缀**兜底路由，
  新任务类型不得与其他已注册类型互为前缀（"post.publish.sweep" 曾被前缀
  路由到 "post.publish" 的 handler 导致 payload 解码失败被吞）。

## 文档索引

- 根 README（三语）：[README.md](./README.md) / [README.en-US.md](./README.en-US.md) / [README.ja-JP.md](./README.ja-JP.md)
- 后端：`backend/README.md`（Swagger、proto/ent/OpenAPI 生成、构建、Docker 部署）
- 脚本与工作流：`backend/scripts/WORKFLOWS_AND_BEST_PRACTICES.md`
- 设计语言：`docs/design-language.md`（品牌蓝等视觉唯一权威值，改颜色先改这里再同步各端）
- 内容建模：`docs/content-modeling-design.md`
