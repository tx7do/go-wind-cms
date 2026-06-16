# API 开发规约文档

本文档描述了项目中 API 层的架构设计、使用规范和最佳实践。

---

## 📋 目录

- [架构概览](#架构概览)
- [目录结构](#目录结构)
- [两层架构详解](#两层架构详解)
- [使用指南](#使用指南)
- [开发规范](#开发规范)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

---

## 架构概览

本项目采用 **两层 API 架构 + 统一 ApiClient 单例**，清晰分离关注点：

```
┌─────────────────────────────────────────┐
│         Vue Components (UI层)            │
└──────────────┬──────────────────────────┘
               │ 使用 Composables
┌──────────────▼──────────────────────────┐
│      API Composables Layer (Vue Query层)  │
│  - listXxx(), getXxx() - 业务纯函数       │
│  - useXxx() - Vue Query Hooks            │
│  - fetchXxx() - 非 Hook 数据获取方法       │
└──────────────┬──────────────────────────┘
               │ 调用 ApiClient
┌──────────────▼──────────────────────────┐
│         ApiClient (统一客户端单例)         │
│  - apiClient.xxxService - 懒加载服务      │
│  - 基于 ClientTransport 发送请求          │
└──────────────┬──────────────────────────┘
               │ 调用 Generated Code
┌──────────────▼──────────────────────────┐
│   Generated API Code (自动生成代码)       │
│  - createXxxServiceClient()             │
│  - createApiClient()                    │
│  - TypeScript Types                     │
└─────────────────────────────────────────┘
```

### 核心原则

1. **职责分离**：Composables 层负责业务逻辑 + Vue Query 集成，Generated 层负责协议映射
2. **统一客户端**：`apiClient` 单例统一管理所有 Service Client，懒加载、自动缓存
3. **类型安全**：全程 TypeScript 类型支持
4. **环境隔离**：组件中使用 `useXxx()` Hook，非组件环境中使用 `fetchXxx()` 或纯函数

---

## 目录结构

```
app/api/
├── client.ts                     # ApiClient 单例（统一入口）
├── generated/                    # 自动生成的 API 代码（不要手动修改）
│   └── app/service/v1/
│       └── index.ts              # 所有 API 类型、Service Client 工厂、ApiClient 类
│
├── composables/                  # Composables 层 - 业务逻辑 + Vue Query 集成
│   ├── auth.ts                   # 认证相关（login, logout, refreshToken...）
│   ├── user-profile.ts           # 用户资料（getMe, updateUser, bindPhone...）
│   ├── post.ts                   # 文章管理
│   ├── category.ts               # 分类管理
│   ├── page.ts                   # 页面管理
│   ├── tag.ts                    # 标签管理
│   ├── comment.ts                # 评论管理
│   ├── navigation.ts             # 导航管理
│   ├── file-transfer.ts          # 文件上传/下载
│   └── index.ts                  # 统一导出
│
└── index.ts                      # API 模块总入口
```

---

## 两层架构详解

### Composables 层（业务逻辑 + Vue Query 集成）

**位置**: `app/api/composables/*.ts`

**职责**:

- 导入 `apiClient` 单例，调用对应 Service 方法
- 导出纯函数式 API 方法（处理分页、参数组装等业务逻辑）
- 提供 `useMutation` / `useQuery` 封装（Vue Query Hooks）
- 同时提供非 Hook 的 `fetchXxx()` 方法（供 Store / 插件等非组件环境使用）

**标准模板**:

```typescript
import { useMutation, type UseMutationOptions } from '@tanstack/vue-query';
import { type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';
import { apiClient } from '@/api/client';
import { queryClient } from '@/plugins/vue-query';
import { getCurrentLocale } from '@/utils/locale';

// ==============================
// Service 层方法（使用 apiClient）
// ==============================

/**
 * 查询文章列表
 */
export async function listPost(
  paging?: Paging,
  formValues?: null | object,
  fieldMask?: undefined | string,
  orderBy?: null | string[],
  options?: { isTenantUser?: boolean; locale?: string },
) {
  const merged: Record<string, any> = {
    ...(formValues ?? {}),
    locale: options?.locale,
  };

  const noPaging = paging?.page === undefined && paging?.pageSize === undefined;
  // @ts-ignore proto generated code is error.
  return await apiClient.postService.List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取文章详情
 */
export async function getPost(id: number, locale?: string) {
  if (!id) return null;
  return await apiClient.postService.Get({ id, locale });
}

/**
 * 创建文章
 */
export async function createPost(values: Record<string, any> = {}) {
  return await apiClient.postService.Create({
    // @ts-ignore proto generated code is error.
    data: { ...values },
  });
}

/**
 * 更新文章
 */
export async function updatePost(id: number, values: Record<string, any> = {}) {
  return await apiClient.postService.Update({
    id,
    // @ts-ignore proto generated code is error.
    data: { ...values },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除文章
 */
export async function deletePost(id: number) {
  return await apiClient.postService.Delete({ id });
}

// ==============================
// Vue Query Hooks（用于组件中）
// ==============================

/** 列表查询参数 */
export interface ListPostParams {
  paging?: Paging;
  formValues?: null | object;
  fieldMask?: string;
  orderBy?: null | string[];
  isTenantUser?: boolean;
}

export function useListPost(options?: UseMutationOptions<any, Error, ListPostParams>) {
  return useMutation({
    mutationFn: (params) => {
      const locale = getCurrentLocale();
      return listPost(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser, locale },
      );
    },
    ...options,
  });
}

export function useGetPost(options?: UseMutationOptions<any, Error, number>) {
  return useMutation({
    mutationFn: (id) => {
      const locale = getCurrentLocale();
      return getPost(id, locale);
    },
    ...options,
  });
}

export function useCreatePost(options?: UseMutationOptions<{}, Error, Record<string, any>>) {
  return useMutation({
    mutationFn: (values) => createPost(values),
    ...options,
  });
}

// ==============================
// Fetch 方法（用于 Store / 非组件环境）
// ==============================

export async function fetchListPost(params: ListPostParams) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['listPost', params, locale],
    queryFn: () =>
      listPost(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser, locale },
      ),
    retry: 0,
  });
}

export async function fetchPost(id: number) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['getPost', id, locale],
    queryFn: () => getPost(id, locale),
    retry: 0,
  });
}
```

**关键要点**:

- ✅ 直接使用 `apiClient.xxxService.Method()` 调用后端 API
- ✅ 业务纯函数（`listXxx`, `getXxx`, `createXxx`...）可被 Hook 和 fetch 方法复用
- ✅ `useXxx()` Hook 用于 Vue 组件，享受 Vue Query 的缓存、重试等功能
- ✅ `fetchXxx()` 用于 Pinia Store、工具函数等非组件环境，设置 `retry: 0`
- ✅ 使用 `getCurrentLocale()` 自动注入当前语言
- ✅ 使用有意义的 `queryKey`，便于缓存管理
- ❌ 不要直接暴露 `apiClient` 给 UI 层组件

---

### ApiClient 单例（统一客户端）

**位置**: `app/api/client.ts`

**说明**:

`apiClient` 是全局唯一的 API 客户端实例，基于生成代码中的 `ApiClient` 类构建。它通过 `ClientTransport` 接口发送请求，内部对每个 Service Client 做懒加载缓存：

```typescript
// client.ts
import { type ClientTransport, createApiClient } from '@/api/generated/app/service/v1';
import { requestApi } from '@/core/transport/rest';

const transport: ClientTransport = {
  unary(path, method, body, _meta) {
    return requestApi({ body, method, path });
  },
  // serverStream / duplexStream 由独立模块管理
};

export const apiClient = createApiClient(transport);
```

**使用方式**:

```typescript
import { apiClient } from '@/api/client';

// 每个服务通过懒加载 getter 访问，首次使用时自动创建实例并缓存
apiClient.authenticationService.Login(request);
apiClient.postService.List(params);
apiClient.categoryService.Get({ id });
apiClient.userProfileService.GetUser({});
// ... 等
```

**支持的服务**:

| 属性名 | 服务 |
|--------|------|
| `apiClient.authenticationService` | 认证服务 |
| `apiClient.userProfileService` | 用户资料服务 |
| `apiClient.postService` | 文章服务 |
| `apiClient.categoryService` | 分类服务 |
| `apiClient.pageService` | 页面服务 |
| `apiClient.tagService` | 标签服务 |
| `apiClient.commentService` | 评论服务 |
| `apiClient.navigationService` | 导航服务 |
| `apiClient.fileTransferService` | 文件传输服务 |

---

### Generated 层（自动生成代码）

**位置**: `app/api/generated/app/service/v1/`

**说明**:

- ⚠️ **此目录由工具自动生成，不要手动修改**
- 包含所有 API 的 TypeScript 类型定义
- 包含每个 Service Client 的工厂函数（`createXxxServiceClient`）
- 包含统一的 `ApiClient` 类和 `createApiClient` 工厂
- 当后端 API 变更时，重新生成此目录

**主要导出**:

```typescript
// 类型定义
export type contentservicev1_Post = { ... }
export type authenticationservicev1_LoginRequest = { ... }

// Service Client 工厂
export function createPostServiceClient(transport: ClientTransport): PostService
export function createCategoryServiceClient(transport: ClientTransport): CategoryService

// 统一 ApiClient 类
export class ApiClient {
  get postService(): PostService { ... }
  get categoryService(): CategoryService { ... }
  // ...
}

export function createApiClient(transport: ClientTransport): ApiClient
```

---

## 使用指南

### 场景 1：在 Vue 组件中使用

```vue
<script setup lang="ts">
import { useListPost, useGetPost } from '@/api/composables/post';

const listPostMutation = useListPost();
const getPostMutation = useGetPost();

// 加载文章列表
const handleLoadPosts = async () => {
  const result = await listPostMutation.mutateAsync({
    paging: { page: 1, pageSize: 10 },
  });
  console.log('文章列表:', result);
};

// 获取文章详情
const handleGetPost = async (postId: number) => {
  const post = await getPostMutation.mutateAsync(postId);
  console.log('文章详情:', post);
};
</script>
```

### 场景 2：在 Pinia Store 中使用

```typescript
import { defineStore } from 'pinia';
import { fetchMe, fetchUpdateUser } from '@/api/composables/user-profile';
import { fetchLogin } from '@/api/composables/auth';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    userInfo: null as any,
  }),

  actions: {
    // 登录：使用 fetch 方法（非 Hook）
    async login(username: string, password: string) {
      try {
        const loginResult = await fetchLogin({
          username,
          password,
          grant_type: 'password',
        });

        // 获取用户信息
        this.userInfo = await fetchMe();
      } catch (error) {
        console.error('登录失败:', error);
        throw error;
      }
    },

    // 更新用户资料：使用纯函数
    async updateProfile(id: number, values: Record<string, any>) {
      await fetchUpdateUser(id, values);
    },
  },
});
```

### 场景 3：在 Nuxt 插件中使用

```typescript
// plugins/01.init-client.ts
import { refreshToken } from '@/api/composables/auth';
import { fetchMe } from '@/api/composables/user-profile';

export default defineNuxtPlugin(async (nuxtApp) => {
  // Token 刷新回调中使用纯函数
  const newToken = await refreshToken(refreshTokenValue);
  // 预加载用户信息
  const user = await fetchMe();
});
```

### 场景 4：直接使用纯函数

```typescript
import { createComment, deleteComment } from '@/api/composables/comment';

// 不需要 Vue Query 缓存的场景，直接调用纯函数
await createComment({ objectId: 1, content: '评论内容' });
await deleteComment(commentId);
```

---

## 开发规范

### 1. 命名规范

#### 业务纯函数命名

```typescript
// Composables 文件中导出的业务方法
export async function listXxx() { ... }      // 列表查询
export async function getXxx() { ... }       // 获取单个
export async function createXxx() { ... }    // 创建
export async function updateXxx() { ... }    // 更新
export async function deleteXxx() { ... }    // 删除
```

#### Vue Query Hooks 命名

```typescript
// 组合式函数（用于组件中）
export function useListXxx() { ... }
export function useGetXxx() { ... }
export function useCreateXxx() { ... }
export function useUpdateXxx() { ... }
export function useDeleteXxx() { ... }
```

#### Fetch 方法命名

```typescript
// 非 Hook 数据获取（供 Store / 插件使用）
export async function fetchListXxx() { ... }
export async function fetchXxx() { ... }
```

### 2. 类型规范

```typescript
// ✅ 正确：使用生成的类型
import type { contentservicev1_Post } from '@/api/generated/app/service/v1';

// ✅ 正确：明确标注返回类型
export async function getPost(id: number): Promise<contentservicev1_Post | null> {
  return apiClient.postService.Get({ id });
}

// ✅ 正确：Hook 中使用 UseMutationOptions
export function useGetPost(
  options?: UseMutationOptions<any, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getPost(id),
    ...options,
  });
}

// ❌ 错误：使用 any
export async function getPost(id: any): Promise<any> { ... }

// ❌ 错误：缺少类型标注
export function useGetPost(options?) { ... }
```

### 3. 错误处理规范

```typescript
// ✅ 纯函数：让错误自然抛出
export async function getPost(id: number) {
  return apiClient.postService.Get({ id, locale });
}

// ✅ Hook：通过 Vue Query 的 onError 处理
export function useGetPost(options?: UseMutationOptions<any, Error, number>) {
  return useMutation({
    mutationFn: (id) => getPost(id),
    onError: (error) => {
      console.error('获取文章失败:', error);
    },
    ...options,
  });
}

// ✅ Fetch 方法：由调用方决定如何处理
export async function fetchPost(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getPost', id, locale],
    queryFn: () => getPost(id),
    retry: 0,
  });
}

// 调用方处理错误
try {
  const post = await fetchPost(1);
} catch (error) {
  console.error('获取文章失败');
}
```

### 4. 注释规范

```typescript
/**
 * 查询文章列表
 */
export async function listPost(
  paging?: Paging,
  formValues?: null | object,
  fieldMask?: undefined | string,
  orderBy?: null | string[],
  options?: { isTenantUser?: boolean; locale?: string },
) { ... }

/**
 * 获取文章列表 Hook
 * @param options Vue Query 配置选项
 * @returns Mutation 对象
 */
export function useListPost(options?: UseMutationOptions<any, Error, ListPostParams>) { ... }

/**
 * 获取文章列表【给 Store / 外部调用】不带 Hook 的方法
 */
export async function fetchListPost(params: ListPostParams) { ... }
```

---

## 最佳实践

### 1. 选择合适的调用方式

| 使用场景 | 推荐方式 | 示例 |
|----------|----------|------|
| Vue 组件中 | `useXxx()` Hook | `const mutation = useGetPost()` |
| Pinia Store 中 | `fetchXxx()` 方法 | `await fetchPost(id)` |
| Nuxt 插件中 | `fetchXxx()` 或纯函数 | `await fetchLogin(req)` |
| 工具函数中 | `fetchXxx()` 方法 | `await fetchNavigation()` |
| 无需缓存的场景 | 纯函数直接调用 | `await createComment(values)` |

### 2. 合理配置 Vue Query 选项

```typescript
// ✅ 列表查询：启用缓存
export function useListPost(options?: UseMutationOptions<any, Error, ListPostParams>) {
  return useMutation({
    mutationFn: (params) => listPost(/* ... */),
    gcTime: 5 * 60 * 1000,  // 缓存 5 分钟
    ...options,
  });
}

// ✅ 创建/更新/删除：不缓存
export function useCreatePost(options?: UseMutationOptions<{}, Error, Record<string, any>>) {
  return useMutation({
    mutationFn: (values) => createPost(values),
    gcTime: 0,
    ...options,
  });
}
```

### 3. 利用 Vue Query 的缓存失效

```typescript
export function useCreatePost() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (values) => createPost(values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['listPost'] });
    },
  });
}
```

---

## 常见问题

### Q1: 为什么不再需要独立的 Service 层？

**A**: 原来的 Service 层负责两个职责：创建 Service Client 单例 + 封装业务参数。现在：
- **Service Client 单例**由 `apiClient` 统一管理（懒加载 getter，自动缓存）
- **业务参数封装**直接内联到 Composables 文件中的纯函数

这样减少了文件层级和样板代码，同时保持了同样的职责分离——业务逻辑与 Vue Query 集成是分开的（纯函数 vs Hook）。

### Q2: 什么时候用 `useXxx()`，什么时候用 `fetchXxx()`？

**A**:

- **`useXxx()`**：在 Vue 组件中使用，享受 Vue Query 的所有功能（缓存、重试、loading 状态等）
- **`fetchXxx()`**：在非组件环境（Pinia Store、Nuxt 插件、工具函数）中使用
- **纯函数**（如 `createComment`）：当不需要 Vue Query 缓存时，直接调用

### Q3: 如何添加新的 API 模块？

**A**: 遵循以下步骤：

1. **生成 API 代码**（后端提供 proto 文件后）
   ```bash
   npm run generate:api
   ```

2. **确认 apiClient 支持新服务**

   生成代码中的 `ApiClient` 类会自动包含新服务的 getter。确保 `createApiClient` 可用。

3. **创建 Composables 文件** (`app/api/composables/xxx.ts`)
   ```typescript
   import { useMutation, type UseMutationOptions } from '@tanstack/vue-query';
   import { type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';
   import { apiClient } from '@/api/client';
   import { queryClient } from '@/plugins/vue-query';
   import { getCurrentLocale } from '@/utils/locale';

   // 纯函数
   export async function listXxx(paging?, formValues?, fieldMask?, orderBy?, options?) {
     return await apiClient.xxxService.List({ /* 组装参数 */ });
   }

   export async function getXxx(id: number) {
     return await apiClient.xxxService.Get({ id });
   }

   // Hook
   export function useListXxx(options?: UseMutationOptions<any, Error, ListXxxParams>) {
     return useMutation({
       mutationFn: (params) => listXxx(/* ... */),
       ...options,
     });
   }

   // Fetch 方法
   export async function fetchListXxx(params: ListXxxParams) {
     return queryClient.fetchQuery({
       queryKey: ['listXxx', params],
       queryFn: () => listXxx(/* ... */),
       retry: 0,
     });
   }
   ```

4. **导出模块** (`app/api/composables/index.ts`)
   ```typescript
   export * from './xxx';
   ```

### Q4: 如何处理分页查询？

**A**: 使用 `Paging` 类型和工具函数：

```typescript
import { type Paging, makeOrderBy, makeQueryString } from '@/core/transport/rest';

// 纯函数中处理分页
export async function listPost(
  paging?: Paging,
  formValues?: null | object,
  fieldMask?: undefined | string,
  orderBy?: null | string[],
  options?: { isTenantUser?: boolean; locale?: string },
) {
  const merged: Record<string, any> = {
    ...(formValues ?? {}),
    locale: options?.locale,
  };

  const noPaging = paging?.page === undefined && paging?.pageSize === undefined;
  return await apiClient.postService.List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

// Hook 中使用
export function useListPost(options?: UseMutationOptions<any, Error, ListPostParams>) {
  return useMutation({
    mutationFn: (params) => {
      const locale = getCurrentLocale();
      return listPost(
        params.paging, params.formValues, params.fieldMask,
        params.orderBy, { isTenantUser: params.isTenantUser, locale },
      );
    },
    ...options,
  });
}
```

### Q5: 如何调试 API 调用？

**A**:

1. **使用 Vue Query Devtools**
   ```typescript
   // 开发模式下自动加载，按 Shift + Alt + D 打开
   ```

2. **查看网络请求**
   - 浏览器开发者工具 → Network 标签
   - 过滤 XHR/Fetch 请求

3. **日志输出**
   ```typescript
   export async function listPost(paging?, formValues?, ...) {
     console.log('[API] listPost called with:', { paging, formValues });
     const result = await apiClient.postService.List({ /* ... */ });
     console.log('[API] listPost result:', result);
     return result;
   }
   ```

---

## 附录

### 相关资源

- [TanStack Query 官方文档](https://tanstack.com/query/latest)
- [TypeScript 官方文档](https://www.typescriptlang.org/docs/)

### 项目特定资源

- 传输层实现：`app/core/transport/`
- 偏好设置系统：`app/core/preferences/`
- 状态管理规范：`app/stores/`
