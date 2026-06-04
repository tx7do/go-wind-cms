import {defineConfig, type UserConfigExport} from '@tarojs/cli'
import path from 'path'
import {UnifiedViteWeappTailwindcssPlugin} from 'weapp-tailwindcss/vite'

import devConfig from './dev'
import prodConfig from './prod'

// https://taro-docs.jd.com/docs/next/config#defineconfig-辅助函数
export default defineConfig<'vite'>(async (merge, {}) => {
  const baseConfig: UserConfigExport<'vite'> = {
    projectName: 'taro',
    date: '2026-3-17',
    designWidth: 750,
    deviceRatio: {
      640: 2.34 / 2,
      750: 1,
      375: 2,
      828: 1.81 / 2
    },
    sourceRoot: 'src',
    outputRoot: 'dist',
    stylelint: true,
    plugins: [
      "@tarojs/plugin-generator"
    ],
    defineConstants: {
      // 注入环境变量到客户端代码
      API_BASE_URL: JSON.stringify(process.env.TARO_APP_API_BASE_URL || ''),
      APP_TITLE: JSON.stringify(process.env.TARO_APP_TITLE || ''),
      ENABLE_MOCK: JSON.stringify(process.env.TARO_ENABLE_MOCK || 'false'),
      DEFAULT_LOCALE: JSON.stringify(process.env.TARO_DEFAULT_LOCALE || 'zh-CN'),
      TOKEN_KEY: JSON.stringify(process.env.TARO_TOKEN_KEY || ''),
      REFRESH_TOKEN_KEY: JSON.stringify(process.env.TARO_REFRESH_TOKEN_KEY || ''),
      APP_NAMESPACE: JSON.stringify(process.env.TARO_APP_NAMESPACE || 'gowind-cms-app'),
      AES_KEY: JSON.stringify(process.env.TARO_APP_AES_KEY || ''),
    },
    copy: {
      patterns: [],
      options: {}
    },
    framework: 'react',
    compiler: 'vite',
    // weapp-tailwindcss：仅小程序端启用（规范 §8.3 插件环境隔离）
    modifyViteConfig(config) {
      config.plugins = config.plugins || [];
      // 仅小程序端注入 weapp-tailwindcss 转换插件
      if (process.env.TARO_ENV === 'weapp') {
        config.plugins.push(
          UnifiedViteWeappTailwindcssPlugin({
            rem2rpx: true,
          }),
        );
      }
    },
    alias: {
      '@': path.resolve(__dirname, '..', 'src'),
    },
    mini: {
      postcss: {
        pxtransform: {
          enable: true,
          config: {}
        },
        cssModules: {
          enable: false,
          config: {
            namingPattern: 'module',
            generateScopedName: '[name]__[local]___[hash:base64:5]'
          }
        },
        tailwindcss: {
          enable: true,
          config: require('../tailwind.config.js'),
        },
      },
    },
    h5: {
      publicPath: '/',
      staticDirectory: 'static',
      router: {
        mode: 'browser',
        customRoutes: {
          'pages/index/index': '/',
          'pages/about/index': '/about',
          'pages/contact/index': '/contact',
          'pages/disclaimer/index': '/disclaimer',
          'pages/privacy/index': '/privacy',
          'pages/terms/index': '/terms',
          'pages/settings/index': '/settings',
          'pages/user/index': '/user',
          'pages/register/index': '/register',
          'pages/login/index': '/login',
          'pages/post/index': '/post',
          'pages/post/detail/index': '/post/:id',
          'pages/tag/index': '/tag',
          'pages/tag/detail/index': '/tag/:id',
          'pages/category/index': '/category',
          'pages/category/detail/index': '/category/:id',
          'pages/search/index': '/search',
          'pages/404/index': '/404'
        }
      },
      miniCssExtractPluginOption: {
        ignoreOrder: true,
        filename: 'css/[name].[hash].css',
        chunkFilename: 'css/[name].[chunkhash].css'
      },
      postcss: {
        autoprefixer: {
          enable: true,
          config: {}
        },
        cssModules: {
          enable: false,
          config: {
            namingPattern: 'module',
            generateScopedName: '[name]__[local]___[hash:base64:5]'
          }
        },
        tailwindcss: {
          enable: true,
          config: require('../tailwind.config.js'),
        },
      },
    },

  }

  // 生产环境才设置 browserslist 环境变量
  if (process.env.NODE_ENV === 'production') {
    process.env.BROWSERSLIST_ENV = 'production'
  }

  if (process.env.NODE_ENV === 'development') {
    // 本地开发构建配置（不混淆压缩）
    return merge({}, baseConfig, devConfig)
  }
  // 生产构建配置（默认开启压缩混淆等）
  return merge({}, baseConfig, prodConfig)
})
