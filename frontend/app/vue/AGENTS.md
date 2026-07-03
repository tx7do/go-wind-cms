# AGENTS.md — Vue (Nuxt) 前端开发指南

> 本文件是 `frontend/app/vue` 子项目的 AI 编码规范单一事实源，适用于所有支持 AGENTS.md 的 AI 编码工具。Claude Code 通过同级 `CLAUDE.md` 中的 `@AGENTS.md` 引用加载。

## 项目概览

基于 **Nuxt 4 (Vue 3)** 的现代 Headless CMS 内容展示前端，支持 SSR/SSG 双模式部署。

**核心技术栈**：Nuxt 4 + Vue 3.5 + TypeScript + Tailwind CSS v4 + shadcn-vue (Reka UI) + Pinia + TanStack Vue Query + @nuxtjs/i18n + Axios + marked/Shiki/KaTeX/Mermaid/Tiptap

**包管理器**：pnpm

## 关键架构认知

### SSR/SSG 双模式

```ts
// nuxt.config.ts
export default defineNuxtConfig({
  ssr: true,
  nitro: {
    prerender: { routes: ['/'], crawlLinks: true },  // SSG 预渲染首页 + 爬取内部链接
  },
})
```

**SSG 部署策略**：Nitro 预渲染 `/`，动态路由（`/post/:id` 等）由客户端 SPA fallback；构建后自动生成根 `index.html` 重定向到 `/zh-CN/`。

### API 两层架构 + 统一 ApiClient

```
api/
├── client.ts            ← ApiClient 单例（createApiClient(transport)，懒加载各 Service）
├── generated/           ← [自动生成] protoc-gen-typescript-http（禁止手改）
│   └── index.ts             # ApiClient 类（get postService() 懒加载缓存）
└── composables/         ← [业务逻辑 + Vue Query] 纯函数 + use* Hook + fetch* 三合一
```

