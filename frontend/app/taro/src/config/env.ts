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

export const env = {
    /**
     * API 基础地址
     */
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000',

    /**
     * 是否启用 Mock 数据
     */
    enableMock: import.meta.env.VITE_ENABLE_MOCK === 'true',

    /**
     * 应用标题
     */
    appTitle: import.meta.env.VITE_APP_TITLE || 'GoWind Content Hub',

    /**
     * 应用描述
     */
    appDescription: import.meta.env.VITE_APP_DESCRIPTION || 'A modern Content Hub built with Taro',

    /**
     * 默认语言
     */
    defaultLocale: import.meta.env.VITE_DEFAULT_LOCALE || 'zh-CN',

    /**
     * Token 存储键名
     */
    tokenKey: import.meta.env.VITE_TOKEN_KEY || 'access_token',

    /**
     * 刷新 Token 存储键名
     */
    refreshTokenKey: import.meta.env.VITE_REFRESH_TOKEN_KEY || 'refresh_token',

    /**
     * 当前环境
     */
    nodeEnv: import.meta.env.MODE || 'development',

    /**
     * 是否为开发环境
     */
    isDev: import.meta.env.MODE === 'development',

    /**
     * 是否为生产环境
     */
    isProd: import.meta.env.MODE === 'production',
} as const;

/**
 * 验证必需的环境变量
 */
export function validateEnv(): void {
    if (env.isProd && !import.meta.env.VITE_API_BASE_URL) {
        throw new Error('生产环境必须配置 VITE_API_BASE_URL');
    }
}

