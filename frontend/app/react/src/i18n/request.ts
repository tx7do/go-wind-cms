import {getRequestConfig} from 'next-intl/server';
import {routing} from '../app/[locale]/routing';

// 导入所有翻译文件
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

export default getRequestConfig(async ({requestLocale}) => {
    let locale = await requestLocale;

    // 确保使用有效的 locale
    if (!locale || !routing.locales.includes(locale as 'zh-CN' | 'en-US')) {
        locale = routing.defaultLocale;
    }

    // 根据 locale 返回对应的翻译消息
    const messages = locale === 'zh-CN' ? {
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
    } : {
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
    };

    return {
        locale,
        messages
    };
});
