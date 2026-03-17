## Taro CMS 前端项目

本项目基于 [Taro](https://taro.jd.com/) 框架，采用 React + TypeScript 构建，适用于多端（微信小程序、H5、支付宝等）内容管理系统（CMS）前端。

---

### 目录结构

```
taro/
├── babel.config.js           # Babel 配置
├── commitlint.config.mjs     # Commitlint 配置
├── package.json              # 项目依赖与脚本
├── pnpm-lock.yaml            # pnpm 锁定文件
├── project.config.json       # 小程序项目配置
├── README.md                 # 项目说明文档
├── stylelint.config.mjs      # Stylelint 配置
├── tsconfig.json             # TypeScript 配置
├── config/                   # 构建与环境配置
│   ├── dev.ts                # 开发环境配置
│   ├── index.ts              # 主配置文件
│   └── prod.ts               # 生产环境配置
├── src/                      # 源码目录
│   ├── app.config.ts         # 应用配置
│   ├── app.scss              # 全局样式
│   ├── app.ts                # 应用入口
│   ├── index.html            # H5 入口页面
│   └── pages/                # 页面目录
│       └── index/            # 首页
│           ├── index.config.ts
│           ├── index.scss
│           └── index.tsx
├── types/                    # 类型声明
│   └── global.d.ts
```

---

### 安装与启动

1. **环境要求**
  - Node.js >= 16
  - pnpm >= 8
  - 推荐 Windows 10/11 或 macOS

2. **安装依赖**
   ```bash
   pnpm install
   ```

3. **开发模式启动**
  - 微信小程序：
    ```bash
    pnpm run dev:weapp
    ```
  - H5：
    ```bash
    pnpm run dev:h5
    ```
  - 其他平台见 package.json scripts

4. **构建产物**
  - 微信小程序：
    ```bash
    pnpm run build:weapp
    ```
  - H5：
    ```bash
    pnpm run build:h5
    ```

---

### 常用命令

| 命令               | 说明      |
|------------------|---------|
| pnpm run dev:*   | 各端开发模式  |
| pnpm run build:* | 各端构建产物  |
| pnpm run new     | 新建页面/组件 |

---

### 常见问题

1. **依赖安装失败**
  - 请确保 Node.js 和 pnpm 版本符合要求。
  - 删除 `node_modules` 和 `pnpm-lock.yaml` 后重新安装。

2. **小程序编译异常**
  - 检查 `project.config.json` 配置。
  - 确认微信开发者工具已正确导入 `dist` 目录。

3. **H5 构建失败**
  - 检查端口占用或依赖缺失。
  - 查看终端报错信息，按提示修复。

---

### 联系方式

如有问题或建议，请联系项目维护者。

---

> 本项目遵循 MIT 协议，欢迎贡献代码。

