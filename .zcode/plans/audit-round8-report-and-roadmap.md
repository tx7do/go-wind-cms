# go-wind-cms 第 8 轮综合审计报告与修复路图

- **审计日期**：2026-08-10
- **审计范围**：综合（安全遗留项复核 + 高危安全维度扫描 + 工程质量/供应链 + 三服务架构信任边界）
- **审计方法**：三路并行只读 agent 调查 + 主审计员对关键事实亲验（git/文件系统）
- **产出性质**：报告 + 分阶段路图，**本轮不改代码**，供评审后决定实施
- **前置**：延续 `.zcode/plans/plan-sess_05110011-...md`（第 7 轮修复计划）的"范围决策"延后清单

---

## 0. 关键事实纠正（亲验）

落盘前对若干前提做了实测，以下纠正前期描述：

| 前期描述 | 实测结论 | 验证方式 |
|---|---|---|
| `backend/server.exe`（~158MB）被提交进仓库 | **未入库**。被 `backend/.gitignore:3 *.exe` 忽略，`git ls-files --error-unmatch` 返回未跟踪 | `git check-ignore -v` / `git ls-files` |
| `backend/.env` 被 git 跟踪 | **未跟踪**。`git ls-files` 无 `.env` 条目 | `git ls-files \| grep .env` |
| 入库的大体积二进制是 server.exe | 实际入库的大文件是 `frontend/app/flutter_app/.tools/openapi-generator-cli.jar`（生成器 CLI） | `git ls-files \| grep -iE '\.(exe\|jar)$'` |
| `go.yml__` 是 CI 配置 | 文件存在但 `.yml__` 后缀不被 GitHub Actions 识别 → **不执行**；且内容陈旧（`./blog-backend`、Go 1.19.x、`setup-go@v3`/`codecov-action@v3` 已弃用） | `ls .github/workflows/` + 文件内容 |

下文发现项中，凡涉上述事实的，定级与处置均以实测为准。

---

## 1. 审计结论摘要

信任边界**完全压在 BFF**：core 无独立鉴权，无条件信任 BFF 传播的 `x-md-global-operator`（base64 protobuf，**无签名**）。当前唯一有效防线是 (1) 网络隔离 core、(2) BFF 的 JWT 校验 + ent.Server fail-closed。RBAC（noop）、审计日志（app 不落盘 / core 无中间件 / 落盘含密码明文）均未形成有效纵深。叠加 CI 禁用 + 安全路径零回归测试，第 3~7 轮修复的有效性**无法被持续保证**。

本轮发现统计：

| 严重度 | 数量 | 编号 |
|---|---|---|
| Critical | 3 | AUD-C1~C3 |
| High | 6 | AUD-H1~H6 |
| Medium | 6 | AUD-M1~M6 |
| Low | 7 | AUD-L1~L7 |

> 验证状态标注：`[亲验]` = 主审计员实测确认；`[审计调查]` = 由审计 agent 读码得出、本轮未独立复核。生产部署配置（环境变量、网络隔离）未动态验证，相关可利用性结论以"假设默认部署"为前提。

---

## 2. Critical 发现

### AUD-C1 — core gRPC 无鉴权 + OperatorMetadata 无签名校验，直连可越权全租户 `[审计调查]`
- **位置**：`backend/app/core/service/internal/server/grpc_server.go:29-34`（中间件链仅 logging + ent.Server，无 accessTokenChecker）；`backend/pkg/middleware/ent/ent.go:56-76`（fail-closed）；`backend/pkg/metadata/operator_metadata.go:20-21,34-73`（`x-md-global-operator` 为 base64 protobuf，`mdSignatureKey` 定义但从未校验）；`backend/pkg/metadata/errors.go:20`（`ErrSignatureMismatch` 定义但全仓无引用）；`backend/app/core/service/configs/server.yaml:2-3`（gRPC `addr: 0.0.0.0:0`，无 TLS/mTLS）
- **现状**：core 进程 `DefaultWhiteList` 为空（core 从不调 `AddWhiteList`），`ent.Server()` 对 md==nil 的非白名单请求 fail-closed 拒绝。**但**只要请求携带任意 `x-md-global-operator` 头，core 即反序列化为 `OperatorMetadata` 并构建 `UserViewer`，无来源校验、无签名、无 mTLS。
- **可利用路径**：能直连 core gRPC 端口者（经 etcd 服务发现即可拿到地址）伪造头 `{UserId:任意, TenantId:0, DataScope:ALL}` → `IsPlatformContext()`=true → `TenantPrivacy.EvalQuery`/`AllowIfAdmin` 均 skip 租户谓词 → 全租户越权读写；`EvalMutation` 对 platform context 不强制覆盖 tenant_id，写入亦可任意指定。
- **前提**：core 绑定 `0.0.0.0` 无 TLS，当前唯一防线是网络隔离（生产是否用 Sidecar/网络策略未核验）。
- **建议**：(1) 给 `x-md-global-operator` 增加 HMAC 签名（`mdSignatureKey` 已预留）并在 core 侧校验；(2) core gRPC 启用 mTLS / 限内部网络；(3) 绑定地址改内网，移除 0.0.0.0 公开。

