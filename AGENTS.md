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
- **换代码必须重建对应 exe 再重启**：只重启不重建=旧二进制静默服务（ent 新表不建/新 RPC Unimplemented/新字段运行时 nil）。判定：`ls --time-style=full-iso build/x.exe` 对时间线，或 `go version -m build/x.exe`；strings 工具对 Go Windows PE 全零命中，不可作二进制新旧证据。core 的 etcd 租约长跑会静默失效（admin 报 Zero endpoint found）——重启 core 即恢复。
- **ygo 协同引擎坑**（内嵌于 admin :6602）：Transact 内禁 doc 级读方法（GetXmlFragment 等，会死锁且无日志——引用事务外取好）；同一容器的子节点须逐个独立事务前向插入（同事务 prelim 混排 YText/YXmlElement 重放乱序）；YText 插入无格式文本也须显式空 Attributes（nil/省略会继承相邻格式）；listItem 等内容模型要求块级子节点（裸文本会被 y-tiptap 整树拒绝→编辑器空白）。种子化单测见 collab 包 seed_test.go。
- **AI 改 Go 代码的 python replace 静默 no-op**：赋值行对齐空格不匹配即不生效且无报错——replace 后必须 grep 验证或带 assert；曾致构造函数漏赋值（编译过、wiring 对、运行时字段全 nil）。
- **实验转化归因口径**：表单提交 → gwc_vid → (experiment_id, visitor_hash) 定位参与记录标 converted；不做 viewer 租户过滤（hash 服务端派生，与参与记录口径一致）。纯转化型实验（变体无 sections 覆盖）也须先记参与再判断是否替换内容。
- **本地 exe 部署模式**（backend/build/{core,admin,app}.exe --conf ...）有守护进程在进程死后数秒内用同路径 exe 重启：换代码必须按「kill → build → start」顺序，或构建后再 kill（守护会用旧镜像拉起并占住端口，新进程静默退出）。服务日志若重定向到 build/*.log 可直接读。
- **`go test ./...` 偶发** `fork/exec %TEMP%\go-build...\x.test.exe: Access is denied.`（Windows 对刚链接好的测试二进制执行被拦，疑似杀软实时扫描）：属环境问题、**不是代码失败**。复验办法是绕开 Temp——`go test -c -o <工作区内路径>/x.test.exe ./pkg/x` 后直接跑该 exe；判成"测试挂了"之前先这样确认一次。

- **新增环境变量（2026-10 批次）**：`SITE_BASE_URL`+`SEO_POST_PATH`/`SEO_PAGE_PATH`（app 服务，sitemap URL 生成，未配置返回空 sitemap）；`SMTP_HOST/PORT/USERNAME/PASSWORD/FROM/TLS`（邮件出站回退配置，email_providers 表可视化配置优先）；`MCP_TOKEN`（MCP 鉴权 Token，必填；凭据不入配置文件）；`GRAPHQL_ADDR`（app 服务 GraphQL 监听地址，默认 `:6702`，off/disabled 停用）；`COLLAB_ADDR`（admin 协同 ws，默认 `:6602`，off/disabled 停用；yaml `server.collab` 段优先、env 兜底，见下条）+ `COLLAB_ALLOWED_ORIGINS`（允许来源，默认 `http://localhost:15999`，同被 `server.collab.allowed_origins` 覆盖）+ `VITE_COLLAB_WS_URL`（前端协同 ws 地址，默认 `ws://localhost:6602`）。协同引擎依赖固定 fork：`replace github.com/reearth/ygo => github.com/tx7do/ygo v1.51.11`（上游更新按需同步）。协同 token 走 query 参数（与 SSE 同款先例），会进反代访问日志——生产日志留存策略注意脱敏。预览令牌签名复用 `ENCRYPTION_KEY`。
- **MCP 服务器配置化（admin 服务）**：监听地址与服务器参数走 config `server.mcp` 段（`conf.Server_MCP`，经 kratos-bootstrap/transport/mcp 适配器映射，字段含 address/server_name/server_version/server_type/instructions/能力开关），示例见 `app/admin/service/configs/server.yaml` 注释段；未配置回退 `MCP_ADDR` 环境变量（等价 HTTP+Recovery+默认名版本），off/disabled 停用，两处都没有则不启动。依赖 kratos-bootstrap `api v0.0.45` + `transport/mcp v0.0.5`（2026-10-16 发布，注意上游 tag.bat 有 `--force` 重打旧版本号惯例——同版本号重打会被 Go proxy 缓存旧内容，跟进时认准真正的新版本号）。
- **协同 ws 配置化（admin 服务）**：监听地址与 Origin 白名单走 config `server.collab` 段（addr/allowed_origins，示例见 server.yaml 注释段），未配置回退 `COLLAB_ADDR`/`COLLAB_ALLOWED_ORIGINS`。上游 `conf.Server` 无 collab 字段，经项目自定义 proto 扩展 `api/protos/conf/v1/bootstrap_extension.proto`（`BootstrapExtension` 壳消息，protojson 整树 Scan 同构映射）接入，且**必须在 `LoadBootstrapConfig` 前注册**（`main.go` runApp 里 `server.RegisterCollabConfigExt`，晚了静默失效）——上游 conf 增设 Collab 段后该扩展可整体删除，yaml 形状不变。
- **OIDC 登录（core 服务，授权码模式）**：core 读取 `OIDC_ISSUER`/`OIDC_CLIENT_ID`/`OIDC_CLIENT_SECRET`/`OIDC_AUTHORIZE_URL`/`OIDC_TOKEN_URL`/`OIDC_USERINFO_URL`（六者齐备才启用；`OIDC_SCOPES` 默认 `openid profile email`，`OIDC_AUTO_PROVISION` 默认 true，`OIDC_DEFAULT_ROLE_CODE` 可选——首登建号绑定的平台级角色 code，未设则建号无角色、登录被授权链拒）。前端入口：登录页「企业单点登录」→ `POST /admin/v1/oidc/authorize-url`（免鉴权白名单）取授权地址整页跳 IdP → 回调 `/auth/login?code=&state=` 前端以 `grant_type=authorization_code` 调 Login。state 由 core 生成落 Redis（10 分钟、一次性消费，Redis 缺位 fail-closed）；用户映射按 (issuer, sub) 查 ENTERPRISE_SSO 凭证（identifier=`issuer|sub`）。C 端面：app 服务 `POST /app/v1/oidc/authorize-url`（免鉴权）+ react/vue 登录页 SSO 按钮与 code/state 回调消费；taro/flutter 的详情页自定义字段渲染已落地（taro RichText/公开媒体端点换 URL、flutter_widget_from_html）；SSO 仍暂缓（小程序宿主限制外部 IdP 跳转、flutter 需 WebView/Custom Tab），单独立项。E2E：`backend/scripts/e2e/test_oidc_login.py`（内置 mock IdP :6900，core/admin/app 须带上述 env 启动；19 断言含 C 端面）。注意 `gow api` 会清掉 `api/gen/go/graphql/` 的 gqlgen 产物——生成后须在 `backend/api/graphql` 补跑 `go run github.com/99designs/gqlgen generate`（resolver.go 是手写实现，gqlgen 在包缺失时重写会丢实现，须 git checkout 恢复）。
- **内容建模前台读侧（2026-10 批次）**：core `ContentModelService.GetModelByPost`（文章→分类→模型 服务端解析，调用方无法指定）+ app 免鉴权 `GET /app/v1/posts/{id}/content-model` 与 `GET /app/v1/media-assets/{id}`（customFields IMAGE 值为 `media_asset:<id>` 引用，渲染链经此换公开 URL；core Get 会被软删资产 404）。前端 react/vue 文章详情 `CustomFieldsSection`：按模型字段定义渲染 post.customFields（值统一字符串；RICHTEXT 注入、IMAGE/FILE 形似 URL 时图/链接、`media_asset:<id>` 经公开媒体端点换 URL、其余原样），无模型/加载失败静默不渲染。E2E：test_media_references.py 第 11 段（27 断言）。
- **RELATION 字段前台深度解析（2026-10）**：app 免鉴权 `GET /app/v1/content-references/resolve?refs=post:12&refs=page:7&languageCode=zh-CN`（proto 面走 content `ContentModelService.ResolveContentReferences`，app 面挂在 PostService 下，走完整生成链）。core 只解析 **PUBLISHED** 的 post/page（草稿/软删/非法形态/`category:`/`tag:` 静默缺失，宁缺毋滥；allow_cross_tenant v1 不放行跨租户）；≤20 refs、去重；标题/slug 按 languageCode 命中、独立回退首个非空；响应含 urlPath（`/post/{slug}`）但**前端不直接用它**——react 用 `/{locale}/post/detail?slug=`、vue post 走 id 路由 page 走 slug、taro 仅 post 可跳（无页面详情路由）、flutter post/page 均按 id。四端 CustomFieldsSection 把 RELATION 渲染成链接卡片（解析缺失回退原样字符串）；react 的 resolve useQuery 必须在早退 return 之前调用（rules-of-hooks）。**gqlgen 生成坑**：`gow api` 会清空 `api/gen/go/graphql/`，恢复姿势 = `git checkout -- api/gen/go/graphql/` 后 `go run -mod=mod github.com/99designs/gqlgen generate`（gqlgen 是 go.mod 直接依赖但 CLI 传递依赖不在 go.sum；`go run` 裸跑会报 missing go.sum）——**它会重写 `app/app/service/internal/graphql/resolver.go` 丢手写实现**，生成完必须再 `git checkout -- app/app/service/internal/graphql/resolver.go`。E2E：`backend/scripts/e2e/test_content_references.py`（15 断言；发布必须走 submit+approve）。
- **即时图片变换（app 服务，2026-10）**：免鉴权 `GET /app/v1/media-assets/{id}/content?w=&h=&q=&fit=`——非 proto 手工挂载路由（操作名 `MediaAssetService_GetContent`，同 sitemap/robots 先例，规避生成链）。实现全在 app 门面（`media_asset_service.go` + `image_transform.go`）：可见性复用 core Get（租户隔离/软删 404）后才经 MinIO 读字节；**只缩小不放大**（目标≥源直出原字节），源限 jpeg/png/webp（gif 丢帧、svg XSS → 415，走元数据 url 兜底），输出统一 JPEG（Go 纯实现无 webp 编码，x/image/webp blank import 供解码）；参数钳位（边≤4096、q∈1..100 缺省 80、fit=contain|cover 且 cover 需双边）+ `image.DecodeConfig` 先行防解码炸弹（像素总量 ≤4000 万）；变换结果按参数键落 MinIO 派生对象（`{对象名}-w{w}h{h}q{q}-{fit}.jpg`，与预生成三档 `-thumbnail/-medium/-large` 后缀不撞名），命中免解码直出，写缓存 best-effort 失败不阻塞响应；`Cache-Control: immutable` + ETag/304。E2E：`backend/scripts/e2e/test_image_transform.py`（20 断言，真实上传字节落 MinIO；**测试图片内容须逐轮唯一**——固定内容会撞 FindByHash 查重 409）。C 端四端尚未接入该端点（仍用元数据 url），接入属后续切片。
- **自动化 Flow（core 服务）**：`flow_rules` 表（trigger_event/condition_status/actions JSON/enabled/match_count）+ `/admin/v1/flow-rules` CRUD（/system/flow-rule 页）。事件族同 webhook（post.*/page.*），触发点与 webhook 订阅并列（post/page 增删改 + 审核通过 + 定时转正）。worker `flow.execute`（asynq，MaxRetry 3）：条件 v1 仅状态等值（condition_status 空=无条件），动作 v1：webhook（HMAC 签名同订阅，envelope 含 rule_id）与站内通知（title/content 支持 {event}/{entity_id}/{rule_id} 占位符；经 InternalMessageService.SendSystemNotification——系统身份 created_by=0、逐收件人落收件箱、不发邮件，跨租户不拦截）。**发布语义坑**：直接 PUT status=PUBLISHED 会被审批工作流吞掉（status 不落库、事件走 updated），发布事件只能由审核 approve（或定时转正）产生——E2E 断言必须走 submit+approve。**gRPC 代理坑**：从带 http_proxy 环境变量的 shell 启动 admin 会让 grpc-go 经 7890 代理拨 core（503 connection error via 127.0.0.1:7890），启动服务前须 env -u 清代理。E2E：`backend/scripts/e2e/test_flow_rules.py`（mock 接收端 :6901，11 断言）。**asynq 队列隔离（2026-10-17 修正）**：core 的 asynq 已从 redis db1 迁到 **db2 专库**（server.yaml `server.asynq.uri`）——db1 的 `asynq:{default}` 被本机多个项目共用，外来 asynq worker 会抢走本项目的任务（`handler not found` → 重试 3 次归档，任务静默丢失；此前的"~20s 队列滞后"玄学多半是这场抢夺）。E2E 排查 asynq 队列时认准 `-n 2`。同期修复存量 bug：`post.scheduled.sweep`/`social.distribution.sweep` 两个周期兜底任务**从未注册订阅者**（handler 在、注册行丢），每分钟入队一条死任务、retry 积压 7400+ 条，兜底扫描从未生效——订阅已在 asynq_server.go 补上，db1 历史积压已清。**Flow 表达力扩展（2026-10）**：**Flow 表达力扩展（2026-10）**：动作新增 `delay`（`{"type":"delay","seconds":N}`，1..604800 钳位）——动作链遇 delay 把剩余链交给 `flow.resume` DELAY 任务（ProcessIn）续跑，**任务类型不得叫 flow.execute.***（前缀兜底路由会被 flow.execute 吞掉）；续跑跳过条件判定、命中计数不重复。新增定时触发：`triggerCron`（5 段 cron，robfig 解析校验）非空即纯定时规则——被事件扇出排除（`TriggerCronIsNil`），由每分钟 `flow.schedule.tick` 周期任务按其租户直派（Event 固定 `schedule`、EntityID=0、直派分支在 ExecuteWorker 里必须走 Get 规则而非事件扇出查询）；tick 丢失的分钟不补跑；多实例部署会重复触发（按单调度实例部署）。admin /system/flow-rule 页支持触发方式切换（事件/定时）+ delay 动作编辑。E2E：test_flow_rules.py 第 7/8 段。
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
