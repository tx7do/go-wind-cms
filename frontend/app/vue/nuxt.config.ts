// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: {enabled: true},
    css: [
    ],
    modules: [
        '@nuxtjs/i18n',
    ],
    components: [
        {
            path: 'node_modules/naive-ui',
            prefix: 'N',
            extensions: [
                '.vue',
                '.ts'
            ],
        },
    ],
    build: {
        transpile: [
            'naive-ui',
            'vueuc',
        ]
    },
    vite: {
        optimizeDeps: {
            include: [
                'naive-ui',
                'vueuc',
                'date-fns-tz',
                'highlight.js/lib/core'
            ],
        },
    },
    i18n: {
        locales: [
            {code: 'zh', iso: 'zh-CN', name: '中文', file: 'zh.json'},
            {code: 'en', iso: 'en-US', name: 'English', file: 'en.json'}
        ],
        defaultLocale: 'zh',
        strategy: 'prefix_except_default', // /en/about 或 /about
        detectBrowserLanguage: false, // CMS 前台通常由路由或用户手动切换
    }
})
