# 基于 Taro 的 Headless CMS 多端前端架构：技术解析与二次开发导引

> 本文面向希望基于此项目进行二次开发的前端工程师，从技术栈选型、核心架构设计、关键模块实现到二开实践路径，提供一份完整的技术地图。

---

## 一、技术栈总览

本项目是一个 **Taro 多端** CMS 内容展示前端，一套代码同时编译为 H5、微信/支付宝/抖音小程序等多端应用：

| 层面       | 技术                       | 版本          | 用途                       |
|----------|--------------------------|-------------|--------------------------|
| 框架       | Taro                     | 4.1.x       | 多端统一框架                   |
| UI 库     | React                    | 18.x        | 视图层                      |
| 语言       | TypeScript               | 5.x         | 类型安全                     |
| 构建       | Vite                     | 4.x         | 编译打包                     |
| 样式       | Tailwind CSS             | 3.x         | 原子化 CSS                  |
| 小程序适配    | weapp-tailwindcss        | 4.x         | 小程序端 Tailwind rem→rpx 转换 |
| 状态管理     | Zustand                  | 5.x         | 轻量响应式 Store（Context 模式）  |
| 数据层      | TanStack React Query     | 5.x         | 服务端状态管理                  |
| 国际化      | i18next + react-i18next  | 25.x / 16.x | 多语言翻译                    |
| HTTP 客户端 | Axios                    | 1.x         | REST 通信                  |
| 代码高亮     | highlight.js             | 11.x        | 语法着色                     |
| Markdown | marked                   | 17.x        | 内容解析                     |
| 数学公式     | KaTeX                    | 0.16.x      | LaTeX 渲染                 |
| 流程图      | Mermaid                  | 11.x        | 图表渲染                     |
| 实时通信     | SSE (fetch-event-source) | 2.x         | 服务端推送                    |

**包管理器：** pnpm

---

## 二、核心架构设计

### 2.1 多端编译模式

项目基于 Taro v4，通过 `@tarojs/cli` 将同一份 React 代码编译为多端产物：

```bash
# H5 开发
pnpm dev:h5

# 微信小程序开发
pnpm dev:weapp

# 构建产物
pnpm build:h5     # → dist/ (SPA)
pnpm build:weapp  # → dist/ (小程序包)
```

**Taro 配置要点**（`config/index.ts`）：

```typescript
export default defineConfig<'vite'>(async (merge, {}) => ({
  framework: 'react',
  compiler: 'vite',
  designWidth: 750,       // 设计稿 750px，rpx 单位基准
  sourceRoot: 'src',
  outputRoot: 'dist',

  // H5 路由：browser mode + customRoutes（美化 URL）
  h5: {
    router: {
      mode: 'browser',
      customRoutes: {
        'pages/index/index': '/',           // 首页 → /
        'pages/post/detail/index': '/post/detail',
        // ... 更多路由映射
      }
    }
  },

  // 小程序端：注入 weapp-tailwindcss rem→rpx 转换
  modifyViteConfig(config) {
    if (process.env.TARO_ENV === 'weapp') {
      config.plugins.push(UnifiedViteWeappTailwindcssPlugin({ rem2rpx: true }));
    }
  },
}));
```

**customRoutes 机制**：Taro 内部页面路径如 `/pages/post/detail/index` 会映射为浏览器友好的短路径 `/post/detail`，使 H5 端 URL 更加简洁。

### 2.2 多端路由架构

Taro 使用**页面栈**路由模型，不同于 Next.js 的文件路由：

```
src/pages/
├── index/index.tsx              # 首页（pages 配置的第一项）
├── post/
│   ├── index/index.tsx          # 文章列表
│   └── detail/index.tsx         # 文章详情（接收 ?id=xxx 参数）
├── category/
│   ├── index/index.tsx          # 分类列表
│   └── detail/index.tsx         # 分类详情（接收 ?id=xxx 参数）
├── tag/
│   ├── index/index.tsx          # 标签列表
│   └── detail/index.tsx         # 标签详情（接收 ?id=xxx 参数）
├── about/index.tsx              # 关于页
├── settings/index.tsx           # 设置
├── user/index.tsx               # 用户中心
└── ...                          # 其他页面
```

