# API 开发规约文档

本文档描述了项目中 API 层的架构设计、使用规范和最佳实践。

---

## 📋 目录

- [架构概览](#架构概览)
- [目录结构](#目录结构)
- [三层架构详解](#三层架构详解)
- [使用指南](#使用指南)
- [开发规范](#开发规范)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

---

## 架构概览

本项目采用**三层 API 架构**，清晰分离关注点：

```
┌─────────────────────────────────────────┐
│         React Components (UI层)          │
└──────────────┬──────────────────────────┘
               │ 使用 Hooks
┌──────────────▼──────────────────────────┐
│      API Hooks Layer (React Query层)     │
│  - useXxx() - React Hooks               │
│  - fetchXxx() - 非 Hook 方法             │
└──────────────┬──────────────────────────┘
               │ 调用 Service
┌──────────────▼──────────────────────────┐
│      API Service Layer (服务封装层)       │
│  - getXxxService() - Service Client     │
│  - listXxx() / listXxxRaw() - 纯函数     │
│  - getXxx() / createXxx() 等             │
└──────────────┬──────────────────────────┘
               │ 调用 Generated Code
┌──────────────▼──────────────────────────┐
│   Generated API Code (自动生成代码)       │
│  - createXxxServiceClient()             │
│  - TypeScript Types                     │
└─────────────────────────────────────────┘
```

### 核心原则

1. **职责分离**：每层只负责自己的职责
2. **单向依赖**：上层依赖下层，不反向依赖
3. **类型安全**：全程 TypeScript 类型支持
4. **环境隔离**：区分 React 组件环境和非 React 环境

---

## 目录结构

```
src/api/
├── generated/                    # 自动生成的 API 代码（不要手动修改）
│   └── app/service/v1/
│       └── index.ts              # 所有 API 类型和 Client 工厂
│
├── service/                      # Service 层 - 服务封装
│   ├── auth.ts                   # 认证服务（login / logout / refreshToken）
│   ├── user-profile.ts           # 用户资料服务（getMe）
│   ├── post.ts                   # 文章管理服务
│   ├── category.ts               # 分类管理服务
│   ├── tag.ts                    # 标签管理服务
│   ├── comment.ts                # 评论管理服务
│   ├── page.ts                   # 页面管理服务
│   ├── navigation.ts             # 导航管理服务
│   ├── file-transfer.ts          # 文件传输服务（上传 / 下载）
│   └── index.ts                  # 统一导出
│
├── hooks/                        # Hooks 层 - React Query 集成
│   ├── auth.ts                   # 认证相关 Hooks + useAuth 编排 Hook
│   ├── user-profile.ts           # 用户资料 Hooks（useQuery 形式）
│   ├── post.ts                   # 文章管理 Hooks + 辅助函数
│   ├── category.ts               # 分类管理 Hooks + 辅助函数
│   ├── tag.ts                    # 标签管理 Hooks + 辅助函数
│   ├── comment.ts                # 评论管理 Hooks
│   ├── page.ts                   # 页面管理 Hooks
│   ├── navigation.ts             # 导航管理 Hooks
│   ├── file-transfer.ts          # 文件传输 Hooks（含上传进度）
│   └── index.ts                  # 统一导出
│
└── index.ts                      # API 模块总入口（重导出 generated + hooks）
```

---

## 三层架构详解

### Service 层（服务封装层）

**位置**: `src/api/service/*.ts`

**职责**:

- 封装生成的 API Client
- 提供单例模式的 Service Client
- 导出纯函数式的 API 调用方法
- 自动注入 `locale` 到查询参数
- 提供两套列表方法：`listXxx(query)` 和 `listXxxRaw(params)`

**标准模板**（以 `tag.ts` 为例）:

```typescript
import {
  type contentservicev1_Tag,
  type contentservicev1_ListTagResponse,
  createTagServiceClient,
} from '@/api/generated/app/service/v1';
import { type PaginationQuery, requestApi } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

// 单例模式：缓存 Service Client 实例
let _instance: ReturnType<typeof createTagServiceClient> | null = null;

/**
 * 获取标签服务单例（延迟初始化）
 */
export function getTagService() {
  if (!_instance) {
    _instance = createTagServiceClient(requestApi);
  }
  return _instance;
}

// ==============================
// 标签管理 API
// ==============================

/**
 * 获取标签列表（PaginationQuery 方式）
 */
export async function listTags(query: PaginationQuery) {
  const params = query.toRawParams();
  const locale = currentLocaleLanguageCode();
  const formValues = query.formValues ? {...query.formValues, locale} : {locale};

  return getTagService().List({
    ...params,
    query: JSON.stringify(formValues),
  });
}

/**
 * 兼容旧调用方式 - 通过原始参数获取标签列表
 */
export async function listTagsRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<contentservicev1_ListTagResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return getTagService().List({
    fieldMask: params.fieldMask,
    orderBy: params.orderBy ? JSON.stringify(params.orderBy) : undefined,
    sorting: Array.isArray(params.orderBy)
      ? params.orderBy.map((o) => ({field: o, direction: 'ASC'}))
      : undefined,
    query: JSON.stringify(formValues),
    page: params.paging?.page,
    pageSize: params.paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取单个标签
 */
export async function getTag(id: number) {
  return getTagService().Get({id});
}

/**
 * 创建标签
 */
export async function createTag(values: Partial<contentservicev1_Tag>) {
  return getTagService().Create({
    data: values as contentservicev1_Tag,
  });
}

/**
 * 更新标签
 */
export async function updateTag(params: {
  id: number;
  values: Partial<contentservicev1_Tag>;
}) {
  return getTagService().Update({
    id: params.id,
    data: params.values as contentservicev1_Tag,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除标签
 */
export async function deleteTag(id: number) {
  return getTagService().Delete({id});
}
```

**关键要点**:

- ✅ 使用单例模式缓存 Service Client
- ✅ 导出 `getXxxService()` 函数供外部使用
- ✅ 导出纯函数式 API 方法（`async/await`）
- ✅ `listXxx()` 接收 `PaginationQuery` 对象，调用 `query.toRawParams()`
- ✅ `listXxxRaw()` 接收原始参数对象，Hooks 层主要调用此方法
- ✅ 自动注入 `currentLocaleLanguageCode()` 到 `formValues`（文章、分类、标签、页面）
- ✅ `updateXxx()` 接收 `{ id, values }` 参数，自动生成 `updateMask`
- ✅ `deleteXxx()` 接收 `id: number`
- ❌ 不要在 Service 层使用 React Hooks
- ❌ 不要直接暴露 Service Client 给 UI 层

---

### Hooks 层（React Query 集成层）

**位置**: `src/api/hooks/*.ts`

**职责**:

- 将 Service 层的函数包装为 React Hooks（`useMutation` / `useQuery`）
- 同时提供非 Hook 的 `fetchXxx()` 方法（基于 `queryClient.fetchQuery`）
- 提供辅助函数（如 `getPostTitle`、`getCategoryName` 等翻译提取函数）

**标准模板（useMutation 形式）**（以 `tag.ts` 为例）:

```typescript
import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_Tag,
  type contentservicev1_TagTranslation,
  type contentservicev1_ListTagResponse,
} from '@/api/generated/app/service/v1';
import {
  listTagsRaw,
  getTag,
  createTag,
  updateTag,
  deleteTag,
} from '@/api/service/tag';
import { queryClient } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

// ==============================
// 标签列表 Hook
// ==============================
export function useListTags(
  options?: UseMutationOptions<
    contentservicev1_ListTagResponse,
    Error,
    {
      paging?: { page?: number; pageSize?: number };
      formValues?: object | undefined;
      fieldMask?: string | undefined;
      orderBy?: string[] | undefined;
    }
  >,
) {
  return useMutation({
    mutationFn: (params) => listTagsRaw(params),
    ...options,
  });
}

// ==============================================
// 获取标签列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListTags(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}) {
  return queryClient.fetchQuery({
    queryKey: ['listTags', params],
    queryFn: () => listTagsRaw(params),
    retry: 0,
  });
}

// ==============================
// 获取单个标签 Hook
// ==============================
export function useGetTag(
  options?: UseMutationOptions<contentservicev1_Tag, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getTag(id),
    ...options,
  });
}

// ==============================================
// 获取单个标签 【给 Store / 外部调用】
// ==============================================
export async function fetchTag(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getTag', id],
    queryFn: () => getTag(id),
    retry: 0,
  });
}

// ==============================
// CRUD Hooks（create / update / delete）
// ==============================
export function useCreateTag(
  options?: UseMutationOptions<contentservicev1_Tag, Error, Partial<contentservicev1_Tag>>,
) {
  return useMutation({
    mutationFn: (data) => createTag(data),
    ...options,
  });
}

export function useUpdateTag(
  options?: UseMutationOptions<
    contentservicev1_Tag,
    Error,
    { id: number; values: Partial<contentservicev1_Tag> }
  >,
) {
  return useMutation({
    mutationFn: (params) => updateTag(params),
    ...options,
  });
}

export function useDeleteTag(options?: UseMutationOptions<Record<string, never>, Error, number>) {
  return useMutation({
    mutationFn: (id) => deleteTag(id),
    ...options,
  });
}

// ==============================
// 辅助函数（纯工具，不依赖 Store）
// ==============================

/**
 * 获取标签的翻译（根据当前 locale 匹配）
 */
export function getTranslation(tag: contentservicev1_Tag | null) {
  if (!tag || !tag?.translations || tag.translations.length === 0) return null;

  const locale = currentLocaleLanguageCode();
  const translation = tag.translations?.find(
    (t: contentservicev1_TagTranslation) => t.languageCode === locale
  );
  return translation || tag.translations?.[0];
}
```

**useQuery 形式**（以 `user-profile.ts` 为例）:

```typescript
import {
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type identityservicev1_User,
} from '@/api/generated/app/service/v1';
import { getMe } from '@/api/service/user-profile';
import { queryClient } from '@/core';

// ==============================
// 获取当前用户信息 Hook（useQuery 形式）
// ==============================
export function useGetUserProfile(options?: UseQueryOptions<identityservicev1_User | null, Error>) {
  return useQuery({
    queryKey: ['getMe'],
    queryFn: () => getMe(),
    ...options,
  });
}

// ==============================================
// 获取当前用户信息 【给 Store / 外部调用】
// ==============================================
export async function fetchUserProfile() {
  return queryClient.fetchQuery({
    queryKey: ['userProfile'],
    queryFn: () => getMe(),
    retry: 0,
  });
}
```

**编排 Hook 形式**（以 `auth.ts` 的 `useAuth` 为例）:

```typescript
// auth.ts 提供了高层编排 Hook，协调多个 API 调用和 Store
export function useAuth() {
    const accessStore = useAccessStore();
    const userStore = useUserStore();
    const router = useI18nRouter();

    async function login(params: LoginParams, onSuccess?: () => Promise<void> | void) {
        // 1. 调用登录 API
        const response = await fetchLogin({...});
        // 2. 存储 token
        accessStore.setAccessToken({...});
        // 3. 获取用户信息
        const userInfo = await fetchUserProfile();
        userStore.setUser(userInfo);
        // 4. 路由跳转
        router.push(userInfo.homePage || '/');
    }

    async function logout(redirect: boolean = true) { ... }
    async function reauthenticate() { ... }
    async function refreshToken() { ... }

    return { login, logout, refreshToken, reauthenticate, fetchUserProfile };
}
```

**关键要点**:

- ✅ 列表 Hook 统一使用 `useMutation`（非 `useQuery`），参数类型为原始对象
- ✅ `user-profile` 使用 `useQuery`（适合"获取单条数据并缓存"的场景）
- ✅ `useMutation` 的参数类型是 `{ paging?, formValues?, fieldMask?, orderBy? }`
- ✅ `fetchXxx()` 方法基于 `queryClient.fetchQuery`，`retry: 0`
- ✅ `useAuth()` 是编排 Hook，协调 `fetchLogin` + `fetchUserProfile` + Store
- ❌ Hooks 层不直接创建 Service Client
- ❌ Hooks 层不直接使用 `requestApi`

---

### Generated 层（自动生成代码）

**位置**: `src/api/generated/app/service/v1/`

**说明**:

- ⚠️ **此目录由工具自动生成，不要手动修改**
- 包含所有 API 的 TypeScript 类型定义
- 包含所有 Service Client 的工厂函数
- 当后端 API 变更时，重新生成此目录

**主要导出**:

```typescript
// 类型定义
export type contentservicev1_Post = { ... }
export type contentservicev1_Category = { ... }
export type contentservicev1_Tag = { ... }
export type identityservicev1_User = { ... }
export type authenticationservicev1_LoginRequest = { ... }

// Service Client 工厂
export function createPostServiceClient(requestApi: RequestApi): PostServiceClient
export function createCategoryServiceClient(requestApi: RequestApi): CategoryServiceClient
export function createTagServiceClient(requestApi: RequestApi): TagServiceClient
export function createAuthenticationServiceClient(requestApi: RequestApi): AuthenticationServiceClient
// ... 更多 Service Client
```

---

## 使用指南

### 场景 1：在 React 组件中使用（fetchXxx + useEffect）

这是项目中最常见的模式，适合需要精细控制加载时机的场景：

```tsx
import {useState, useEffect} from 'react';
import {fetchListPosts} from '@/api/hooks/post';
import type {contentservicev1_Post, contentservicev1_ListPostResponse} from '@/api/generated/app/service/v1';

function PostListPage() {
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        let cancelled = false;
        (async () => {
            setLoading(true);
            try {
                const res = await fetchListPosts({
                    paging: {page: 1, pageSize: 10},
                    formValues: {status: 'POST_STATUS_PUBLISHED'},
                    orderBy: ['-createdAt'],
                }) as unknown as contentservicev1_ListPostResponse;

                if (!cancelled) setPosts(res.items || []);
            } catch (error) {
                console.error('Load failed:', error);
            } finally {
                if (!cancelled) setLoading(false);
            }
        })();
        return () => { cancelled = true; };
    }, []);

    return <div>...</div>;
}
```

### 场景 2：使用 useAuth 编排 Hook

```tsx
import {useAuth} from '@/api/hooks/auth';

function LoginPage() {
    const {login} = useAuth();

    const handleLogin = async () => {
        try {
            await login({
                username: 'admin',
                password: '123456',
                grant_type: 'password',
            });
            // login 内部已处理 token 存储、用户信息获取、路由跳转
        } catch (error) {
            console.error('登录失败:', error);
        }
    };

    return <button onClick={handleLogin}>登录</button>;
}
```

### 场景 3：使用辅助函数提取翻译

```tsx
import {getPostTitle, getPostSummary, getPostThumbnail} from '@/api/hooks/post';
import {getCategoryName, getCategoryDescription} from '@/api/hooks/category';
import {getTranslation as getTagTranslation} from '@/api/hooks/tag';

// 根据当前 locale 自动匹配翻译
const title = getPostTitle(post);
const categoryName = getCategoryName(category);
```

### 场景 4：使用 useQuery Hook（用户资料场景）

```tsx
import {useGetUserProfile} from '@/api/hooks/user-profile';

function UserAvatar() {
    const {data: user, isLoading} = useGetUserProfile();

    if (isLoading) return <Skeleton/>;
    return <span>{user?.nickname}</span>;
}
```

### 场景 5：文件上传（带进度）

```tsx
import {useUploadFile} from '@/api/hooks/file-transfer';

function UploadButton() {
    const {mutateAsync: upload, isPending} = useUploadFile();

    const handleUpload = async (file: File) => {
        await upload({
            bucketName: 'public',
            fileDirectory: 'images',
            file,
            method: 'post',
            onUploadProgress: (progressEvent) => {
                const percent = Math.round(
                    (progressEvent.loaded * 100) / (progressEvent.total || 1)
                );
                console.log(`上传进度: ${percent}%`);
            },
        });
    };

    return <input type="file" onChange={(e) => e.target.files?.[0] && handleUpload(e.target.files[0])}/>;
}
```

---

## 开发规范

### 1. 命名规范

#### Service 层命名

```typescript
// Service Client 获取函数
export function getTagService() { ... }
export function getPostService() { ... }

// 列表方法（两套）
export async function listTags(query: PaginationQuery) { ... }    // PaginationQuery 方式
export async function listTagsRaw(params: RawParams) { ... }      // 原始参数方式（Hooks 主要调用）

// 单个 / CRUD
export async function getTag(id: number) { ... }
export async function createTag(values: Partial<T>) { ... }
export async function updateTag(params: { id: number; values: Partial<T> }) { ... }
export async function deleteTag(id: number) { ... }
```

#### Hooks 层命名

```typescript
// useMutation Hooks
export function useListTags(options?) { ... }
export function useGetTag(options?) { ... }
export function useCreateTag(options?) { ... }
export function useUpdateTag(options?) { ... }
export function useDeleteTag(options?) { ... }

// useQuery Hooks（适合单条数据自动缓存）
export function useGetUserProfile(options?) { ... }

// Fetch 方法（非 Hook，用于 useEffect / Store）
export async function fetchListTags(params) { ... }
export async function fetchTag(id: number) { ... }

// 辅助函数
export function getTranslation(entity) { ... }
export function getPostTitle(post) { ... }
```

### 2. 类型规范

```typescript
// ✅ 正确：使用生成的类型
import type { contentservicev1_Tag } from '@/api/generated/app/service/v1';

// ✅ 正确：列表 Hook 的参数类型
export function useListTags(
  options?: UseMutationOptions<
    contentservicev1_ListTagResponse,
    Error,
    {
      paging?: { page?: number; pageSize?: number };
      formValues?: object | undefined;
      fieldMask?: string | undefined;
      orderBy?: string[] | undefined;
    }
  >,
) { ... }

// ✅ 正确：更新操作参数
export function useUpdateTag(
  options?: UseMutationOptions<
    contentservicev1_Tag,
    Error,
    { id: number; values: Partial<contentservicev1_Tag> }
  >,
) { ... }

// ❌ 错误：使用 any
export async function getTag(id: any): Promise<any> { ... }

// ❌ 错误：缺少类型标注
export function useGetTag(options?) { ... }
```

### 3. 错误处理规范

```typescript
// ✅ Service 层：让错误自然抛出
export async function getTag(id: number) {
  return getTagService().Get({ id });  // 错误会向上传播
}

// ✅ React 组件：在 useEffect 中 try/catch
useEffect(() => {
    (async () => {
        try {
            const res = await fetchListTags({...});
            setTags(res.items || []);
        } catch (error) {
            console.error('Failed to load tags:', error);
            setTags([]);
        } finally {
            setLoading(false);
        }
    })();
}, []);

// ✅ Hook 调用：使用 mutateAsync + try/catch
const { mutateAsync } = useCreateTag();
try {
    await mutateAsync({ name: '新标签' });
} catch (error) {
    console.error('创建失败:', error);
}
```

### 4. 注释规范

```typescript
/**
 * 获取标签服务单例（延迟初始化）
 */
export function getTagService() { ... }

/**
 * 兼容旧调用方式 - 通过原始参数获取标签列表
 */
export async function listTagsRaw(params: { ... }): Promise<contentservicev1_ListTagResponse> { ... }

// ==============================================
// 获取标签列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListTags(params: { ... }) { ... }
```

---

## 最佳实践

### 1. 选择合适的调用方式

| 使用场景 | 推荐方式 | 示例 |
|---------|---------|------|
| React 组件 + useEffect | `fetchXxx()` | `await fetchListPosts({...})` |
| React 组件 + 自动缓存 | `useQuery` Hook | `useGetUserProfile()` |
| React 组件 + 手动触发 | `useMutation` Hook | `const {mutateAsync} = useCreateTag()` |
| Zustand Store / 非组件 | `fetchXxx()` 方法 | `await fetchUserProfile()` |
| 登录 / 认证流程 | 编排 Hook | `const {login} = useAuth()` |
| 文件上传 | `useUploadFile` Hook | 支持上传进度回调 |

### 2. QueryClient 默认配置

项目在 `src/core/query-client.ts` 中配置了全局默认值：

```typescript
export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 5 * 60 * 1000,       // 数据新鲜度：5分钟
            gcTime: 30 * 60 * 1000,          // 缓存时间：30分钟
            refetchOnWindowFocus: false,      // 窗口聚焦时不自动刷新
            retry: (failureCount, error) => {
                if (process.env.NODE_ENV === 'production') {
                    return failureCount < 1 && !!error.message?.includes('Network');
                }
                return failureCount < 3;
            },
            retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
        },
        mutations: {
            onError: (error) => {
                console.error('Mutation error:', error);
            },
        },
    },
});
```

### 3. 语言切换时清除缓存

项目通过 `useI18n.changeLocale` 在语言切换前清除所有 React Query 缓存，确保用新语言重新请求：

```typescript
// src/i18n/helpers/useI18n.ts
const changeLocale = (targetLocale: string) => {
    setLanguage(targetLocale);
    queryClient.clear();  // 清除所有缓存
    router.replaceWithoutLocale(newPath);
};
```

同时 `ClientLocaleLayout` 使用 `key={locale}` 强制内容区重挂载。

### 4. 调试 API 调用

1. **React Query DevTools**

   已集成 `@tanstack/react-query-devtools`（仅开发环境），位于页面右下角，可查看所有 query 的缓存状态、queryKey、数据等。

2. **查看网络请求**
   - 浏览器开发者工具 → Network 标签
   - 所有请求自动携带 `Accept-Language` header

3. **日志输出**

   ```typescript
   export async function listTagsRaw(params) {
     console.log('[API] listTagsRaw called with:', params);
     const result = await getTagService().List({...});
     console.log('[API] listTagsRaw result:', result);
     return result;
   }
   ```

---

## 常见问题

### Q1: 为什么需要 Service 层和 Hooks 层？

**A**:

- **Service 层**：封装纯函数式 API 调用，自动注入 locale，处理参数转换，可在任何环境中使用
- **Hooks 层**：集成 React Query，提供声明式数据获取、缓存、重试等功能
- **分离好处**：职责清晰、易于测试、可复用性强

### Q2: `listXxx()` 和 `listXxxRaw()` 的区别？

**A**:

- **`listXxx(query: PaginationQuery)`**：接收 `PaginationQuery` 对象，通过 `query.toRawParams()` 转换参数
- **`listXxxRaw(params: RawParams)`**：接收原始参数对象 `{ paging?, formValues?, fieldMask?, orderBy? }`，**Hooks 层主要调用此方法**
- 两者都自动注入 `currentLocaleLanguageCode()` 到 `formValues`（文章、分类、标签、页面模块）

### Q3: 什么时候用 `useMutation`，什么时候用 `useQuery`？

**A**:

- **`useMutation`**：项目中的**默认选择**，适合列表查询（需手动触发）、创建、更新、删除
- **`useQuery`**：适合"自动获取并缓存"的场景，目前仅 `useGetUserProfile` 使用
- 列表查询统一使用 `useMutation`（非 `useQuery`），在 `useEffect` 中通过 `fetchXxx()` 手动触发

### Q4: 如何添加新的 API 模块？

**A**: 遵循以下步骤：

1. **生成 API 代码**（后端提供 proto 文件后）
   ```bash
   npm run generate:api
   ```

2. **创建 Service 文件** (`src/api/service/bookmark.ts`)
   ```typescript
   import { createBookmarkServiceClient } from '@/api/generated/app/service/v1';
   import { type PaginationQuery, requestApi } from '@/core';
   import { currentLocaleLanguageCode } from '@/i18n';

   let _instance: ReturnType<typeof createBookmarkServiceClient> | null = null;

   export function getBookmarkService() {
     if (!_instance) {
       _instance = createBookmarkServiceClient(requestApi);
     }
     return _instance;
   }

   export async function listBookmarksRaw(params: {
     paging?: { page?: number; pageSize?: number };
     formValues?: object | undefined;
     fieldMask?: string | undefined;
     orderBy?: string[] | undefined;
   }) { ... }

   export async function getBookmark(id: number) { ... }
   export async function createBookmark(values: Partial<T>) { ... }
   export async function updateBookmark(params: { id: number; values: Partial<T> }) { ... }
   export async function deleteBookmark(id: number) { ... }
   ```

3. **创建 Hooks 文件** (`src/api/hooks/bookmark.ts`)
   ```typescript
   import { useMutation, type UseMutationOptions } from '@tanstack/react-query';
   import { listBookmarksRaw, getBookmark, createBookmark, updateBookmark, deleteBookmark } from '@/api/service/bookmark';
   import { queryClient } from '@/core';

   export function useListBookmarks(options?) {
     return useMutation({ mutationFn: (params) => listBookmarksRaw(params), ...options });
   }

   export async function fetchListBookmarks(params) {
     return queryClient.fetchQuery({
       queryKey: ['listBookmarks', params],
       queryFn: () => listBookmarksRaw(params),
       retry: 0,
     });
   }
   // ... useGetBookmark, useCreateBookmark, useUpdateBookmark, useDeleteBookmark
   ```

4. **导出模块** (`src/api/service/index.ts` 和 `src/api/hooks/index.ts`)
   ```typescript
   export * from './bookmark';
   ```

### Q5: 为什么列表查询用 `useMutation` 而不是 `useQuery`？

**A**: 这是项目早期的设计选择，主要原因：

- **灵活控制**：列表查询通常需要配合分页、筛选等交互，`useMutation` 提供手动触发（`mutate`），不受 `useQuery` 的自动 refetch 机制影响
- **配合 `fetchXxx`**：实际使用中，组件主要通过 `fetchListXxx()` + `useEffect` 获取列表数据（利用 React Query 的缓存能力），而非直接使用 `useListXxx` Hook
- **历史兼容**：从旧版迁移时保留了 `useMutation` 的模式

---

## 附录

### 相关资源

- [TanStack Query 官方文档](https://tanstack.com/query/latest)
- [Zustand 官方文档](https://zustand-demo.pmnd.rs/)
- [TypeScript 官方文档](https://www.typescriptlang.org/docs/)
- [Next.js App Router 文档](https://nextjs.org/docs/app)

### 项目特定资源

- Core 模块：`src/core/`（包含 QueryClient、RequestClient、PaginationQuery）
- 状态管理：`src/store/`（Zustand stores：access、user、loading、preferences）
- 主题与偏好：`src/core/preferences/`（语言、主题、样式偏好管理）
