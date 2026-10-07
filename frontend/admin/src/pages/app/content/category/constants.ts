/**
 * 分类模块常量（状态枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const CATEGORY_STATUS_COLOR_MAP: Record<string, string> = {
  CATEGORY_STATUS_ACTIVE: 'green',
  CATEGORY_STATUS_HIDDEN: 'orange',
  CATEGORY_STATUS_ARCHIVED: 'gold',
  DEFAULT: 'default',
};

const CATEGORY_STATUS_VALUES = [
  'CATEGORY_STATUS_ACTIVE',
  'CATEGORY_STATUS_HIDDEN',
  'CATEGORY_STATUS_ARCHIVED',
] as const;

export function getCategoryStatusColor(status: string): string {
  return CATEGORY_STATUS_COLOR_MAP[status] || CATEGORY_STATUS_COLOR_MAP.DEFAULT;
}

export function getCategoryStatusLabel(t: TFn, status: string): string {
  // i18next 未开 returnObjects，取对象需走点路径（对齐 dict 模块约定）
  return t(`statusMap.${status}`, { defaultValue: status });
}

export function categoryStatusOptions(t: TFn) {
  return CATEGORY_STATUS_VALUES.map((value) => ({
    label: t(`statusMap.${value}`),
    value,
  }));
}

export function getBoolLabel(t: TFn, value: boolean): string {
  return t(value ? 'boolMap.true' : 'boolMap.false');
}

/** 从 translations 里取当前语言的翻译，缺失时回退第一条（对齐 vben 侧） */
export function pickTranslation(row: any, lang: string) {
  return (
    row?.translations?.find((tr: any) => tr.languageCode === lang) ||
    row?.translations?.[0] ||
    undefined
  );
}
