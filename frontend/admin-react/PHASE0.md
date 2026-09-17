# admin-react Phase 0 冒烟报告（2026-09-16，完整版）

基座（go-wind-admin `frontend/admin/react`）复制为 `frontend/admin-react/`，对接 CMS 后端（127.0.0.1:6600 / SSE 6601）。
本轮覆盖：登录（平台+租户）、动态路由/菜单、读 CRUD、写 CRUD 回环、登出、租户数据隔离、SSE 事件流，以及静态漂移全量扫描（缺失服务清单 + tsc 184 错误清单）。
结论：**主链路全部打通**；发现的契约差异 D1–D14 全部记录在案，其中 D1/D5/D8 已修并实测验证，品牌文案已对齐，其余归入 Phase 1。

## 冒烟结果（全部经 vite 代理，浏览器内实测）

| 测试点 | 结果 | 证据 |
|---|---|---|
| 登录页渲染（品牌/i18n/验证码图） | ✅ | gui-test-screenshots/smoke-01-login-initial.png |
| 图形验证码（GET /admin/v1/captcha，经 vite 代理） | ✅ 200 | 网络层 200，img naturalWidth=240 |
| AES 登录（CBC, key=iv=VITE_AES_KEY, base64）平台账号 admin/admin | ✅ | 返回 bearer JWT（roles: platform:admin） |
| **租户登录**（租户编号 super + tenant_admin/admin） | ✅ | 会话身份 `{id:2, tenantId:1, tenantName:测试租户}`；干净态下头部身份标签正确显示 tenant_admin |
| 登录后链路（me / perm-codes / dict / internal-message inbox） | ✅ 全 200 | 网络层 |
| 动态菜单（7 模块 + 系统管理 11 子项） | ✅ | DOM 快照（注意 D11：菜单来自本地路由模块，未经权限过滤） |
| 字典管理主从列表（读路径） | ✅ | smoke-02-dict-list.png；`total:"3"` 字符串契约正确渲染 |
| **写路径 CRUD 回环**（新建类型→列表验证→删除→验证消失） | ✅ | 抽屉表单（名称/编码/排序/状态）提交 200，列表总数 5→6、alert“创建成功”；删除确认后 6→5、alert“删除成功”。测试行已清理，库回到 5 行 demo 态 |
| **租户数据隔离** | ✅ | 租户 token 下 dict List 返回 `{"items":[],"total":"0"}`（服务端 TenantPrivacy 过滤 tenant_id=0 平台行）；同一接口平台 token 返回 5 行 |
| **登出流** | ✅ | POST /admin/v1/logout → 200；会话态清理并跳回登录页 |
| **SSE 事件流**（D8 修复后） | ✅ | 登录即建立连接：GET 6601/events?stream=<JWT> → 200，`content-type: text/event-stream, connection: keep-alive` |
| 控制台错误 | ✅ 0（登录页）| MFA 补丁后；仪表盘 1 处已知告警（antd v6 Drawer width 弃用，D6） |

## 契约差异与处理（D1–D14 全量）

