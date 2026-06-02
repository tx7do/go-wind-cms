# 存储管理模块 (StorageManager)

## 概述

本模块提供了一套增强型的浏览器本地存储管理方案，在原生 `localStorage` / `sessionStorage` 的基础上封装了 **TTL 过期控制**、**智能缓存驱逐**、**批量操作**、**多标签页同步** 和 **性能指标监控** 等能力。

### 核心特性

- **自动过期**：支持全局默认 TTL 和单条覆盖，过期数据读取时自动清理
- **缓存驱逐**：内置 LRU / LFU / Hybrid 三种策略，超限时自动淘汰低优先级数据
- **批量操作**：`getItems` / `setItems` / `removeItems` / `hasItems` 批量接口，大数量处理时自动让出主线程
- **多标签页同步**：基于 `BroadcastChannel` 实现跨 Tab 实时同步
- **性能监控**：内置命中/未命中/过期/错误等指标采集，支持回调上报
- **配额保护**：`QuotaExceededError` 时自动清理过期项并重试
- **键名隔离**：通过 `prefix` 前缀机制实现模块间数据隔离
- **SSR 兼容**：在 `window` 不存在的环境下安全降级

---

## 模块结构

```
src/core/storage/
├── index.ts              # 统一导出
├── storage.class.ts      # StorageManager 类实现
├── storage.types.ts      # 所有 TypeScript 类型定义
└── README.md             # 本文档
```

### 导入方式

```typescript
// 从 core 层统一导入（推荐）
import { StorageManager } from '@/core';

// 或从模块直接导入
import { StorageManager } from '@/core/storage';
import type { StorageManagerOptions, StorageMetrics } from '@/core/storage';
```

---

## 快速开始

### 创建实例

```typescript
import { StorageManager } from '@/core';

// 最简用法 — 使用 localStorage，无前缀
const storage = new StorageManager();

// 指定前缀和过期时间
const userCache = new StorageManager({
  prefix: 'user',           // 键名前缀：实际存储为 "user:xxx"
  defaultTTL: 30 * 60 * 1000,  // 默认 30 分钟过期
});

// 使用 sessionStorage
const sessionCache = new StorageManager({
  prefix: 'session',
  storageType: 'sessionStorage',
});
```

### 基本读写

```typescript
// 写入（永久有效）
userCache.setItem('profile', { name: '张三', age: 28 });

// 写入（指定 TTL，单位毫秒）
userCache.setItem('token', 'abc123', 2 * 60 * 60 * 1000); // 2 小时后过期

// 读取
const profile = userCache.getItem<{ name: string; age: number }>('profile');
// → { name: '张三', age: 28 }

// 读取（键不存在或已过期时返回默认值）
const token = userCache.getItem<string>('token', 'default-token');

// 检查是否存在（且未过期）
if (userCache.hasItem('profile')) {
  // ...
}

// 删除
userCache.removeItem('token');

// 清除当前前缀下的所有项
userCache.clear();
```

> **注意**：`getItem` 返回值类型为 `T | null`。如果数据已过期或不存在，返回 `null`（或你传入的 `defaultValue`）。无需手动检查过期状态。

---

## 配置选项

创建 `StorageManager` 实例时传入 `StorageManagerOptions`：

```typescript
const storage = new StorageManager({
  prefix: 'my-module',         // 键名前缀，用于模块隔离
  storageType: 'localStorage', // 'localStorage' | 'sessionStorage'
  compressionThreshold: 1024,  // 压缩阈值（bytes），当前版本预留
  evictionStrategy: 'hybrid',  // 'lru' | 'lfu' | 'hybrid'
  maxItems: 0,                 // 最大存储项数（0 = 不限）
  maxUsageMB: 10,              // 最大使用量（MB）
  defaultTTL: null,            // 默认过期时间（ms），null = 永不过期
  enableSync: true,            // 是否开启多标签页同步
  onMetrics: (m) => {},        // 性能指标变化回调
  onSync: (msg) => {},         // 跨标签页同步事件回调
});
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `prefix` | `string` | `""` | 键名前缀，避免不同模块数据冲突 |
| `storageType` | `'localStorage' \| 'sessionStorage'` | `'localStorage'` | 底层存储引擎 |
| `compressionThreshold` | `number` | `1024` | 触发压缩的阈值（bytes），当前版本为预留参数 |
| `evictionStrategy` | `'lru' \| 'lfu' \| 'hybrid'` | `'hybrid'` | 缓存驱逐策略 |
| `maxItems` | `number` | `0` | 最大存储项数，`0` 表示不限制 |
| `maxUsageMB` | `number` | `10` | 最大使用量（MB） |
| `defaultTTL` | `number \| null` | `null` | 默认 TTL（毫秒），`null` 为永不过期 |
| `enableSync` | `boolean` | `true` | 是否启用 BroadcastChannel 多标签页同步 |
| `onMetrics` | `(metrics: StorageMetrics) => void` | - | 指标变化回调（防抖 500ms） |
| `onSync` | `(event: SyncMessage) => void` | - | 同步事件回调 |

---

## 缓存驱逐策略

当存储项数或容量超过 `maxItems` / `maxUsageMB` 限制时，自动触发驱逐：

| 策略 | 说明 | 适用场景 |
|------|------|----------|
| `lru` | 淘汰最久未被访问的数据 | 侧重访问时效性 |
| `lfu` | 淘汰访问次数最少的数据 | 侧重访问频率 |
| `hybrid` | 频率权重 60% + 新鲜度 40% | **默认推荐**，综合平衡 |

驱逐优先级：
1. **已过期的数据**优先被淘汰
2. 按策略计算权重分数，分数低的优先淘汰

---

## 批量操作

适用于需要一次性处理大量键值对的场景，每处理 50 项自动让出主线程避免阻塞 UI。

### 批量获取

```typescript
const result = await storage.getItems<string>(['key1', 'key2', 'key3'], {
  skipExpired: true,   // 跳过过期项并自动清理（默认 true）
  includeRaw: false,   // 返回原始解析结果（调试用，默认 false）
});
// result: { key1: 'value1', key2: null, key3: 'value3' }
```

### 批量设置

```typescript
// 简单值
const result = await storage.setItems({
  key1: 'value1',
  key2: 'value2',
});

