/**
 * 评论模块常量（状态/内容类型/作者类型枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const STATUS_COLOR_MAP: Record<string, string> = {
  STATUS_PENDING: 'geekblue',
  STATUS_APPROVED: 'green',
  STATUS_REJECTED: 'orange',
  STATUS_SPAM: 'red',
  DEFAULT: 'default',
};

const STATUS_VALUES = [
  'STATUS_PENDING',
  'STATUS_APPROVED',
  'STATUS_REJECTED',
  'STATUS_SPAM',
] as const;

const CONTENT_TYPE_COLOR_MAP: Record<string, string> = {
  CONTENT_TYPE_POST: 'geekblue',
  CONTENT_TYPE_PAGE: 'purple',
  CONTENT_TYPE_PRODUCT: 'volcano',
  DEFAULT: 'geekblue',
};

const CONTENT_TYPE_VALUES = [
  'CONTENT_TYPE_POST',
  'CONTENT_TYPE_PAGE',
  'CONTENT_TYPE_PRODUCT',
] as const;

const AUTHOR_TYPE_COLOR_MAP: Record<string, string> = {
  AUTHOR_TYPE_GUEST: 'default',
  AUTHOR_TYPE_USER: 'green',
  AUTHOR_TYPE_ADMIN: 'blue',
  AUTHOR_TYPE_MODERATOR: 'orange',
  DEFAULT: 'default',
};

const AUTHOR_TYPE_VALUES = [
  'AUTHOR_TYPE_GUEST',
  'AUTHOR_TYPE_USER',
  'AUTHOR_TYPE_ADMIN',
  'AUTHOR_TYPE_MODERATOR',
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

export function getContentTypeColor(contentType: string): string {
  return CONTENT_TYPE_COLOR_MAP[contentType] || CONTENT_TYPE_COLOR_MAP.DEFAULT;
}

export function getContentTypeLabel(t: TFn, contentType: string): string {
  return t(`contentTypeMap.${contentType}`, { defaultValue: contentType });
}

export function contentTypeOptions(t: TFn) {
  return CONTENT_TYPE_VALUES.map((value) => ({
    label: t(`contentTypeMap.${value}`),
    value,
  }));
}

export function getAuthorTypeColor(authorType: string): string {
  return AUTHOR_TYPE_COLOR_MAP[authorType] || AUTHOR_TYPE_COLOR_MAP.DEFAULT;
}

export function getAuthorTypeLabel(t: TFn, authorType: string): string {
  return t(`authorTypeMap.${authorType}`, { defaultValue: authorType });
}

export function authorTypeOptions(t: TFn) {
  return AUTHOR_TYPE_VALUES.map((value) => ({
    label: t(`authorTypeMap.${value}`),
    value,
  }));
}