**页面注册**（`app.config.ts`）：

```typescript
export default defineAppConfig({
  pages: [
    'pages/index',              // 首页（第一个为启动页）
    'pages/category/index',
    'pages/category/detail/index',
    'pages/post/index',
    'pages/post/detail/index',
    // ...
  ],
  window: {
    navigationBarBackgroundColor: '#fff',
    navigationBarTitleText: 'Taro CMS',
  }
});
```

**路由导航**（`useI18nRouter`）：

由于 Taro 不使用 locale 前缀路由（与 Next.js 版不同），提供了 `useI18nRouter` Hook 封装 Taro 路由 API：

```typescript
const router = useI18nRouter();

// 短路径自动转换为 Taro 内部路径
router.push('/');                    // → Taro.navigateTo('/pages/index/index')
router.push('/post/123');            // → Taro.navigateTo('/pages/post/detail/index?id=123')
router.push('/category/detail?id=5');// → Taro.navigateTo('/pages/category/detail/index?id=5')
router.back();                       // → Taro.navigateBack()
```

**页面刷新机制**：Taro 页面使用 `useDidShow` 生命周期（而非 `useEffect`）监听页面显示，确保从子页面返回时也能刷新数据：

```typescript
export default function Home() {
  const [refreshKey, setRefreshKey] = useState(0);

  useDidShow(() => {
    setRefreshKey(prev => prev + 1);  // 每次显示都刷新
  });

  return (
    <View>
      <LatestPostsSection key={`latest-${refreshKey}`} />
      <CategoryListSection key={`cat-${refreshKey}`} />
    </View>
  );
}
```

### 2.3 三层 API 架构

API 层遵循**生成层 → 服务层 → Hook 层**的三层分离架构：

```
src/api/
├── generated/          # [自动生成] protoc-gen-typescript-http 产出
├── service/            # [服务封装] 业务逻辑、参数转换、单例管理
├── hooks/              # [React Hook] useMutation/useQuery 封装 + 辅助函数
└── index.ts            # 统一导出
```

**第一层 — 自动生成的客户端**（`generated/`）：由 protobuf 定义自动生成的 HTTP 客户端代码，不应手动编辑。

**第二层 — 服务封装**（`service/`）：基于生成的客户端封装业务方法，注入 locale、分页参数等：

```typescript
// service/post.ts
export async function listPostsRaw(params) {
    const locale = currentLocaleLanguageCode();
    const formValues = {...(params.formValues || {}), locale};
    return getPostService().List({
        query: JSON.stringify(formValues),
        page: params.paging?.page,
        pageSize: params.paging?.pageSize,
    });
}
```

**第三层 — React Hook**（`hooks/`）：封装为 `useMutation` / `useQuery` Hook 和纯函数 `fetch*` 两种形态：

- `useListPosts()` — 组件内使用的 React Hook
- `fetchListPosts()` — Store / 非 React 上下文中使用的纯异步函数

每层职责清晰，二开时只需关注 service 和 hooks 层。

### 2.4 RequestClient — HTTP 通信内核

`RequestClient` 是基于 Axios 封装的全局单例，通过拦截器链实现完整的认证生命周期：

```
请求拦截链：Token 注入 → Request-ID 注入 → Locale 注入
响应拦截链：数据解构 → 401 认证处理 → 错误消息提取
```

**关键特性：**

- **Token 自动刷新**：401 时自动调用 refresh token 接口，刷新期间后续请求排队等待
- **请求重认证**：刷新失败时清除凭证并重定向至登录页
- **语言感知**：自动注入 `Accept-Language` 头

