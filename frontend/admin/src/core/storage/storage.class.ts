import type {
  BatchGetOptions,
  BatchProgressCallback,
  BatchResult,
  BatchSetItem,
  IStorageCache,
  StorageManagerOptions,
  StorageMeta,
  StorageMetrics,
  StorageValue,
  SyncMessage,
} from "./storage.types";

class StorageManager implements IStorageCache {
  private readonly prefix: string | undefined;
  private readonly storage: Storage;
  private readonly options: Required<Omit<StorageManagerOptions, "onMetrics" | "onSync">>;
  private readonly defaultTTL?: number | null;
  private readonly tabId: string;
  private channel: BroadcastChannel | null = null;
  private metrics: StorageMetrics = {
    hits: 0,
    misses: 0,
    expirations: 0,
    sets: 0,
    errors: 0,
    quotaUsageBytes: 0,
    hitRate: 0,
  };
  private metricsDebounceTimer: number | null = null;
  private readonly onMetricsCallback?: (metrics: StorageMetrics) => void;
  private onSyncCallback?: (event: SyncMessage) => void;

  constructor(options: StorageManagerOptions = {}) {
    const {
      prefix = "",
      storageType = "localStorage",
      compressionThreshold = 1024,
      evictionStrategy = "hybrid",
      maxItems = 0,
      maxUsageMB = 10,
      defaultTTL = null,
      enableSync = true,
      onMetrics,
      onSync,
    } = options;

    this.prefix = prefix;
    this.options = {
      prefix,
      storageType,
      compressionThreshold,
      evictionStrategy,
      maxItems,
      maxUsageMB,
      defaultTTL,
      enableSync,
    };
    this.onMetricsCallback = onMetrics;
    this.onSyncCallback = onSync;
    this.tabId =
      typeof window !== "undefined"
        ? `${prefix}_${Math.random().toString(36).substring(2, 10)}`
        : "ssr";

    // SSR 兼容
    if (typeof window === "undefined") {
      this.storage = {} as Storage;
      return;
    }

    this.storage = storageType === "localStorage" ? window.localStorage : window.sessionStorage;
    this.defaultTTL = defaultTTL;

    if (enableSync) this.initSync();
    this.updateQuotaUsage().catch(() => {});
  }

  /**
   * 获取完整存储键（安全添加前缀）
   */
  private getFullKey(key: string): string {
    return this.prefix ? `${this.prefix}:${key}` : key;
  }

  private getMetaKey(key: string): string {
    return this.prefix ? `${this.prefix}:__meta:${key}` : `__meta:${key}`;
  }

  private readMeta(key: string): StorageMeta | null {
    const raw = this.storage.getItem(this.getMetaKey(key));
    if (!raw) return null;
    try {
      return JSON.parse(raw);
    } catch {
      return null;
    }
  }

  private writeMeta(key: string, meta: StorageMeta): void {
    try {
      this.storage.setItem(this.getMetaKey(key), JSON.stringify(meta));
    } catch (error) {
      console.warn("[StorageManager] Failed to write meta:", error);
    }
  }

  /**
   * 访问时更新元数据（访问次数和最后访问时间）
   */
  private updateMetaOnAccess(key: string): void {
    const meta = this.readMeta(key);
    if (meta) {
      meta.a += 1; // 增加访问次数
      meta.t = Date.now(); // 更新最后访问时间
      this.writeMeta(key, meta);
    } else {
      // 如果没有元数据，创建新的
      this.writeMeta(key, {
        e: null,
        s: 0,
        a: 1,
        t: Date.now(),
        c: false,
      });
    }
  }

  /**
   * 设置时更新元数据
   */
  private updateMetaOnSet(key: string, serialized: string, ttl?: number | null): void {
    const size = new Blob([serialized]).size;
    const expiry =
      ttl !== undefined && ttl !== null
        ? Date.now() + ttl
        : this.defaultTTL
          ? Date.now() + this.defaultTTL
          : null;

    const existingMeta = this.readMeta(key);
    const meta: StorageMeta = {
      e: expiry,
      s: size,
      a: existingMeta?.a || 0, // 保留访问次数
      t: Date.now(), // 更新最后修改时间
      c: false,
    };
    this.writeMeta(key, meta);
  }