### AUD-C2 — CI 构建/测试流水线禁用且失效，无任何自动化测试门禁 `[亲验]`
- **位置**：`.github/workflows/go.yml__`（1241B，双下划线不执行，内容陈旧）；`.github/workflows/codeql-analysis.yml`（717B，唯一活跃 CI，仅 Go SAST）
- **现状**：`go.yml__` 因后缀不被 Actions 识别 → 从不执行；即便启用也失效（指向旧项目 `blog-backend`、Go 1.19.x、弃用 actions）。原承担的 `go build ./...` / `go test -race -coverprofile` / API 生成 / Codecov **当前全部不运行**。CodeQL 仅 `language: ['go']`，前端 JS/TS/Vue/React/Taro/Flutter 与 Actions 自身不在 SAST 范围，且 SAST ≠ build/test/lint。无 Jenkinsfile / GitLab CI / 其他流水线。前端 `frontend/admin/package.json` 定义的 `lint`/`check:type`/`test:unit`(vitest)/`test:e2e`(playwright)/`build` **无任何 CI 调用**。
- **风险**：安全修复（中间件链、租户隔离、SSE 鉴权、mask 过滤）无自动化回归保护；任何回归、重构或依赖升级都可能静默破坏隔离性。PR 合并无门禁 = "信任提交"。
- **建议**：新建 `.github/workflows/ci.yml`，矩阵覆盖后端 `go vet && go build ./... && go test -race -coverprofile` 与前端 `pnpm install --frozen-lockfile && pnpm lint && pnpm check:type && pnpm test:unit`，配置为 PR 必过检查；删除 `go.yml__`。

### AUD-C3 — 安全关键路径回归测试覆盖接近零 `[审计调查]`
- **现状（逐项核实）**：
  - 鉴权中间件链 `backend/pkg/middleware/auth/*.go`（~458 行）→ 测试中**零命中**
  - 租户隔离中间件 `backend/pkg/middleware/ent/{ent,broker}.go`（~261 行）→ 测试中**零命中**（grep `viewer.|TenantPrivacy|IsSystemContext|SystemViewer|TenantID()` 均无）
  - 57 个 `*_repo.go` → 仅 2 个有 co-located 测试（`menu_repo_test.go`、`user_repo_test.go`），且为 protobuf FieldMask 工具测试，**不涉及 DB 查询作用域或跨租户拒绝**
  - updateMask 过滤 → 仅上述 2 测试出现，未断言越权字段剥离
  - SSE 鉴权 `sse_server.go` + `internal_message_service.HandleAuthorize` → **无专属测试**
- **现存 30 个 `*_test.go`** 主要覆盖 `pkg/crypto`、`pkg/lua`、`pkg/oss`、`pkg/jwt`(payload)、`pkg/content/{count,summary}`、`pkg/utils/converter` 等工具/旁路模块，未覆盖 auth/ent 中间件、租户谓词、mask 服务侧行为、SSE 鉴权。
- **风险**：与 AUD-C2 叠加 → 第 3~7 轮修复既无 CI 守护又无测试锁定，一次 `Update` 重构/ent 升级/中间件顺序调整都可能让租户隔离失效而无人知晓。
- **建议**：为 (a) auth 中间件 token 校验/过期/黑名单分支，(b) ent 中间件 viewer/TenantID 注入与白名单兜底，(c) 高风险 repo（user/tenant/post/media_asset/file/permission/role）跨租户读写拒绝，(d) SSE `HandleAuthorize` 有效/无效/黑名单 token 三分支，(e) updateMask 对越权字段（tenantId、role）剥离，编写表驱动测试 + sqlite/miniredis 集成测试。

---

## 3. High 发现

### AUD-H1 — RBAC 鉴权引擎为 noop，权限校验形同虚设 `[审计调查]`
- **位置**：`backend/app/admin/service/internal/data/data.go:100-102`、`backend/app/app/service/internal/data/data.go:92-94`（`NewAuthorizer()` 返回 `noop.State{}`）；`kratos-authz/engine/noop/noop.go:39-41`（`IsAuthorized` 恒 true）；`rest_server.go:70-77`(admin)/`:72-79`(app)（注入 noop）；`grpc_server.go:29-34`（core gRPC 无 authz 中间件）
- **现状**：admin/app BFF authorizer 是 `noop.State{}`，`IsAuthorized` 恒 true。core 虽有"真"`authorizer.Authorizer` 包装器（`backend/pkg/authorizer/authorizer.go`），但 core 无 `authz.yaml` → `Engine()` 为 nil，且 gRPC 链未挂 authz 中间件。`permission_service.go`/`policy_evaluation_log_service.go` 仅做 RBAC 数据 CRUD，不影响决策。`LoginPolicyService` 仅策略配置 CRUD，登录流程（`authentication_service.go:73-93`）从不调用它。
- **风险**：JWT 内任意角色即获全权。
- **建议**：启用 casbin/opa 引擎（`backend/pkg/authorizer/authorizer.go` 已实现策略生成），BFF `NewAuthorizer` 改用该包装器并配 `authz.yaml`；core gRPC 也挂 authz 中间件做纵深防御。

