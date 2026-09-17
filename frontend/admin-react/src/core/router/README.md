# 路由系统模块 (Router)

## 概述

本模块提供了基于 React Router v6 的增强路由解决方案，支持**前端静态路由**和**后端动态路由**两种模式，集成了权限过滤、菜单生成、自动重定向、懒加载等能力。

### 核心特性

- **双模式路由生成**：前端模式（静态路由 + 权限过滤）和后端模式（API 拉取 + 组件映射）
- **权限集成**：路由级别权限控制，通过 `meta.authority` 字段与用户权限码匹配
- **自动菜单生成**：从路由配置自动生成侧边栏菜单树
- **懒加载**：`createLazyRoute` 工具函数，自动包裹 Suspense
- **路由守卫**：`AuthGuard`（需登录）和 `GuestGuard`（仅未登录）
- **元数据传递**：`meta` 自动转换为 React Router 的 `handle`，组件中通过 `useMatches` 获取

---

## 模块结构

```
src/core/router/                    # 核心路由引擎（通用、可复用）
├── index.ts                        # 统一导出
├── factory.ts                      # 路由工厂：createAccessibleRouter
├── types/index.ts                  # 所有类型定义（AppRouteObject、RouteMeta、AppMenu 等）
├── generators/                     # 路由/菜单生成器
│   ├── generate-routes-frontend.ts # 前端模式：静态路由 + 权限过滤
│   ├── generate-routes-backend.ts  # 后端模式：API 路由 + 组件映射
│   └── generate-menus.ts           # 菜单树生成
└── utils/                          # 工具函数
    ├── inject-redirect.ts          # 自动注入重定向
    ├── loader.ts                   # 模块加载与路由过滤
    ├── menu.ts                     # 路由转 ProLayout 菜单格式
    ├── merge-route-modules.ts      # 合并 Vite glob 路由模块
    ├── reset-routes.ts             # 重置可访问路由（白名单机制）
    ├── sort-routes.ts              # 按 meta.order 排序路由树
    └── transform-meta-to-handle.ts # meta → handle 转换

src/router/                         # 业务路由配置（项目特定）
├── index.tsx                       # AppRouter 组件（路由初始化入口）
├── config/                         # 基础路由配置
│   ├── static.tsx                  # 静态路由（主布局 + 根路径重定向）
│   ├── auth.tsx                    # 认证路由（登录/注册，不受 AuthGuard 保护）
│   └── error-routes.tsx            # 错误页路由（401/403/500/404 等）
├── guards/                         # 路由守卫
│   ├── AuthGuard.tsx               # 认证守卫：未登录 → 跳转登录页
│   └── GuestGuard.tsx              # 访客守卫：已登录 → 跳转首页
├── modules/                        # 业务路由模块（自动导入）
│   ├── dashboard.tsx               # 仪表盘
│   ├── system.tsx                  # 系统管理
│   ├── permissions.tsx             # 权限管理
│   ├── opm.tsx                     # 组织架构管理
│   ├── tenant.tsx                  # 租户管理
│   ├── log.tsx                     # 审计日志
│   └── internal-message.tsx        # 内部消息
└── utils/
    └── lazy.tsx                    # createLazyRoute 工具
```

### 导入方式

```typescript
// 从 core/router 导入类型和工具
import type { AppRouteObject, RouteMeta, AppMenu } from '@/core/router';
import { createAccessibleRouter } from '@/core/router';
import { generateMenus } from '@/core/router';
```

---

## 快速开始

### 新增一个业务页面

假设要新增 `/system/users` 用户管理页面：

**第一步**：创建页面组件 `src/pages/app/system/users/index.tsx`

**第二步**：在 `src/router/modules/system.tsx` 中添加路由配置：

```tsx
import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/router/utils/lazy';

// 在 systemRoutes 的 children 数组中添加：
{
  name: 'users',                                  // 唯一标识（必须）
  path: 'users',                                  // 相对路径，最终为 /system/users
  element: createLazyRoute(() => import('@/pages/app/system/users')),
  meta: {
    title: 'routes:users',                        // i18n 翻译键
    icon: 'lucide:users',                         // 图标（Iconify 格式）
    order: 2,                                     // 菜单排序（越小越靠前）
    authority: ['sys:user:view'],                 // 权限码（可选）
  },
},
```

**第三步**：在 `src/locales/zh-CN/_modules/routes.json` 中添加翻译：

```json
{
  "users": "用户管理"
}
```

