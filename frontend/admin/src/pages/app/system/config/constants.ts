/**
 * 参数模块枚举映射常量
 * 页面和 Drawer 共用
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 值类型映射 ==========

/** 参数值类型下拉选项 */
export function getValueTypeOptions(t: TFn) {
  return [
    { label: t('valueTypeString'), value: 'STRING' },
    { label: t('valueTypeBool'), value: 'BOOL' },
    { label: t('valueTypeInt'), value: 'INT' },
  ];
}

/** 参数值类型 → Tag 颜色 */
export function valueTypeToColor(valueType?: string): string {
  switch (valueType) {
    case 'BOOL':
      return 'purple';
    case 'INT':
      return 'gold';
    default:
      return 'blue';
  }
}

/** 参数值类型 → 展示文本 */
export function getValueTypeLabel(t: TFn, valueType?: string): string {
  switch (valueType) {
    case 'BOOL':
      return t('valueTypeBool');
    case 'INT':
      return t('valueTypeInt');
    default:
      return t('valueTypeString');
  }
}

// ========== 布尔值映射 ==========

/** 布尔值颜色映射 */
export function enableBoolToColor(value: boolean): string {
  return value ? 'success' : 'error';
}
