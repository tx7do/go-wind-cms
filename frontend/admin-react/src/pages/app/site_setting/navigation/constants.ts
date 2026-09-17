/**
 * 导航模块常量（位置/链接类型枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const LOCATION_COLOR_MAP: Record<string, string> = {
  HEADER: 'geekblue',
  FOOTER: 'green',
  SIDEBAR: 'gold',
  MOBILE: 'red',
  TOP_BAR: 'red',
  OFFCANVAS: 'red',
  DEFAULT: 'default',
};

const LOCATION_VALUES = [
  'HEADER',
  'FOOTER',
  'SIDEBAR',
  'MOBILE',
  'TOP_BAR',
  'OFFCANVAS',
] as const;

const LINK_TYPE_COLOR_MAP: Record<string, string> = {
  LINK_TYPE_CUSTOM: 'geekblue',
  LINK_TYPE_POST: 'green',
  LINK_TYPE_PAGE: 'gold',
  LINK_TYPE_CATEGORY: 'red',
  LINK_TYPE_EXTERNAL: 'red',
  DEFAULT: 'default',
};

const LINK_TYPE_VALUES = [
  'LINK_TYPE_CUSTOM',
  'LINK_TYPE_POST',
  'LINK_TYPE_PAGE',
  'LINK_TYPE_CATEGORY',
  'LINK_TYPE_EXTERNAL',
] as const;

export function getLocationColor(location: string): string {
  return LOCATION_COLOR_MAP[location] || LOCATION_COLOR_MAP.DEFAULT;
}

export function getLocationLabel(t: TFn, location: string): string {
  // i18next 未开 returnObjects，取对象需走点路径（对齐 dict 模块约定）
  return t(`locationMap.${location}`, { defaultValue: location });
}

export function locationOptions(t: TFn) {
  return LOCATION_VALUES.map((value) => ({
    label: t(`locationMap.${value}`),
    value,
  }));
}

export function getLinkTypeColor(linkType: string): string {
  return LINK_TYPE_COLOR_MAP[linkType] || LINK_TYPE_COLOR_MAP.DEFAULT;
}

export function getLinkTypeLabel(t: TFn, linkType: string): string {
  return t(`linkTypeMap.${linkType}`, { defaultValue: linkType });
}

export function linkTypeOptions(t: TFn) {
  return LINK_TYPE_VALUES.map((value) => ({
    label: t(`linkTypeMap.${value}`),
    value,
  }));
}

export function getBoolLabel(t: TFn, value: boolean): string {
  return t(value ? 'boolMap.true' : 'boolMap.false');
}
