export interface ApplicationConfig {
    apiURL: string
    enableMock: boolean
    appTitle: string
    appDescription: string
    defaultLocale: string
    tokenKey: string
    refreshTokenKey: string
    aesKey: string
    nodeEnv: string
    isDev: boolean
    isProd: boolean
}

/**
 * 应用配置 composable
 * 基于 Nuxt runtimeConfig（对应 React 的 env.ts）
 */
export function useAppConfig(): ApplicationConfig {
    const config = useRuntimeConfig()
    const nodeEnv = import.meta.env.MODE || 'development'

    return {
        apiURL: config.public.apiBaseUrl as string,
        enableMock: config.public.enableMock as boolean,
        appTitle: config.public.appTitle as string,
        appDescription: config.public.appDescription as string,
        defaultLocale: config.public.defaultLocale as string,
        tokenKey: config.public.tokenKey as string,
        refreshTokenKey: config.public.refreshTokenKey as string,
        aesKey: config.public.aesKey as string,
        nodeEnv,
        isDev: nodeEnv === 'development',
        isProd: nodeEnv === 'production',
    }
}

/**
 * 验证必需的环境变量
 */
export function validateEnv(): void {
    const config = useRuntimeConfig()
    if (import.meta.env.PROD && !config.public.apiBaseUrl) {
        throw new Error('生产环境必须配置 NUXT_PUBLIC_API_BASE_URL')
    }
}
