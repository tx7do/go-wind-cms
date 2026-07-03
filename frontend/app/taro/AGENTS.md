# AGENTS.md — Taro 多端前端开发指南

> 本文件是 `frontend/app/taro` 子项目的 AI 编码规范单一事实源，适用于所有支持 AGENTS.md 的 AI 编码工具。Claude Code 通过同级 `CLAUDE.md` 中的 `@AGENTS.md` 引用加载。

## 项目概览

基于 **Taro** 的多端统一 CMS 内容展示前端，一套代码同时编译为 H5、微信/支付宝/抖音小程序。

**核心技术栈**：Taro 4.1 + React 18 + TypeScript + Vite + Tailwind CSS 3 + weapp-tailwindcss（小程序 rem→rpx）+ Zustand + TanStack React Query + i18next/react-i18next + Axios + highlight.js/marked/KaTeX/Mermaid

**包管理器**：pnpm

## 关键架构认知

### 多端编译模式

```bash
pnpm dev:h5       # H5 开发 → http://localhost:10086
pnpm dev:weapp    # 微信小程序开发
pnpm build:h5     # → dist/ (SPA)
pnpm build:weapp  # → dist/ (小程序包)
```

**Taro 配置要点**（`config/index.ts`）：

```typescript
export default defineConfig<'vite'>(async (merge, {}) => ({
  framework: 'react',
  compiler: 'vite',
  designWidth: 750,       // 设计稿 750px，rpx 基准
  h5: {
    router: {
      mode: 'browser',
      customRoutes: {     // H5 美化路由
        'pages/index/index': '/',
        'pages/post/detail/index': '/post/detail',
      }
    }
  },
  // 小程序端注入 weapp-tailwindcss rem→rpx 转换
  modifyViteConfig(config) {
    if (process.env.TARO_ENV === 'weapp') {
      config.plugins.push(UnifiedViteWeappTailwindcssPlugin({ rem2rpx: true }));
    }
  },
}));
```

### 多端路由（页面栈模型）

Taro 使用**页面栈**路由模型（非文件路由），页面注册在 `app.config.ts`：

```
src/pages/
├── index/index.tsx              # 首页（pages 第一项 = 启动页）
├── post/
│   ├── index/index.tsx          # 文章列表
│   └── detail/index.tsx         # 文章详情（?id=xxx）
├── category/ tag/ about/ login/ settings/ user/
```

```typescript
// app.config.ts
export default defineAppConfig({
  pages: ['pages/index', 'pages/category/index', 'pages/post/detail/index', /* ... */],
});
```

**路由导航**用 `useI18nRouter()` Hook（封装 Taro 路由 API，自动短路径→内部路径转换）：

```typescript
const router = useI18nRouter();
router.push('/');                    // → Taro.navigateTo('/pages/index/index')
router.push('/post/123');            // → Taro.navigateTo('/pages/post/detail/index?id=123')
router.back();                       // → Taro.navigateBack()
```

**页面刷新**用 `useDidShow` 生命周期（非 `useEffect`），确保子页面返回时也刷新：

```typescript
useDidShow(() => { setRefreshKey(prev => prev + 1); });
```

### 两层 API 架构

```
src/api/
├── client.ts            # ApiClient 单例（createApiClient(transport)，懒加载各 Service）
├── generated/           # [自动生成] protoc-gen-typescript-http（禁止手改）
├── hooks/               # [业务封装 + Hook] 直接用 apiClient.xxxService
└── index.ts
```

不同于 React 版的三层，Taro 版简化为**两层**（无独立 service 层）：hooks 文件内直接调用 `apiClient`，封装业务逻辑（locale 注入、分页转换）+ `useMutation`/`useQuery` + `fetch*` 纯函数。

```typescript
// hooks/product.ts
import { apiClient } from '@/api/client';

export async function listProductsRaw(params) {
  return apiClient.productService.List({ page: params.paging?.page, pageSize: params.paging?.pageSize });
}
export function useListProducts() {
  return useMutation({ mutationFn: (params) => listProductsRaw(params) });
}
export async function fetchListProducts(params) {
  return queryClient.fetchQuery({ queryKey: ['listProducts', params], queryFn: () => listProductsRaw(params), retry: 0 });
}
```

### RequestClient / 状态管理 / 偏好系统

与 React 版一致：Axios 单例（Token 自动刷新排队）+ Zustand Context 模式 + `core/preferences` 偏好模块。差别：Taro 版用 **i18next**（非 next-intl），**不使用** URL 语言前缀路由。

### 主题系统（内联 CSS 变量）

采用**内联 CSS 变量**注入（非 CSS 文件），由 `ThemeClientProvider` 根据 `isDark` 动态切换色板，全平台兼容：

