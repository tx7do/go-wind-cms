import {useTranslation} from 'react-i18next';
import {useDispatch} from 'react-redux';

import { useLanguageStore } from '@/store/core/language/hooks';
import { startLoading, finishLoading } from '@/store/core/loading/slice';


/**
 * useI18n hook - 增强的多语言支持
 * @param namespace 翻译模块名（如 'common', 'article', 'category'）
 */
export function useI18n(namespace: string = 'common') {
  const {t, i18n} = useTranslation(namespace);
  const dispatch = useDispatch();

  // 同步 locale 到 Redux store
  const languageStore = useLanguageStore();
  if (languageStore.language.locale !== i18n.language) {
    languageStore.setLocale(i18n.language);
  }

  // 切换语言
  const changeLocale = (targetLocale: string) => {
    console.log('[useI18n] Start locale change, showing loading...');
    dispatch(startLoading());

    // 切换语言
    i18n.changeLanguage(targetLocale).then(() => {
      // 保存到 localStorage
      if (typeof window !== 'undefined') {
        localStorage.setItem('locale', targetLocale);
      }
      dispatch(finishLoading());
      console.log('[useI18n] Locale changed successfully');
    });
  };

  return {t, locale: i18n.language, changeLocale};
}
