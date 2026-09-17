/**
 * 站点配置模块常量（类型枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const SETTING_TYPE_COLOR_MAP: Record<string, string> = {
  SETTING_TYPE_TEXT: 'blue',
  SETTING_TYPE_TEXTAREA: 'geekblue',
  SETTING_TYPE_NUMBER: 'green',
  SETTING_TYPE_BOOLEAN: 'gold',
  SETTING_TYPE_URL: 'cyan',
  SETTING_TYPE_EMAIL: 'purple',
  SETTING_TYPE_IMAGE: 'orange',
  SETTING_TYPE_SELECT: 'magenta',
  SETTING_TYPE_JSON: 'blue',
  DEFAULT: 'default',
};

const SETTING_TYPE_VALUES = [
  'SETTING_TYPE_TEXT',
  'SETTING_TYPE_TEXTAREA',
  'SETTING_TYPE_NUMBER',
  'SETTING_TYPE_BOOLEAN',
  'SETTING_TYPE_URL',
  'SETTING_TYPE_EMAIL',
  'SETTING_TYPE_IMAGE',
  'SETTING_TYPE_SELECT',
  'SETTING_TYPE_JSON',
] as const;

export function getTypeColor(type: string): string {
  return SETTING_TYPE_COLOR_MAP[type] || SETTING_TYPE_COLOR_MAP.DEFAULT;
}

export function getTypeLabel(t: TFn, type: string): string {
  // i18next 未开 returnObjects，取对象需走点路径（对齐 dict 模块约定）
  return t(`typeMap.${type}`, { defaultValue: type });
}

export function typeOptions(t: TFn) {
  return SETTING_TYPE_VALUES.map((value) => ({
    label: t(`typeMap.${value}`),
    value,
  }));
}
