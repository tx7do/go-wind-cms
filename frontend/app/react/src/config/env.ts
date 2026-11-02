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
     * API 基础地址：缺省同源相对路径（浏览器请求 /app/v1/... 由部署层 nginx
     * 转发到 app 服务，不经独立 API 域名）。开发态由 .env.development 指向
     * http://localhost:6700 直连。静态导出后全站客户端取数，无服务端直连路径。
     */
    apiBaseUrl: process.env.NEXT_PUBLIC_API_BASE_URL || '',

    /**
     * 是否启用 Mock 数据
     */
    enableMock: process.env.NEXT_PUBLIC_ENABLE_MOCK === 'true',

    /**
     * 应用标题
     */
    appTitle: process.env.NEXT_PUBLIC_APP_TITLE || 'GoWind Content Hub',

    /**
     * 应用描述
     */
    appDescription: process.env.NEXT_PUBLIC_APP_DESCRIPTION || 'A modern Content Hub built with Next.js',

    /**
     * 默认语言
     */
    defaultLocale: process.env.NEXT_PUBLIC_DEFAULT_LOCALE || 'zh-CN',

    /**
     * Token 存储键名
     */
    tokenKey: process.env.NEXT_PUBLIC_TOKEN_KEY || 'access_token',

    /**
     * 刷新 Token 存储键名
     */
    refreshTokenKey: process.env.NEXT_PUBLIC_REFRESH_TOKEN_KEY || 'refresh_token',

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
    if (env.isProd && !process.env.NEXT_PUBLIC_API_BASE_URL) {
        throw new Error('生产环境必须配置 NEXT_PUBLIC_API_BASE_URL');
    }
}
