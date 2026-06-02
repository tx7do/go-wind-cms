// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: {enabled: true},
    components: [
        {
            path: '~/components',
            ignore: ['**/index.ts'],
        },
    ],
    css: [
        '~/assets/css/main.css',
    ],
    modules: [
        '@nuxtjs/i18n',
        '@pinia/nuxt',
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
        langDir: '../locales',
        locales: [
            {code: 'zh-CN', iso: 'zh-CN', name: '中文', file: 'zh-CN/index.ts'},
            {code: 'en-US', iso: 'en-US', name: 'English', file: 'en-US/index.ts'}
        ],
        defaultLocale: 'zh-CN',
        strategy: 'prefix',
        detectBrowserLanguage: false,
    },
    runtimeConfig: {
        public: {
            apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8000',
            enableMock: process.env.NUXT_PUBLIC_ENABLE_MOCK === 'true',
            appTitle: process.env.NUXT_PUBLIC_APP_TITLE || 'GoWind Content Hub',
            appDescription: process.env.NUXT_PUBLIC_APP_DESCRIPTION || 'A modern Content Hub built with Nuxt',
            defaultLocale: process.env.NUXT_PUBLIC_DEFAULT_LOCALE || 'zh-CN',
            tokenKey: process.env.NUXT_PUBLIC_TOKEN_KEY || 'access_token',
            refreshTokenKey: process.env.NUXT_PUBLIC_REFRESH_TOKEN_KEY || 'refresh_token',
            aesKey: process.env.NUXT_PUBLIC_AES_KEY || '',
        },
    },
    pinia: {
        storesDirs: ['./app/stores/**'],
    },
})
