/**
 * 应用核心常量
 */

/**
 * 登录页面 url 地址
 */
export const LOGIN_PATH = '/login'

/**
 * 默认首页地址
 */
export const DEFAULT_HOME_PATH = '/'

/**
 * Token 存储命名空间
 */
export const APP_NAMESPACE = 'gowind'

/**
 * Access Token 存储 key
 */
export const ACCESS_TOKEN_KEY = 'access-token'

/**
 * Refresh Token 存储 key
 */
export const REFRESH_TOKEN_KEY = 'refresh-token'

/**
 * Preferences 存储 key
 */
export const PREFERENCES_STORAGE_KEY = 'app-preferences'

/**
 * 默认时区
 */
export const DEFAULT_TIME_ZONE = 'Asia/Shanghai'

/**
 * 支持的语言列表
 */
export const SUPPORTED_LOCALES = ['zh-CN', 'en-US'] as const

/**
 * 默认语言
 */
export const DEFAULT_LANGUAGE = 'zh-CN'

/**
 * 支持的语言类型
 */
export type SupportedLocale = (typeof SUPPORTED_LOCALES)[number]
