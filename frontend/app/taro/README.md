## Taro CMS 前端项目

基于 [Taro 4](https://taro.jd.com/) + React + TypeScript 的多端 CMS 前端，支持 H5、微信/支付宝/抖音小程序。

---

### 快速开始

```bash
# 安装依赖
pnpm install

# H5 开发（默认 http://localhost:10086）
pnpm dev:h5

# 微信小程序开发（dist/ 导入微信开发者工具）
pnpm dev:weapp
```

### 构建产物

```bash
pnpm build:h5      # H5 生产构建 → dist/
pnpm build:weapp   # 微信小程序 → dist/
```

### 环境变量

编辑 `.env.development`：

```env
TARO_APP_API_BASE_URL=https://api.cms.gowind.cloud   # 后端 API
TARO_APP_TITLE='GoWind Content Hub'                    # 应用标题
TARO_DEFAULT_LOCALE=zh-CN                              # 默认语言
```

---

### 常用命令

| 命令                          | 说明       |
|-----------------------------|----------|
| `pnpm dev:h5`               | H5 开发模式  |
| `pnpm dev:weapp`            | 微信小程序开发  |
| `pnpm dev:alipay`           | 支付宝小程序开发 |
| `pnpm dev:tt`               | 抖音小程序开发  |
| `pnpm dev:swan`             | 百度小程序开发  |
| `pnpm dev:qq`               | QQ 小程序开发 |
| `pnpm dev:jd`               | 京东小程序开发  |
| `pnpm dev:harmony-hybrid`   | 鸿蒙混合开发   |
| `pnpm build:h5`             | H5 生产构建  |
| `pnpm build:weapp`          | 微信小程序构建  |
| `pnpm build:alipay`         | 支付宝小程序构建 |
| `pnpm build:tt`             | 抖音小程序构建  |
| `pnpm build:swan`           | 百度小程序构建  |
| `pnpm build:qq`             | QQ 小程序构建 |
| `pnpm build:jd`             | 京东小程序构建  |
| `pnpm build:harmony-hybrid` | 鸿蒙混合构建   |
| `pnpm run new`              | 新建页面/组件  |

---

> 详细技术文档见 [TECHNICAL_GUIDE.md](./TECHNICAL_GUIDE.md)
