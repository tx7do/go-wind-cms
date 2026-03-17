# GoWind Content Hub Vue3 App

## 项目简介

这是一个前后端分离的内容中台系统（Headless CMS），后端基于 Go（go-kratos），后台前端基于 Vue3（Vben），前台提供 React 与 Vue3 双版本，支持多端渲染与跨平台发布。

- 在线 Demo：[https://cms.gowind.cloud](https://cms.gowind.cloud)

## 核心能力

- 多租户、数据权限、动态 API、任务调度、多数据库兼容、高级查询
- 集成统一认证、国际化、接口文档等企业级特性
- 支持富文本、Markdown、JSON 编辑器
- 多语言内容管理与多端数据同步
- 快速构建内容管理、门户、社区、知识库等应用

## 前端主要特性

- ⚡️ Vue 3 + Vite 5 + pnpm + esbuild，极速开发体验
- 🗂 文件路由自动生成
- 📦 组件自动加载
- 🍍 Pinia 状态管理
- 📲 PWA 支持
- 🎨 UnoCSS 原子化 CSS
- 🌍 国际化（I18n）开箱即用
- 🔥 `<script setup>` 语法
- 📥 API 自动加载
- 💪 TypeScript 全面支持
- ⚙️ Vitest 单元测试
- 💾 本地数据模拟（Mock）
- 🌈 Git hooks（eslint 检测、提交规范）
- 🪶 Naive UI 组件库
- 🔭 vConsole 移动端调试
- 🌓 深色模式支持
- ☁️ 零配置 Netlify 部署

## 技术栈与插件

- UI：UnoCSS、Naive UI
- 路由：Vue Router、unplugin-vue-router
- 状态管理：Pinia、pinia-plugin-persistedstate
- 国际化：Vue I18n、unplugin-vue-i18n
- 组件/自动导入：unplugin-vue-components、unplugin-auto-import
- Mock：vite-plugin-mock-dev-server
- PWA：vite-plugin-pwa
- 头信息：@unhead/vue
- 测试：Vitest

## 开发规范

- 使用 `<script setup>` SFC 语法
- ESLint 配置为 @antfu/eslint-config
- TypeScript 支持

## 快速开始

### 安装依赖

```bash
pnpm i # 没装过 pnpm 可先 npm install -g pnpm
```

### 本地开发

```bash
pnpm dev
```

访问：<http://localhost:3000>

### 构建

```bash
pnpm build
```

构建产物在 dist 文件夹。

### 部署

推荐使用 Netlify，支持零配置部署。

## 说明

- Node 18+ 环境
- 修改 LICENSE、index.html、vite.config.ts、public/favicon 等信息以适配你的项目

## License

MIT License
