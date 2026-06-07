# Go-Wind-CMS 前端应用

本项目是 Go-Wind-CMS 的前端部分，包含四个不同技术栈的实现版本：

## 📱 版本介绍

### 1. React 版本 (`/react`)
基于 Next.js 框架构建的现代化 Web 应用。
- **技术栈**: React + Next.js + TypeScript
- **特点**:
  - 服务端渲染 (SSR) 和静态生成 (SSG)
  - 国际化路由 (`/[locale]/...`)
  - 自动代码分割与内置 API 路由
  - shadcn/ui 组件库
- **页面模块**: 首页、文章详情、分类列表、标签管理、用户中心、登录/注册、关于、联系我们、免责声明、隐私政策、使用条款、设置
- **适用场景**: PC 端管理后台、内容展示网站

### 2. Taro 版本 (`/taro`)
基于 Taro 框架的多端统一解决方案。
- **技术栈**: Taro + React + TypeScript
- **特点**:
  - 一套代码多端运行（微信小程序、H5、React Native 等）
  - 完整的组件化开发体验
  - 支持 npm/yarn 包管理
  - 丰富的插件生态
- **页面模块**: 首页、文章详情、分类列表、标签管理、用户中心、关于、联系我们、免责声明、隐私政策、使用条款、设置
- **适用场景**: 移动端应用、小程序开发

### 3. Vue 版本 (`/vue`)
基于 Nuxt.js 框架构建的 Vue 应用。
- **技术栈**: Vue 3 + Nuxt.js + TypeScript
- **特点**:
  - 服务端渲染能力
  - 自动导入组件和 API
  - 模块化架构
  - 良好的开发体验
- **页面模块**: 首页、文章详情、分类列表、标签管理、用户中心、登录/注册、关于、联系我们、免责声明、隐私政策、使用条款、设置
- **适用场景**: 快速原型开发、企业级应用

### 4. Flutter 版本 (`/flutter_app`)
基于 Flutter 框架构建的跨平台原生应用。
- **技术栈**: Flutter + Dart
- **架构**: BLoC 状态管理 + go_router 路由 + 特性模块化 (Feature-First)
- **核心依赖**: dio (网络)、sqflite (本地存储)、flutter_bloc (状态管理)、go_router (路由)、retrofit (API 生成)、cached_network_image (图片缓存)
- **特点**:
  - 真正的跨平台原生体验 (Android / iOS / Web / macOS / Windows / Linux)
  - Feature-First 模块化架构
  - BLoC/Cubit 状态管理模式
  - 支持国际化 (中/英)、主题切换 (亮色/暗色)
  - OpenAPI 代码自动生成 (swagger_parser + retrofit)
- **页面模块**: 首页、探索、文章详情、文章列表、分类列表、标签列表/标签动态、搜索、书签、我的评论、个人资料、关于、法律信息、设置
- **适用场景**: 原生移动应用、跨平台桌面应用

## 🚀 快速开始

### React 版本
```bash
cd react
pnpm install
pnpm dev
```

### Taro 版本
```bash
cd taro
yarn install
yarn dev:weapp  # 或其他平台
```

### Vue 版本
```bash
cd vue
pnpm install
pnpm dev
```

### Flutter 版本
```bash
cd flutter_app
flutter pub get
flutter run  # 或 flutter run -d chrome / flutter run -d windows
```

## 📁 目录结构

```
frontend/app/
├── react/          # React (Next.js) 版本 - PC 端 Web
├── taro/           # Taro 多端版本 - 小程序/H5
├── vue/            # Vue (Nuxt.js) 版本 - Web 应用
├── flutter_app/    # Flutter 版本 - 跨平台原生应用
└── README.md       # 项目说明文档
```

## 🛠️ 技术特性

- **国际化支持**: 所有版本均支持多语言（中文/英文）
- **状态管理**: 各自采用对应的最佳实践方案
- **API 集成**: 统一的 RESTful API 接口调用
- **响应式设计**: 适配不同屏幕尺寸
- **类型安全**: 全面使用 TypeScript / Dart 强类型
- **主题支持**: 亮色/暗色主题切换

## 📝 开发规范

- 遵循各框架的最佳实践
- 统一的代码风格和质量标准
- 完善的注释和文档
- 模块化和组件化开发

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进项目！