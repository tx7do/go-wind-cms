// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: {enabled: true},
    css: [
        '~/assets/css/main.css',
    ],
    modules: [
        '@nuxtjs/i18n',
    ],
    vite: {
        plugins: [
            tailwindcss(),
        ],
        optimizeDeps: {
            include: [
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
        strategy: 'prefix_except_default',
        detectBrowserLanguage: false,
    }
})
