/**
 * 站点模块常量（状态枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const STATUS_COLOR_MAP: Record<string, string> = {
  SITE_STATUS_ACTIVE: 'blue',
  SITE_STATUS_INACTIVE: 'geekblue',
  SITE_STATUS_MAINTENANCE: 'green',
  DEFAULT: 'default',
};

const STATUS_VALUES = [
  'SITE_STATUS_ACTIVE',
  'SITE_STATUS_INACTIVE',
  'SITE_STATUS_MAINTENANCE',
] as const;

export function getStatusColor(status: string): string {
  return STATUS_COLOR_MAP[status] || STATUS_COLOR_MAP.DEFAULT;
}

export function getStatusLabel(t: TFn, status: string): string {
  // i18next 未开 returnObjects，取对象需走点路径（对齐 dict 模块约定）
  return t(`statusMap.${status}`, { defaultValue: status });
}

export function statusOptions(t: TFn) {
  return STATUS_VALUES.map((value) => ({
    label: t(`statusMap.${value}`),
    value,
  }));
}

export function getBoolLabel(t: TFn, value: boolean): string {
  return t(value ? 'boolMap.true' : 'boolMap.false');
}
