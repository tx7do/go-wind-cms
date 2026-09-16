// zh-CN imports
import zhCN_authentication from '../../messages/zh-CN/authentication.json';
import zhCN_app from '../../messages/zh-CN/app.json';
import zhCN_page from '../../messages/zh-CN/page.json';
import zhCN_navbar from '../../messages/zh-CN/navbar.json';
import zhCN_menu from '../../messages/zh-CN/menu.json';
import zhCN_enum from '../../messages/zh-CN/enum.json';
import zhCN_error from '../../messages/zh-CN/error.json';
import zhCN_component from '../../messages/zh-CN/component.json';
import zhCN_common from '../../messages/zh-CN/common.json';
import zhCN_comment from '../../messages/zh-CN/comment.json';
import zhCN_cms from '../../messages/zh-CN/cms.json';
import zhCN_ui from '../../messages/zh-CN/ui.json';
import zhCN_settings from '../../messages/zh-CN/settings.json';
import zhCN_preferences from '../../messages/zh-CN/preferences.json';
// en-US imports
import enUS_cms from '../../messages/en-US/cms.json';
import enUS_authentication from '../../messages/en-US/authentication.json';
import enUS_app from '../../messages/en-US/app.json';
import enUS_preferences from '../../messages/en-US/preferences.json';
import enUS_page from '../../messages/en-US/page.json';
import enUS_navbar from '../../messages/en-US/navbar.json';
import enUS_menu from '../../messages/en-US/menu.json';
import enUS_enum from '../../messages/en-US/enum.json';
import enUS_error from '../../messages/en-US/error.json';
import enUS_component from '../../messages/en-US/component.json';
import enUS_common from '../../messages/en-US/common.json';
import enUS_settings from '../../messages/en-US/settings.json';
import enUS_ui from '../../messages/en-US/ui.json';
import enUS_comment from '../../messages/en-US/comment.json';

export const locales = ['zh-CN', 'en-US'] as const;

export const DEFAULT_LANGUAGE = 'zh-CN';
export const DEFAULT_TIME_ZONE = 'Asia/Shanghai';

export type Locale = (typeof locales)[number];
export const defaultLocale: Locale = DEFAULT_LANGUAGE as Locale;

export function validateLocale(locale?: string): Locale {
  if (!locale || !locales.includes(locale as Locale)) {
    return defaultLocale;
  }
  return locale as Locale;
}


export const allMessages: Record<Locale, Record<string, Record<string, unknown>>> = {
  'zh-CN': {
    authentication: zhCN_authentication,
    app: zhCN_app,
    page: zhCN_page,
    navbar: zhCN_navbar,
    menu: zhCN_menu,
    enum: zhCN_enum,
    error: zhCN_error,
    component: zhCN_component,
    common: zhCN_common,
    comment: zhCN_comment,
    cms: zhCN_cms,
    ui: zhCN_ui,
    settings: zhCN_settings,
    preferences: zhCN_preferences,
  },
  'en-US': {
    cms: enUS_cms,
    authentication: enUS_authentication,
    app: enUS_app,
    preferences: enUS_preferences,
    page: enUS_page,
    navbar: enUS_navbar,
    menu: enUS_menu,
    enum: enUS_enum,
    error: enUS_error,
    component: enUS_component,
    common: enUS_common,
    settings: enUS_settings,
    ui: enUS_ui,
    comment: enUS_comment,
  }
};

/**
 * 递归将嵌套对象的所有叶子节点路径提取出来
 * 例如：{ login: { brand_title: '标题' } } => { 'login.brand_title': '标题' }
 */
function flattenObject(obj: Record<string, any>, prefix = ''): Record<string, string> {
  return Object.keys(obj).reduce((acc: Record<string, string>, key: string) => {
    const newKey = prefix ? `${prefix}.${key}` : key;
    if (typeof obj[key] === 'object' && obj[key] !== null && !Array.isArray(obj[key])) {
      // 递归处理嵌套对象
      Object.assign(acc, flattenObject(obj[key], newKey));
    } else {
      acc[newKey] = obj[key];
    }
    return acc;
  }, {});
}

/**
 * 将资源转换为 react-i18next 支持的格式
 * 将所有消息扁平化到一个命名空间中，支持直接访问（如 t('login.brand_title')）
 */
export function getFlattenedMessages(): Record<Locale, Record<string, string>> {
  const result: Record<Locale, Record<string, string>> = {
    'zh-CN': {},
    'en-US': {},
  };

  (locales as unknown as Locale[]).forEach(locale => {
    const messages = allMessages[locale];
    // 将整个消息对象扁平化到一个命名空间中
    result[locale] = flattenObject(messages);
  });

  return result;
}
