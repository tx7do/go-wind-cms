import {useLocale} from 'next-intl';

export * from './config';
export * from './helpers/useI18n';
export * from './helpers/useApiError';

// 导出语言切换组件
export {LocaleSwitcher} from '@/components/layout/LocaleSwitcher';


// 获取当前语言代码（如 zh-CN、en-US）
export function useCurrentLocaleLanguageCode(): string {
    const locale = useLocale();
    if (locale.startsWith('zh')) return 'zh-CN';
    if (locale.startsWith('en')) return 'en-US';
    return locale;
}

// 获取当前语言代码（如 zh-CN、en-US）
export function currentLocaleLanguageCode(): string {
    if (typeof window !== 'undefined') {
        // 优先从 window.__NEXT_LOCALE__ 获取
        let locale = (window as any).__NEXT_LOCALE__;
        // 如果没有，则尝试从 URL 路径获取
        if (!locale) {
            const match = window.location.pathname.match(/^\/(zh-CN|en-US)/);
            if (match) locale = match[1];
        }
        // fallback 到浏览器语言
        if (!locale) locale = navigator.language || 'zh-CN';
        if (typeof locale === 'string') {
            if (locale.startsWith('zh')) return 'zh-CN';
            if (locale.startsWith('en')) return 'en-US';
            return locale;
        }
    }
    return 'zh-CN';
}
