import {getRequestConfig} from 'next-intl/server';

// zh-CN imports
import zhCN_authentication from './locales/zh-CN/authentication.json';
import zhCN_app from './locales/zh-CN/app.json';
import zhCN_page from './locales/zh-CN/page.json';
import zhCN_navbar from './locales/zh-CN/navbar.json';
import zhCN_menu from './locales/zh-CN/menu.json';
import zhCN_enum from './locales/zh-CN/enum.json';
import zhCN_component from './locales/zh-CN/component.json';
import zhCN_common from './locales/zh-CN/common.json';
import zhCN_comment from './locales/zh-CN/comment.json';
import zhCN_cms from './locales/zh-CN/cms.json';
import zhCN_ui from './locales/zh-CN/ui.json';
import zhCN_settings from './locales/zh-CN/settings.json';
import zhCN_preferences from './locales/zh-CN/preferences.json';
// en-US imports
import enUS_cms from './locales/en-US/cms.json';
import enUS_authentication from './locales/en-US/authentication.json';
import enUS_app from './locales/en-US/app.json';
import enUS_preferences from './locales/en-US/preferences.json';
import enUS_page from './locales/en-US/page.json';
import enUS_navbar from './locales/en-US/navbar.json';
import enUS_menu from './locales/en-US/menu.json';
import enUS_enum from './locales/en-US/enum.json';
import enUS_component from './locales/en-US/component.json';
import enUS_common from './locales/en-US/common.json';
import enUS_settings from './locales/en-US/settings.json';
import enUS_ui from './locales/en-US/ui.json';

export const locales = ['zh-CN', 'en-US'] as const;
export type Locale = (typeof locales)[number];
export const DEFAULT_LANGUAGE = 'zh-CN';
export const defaultLocale: Locale = DEFAULT_LANGUAGE as Locale;

export function validateLocale(locale?: string): Locale {
    if (!locale || !locales.includes(locale as Locale)) {
        return defaultLocale;
    }
    return locale as Locale;
}

const allMessages: Record<Locale, Record<string, Record<string, unknown>>> = {
    'zh-CN': {
        authentication: zhCN_authentication,
        app: zhCN_app,
        page: zhCN_page,
        navbar: zhCN_navbar,
        menu: zhCN_menu,
        enum: zhCN_enum,
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
        component: enUS_component,
        common: enUS_common,
        settings: enUS_settings,
        ui: enUS_ui,
    }
};

export default getRequestConfig(async ({locale}) => {
    const validatedLocale = validateLocale(locale);
    const messages = allMessages[validatedLocale] || allMessages[defaultLocale];
    return {locale: validatedLocale, messages};
});
