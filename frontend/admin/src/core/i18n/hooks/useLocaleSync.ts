import { useCallback, useEffect, useRef } from 'react';

import { useI18n } from './useI18n';
import { usePreferencesStore } from '@/core/preferences';
import type { SupportedLocale } from '@/locales';

export const useLocaleSync = () => {
  const { i18n } = useI18n();
  // 使用稳定的选择器，避免每次渲染创建新对象
  const locale = usePreferencesStore((state) => state.preferences.app.locale);
  const setPreferences = usePreferencesStore((state) => state.setPreferences);

  // 标记是否已初始化
  const initialized = useRef(false);

  // preferences 变更 → 切换 i18n 语言
  useEffect(() => {
    if (locale && i18n.language !== locale) {
      i18n.changeLanguage(locale).catch((error) => {
        console.error('Failed to switch language:', error);
      });
    }
    // 同步 <html lang/dir>：切换语言后必须更新，否则屏幕阅读器/浏览器翻译/a11y
    // 仍按旧语言工作。此前 react 端缺失这一步（vue-element 在 setI18nLanguage 有做）。
    // locale 取值仅 'zh-CN'/'en-US'，均为 LTR，dir 始终为 'ltr'。
    if (locale) {
      document.documentElement.lang = locale;
      document.documentElement.dir = 'ltr';
    }
  }, [locale]);

  // i18n 初始化后，只在第一次反向同步到 store
  useEffect(() => {
    if (i18n.isInitialized && !locale && !initialized.current) {
      initialized.current = true;
      const detected = i18n.language as SupportedLocale;
      setPreferences({ app: { locale: detected } });
    }
  }, [i18n.isInitialized, locale, setPreferences]);

  // 便捷切换方法
  const changeLocale = useCallback(
    (newLocale: SupportedLocale) => {
      setPreferences({ app: { locale: newLocale } });
    },
    [setPreferences],
  );

  return { changeLocale };
};