  /**
   * 移除元数据
   */
  private removeMeta(key: string): void {
    try {
      this.storage.removeItem(this.getMetaKey(key));
    } catch (error) {
      console.warn("[StorageManager] Failed to remove meta:", error);
    }
  }

  private getAllMeta(): Map<string, StorageMeta> {
    const map = new Map<string, StorageMeta>();
    const prefix = this.prefix ? `${this.prefix}:__meta:` : "__meta:";

    for (let i = 0; i < this.storage.length; i++) {
      const k = this.storage.key(i);
      if (k?.startsWith(prefix)) {
        const shortKey = k.slice(prefix.length);
        const meta = this.readMeta(shortKey);
        if (meta) map.set(shortKey, meta);
      }
    }
    return map;
  }

  /**
   * 安全解析存储值
   */
  private parseValue<T>(raw: string | null): T | null {
    if (!raw) return null;
    try {
      const parsed: StorageValue<T> = JSON.parse(raw);
      // 检查过期
      if (parsed.expiry !== null && Date.now() > parsed.expiry) {
        this.trackMetric("expirations");
        return null;
      }
      return parsed.data;
    } catch (error) {
      console.error(`[StorageManager] Parse error:`, error);
      this.trackMetric("errors");
      return null;
    }
  }

  /**
   * 序列化存储值
   */
  private serializeValue<T>(data: T, ttl?: number | null): string {
    const expiry =
      ttl !== undefined && ttl !== null
        ? Date.now() + ttl
        : this.defaultTTL
          ? Date.now() + this.defaultTTL
          : null;

    const value: StorageValue<T> = { data, expiry };
    return JSON.stringify(value);
  }

  /**
   * 获取存储项
   */
  getItem<T>(key: string, defaultValue: T | null = null): T | null {
    if (!this.storage) return defaultValue;

    const fullKey = this.getFullKey(key);
    const raw = this.storage.getItem(fullKey);

    // 如果键不存在，记录未命中
    if (raw === null) {
      this.trackMetric("misses");
      return defaultValue;
    }

    const parsed = this.parseValue<T>(raw);

    // 如果解析失败或过期，清理脏数据
    if (parsed === null) {
      this.storage.removeItem(fullKey);
      this.removeMeta(key);
      // parseValue 已经追踪了 expirations 或 errors
      return defaultValue;
    }

    // 记录命中并更新元数据
    this.trackMetric("hits");
    this.updateMetaOnAccess(key);

    return parsed;
  }

  /**
   * 设置存储项
   */
  setItem<T>(key: string, value: T, ttl?: number | null): void {
    if (!this.storage) return;

    const fullKey = this.getFullKey(key);
    try {
      const serialized = this.serializeValue(value, ttl);
      this.storage.setItem(fullKey, serialized);

      // 更新元数据
      this.updateMetaOnSet(key, serialized, ttl);

      // 记录设置操作
      this.trackMetric("sets");

      // 检查是否需要驱逐
      this.checkAndEvict();

      // 广播同步消息
      this.broadcast({ type: "set", key, ts: Date.now() });
    } catch (error) {
      this.trackMetric("errors");
      // 处理 QuotaExceededError
      if (error instanceof DOMException && error.name === "QuotaExceededError") {
        console.warn("[StorageManager] Storage quota exceeded, clearing expired items...");
        this.clearExpiredItems();
        // 重试一次
        try {
          const serialized = this.serializeValue(value, ttl);
          this.storage.setItem(fullKey, serialized);
          this.updateMetaOnSet(key, serialized, ttl);
          this.trackMetric("sets");
          this.broadcast({ type: "set", key, ts: Date.now() });
        } catch (retryError) {
          console.error("[StorageManager] Failed to set item after cleanup:", retryError);
          this.trackMetric("errors");
        }
      } else {
        console.error("[StorageManager] Set item error:", error);
      }
    }
  }

