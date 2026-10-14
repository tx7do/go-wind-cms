/**
 * 社交账号模块常量（平台枚举展示 + 按平台动态凭据字段配置）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const PLATFORM_COLOR_MAP: Record<string, string> = {
  PLATFORM_BLUESKY: 'geekblue',
  PLATFORM_MASTODON: 'purple',
  PLATFORM_DEVTO: 'volcano',
  PLATFORM_TELEGRAM: 'blue',
  PLATFORM_WORDPRESS: 'cyan',
  PLATFORM_WECHAT_OA: 'green',
  PLATFORM_NOSTR: 'orange',
  PLATFORM_WEIBO: 'red',
  PLATFORM_LINKEDIN: 'geekblue',
  PLATFORM_METAWEBLOG: 'gold',
  DEFAULT: 'default',
};

const PLATFORM_VALUES = [
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

/** 平台凭据字段定义：drawer 按所选平台动态渲染表单 */
export interface CredentialField {
  /** 凭据 JSON 键 */
  key: string;
  /** 是否必填（仅新建时强校验；编辑留空=不修改） */
  required: boolean;
  /** 密码类字段（Password 输入框） */
  secret?: boolean;
  /** 该键的值同时回填到账号 handle 展示列 */
  alsoHandle?: boolean;
}

export const PLATFORM_CREDENTIAL_FIELDS: Record<string, CredentialField[]> = {
  PLATFORM_BLUESKY: [{ key: 'app_password', required: true, secret: true }],
  PLATFORM_MASTODON: [
    { key: 'server', required: true, alsoHandle: true },
    { key: 'access_token', required: true, secret: true },
  ],
  PLATFORM_DEVTO: [{ key: 'api_key', required: true, secret: true }],
  PLATFORM_TELEGRAM: [
    { key: 'bot_token', required: true, secret: true },
    { key: 'chat_id', required: true, alsoHandle: true },
  ],
  PLATFORM_WORDPRESS: [
    { key: 'server', required: true, alsoHandle: true },
    { key: 'username', required: true },
    { key: 'app_password', required: true, secret: true },
  ],
  PLATFORM_WECHAT_OA: [
    { key: 'app_id', required: true, alsoHandle: true },
    { key: 'app_secret', required: true, secret: true },
    { key: 'thumb_media_id', required: true },
  ],
  PLATFORM_NOSTR: [
    { key: 'private_key', required: true, secret: true },
    { key: 'relays', required: true },
  ],
  PLATFORM_WEIBO: [{ key: 'access_token', required: true, secret: true }],
  PLATFORM_LINKEDIN: [{ key: 'access_token', required: true, secret: true }],
  PLATFORM_METAWEBLOG: [
    { key: 'rpc_url', required: true, alsoHandle: true },
    { key: 'username', required: true },
    { key: 'password', required: true, secret: true },
  ],
};

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
