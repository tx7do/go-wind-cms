/**
 * 环境变量配置工具类
 *
 * @example
 * ```typescript
 * // 在客户端或服务端使用
 * import { env } from '@/config/env';
 *
 * console.log(env.apiBaseUrl);  // API 地址
 * console.log(env.appTitle);    // 应用标题
 * ```
 */

// 声明全局常量类型
declare const API_BASE_URL: string;
declare const APP_TITLE: string;
declare const ENABLE_MOCK: string;
declare const DEFAULT_LOCALE: string;
declare const TOKEN_KEY: string;
declare const REFRESH_TOKEN_KEY: string;
declare const APP_NAMESPACE: string;
declare const AES_KEY: string;

export const env = {
    /**
     * API 基础地址
     */
    apiBaseUrl: typeof API_BASE_URL !== 'undefined' ? API_BASE_URL : '',

    /**
     * 是否启用 Mock 数据
     */
    enableMock: typeof ENABLE_MOCK !== 'undefined' && ENABLE_MOCK === 'true',

    /**
     * 应用标题
     */
    appTitle: typeof APP_TITLE !== 'undefined' ? APP_TITLE : 'GoWind Content Hub',

    /**
     * 应用描述
     */
    appDescription: typeof APP_TITLE !== 'undefined' ? APP_TITLE : 'A modern Content Hub built with Taro',

    /**
     * 默认语言
     */
    defaultLocale: typeof DEFAULT_LOCALE !== 'undefined' ? DEFAULT_LOCALE : 'zh-CN',

    /**
     * Token 存储键名
     */
    tokenKey: typeof TOKEN_KEY !== 'undefined' ? TOKEN_KEY : 'access_token',

    /**
     * 刷新 Token 存储键名
     */
    refreshTokenKey: typeof REFRESH_TOKEN_KEY !== 'undefined' ? REFRESH_TOKEN_KEY : 'refresh_token',

    /**
     * 当前环境
     */
    nodeEnv: process.env.NODE_ENV || 'development',

    /**
     * 是否为开发环境
     */
    isDev: process.env.NODE_ENV === 'development',

    /**
     * 是否为生产环境
     */
    isProd: process.env.NODE_ENV === 'production',
} as const;

/**
 * 验证必需的环境变量
 */
export function validateEnv(): void {
    if (env.isProd && !env.apiBaseUrl) {
        throw new Error('生产环境必须配置 VITE_API_BASE_URL');
    }
}