  /**
   * 移除存储项
   */
  removeItem(key: string): void {
    if (!this.storage) return;
    const fullKey = this.getFullKey(key);
    this.storage.removeItem(fullKey);

    // 清理元数据
    this.removeMeta(key);

    // 广播同步消息
    this.broadcast({ type: "remove", key, ts: Date.now() });
  }

  /**
   * 检查键是否存在且未过期
   */
  hasItem(key: string): boolean {
    return this.getItem(key) !== null;
  }

  /**
   * 清除所有带前缀的项（使用倒序遍历避免索引错位）
   */
  clear(): void {
    if (!this.storage) return;

    const keysToRemove: string[] = [];

    // 倒序遍历，收集要删除的 key
    for (let i = this.storage.length - 1; i >= 0; i--) {
      const key = this.storage.key(i);
      if (key && (this.prefix === "" || key.startsWith(`${this.prefix}:`))) {
        // 排除元数据键，单独处理
        if (!key.includes(":__meta:")) {
          keysToRemove.push(key);
        }
        this.storage.removeItem(key);
      }
    }

    // 清理元数据
    keysToRemove.forEach((key) => {
      const shortKey =
        this.prefix && key.startsWith(`${this.prefix}:`) ? key.slice(this.prefix.length + 1) : key;
      this.removeMeta(shortKey);
    });

    // 广播同步消息
    this.broadcast({ type: "clear", key: "*", ts: Date.now() });
  }

  /**
   * 清除所有过期项
   */
  clearExpiredItems(): void {
    if (!this.storage) return;

    const keysToRemove: string[] = [];

    // 先收集要删除的 key，避免遍历中修改集合
    for (let i = 0; i < this.storage.length; i++) {
      const key = this.storage.key(i);
      if (!key) continue;

      // 跳过元数据键
      if (key.includes(":__meta:")) {
        continue;
      }

      // 只检查带前缀的项
      if (this.prefix && !key.startsWith(`${this.prefix}:`)) {
        continue;
      }

      const raw = this.storage.getItem(key);
      if (raw) {
        try {
          const parsed: StorageValue<unknown> = JSON.parse(raw);
          if (parsed.expiry !== null && Date.now() > parsed.expiry) {
            keysToRemove.push(key);
          }
        } catch {
          // 解析失败的脏数据也清理
          keysToRemove.push(key);
        }
      }
    }

    // 统一删除
    keysToRemove.forEach((key) => {
      this.storage?.removeItem(key);
      // 清理对应的元数据
      const shortKey =
        this.prefix && key.startsWith(`${this.prefix}:`) ? key.slice(this.prefix.length + 1) : key;
      this.removeMeta(shortKey);
      this.trackMetric("expirations");
    });
  }

  /**
   * 获取存储键（带前缀过滤）
   */
  key(index: number): string | null {
    if (!this.storage) return null;

    const rawKey = this.storage.key(index);
    if (!rawKey) return null;

    // 返回不带前缀的 key
    if (this.prefix && rawKey.startsWith(`${this.prefix}:`)) {
      return rawKey.slice(this.prefix.length + 1);
    }
    return this.prefix ? null : rawKey;
  }

  /**
   * 获取有效存储项数量（带前缀过滤）
   */
  length(): number {
    if (!this.storage) return 0;

    if (!this.prefix) {
      return this.storage.length;
    }

    let count = 0;
    const prefix = `${this.prefix}:`;
    for (let i = 0; i < this.storage.length; i++) {
      if (this.storage.key(i)?.startsWith(prefix)) {
        count++;
      }
    }
    return count;
  }

  /**
   * 获取底层 Storage 实例（用于高级操作）
   */
  getRawStorage(): Storage | null {
    return this.storage;
  }