// 带独立 TTL
const result = await storage.setItems< string>({
  token: { value: 'abc', ttl: 7200000 },      // 2 小时
  temp:  { value: 'xyz', ttl: 60000 },         // 1 分钟
  perm:  'persistent',                          // 使用默认 TTL
});

// result.success → { token: 'abc', temp: 'xyz', perm: 'persistent' }
// result.failed  → { key: Error }  失败的项
// result.expired → ['key']         被清理的过期键
```

### 批量移除

```typescript
const { removed, failed } = await storage.removeItems(['key1', 'key2']);
// removed: ['key1', 'key2']
// failed:  { key3: Error }
```

### 批量检查存在

```typescript
const exists = await storage.hasItems(['key1', 'key2', 'key3']);
// { key1: true, key2: false, key3: true }
```

### 进度回调

批量操作支持 `onProgress` 回调实时监控执行进度：

```typescript
await storage.setItems(items, (progress) => {
  console.log(`[${progress.current}/${progress.total}] ${progress.key}: ${progress.status}`);
  // status: 'success' | 'failed' | 'expired'
});
```

---

## 多标签页同步

当 `enableSync: true`（默认开启）时，同一前缀的 `StorageManager` 实例会通过 `BroadcastChannel` 自动同步：

- **同一浏览器**的不同标签页共享同一份存储数据
- 某一标签页执行 `setItem` / `removeItem` / `clear` 后，其他标签页会收到同步通知
- 监听同步事件：

```typescript
const storage = new StorageManager({
  prefix: 'app',
  onSync: (msg) => {
    console.log(`Tab ${msg.tabId} did ${msg.type} on key "${msg.key}"`);
  },
});

// 或动态注册
const unsubscribe = storage.onSync((msg) => {
  // 处理同步事件
});
// 取消监听
unsubscribe();
```

---

## 性能指标

内置指标采集系统，用于监控存储的健康状况：

```typescript
const storage = new StorageManager({
  prefix: 'app',
  onMetrics: (metrics) => {
    console.log('命中率:', (metrics.hitRate * 100).toFixed(1) + '%');
    console.log('配额用量:', metrics.quotaUsageBytes, 'bytes');
  },
});

// 手动获取当前指标
const metrics = storage.getMetrics();
// {
//   hits: 120,              // 命中次数
//   misses: 15,             // 未命中次数
//   expirations: 3,         // 过期清理次数
//   sets: 200,              // 设置操作次数
//   errors: 0,              // 错误次数
//   quotaUsageBytes: 51200, // 配额使用量 (bytes)
//   hitRate: 0.889,         // 命中率 (0-1)
// }

// 重置指标
storage.resetMetrics();
```

> 指标回调内置 500ms 防抖，避免高频操作时频繁触发。

---

## API 参考

### 基础操作（IStorageCache 接口）

| 方法 | 签名 | 说明 |
|------|------|------|
| `getItem` | `getItem<T>(key, defaultValue?): T \| null` | 获取存储项，过期返回 `null` |
| `setItem` | `setItem<T>(key, value, ttl?): void` | 设置存储项，可选 TTL |
| `removeItem` | `removeItem(key): void` | 移除指定项 |
| `hasItem` | `hasItem(key): boolean` | 是否存在且未过期 |
| `clear` | `clear(): void` | 清除当前前缀下所有项 |
| `clearExpiredItems` | `clearExpiredItems(): void` | 清除所有过期项 |
| `key` | `key(index): string \| null` | 获取指定索引的键名（去前缀） |
| `length` | `length(): number` | 获取当前前缀下的有效项数 |

### 批量操作

| 方法 | 签名 | 说明 |
|------|------|------|
| `getItems` | `getItems<T>(keys, options?, onProgress?): Promise<Record<string, T \| null>>` | 批量获取 |
| `setItems` | `setItems<T>(items, onProgress?): Promise<BatchResult<T>>` | 批量设置 |
| `removeItems` | `removeItems(keys): Promise<{ removed, failed }>` | 批量移除 |
| `hasItems` | `hasItems(keys): Promise<Record<string, boolean>>` | 批量检查存在 |

### 高级方法

| 方法 | 签名 | 说明 |
|------|------|------|
| `getMetrics` | `getMetrics(): StorageMetrics` | 获取当前性能指标快照 |
| `resetMetrics` | `resetMetrics(): void` | 重置所有指标归零 |
| `onSync` | `onSync(callback): () => void` | 注册同步事件监听，返回取消函数 |
| `getRawStorage` | `getRawStorage(): Storage \| null` | 获取底层 Storage 实例（高级用途） |

---

## 典型场景

### 场景一：字典数据缓存

缓存接口返回的字典数据，设置 10 分钟过期，避免频繁请求：

```typescript
const dictCache = new StorageManager({
  prefix: 'dict',
  defaultTTL: 10 * 60 * 1000,
});

