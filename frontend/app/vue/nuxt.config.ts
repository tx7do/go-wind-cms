// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
    compatibilityDate: '2025-07-15',
    devtools: {enabled: true},

    // 静态站点生成配置
    ssr: true,
    nitro: {
        prerender: {
            // 预渲染静态页面，动态路由页面由客户端渲染
            routes: ['/'],
            crawlLinks: true,
        },
    },
    // SPA fallback：未预渲染的路由返回 200 + index.html（由 nginx 配合）
    app: {
        baseURL: '/',
        buildAssetsDir: '_nuxt',
    },
    components: [
        {
            path: '~/components',
            extensions: ['.vue'],
            ignore: ['**/index.ts'],
        },
    ],
    css: [
        '~/assets/css/main.css',
        'katex/dist/katex.min.css',
    ],
    modules: [
        '@nuxtjs/i18n',
        '@pinia/nuxt',
        'shadcn-nuxt'
    ],
    vite: {
        plugins: [
            tailwindcss(),
        ],
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
    typescript: {
        tsConfig: {
            compilerOptions: {
                resolveJsonModule: true,
            },
        },
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