```typescript
// 初始化（在 StoreProvider 中执行）
RequestClient.init(env.apiBaseUrl, {
  getToken: () => accessStore.getState().accessToken?.value,
  getLocale: () => preferencesStore.getState().preferences.app.locale,
  refreshToken: async () => { /* ... */ },
  onReAuthenticate: async (redirect) => { /* ... */ },
});
```

二开时如需对接不同后端，只需修改 `StoreProvider` 中 `RequestClient.init()` 的回调参数。

### 2.5 状态管理 — Zustand + React Context

采用 Zustand 配合 React Context 的混合模式，避免全局单例在多实例场景下的数据泄漏：

```
src/store/
├── StoreProvider.tsx           # 聚合 Provider（创建所有 store 实例）
└── core/
    ├── access/                 # 认证凭证（accessToken、refreshToken）
    ├── user/                   # 用户信息
    └── loading/                # 全局加载状态
```

**设计要点：**

- 每个 Store 通过 `create*Store()` 工厂函数创建独立实例
- 通过 `useMemo` 确保 store 实例在组件生命周期内稳定
- Context Provider 嵌套提供 store 给子树
- `RequestClient` 通过 `storeApi.getState()` 桥接 Context-based stores 到拦截器

```typescript
// StoreProvider.tsx — 聚合所有 store 的 Provider
export default function StoreProvider({children}) {
  const accessStore = useMemo(() => createAccessStore(), []);
  const userStore = useMemo(() => createUserStore(), []);
  const loadingStore = useMemo(() => createLoadingStore(), []);
  const preferencesStore = useMemo(() => createPreferencesStore(), []);

  // 初始化 RequestClient 拦截器
  useMemo(() => { /* ... */ }, [accessStore, userStore, preferencesStore]);

  return (
    <PreferencesStoreContext.Provider value={preferencesStore}>
      <AccessStoreContext.Provider value={accessStore}>
        <UserStoreContext.Provider value={userStore}>
          <LoadingStoreContext.Provider value={loadingStore}>
            {children}
          </LoadingStoreContext.Provider>
        </UserStoreContext.Provider>
      </AccessStoreContext.Provider>
    </PreferencesStoreContext.Provider>
  );
}
```

### 2.6 偏好系统 — Preferences

`core/preferences` 是一个独立的偏好管理模块，管理主题模式、语言、内容偏好、功能部件等用户配置：

```
src/core/preferences/
├── index.ts         # 统一导出（types + config + hooks + store）
├── types/           # 类型定义（6 个分组：app/theme/content/copyright/widget/transition）
├── config/          # 默认配置（defaultPreferences）
├── store/           # Zustand Store（createPreferencesStore + Context）
├── hooks/           # usePreferences / useThemeConfig / usePreferencesLocale
└── utils/           # mergeDeep 深度合并工具
```

主题支持三种模式：`light`（亮色）、`dark`（暗色）、`auto`（跟随系统），通过内联 CSS 变量注入色板。

### 2.7 国际化体系

Taro 版使用 **i18next + react-i18next**（不同于 Next.js 版的 next-intl），不使用 URL 语言前缀路由：

```typescript
// i18n/config.ts
export const locales = ['zh-CN', 'en-US'] as const;
export const defaultLocale = 'zh-CN';

// 翻译文件扁平化后注入 i18next
export function getFlattenedMessages(): Record<Locale, Record<string, string>> {
  // 将 messages/zh-CN/*.json 所有文件合并为扁平键值对
  // 如 { page: { home: { title: '首页' } } } → { 'page.home.title': '首页' }
}
```

**翻译文件结构：**

```
messages/
├── zh-CN/
│   ├── app.json           # 应用级文案
│   ├── navbar.json         # 导航栏
│   ├── page.json           # 页面文案
│   ├── cms.json            # CMS 业务文案
│   ├── authentication.json # 认证相关
│   ├── comment.json        # 评论
│   ├── enum.json           # 枚举翻译
│   ├── settings.json       # 设置
│   └── ...
└── en-US/
    └── ...（同构文件）
```

**多语言内容获取**：后端返回的实体（Post、Category 等）携带 `translations[]` 数组，前端通过辅助函数按当前 locale 提取对应翻译：