```typescript
const themeStyle = useMemo(() => isDark ? {
  '--color-text-main': '#ffffffe6', '--color-page-bg': '#17171a', /* ... */
} : {
  '--color-text-main': '#1d2129', '--color-page-bg': '#ebedf0', /* ... */
}, [isDark]);
return <View style={themeStyle}>{children}</View>;
```

CSS 变量与 Tailwind 映射：`--color-text-main` → `text-textMain`，`--color-page-bg` → `bg-pageBg`。

## 目录结构

```
src/
├── api/                    # 两层架构（client + generated + hooks）
├── app.config.ts           # Taro 页面注册 + 全局窗口配置
├── app.ts                  # 应用入口（StoreProvider + Layout）
├── components/             # layout/category/comment/content/post/Pagination/ui
├── config/                 # 环境变量（env.ts）
├── core/                   # preferences / storage / transport(rest+sse) / query-client
├── hooks/ i18n/ lib/ pages/ plugins/ store/ utils/
└── messages/               # i18n 翻译（zh-CN/en-US: app/navbar/page/cms/authentication/comment/enum/settings.json）
```

## 关键约定（必须遵守）

1. **必须用 Taro 跨端组件** — `<View>`/`<Text>`/`<Image>` 替代 `<div>`/`<span>`/`<img>`（小程序不支持 HTML 标签）
2. **尺寸用 rpx 单位** — 设计稿基于 750px，Tailwind 中用方括号语法 `px-[24rpx]`/`text-[28rpx]`
3. **禁止手改 `api/generated/`** — protobuf 自动生成
4. **API 双形态** — 组件内用 `use*`；Store/非 React 上下文用 `fetch*`
5. **多语言内容用辅助函数** — `getPostTitle(post)` 等，不直接访问 `translations`
6. **环境变量加 `TARO_APP_` 前缀** — Taro 规范，通过 `config/index.ts` 的 `defineConstants` 注入，改 `.env` 后需重启
7. **页面刷新用 `useDidShow`** — 不用 `useEffect`（确保子页面返回时刷新）
8. **路由用 `useI18nRouter()`** — 自动处理短路径到 Taro 内部路径转换

## 开发命令与环境变量

```bash
pnpm install
pnpm dev:h5        # H5 开发
pnpm dev:weapp     # 微信小程序
pnpm build:h5      # H5 生产 → dist/
pnpm build:weapp   # 小程序生产 → dist/
pnpm lint
```

```env
# .env.development
TARO_APP_NAMESPACE=gowind-cms-app
TARO_APP_API_BASE_URL=https://api.cms.gowind.cloud
TARO_APP_AES_KEY="f51d66a73d8a0927"
TARO_APP_TITLE='GoWind Content Hub'
TARO_DEFAULT_LOCALE=zh-CN
TARO_ENABLE_MOCK=false
```

## 新增业务模块 Checklist（以"产品"为例）

```
- [ ] Step 1: 封装 Hook 层（api/hooks/product.ts：业务函数 + use* + fetch*）
- [ ] Step 2: 创建页面（src/pages/product/index.tsx，用 View/Text 跨端组件）
- [ ] Step 3: 注册页面路由（app.config.ts 的 pages 数组）
- [ ] Step 4: H5 美化路由（config/index.ts 的 h5.router.customRoutes）
- [ ] Step 5: 添加导航入口（components/layout/MobileNav.tsx）
```

## 新增语言

1. 在 `messages/` 下创建新语言目录（如 `ja-JP/`），复制翻译
2. 在 `i18n/config.ts` 导入翻译、加入 `locales` 数组和 `allMessages` 对象
3. 在 `core/preferences/types/layout.ts` 扩展 `SupportedLanguagesType`

## 部署

**H5**：`pnpm build:h5` → SPA 静态文件，Nginx 配 `try_files $uri $uri/ /index.html`
**微信小程序**：`pnpm build:weapp` → `dist/` 上传到微信开发者工具

## 常见错误与纠正

| 错误做法 | 正确做法 |
|---|---|
| `<div>`/`<span>`/`<img>` | `<View>`/`<Text>`/`<Image>`（跨端组件） |
| 用 `px`/`rem` 固定尺寸 | 用 `rpx`（`px-[24rpx]`） |
| 手改 `api/generated/` | 改 proto 源后重新生成 |
| 页面用 `useEffect` 监听显示刷新 | 用 `useDidShow` 生命周期 |
| 用 `Taro.navigateTo` 手拼路径 | 用 `useI18nRouter()` 短路径 |
| 直接访问 `post.translations` | 用 `getPostTitle(post)` 辅助函数 |
| 环境变量不加 `TARO_APP_` 前缀 | 必须加前缀 |
| H5 回首页用 Taro SPA 路由 | 用 `window.location.href` 强制原生导航 |
