import {unref} from "vue";

import {
  i18n,
  languageColumns,
  languageShorts,
  loadLocaleMessages,
  loadLocalesMap,
  loadLocalesMapFromDir,
  setupI18n,
  setI18nLanguage,
} from './i18n';

export {
  type ImportLocaleFn,
  type LocaleSetupOptions,
  type SupportedLanguagesType,
} from './typing';

export {
  localeToLanguageCode,
  languageCodeToLocale,
} from './utils';

export type {CompileError} from '@intlify/core-base';

const $t = i18n.global.t;

/**
 * Get the current locale's API language code (zh-CN or en-US)
 * @returns API language code (zh-CN or en-US)
 */
export function currentLocaleLanguageCode(): string {
  const currentLocale = unref(i18n.global.locale);
  return currentLocale === 'zhCN' ? 'zh-CN' : 'en-US';
}

export {
  $t,
  i18n,
  languageColumns,
  languageShorts,
  loadLocaleMessages,
  loadLocalesMap,
  loadLocalesMapFromDir,
  setupI18n,
  setI18nLanguage,
};
