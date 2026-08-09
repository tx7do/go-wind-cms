import type { SupportedLanguagesType } from '@/core/preferences/types';

/**
 * 直接从 Nuxt i18n 运行时获取当前 locale（实时值）
 * 适用于 API composables、service 层等非 setup 上下文
 */
export function getCurrentLocale(): SupportedLanguagesType {
  if (import.meta.server) return 'zh-CN';
  try {
    const nuxtApp = useNuxtApp();
    const locale = (nuxtApp.$i18n as any)?.locale?.value as string;
    if (locale?.includes('-')) return locale as SupportedLanguagesType;
    // 注意：'ar' 为无地区后缀的宏语言代码，需在此显式映射，否则会落入
    // 下方 fallback 被错误地拼成 'ar-CN'。
    const map: Record<string, SupportedLanguagesType> = { zh: 'zh-CN', en: 'en-US', ar: 'ar' };
    return (map[locale] || 'zh-CN') as SupportedLanguagesType;
  } catch {
    return 'zh-CN';
  }
}
