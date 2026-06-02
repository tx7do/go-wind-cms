'use client';

import {useEffect, useRef} from 'react';
import {useLocale} from 'next-intl';

/**
 * 语言切换监听器 Hook - 类似 Vue3 的 useLanguageChangeEffect
 * 基于 next-intl 的 locale 变化触发回调
 * 
 * @param callback 语言变化时执行的回调函数
 * @param options 配置选项
 *   - immediate: 是否在首次渲染时立即执行（默认 false）
 *   - autoCleanup: 是否自动清理副作用（默认 true）
 */
interface UseLanguageChangeEffectOptions {
    immediate?: boolean;
    autoCleanup?: boolean;
}

export function useLanguageChangeEffect(
    callback: () => void | (() => void),
    options: UseLanguageChangeEffectOptions = {}
) {
    const {immediate = false, autoCleanup = true} = options;
    
    const locale = useLocale();
    const prevLocaleRef = useRef(locale);

    useEffect(() => {
        // 执行回调
        let cleanupFn: (() => void) | undefined;
        
        const executeCallback = async () => {
            try {
                const result = callback();
                if (typeof result === 'function') {
                    cleanupFn = result;
                }
            } catch (error) {
                console.error('[useLanguageChangeEffect] Callback error:', error);
            }
        };

        const localeChanged = prevLocaleRef.current !== locale;

        if (immediate || localeChanged) {
            void executeCallback();
            prevLocaleRef.current = locale;
        }

        // 清理函数
        return () => {
            if (autoCleanup && cleanupFn) {
                cleanupFn();
            }
        };
    }, [locale, callback, immediate, autoCleanup]);
}
