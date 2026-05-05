import type {Locale} from 'vue-i18n';

import type {
  ImportLocaleFn,
  LoadMessageFn,
  LocaleSetupOptions,
  SupportedLanguagesType,
} from './typing';

import {createI18n} from 'vue-i18n'

import type {App} from "vue";
import {unref} from "vue";
import {preferences} from "@/preferences";

const i18n = createI18n<false>({
  legacy: false,
  globalInjection: true,
  locale: 'zhCN',
  messages: {},
})

const modules = import.meta.glob('./langs/**/*.json');

const localesMap = loadLocalesMapFromDir(
  /\.\/langs\/([^/]+)\/(.*)\.json$/,
  modules,
);
let loadMessages: LoadMessageFn = async () => ({});

/**
 * Load locale modules
 * @param modules
 */
function loadLocalesMap(modules: Record<string, () => Promise<unknown>>) {
  const localesMap: Record<Locale, ImportLocaleFn> = {};

  for (const [path, loadLocale] of Object.entries(modules)) {
    const key = path.match(/([\w-]*)\.(json)/)?.[1];
    if (key) {
      localesMap[key] = loadLocale as ImportLocaleFn;
    }
  }
  return localesMap;
}

/**
 * Load locale modules with directory structure
 * @param regexp - Regular expression to match language and file names
 * @param modules - The modules object containing paths and import functions
 * @returns A map of locales to their corresponding import functions
 */
function loadLocalesMapFromDir(
  regexp: RegExp,
  modules: Record<string, () => Promise<unknown>>,
): Record<Locale, ImportLocaleFn> {
  const localesRaw: Record<Locale, Record<string, () => Promise<unknown>>> = {};
  const localesMap: Record<Locale, ImportLocaleFn> = {};

  // Iterate over the modules to extract language and file names
  for (const path in modules) {
    const match = path.match(regexp);

    if (match) {
      const [_, locale, fileName] = match;
      if (locale && fileName) {
        if (!localesRaw[locale]) {
          localesRaw[locale] = {};
        }
        if (modules[path]) {
          localesRaw[locale][fileName] = modules[path];
        }
      }
    }
  }

  // Convert raw locale data into async import functions
  for (const [locale, files] of Object.entries(localesRaw)) {
    localesMap[locale] = async () => {
      const messages: Record<string, any> = {};
      for (const [fileName, importFn] of Object.entries(files)) {
        messages[fileName] = ((await importFn()) as any)?.default;
      }
      return {default: messages};
    };
  }

  return localesMap;
}

/**
 * Set i18n language
 * @param locale
 */
function setI18nLanguage(locale: Locale) {
  i18n.global.locale.value = locale;

  document?.querySelector('html')?.setAttribute('lang', locale);
}

async function setupI18n(app: App, options: LocaleSetupOptions = {}) {
  let {defaultLocale = 'zhCN'} = options;

  // Try to load saved locale from cache
  try {
    const savedLocale = preferences.app.locale;

    if (savedLocale && (savedLocale === 'zhCN' || savedLocale === 'enUS')) {
      defaultLocale = savedLocale as SupportedLanguagesType;
    }
  } catch (error) {
    console.warn('[i18n] Failed to load saved locale from cache:', error);
  }

  // app可以自行扩展一些第三方库和组件库的国际化
  loadMessages = options.loadMessages || (async () => ({}));

  app.use(i18n);

  await loadLocaleMessages(defaultLocale);

  // Missing key warnings are handled by vue-i18n's built-in missingWarn option
  if (options.missingWarn) {
    console.log('[i18n] Missing key warnings are enabled');
  }
}

/**
 * Load locale messages
 * @param lang
 */
async function loadLocaleMessages(lang: SupportedLanguagesType) {
  // Check if messages are already loaded
  const existingMessages = i18n.global.getLocaleMessage(lang);
  const hasMessages = existingMessages && Object.keys(existingMessages).length > 0;

  // If already loaded and locale matches, just set the language
  if (hasMessages && unref(i18n.global.locale) === lang) {
    return setI18nLanguage(lang);
  }

  const message = await localesMap[lang]?.();

  if (message?.default) {
    i18n.global.setLocaleMessage(lang, message.default);
  }

  const mergeMessage = await loadMessages(lang);
  i18n.global.mergeLocaleMessage(lang, mergeMessage);

  return setI18nLanguage(lang);
}

/** 多语言 picker columns */
const languageColumns = [
  {label: '简体中文', key: 'zhCN'},
  {label: 'English', key: 'enUS'},
];

const languageShorts = [
  {label: '简', key: 'zhCN'},
  {label: 'EN', key: 'enUS'},
];

export {
  i18n,
  languageColumns,
  languageShorts,
  loadLocaleMessages,
  loadLocalesMap,
  loadLocalesMapFromDir,
  setupI18n,
  setI18nLanguage,
};
