# 基于 Next.js 的 Headless CMS 前端架构：技术解析与二次开发导引

> 本文面向希望基于此项目进行二次开发的前端工程师，从技术栈选型、核心架构设计、关键模块实现到二开实践路径，提供一份完整的技术地图。

---

## 一、技术栈总览

本项目是一个**静态导出型** CMS 内容展示前端，采用以下核心技术栈：

| 层面 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 框架 | Next.js (App Router) | 16.x | 静态导出、路由、SSG |
| UI 库 | React | 19.x | 视图层 |
| 语言 | TypeScript | 5.x | 类型安全 |
| 样式 | Tailwind CSS | 4.x | 原子化 CSS |
| 组件库 | shadcn/ui (Radix UI) | 最新 | 无障碍 UI 原语 |
| 状态管理 | Zustand | 5.x | 轻量响应式 Store |
| 数据层 | TanStack React Query | 5.x | 服务端状态管理 |
| 国际化 | next-intl | 4.x | 多语言路由与翻译 |
| HTTP 客户端 | Axios | 1.x | REST 通信 |
| 代码高亮 | Shiki | 4.x | 双主题语法着色 |
| Markdown | marked | 17.x | 内容解析 |
| 数学公式 | KaTeX | 0.16.x | LaTeX 渲染 |
| 流程图 | Mermaid | 11.x | 图表渲染 |
| 富文本编辑 | Tiptap | 3.x | 评论编辑器 |
| 实时通信 | SSE (fetch-event-source) | 2.x | 服务端推送 |

**包管理器：** pnpm

---

## 二、核心架构设计

### 2.1 静态导出模式

项目配置为**完全静态导出**（`output: 'export'`），最终产物为纯 HTML/CSS/JS 文件，可部署到 Nginx、CDN 或对象存储：

```typescript
// next.config.ts
const nextConfig: NextConfig = {
    output: 'export',          // 静态导出
    trailingSlash: true,       // 生成 /path/index.html 目录结构
    distDir: 'dist',           // 输出到 dist 目录
};
```

`trailingSlash: true` 确保 Nginx 可通过 `try_files` 做 fallback 路由，无需额外配置 SPA fallback。

### 2.2 国际化路由架构

采用 `[locale]` 动态段实现基于 URL 前缀的多语言路由：

```
src/app/[locale]/
├── layout.tsx          # 语言布局（SSG 入口）
├── ClientLocaleLayout  # 客户端布局（Provider 注入）
├── page.tsx            # 首页
├── post/[id]/          # 文章详情
├── category/[slug]/    # 分类页
├── tag/[slug]/         # 标签页
├── login/              # 登录
├── register/           # 注册
├── settings/           # 设置
└── ...
```

**路由配置** (`routing.ts`)：

```typescript
export const routing = defineRouting({
    locales: ['zh-CN', 'en-US'],
    defaultLocale: 'zh-CN',
    localePrefix: 'always',  // URL 中始终包含语言前缀
});
```

`generateStaticParams()` 在构建时为每个语言预生成静态页面，实现完全 SSG。

**语言切换时的数据刷新**采用 `key={locale}` 强制重挂载内容区，触发所有 `useEffect` 重新加载数据，配合 `queryClient.clear()` 清除缓存。

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

二开时如需对接不同后端，只需修改 `RequestClient.init()` 的回调参数。

### 2.5 状态管理 — Zustand + React Context

采用 Zustand 配合 React Context 的混合模式，避免全局单例在 SSR 场景下的数据泄漏：

```
src/store/
├── StoreProvider.tsx           # 聚合 Provider
└── core/
    ├── access/                 # 认证凭证（accessToken、refreshToken）
    ├── user/                   # 用户信息
    └── loading/                # 全局加载状态
```

**设计要点：**

- 每个 Store 通过 `create*Store()` 工厂函数创建独立实例
- 通过 `useMemo` 确保 store 实例在组件生命周期内稳定
- Context Provider 嵌套提供 store 给子树
- `RequestClient` 通过 `store.getState()` 桥接 Context-based stores 到拦截器