```typescript
export function getPostTitle(post: contentservicev1_Post): string {
    const translation = getTranslation(post);  // 匹配当前语言
    return translation?.title || '';
}
```

---

## 三、关键模块深度解析

### 3.1 内容渲染管线

`ContentViewer` 组件实现了一条完整的 Markdown → HTML 渲染管线：

```
Markdown 源文
  ↓ marked（自定义 Renderer）
  ├── 代码块 → highlight.js 语法高亮
  ├── 数学公式 → KaTeX 渲染（行内 + 块级）
  ├── 流程图 → Mermaid 渲染
  ├── 表格 → 响应式容器包装
  ├── 图片 → figure/figcaption 语义化
  └── 链接 → 外部链接自动新窗口
  ↓ DOMPurify（XSS 清洗）
  ↓ 安全 HTML 输出
```

**安全策略**：DOMPurify 白名单严格限制允许的标签和属性，防止 XSS 攻击。

### 3.2 主题系统

采用**内联 CSS 变量**注入方案（而非 CSS 文件定义），由 `ThemeClientProvider` 组件根据 `isDark` 状态动态切换色板：

```typescript
// ThemeClientProvider.tsx
const themeStyle = useMemo(() => {
  if (isDark) {
    return {
      '--color-text-main': '#ffffffe6',
      '--color-text-sec': '#ffffffcc',
      '--color-page-bg': '#17171a',
      '--color-card-bg': '#232326',
      '--color-split-line': '#3a3a3c',
      '--color-bar-bg': 'rgba(35, 35, 38, 0.92)',
    };
  }
  return {
    '--color-text-main': '#1d2129',
    '--color-text-sec': '#4e5969',
    '--color-page-bg': '#ebedf0',
    '--color-card-bg': '#ffffff',
    '--color-split-line': '#e8eaed',
    '--color-bar-bg': 'rgba(255, 255, 255, 0.92)',
  };
}, [isDark]);

return <View style={themeStyle}>{children}</View>;
```

**设计要点：**

- **全平台兼容**：内联 CSS 变量方案在 H5 和小程序端均可靠工作
- **暗色模式检测**：H5 端监听 `window.matchMedia('(prefers-color-scheme: dark)')`
- **亮色/暗色类切换**：H5 端在 `<html>` 上切换 `.dark` / `.light` class，小程序端使用 `page.dark` 选择器
- **CSS 变量与 Tailwind 映射**：`--color-text-main` → `text-textMain`，`--color-page-bg` → `bg-pageBg`

### 3.3 布局系统

`Layout` 组件是所有页面的外层容器，统一提供 Header、Footer、BackToTop、NavigationProgress：

```typescript
// Layout.tsx
function Layout({children}) {
  return (
    <ThemeClientProvider>
      <MobileNavProvider>
        <NavigationProgress />
        <Header />
        <View className='min-h-screen pt-[88rpx]'>
          {children}
        </View>
        <Footer />
        <BackToTop />
      </MobileNavProvider>
    </ThemeClientProvider>
  );
}
```

**应用入口**（`app.ts`）：

```typescript
function App({children}) {
  useLaunch(() => { console.log('App launched.'); });
  return createElement(
    StoreProvider, null,
    createElement(Layout, null, children)
  );
}
```

层级关系：`App → StoreProvider → Layout (ThemeClientProvider → MobileNavProvider → Header/Footer/Content)`

### 3.4 认证流程

```
用户登录 → 存储 accessToken + refreshToken（Zustand Store + AES 加密持久化）
    ↓
每次请求 → Token 拦截器注入 Authorization 头
    ↓
401 响应 → 自动调用 refreshToken 接口
    ↓ 成功               ↓ 失败
更新 Token 继续请求    清除凭证 → Taro.redirectTo 用户页
```

Token 存储在 Zustand Store 中（内存态），通过 AES 加密后持久化到 localStorage，实现「刷新页面不丢失登录态」。

