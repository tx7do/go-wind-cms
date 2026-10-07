/**
 * 语言模块枚举映射常量
 * 页面和 Drawer 共用
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 布尔值映射 ==========

/** 布尔值颜色映射 */
export function enableBoolToColor(value: boolean): string {
  return value ? 'success' : 'error';
}

/** 获取布尔值映射（启用/禁用） */
export function getEnableBoolMap(t: TFn) {
  return {
    true: { text: t('yes'), color: 'success' },
    false: { text: t('no'), color: 'error' },
  };
}
