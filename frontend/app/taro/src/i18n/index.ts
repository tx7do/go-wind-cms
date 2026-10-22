import i18n from 'i18next';
import Taro from '@tarojs/taro';
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

/**
 * 从持久化的偏好（core/preferences 的 zustand persist，键 app-preferences）里
 * 读取上次使用的语言。两端存储介质不同：
 * - H5：localStorage
 * - 小程序：Taro.getStorageSync（preferences store 经 storage-adapter 写入）
 * 两侧都是同一份 JSON 结构 {state:{preferences:{app:{locale}}}}。
 */
function getPersistedLocale(): string {
  const readParsed = (raw: string | undefined): string | null => {
    if (!raw) return null;
    try {
      const parsed = JSON.parse(raw);
      const locale = parsed?.state?.preferences?.app?.locale;
      if (locale && (locale === 'zh-CN' || locale === 'en-US')) return locale;
    } catch {
      // 存储内容损坏时按无持久化处理，走语言推断
    }
    return null;
  };

  // 小程序：同步存储直读；H5 上 Taro 也代理了 storage，但直读 localStorage 更快
  if (typeof window === 'undefined') {
    try {
      const fromTaro = readParsed(Taro.getStorageSync('app-preferences') || undefined);
      if (fromTaro) return fromTaro;
    } catch {
      // 存储不可用时继续走推断
    }
  } else {
    try {
      const fromLocal = readParsed(window.localStorage.getItem('app-preferences') || undefined);
      if (fromLocal) return fromLocal;
    } catch {
      // localStorage 不可用（隐私模式等）时继续走推断
    }
  }

  // fallback: 从系统语言推断
  let navLanguage = '';
  try {
    navLanguage = (typeof window !== 'undefined' ? window.navigator?.language : Taro.getSystemInfoSync().language) || '';
  } catch {
    navLanguage = '';
  }
  if (navLanguage.startsWith('zh')) return 'zh-CN';
  if (navLanguage.startsWith('en')) return 'en-US';
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