### 3.5 评论系统

评论模块支持多级嵌套回复，包含：

- **评论列表**（`CommentTree`）：树形结构展示，博主回复标签高亮
- **评论输入**（`CommentSection`）：支持回复指定评论，暗黑模式适配
- **SSE 实时推送**：新评论通过 Server-Sent Events 实时通知

---

## 四、项目目录结构与职责

```
src/
├── api/                    # API 三层架构
│   ├── generated/          #   自动生成的客户端代码
│   ├── service/            #   业务服务封装
│   └── hooks/              #   React Hook + 辅助函数
├── app.config.ts           # Taro 页面注册 + 全局窗口配置
├── app.ts                  # 应用入口（StoreProvider + Layout）
├── assets/                 # 静态资源
│   ├── fonts/              #   字体文件
│   ├── icons/              #   图标
│   └── images/             #   图片（logo 等）
├── components/             # UI 组件
│   ├── layout/             #   布局组件（Header、Footer、MobileNav、BackToTop、ThemeClientProvider）
│   ├── category/           #   分类组件（Card、List、Filter、Tree）
│   ├── comment/            #   评论组件（CommentSection、CommentTree）
│   ├── content/            #   内容渲染器（ContentViewer）
│   ├── post/               #   文章组件（PostCard、PostList）
│   ├── Pagination/         #   分页组件
│   ├── otp-input/          #   验证码输入组件
│   └── ui/                 #   基础 UI 组件（AppEmpty、Button、Input）
├── config/                 # 环境变量配置（env.ts）
├── core/                   # 核心基础设施
│   ├── preferences/        #   偏好系统（主题、语言、内容偏好、部件开关）
│   ├── storage/            #   存储抽象（localStorage 封装 + AES 加密）
│   ├── transport/          #   通信层
│   │   ├── rest/           #     REST（RequestClient）
│   │   └── sse/            #     SSE（实时推送）
│   └── query-client.ts     #   React Query 全局配置
├── hooks/                  # 通用自定义 Hook
├── i18n/                   # 国际化配置与工具
│   ├── config.ts           #   语言定义 + 翻译资源
│   ├── index.ts            #   i18next 初始化
│   └── helpers/            #   useI18nRouter 路由 Hook
├── lib/                    # 工具库（next-intl-compat、cn 等）
├── pages/                  # Taro 页面
│   ├── index/              #   首页
│   ├── post/               #   文章（列表 + 详情）
│   ├── category/           #   分类（列表 + 详情）
│   ├── tag/                #   标签（列表 + 详情）
│   ├── about/              #   关于
│   ├── login/              #   登录
│   ├── register/           #   注册
│   ├── settings/           #   设置
│   ├── user/               #   用户中心
│   └── ...                 #   其他页面
├── plugins/                # 插件（XIcon 图标组件）
├── store/                  # Zustand Store（Provider 模式）
│   ├── StoreProvider.tsx   #   聚合 Provider
│   └── core/               #   access / user / loading
└── utils/                  # 通用工具函数
```

---

## 五、二次开发导引

### 5.1 环境搭建

```bash
# 安装依赖
pnpm install

# 启动 H5 开发服务器
pnpm dev:h5

# 启动微信小程序开发
pnpm dev:weapp

# 构建 H5 生产产物（输出到 dist/）
pnpm build:h5

# 构建微信小程序生产产物
pnpm build:weapp

# 类型检查
pnpm lint
```

**环境变量配置**（`.env.development`）：

```env
TARO_APP_NAMESPACE=gowind-cms-app            # 应用命名空间
TARO_APP_API_BASE_URL=https://api.cms.gowind.cloud  # 后端 API 地址
TARO_APP_AES_KEY="f51d66a73d8a0927"           # AES 加密密钥
TARO_APP_TITLE='GoWind Content Hub'           # 应用标题
TARO_DEFAULT_LOCALE=zh-CN                     # 默认语言
TARO_ENABLE_MOCK=false                        # 是否启用 Mock
```

