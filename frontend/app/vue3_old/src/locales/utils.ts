/**
 * Locale format conversion utilities
 * Converts between i18n locale format (zhCN, enUS) and API languageCode format (zh-CN, en-US)
 */

import type { SupportedLanguagesType } from '@/locales/typing'

/**
 * Convert i18n locale to API languageCode
 * @param locale - i18n locale (zhCN, enUS)
 * @returns API languageCode (zh-CN, en-US)
 */
export function localeToLanguageCode(locale: SupportedLanguagesType): string {
  const localeMap: Record<SupportedLanguagesType, string> = {
    zhCN: 'zh-CN',
    enUS: 'en-US',
  }
  return localeMap[locale] || 'zh-CN'
}

/**
 * Convert API languageCode to i18n locale
 * @param languageCode - API languageCode (zh-CN, en-US)
 * @returns i18n locale (zhCN, enUS)
 */
export function languageCodeToLocale(languageCode: string): SupportedLanguagesType {
  const codeMap: Record<string, SupportedLanguagesType> = {
    'zh-CN': 'zhCN',
    'en-US': 'enUS',
  }
  return codeMap[languageCode] || 'zhCN'
}

