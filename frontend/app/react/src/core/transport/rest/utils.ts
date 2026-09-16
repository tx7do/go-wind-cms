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
 * 从请求层错误中提取 API 错误的 reason/code。
 * 兼容两类错误形态：
 * - AxiosError（拦截器链上，响应体在 error.response.data）
 * - 已被 request() 解包的响应体（{code, reason, message, metadata}）
 */
export function extractApiError(error: unknown): {
  code?: number | string;
  reason?: string;
} {
  if (!error || typeof error !== 'object') return {};
  const source = error as Record<string, any>;
  const data = source.response?.data ?? source;
  if (!data || typeof data !== 'object') return {};
  return {
    code: data.code,
    reason:
      typeof data.reason === 'string' && data.reason !== ''
        ? data.reason
        : undefined,
  };
}

/**
 * API 错误文案解析：不再使用后端 message 字段（开发用英文文本，无法本地化），
 * 改用 reason 查询 i18n 错误文案（error.<REASON> 命名空间），
 * 查不到时回退到默认错误文案（DEFAULT）。
 * @param error 请求层抛出的错误
 * @param dict 当前语言的 error 命名空间文案表
 */
export function resolveApiErrorMsg(
  error: unknown,
  dict?: Record<string, string>,
): string {
  // 网络错误
  const errStr = String(error ?? '');
  if (errStr.includes('Network Error')) {
    return dict?.NETWORK_ERROR ?? 'Network Error';
  }

  // 超时
  if (
    error &&
    typeof error === 'object' &&
    'message' in error &&
    String((error as Record<string, any>).message).includes('timeout')
  ) {
    return dict?.TIMEOUT ?? 'Request Timeout';
  }

  const { reason } = extractApiError(error);

  // reason 查 i18n 文案
  if (reason && dict?.[reason]) {
    return dict[reason] as string;
  }

  // 默认错误信息
  return dict?.DEFAULT ?? 'Unknown Error';
}

/**
 * 默认错误消息提取（无 i18n 字典时的兜底；
 * 业务侧应通过 RequestClientCallbacks.getErrorMsg 注入带字典的 resolveApiErrorMsg）
 */
export function getDefaultErrorMsg(error: unknown): string {
  return resolveApiErrorMsg(error);
}

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