  /**
   * 批量获取存储项
   * @param keys 要获取的键数组
   * @param options 获取选项
   * @param onProgress 进度回调（每处理完一个项触发）
   * @returns Promise<Record<string, T | null>> 键值映射
   */
  async getItems<T>(
    keys: string[],
    options: BatchGetOptions = {},
    onProgress?: BatchProgressCallback
  ): Promise<Record<string, T | null>> {
    const { skipExpired = true, includeRaw = false } = options;
    const result: Record<string, T | null> = {};

    if (!this.storage || keys.length === 0) {
      return result;
    }

    for (let i = 0; i < keys.length; i++) {
      const key = keys[i];
      const fullKey = this.getFullKey(key);

      try {
        const raw = this.storage.getItem(fullKey);

        if (!raw) {
          result[key] = null;
          this.trackMetric("misses");
          onProgress?.({ current: i + 1, total: keys.length, key, status: "failed" });
          continue;
        }

        const parsed: StorageValue<T> = JSON.parse(raw);

        // 检查过期
        if (parsed.expiry !== null && Date.now() > parsed.expiry) {
          this.trackMetric("expirations");
          if (skipExpired) {
            this.storage.removeItem(fullKey);
            this.removeMeta(key);
            result[key] = null;
            onProgress?.({ current: i + 1, total: keys.length, key, status: "expired" });
          } else {
            // 不过滤过期项，返回原始数据（调用方自行处理）
            result[key] = includeRaw ? (parsed as unknown as T) : null;
          }
          continue;
        }

        result[key] = parsed.data;
        this.trackMetric("hits");
        this.updateMetaOnAccess(key);
        onProgress?.({ current: i + 1, total: keys.length, key, status: "success" });
      } catch (error) {
        console.error(`[StorageManager] getItems error for key "${key}":`, error);
        this.trackMetric("errors");
        // 解析失败则清理脏数据
        this.storage.removeItem(fullKey);
        this.removeMeta(key);
        result[key] = null;
        onProgress?.({ current: i + 1, total: keys.length, key, status: "failed" });
      }

      // 微任务让出主线程，避免阻塞 UI
      if (i % 50 === 0) {
        await new Promise((resolve) => setTimeout(resolve, 0));
      }
    }

    return result;
  }

  /**
   * 批量设置存储项
   * @param items 键值对映射，值可以是普通值或 BatchSetItem 配置
   * @param onProgress 进度回调
   * @returns Promise<BatchResult<T>> 操作结果统计
   */
  async setItems<T>(
    items: Record<string, T | BatchSetItem<T>>,
    onProgress?: BatchProgressCallback
  ): Promise<BatchResult<T>> {
    const result: BatchResult<T> = {
      success: {},
      failed: {},
      expired: [],
    };

    if (!this.storage) {
      // 将所有项标记为失败
      for (const key of Object.keys(items)) {
        result.failed[key] = new Error("Storage not available");
      }
      return result;
    }

    const entries = Object.entries(items);

    for (let i = 0; i < entries.length; i++) {
      const [key, item] = entries[i];
      const fullKey = this.getFullKey(key);

      try {
        // 解析值配置
        const { value, ttl } = this._normalizeBatchValue(item);
        const serialized = this.serializeValue(value, ttl);

        this.storage.setItem(fullKey, serialized);
        this.updateMetaOnSet(key, serialized, ttl);
        result.success[key] = value;
        this.trackMetric("sets");
        onProgress?.({ current: i + 1, total: entries.length, key, status: "success" });
      } catch (error) {
        this.trackMetric("errors");
        const err = error instanceof Error ? error : new Error(String(error));

        // 处理配额超限：自动清理过期项后重试
        if (err.name === "QuotaExceededError") {
          console.warn(
            "[StorageManager] Quota exceeded during batch set, cleaning expired items..."
          );
          this.clearExpiredItems();

          try {
            const { value, ttl } = this._normalizeBatchValue(item);
            const serialized = this.serializeValue(value, ttl);
            this.storage.setItem(fullKey, serialized);
            this.updateMetaOnSet(key, serialized, ttl);
            result.success[key] = value;
            this.trackMetric("sets");
            onProgress?.({ current: i + 1, total: entries.length, key, status: "success" });
            continue;
          } catch (retryError) {
            result.failed[key] =
              retryError instanceof Error ? retryError : new Error("Retry failed");
            this.trackMetric("errors");
            onProgress?.({ current: i + 1, total: entries.length, key, status: "failed" });
            continue;
          }
        }

        result.failed[key] = err;
        onProgress?.({ current: i + 1, total: entries.length, key, status: "failed" });
      }

      // 每 50 项让出主线程
      if (i % 50 === 0) {
        await new Promise((resolve) => setTimeout(resolve, 0));
      }
    }

    // 检查是否需要驱逐
    this.checkAndEvict();

    // 广播同步消息
    this.broadcast({ type: "set", key: "*", ts: Date.now() });

    return result;
  }

