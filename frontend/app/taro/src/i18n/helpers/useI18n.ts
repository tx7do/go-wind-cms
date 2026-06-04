import {useTranslation} from 'react-i18next';
import {usePreferencesStore} from '@/core/preferences/store';
import {useLoadingStore} from '@/store/core/loading/store';
import type {SupportedLanguagesType} from '@/core/preferences/types/layout';


/**
 * useI18n hook - 增强的多语言支持
 * @param namespace 翻译模块名（如 'common', 'article', 'category'）
 */
export function useI18n(namespace: string = 'common') {
  const {t, i18n} = useTranslation(namespace);
  const preferencesStore = usePreferencesStore();
  const loadingStore = useLoadingStore();

  // 同步 locale 到 preferences store
  if (preferencesStore.preferences.app.locale !== i18n.language) {
    preferencesStore.setPreferences({app: {locale: i18n.language as SupportedLanguagesType}});
  }

  // 切换语言
  const changeLocale = (targetLocale: string) => {
    console.log('[useI18n] Start locale change, showing loading...');
    loadingStore.start();

    // 切换语言
    i18n.changeLanguage(targetLocale).then(() => {
      // 保存到 preferences store（自动持久化）
      preferencesStore.setPreferences({app: {locale: targetLocale as SupportedLanguagesType}});
      loadingStore.finish();
      console.log('[useI18n] Locale changed successfully');
    });
  };

  return {t, locale: i18n.language, changeLocale};
}
