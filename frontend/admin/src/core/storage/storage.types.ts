/**
 * 存储类型
 */
export type StorageType = "localStorage" | "sessionStorage";

/**
 * 缓存驱逐策略
 * - lru: 最近最少使用 (Least Recently Used)
 * - lfu: 最不经常使用 (Least Frequently Used)
 * - hybrid: 混合策略（结合 LRU 和 LFU）
 */
export type EvictionStrategy = "lru" | "lfu" | "hybrid";

/**
 * 存储元数据
 * @description 用于跟踪每个存储项的状态信息，支持智能缓存管理
 */
export interface StorageMeta {
  /** 过期时间戳 (null=永久有效) */
  e: number | null;
  /** 原始数据大小 (bytes) */
  s: number;
  /** 访问次数 (用于 LFU 策略) */
  a: number;
  /** 最后访问时间戳 (用于 LRU 策略) */
  t: number;
  /** 是否已压缩 */
  c: boolean;
}

/**
 * 存储性能指标
 * @description 用于监控存储使用情况和分析性能瓶颈
 */
export interface StorageMetrics {
  /** 命中次数 */
  hits: number;
  /** 未命中次数 */
  misses: number;
  /** 过期清理次数 */
  expirations: number;
  /** 设置操作次数 */
  sets: number;
  /** 错误次数 */
  errors: number;
  /** 配额使用量 (bytes) */
  quotaUsageBytes: number;
  /** 命中率 (0-1) */
  hitRate: number;
}

/**
 * 多标签页同步消息
 * @description 用于在不同标签页之间同步存储状态
 */
export interface SyncMessage {
  /** 操作类型 */
  type: "set" | "remove" | "clear";
  /** 存储键名 */
  key: string;
  /** 标签页 ID */
  tabId: string;
  /** 时间戳 */
  ts: number;
}

/**
 * 存储管理器配置选项
 * @description 用于初始化 StorageManager 实例时的配置参数
 */
export interface StorageManagerOptions {
  /** 存储键名前缀，用于隔离不同模块的数据 */
  prefix?: string;
  /** 存储类型：localStorage 或 sessionStorage */
  storageType?: StorageType;
  /** 触发压缩的阈值 (bytes)，默认 1024 */
  compressionThreshold?: number;
  /** 缓存驱逐策略，默认 "hybrid" */
  evictionStrategy?: EvictionStrategy;
  /** 最大存储项数 (0=不限)，默认 0 */
  maxItems?: number;
  /** 最大使用量 (MB)，默认 10 */
  maxUsageMB?: number;
  /** 默认过期时间 (ms)，null 表示永不过期 */
  defaultTTL?: number | null;
  /** 是否开启多标签页同步，默认 true */
  enableSync?: boolean;
  /** 指标变化回调函数 */
  onMetrics?: (metrics: StorageMetrics) => void;
  /** 同步事件回调函数 */
  onSync?: (event: SyncMessage) => void;
}

/**
 * 存储项结构（旧版，保留用于兼容）
 * @deprecated 请使用 StorageValue 替代
 */
export interface StorageItem<T> {
  /** 过期时间戳 */
  expiry?: number;
  /** 存储的值 */
  value: T;
}

/**
 * 存储值结构
 * @description 封装数据和过期信息的标准格式
 */
export interface StorageValue<T> {
  /** 存储的实际数据 */
  data: T;
  /** 过期时间戳 (null=永不过期) */
  expiry: number | null;
}

/**
 * 批量设置项的配置
 * @description 用于在批量设置时为每个项指定独立的 TTL
 */
export interface BatchSetItem<T> {
  /** 要存储的值 */
  value: T;
  /** 过期时间 (ms)，覆盖默认 TTL */
  ttl?: number | null;
}

/**
 * 批量获取的选项
 * @description 控制批量获取操作的行为
 */
export interface BatchGetOptions {
  /** 是否跳过过期项（默认 true，自动清理过期数据） */
  skipExpired?: boolean;
  /** 是否返回原始解析结果（用于调试，默认 false） */
  includeRaw?: boolean;
}

/**
 * 批量操作结果
 * @description 包含批量操作的执行结果统计
 */
export interface BatchResult<T> {
  /** 成功设置的项 */
  success: Record<string, T>;
  /** 失败的项及对应的错误信息 */
  failed: Record<string, Error>;
  /** 因过期被清理的键列表 */
  expired: string[];
}

/**
 * 批量操作进度回调
 * @description 用于实时监控批量操作的执行进度
 */
export type BatchProgressCallback = (progress: {
  /** 当前已处理的项目数 */
  current: number;
  /** 总项目数 */
  total: number;
  /** 当前处理的键名 */
  key: string;
  /** 处理状态 */
  status: "success" | "failed" | "expired";
}) => void;

/**
 * 存储缓存接口
 * @description 定义存储管理器的基本操作契约
 */
export interface IStorageCache {
  /** 清除所有存储项 */
  clear(): void;

  /** 清除所有过期项 */
  clearExpiredItems(): void;

  /**
   * 获取指定索引的键名
   * @param index - 索引位置
   * @returns 键名或 null
   */
  key(index: number): string | null;

  /**
   * 获取存储项数量
   * @returns 存储项总数
   */
  length(): number;

  /**
   * 移除指定键的存储项
   * @param key - 存储键名
   */
  removeItem(key: string): void;

  /**
   * 设置存储项
   * @param key - 存储键名
   * @param value - 要存储的值
   * @param ttl - 过期时间 (ms)，可选
   */
  setItem<T>(key: string, value: T, ttl?: number | null): void;

  /**
   * 获取存储项
   * @param key - 存储键名
   * @param defaultValue - 默认值，当键不存在时返回
   * @returns 存储的值或默认值
   */
  getItem<T>(key: string, defaultValue?: T | null): T | null;

  /**
   * 检查键是否存在且未过期
   * @param key - 存储键名
   * @returns 是否存在
   */
  hasItem(key: string): boolean;
}