  /**
   * 批量移除存储项
   * @param keys 要移除的键数组
   * @returns Promise<{ removed: string[]; failed: Record<string, Error> }>
   */
  async removeItems(keys: string[]): Promise<{ removed: string[]; failed: Record<string, Error> }> {
    const removed: string[] = [];
    const failed: Record<string, Error> = {};

    if (!this.storage) {
      for (const key of keys) {
        failed[key] = new Error("Storage not available");
      }
      return { removed, failed };
    }

    for (const key of keys) {
      try {
        const fullKey = this.getFullKey(key);
        this.storage.removeItem(fullKey);
        this.removeMeta(key);
        removed.push(key);
      } catch (error) {
        failed[key] = error instanceof Error ? error : new Error(String(error));
      }
    }

    // 广播同步消息
    if (removed.length > 0) {
      this.broadcast({ type: "remove", key: "*", ts: Date.now() });
    }

    return { removed, failed };
  }

  /**
   * 批量检查键是否存在（且未过期）
   * @param keys 要检查的键数组
   * @returns Promise<Record<string, boolean>>
   */
  async hasItems(keys: string[]): Promise<Record<string, boolean>> {
    const result: Record<string, boolean> = {};

    if (!this.storage) {
      for (const key of keys) result[key] = false;
      return result;
    }

    for (const key of keys) {
      result[key] = this.getItem(key) !== null;
    }

    return result;
  }

  /**
   * 内部方法：标准化批量设置的值配置
   */
  private _normalizeBatchValue<T>(item: T | BatchSetItem<T>): { value: T; ttl?: number | null } {
    if (item && typeof item === "object" && "value" in item) {
      const batchItem = item as BatchSetItem<T>;
      return { value: batchItem.value, ttl: batchItem.ttl };
    }
    return { value: item as T };
  }

  private checkAndEvict(): void {
    const { maxItems, maxUsageMB } = this.options;
    if (maxItems === 0 && maxUsageMB === 0) return;

    const allMeta = this.getAllMeta();
    let currentItems = allMeta.size;
    let currentUsage = Array.from(allMeta.values()).reduce((sum, m) => sum + m.s, 0);

    if (currentItems <= maxItems && currentUsage <= maxUsageMB * 1024 * 1024) return;

    // 排序策略
    const sortedKeys = this.sortByEvictionPriority(allMeta);
    let evicted = 0;
    const keysToEvict: string[] = [];

    for (const key of sortedKeys) {
      if (currentItems <= maxItems && currentUsage <= maxUsageMB * 1024 * 1024) break;
      keysToEvict.push(key);
      currentItems--;
      currentUsage -= allMeta.get(key)!.s;
      evicted++;
    }

    // 异步执行删除，避免阻塞主线程
    if (keysToEvict.length > 0) {
      setTimeout(() => {
        keysToEvict.forEach((k) => this.removeItem(k));
        console.log(
          `[StorageManager] Evicted ${evicted} items using ${this.options.evictionStrategy} strategy`
        );
      }, 0);
    }
  }