> 环境变量通过 `config/index.ts` 的 `defineConstants` 注入到客户端代码，在 `src/config/env.ts` 中统一导出。

### 5.2 新增一个业务页面

以「产品」模块为例：

**Step 1 — 封装服务层**

创建 `api/service/product.ts`：

```typescript
import { requestApi } from '@/core';

export async function listProducts(params: { page?: number; pageSize?: number }) {
    return requestApi.get('/api/v1/products', { params });
}
```

**Step 2 — 封装 Hook 层**

创建 `api/hooks/product.ts`：

```typescript
import { useMutation } from '@tanstack/react-query';
import { listProducts } from '@/api/service/product';

export function useListProducts() {
    return useMutation({ mutationFn: (params) => listProducts(params) });
}
```

**Step 3 — 创建页面**

创建 `src/pages/product/index.tsx`：

```typescript
import {View} from '@tarojs/components';
import {useListProducts} from '@/api/hooks/product';

export default function ProductPage() {
    const {mutate: fetchProducts, data} = useListProducts();
    // ...
}
```

**Step 4 — 注册页面路由**

在 `app.config.ts` 中添加：

```typescript
export default defineAppConfig({
  pages: [
    'pages/index',
    // ...
    'pages/product/index',    // 新增
  ],
});
```

在 `config/index.ts` 的 `h5.router.customRoutes` 中添加 H5 美化路由：

```typescript
customRoutes: {
  // ...
  'pages/product/index': '/product',
}
```

**Step 5 — 添加导航入口**

在 `components/layout/MobileNav.tsx` 的导航菜单中添加产品页入口。

### 5.3 新增一种语言

**Step 1 — 创建翻译文件**

在 `messages/` 下创建新语言目录（如 `ja-JP/`），复制现有 JSON 文件并翻译。

**Step 2 — 注册语言**

在 `i18n/config.ts` 中：

```typescript
// 1. 导入翻译文件
import jaJP_app from '../../messages/ja-JP/app.json';
// ... 其他文件

// 2. 添加到 locales 数组
export const locales = ['zh-CN', 'en-US', 'ja-JP'] as const;

// 3. 添加到 allMessages 对象
export const allMessages = {
  'zh-CN': { /* ... */ },
  'en-US': { /* ... */ },
  'ja-JP': { app: jaJP_app, /* ... */ },
};
```

**Step 3 — 更新类型定义**

在 `core/preferences/types/layout.ts` 中扩展 `SupportedLanguagesType`：

```typescript
export type SupportedLanguagesType = 'zh-CN' | 'en-US' | 'ja-JP';
```

### 5.4 自定义主题配色

修改 `ThemeClientProvider.tsx` 中的 CSS 变量色板即可：

```typescript
// 亮色色板
const lightTheme = {
  '--color-text-main': '#1d2129',
  '--color-page-bg': '#ebedf0',
  '--color-card-bg': '#ffffff',
  '--color-split-line': '#e8eaed',
  // ...
};

// 暗色色板
const darkTheme = {
  '--color-text-main': '#ffffffe6',
  '--color-page-bg': '#17171a',
  '--color-card-bg': '#232326',
  '--color-split-line': '#3a3a3c',
  // ...
};
```

所有使用 `bg-pageBg`、`bg-cardBg`、`text-textMain` 等 Tailwind 类的组件会自动跟随变化。

### 5.5 对接不同后端

本项目前端与后端通过 REST API 通信，对接不同后端的核心修改点：

1. **`config/env.ts`** / **`.env.development`** — 修改 `TARO_APP_API_BASE_URL`
2. **`api/service/*.ts`** — 调整请求参数格式和响应结构
3. **`api/hooks/*.ts`** — 调整类型定义
4. **`api/generated/`** — 如后端使用 protobuf，重新生成；否则手动定义类型

认证流程可通过修改 `StoreProvider.tsx` 中 `RequestClient.init()` 的回调来自定义。

### 5.6 部署

**H5 部署：**

