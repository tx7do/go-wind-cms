# AGENTS.md — React (Next.js) 前端开发指南

> 本文件是 `frontend/app/react` 子项目的 AI 编码规范单一事实源，适用于所有支持 AGENTS.md 的 AI 编码工具。Claude Code 通过同级 `CLAUDE.md` 中的 `@AGENTS.md` 引用加载。

## 项目概览

基于 **Next.js (App Router)** 的**静态导出型** Headless CMS 内容展示前端。

**核心技术栈**：Next.js 16 + React 19 + TypeScript + Tailwind CSS v4 + shadcn/ui (Radix) + Zustand + TanStack React Query + next-intl + Axios + Shiki/marked/KaTeX/Mermaid/Tiptap

**包管理器**：pnpm

## 关键架构认知

### 静态导出模式

项目配置为**完全静态导出**，产物为纯 HTML/CSS/JS：

```typescript
// next.config.ts
const nextConfig: NextConfig = {
  output: 'export',          // 静态导出
  trailingSlash: true,       // 生成 /path/index.html 目录结构
  distDir: 'dist',           // 输出到 dist 目录
};
```

### 国际化路由

采用 `[locale]` 动态段实现 URL 前缀多语言路由（`localePrefix: 'always'`）：

```
src/app/[locale]/
├── layout.tsx          # 语言布局（SSG 入口）
├── ClientLocaleLayout.tsx  # 客户端布局（Provider 注入）
├── page.tsx            # 首页
├── post/[id]/          # 文章详情
├── category/[slug]/    # 分类页
└── ...
```

`generateStaticParams()` 为每个语言预生成静态页面。语言切换通过 `key={locale}` 强制重挂载 + `queryClient.clear()` 清缓存。

### 三层 API 架构

```
src/api/
├── generated/          # [自动生成] protoc-gen-typescript-http 产出（禁止手改）
├── service/            # [服务封装] 业务逻辑、参数转换、locale/分页注入
├── hooks/              # [React Hook] useMutation/useQuery + fetch* 纯函数 + 多语言辅助函数
└── index.ts            # 统一导出
```