// 写入
dictCache.setItem('user-status', [{ label: '启用', value: 1 }]);

// 读取（过期自动返回 null）
const statusOptions = dictCache.getItem<{ label: string; value: number }[]>('user-status');
if (!statusOptions) {
  // 缓存失效，重新请求接口
}
```

### 场景二：模块级独立存储

不同业务模块使用不同前缀，互不干扰：

```typescript
const userStorage = new StorageManager({ prefix: 'user' });
const roleStorage = new StorageManager({ prefix: 'role' });

userStorage.setItem('list', [...]);  // 实际键名: "user:list"
roleStorage.setItem('list', [...]);  // 实际键名: "role:list"

userStorage.clear();  // 只清除 "user:" 前缀的项，不影响 "role:"
```

### 场景三：配额受限的临时缓存

限制缓存数量和容量，自动淘汰旧数据：

```typescript
const tempCache = new StorageManager({
  prefix: 'temp',
  maxItems: 100,
  maxUsageMB: 2,
  evictionStrategy: 'lru',
  defaultTTL: 5 * 60 * 1000,
});

// 超过 100 项或 2MB 时，最久未访问的数据会被自动淘汰
```

---

## 常见问题

### Q1: TTL 过期是被动清理还是主动清理？

**被动清理**。数据在 `getItem` 时检查是否过期，过期则返回 `null` 并删除。不会启动定时器主动扫描。如需手动触发全量清理过期数据，调用 `clearExpiredItems()`。

### Q2: 多个 StorageManager 实例的 prefix 可以相同吗？

可以。相同 prefix 的实例共享相同的键空间，且 BroadcastChannel 同步也是按 prefix 隔离的。通常建议一个模块只创建一个实例。

### Q3: 和直接使用 localStorage 有什么区别？

| 特性 | 原生 localStorage | StorageManager |
|------|-------------------|----------------|
| 过期控制 | 需手动实现 | 内置 TTL |
| 数据隔离 | 需手动拼前缀 | prefix 自动管理 |
| 类型安全 | 需手动 JSON.parse | 泛型自动推导 |
| 配额溢出 | 抛出异常 | 自动清理 + 重试 |
| 跨 Tab 同步 | 需手动监听 storage 事件 | 内置 BroadcastChannel |
| 性能监控 | 无 | 内置指标采集 |

### Q4: 存储的数据格式是什么？

数据以 `StorageValue<T>` 格式序列化存储：

```json
{
  "data": { "name": "张三" },
  "expiry": 1716873600000
}
```

`expiry` 为 `null` 表示永不过期。元数据（访问次数、最后访问时间等）存储在独立的 `__meta:` 前缀键中，与业务数据分离。

### Q5: 在 SSR 环境下安全吗？

安全。当 `window` 不存在时，构造函数不会访问 `localStorage`，所有操作静默降级（`getItem` 返回 `defaultValue`，`setItem` 无操作）。

---

## 注意事项

1. **前缀隔离**：`clear()` 只清除当前前缀下的数据，不会影响其他模块。未设置 `prefix` 时清除所有数据。
2. **泛型类型**：`getItem<T>` 的泛型参数仅在编译期提供类型提示，运行时不做类型校验。请确保存入和取出的类型一致。
3. **配额保护**：写入时如果触发 `QuotaExceededError`，会自动执行一次过期清理后重试。如果仍然失败，会静默记录错误指标。
4. **驱逐异步化**：缓存驱逐通过 `setTimeout` 异步执行，不会阻塞当前的 `setItem` 操作。
5. **序列化限制**：值通过 `JSON.stringify` / `JSON.parse` 处理，不支持 `undefined`、`Function`、`Symbol` 等不可序列化的值。
6. **批量操作**：返回 `Promise`，需要 `await`。批量设置中的值可以是普通值或 `{ value, ttl }` 配置对象。

---

## 相关文件

- [core/index.ts](../index.ts) - 统一导出入口
- [storage.class.ts](./storage.class.ts) - StorageManager 类实现
- [storage.types.ts](./storage.types.ts) - 类型定义
