import { resolveApiErrorMsg } from '@/core/transport/rest/utils'

/**
 * API 错误文案解析：不再使用后端 message 字段，
 * 用 reason 查询 i18n 错误文案（error.<REASON> 命名空间），
 * 查不到时回退到默认错误文案。
 */
export function useApiError() {
    const i18n = useNuxtApp().$i18n as any

    const resolveApiError = (error: unknown): string => {
        const locale = i18n?.locale?.value || 'zh-CN'
        const dict = i18n?.getLocaleMessage?.(locale)?.error as Record<string, string> | undefined
        return resolveApiErrorMsg(error, dict)
    }

    return { resolveApiError }
}