| # | 状态 | 事项 |
|---|---|---|
| D1 | **已修** | `api/hooks/auth.ts` 引用 `apiClient.mfaService.VerifyMFAChallenge`，CMS 生成客户端无独立 MFA 服务 → 模块求值即崩。已加可选守卫兜底（缺服务时 mutation 挂拒绝函数）。MFA 挑战页暂不可用（CMS 登录响应无 mfa_operation_id，不会触发） |
| D2 | 已配 | `.env`：标题/namespace 改 CMS（gowind-cms）；`.env.development`：代理 `/admin`→6600、`VITE_SSE_URL`→6601/events。AES key 与基座一致（f51d66a73d8a0927），无需改。注：SSE 跨源 CORS 疑虑已证伪——6601 预检返回 `Access-Control-Allow-Origin: *`（含 Authorization 头），直连可用，无需代理化 |
| D3 | 已换 | `src/api/generated` 整体替换为 CMS vben 同款生成客户端（421 导出），`ClientTransport` 接口形状兼容，业务 hooks 零改动 |
| D4 | 环境坑 | **6600 端口双监听**：go-wind-admin 后端（DB=gwa）与 CMS 后端共存时，Go 的 SO_REUSEADDR 使连接随机分发——表现为“时好时坏的密码错误”。若再出现随机 401/INVALID_PASSWORD，先查 `netstat -ano | grep 6600` 是否双 LISTENING |
| D5 | **已修（本轮）** | JsonEditor.tsx 残留 `require('jsoneditor/dist/jsoneditor.min.css')` 死壳（ESM 下 require 未定义，ReferenceError 被静默吞掉，CSS 从未加载）+ package.json 声明的 `json-editor-vue`（Vue 专用组件库，React 侧零引用）。两者已移除并 `pnpm install` 同步锁文件。**遗留 parity 缺口**：vben 侧 JSON 编辑为结构化 tree/form 视图（json-editor-vue），React 侧为带解析校验的 textarea 兜底；当前 React 无页面消费 EDITOR_TYPE_JSON_BLOCK，Phase 1 若页面编辑器移植需要，可改用框架无关的 vanilla-jsoneditor 命令式挂载 |
| D6 | 待 Phase 1 | antd v6：Drawer `width` 弃用（改 `size`）等 API 迁移告警，逐页清理 |
| D7 | **已实证+映射（本轮）** | 仪表盘四张统计卡（用户总数/角色总数/今日登录次数/今日操作审计条数）与三张图为**本地假 0/空**：`api/hooks/dashboard.ts` 调 `dashboardService.*`，CMS 客户端无此服务（queryFn 抛错被吞，无任何 stats 请求发出，网络层可证）。CMS 侧真服务为 `statsService`（vben 分析页在用），映射：`GetOverview→GetDashboardOverview{userCount,postCount,commentCount,interactionCount,...}`（用户总数有对应；角色总数/今日操作审计条数无对应）；`GetLoginTrend→GetLoginActivity{success:DailyCount[],failed:DailyCount[]}`（登录趋势图直接对应，成功/失败占比可由其聚合）；`GetOperationActionDistribution→无对应`。Phase 1 决策：重接 statsService 并裁掉无对应的卡片，或整页换 CMS 分析页设计 |
| D8 | **已修+实测（本轮）** | 两处漂移：① 登录成功流（`applySuccessfulLogin`）从不调 `connectSSEServer`（vben 在 authentication.store.ts:477 登录即连），唯一调用点在 bootstrap 静默恢复分支而该分支本身被 D12 废掉 → SSE 实际永远不连；② stream 参数传 `userInfo.id`，而 CMS 后端 `HandleAuthorize` 强制 `stream==token`（internal_message_service.go:107-113）否则 403。已对齐 vben：登录流接入连接 + stream 改传 token（`useTokenRefresh.ts` 两处）。实测：登录即建立事件流（200，text/event-stream） |
| D9 | 待 Phase 1 | **静态漂移：CMS 客户端缺失服务 13 个**——基座 hooks 引用而 CMS 生成客户端不存在：scriptService、planService、planQuotaService、planModuleService、mfaService、notificationChannelService、configService、onlineSessionService、accessKeyService、dashboardService、scriptLogService、serverMonitorService、redisCacheMonitorService。对应 6 个上游独有模块（脚本/套餐/参数/在线用户/服务监控/通知渠道）+ MFA + 仪表盘。这些页面在本仓为死页（路由被权限过滤，D11），但代码与菜单残留 |
| D10 | 待 Phase 1 | **静态漂移：tsc 184 错误 / 41 文件**，构成：TS2339 属性缺失×56、TS2305/TS2724 缺失导出×84、TS2551×20、TS7006 隐式 any×20 等。热点文件：plan.ts(26)、script.ts(16)、notification-channel.ts(12)、mfa.ts(11)、config.ts(11)、online-session.ts(10)、access-key.ts(10)、dashboard.ts(8)、redis-cache-monitor.ts(6)、script-log.ts(6)——**全部落在 D9 缺失服务域**，随模块删除/裁剪自然消除。**共有服务真漂移 3 处**：tenant.ts（TenantUsage/GetUsage/CleanupData 缺）、user-profile.ts（UploadAvatar/DeleteAvatar 缺）、stores/auth.ts（LoginResponse 无 mfa_operation_id 字段）。清单项：全量清单在 /tmp/adminreact-tsc.log（会话产物，需重跑时用 `npx tsc --noEmit` 于 frontend/admin-react） |
| D11 | 已实证（本轮） | 菜单生成（generate-menus.ts，本地路由模块）不做权限过滤，路由注册（generate-routes-frontend.ts）按 perm-codes 过滤——上游 6 个独有模块（脚本管理/参数管理/在线用户/服务监控/通知渠道/OpenAPI 凭证）**菜单可见但路由落空**，点击落仪表盘（/system/scripts 实测）。另有路由模块内 `permission` 字段全部被“开发阶段暂时注释”掉的基座遗留，需在 Phase 1 恢复并按 CMS 权限码（sys:platform_admin 等）重标 |
| D12 | 已实证（本轮） | 页面硬刷新后 refresh cookie 不恢复会话——bootstrap 静默恢复分支（bootstrap.ts:79-99）实测未生效，一律弹回 /auth/login?redirect=...。vben 侧同款恢复逻辑可用，需对照排查（refresh_exp cookie 读取/refresh-token 端点/dev 代理路径三者之一断链） |
| D13 | **已修+实测（本轮）** | **跨登录权限/身份残留（安全相关）**：登出只清 auth store 与 localStorage（`logout` 移除 user-storage 键），**从不调 `useUserStore.$reset()`**，内存态 userInfo/roles/accessCodes 原样保留；重登时路由守卫经 `useAuth.ts:47` 缓存跳过分支（userInfo 非空且 roles 非空即不拉新）→ 新会话整体继承前一账号的身份标签、角色与权限码。实测：平台 admin 登出后同 tab 登录 tenant_admin，头部身份标签仍显示 "admin"（干净态对照实验显示 "tenant_admin" 证明为残留）；租户会话以 platform:admin 权限码运行。**已修**：logout/forceLogout 补 `useUserStore.getState().$reset()`（router/index.tsx:122 同款先例）。**回归验证通过**：同序列（平台登录→登出→同 tab 租户登录）后标签正确显示 tenant_admin |
| D14 | 待修（本轮发现） | 重登后 tab 系统残留失效：tab 列表跨登出/登录存续（布局层内存态 store 未随登出清理，与 D13 同模式但不同 store），残留 tab 无法激活/关闭，全页刷新后才恢复。D13 修复后复查仍复现（残留"字典管理"tab）。修法：登出时同步重置布局 tab store |

