/**
 * 消息分类模块枚举映射常量
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 布尔值映射（isEnabled） ==========

export function enableBoolToColor(value: boolean): string {
  return value ? 'success' : 'error';
}

export function getEnableBoolMap(t: TFn) {
  return {
    true: { text: t('yes'), color: 'success' },
    false: { text: t('no'), color: 'error' },
  };
}
