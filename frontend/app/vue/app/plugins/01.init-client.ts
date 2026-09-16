import {RequestClient} from '@/core/transport/rest/request-client'
import {resolveApiErrorMsg} from '@/core/transport/rest/utils'
import {useAccessStore} from '@/stores/modules/core/access.state'
import {useUserStore} from '@/stores/modules/core/user.state'
import {refreshToken as apiRefreshToken} from '@/api/composables/auth'
import {fetchMe} from '@/api/composables/user-profile'
import {useAppConfig} from '@/hooks/use-app-config'

export default defineNuxtPlugin(async (nuxtApp) => {
    // 服务端也必须初始化：useAsyncData 的 SSR 数据获取走同一 RequestClient，
    // 若 server 端不 init，getInstance() 直接抛错，导致搜索页等 SSR 首屏失败。
    // callbacks 中的 store/i18n 访问在 SSR 下均可用；onReAuthenticate 内部已有
    // window 守卫，服务端不会触发导航。
    const config = useAppConfig()
    const i18n = nuxtApp.$i18n

    const accessStore = useAccessStore()
    const userStore = useUserStore()

    // Initialize RequestClient
    RequestClient.init(config.apiURL, {
        getToken: () => {
            const token = accessStore.accessToken
            return token?.value ?? null
        },

        getLocale: () => {
            return (i18n as any)?.locale?.value || 'zh-CN'
        },

        refreshToken: async () => {
            const refreshTokenValue = accessStore.refreshToken?.value ?? ''
            if (!refreshTokenValue) return ''

            try {
                const resp = await apiRefreshToken(refreshTokenValue)
                const newToken = resp.access_token || ''

                accessStore.setAccessToken(newToken, resp.expires_in && resp.expires_in !== 0
                    ? Date.now() + resp.expires_in * 1000
                    : undefined)

                if (resp.refresh_token) {
                    accessStore.setRefreshToken(resp.refresh_token, resp.refresh_expires_in && resp.refresh_expires_in !== 0
                        ? Date.now() + resp.refresh_expires_in * 1000
                        : undefined)
                }

                return newToken
            } catch {
                accessStore.setLoginExpired(true)
                return ''
            }
        },

        onReAuthenticate: async (redirect?: boolean) => {
            accessStore.setAccessToken('')
            accessStore.setRefreshToken('')
            accessStore.setLoginExpired(true)
            userStore.setUserInfo(null as any)

            if (redirect && typeof window !== 'undefined') {
                const currentPath = window.location.pathname
                window.location.href = `/${(i18n as any)?.locale?.value || 'zh-CN'}/login?redirect=${encodeURIComponent(currentPath)}`
            }
        },

        onError: (message: string) => {
            if (config.isDev) {
                console.error('[RequestClient Error]', message)
            }
        },

        // 错误文案统一走 reason → i18n（error.<REASON>），查不到回退默认文案
        getErrorMsg: (error: unknown) => {
            const locale = (i18n as any)?.locale?.value || 'zh-CN'
            const dict = (i18n as any)?.getLocaleMessage?.(locale)?.error as Record<string, string> | undefined
            return resolveApiErrorMsg(error, dict)
        },
    })

    // Preload user info if token exists
    if (accessStore.accessToken?.value && !accessStore.loginExpired) {
        const expiresAt = accessStore.accessToken.expiresAt
        const isValid = !expiresAt || expiresAt > Date.now()
        if (isValid) {
            try {
                const user = await fetchMe()
                userStore.setUserInfo(user as any)
            } catch {
                // Silent - token may be expired
            }
        }
    }
})