## 品牌文案对齐（本轮完成并验证）

对齐参照为 vben 侧 2026-09-16 品牌同步落地（GoWind 风行 / 风行CMS后台管理系统）：

- `src/core/preferences/config/default.ts`：`name` 由硬编码 `"GoWind Admin"` 改读 `import.meta.env.VITE_APP_TITLE`（= "GoWind CMS Admin"，同 vben preferences.ts:11）——侧边栏品牌字与 document.title 均走此值
- `src/locales/zh-CN/_core/auth.json` systemTitle → “风行CMS后台管理系统”（登录页实测生效）
- `src/locales/en-US/_core/auth.json` systemTitle → "Plug-and-play Admin system"（对齐 vben en-US pageTitle）
- `index.html` meta keywords → "GoWind CMS React AntD Vite"
- logo.png 与 vben 侧 md5 一致（品牌图本就同源，无需改）

## 本机环境备忘

- Node 须 ≥20.19（vite 8/rolldown）：系统默认 20.10 太老，用 `export PATH="/c/Users/yangl/scoop/apps/nvm/current/nodejs/nodejs:$PATH"`（20.20）再 `pnpm dev`
- 首次 `pnpm install` 后 rolldown 平台绑定（@rolldown/binding-win32-x64-msvc）可能缺失报 MODULE_NOT_FOUND：`pnpm install --force` 恢复
- pnpm 裸跑会命中仓根 package.json——务必在 `frontend/admin-react/` 目录内执行
- 测试便利：验证码明文可在 Redis db0 `gowind-cms:captcha:<captchaId>` 读取（TTL 10min，验证即删）
- 账号：平台 admin/admin；租户 tenant_admin/admin（租户编号 super，tenantId=1，角色 template:tenant:manager，真实权限码 3 个：sys:access_app/sys:access_backend/sys:tenant_manager——无 sys:manage_tenants/sys:platform_admin，可作 D13 修复后的回归基线）
- D13 复现步骤：admin 登录 → 登出 → 同 tab 登录 tenant_admin → 看头部身份标签（残留时显示 admin，修复后显示 tenant_admin）
- 6601 SSE 预检 CORS 全开（ACAO:*，含 Authorization），生产 nginx 需保留 SSE host 透传

## 后续（Phase 1 预告）

**前置阻断项（须先于模块移植修复）**：D14（重登后 tab 残留失效，与 D13 同模式不同 store，D13 已修但 tab store 仍未随登出重置）；D12（刷新弹登录，影响联调效率）；D11 菜单/路由权限过滤 + 路由模块 permission 字段恢复重标；D7 仪表盘重接 statsService。

按 site_setting → media → engagement → content 顺序移植 CMS 四模块；6 个共有模块逐页对齐双向补丁（CMS 侧 15 CRUD bug/24 drawer 审计结论 ↔ 基座侧分栏/序列化守卫修复）；D9/D10 死模块随裁剪决定（删除或保菜单占位）一并清掉 tsc 基线。
