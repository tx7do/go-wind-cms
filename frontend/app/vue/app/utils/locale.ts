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
    const map: Record<string, SupportedLanguagesType> = { zh: 'zh-CN', en: 'en-US' };
    return (map[locale] || 'zh-CN') as SupportedLanguagesType;
  } catch {
    return 'zh-CN';
  }
}
