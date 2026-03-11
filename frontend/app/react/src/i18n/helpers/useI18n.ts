'use client';

import {useTranslations, useLocale} from 'next-intl';
import {useRouter} from 'next/navigation';
import {useLanguageStore} from '@/store/core/language/hooks';

/**
 * useI18n hook - 增强的多语言支持
 * @param namespace 翻译模块名（如 'common', 'article', 'category'）
 */
export function useI18n(namespace: string = 'common') {
    const t = useTranslations(namespace);
    const locale = useLocale();
    const router = useRouter();
    const languageStore = useLanguageStore();

    // 切换语言（保留当前路径）
    const changeLocale = (targetLocale: string) => {
        if (typeof window !== 'undefined') {
            const currentPath = window.location.pathname;
            
            // 移除旧的 locale 前缀
            const pathWithoutLocale = currentPath.replace(/^\/[a-zA-Z]+(-[a-zA-Z]+)?/, '');
            
            // 更新 Redux store
            languageStore.setLocale(targetLocale);
            
            // 添加新的 locale 前缀并跳转
            const newPath = `/${targetLocale}${pathWithoutLocale}`;
            router.push(newPath);
        }
    };

    return {t, locale, changeLocale};
}
