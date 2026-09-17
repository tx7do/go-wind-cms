import type { SupportedLocale } from '@/locales';

// 配置常量（和你的 i18n 配置保持一致）
const STORAGE_KEY = 'app_locale';
const DEFAULT_LOCALE: SupportedLocale = 'zh-CN';
const SUPPORTED_LOCALES: readonly SupportedLocale[] = ['zh-CN', 'en-US'] as const;

/**
 * 从浏览器环境自动检测首选语言
 * @returns 检测到的语言代码
 */
export const detectBrowserLocale = (): SupportedLocale => {
  // 1. 优先读取浏览器首选语言列表
  const browserLanguages = navigator.languages ?? [navigator.language];

  for (const lang of browserLanguages) {
    // 完全匹配（zh-CN → zh-CN）
    const exactMatch = SUPPORTED_LOCALES.find((supported) => supported === lang);
    if (exactMatch) return exactMatch;

    // 前缀匹配（zh → zh-CN，en → en-US）
    const prefixMatch = SUPPORTED_LOCALES.find((supported) =>
      supported.startsWith(lang.split('-')[0]),
    );
    if (prefixMatch) return prefixMatch;
  }

  // 2. 无匹配时返回默认语言
  return DEFAULT_LOCALE;
};

/**
 * 从本地存储读取用户上次选择的语言
 * @returns 存储的语言代码 | null
 */
export const getStoredLocale = (): SupportedLocale | null => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored && SUPPORTED_LOCALES.includes(stored as SupportedLocale)) {
      return stored as SupportedLocale;
    }
    return null;
  } catch (e) {
    // 隐私模式/存储禁用时静默失败
    console.warn('[i18n detector] 读取本地存储失败', e);
    return null;
  }
};

/**
 * 存储用户选择的语言到本地存储
 * @param locale 要存储的语言代码
 */
export const setStoredLocale = (locale: SupportedLocale): void => {
  try {
    localStorage.setItem(STORAGE_KEY, locale);
  } catch (e) {
    console.warn('[i18n detector] 写入本地存储失败', e);
  }
};

/**
 * 综合检测：优先级 存储语言 > 浏览器语言 > 默认语言
 * @returns 最终生效的语言代码
 */
export const getFinalLocale = (): SupportedLocale => {
  return getStoredLocale() ?? detectBrowserLocale();
};

/**
 * 校验语言是否为项目支持的语言
 * @param locale 待校验的语言代码
 */
export const isSupportedLocale = (locale: string): locale is SupportedLocale => {
  return SUPPORTED_LOCALES.includes(locale as SupportedLocale);
};