### 2.6 偏好系统 — Preferences

`core/preferences` 是一个独立的偏好管理模块，管理主题模式、语言等用户偏好：

```
src/core/preferences/
├── store/          # Zustand Store
├── hooks/          # usePreferences 等 React Hook
├── components/     # 偏好相关 UI 组件
├── config/         # 默认配置
├── types/          # 类型定义
└── utils/          # 工具函数
```

主题支持三种模式：`light`（亮色）、`dark`（暗色）、`auto`（跟随系统），通过 `<html>` 上的 `.dark` 类切换。

---

## 三、关键模块深度解析

### 3.1 内容渲染管线

`ContentViewer` 组件实现了一条完整的 Markdown → HTML 渲染管线：

```
Markdown 源文
  ↓ marked（自定义 Renderer）
  ├── 代码块 → Shiki 双主题高亮
  ├── 数学公式 → KaTeX 渲染（行内 + 块级）
  ├── 流程图 → Mermaid 渲染
  ├── 表格 → 响应式容器包装
  ├── 图片 → figure/figcaption 语义化
  └── 链接 → 外部链接自动新窗口
  ↓ DOMPurify（XSS 清洗）
  ↓ 安全 HTML 输出
```

**Shiki 双主题**：使用 `github-light` / `github-dark` 主题，通过 CSS 变量 `--shiki-dark` 实现主题切换时无需重新渲染。

**安全策略**：DOMPurify 白名单严格限制允许的标签和属性，防止 XSS 攻击。

### 3.2 主题系统

基于 CSS 变量的 HSL 色板系统，亮色/暗色两套完整变量定义在 `globals.css` 中：

```css
:root {
    --primary: 142.1 76.2% 36.3%;       /* 主色 */
    --background: 210 40% 98%;           /* 背景色 */
    --card: 0 0% 100%;                   /* 卡片色 */
    --radius: 0.6rem;                    /* 全局圆角 */
    --layout-header-height: 64px;        /* 布局常量 */
    --layout-max-width: 1200px;
}

.dark {
    --primary: 142.1 86.2% 50.3%;
    --background: 224 45% 6%;
    --card: 222.2 47.4% 11%;
}
```

通过 `@theme inline` 指令将 CSS 变量映射为 Tailwind 的颜色 token（`bg-primary`、`text-foreground` 等），实现设计系统与组件的解耦。

**防闪烁**：`<head>` 中注入内联脚本（`initThemeScript`），在首帧渲染前读取 localStorage 并设置 `.dark` 类，避免主题闪烁。

### 3.3 国际化体系

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

### 3.4 认证流程

```
用户登录 → 存储 accessToken + refreshToken
    ↓
每次请求 → Token 拦截器注入 Authorization 头
    ↓
401 响应 → 自动调用 refreshToken 接口
    ↓ 成功               ↓ 失败
更新 Token 继续请求    清除凭证 → 重定向登录页
```

Token 存储在 Zustand Store 中（内存态），通过 AES 加密后持久化到 localStorage，实现「刷新页面不丢失登录态」。

---

## 四、项目目录结构与职责

