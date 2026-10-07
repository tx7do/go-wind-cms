/**
 * 帖子模块常量（状态枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const POST_STATUS_COLOR_MAP: Record<string, string> = {
  POST_STATUS_DRAFT: 'purple',
  POST_STATUS_PUBLISHED: 'green',
  POST_STATUS_SCHEDULED: 'gold',
  POST_STATUS_TRASHED: 'red',
  DEFAULT: 'default',
};

const POST_STATUS_VALUES = [
  'POST_STATUS_DRAFT',
  'POST_STATUS_PUBLISHED',
  'POST_STATUS_SCHEDULED',
  'POST_STATUS_TRASHED',
] as const;

export const EDITOR_TYPE_COLOR_MAP: Record<string, string> = {
  EDITOR_TYPE_MARKDOWN: 'geekblue',
  EDITOR_TYPE_RICH_TEXT: 'cyan',
  DEFAULT: 'default',
};

export function getPostStatusColor(status: string): string {
  return POST_STATUS_COLOR_MAP[status] || POST_STATUS_COLOR_MAP.DEFAULT;
}

export function getPostStatusLabel(t: TFn, status: string): string {
  // i18next 未开 returnObjects，取对象需走点路径（对齐 dict 模块约定）
  return t(`statusMap.${status}`, { defaultValue: status });
}

export function postStatusOptions(t: TFn) {
  return POST_STATUS_VALUES.map((value) => ({
    label: t(`statusMap.${value}`),
    value,
  }));
}

export function getEditorTypeColor(editorType: string): string {
  return EDITOR_TYPE_COLOR_MAP[editorType] || EDITOR_TYPE_COLOR_MAP.DEFAULT;
}

export function getEditorTypeLabel(t: TFn, editorType: string): string {
  return t(`editorTypeMap.${editorType}`, { defaultValue: editorType });
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
