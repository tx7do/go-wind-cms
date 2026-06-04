/**
 * next-intl 兼容层 - 将 next-intl API 映射到 react-i18next
 * 
 * 用法：将 `import {useTranslations} from '@/lib/next-intl-compat'` 
 *      替换为 `import {useTranslations} from '@/lib/next-intl-compat'`
 */
import {useTranslation} from 'react-i18next';

/**
 * 兼容 next-intl 的 useTranslations
 * next-intl: useTranslations('namespace') → t('key')
 * react-i18next: useTranslation() → t('namespace.key')
 */
export function useTranslations(namespace: string) {
    const {t} = useTranslation();
    
    return (key: string, params?: Record<string, any>) => {
        // next-intl 使用 'key' 格式，react-i18next 需要 'namespace.key'
        const fullKey = `${namespace}.${key}`;
        return t(fullKey, params) as string;
    };
}

/**
 * 兼容 next-intl 的 useLocale
 */
export function useLocale(): string {
    const {i18n} = useTranslation();
    return i18n.language || 'zh-CN';
}