> 添加路由后无需额外注册，`modules/` 目录下的文件会通过 `import.meta.glob` 自动导入并合并到主路由中。

---

## 路由配置详解

### AppRouteObject

项目扩展的路由对象，在 React Router 原生 `RouteObject` 基础上增加了业务字段：

```typescript
interface AppRouteObject extends Omit<RouteObject, 'children'> {
  name?: string;           // 路由唯一标识（权限匹配关键，index 路由可省略）
  path?: string;           // 路由路径（相对路径）
  element?: ReactNode;     // 渲染的 React 元素
  redirect?: string;       // 重定向目标
  meta?: RouteMeta;        // 业务元数据
  children?: AppRouteObject[];  // 子路由
  componentPath?: string;  // 后端模式：组件路径字符串
  parent?: string;         // 父级路径
  parents?: string[];      // 所有父级路径（面包屑用）
  order?: number;          // 排序
  disabled?: boolean;      // 是否禁用
  show?: boolean;          // 是否显示
}
```

### RouteMeta 元数据

`meta` 是路由配置的核心扩展字段，控制菜单、标签页、权限、面包屑等行为：

```typescript
interface RouteMeta {
  // === 基础信息 ===
  title?: string;             // 菜单/标签页/面包屑标题（支持 i18n 键名，如 'routes:dashboard'）
  icon?: IconType;            // 菜单图标
  activeIcon?: IconType;      // 激活时图标

  // === 权限控制 ===
  authority?: string[];       // 权限码数组（角色码和权限码混合，任一匹配即通过）
  ignoreAccess?: boolean;     // 忽略权限检查（公开访问）

  // === 显示控制 ===
  hideInMenu?: boolean;       // 在侧边菜单中隐藏
  hideInTab?: boolean;        // 在标签页中隐藏
  hideInBreadcrumb?: boolean; // 在面包屑中隐藏
  hideChildrenInMenu?: boolean; // 子路由不在菜单中展开

  // === 特殊行为 ===
  menuVisibleWithForbidden?: boolean;  // 菜单可见但无权限时显示 403
  activePath?: string;        // 手动指定菜单高亮路径
  affixTab?: boolean;         // 固定标签页（不可关闭）
  affixTabOrder?: number;     // 固定标签页排序
  keepAlive?: boolean;        // 页面缓存（预留）

  // === 外部链接/iframe ===
  link?: string;              // 外链跳转 URL
  iframeSrc?: string;         // iframe 嵌入地址
  openInNewWindow?: boolean;  // 新窗口打开

  // === 排序与参数 ===
  order?: number;             // 菜单排序（越小越靠前）
  query?: Record<string, any>; // 路由携带参数
  maxNumOfOpenTab?: number;   // 标签页最大打开数

  // === 菜单徽标 ===
  badge?: string;             // 徽标内容
  badgeType?: 'dot' | 'normal';     // 徽标类型
  badgeVariants?: BadgeVariant;      // 徽标颜色

  // === 扩展 ===
  [key: string]: any;         // 保留灵活性（兼容后端动态路由自定义属性）
}
```

### meta 常用配置速查

| 场景 | 配置 |
|------|------|
| 公开页面（不需要登录） | `ignoreAccess: true` |
| 需要特定权限 | `authority: ['sys:user:view']` |
| 菜单中隐藏但路由可访问 | `hideInMenu: true` |
| 标签页中隐藏 | `hideInTab: true` |
| 面包屑中隐藏 | `hideInBreadcrumb: true` |
| 固定标签页（常驻） | `affixTab: true, affixTabOrder: 1` |
| 菜单可见但点击显示 403 | `menuVisibleWithForbidden: true` + `authority` |
| 外链跳转 | `link: 'https://example.com'` |
| iframe 内嵌 | `iframeSrc: 'https://example.com'` |
| 手动指定高亮路径 | `activePath: '/system/users'` |

---

## 两种路由模式

### 前端模式（默认）

路由在前端静态定义（`src/router/modules/`），运行时根据用户权限码过滤：

```
modules/*.tsx 静态路由
  → mergeRouteModules 合并
  → generateRoutesByFrontend 权限过滤
  → injectRedirects 自动重定向
  → sortRoutes 排序
  → createBrowserRouter 创建
```

**特点**：
- 路由结构在构建时确定，运行时只做过滤
- 权限变化后需重新生成路由
- 适合权限结构相对固定的系统

### 后端模式

路由由后端 API 返回，前端通过组件映射表动态渲染：

