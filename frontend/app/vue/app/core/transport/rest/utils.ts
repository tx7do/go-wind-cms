import type { HttpResponse } from './types';

/**
 * 分页参数类型
 */
export interface Paging {
  page?: number;
  pageSize?: number;
}

/**
 * 从对象中移除指定字段
 * @param obj 原始对象
 * @param keys 要移除的字段名（支持单个或多个）
 * @returns 移除字段后的新对象
 */
export function omit<T extends Record<string, any>, K extends keyof T>(
  obj: T,
  ...keys: (K | string)[]
): Omit<T, K> {
  if (!obj || typeof obj !== 'object') return {} as Omit<T, K>;
  const keySet = new Set<string>(keys.map((k) => String(k)));
  const result: Record<string, any> = {};
  for (const key of Object.keys(obj)) {
    if (!keySet.has(key)) {
      result[key] = obj[key];
    }
  }
  return result as Omit<T, K>;
}

/**
 * 创建列表查询 JSON 过滤字符串
 * @param formValues 查询表单值
 * @param needCleanTenant 是否清理租户字段
 * @returns JSON 字符串或 undefined
 */
export function makeQueryString(
  formValues?: null | object | Record<string, unknown>,
  needCleanTenant: boolean = false,
): string | undefined {
  if (formValues === null || formValues === undefined) {
    return undefined;
  }

  // 去除掉空值
  const cleaned = Object.fromEntries(
    Object.entries(formValues).filter(
      ([, v]) => v !== null && v !== undefined && v !== '',
    ),
  );

  if (Object.keys(cleaned).length === 0) {
    return undefined;
  }

  if (needCleanTenant) {
    // 删除租户相关字段 tenant_id 和 tenantId
    const { tenant_id, tenantId, ...rest } = cleaned as Record<string, unknown>;
    if (Object.keys(rest).length === 0) {
      return undefined;
    }
    return JSON.stringify(rest);
  }

  return JSON.stringify(cleaned);
}

/**
 * 创建排序字符串
 * @param orderBy 排序字段数组（带 - 前缀表示降序）
 * @returns JSON 字符串或 undefined
 */
export function makeOrderBy(orderBy?: null | string[]): string | undefined {
  const arr = orderBy ?? ['-created_at'];
  return JSON.stringify(arr) ?? undefined;
}


/**
 * 创建更新掩码字符串
 * @param keys - 字段键名数组
 * @returns 逗号分隔的字符串
 */
export function makeUpdateMask(keys: string[]): string {
  return [...keys, 'id'].join(',');
}

/**
 * 默认的请求 ID 生成器
 */
export function defaultIdGenerator(): string {
  try {
    // 优先使用标准 API
    const rnd = (
      globalThis as unknown as { crypto?: { randomUUID?: () => string } }
    )?.crypto?.randomUUID?.();
    if (typeof rnd === 'string' && rnd.length > 0) return rnd;
  } catch {
    // ignore
  }
  // 降级方案
  return Math.random().toString(36).slice(2) + Date.now().toString(36);
}

/**
 * 默认的错误消息提取（不依赖 i18n，纯逻辑 fallback）
 * 按优先级：reason → message → status code → 兜底
 * 如需 i18n 翻译，通过 RequestClientCallbacks.getErrorMsg 注入
 */
export function getDefaultErrorMsg(error: unknown): string {
  // 网络错误
  const errStr = String(error ?? '');
  if (errStr.includes('Network Error')) {
    return 'Network Error';
  }

  // 超时
  if (
    error &&
    typeof error === 'object' &&
    'message' in error &&
    String(error.message).includes('timeout')
  ) {
    return 'Request Timeout';
  }

  // 获取后端返回数据
  const resData =
    error &&
    typeof error === 'object' &&
    'response' in error &&
    error.response &&
    typeof error.response === 'object' &&
    'data' in error.response
      ? (error.response.data as HttpResponse)
      : undefined;

  if (!resData) {
    return 'Unknown Error';
  }

  const { reason, message, code } = resData;

  // 1. 优先使用 reason
  if (reason) {
    return reason;
  }

  // 2. 使用后端 message
  if (message?.trim()) {
    return message.trim();
  }

  // 3. 使用 code
  if (code) {
    return `Error ${code}`;
  }

  // 4. 兜底
  return 'Unknown Error';
}

/**
 * 绑定实例方法
 * @param instance 实例
 */
export function bindMethods<T extends object>(instance: T): void {
  const prototype = Object.getPrototypeOf(instance);
  const propertyNames = Object.getOwnPropertyNames(prototype);

  propertyNames.forEach((propertyName) => {
    const descriptor = Object.getOwnPropertyDescriptor(prototype, propertyName);
    const propertyValue = instance[propertyName as keyof T];

    if (
      typeof propertyValue === 'function' &&
      propertyName !== 'constructor' &&
      descriptor &&
      !descriptor.get &&
      !descriptor.set
    ) {
      instance[propertyName as keyof T] = propertyValue.bind(instance);
    }
  });
}
