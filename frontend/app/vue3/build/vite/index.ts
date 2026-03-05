import {dirname, resolve} from 'node:path'
import {fileURLToPath} from 'node:url'

import {unheadVueComposablesImports} from '@unhead/vue'
import legacy from '@vitejs/plugin-legacy'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'

import Components from 'unplugin-vue-components/vite'
import {NaiveUiResolver} from 'unplugin-vue-components/resolvers'
import {VueRouterAutoImports} from 'unplugin-vue-router'
import VueRouter from 'unplugin-vue-router/vite'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'

import mockDevServerPlugin from 'vite-plugin-mock-dev-server'
import {VitePWA} from 'vite-plugin-pwa'
import Sitemap from 'vite-plugin-sitemap'
import VueDevTools from 'vite-plugin-vue-devtools'

import {createViteVConsole} from './vconsole'

export function createVitePlugins() {
  return [
    // https://github.com/posva/unplugin-vue-router
    // ⚠️ 必须要放到Vue()之前
    VueRouter({
      routesFolder: 'src/pages',
      exclude: ['**/components/*.vue', '**/components/**/*.vue', '**/*.components.vue'],
      exttensions: ['vue'],
      dts: 'src/typed-router.d.ts',
    }),

    vue(),

    // https://github.com/jbaubree/vite-plugin-sitemap
    Sitemap(),

    // https://github.com/pengzhanbo/vite-plugin-mock-dev-server
    mockDevServerPlugin(
      {
        prefix: '/app/v1',
      }
    ),

    // https://github.com/antfu/unplugin-vue-components
    Components({
      resolvers: [NaiveUiResolver()],
      extensions: ['vue'],
      deep: true,
      dirs: ['src/components'],
      // globs: ['src/components/*.{vue}'],
      include: [/\.vue$/, /\.vue\?vue/],
      // exclude: [/[\\/]node_modules[\\/]/, /[\\/]\.git[\\/]/, /[\\/]\.nuxt[\\/]/],
      dts: 'src/components.d.ts',
    }),

    // https://github.com/antfu/unplugin-auto-import
    // AutoImport({
    //   include: [
    //     /\.[tj]sx?$/,
    //     /\.vue$/,
    //     /\.vue\?vue/,
    //   ],
    //   imports: [
    //     'vue',
    //     'vitest',
    //     '@vueuse/core',
    //     VueRouterAutoImports,
    //     {
    //       'vue-router/auto': ['useLink'],
    //       '@/utils/i18n': ['i18n', 'locale'],
    //       'vue-i18n': ['useI18n'],
    //     },
    //     {
    //       'naive-ui': [
    //         'useDialog',
    //         'useMessage',
    //         'useNotification',
    //         'useLoadingBar',
    //       ],
    //     },
    //     unheadVueComposablesImports,
    //   ],
    //   dts: 'src/auto-imports.d.ts',
    //   dirs: [
    //     'src/composables',
    //   ],
    // }),

    // https://github.com/intlify/bundle-tools/tree/main/packages/unplugin-vue-i18n
    VueI18nPlugin({
      compositionOnly: true,
      fullInstall: true,
      runtimeOnly: true,
    }),

    legacy({
      targets: ['defaults', 'not IE 11'],
    }),

    // https://github.com/antfu/unocss
    // see uno.config.ts for config
    UnoCSS(),

    // https://github.com/vadxq/vite-plugin-vconsole
    createViteVConsole(),

    // https://github.com/vuejs/devtools-next
    VueDevTools(),

    // https://github.com/antfu/vite-plugin-pwa
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.svg', 'safari-pinned-tab.svg'],
      workbox: {
        // Allow large editor chunks/workers to be precached without aborting build.
        maximumFileSizeToCacheInBytes: 8 * 1024 * 1024,
      },
      manifest: {
        name: 'vue3-naiveui-pc',
        short_name: 'vue3-naiveui-pc',
        theme_color: '#ffffff',
        icons: [
          {
            src: '/pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: '/pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
          {
            src: '/pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any maskable',
          },
        ],
      },
    }),
  ]
}