### AUD-H2 — 登录无限流/无锁定/captcha fail-open，易暴力破解 `[审计调查]`
- **位置**：`backend/app/core/service/internal/service/authentication_service.go:193-273`（`doGrantTypePassword` 无失败计数/锁定/IP 限流）；`backend/app/app/.../authentication_service.go:36-54`（app 登录无验证码、直接转发）；`backend/app/admin/.../authentication_service.go:56-80`（admin 验证码 fail-open，`captchaClient==nil` 时放行，注释明示）；三服务 `server.yaml` 均无 `Limiter` 配置；`initGrpcServerConfig`(`grpc.go:147-154`) 仅配置非 nil 时加 bbr，实际未启用
- **现状**：全链路无 rate limit 中间件；core 登录无失败计数/锁定；`LoginPolicyService` 数据存在但登录时不查询。admin captcha fail-open，app 登录完全无验证码。
- **可利用路径**：app 登录端点（白名单 `OperationAuthenticationServiceLogin`）+ core 无锁定 → 可对 `/app.service.v1.AuthenticationService/Login` 无限制撞库。
- **建议**：(1) 登录流程接入 `LoginPolicy`（锁定阈值/时长）；(2) 加 IP+账号维度 rate limit；(3) captcha 改 fail-closed；(4) app 登录也加验证码。

### AUD-H3 — 审计日志体系缺陷：app BFF 不落日志 + core gRPC 无审计中间件 `[审计调查]`
- **位置**：`backend/app/app/service/internal/server/rest_server.go:60-67`（`WithWriteApiLogFunc`/`WithWriteLoginLogFunc` 均 `return nil`）；`backend/app/admin/.../rest_server.go:52-65`（admin 会通过 gRPC 写 core，对比可见 app 是空实现）；`backend/app/core/.../grpc_server.go:29-34`（core gRPC 仅 logging+ent，无审计中间件）；`backend/pkg/middleware/logging/api_audit_log.go:84-87`（仅在 BFF HTTP 层捕获）
- **现状**：app BFF 调用审计中间件但丢弃所有日志；core gRPC 完全无审计中间件。即 app 站点登录、公开内容访问、任何经 app BFF 的操作均无审计；直连 core gRPC 的操作也无审计。
- **风险**：无法追溯 app 端攻击/越权；core 直连路径完全黑箱。
- **建议**：app BFF 接入真实审计写入；core gRPC 增加审计中间件覆盖写操作。
- **关联**：与 AUD-H4（落盘内容含密码）、AUD-L1（签名密钥/跳过逻辑）同属审计日志体系，路图中合并为"审计日志体系加固"工作流。

### AUD-H4 — API 审计日志记录完整请求体，含密码等敏感字段 `[审计调查]`
- **位置**：`backend/pkg/middleware/logging/api_audit_log.go:47,56`
- **现状**：`ApiAuditLogMiddleware.Handle` 对**所有**非登录请求 `io.ReadAll(htr.Request().Body)` 并将原始字符串写入 `apiAuditLog.RequestBody` 持久化。仅 `loginOperation` 被排除（`:38-39`）。
- **风险**：`/admin/v1/me/password`(ChangePassword)、`/admin/v1/users/{user_id}/password`(EditUserPassword) 等密码类端点不在排除列表 → 新/旧密码明文落 `api_audit_log.request_body`。任何能读该表的角色（只读 DBA、备份恢复、日志聚合）可直接读密码。同理波及含 token 的内部消息等端点。
- **建议**：对 password/secret/token 类端点加排除名单或字段级脱敏（正则替换 `password`/`old_password`/`new_password`/`token` 值为 `***`）后落库。

### AUD-H5 — 配置散布弱默认密钥且无 fail-fast（含 oss.yaml 硬编码 access_key） `[审计调查]`
- **位置**：`backend/app/core/service/configs/authenticator.yaml:4,7,11,14`（`${jwt_signing_key:dev_only_change_me_in_prod}` / `${aes_key:dev_only_change_me}`）；`data.yaml:4,9,18`（redis/db/es 密码同模式）；`backend/app/{admin,app,core}/service/configs/{server,client}.yaml`（JWT key 同默认）；`backend/app/{admin,app}/service/configs/oss.yaml:7`（`access_key:"root"` 硬编码 + `${minio_secret_key:dev_only_change_me}`）
- **现状**：采用 kratos `${env:default}` 模式，默认值统一弱串 `dev_only_change_me`。模式本身合理，但 (1) 弱默认落在仓库；(2) 生产未设环境变量时 JWT 签名密钥退化为已知串 → 任意租户 token 可伪造；(3) 无 CI、无启动校验强制覆盖；(4) `oss.yaml` 硬编码 `access_key:"root"`。
- **部分良好**：`NewAuthenticator` 会 panic 拒绝 `dev_only_change_me_in_prod` 默认 JWT 密钥（`backend/app/*/service/internal/data/data.go`，审计已确认）——但 AES/DB/Redis/Elastic/Minio 凭证无此校验。
- **建议**：启动时对关键 key 做"非默认值+长度+随机性"校验，生产 profile 命中默认值即 panic；移除 `oss.yaml` 硬编码 `access_key:"root"`；`.gitignore` 增加 `.env*`（注：实测 `.env` 当前未跟踪，此为预防）。

