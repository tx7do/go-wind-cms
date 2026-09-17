import i18n from 'i18next';
import type { SupportedLocale } from '@/locales';
import { setStoredLocale, isSupportedLocale } from './detector';

/**
 * 同步语言到浏览器环境
 * 1. 更新 HTML lang 属性
 * 2. 更新 HTML dir 属性（RTL 语言支持）
 * 3. 持久化到本地存储
 * @param locale 目标语言
 */
export const syncLocaleToBrowser = (locale: SupportedLocale): void => {
  // 1. 更新 HTML lang 属性
  document.documentElement.lang = locale;

  // 2. 支持 RTL 语言（如阿拉伯语、希伯来语）
  const rtlLocales = ['ar-SA', 'he-IL'] as const;
  document.documentElement.dir = (rtlLocales as readonly string[]).includes(locale) ? 'rtl' : 'ltr';

  // 3. 持久化存储
  setStoredLocale(locale);
};

/**
 * 监听 i18n 语言变化，自动触发同步
 * @param callback 语言变化时的额外回调
 * @returns 取消监听的函数
 */
export const setupLocaleSync = (callback?: (locale: SupportedLocale) => void): (() => void) => {
  const handleLanguageChange = (newLocale: string) => {
    if (isSupportedLocale(newLocale)) {
      // 执行核心同步
      syncLocaleToBrowser(newLocale);
      // 执行自定义回调
      callback?.(newLocale);
    }
  };

  // 监听 i18n 语言变化事件
  i18n.on('languageChanged', handleLanguageChange);

  // 立即同步一次当前语言
  const currentLocale = i18n.language;
  if (isSupportedLocale(currentLocale)) {
    syncLocaleToBrowser(currentLocale);
  }

  // 返回取消监听的函数
  return () => {
    i18n.off('languageChanged', handleLanguageChange);
  };
};

/**
 * 同步页面标题（配合 i18n 实现标题多语言）
 * @param t 翻译函数
 * @param titleKey 翻译 key
 * @param params 翻译参数
 */
export const syncPageTitle = (
  t: (key: string, params?: Record<string, any>) => string,
  titleKey: string,
  params?: Record<string, any>,
): void => {
  document.title = t(titleKey, params);
};

/**
 * 初始化语言同步（项目启动时调用一次）
 * 自动读取存储/浏览器语言，同步到 i18n 和浏览器环境
 */
export const initLocaleSync = async (): Promise<void> => {
  // 避免循环依赖，这里动态导入 detector 的 getFinalLocale
  const { getFinalLocale } = await import('./detector');
  const finalLocale = getFinalLocale();

  // 同步到 i18n
  if (i18n.language !== finalLocale) {
    await i18n.changeLanguage(finalLocale);
  } else {
    // 语言一致时，直接同步到浏览器
    syncLocaleToBrowser(finalLocale);
  }
};
