// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'
import { writeFileSync } from 'node:fs'
import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))

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
        // SSG 构建后自动生成根路径 index.html 重定向页
        hooks: {
            'prerender:done'() {
                const outputDir = resolve(__dirname, '.output', 'public')
                writeFileSync(resolve(outputDir, 'index.html'), [
                    '<!DOCTYPE html>',
                    '<html><head>',
                    '<meta charset="utf-8">',
                    '<meta http-equiv="refresh" content="0;url=/zh-CN/">',
                    '<title>Redirecting...</title>',
                    '</head><body>',
                    '<script>location.replace("/zh-CN/"+location.search+location.hash)</script>',
                    '<noscript><meta http-equiv="refresh" content="0;url=/zh-CN/"></noscript>',
                    '<p><a href="/zh-CN/">Click here</a></p>',
                    '</body></html>',
                ].join('\n'))
            },
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
        'pinia-plugin-persistedstate/nuxt',
        'shadcn-nuxt'
    ],
    vite: {
        plugins: [
            tailwindcss(),
        ],
    },
    i18n: {
        langDir: '../locales',
        // seo:false（见 app.vue 的 useLocaleHead）下 baseUrl 不参与 hreflang/canonical 生成，
        // 但 i18n head 守卫要求其非空，否则每轮 SSR 都打 warn；可由 NUXT_PUBLIC_I18N_BASE_URL 覆盖，
        // 以备将来启用 SEO 时注入正确生产源。
        baseUrl: process.env.NUXT_PUBLIC_I18N_BASE_URL || 'http://localhost:3000',
        // 注意：dir 由 @nuxtjs/i18n 运行时读取 LocaleObject.dir 并自动写入 <html dir>。
        //   - 未声明 dir 的 locale 回落到 defaultDirection: 'ltr'（见 @nuxtjs/i18n 模块默认值）
        //   - 阿拉伯语必须显式声明 dir: 'rtl'，否则 HTML 方向不会翻转
        // `language` 为 BCP47 语言标签（@nuxtjs/i18n v10 字段名，旧版叫 `iso`），
        // 运行时用它写入 <html lang>，因此必须设置，否则 <html lang> 不会注入。
        locales: [
            {code: 'zh-CN', language: 'zh-CN', name: '中文', file: 'zh-CN/index.ts'},
            {code: 'en-US', language: 'en-US', name: 'English', file: 'en-US/index.ts'},
            {code: 'ar', language: 'ar', name: 'العربية', file: 'ar/index.ts', dir: 'rtl'},
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
    piniaPluginPersistedstate: {
        // 使用 localStorage 而非默认的 cookies：
        // 1) cookies 4KB 上限会截断 accessMenus/accessRoutes 等大体积状态
        // 2) cookies 随每个 HTTP 请求自动发送，扩大 token 泄露面
        // 3) localStorage 仅客户端可读，避免 SSR 上下文暴露凭证
        storage: 'localStorage',
    },
})
