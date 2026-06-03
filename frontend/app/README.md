# Go-Wind-CMS 前端应用

本项目是 Go-Wind-CMS 的前端部分，包含三个不同技术栈的实现版本：

## 📱 版本介绍

### 1. React 版本 (`/react`)
基于 Next.js 框架构建的现代化 Web 应用。
- **技术栈**: React + Next.js + TypeScript
- **特点**: 
  - 服务端渲染 (SSR) 和静态生成 (SSG)
  - 优秀的 SEO 支持
  - 自动代码分割
  - 内置 API 路由
- **适用场景**: PC 端管理后台、内容展示网站

### 2. Taro 版本 (`/taro`)
基于 Taro 框架的多端统一解决方案。
- **技术栈**: Taro + React + TypeScript
- **特点**:
  - 一套代码多端运行（微信小程序、H5、React Native 等）
  - 完整的组件化开发体验
  - 支持 npm/yarn 包管理
  - 丰富的插件生态
- **适用场景**: 移动端应用、小程序开发

### 3. Vue 版本 (`/vue`)
基于 Nuxt.js 框架构建的 Vue 应用。
- **技术栈**: Vue 3 + Nuxt.js + TypeScript
- **特点**:
  - 服务端渲染能力
  - 自动导入组件和 API
  - 模块化架构
  - 良好的开发体验
- **适用场景**: 快速原型开发、企业级应用

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

## 📁 目录结构

```
frontend/app/
├── react/          # React (Next.js) 版本
├── taro/           # Taro 多端版本
├── vue/            # Vue (Nuxt.js) 版本
└── README.md       # 项目说明文档
```

## 🛠️ 技术特性

- **国际化支持**: 所有版本均支持多语言（中文/英文）
- **状态管理**: 各自采用对应的最佳实践方案
- **API 集成**: 统一的 RESTful API 接口调用
- **响应式设计**: 适配不同屏幕尺寸
- **类型安全**: 全面使用 TypeScript

## 📝 开发规范

- 遵循各框架的最佳实践
- 统一的代码风格和质量标准
- 完善的注释和文档
- 模块化和组件化开发

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进项目！