### AUD-H6 — 供应链依赖治理不足：Dependabot 仅 github-actions + tx7do/* 低成熟度包群集中在 authn/authz/transport `[审计调查]`
- **位置**：`.github/dependabot.yml`（仅 `package-ecosystem: github-actions`，monthly，未配 gomod/npm/pub/docker）；`backend/go.mod`
- **现状 (Dependabot)**：未覆盖后端 `go.mod`（260 直接依赖）、前端 `frontend/admin` + `frontend/app/{react,taro,vue}` 多个 `package.json`、flutter `pub`、docker。结合 AUD-C2（无 CI），依赖升级既不自动也不被测试覆盖。
- **现状 (tx7do/*)**：依赖大量 `github.com/tx7do/*` 包，多为 `v0.0.x`/`v0.1.x`，占据鉴权/鉴权决策/引导/传输安全敏感位置：`kratos-authn`、`kratos-authn/engine/jwt`、`kratos-authz`、`kratos-authz/engine/{casbin,opa}`、`kratos-authz/middleware`、`kratos-bootstrap/*`、`kratos-transport/transport/{sse,asynq}`、`kratos-swagger-ui`、`go-crud/entgo`、`go-utils/crypto`。单一维护者低版本包群集中 = 作者停维/被投毒/语义变更的集中性供应链风险。
- **主链路依赖**（已 pinned，未见明显未修补高危 CVE）：`go-kratos/v2 v2.9.2`、`entgo.io/ent v0.14.6`、`golang-jwt/v5 v5.3.1`、`pgx/v5 v5.9.2`、`go-redis/v9 v9.19.0`、`minio-go/v7 v7.0.99`。**未跑 govulncheck/osv-scanner，CVE 核对不完整**。
- **建议**：(1) Dependabot 扩展 `gomod`(`/backend`)、`npm`(`/frontend/admin`、`/frontend/app/{react,taro,vue}`)、`pub`(`/frontend/app/flutter_app`)、`docker`，weekly，自动 PR；(2) CI 加 `govulncheck ./...` + `osv-scanner` + `pnpm audit --prod`；(3) 对 tx7do/* 做集中评估（star/提交频率/可替换性），pin 到可信 commit 或自建 fork，长期评估迁移官方包；(4) 为 authn/authz 链路补集成测试锁定行为。

---

## 4. Medium 发现

### AUD-M1 — 公开 app 内容路由无 host→租户解析，当前为功能断链，修复不当则跨租户泄露 `[审计调查]`
- **位置**：`backend/app/app/.../rest_server.go:40-58`（14 个白名单 List/Get：Post/Page/Category/Comment/Tag/Section/Navigation）；`backend/pkg/middleware/ent/ent.go:66-69`（白名单 md==nil 注入 SystemViewer）；`app/.../service/{post,page,category}_service.go:29-58`（仅按 status 过滤，无租户过滤，透传 core）；`go-crud/entgo/rule/tenant.go:76-79`（System/Platform context skip 租户谓词）
- **现状**：app BFF 对 14 个操作走白名单 → auth 中间件被 selector 跳过 → 不注入 OperatorMetadata → BFF ent.Server 注入 SystemViewer。但 BFF→core 的 gRPC 调用不传播 operator metadata（只有 `auth.Server` 调 `metadata.NewContext`，而它被跳过）。core 侧 `ent.Server()` 见 md==nil 且 core `DefaultWhiteList` 为空 → fail-closed 返回 `ErrMissingOperatorMetadata`。
- **结论**：当前这 14 个公开内容路由**会在 core 处报错**（功能异常），而非返回跨租户数据。跨租户信息枚举是**潜在**风险——一旦 (a) BFF 为白名单注入并传播 metadata，或 (b) core `DefaultWhiteList` 被填充对应操作名，SystemViewer 即在 core 放行全租户已发布内容。无 host→租户解析中间件的问题确实存在。
- **建议**：引入 host→tenant 解析中间件，为公开内容显式注入目标 tenantId 的 UserViewer（而非 SystemViewer），core 侧据此做租户谓词过滤。**修复方向不当会直接升级为跨租户泄露**。

### AUD-M2 — WithInjectTenantId 未启用，file_transfer 信任客户端传入 tenantId/userId `[审计调查]`
- **位置**：`backend/app/{admin,app}/.../rest_server.go:72-74`/`:74-76`（仅 `WithInjectMetadata(true)`+`WithInjectEnt(true)`，**未**调 `WithInjectTenantId`）；`backend/pkg/middleware/auth/auth.go:75-80`（`injectTenantId` 默认 false，`setRequestTenantId` 永不执行）；`backend/app/core/.../file_transfer_service.go:170-177,115-116`（`recordFile` 用 `req.GetTenantId()`/`req.GetUserId()` 落库）
- **现状**：`WithInjectTenantId` 全仓无调用方。`recordFile` 直接用请求体 tenantId/userId 写文件元数据。`TenantPrivacy.EvalMutation` 对 UserViewer 会强制覆盖 tenant_id（`tenant.go:114-127`），故 tenant_id 最终由 viewer 修正，但 `CreatedBy`（`created_by` 字段）未被强制覆盖——若 File 实体挂 `OwnerOnlyRule`，伪造的 CreatedBy 会扭曲归属。
- **建议**：`recordFile` 改用 viewer 的 `UserID()/TenantID()`，不从 `req` 取；或启用 `WithInjectTenantId` 让请求字段先被 JWT 值覆盖。

### AUD-M3 — 前端三应用仍 console.log 明文密码/Token（R7-M7 未修） `[审计调查]`
- **位置**：`frontend/app/react/src/app/[locale]/login/components/AccountLoginPage.tsx:19`、`EmailLoginPage.tsx:19`（打印明文密码）；`frontend/app/vue/app/components/auth/AccountRegister.vue:25`（打印明文密码）；`frontend/app/taro/src/store/core/access/store.ts:57,62`（打印 access/refresh token）
- **现状**：第 7 轮 R7-M7 要求移除前端凭证日志，**至今未修**。React/Vue 登录注册页直接 `console.log` 明文密码；Taro 端打印 token 明文。生产构建若未剥离 console，凭证将进入浏览器控制台、错误监控（Sentry 等）、客服截图。
- **建议**：删除这 4 处 `console.log`，或改 `import.meta.env.DEV` 守卫；配置生产构建 `drop_console`。

### AUD-M4 — 构建产物/大体积资产入库 + .dockerignore 不充分 `[审计调查]`（大文件清单亲验）
- **位置/事实**：`frontend/app/flutter_app/.tools/openapi-generator-cli.jar`（约 28MB，入库 `[亲验]`）；`frontend/app/flutter_app/assets/fonts/Noto*`（约 57MB 字体）；`frontend/app/{vue,taro,react}/public/placeholder.png` 各 3.6MB（重复大占位图）；`docs/images/*.png` 多个 1~3.6MB 截图；`backend/app/core/service/internal/data/ent/mutation.go`（2.9MB ent 生成代码）；`.dockerignore`（已读全文）仅排除 `.git/`、`bazel-*`、`bin/`、`*.test.bin`、`cmd/`、`Dockerfile`，**未排除** `configs/`（含弱默认密钥）、`node_modules/`、`*.jar`、`docs/`、前端构建产物
- **风险**：镜像构建上下文可能把含弱密钥的 configs 与大文件一并打入镜像，增大攻击面与镜像体积；仓库膨胀拖慢克隆。
- **建议**：补全 `.dockerignore`（configs/、*.jar、docs/、各前端 node_modules/dist）；评估 `openapi-generator-cli.jar` 改构建时按版本下载并校验 checksum；前端大占位图改用压缩版或 CDN。

### AUD-M5 — 缺少 SBOM/签名/来源证明，无依赖校验步骤 `[审计调查]`
- **现状**：`backend/go.sum`（82KB）与 `frontend/admin/pnpm-lock.yaml`（724KB）均在，依赖锁定良好。但无 SBOM（spdx/cyclonedx）、无 sigstore/cosign 签名、无 SLSA 来源，CI 中无 `go mod verify`/`osv-scanner`/`govulncheck`（因 CI 测试步骤本就不存在，见 AUD-C2）。默认从公共 GOPROXY/npm registry 拉取，仅靠 lockfile 校验。
- **建议**：CI 加 `go mod verify` + `govulncheck ./...` + `osv-scanner`；前端加 `pnpm audit --prod`；构建产物生成 SBOM 并（理想情况）cosign 签名。

### AUD-M6 — CodeQL 扫描范围仅限 Go，前端与 Actions 未纳入 `[审计调查]`
- **位置**：`.github/workflows/codeql-analysis.yml`（matrix `language: ['go']`）
- **现状**：前端 4 个工程（vben/vue、react、taro、flutter）及 GitHub Actions 自身 yaml 不在 SAST 范围。
- **建议**：matrix 增加 `javascript-typescript`、`actions`；对 `configs/` 做密钥扫描（github/ossar 或 secret-scanning）。

---

## 5. Low 发现

### AUD-L1 — SystemViewer 审计跳过 + 审计日志签名密钥进程内临时生成 `[审计调查]`
- **位置**：`backend/pkg/entgo/viewer/system_viewer.go:21-23`（`ShouldAudit()` 返回 false）；`backend/pkg/middleware/logging/logging.go:24-26`（`ecPrivateKey/ecPublicKey` 在 nil 时每次启动 `generateECDSAKeyPair()` 临时生成）；`api_audit_log.go:85`（审计写入用 `NewSystemViewerContext`）
- **风险**：审计日志签名密钥每次进程重启重新生成，无法跨重启验签；SystemViewer 标记不审计。影响审计链完整性而非直接越权。
- **建议**：签名密钥持久化；审计写入不应依赖 SystemViewer 跳过逻辑。

### AUD-L2 — MinIO HMAC 密钥硬编码（死代码） `[审计调查]`
- **位置**：`backend/pkg/oss/utils.go:35`（`var staticHMACSecret = []byte("0123456789abcdef0123456789abcdef")`）
- **现状**：所有活跃上传路径调用 `EnsureObjectName(..., oss.GenerateFileNameTypeUUID)`，从不使用 `GenerateFileNameTypeHMACContent`。`staticHMACSecret` 仅被无活跃调用方的 `GenerateHMACContentFileName` 引用。
- **风险**：当前为零，但未来切换 HMACContent 命名策略时硬编码密钥使文件名可被离线伪造/探测。
- **建议**：改配置/环境变量注入，或直接删除未用 HMAC 分支。

### AUD-L3 — JWT 缺少 issuer/audience 校验，admin/app 共用同一签名密钥 `[审计调查]`
- **位置**：`backend/pkg/jwt/user_token_payload.go:62-80`（claims 仅 sub/uid/tid/iat/exp/jti）；`backend/app/core/.../authenticator.go:117-206`（验证流程）
- **现状**：token 无 iss/aud，验证端不校验。签名算法固定 HS256，`authnJwt` 引擎在 `jwt.go:68` 显式校验 `jwtToken.Method == signingMethod`，`keyFunc` 始终返回非空 HMAC key → `none` 算法与算法混淆不可绕过。exp + Redis jti 白/黑名单均生效。故当前影响有限，但 admin 与 app 共用同一 `${jwt_signing_key}` 且无 audience 区分，理论上 admin token 可被 app 侧校验接受（clientType 在 `ValidateTokenRequest` 中区分，但密钥同源）。
- **建议**：为 admin/app 分别注入 audience claim 并在验证时校验；或为两者用不同密钥。

### AUD-L4 — Lua 引擎沙箱良好但完全未接线（死代码） `[审计调查]`
- **位置**：`backend/pkg/lua/engine.go:99-115,117-185`；`backend/pkg/lua/api/{oss,cache,crypto,eventbus,task,logger,util}.go`
- **现状**：`pkg/lua` 沙箱质量好（`SkipOpenLibs:true`，仅开 base/table/string/math，`dofile/loadfile/load/loadstring` 置 nil，自定义最小 `require`，无 os/io）。但全仓无任何 `app/` 下服务 import `pkg/lua`，引擎从未被实例化、无 HTTP/RPC 入口暴露 Lua 执行。
- **风险**：当前无注入面。未来接线（脚本化任务/动态 hook）需注意：OSS API 以 `context.Background()` 调用、无租户隔离、无大小限制；crypto/task API 需资源限制。
- **建议**：接线时为 Lua OSS API 注入请求 ctx + 租户作用域 + bucket 白名单 + 资源限制；或确认不接线则移除。

### AUD-L5 — 失效的 go.yml__ 残留仓库 `[亲验]`
- **位置**：`.github/workflows/go.yml__`
- **建议**：删除（与 AUD-C2 修复一并）。

### AUD-L6 — server.exe（151MB）残留工作树 `[亲验，纠正定级]`
- **事实**：实测**未入库**（被 `backend/.gitignore:3 *.exe` 忽略，`git ls-files --error-unmatch` 未跟踪，本地 ref `git log --all --follow` 无历史）。仅本地构建产物残留。
- **风险**：低——常驻源码树拖慢备份/IDE 索引，存在 `git add -f` 误提交风险。
- **建议**：删除本地 `backend/server.exe`，构建产物统一到已被 ignore 的 `bin/`。（注：远端历史是否曾提交后改写未核验，据需可 `git log --all --source` 远端再确认。）

### AUD-L7 — 无自动化 license 检查与 THIRD-PARTY 清单 `[审计调查]`
- **现状**：根 `LICENSE` = MIT（Copyright 2022 Bobo）。关键依赖多为 MIT/Apache-2/BSD（ent、opa、kratos 均 Apache-2），未见明显 copyleft 冲突。但无自动化 license 检查与 THIRD-PATY 清单。
- **建议**：加 `go-licenses`/`license-checker` 步骤并生成 NOTICES。

---

## 6. 已确认良好项（快照判断，但因 AUD-C2/C3 无法持续保证）

| 维度 | 结论 | 依据 |
|---|---|---|
| SQL 注入 | 无注入点 | 全仓无 GORM；Ent `b.Arg()` 参数化；字段名 `s.C(field)`；jsonKeyPattern 白名单严格 |
| SSRF | 防护良好 | `netutil.SafeHTTPClient` dial 时拒绝 loopback/private/link-local/unspecified/multicast，封堵 DNS rebinding；body 限 10MiB，超时 30s；无 webhook/OAuth 回调/图片代理 |
| CORS | 安全 | `gorilla/handlers.CORS` + 显式 origin 列表，无 `*`，无 `AllowOriginFunc` |
| 路径遍历 | 无穿越面 | 无 `http.FileServer/ServeFile`；MinIO 扁平 KV；`DownloadFile` 用 DB 存值非用户输入 |
| 反序列化 | 安全 | 无 gob；json.Unmarshal 到具体结构体；protobuf 类型受限 |
| 文件上传 | 基本良好 | UUID 化对象名；预签名直传已禁用（not implemented）；bucket 按 MIME 分类；孤儿对象回滚。**缺大小/扩展名黑名单**（对象名 UUID 化使路径穿越影响有限） |
| SSE 鉴权（R7-H7） | **已确认修复** | app 侧 `sse_server.go:35-56` 加 `WithAuthorizeFunc` 调 `ValidateToken(ClientType_app, ACCESS)`；admin 侧 `internal_message_service.go:81-118` 同样鉴权 + `stream==token` 绑定校验 + 日志 `hashToken`(SHA256 前16) 脱敏 |
| DTM 滥用 | 低风险 | 全仓无业务 saga 注册；`*WithTx`/`Transaction` 实为单库 ent Tx；ctx 贯穿，租户谓词事务内持续生效 |
| JWT 算法 | 良好 | 固定 HS256，校验 method，`none`/算法混淆不可绕过；exp + Redis jti 黑/白名单生效 |
| 默认 JWT 密钥拒绝 | 良好 | `NewAuthenticator` panic 拒绝 `dev_only_change_me_in_prod`（但仅 JWT，见 AUD-H5） |

> 重要：以上"良好"是基于当前代码快照的判断。因 AUD-C2（无 CI）+ AUD-C3（无回归测试），这些属性**无法在后续变更中被自动保持**。建立 CI 与测试是保证它们持续良好的前提。

---

## 7. 本轮审计局限

1. **静态快照**：基于代码静态分析，未动态验证（未跑服务、未做渗透、未压测）。
2. **生产部署未核验**：环境变量是否设置、core gRPC 是否网络隔离、TLS 配置均未核验。AUD-C1 的实际可利用性强依赖生产网络隔离。
3. **第 7 轮已做项未逐一复核**：本轮聚焦遗留项 + 新维度 + 工程侧，未复核第 7 轮 R7-M1/M2/M3/M4/M6（24 个 repo updated_by、SyncPermissions、updateMask translations、前端 14 drawer）的落地情况。
4. **前端逐页安全未深扫**：4 个前端应用的逐页 XSS（dangerouslySetInnerHTML/v-html）、状态泄露、存储型 XSS 未逐页审计。
5. **移动端特有安全未深扫**：taro（小程序）/flutter 的本地存储、网络证书校验、deeplink 未深扫。
6. **CVE 核对不完整**：未跑 govulncheck/osv-scanner/npm audit，依赖 CVE 判断基于版本号经验，不完整。

---

# 修复路图（分阶段）

> 原则：先止血（直接可利用 + 回归保护基线）→ 再收敛（剩余 High/Medium + 纵深）→ 再加固（工程卫生 + 供应链完整性）→ 持续治理。每阶段交付可验证，验收以测试 + CI 通过为准（依赖 AUD-C2 先落地）。

## P0 — 止血（0~2 周）：消除可被直接利用的 Critical/High + 启用回归保护基线

| 工作项 | 映射发现 | 交付物 | 验收 |
|---|---|---|---|
| core gRPC 鉴权加固 | AUD-C1 | `x-md-global-operator` HMAC 签名（启用预留 `mdSignatureKey`），core 侧校验；core gRPC 绑定内网地址 + mTLS；移除 `0.0.0.0` 公开 | 伪造/无签名头被 core 拒；集成测试覆盖 |
| 审计日志体系修复（写入侧） | AUD-H3, AUD-H4 | app BFF 接真实审计写入；core gRPC 加审计中间件覆盖写操作；`api_audit_log` 请求体字段级脱敏（password/secret/token → `***`）+ 密码类端点排除 | app/core 写操作有审计；密码不落库；测试覆盖 |
| 启用真 RBAC 引擎 | AUD-H1 | BFF `NewAuthorizer` 改用 `authorizer.Authorizer` 包装器 + 配 `authz.yaml`；core gRPC 挂 authz 中间件 | noop 替换；越权请求被拒；测试覆盖 |
| 登录限流 + 锁定 + captcha fail-closed | AUD-H2 | 登录流程接入 `LoginPolicy`（锁定阈值/时长）；IP+账号 rate limit；captcha `nil` 时 fail-closed；app 登录加验证码 | 暴力破解被限；captcha 缺失时拒绝；测试覆盖 |
| 恢复 CI（基线） | AUD-C2 | 新建 `.github/workflows/ci.yml`：后端 `go vet && go build ./... && go test -race -coverprofile`，前端 `pnpm install --frozen-lockfile && pnpm lint && pnpm check:type && pnpm test:unit`，PR 必过；删除 `go.yml__`（AUD-L5） | PR 无绿色 CI 不合并；go.yml__ 删除 |
| 安全路径回归测试基线 | AUD-C3 | 表驱动 + sqlite/miniredis 集成测试，覆盖 P0 修复项：auth 中间件、ent 中间件 viewer/TenantID、高风险 repo 跨租户拒绝、SSE `HandleAuthorize` 三分支、updateMask 越权字段剥离 | 覆盖率门槛（如关键包 ≥60%）；CI 跑通 |
| admin/app JWT audience 区分 | AUD-L3 | 为 admin/app 分别注入 audience claim 并在验证时校验（顺带，因同批改 JWT 侧） | admin token 不被 app 接受；测试覆盖 |

P0 完成判定：Critical 3 项 + High 中 H1/H2/H3/H4 全部关闭，CI 与回归测试基线上线。AUD-H5/H6 进入 P1。

## P1 — 收敛（2~4 周）：剩余 High/Medium，建立纵深

| 工作项 | 映射发现 | 交付物 | 验收 |
|---|---|---|---|
| 密钥 fail-fast | AUD-H5 | 启动时校验生产 profile 关键 key（AES/DB/Redis/ES/Minio）非默认+长度+随机性，命中 panic；移除 `oss.yaml` 硬编码 `access_key:"root"`；`.gitignore` 加 `.env*` | 默认密钥生产启动失败；root access_key 移除 |
| 供应链治理 | AUD-H6 | Dependabot 扩展 gomod/npm/pub/docker，weekly，自动 PR；CI 加 `govulncheck`/`osv-scanner`/`pnpm audit`；tx7do/* 集中评估报告（pin commit/自建 fork/迁移路径） | Dependabot 覆盖全生态；漏洞扫描入 CI；评估报告产出 |
| 公开内容路由 host→租户解析 | AUD-M1 | host→tenant 解析中间件，为公开内容注入目标 tenantId 的 UserViewer（非 SystemViewer），core 据此做租户谓词 | 14 路由按 host 返回对应租户已发布内容；跨租户不泄露；功能断链修复；测试覆盖 |
| file_transfer 改用 viewer 来源 | AUD-M2 | `recordFile` 改用 viewer `UserID()/TenantID()`，或启用 `WithInjectTenantId` | tenant_id/CreatedBy 由 viewer 决定；测试覆盖 |
| 前端凭证 log 移除 | AUD-M3 | 删除 4 处 `console.log`；生产构建 `drop_console` | grep 无凭证 log；构建配置生效 |
| CodeQL 扩展 | AUD-M6 | matrix 增加 `javascript-typescript`、`actions`；configs 密钥扫描 | CodeQL 覆盖前端+Actions |

P1 完成判定：H5/H6 关闭，M1/M2/M3/M6 关闭，纵深防御（RBAC、限流、审计、密钥）就位。

## P2 — 加固（4~8 周）：工程卫生与供应链完整性

| 工作项 | 映射发现 | 交付物 | 验收 |
|---|---|---|---|
| 清理入库大文件 + .dockerignore | AUD-M4, AUD-L5, AUD-L6 | `openapi-generator-cli.jar` 改构建时下载校验 checksum；字体 CDN；占位图压缩；补全 `.dockerignore`（configs/、*.jar、docs/、前端 node_modules/dist）；删除本地 `server.exe`，产物统一 `bin/` | 仓库瘦身；镜像上下文不含弱密钥 |
| SBOM + 签名 + 依赖校验 | AUD-M5 | CI 加 `go mod verify`；构建产物生成 SBOM（spdx/cyclonedx）；cosign 签名（理想） | SBOM 产出；mod verify 入 CI |
| 审计日志完整性加固 | AUD-L1 | 审计日志签名密钥持久化；审计写入不依赖 SystemViewer 跳过逻辑 | 跨重启可验签；审计覆盖 SystemViewer 路径 |
| MinIO HMAC 死代码清理 | AUD-L2 | 移除 `staticHMACSecret` 或配置化 | 死代码移除 |
| Lua 接线安全准备 | AUD-L4 | 接线前为 Lua OSS API 注入 ctx+租户作用域+bucket 白名单+资源限制；或确认不接线则移除 | 决策记录；若保留则沙箱增强 |
| CodeQL 自定义查询 | AUD-M6（深化） | 租户谓词缺失、SystemViewer 滥用、敏感字段落库 的自定义查询 | 自定义查询入 CodeQL |
| 自动化 license 检查 | AUD-L7 | `go-licenses`/`license-checker` 步骤 + NOTICES 生成 | license 检查入 CI；NOTICES 产出 |

P2 完成判定：M4/M5 + L1/L2/L4/L5/L6/L7 关闭，供应链完整性基线建立。

## P3 — 持续治理

- 测试覆盖扩展至全部 `*_repo.go` 的跨租户拒绝（补 AUD-C3 未覆盖的 55 个 repo）。
- 前端逐页安全审计（XSS/v-html/状态泄露）+ 移动端特有安全（taro/flutter 本地存储/证书校验/deeplink）—— 补本轮局限 3/4/5。
- 第 7 轮已做项（R7-M1/M2/M3/M4/M6）逐一复核落地情况 —— 补本轮局限 3。
- 定期重审（第 9 轮+）+ 第三方渗透测试。
- 安全 RFC：BFF/core 职责边界、metadata 传播契约、租户隔离机制文档化。
- 动态验证：搭测试环境跑服务，做越权/撞库/SSE 伪造的端到端验证（补本轮局限 1/2）。

---

## 附录 A：发现与路图阶段映射总表

| 发现 | 严重度 | 阶段 |
|---|---|---|
| AUD-C1 core gRPC 无鉴权无签名 | Critical | P0 |
| AUD-C2 CI 禁用失效 | Critical | P0 |
| AUD-C3 回归测试缺口 | Critical | P0（基线）/ P3（全覆盖） |
| AUD-H1 RBAC noop | High | P0 |
| AUD-H2 登录无限流/captcha fail-open | High | P0 |
| AUD-H3 审计日志体系（写入缺失） | High | P0 |
| AUD-H4 审计日志含密码明文 | High | P0 |
| AUD-H5 弱默认密钥无 fail-fast | High | P1 |
| AUD-H6 供应链治理不足 | High | P1 |
| AUD-M1 公开内容路由跨租户/断链 | Medium | P1 |
| AUD-M2 WithInjectTenantId 未启用 | Medium | P1 |
| AUD-M3 前端凭证 log（R7-M7 未修） | Medium | P1 |
| AUD-M4 大文件入库/.dockerignore | Medium | P2 |
| AUD-M5 无 SBOM/签名/依赖校验 | Medium | P2 |
| AUD-M6 CodeQL 仅 Go | Medium | P1（扩展）/ P2（自定义查询） |
| AUD-L1 审计签名密钥/SystemViewer 跳过 | Low | P2 |
| AUD-L2 MinIO HMAC 硬编码死代码 | Low | P2 |
| AUD-L3 JWT 无 iss/aud | Low | P0（顺带） |
| AUD-L4 Lua 未接线死代码 | Low | P2 |
| AUD-L5 go.yml__ 残留 | Low | P0（随 CI 删除） |
| AUD-L6 server.exe 本地残留 | Low | P2 |
| AUD-L7 无 license 自动检查 | Low | P2 |

## 附录 B：关键文件索引

- 中间件链：`backend/pkg/middleware/{auth,ent,logging}/*.go`、`backend/app/{admin,app,core}/service/internal/server/{rest_server,grpc_server,sse_server}.go`
- 租户隔离：`backend/app/core/service/internal/data/*_repo.go`、`backend/pkg/entgo/{viewer,rule}/*`、`backend/pkg/metadata/operator_metadata.go`
- 鉴权：`backend/app/core/service/internal/data/authenticator.go`、`backend/pkg/jwt/*`、`backend/pkg/authorizer/*`
- 审计日志：`backend/pkg/middleware/logging/{api_audit_log,login_audit_log,logging}.go`
- 配置：`backend/app/*/service/configs/{authenticator,data,oss,server,client,logger}.yaml`
- CI/供应链：`.github/{workflows,dependabot.yml}`、`backend/go.mod`、`frontend/admin/pnpm-lock.yaml`
- 前端凭证 log：`frontend/app/{react,vue,taro}/src/**`
