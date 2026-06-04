import i18n from 'i18next';
import {initReactI18next} from 'react-i18next';
import {allMessages, defaultLocale, DEFAULT_LANGUAGE, locales, validateLocale, getFlattenedMessages} from './config';

// 获取扁平化后的消息
const flattenedMessages = getFlattenedMessages();

// 初始化 i18next
i18n
  .use(initReactI18next)
  .init({
    resources: {
      'zh-CN': {
        translation: flattenedMessages['zh-CN'],
      },
      'en-US': {
        translation: flattenedMessages['en-US'],
      },
    },
    lng: getPersistedLocale(),
    fallbackLng: 'zh-CN',
    interpolation: {
      escapeValue: false,
    },
    // 支持嵌套的命名空间（如 page.home）
    nsSeparator: '.',
    keySeparator: '.',
  });

// 从 zustand persist 存储中读取持久化的 locale
function getPersistedLocale(): string {
  if (typeof window === 'undefined') return defaultLocale;
  try {
    const raw = localStorage.getItem('app-preferences');
    if (raw) {
      const parsed = JSON.parse(raw);
      const locale = parsed?.state?.preferences?.app?.locale;
      if (locale && (locale === 'zh-CN' || locale === 'en-US')) return locale;
    }
  } catch {}
  // fallback: 从 navigator.language 推断
  const nav = navigator.language || '';
  if (nav.startsWith('zh')) return 'zh-CN';
  if (nav.startsWith('en')) return 'en-US';
  return defaultLocale;
}

// 显式导出 config 中的内容
export {allMessages, defaultLocale, DEFAULT_LANGUAGE, locales, validateLocale};
export type {Locale} from './config';

// 导出 useI18n hook
export {useI18n} from './helpers/useI18n';


// 获取当前语言代码（如 zh-CN、en-US）
export function useCurrentLocaleLanguageCode(): string {
  return i18n.language || 'zh-CN';
}

// 获取当前语言代码（如 zh-CN、en-US）
export function currentLocaleLanguageCode(): string {
  return i18n.language || 'zh-CN';
}
