/**
 * 内容分发模块常量（平台/状态枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const PLATFORM_COLOR_MAP: Record<string, string> = {
  PLATFORM_BLUESKY: 'geekblue',
  PLATFORM_MASTODON: 'purple',
  PLATFORM_DEVTO: 'volcano',
  DEFAULT: 'default',
};

const PLATFORM_VALUES = [
  'PLATFORM_BLUESKY',
  'PLATFORM_MASTODON',
  'PLATFORM_DEVTO',
] as const;

const STATUS_COLOR_MAP: Record<string, string> = {
  STATUS_DRAFT: 'default',
  STATUS_SCHEDULED: 'geekblue',
  STATUS_PUBLISHING: 'processing',
  STATUS_PUBLISHED: 'green',
  STATUS_FAILED: 'red',
  DEFAULT: 'default',
};

const STATUS_VALUES = [
  'STATUS_DRAFT',
  'STATUS_SCHEDULED',
  'STATUS_PUBLISHING',
  'STATUS_PUBLISHED',
  'STATUS_FAILED',
] as const;

/** 可立即接入的平台（其余枚举值预留，不在创建入口展示） */
export const SUPPORTED_PLATFORM_VALUES = [
  'PLATFORM_BLUESKY',
  'PLATFORM_MASTODON',
  'PLATFORM_DEVTO',
  'PLATFORM_TELEGRAM',
  'PLATFORM_WORDPRESS',
  'PLATFORM_WECHAT_OA',
  'PLATFORM_NOSTR',
  'PLATFORM_WEIBO',
  'PLATFORM_LINKEDIN',
] as const;

/** Bluesky 单帖文案上限（按 rune 预检口径） */
export const MAX_POST_CHARS = 300;

export function getPlatformColor(platform: string): string {
  return PLATFORM_COLOR_MAP[platform] || PLATFORM_COLOR_MAP.DEFAULT;
}

export function getPlatformLabel(t: TFn, platform: string): string {
  return t(`platformMap.${platform}`, { defaultValue: platform });
}

export function platformOptions(t: TFn) {
  return PLATFORM_VALUES.map((value) => ({
    label: t(`platformMap.${value}`),
    value,
  }));
}

export function supportedPlatformOptions(t: TFn) {
  return SUPPORTED_PLATFORM_VALUES.map((value) => ({
    label: t(`platformMap.${value}`),
    value,
  }));
}

export function getStatusColor(status: string): string {
  return STATUS_COLOR_MAP[status] || STATUS_COLOR_MAP.DEFAULT;
}

export function getStatusLabel(t: TFn, status: string): string {
  return t(`statusMap.${status}`, { defaultValue: status });
}

export function statusOptions(t: TFn) {
  return STATUS_VALUES.map((value) => ({
    label: t(`statusMap.${value}`),
    value,
  }));
}
