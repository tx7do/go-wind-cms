# GoWind CMS — Vue Frontend

基于 [Nuxt](https://nuxt.com/) 的现代 CMS 前端，使用 Tailwind CSS + shadcn-vue 构建。

## 环境准备

- Node.js >= 18
- pnpm（推荐）

## 安装依赖

```bash
pnpm install
```

## 开发

启动开发服务器（默认端口 `3000`）：

```bash
pnpm dev
```

## 构建与部署

### 静态站点（推荐，部署到 Nginx）

生成纯静态站点到 `.output/public/`：

```bash
pnpm generate
```

构建产物可直接部署到 Nginx 等 Web 服务器。根路径 `/` 会自动重定向到默认语言 `/zh-CN/`。

### Node.js 服务端

构建 Node.js 服务端应用：

```bash
pnpm build
```

本地预览：

```bash
pnpm preview
```

## 项目结构

```
app/
├── api/            # API 服务层
├── assets/css/     # 全局样式
├── components/     # UI 组件
├── composables/    # 组合式函数
├── constants/      # 常量
├── core/           # 核心模块（偏好设置、存储、传输）
├── pages/          # 页面路由
├── plugins/        # Nuxt 插件
└── stores/         # Pinia 状态管理
i18n/locales/       # 多语言翻译文件
public/             # 静态资源
```