**第一层 generated/**：protobuf 自动生成，含 `ApiClient` 类（内置懒加载缓存）。**第二层 composables/**：每个文件三合一——业务纯函数 + Vue Query Hook + fetch 方法，直接调用 `apiClient`：

```ts
// composables/product.ts
import { apiClient } from '@/api/client';

// 业务纯函数
export async function listProduct(paging?, formValues?, options?) {
  const locale = getCurrentLocale();  // 自动注入当前语言
  return await apiClient.productService.List({ /* 组装参数 */ });
}
// Vue Query Hook（组件用）
export function useListProduct(options?) {
  return useMutation({ mutationFn: (params) => listProduct(params.paging, params.formValues), ...options });
}
// fetch 方法（Store/非组件上下文用）
export async function fetchListProduct(params) {
  return queryClient.fetchQuery({ queryKey: ['listProduct', params], queryFn: () => listProduct(...) });
}
```

### RequestClient — HTTP 通信内核

`core/transport/rest/` 基于 Axios 封装，拦截器链：Token 注入 → Request-ID → Locale → 401 自动刷新 Token/跳登录 → 响应解构 → 错误消息。初始化在 Nuxt 插件 `plugins/01.init-client.ts`（仅客户端）。

### 状态管理 — Pinia

```
stores/modules/
├── core/   # access（Token/登录态）/ user / navbar / loading
└── app/    # auth（登录、登出、Token 刷新）
```

密码传输前 AES-CBC 加密（CryptoJS），密钥运行时配置注入。

### 偏好系统（PreferenceManager）

`core/preferences/` 实现响应式单例（`reactive` + `readonly` 暴露）：偏好变更自动同步 CSS 变量 + localStorage（防抖写入）+ `@nuxtjs/i18n`。支持 `light`/`dark`/`auto` 主题。

```ts
const { isDark, theme, locale, toggleTheme, setLocale } = usePreferences();
```

### 国际化（@nuxtjs/i18n）

采用 **prefix 路由策略**（所有路由带语言前缀 `/zh-CN/`、`/en-US/`），翻译文件按模块拆分在 `locales/`。语言感知 API 调用用 `getCurrentLocale()`（`utils/locale.ts`）。

### 主题系统（语义化 CSS 变量）

```css
:root { --primary: 142.1 76.2% 36.3%; --background: 210 40% 98%; /* ... */ }
.dark { --primary: 142.1 86.2% 50.3%; --background: 224 45% 6%; /* ... */ }
```

**关键设计**：禁用 Tailwind `dark:` 变体，统一用语义化 CSS 变量（`hsl(var(--foreground))`）控制样式。设置 `@media (prefers-color-scheme: dark)` 纯 CSS 后备防 FOUC。

## 目录结构

```
app/
├── api/                    # 两层架构（client + generated + composables）
├── assets/css/             # 全局样式 + 主题变量（main.css）
├── components/             # auth/category/comment/content/home/layout/post/ui(shadcn-vue)
├── constants/              # 常量
├── core/                   # preferences / storage(StorageManager) / transport(rest+sse)
├── hooks/                  # 通用 Composables
├── layouts/                # 页面布局
├── pages/                  # 路由页面（index/login/about/post/[id]/category/[id]/tag/[id]...）
├── plugins/                # Nuxt 插件（01.init-client.ts 等）
├── stores/                 # Pinia（modules/core + modules/app）
├── utils/ typings/
locales/                    # i18n 翻译（zh-CN/en-US）
```

## 关键约定（必须遵守）

1. **Composable 命名注意单复数** — 精确匹配（如 `useListPosts` vs `useListPost`）
2. **禁止手改 `api/generated/`** — protobuf 自动生成
3. **组件内 i18n 用 `useI18n()`** — 不要直接引用 `$t`
4. **路由导航用 `navigateTo()`** — 不用 `useRouter().push()`；**Store 中也用 `navigateTo()`**（Pinia 中无法用路由 composables）
5. **暗色样式用 CSS 变量** — `hsl(var(--foreground))`，**禁用** Tailwind `dark:` 变体
6. **图片用 `<UiImage>` 公共组件** — 自动处理加载失败占位
7. **新增组件后需重启 dev server** — 刷新 Nuxt 自动导入缓存
8. **JSON 翻译文件禁止裸花括号 `{}`** — 需用 `{'@'}` 转义特殊字符
9. **API 三形态** — 组件用 `use*` Hook；Store 用 `fetch*`；无需缓存直接调纯函数

## 开发命令

```bash
pnpm dev          # 开发服务器（.env.development）
pnpm build        # 服务端构建（.env.production）
pnpm generate     # 静态站点生成（SSG）
pnpm preview      # 预览构建结果
```

```env
# .env.development
NUXT_PUBLIC_API_BASE_URL=http://localhost:6700
NUXT_PUBLIC_ENABLE_MOCK=false
NUXT_PUBLIC_AES_KEY=
```

## 新增业务模块 Checklist（以"产品"为例）

```
- [ ] Step 1: 确认 API 类型（api/generated/ 已有 createProductServiceClient）
- [ ] Step 2: 编写 composables（api/composables/product.ts：纯函数 + use* + fetch*）
- [ ] Step 3: 注册导出（api/composables/index.ts 添加 export * from './product'）
- [ ] Step 4: 创建页面（pages/product/ 下，用 Composable 取数据）
```

## 新增 UI 组件

```bash
npx shadcn-vue@latest add dialog    # 添加 shadcn-vue 组件
```

自定义业务组件放 `components/` 对应目录，Nuxt 自动扫描注册。**注意**：只扫描 `.vue` 文件，`index.ts` 被忽略以避免命名冲突。

## 新增语言

1. 在 `locales/` 下创建新语言目录（如 `ja-JP/`）复制翻译
2. 在 `nuxt.config.ts` 的 `i18n.locales` 注册
3. 在 `core/preferences/types.ts` 的 `SupportedLanguagesType` 加类型

## 自定义主题

修改 `assets/css/main.css` 中的 CSS 变量（`--primary`/`--accent` 等），所有 `bg-primary`/`text-foreground` 类自动跟随。

## 扩展请求拦截器

```ts
const client = RequestClient.getInstance();
client.addRequestInterceptor({ fulfilled: (config) => { /* ... */ } });
client.addResponseInterceptor({ fulfilled: (response) => { /* ... */ } });
```

## 常见错误与纠正

| 错误做法 | 正确做法 |
|---|---|
| 手改 `api/generated/` | 改 proto 源后重新生成 |
| 直接引用 `$t` | 用 `useI18n()` composable |
| `useRouter().push()` 导航 | 用 `navigateTo()` |
| Store 中用 `useRouter()` | 用 `navigateTo()`（Pinia 中无路由 composables） |
| 用 Tailwind `dark:` 变体 | 用 CSS 变量 `hsl(var(--foreground))` |
| `<img>` 直接用 | 用 `<UiImage>` 公共组件 |
| Composable 单复数命名不一致 | 精确匹配（`useListPosts` vs `useListPost`） |
| JSON 翻译用裸 `{}` | 用 `{'@'}` 转义 |
| 新增组件后不重启 dev | 重启刷新自动导入缓存 |
