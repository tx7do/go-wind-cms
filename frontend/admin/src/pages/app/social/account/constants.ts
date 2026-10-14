/**
 * 社交账号模块常量（平台枚举展示）
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
