# Composables (可组合函数)

## 概述

本目录包含可复用的 Vue 3 Composition API 函数，也称为 Composables。

### 目录结构

```
composables/
├── dark.ts           # 深色模式切换
├── useAutoThemeSwitcher.ts  # 自动主题切换
└── README.md         # 本文档
```

## 什么是 Composables？

Composables 是 Vue 3 中的一种复用逻辑的方式，使用 Composition API 封装可复用的功能。

**优势**：
- ✅ 高度可复用
- ✅ 易于测试
- ✅ 易于理解和维护
- ✅ 避免深度嵌套

## 使用方法

### 基本模式

```typescript
import { ref, computed } from 'vue'

export function useCounter(initialValue = 0) {
  const count = ref(initialValue)
  
  const double = computed(() => count.value * 2)
  
  const increment = () => count.value++
  const decrement = () => count.value--
  
  return {
    count,
    double,
    increment,
    decrement
  }
}
```

### 在组件中使用

```vue
<script setup>
import { useCounter } from '@/composables/useCounter'

const { count, double, increment, decrement } = useCounter(10)
</script>

<template>
  <div>
    <p>Count: {{ count }}</p>
    <p>Double: {{ double }}</p>
    <button @click="increment">+1</button>
    <button @click="decrement">-1</button>
  </div>
</template>
```

## 现有 Composables

### dark.ts
处理深色模式逻辑：
```typescript
import { preferredDark } from '@/composables/dark'

// 获取系统偏好的深色模式
const isDark = preferredDark.value
```

### useAutoThemeSwitcher.ts
自动主题切换功能：
```typescript
import { useAutoThemeSwitcher } from '@/composables'

const { initializeThemeSwitcher } = useAutoThemeSwitcher(appStore)

// 初始化主题切换
initializeThemeSwitcher()
```

## 编写 Composables 的最佳实践

### 1. 清晰的命名
- 使用 `use` 前缀
- 准确描述功能
- 例如：`useLocalStorage`、`useFetch`、`useForm`

### 2. 清晰的入参和返回值
```typescript
/**
 * 获取用户列表
 * @param pageSize - 每页条数
 * @param onError - 错误回调
 * @returns 用户列表和加载状态
 */
export function useUserList(pageSize = 10, onError?: (error: Error) => void) {
  const users = ref<User[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)
  
  const fetchUsers = async (page: number) => {
    loading.value = true
    try {
      // ...
    } catch (e) {
      error.value = e as Error
      onError?.(error.value)
    } finally {
      loading.value = false
    }
  }
  
  return {
    users,
    loading,
    error,
    fetchUsers
  }
}
```

### 3. 正确处理生命周期
```typescript
import { onMounted, onUnmounted } from 'vue'

export function useEventListener(target: EventTarget, event: string, callback: EventListener) {
  onMounted(() => {
    target.addEventListener(event, callback)
  })
  
  onUnmounted(() => {
    target.removeEventListener(event, callback)
  })
}
```

### 4. 类型安全
```typescript
import type { Ref } from 'vue'

interface User {
  id: number
  name: string
  email: string
}

export function useUser(userId: number): {
  user: Ref<User | null>
  loading: Ref<boolean>
  error: Ref<Error | null>
} {
  // 实现...
}
```

### 5. 避免重复订阅
```typescript
export function useWatch(source: any, callback: any) {
  const stopWatch = watch(source, callback)
  
  onUnmounted(() => {
    stopWatch() // 组件卸载时停止观察
  })
  
  return { stopWatch }
}
```

## Composables vs Hooks

| 特性 | Composables | Hooks |
|------|-----------|-------|
| 存储位置 | `composables/` | `hooks/` |
| 使用场景 | 通用可复用逻辑 | Vue 特定的逻辑 |
| 依赖关系 | 轻量级 | 可能依赖业务 |

## 常见 Composables 模式

### 获取数据
```typescript
export function useFetch(url: string) {
  const data = ref(null)
  const loading = ref(true)
  const error = ref(null)
  
  fetch(url)
    .then(r => r.json())
    .then(d => data.value = d)
    .catch(e => error.value = e)
    .finally(() => loading.value = false)
  
  return { data, loading, error }
}
```

### 本地存储
```typescript
export function useLocalStorage(key: string, defaultValue = null) {
  const value = ref(localStorage.getItem(key) ?? defaultValue)
  
  watch(value, (newValue) => {
    localStorage.setItem(key, JSON.stringify(newValue))
  })
  
  return value
}
```

### 表单验证
```typescript
export function useForm(initialValues: any, onSubmit: (values: any) => void) {
  const values = reactive(initialValues)
  const errors = reactive({})
  
  const validate = () => {
    // 验证逻辑
  }
  
  const handleSubmit = () => {
    if (validate()) {
      onSubmit(values)
    }
  }
  
  return { values, errors, handleSubmit }
}
```

## 项目现存 Composables

| 名称 | 功能 | 位置 |
|------|------|------|
| `dark` | 深色模式检测 | `composables/dark.ts` |
| `useAutoThemeSwitcher` | 自动主题切换 | `composables/useAutoThemeSwitcher.ts` |

## 参考链接

- [Vue 3 Composition API](https://vuejs.org/guide/extras/composition-api-faq.html)
- [VueUse - 常用 Composables 库](https://vueuse.org/)
- [推荐的 Composables 命名规范](https://antfu.me/posts/composable-vue)