  private sortByEvictionPriority(metas: Map<string, StorageMeta>): string[] {
    const now = Date.now();
    const keys = Array.from(metas.keys());

    return keys.sort((a, b) => {
      const ma = metas.get(a)!;
      const mb = metas.get(b)!;

      // 1. 已过期的优先淘汰
      if (ma.e !== null && ma.e < now) return -1;
      if (mb.e !== null && mb.e < now) return 1;

      // 2. 按策略计算权重
      const strategy = this.options.evictionStrategy;
      let scoreA: number, scoreB: number;

      if (strategy === "lru") {
        scoreA = ma.t;
        scoreB = mb.t;
      } else if (strategy === "lfu") {
        scoreA = ma.a;
        scoreB = mb.a;
      } else {
        // hybrid: 频率权重 60%，新鲜度 40%
        const freqWeight = 0.6;
        const timeDecay = Math.log(now - ma.t + 1);
        scoreA = ma.a * freqWeight - timeDecay * (1 - freqWeight);
        const timeDecayB = Math.log(now - mb.t + 1);
        scoreB = mb.a * freqWeight - timeDecayB * (1 - freqWeight);
      }

      return scoreA - scoreB; // 分数低的先淘汰
    });
  }

  private initSync(): void {
    if (typeof window === "undefined" || !window.BroadcastChannel) return;

    this.channel = new BroadcastChannel(`__sm_sync__${this.prefix || "default"}`);
    this.channel.onmessage = (e: MessageEvent<SyncMessage>) => {
      const msg = e.data;
      if (msg.tabId === this.tabId) return; // 忽略自身发出的消息

      this.handleSyncMessage(msg);
      this.onSyncCallback?.(msg);
    };
  }

  private handleSyncMessage(msg: SyncMessage): void {
    // 同步当前实例状态，避免脏读
    if (msg.type === "remove" || msg.type === "clear") {
      // 元数据已随 removeItem 清理，无需额外操作
    }
    // set 操作已在 getItem 中通过元数据保证一致性
  }

  private broadcast(msg: Omit<SyncMessage, "tabId">): void {
    if (!this.channel || !this.options.enableSync) return;
    this.channel.postMessage({ ...msg, tabId: this.tabId });
  }

  onSync(callback: (event: SyncMessage) => void): () => void {
    this.onSyncCallback = callback;
    return () => {
      this.onSyncCallback = undefined;
    };
  }

  // ================= 监控埋点 =================

  private trackMetric(key: keyof StorageMetrics): void {
    (this.metrics as any)[key]++;
    this.updateHitRate();
    this.debounceReportMetrics();
  }

  private updateHitRate(): void {
    const total = this.metrics.hits + this.metrics.misses;
    this.metrics.hitRate = total === 0 ? 0 : this.metrics.hits / total;
  }

  private debounceReportMetrics(): void {
    if (typeof window === "undefined") return;
    if (this.metricsDebounceTimer) clearTimeout(this.metricsDebounceTimer);
    this.metricsDebounceTimer = window.setTimeout(() => {
      this.onMetricsCallback?.({ ...this.metrics });
    }, 500);
  }

  private async updateQuotaUsage(): Promise<void> {
    if (typeof window === "undefined" || !navigator.storage?.estimate) return;
    try {
      const est = await navigator.storage.estimate();
      this.metrics.quotaUsageBytes = est.usage || 0;
    } catch {
      // Fallback 估算
      let size = 0;
      for (let i = 0; i < this.storage.length; i++) {
        const k = this.storage.key(i);
        if (k) size += (k.length + (this.storage.getItem(k)?.length || 0)) * 2;
      }
      this.metrics.quotaUsageBytes = size;
    }
  }

  getMetrics(): StorageMetrics {
    return { ...this.metrics };
  }

  resetMetrics(): void {
    this.metrics = {
      hits: 0,
      misses: 0,
      expirations: 0,
      sets: 0,
      errors: 0,
      quotaUsageBytes: 0,
      hitRate: 0,
    };
  }
}

export { StorageManager };