- **第一层 generated/**：protobuf 自动生成的 HTTP 客户端，禁止手动编辑
- **第二层 service/**：封装业务方法，自动注入 locale、分页参数
- **第三层 hooks/**：封装为 `use*` Hook（组件内）和 `fetch*` 纯函数（Store/非 React 上下文）两种形态

```typescript
// service/post.ts — 自动注入 locale
export async function listPostsRaw(params) {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  return getPostService().List({ query: JSON.stringify(formValues), page: params.paging?.page, pageSize: params.paging?.pageSize });
}
```

### RequestClient — HTTP 通信内核

基于 Axios 封装的全局单例，拦截器链：

```
请求拦截：Token 注入 → Request-ID 注入 → Locale 注入
响应拦截：数据解构 → 401 认证处理 → 错误消息提取
```

- **Token 自动刷新**：401 时自动调用 refresh token，刷新期间后续请求排队
- **请求重认证**：刷新失败清除凭证并重定向登录页
- 初始化在 `StoreProvider.tsx` 的 `RequestClient.init(env.apiBaseUrl, {...})`

### 状态管理 — Zustand + React Context

采用 Zustand 配合 React Context 的混合模式，避免 SSR 全局单例数据泄漏：

```
src/store/
├── StoreProvider.tsx           # 聚合 Provider
└── core/
    ├── access/                 # 认证凭证（accessToken、refreshToken）
    ├── user/                   # 用户信息
    └── loading/                # 全局加载状态
```

- 每个 Store 通过 `create*Store()` 工厂函数创建独立实例
- `useMemo` 确保 store 实例稳定
- Token AES 加密后持久化到 localStorage

## 目录结构

```
src/
├── api/                    # API 三层架构（generated / service / hooks）
├── app/                    # Next.js App Router
│   ├── globals.css         #   全局样式 + 主题变量
│   ├── layout.tsx          #   根布局（StoreProvider、ThemeProvider）
│   └── [locale]/           #   多语言路由
│       ├── layout.tsx / ClientLocaleLayout.tsx / routing.ts
│       └── ...业务页面
├── components/             # UI 组件
│   ├── ui/                 #   shadcn/ui 基础组件
│   ├── layout/             #   布局（Header、Footer、Nav）
│   ├── content/            #   内容渲染器（ContentViewer）
│   └── ...                 #   home/post/category/comment/auth
├── config/                 # 环境变量配置
├── core/                   # 核心基础设施
│   ├── preferences/        #   偏好系统（主题、语言）
│   ├── transport/          #   通信层（rest / sse）
│   └── query-client.ts     #   React Query 全局配置
├── hooks/ i18n/ lib/ plugins/ store/ utils/
└── messages/               # i18n 翻译文件（zh-CN/en-US: app/navbar/page/cms/authentication/comment/enum/settings.json）
```

## 关键约定（必须遵守）

1. **客户端组件必须标记 `'use client'`** — 使用 `useState`/`useEffect`/事件处理的组件文件顶部必须加 `'use client'`
2. **禁止手改 `api/generated/`** — 由 protobuf 自动生成，改 proto 后在后端 `make api`/`make ts` 重新生成
3. **API 双形态** — 组件内用 `use*` Hook；Store/事件处理等非 React 上下文用 `fetch*` 纯函数
4. **多语言内容用辅助函数** — 获取后端实体多语言字段用 `getPostTitle(post)` 等辅助函数（自动匹配当前 locale），不要直接访问 `translations` 数组
5. **主题用 CSS 变量** — 修改 `globals.css` 中的 `--primary` 等变量，所有 `bg-primary`/`text-foreground` 等 Tailwind 类自动跟随。**禁用** Tailwind `dark:` 变体，统一用语义化 CSS 变量
6. **环境变量加 `NEXT_PUBLIC_` 前缀** — 客户端可访问的变量必须以此前缀开头，改 `.env` 后需重启 dev server
7. **样式优先用 Tailwind 类名** — 自定义 CSS 在 `globals.css` 用 `@layer base`/`@utility` 指令添加

## 开发命令

```bash
pnpm install      # 安装依赖
pnpm dev          # 开发服务器 → http://localhost:3000
pnpm build        # 构建静态产物 → dist/
pnpm lint         # 类型检查
```

**环境变量**（`.env.development`）：

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:6700    # 后端 API 地址
NEXT_PUBLIC_APP_TITLE='My CMS'
NEXT_PUBLIC_DEFAULT_LOCALE=zh-CN
NEXT_PUBLIC_TOKEN_KEY=access_token
```

## 新增业务模块 Checklist（以"产品"为例）

```
- [ ] Step 1: 确认 API 类型（后端 protobuf 已生成则跳过；否则在 api/generated/ 定义）
- [ ] Step 2: 封装服务层（api/service/product.ts）
- [ ] Step 3: 封装 Hook 层（api/hooks/product.ts：use* + fetch*）
- [ ] Step 4: 创建页面（app/[locale]/product/[id]/page.tsx，按需加 'use client'）
- [ ] Step 5: 添加导航入口（components/layout/TopNavbar.tsx 或 MobileNav.tsx）
```

## 新增语言

1. 在 `messages/` 下创建新语言目录（如 `ja-JP/`），复制现有 JSON 并翻译
2. 在 `i18n/config.ts` 的 `locales` 数组新增，并导入翻译文件加入 `allMessages`
3. 路由基于 `[locale]` + `generateStaticParams()`，新语言自动生成静态页面

## shadcn/ui 组件

```bash
pnpm dlx shadcn@latest add dialog    # 新增组件
```

组件源码在 `components/ui/`，可直接修改。已有：button、input、select、dropdown-menu、dialog、sheet、avatar、toggle、switch、separator、navigation-menu、carousel、pagination、skeleton、spinner 等。

## 内容渲染管线（ContentViewer）

```
Markdown → marked(自定义 Renderer)
  ├ 代码块 → Shiki 双主题高亮
  ├ 数学公式 → KaTeX
  ├ 流程图 → Mermaid
  └ 表格/图片/链接 语义化
  → DOMPurify XSS 清洗 → 安全 HTML
```

## 对接不同后端

1. `config/env.ts` — 修改 `apiBaseUrl`
2. `api/service/*.ts` — 调整请求参数/响应结构
3. `api/hooks/*.ts` — 调整类型定义
4. `api/generated/` — 后端用 protobuf 则重新生成
5. `StoreProvider.tsx` 的 `RequestClient.init()` 回调 — 自定义认证流程

## 常见错误与纠正

| 错误做法 | 正确做法 |
|---|---|
| 客户端组件漏加 `'use client'` | 文件顶部添加 `'use client'` |
| 手改 `api/generated/` | 改 proto 源后重新生成 |
| 直接访问 `post.translations[0].title` | 用 `getPostTitle(post)` 辅助函数 |
| 用 Tailwind `dark:` 变体做暗色 | 用 CSS 变量 `hsl(var(--foreground))` |
| 环境变量不加 `NEXT_PUBLIC_` 前缀 | 客户端变量必须加前缀 |
| 在非 React 上下文用 `use*` Hook | 用 `fetch*` 纯函数 |