```bash
pnpm build:h5   # 输出到 dist/
```

产物为 SPA 静态文件，需配置服务器 fallback 到 `index.html`。

**Nginx 配置示例：**

```nginx
server {
    listen 80;
    root /var/www/cms;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

**微信小程序部署：**

```bash
pnpm build:weapp   # 输出到 dist/
```

将 `dist/` 目录上传到微信开发者工具即可预览和发布。

---

## 六、开发规范与注意事项

### 6.1 Taro 组件限制

Taro 使用 `<View>`、`<Text>`、`<Image>` 等跨端组件替代 HTML 标签：

```tsx
// ✅ 正确（跨端兼容）
<View className='flex items-center'>
  <Text className='text-textMain'>标题</Text>
  <Image src={logo} mode='aspectFit' />
</View>

// ❌ 错误（小程序不支持）
<div className='flex items-center'>
  <span>标题</span>
  <img src={logo} />
</div>
```

### 6.2 rpx 单位

设计稿基于 750px 宽度，使用 `rpx` 单位确保多端自适应：

```tsx
<View className='px-[24rpx] pt-[32rpx]'>  {/* 内边距使用 rpx */}
  <Text className='text-[28rpx]'>文字</Text>  {/* 字号使用 rpx */}
</View>
```

Tailwind CSS 中使用方括号语法（如 `px-[24rpx]`）指定 rpx 值。

### 6.3 API Hooks 双形态

每个业务实体通常提供两种调用形态：

- **Hook 形态**（`use*`）— 用于 React 组件内
- **纯函数形态**（`fetch*`）— 用于 Store、事件处理等非 React 上下文

### 6.4 多语言内容辅助函数

获取后端实体的多语言字段时，使用 `hooks/` 中导出的辅助函数而非直接访问 `translations` 数组：

```typescript
import { getPostTitle, getPostSummary, getPostThumbnail } from '@/api/hooks';

const title = getPostTitle(post);       // 自动匹配当前语言
const summary = getPostSummary(post);
```

### 6.5 环境变量

Taro 环境变量必须以 `TARO_APP_` 前缀开头（Taro 规范），通过 `config/index.ts` 的 `defineConstants` 注入。修改 `.env` 文件后需重启开发服务器。

### 6.6 样式优先级

项目使用 Tailwind CSS + rpx 单位，自定义样式优先使用 Tailwind 类名。全局 CSS 样式在 `app.css` 中定义。

**H5 端特殊处理**：Input / Textarea 组件的背景色和文字颜色需要通过全局 CSS `!important` 覆盖 Taro 默认样式，以确保暗黑模式正确显示。

### 6.7 路由注意事项

- 使用 `useI18nRouter()` Hook 进行路由跳转，它会自动处理短路径到 Taro 内部路径的转换
- H5 端回首页（`reLaunch`）使用 `window.location.href` 强制浏览器原生导航，因为 Taro 的 SPA 路由只做 `history.pushState` 不触发页面重渲染

---

## 七、技术亮点总结

1. **一套代码多端运行**：Taro 框架实现 H5 + 微信/支付宝/抖音小程序的统一代码库
2. **类型安全的 API 层**：protobuf 自动生成 → 三层架构 → 完整 TypeScript 类型贯穿
3. **Tailwind CSS 多端适配**：通过 weapp-tailwindcss 插件实现 rem→rpx 自动转换
4. **内联 CSS 变量主题系统**：全平台兼容的暗色/亮色/自动主题切换
5. **Token 自动刷新**：内置请求排队机制，刷新期间不丢失任何请求
6. **内容渲染管线**：Markdown + 代码高亮 + 数学公式 + 流程图，一条管线处理多种内容格式
7. **Zustand + Context**：工厂函数模式避免全局单例，Provider 嵌套保证 store 隔离
8. **i18next 国际化**：翻译文件扁平化注入，语言切换无需 URL 前缀

---

> **快速开始**：`pnpm install && pnpm dev:h5`，打开 `http://localhost:10086` 即可运行。
