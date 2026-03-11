import {defineRouting} from 'next-intl/routing';

export const routing = defineRouting({
    // 支持的语言列表
    locales: ['zh-CN', 'en-US'],

    // 默认语言
    defaultLocale: 'zh-CN',

    // URL 中始终包含语言前缀
    localePrefix: 'always',
});