```
API 返回路由树（component 为字符串）
  → generateRoutesByBackend
    → normalizeViewPath 标准化路径
    → layoutMap / pageMap 匹配组件
    → Suspense + lazy 包裹
  → createBrowserRouter 创建
```

**组件映射配置**（在 `src/router/index.tsx`）：

```typescript
// 布局映射：后端 component="BasicLayout" → React 组件
const layoutMap = { BasicLayout: AuthenticatedLayout };

// 页面映射：Vite glob 自动扫描 pages/app/**/*.tsx
// 后端 component="dashboard/index" → /dashboard → Dashboard 组件
const pageMap = buildPageMapFromGlob();
```

> 后端模式通过偏好设置 `app.accessMode` 切换，默认为 `'frontend'`。

---

## 路由守卫

### AuthGuard — 认证守卫

包裹需要登录才能访问的路由，未登录自动跳转到登录页并携带 `redirect` 参数：

```tsx
import { AuthGuard } from '@/router/guards';

// 在路由 element 中包裹
<AuthGuard>
  <MainLayout />
</AuthGuard>
```

登录成功后会自动跳回原页面（通过 URL 中的 `redirect` 参数）。

### GuestGuard — 访客守卫

已登录用户不能访问的路由（如登录页、注册页），自动重定向到首页：

```tsx
import { GuestGuard } from '@/router/guards';

<GuestGuard>
  <Login />
</GuestGuard>
```

---

## 懒加载

使用 `createLazyRoute` 实现路由级代码分割：

```tsx
import { createLazyRoute } from '@/router/utils/lazy';

// 自动包裹 Suspense + Spin 加载状态
{
  name: 'users',
  path: 'users',
  element: createLazyRoute(() => import('@/pages/app/system/users')),
}
```

> 路由模块中**必须**使用 `createLazyRoute`，不要直接使用 `lazy()` 或 `import()`。

---

## 菜单生成

### 自动生成流程

菜单从路由树自动提取，无需手动维护：

1. `generateMenus()` 遍历路由树，提取 `meta` 中的菜单信息
2. 转换为 `AppMenu[]` 格式（包含 `name`、`path`、`label`、`icon`、`order` 等）
3. 按 `meta.order` 排序
4. 过滤 `hideInMenu: true` 的路由
5. 处理 `hideChildrenInMenu`（子路由不展开）
6. 计算父子关系（`parent` / `parents`）用于面包屑

### 路径解析规则

子路由使用**相对路径**定义，菜单生成时会自动拼接父路径：

```typescript
// 父路由 path='system'，子路由 path='users'
// 菜单实际路径：/system/users
```

> **注意**：`selectedKeys` 使用完整绝对路径（如 `/dashboard`），菜单生成时已自动将相对路径转为绝对路径，确保菜单高亮正确。

### 菜单过滤

菜单支持权限过滤：

```typescript
import { shouldShowMenu } from '@/core/router';

// 检查菜单是否应显示
if (shouldShowMenu(menuItem, userPermissions)) {
  // 渲染菜单项
}
```

---

## 在组件中获取路由元数据

`meta` 会在路由初始化时自动转换为 React Router 的 `handle`，可在组件中通过 `useMatches` 获取：

```tsx
import { useMatches } from 'react-router-dom';

const MyPage = () => {
  const matches = useMatches();
  const currentMatch = matches[matches.length - 1];
  
  // 获取当前路由的元数据
  const { title, icon, authority } = currentMatch.handle || {};
  
  return <h1>{title}</h1>;
};
```

---

## 路由初始化流程

`AppRouter` 组件（`src/router/index.tsx`）负责整个路由初始化：

```
AppRouter 渲染
  → useEffect(isAuthenticated, accessMode)
    → 1. 获取用户权限码（auth.getUserPermissionCodes）
    → 2. 预加载字典数据（fetchAllDictEntries）
    → 3. 合并所有权限码（getAccessStatic）
    → 4. createAccessibleRouter(mode, options)
         → 根据 mode 选择生成策略
         → injectRedirects（自动重定向）
         → sortRoutes（排序）
         → transformMetaToHandle（meta → handle）
         → createBrowserRouter
    → 5. <RouterProvider router={router} />
```

**触发条件**：`isAuthenticated` 或 `accessMode` 变化时重新初始化。

---

## 工具函数速查