```
src/
├── api/                    # API 三层架构
│   ├── generated/          #   自动生成的客户端代码
│   ├── service/            #   业务服务封装
│   └── hooks/              #   React Hook + 辅助函数
├── app/                    # Next.js App Router 页面
│   ├── globals.css         #   全局样式 + 主题变量
│   ├── layout.tsx          #   根布局（StoreProvider、ThemeProvider）
│   └── [locale]/           #   多语言路由
│       ├── layout.tsx      #     语言布局（SSG）
│       ├── ClientLocaleLayout.tsx  # 客户端布局
│       ├── routing.ts      #     路由配置
│       └── ...             #     各业务页面
├── components/             # UI 组件
│   ├── ui/                 #   shadcn/ui 基础组件
│   ├── layout/             #   布局组件（Header、Footer、Nav）
│   ├── home/               #   首页区块组件
│   ├── post/               #   文章相关组件
│   ├── category/           #   分类组件
│   ├── comment/            #   评论组件（含 Tiptap 编辑器）
│   ├── content/            #   内容渲染器（ContentViewer）
│   └── auth/               #   认证布局
├── config/                 # 环境变量配置
├── core/                   # 核心基础设施
│   ├── preferences/        #   偏好系统（主题、语言）
│   ├── storage/            #   存储抽象（localStorage 封装）
│   ├── transport/          #   通信层
│   │   ├── rest/           #     REST（RequestClient）
│   │   └── sse/            #     SSE（实时推送）
│   └── query-client.ts     #   React Query 全局配置
├── hooks/                  # 通用自定义 Hook
├── i18n/                   # 国际化配置与工具
├── lib/                    # 工具库（cn 等）
├── plugins/                # 插件（图标等）
├── store/                  # Zustand Store（Provider 模式）
└── utils/                  # 通用工具函数
```

---

## 五、二次开发导引

### 5.1 环境搭建

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建静态产物（输出到 dist/）
pnpm build

# 类型检查
pnpm lint
```

**环境变量配置**（`.env.development`）：

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:6700    # 后端 API 地址
NEXT_PUBLIC_APP_TITLE='My CMS'                    # 应用标题
NEXT_PUBLIC_DEFAULT_LOCALE=zh-CN                  # 默认语言
NEXT_PUBLIC_TOKEN_KEY=access_token                # Token 存储键名
```

### 5.2 新增一个业务页面

以「产品」模块为例：

**Step 1 — 定义 API 类型（后端 protobuf 已生成则跳过）**

如果后端使用 protobuf，运行代码生成即可。否则在 `api/generated/` 中手动定义类型。

**Step 2 — 封装服务层**

创建 `api/service/product.ts`：

```typescript
import { requestApi } from '@/core';

export async function listProducts(params: { page?: number; pageSize?: number }) {
    return requestApi.get('/api/v1/products', { params });
}

export async function getProduct(id: number) {
    return requestApi.get(`/api/v1/products/${id}`);
}
```

**Step 3 — 封装 Hook 层**

创建 `api/hooks/product.ts`：

```typescript
import { useMutation } from '@tanstack/react-query';
import { listProducts, getProduct } from '@/api/service/product';

export function useListProducts() {
    return useMutation({ mutationFn: (params) => listProducts(params) });
}

export function useGetProduct() {
    return useMutation({ mutationFn: (id: number) => getProduct(id) });
}
```

**Step 4 — 创建页面**

创建 `app/[locale]/product/[id]/page.tsx`：

```typescript
'use client';

import { useGetProduct } from '@/api/hooks/product';

export default function ProductPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = React.use(params);
    const { mutate: fetchProduct, data } = useGetProduct();
    // ...
}
```

**Step 5 — 添加路由到导航**

在 `components/layout/TopNavbar.tsx` 或 `MobileNav.tsx` 中添加导航链接。

### 5.3 新增一种语言

**Step 1 — 创建翻译文件**

在 `messages/` 下创建新语言目录（如 `ja-JP/`），复制现有 JSON 文件并翻译。

**Step 2 — 注册语言**

在 `i18n/config.ts` 中：

```typescript
export const locales = ['zh-CN', 'en-US', 'ja-JP'] as const;  // 新增
```

并导入新语言的翻译文件，添加到 `allMessages` 对象中。

**Step 3 — 完成**

由于路由基于 `[locale]` 动态段 + `generateStaticParams()`，新语言会自动在构建时生成对应的静态页面。

### 5.4 自定义主题配色

修改 `globals.css` 中的 CSS 变量即可。以替换主色为例：

