import {defineConfig, type UserConfigExport} from '@tarojs/cli'
import path from 'path'

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
          enable: false, // 默认为 false，如需使用 css modules 功能，则设为 true
          config: {
            namingPattern: 'module', // 转换模式，取值为 global/module
            generateScopedName: '[name]__[local]___[hash:base64:5]'
          }
        }
      },
    },
    h5: {
      publicPath: '/',
      staticDirectory: 'static',

      // 路由器配置
      router: {
        mode: 'browser', // 可选：'hash' | 'browser' | 'multi'
        customRoutes: { // 自定义路由映射
          'pages/index/index': '/',
          'pages/about/index': '/about',
          'pages/contact/index': '/contact',
          'pages/disclaimer/index': '/disclaimer',
          'pages/privacy/index': '/privacy',
          'pages/settings/index': '/settings',
          'pages/user/index': '/user',
          'pages/post/index': '/post',
          'pages/post/detail/index': '/post/:id',
          'pages/tag/index': '/tag',
          'pages/tag/detail/index': '/tag/:id',
          'pages/category/index': '/category',
          'pages/category/detail/index': '/category/:id',
          'pages/page/index': '/page',
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
          enable: false, // 默认为 false，如需使用 css modules 功能，则设为 true
          config: {
            namingPattern: 'module', // 转换模式，取值为 global/module
            generateScopedName: '[name]__[local]___[hash:base64:5]'
          }
        }
      },
    },
    rn: {
      appName: 'taroDemo',
      postcss: {
        cssModules: {
          enable: false, // 默认为 false，如需使用 css modules 功能，则设为 true
        }
      },
      output: {
        iosSourceMapUrl: '', // sourcemap 文件url
        iosSourcemapOutput: '../taro-native-shell/ios/main.map', // sourcemap 文件输出路径
        iosSourcemapSourcesRoot: '', // 将 sourcemap 资源路径转为相对路径时的根目录
        androidSourceMapUrl: '',
        androidSourcemapOutput: '../taro-native-shell/android/app/src/main/assets/index.android.map',
        androidSourcemapSourcesRoot: '',
        ios: '../taro-native-shell/ios/main.jsbundle',
        iosAssetsDest: '../taro-native-shell/ios',
        android: '../taro-native-shell/android/app/src/main/assets/index.android.bundle',
        androidAssetsDest: '../taro-native-shell/android/app/src/main/res'
      },
    }
  }

  process.env.BROWSERSLIST_ENV = process.env.NODE_ENV

  if (process.env.NODE_ENV === 'development') {
    // 本地开发构建配置（不混淆压缩）
    return merge({}, baseConfig, devConfig)
  }
  // 生产构建配置（默认开启压缩混淆等）
  return merge({}, baseConfig, prodConfig)
})
