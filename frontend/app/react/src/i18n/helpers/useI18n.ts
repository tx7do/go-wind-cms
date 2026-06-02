'use client';

import {useTranslations, useLocale} from 'next-intl';
import {usePathname} from 'next/navigation';

import {useI18nRouter} from "@/i18n/helpers/useI18nRouter";
import {useLoading} from '@/store/core/loading/store';
import {usePreferences} from '@/core/preferences';


/**
 * useI18n hook - 增强的多语言支持
 * locale 由 next-intl URL 路由管理，preferences zustand 管理存储的偏好
 * @param namespace 翻译模块名（如 'common', 'article', 'category'）
 */
export function useI18n(namespace: string = 'common') {
    const t = useTranslations(namespace);
    const locale = useLocale();
    const pathname = usePathname();
    const router = useI18nRouter();
    const {setLanguage} = usePreferences();
    const {start: startLoading} = useLoading();

    // 切换语言（使用 next-intl 的路由器）
    const changeLocale = (targetLocale: string) => {
        // 同步到 preferences store
        setLanguage(targetLocale as never);

        // 从当前路径中提取不包含 locale 的路径部分
        const segments = pathname.split('/').slice(2).join('/');
        const pathWithoutLocale = segments ? `/${segments}` : '/';

        const newPath = `/${targetLocale}${pathWithoutLocale}`;

        console.log('[useI18n] Start locale change, showing loading...');
        startLoading();

        router.replaceWithoutLocale(newPath);
    };

    return {t, locale, changeLocale};
}