| 函数 | 文件 | 说明 |
|------|------|------|
| `createAccessibleRouter` | `factory.ts` | 根据模式创建完整路由器（主入口） |
| `generateRoutesByFrontend` | `generators/` | 前端模式路由生成 |
| `generateRoutesByBackend` | `generators/` | 后端模式路由生成 |
| `generateMenus` | `generators/` | 从路由树生成菜单 |
| `createLazyRoute` | `router/utils/lazy.tsx` | 创建懒加载路由元素 |
| `injectRedirects` | `utils/` | 自动注入重定向到第一个子路由 |
| `sortRoutes` | `utils/` | 按 `meta.order` 排序路由树 |
| `mergeRouteModules` | `utils/` | 合并 Vite glob 导入的路由模块 |
| `resetAccessibleRoutes` | `utils/` | 按白名单重置可访问路由 |
| `transformRoutesWithHandle` | `utils/` | 批量将 meta 转为 handle |
| `transformRoutesToMenu` | `utils/` | 转为 ProLayout 菜单格式 |
| `shouldShowMenu` | `generators/` | 判断菜单是否应显示 |

---

## 常见问题

### Q1: 新增路由后页面 404？

确认：
1. 路由文件放在 `src/router/modules/` 目录下
2. 文件导出 `AppRouteObject[]` 数组（支持 `default` 导出或具名导出）
3. `path` 使用相对路径（如 `'users'`，不是 `'/users'`）

### Q2: 菜单没有显示？

确认：
1. `meta.hideInMenu` 未设置为 `true`
2. `meta.title` 已配置（支持 i18n 键名）
3. 对应的 i18n 翻译文件中已有该键
4. 如有 `meta.authority`，当前用户拥有对应权限

### Q3: 菜单高亮/选中不正确？

菜单的 `key` 使用绝对路径（如 `/system/users`）。如果菜单自定义渲染，确保 `selectedKeys` 也使用绝对路径。路由的 `path` 应使用相对路径，菜单生成时会自动转为绝对路径。

### Q4: 如何设置只显示菜单但无权限时显示 403？

```typescript
{
  name: 'admin-panel',
  path: 'admin',
  meta: {
    title: 'routes:admin-panel',
    authority: ['sys:admin'],
    menuVisibleWithForbidden: true,  // 菜单可见，但无权限时渲染 Forbidden 页面
  },
}
```

### Q5: 如何添加外链或 iframe 页面？

```typescript
// 外链：点击菜单打开新窗口
{
  name: 'docs',
  path: 'docs',
  meta: {
    title: 'routes:docs',
    icon: 'lucide:book-open',
    link: 'https://docs.example.com',
  },
}

// iframe：在系统内嵌入外部页面
{
  name: 'external-tool',
  path: 'external-tool',
  meta: {
    title: 'routes:external-tool',
    iframeSrc: 'https://tool.example.com',
  },
}
```

### Q6: 如何在非组件环境中获取当前路由信息？

使用 React Router 的 `useMatches` 只能在组件中使用。如需在路由生成阶段做判断，使用静态工具：

```typescript
import { getAccessStatic } from '@/core/access';

const { getAllPermissions, hasAccessByCodes } = getAccessStatic();
```

---

## 注意事项

1. **name 字段必填**：`name` 是路由的唯一标识，权限匹配、菜单生成、路由重置都依赖它
2. **使用相对路径**：`modules/` 中的子路由使用相对路径（如 `'users'`），框架会自动拼接父路径
3. **权限字段统一为 `authority`**：使用 `meta.authority` 数组，不要使用 `meta.permission`（已废弃）
4. **404 路由必须在最后**：`error-routes.tsx` 中的 `path: '*'` 通配符必须放在路由配置的末尾
5. **懒加载必须用 `createLazyRoute`**：它会自动包裹 `Suspense` 和加载动画
6. **路由重置限制**：React Router v6 不支持运行时修改路由，`resetAccessibleRoutes` 返回过滤后的路由树，需重新创建 router
7. **i18n 标题格式**：`meta.title` 使用 `'routes:xxx'` 格式，对应 `locales/xx/_modules/routes.json` 中的键名

---

## 相关文件

- [路由类型定义](./types/index.ts) — AppRouteObject、RouteMeta、AppMenu 等
- [路由工厂](./factory.ts) — createAccessibleRouter 主入口
- [路由生成器](./generators/) — 前端/后端路由生成 + 菜单生成
- [工具函数](./utils/) — 重定向注入、排序、模块合并等
- [业务路由入口](../../router/index.tsx) — AppRouter 组件
- [业务路由模块](../../router/modules/) — 各功能模块路由配置
- [路由守卫](../../router/guards/) — AuthGuard / GuestGuard
- [权限模块](../access/) — 权限码获取与鉴权
