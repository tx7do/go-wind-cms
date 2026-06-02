import { watch, ref, onUnmounted } from 'vue'

/**
 * 语言切换监听器 Hook
 * 基于 @nuxtjs/i18n 的 locale 变化触发回调
 *
 * @param callback 语言变化时执行的回调函数
 * @param options 配置选项
 *   - immediate: 是否在首次渲染时立即执行（默认 false）
 *   - autoCleanup: 是否自动清理副作用（默认 true）
 */
interface UseLanguageChangeEffectOptions {
    immediate?: boolean
    autoCleanup?: boolean
}

export function useLanguageChangeEffect(
    callback: () => void | (() => void),
    options: UseLanguageChangeEffectOptions = {},
) {
    const { immediate = false, autoCleanup = true } = options

    const { locale } = useI18n()
    const prevLocale = ref(locale.value)
    let cleanupFn: (() => void) | undefined

    const stopWatch = watch(
        () => locale.value,
        (newLocale) => {
            const localeChanged = prevLocale.value !== newLocale

            if (immediate || localeChanged) {
                // 清理上一次的副作用
                if (autoCleanup && cleanupFn) {
                    cleanupFn()
                }

                try {
                    const result = callback()
                    if (typeof result === 'function') {
                        cleanupFn = result
                    }
                } catch (error) {
                    console.error('[useLanguageChangeEffect] Callback error:', error)
                }

                prevLocale.value = newLocale
            }
        },
        { immediate },
    )

    onUnmounted(() => {
        stopWatch()
        if (autoCleanup && cleanupFn) {
            cleanupFn()
        }
    })
}
