# Hooks (自定义钩子)

## 概述

本目录包含项目特定的 Vue 3 Hooks，即自定义的业务逻辑钩子函数。

### 目录结构

```
hooks/
├── use-app-config.ts       # 应用配置 hook
├── useAutoThemeSwitcher.ts # 自动主题切换 hook
└── README.md               # 本文档
```

## Hooks vs Composables

| 维度       | Hooks                     | Composables                   |
|----------|---------------------------|-------------------------------|
| **存储位置** | `hooks/`                  | `composables/`                |
| **用途**   | 项目特定业务逻辑                  | 通用可复用逻辑                       |
| **依赖**   | 通常依赖应用状态和服务               | 最小依赖                          |
| **复用性**  | 中等                        | 高                             |
| **示例**   | `useAppConfig`, `useAuth` | `useLocalStorage`, `useFetch` |

## 现有 Hooks

### use-app-config.ts

**用途**: 获取和管理应用配置

```typescript
import {useAppConfig} from '@/hooks/use-app-config'

const {config, updateConfig} = useAppConfig()
console.log(config.value.apiBaseUrl)
updateConfig({apiBaseUrl: 'https://api.example.com'})
```

### useAutoThemeSwitcher.ts

**用途**: 处理自动主题切换逻辑

```typescript
import {useAutoThemeSwitcher} from '@/hooks/useAutoThemeSwitcher'
import {useAppStore} from '@/stores'

const appStore = useAppStore()
const {initializeThemeSwitcher, switchTheme} = useAutoThemeSwitcher(appStore)

// 初始化主题切换监听
initializeThemeSwitcher()

// 手动切换主题
switchTheme('dark')
```

## 编写 Hooks 的最佳实践

### 1. 清晰的职责

```typescript
/**
 * 管理表单提交状态
 * @param onSubmit 提交回调
 * @returns 提交状态和方法
 */
export function useFormSubmit(onSubmit: (data: any) => Promise<void>) {
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const submit = async (data: any) => {
    loading.value = true
    error.value = null
    try {
      await onSubmit(data)
    } catch (e) {
      error.value = e as Error
    } finally {
      loading.value = false
    }
  }

  return {loading, error, submit}
}
```

### 2. 类型安全

```typescript
import type {Ref} from 'vue'

interface AppConfig {
  apiBaseUrl: string
  apiTimeout: number
  enableMock: boolean
}

export function useAppConfig(): {
  config: Ref<AppConfig>
  updateConfig: (updates: Partial<AppConfig>) => void
} {
  // 实现...
}
```

### 3. 依赖注入

```typescript
import {useAppStore} from '@/stores'

export function useTheme(store?: ReturnType<typeof useAppStore>) {
  const appStore = store ?? useAppStore()

  const isDark = computed(() => appStore.mode === 'dark')

  const toggleTheme = () => {
    appStore.setMode(isDark.value ? 'light' : 'dark')
  }

  return {isDark, toggleTheme}
}
```

### 4. 清晰的返回接口

```typescript
interface UseAuthReturn {
  isAuthenticated: Ref<boolean>
  user: Ref<User | null>
  login: (credentials: Credentials) => Promise<void>
  logout: () => void
  loading: Ref<boolean>
}

export function useAuth(): UseAuthReturn {
  // 实现...
}
```

### 5. 生命周期管理

```typescript
import {onMounted, onUnmounted} from 'vue'

export function usePolling(callback: () => void, interval = 5000) {
  let timer: ReturnType<typeof setInterval> | null = null

  const start = () => {
    timer = setInterval(callback, interval)
  }

  const stop = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }

  onMounted(() => start())
  onUnmounted(() => stop())

  return {start, stop}
}
```

## 常见 Hooks 模式

### 认证相关

```typescript
export function useAuth() {
  const authStore = useAuthStore()

  const login = async (credentials: Credentials) => {
    // 调用 API
  }

  const logout = () => {
    // 清除状态
  }

  const isAuthenticated = computed(() => !!authStore.user)

  return {login, logout, isAuthenticated}
}
```

### 页面权限

```typescript
export function usePageAccess() {
  const userStore = useUserStore()

  const hasAccess = (permission: string): boolean => {
    return userStore.permissions.includes(permission)
  }

  const canAccess = (path: string): boolean => {
    // 检查路由权限
  }

  return {hasAccess, canAccess}
}
```

### 业务数据加载

```typescript
export function useUserList() {
  const users = ref<User[]>([])
  const loading = ref(false)
  const total = ref(0)

  const fetchUsers = async (page: number, pageSize: number) => {
    loading.value = true
    try {
      const response = await getUsersAPI({page, pageSize})
      users.value = response.data
      total.value = response.total
    } finally {
      loading.value = false
    }
  }

  return {users, loading, total, fetchUsers}
}
```

## 使用指南

### ✅ 何时使用 Hooks

- 处理项目特定的业务逻辑
- 组合多个 Pinia stores
- 调用业务 API
- 处理复杂的应用状态

### ❌ 何时不使用 Hooks

- 通用的 DOM 操作 → 使用 Composables
- 纯逻辑函数 → 使用 utils
- 简单的数据转换 → 使用 computed

## 最佳实践清单

- [ ] 使用 `use` 前缀命名
- [ ] 提供完整的 TypeScript 类型
- [ ] 编写清晰的 JSDoc 注释
- [ ] 返回值结构化（避免数组）
- [ ] 处理错误和加载状态
- [ ] 清理副作用（onUnmounted）
- [ ] 编写单元测试
- [ ] 将逻辑与组件分离

## 常见问题

### Q: Hook 和 Composable 有什么区别？

A: Hooks 是项目特定的，通常依赖应用状态；Composables 是通用的，可以在任何项目中使用。

### Q: 为什么要分离 hooks 和 composables？

A: 便于维护和复用。Composables 可以提取为独立的库，而 Hooks 则是项目特定的。

### Q: Hook 中如何处理副作用？

A: 在 `onUnmounted` 或 `onBeforeUnmount` 中清理，使用 `watch` 的返回函数停止观察。

## 参考链接

- [Vue 3 生命周期钩子](https://vuejs.org/guide/essentials/lifecycle.html)
- [Vue 3 watchEffect](https://vuejs.org/api/reactivity-core.html#watcheffect)
- [Vue 3 Composition API](https://vuejs.org/guide/extras/composition-api-faq.html)