```css
:root {
    --primary: 220 90% 56%;            /* 改为蓝色主色 */
    --primary-foreground: 0 0% 100%;
}

.dark {
    --primary: 220 90% 65%;
    --primary-foreground: 0 0% 100%;
}
```

所有使用 `bg-primary`、`text-primary` 等 Tailwind 类的组件会自动跟随变化。

### 5.5 替换或扩展 UI 组件

项目使用 shadcn/ui，组件源码位于 `components/ui/`，可直接修改。

**新增 shadcn/ui 组件：**

```bash
pnpm dlx shadcn@latest add dialog
```

**已有组件列表：** button、input、select、dropdown-menu、dialog、sheet、avatar、toggle、switch、separator、navigation-menu、carousel、pagination、skeleton、spinner 等。

### 5.6 对接不同后端

本项目前端与后端通过 REST API 通信，对接不同后端的核心修改点：

1. **`config/env.ts`** — 修改 `apiBaseUrl`
2. **`api/service/*.ts`** — 调整请求参数格式和响应结构
3. **`api/hooks/*.ts`** — 调整类型定义
4. **`api/generated/`** — 如后端使用 protobuf，重新生成；否则手动定义类型

认证流程可通过修改 `StoreProvider.tsx` 中 `RequestClient.init()` 的回调来自定义。

### 5.7 部署

构建后产物为纯静态文件，部署方式：

```bash
pnpm build   # 输出到 dist/
```

**Nginx 配置示例：**

```nginx
server {
    listen 80;
    root /var/www/cms;
    index index.html;

    location / {
        try_files $uri $uri/ $uri/index.html =404;
    }

    # SPA fallback for client-side routing
    location ~ ^/(zh-CN|en-US)/ {
        try_files $uri $uri/ $uri/index.html /index.html;
    }
}
```

---

## 六、开发规范与注意事项

### 6.1 客户端组件标记

Next.js App Router 下，使用 `useState`、`useEffect`、事件处理等客户端功能的组件，必须在文件顶部添加：

```typescript
'use client';
```

### 6.2 API Hooks 双形态

每个业务实体通常提供两种调用形态：

- **Hook 形态**（`use*`）— 用于 React 组件内
- **纯函数形态**（`fetch*`）— 用于 Store、事件处理等非 React 上下文

### 6.3 多语言内容辅助函数

获取后端实体的多语言字段时，使用 `hooks/` 中导出的辅助函数而非直接访问 `translations` 数组：

```typescript
import { getPostTitle, getPostSummary, getPostThumbnail } from '@/api/hooks';

const title = getPostTitle(post);       // 自动匹配当前语言
const summary = getPostSummary(post);
```

### 6.4 环境变量

所有客户端可访问的环境变量必须以 `NEXT_PUBLIC_` 前缀开头。修改 `.env` 文件后需重启开发服务器。

### 6.5 样式优先级

项目使用 Tailwind CSS v4，自定义样式优先使用 Tailwind 类名。如需自定义 CSS，在 `globals.css` 中通过 `@layer base`、`@utility` 等指令添加，确保优先级正确。

---

## 七、技术亮点总结

1. **零服务器部署**：完全静态导出，可部署到任何静态托管环境，降低运维成本
2. **类型安全的 API 层**：protobuf 自动生成 → 三层架构 → 完整 TypeScript 类型贯穿
3. **多语言全链路**：路由级国际化 + 内容级翻译 + UI 级文案，三层 i18n 完整覆盖
4. **Token 自动刷新**：内置请求排队机制，刷新期间不丢失任何请求
5. **内容渲染管线**：Markdown + 代码高亮 + 数学公式 + 流程图，一条管线处理多种内容格式
6. **主题防闪烁**：内联脚本 + CSS 变量，确保首帧即正确主题
7. **Zustand + Context**：避免全局单例的 SSR 陷阱，同时保持 Zustand 的简洁 API

---

> **快速开始**：`pnpm install && pnpm dev`，打开 `http://localhost:3000` 即可运行。
