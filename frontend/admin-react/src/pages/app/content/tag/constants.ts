/**
 * 标签模块常量（状态枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const TAG_STATUS_COLOR_MAP: Record<string, string> = {
  TAG_STATUS_ACTIVE: 'green',
  TAG_STATUS_HIDDEN: 'orange',
  TAG_STATUS_ARCHIVED: 'gold',
  DEFAULT: 'default',
};

const TAG_STATUS_VALUES = [
  'TAG_STATUS_ACTIVE',
  'TAG_STATUS_HIDDEN',
  'TAG_STATUS_ARCHIVED',
] as const;

export function getTagStatusColor(status: string): string {
  return TAG_STATUS_COLOR_MAP[status] || TAG_STATUS_COLOR_MAP.DEFAULT;
}

export function getTagStatusLabel(t: TFn, status: string): string {
  // i18next 未开 returnObjects，取对象需走点路径（对齐 dict 模块约定）
  return t(`statusMap.${status}`, { defaultValue: status });
}

export function tagStatusOptions(t: TFn) {
  return TAG_STATUS_VALUES.map((value) => ({
    label: t(`statusMap.${value}`),
    value,
  }));
}

export function getBoolLabel(t: TFn, value: boolean): string {
  return t(value ? 'boolMap.true' : 'boolMap.false');
}

/** 图标缺前缀时补 carbon:（对齐 vben 侧约定） */
export function getIconName(icon: string): string {
  if (!icon) return '';
  if (icon.includes(':')) return icon;
  return `carbon:${icon}`;
}

/** 从 translations 里取当前语言的翻译，缺失时回退第一条（对齐 vben 侧） */
export function pickTranslation(
  row: any,
  lang: string,
): { name?: string; description?: string; slug?: string } | undefined {
  return (
    row?.translations?.find((t: any) => t.languageCode === lang) ||
    row?.translations?.[0] ||
    undefined
  );
}
