/**
 * 页面模块常量（状态/类型/编辑器/区块类型枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const PAGE_STATUS_COLOR_MAP: Record<string, string> = {
  PAGE_STATUS_DRAFT: 'purple',
  PAGE_STATUS_PUBLISHED: 'green',
  PAGE_STATUS_ARCHIVED: 'gold',
  DEFAULT: 'default',
};

const PAGE_STATUS_VALUES = [
  'PAGE_STATUS_DRAFT',
  'PAGE_STATUS_PUBLISHED',
  'PAGE_STATUS_ARCHIVED',
] as const;

const PAGE_TYPE_COLOR_MAP: Record<string, string> = {
  PAGE_TYPE_DEFAULT: 'purple',
  PAGE_TYPE_HOME: 'blue',
  PAGE_TYPE_ERROR_404: 'orange',
  PAGE_TYPE_ERROR_500: 'red',
  PAGE_TYPE_CUSTOM: 'lime',
  DEFAULT: 'default',
};

const PAGE_TYPE_VALUES = [
  'PAGE_TYPE_DEFAULT',
  'PAGE_TYPE_HOME',
  'PAGE_TYPE_ERROR_404',
  'PAGE_TYPE_ERROR_500',
  'PAGE_TYPE_CUSTOM',
] as const;

export const EDITOR_TYPE_COLOR_MAP: Record<string, string> = {
  EDITOR_TYPE_MARKDOWN: 'geekblue',
  EDITOR_TYPE_RICH_TEXT: 'cyan',
  EDITOR_TYPE_JSON_BLOCK: 'magenta',
  EDITOR_TYPE_PLAIN_TEXT: 'default',
  EDITOR_TYPE_CODE: 'blue',
  EDITOR_TYPE_VISUAL_BUILDER: 'lime',
  DEFAULT: 'default',
};

export const EDITOR_TYPE_VALUES = [
  'EDITOR_TYPE_MARKDOWN',
  'EDITOR_TYPE_RICH_TEXT',
  'EDITOR_TYPE_JSON_BLOCK',
  'EDITOR_TYPE_CODE',
  'EDITOR_TYPE_PLAIN_TEXT',
] as const;

export const SECTION_TYPE_VALUES = [
  'SECTION_TYPE_RICH_TEXT',
  'SECTION_TYPE_MARKDOWN',
  'SECTION_TYPE_TITLE',
  'SECTION_TYPE_IMAGE',
  'SECTION_TYPE_GALLERY',
  'SECTION_TYPE_VIDEO',
  'SECTION_TYPE_BUTTON',
  'SECTION_TYPE_DIVIDER',
  'SECTION_TYPE_SPACER',
  'SECTION_TYPE_CODE',
  'SECTION_TYPE_HTML',
  'SECTION_TYPE_FORM',
  'SECTION_TYPE_CAROUSEL',
  'SECTION_TYPE_CUSTOM',
] as const;

export function getPageStatusColor(status: string): string {
  return PAGE_STATUS_COLOR_MAP[status] || PAGE_STATUS_COLOR_MAP.DEFAULT;
}

export function getPageStatusLabel(t: TFn, status: string): string {
  // i18next 未开 returnObjects，取对象需走点路径（对齐 dict 模块约定）
  return t(`statusMap.${status}`, { defaultValue: status });
}

export function pageStatusOptions(t: TFn) {
  return PAGE_STATUS_VALUES.map((value) => ({
    label: t(`statusMap.${value}`),
    value,
  }));
}

export function getPageTypeColor(type: string): string {
  return PAGE_TYPE_COLOR_MAP[type] || PAGE_TYPE_COLOR_MAP.DEFAULT;
}

export function getPageTypeLabel(t: TFn, type: string): string {
  return t(`typeMap.${type}`, { defaultValue: type });
}

export function pageTypeOptions(t: TFn) {
  return PAGE_TYPE_VALUES.map((value) => ({
    label: t(`typeMap.${value}`),
    value,
  }));
}

export function getEditorTypeColor(editorType: string): string {
  return EDITOR_TYPE_COLOR_MAP[editorType] || EDITOR_TYPE_COLOR_MAP.DEFAULT;
}

export function getEditorTypeLabel(t: TFn, editorType: string): string {
  return t(`editorTypeMap.${editorType}`, { defaultValue: editorType });
}

export function editorTypeOptions(t: TFn) {
  return EDITOR_TYPE_VALUES.map((value) => ({
    label: t(`editorTypeMap.${value}`),
    value,
  }));
}

export function getSectionTypeLabel(t: TFn, type: string): string {
  return t(`sectionTypeMap.${type}`, { defaultValue: type });
}

export function sectionTypeOptions(t: TFn) {
  return SECTION_TYPE_VALUES.map((value) => ({
    label: t(`sectionTypeMap.${value}`),
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
