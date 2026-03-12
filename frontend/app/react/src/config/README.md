# 环境变量配置说明

本项目使用环境变量来管理不同环境的配置，支持开发环境和生产环境的自动切换。

## 📋 环境变量列表

### 必需的环境变量

| 变量名 | 说明 | 示例值 |
|--------|------|--------|
| `NEXT_PUBLIC_API_BASE_URL` | API 接口基础地址 | `http://localhost:8000` |
| `NEXT_PUBLIC_APP_TITLE` | 应用标题 | `Kratos CMS` |
| `NEXT_PUBLIC_DEFAULT_LOCALE` | 默认语言 | `zh-CN` |

### 可选的环境变量

| 变量名 | 说明 | 默认值 | 示例值 |
|--------|------|--------|--------|
| `NEXT_PUBLIC_ENABLE_MOCK` | 是否启用 Mock 数据 | `false` | `true` / `false` |
| `NEXT_PUBLIC_TOKEN_KEY` | Token 存储键名 | `access_token` | `my_app_token` |
| `NEXT_PUBLIC_REFRESH_TOKEN_KEY` | Refresh Token 存储键名 | `refresh_token` | `my_app_refresh_token` |

## 🔧 配置文件说明

项目包含以下环境配置文件：

- `.env.development` - 开发环境配置
- `.env.production` - 生产环境配置
- `.env.local` - 本地个人配置（不提交到 Git）

## 📝 使用示例

### 1. 基础使用

```typescript
import {env} from '@/config/env';

// 读取 API 地址
console.log('API Base URL:', env.apiBaseUrl);

// 读取应用标题
console.log('App Title:', env.appTitle);

// 检查是否启用 Mock
if (env.enableMock) {
  console.log('Mock is enabled');
}
```

### 2. 环境判断

```typescript
import {env} from '@/config/env';

// 开发环境
if (env.isDev) {
  console.log('开发环境 - 启用调试功能');
}

// 生产环境
if (env.isProd) {
  console.log('生产环境 - 禁用调试功能');
}
```

### 3. 创建 API 客户端

```typescript
import {env} from '@/config/env';
import axios from 'axios';

const apiClient = axios.create({
  baseURL: env.apiBaseUrl,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export default apiClient;
```

### 4. 认证服务中使用

```typescript
import {env} from '@/config/env';

class AuthService {
  private tokenKey = env.tokenKey;
  private refreshTokenKey = env.refreshTokenKey;

  // 存储 token
  setToken(token: string, refreshToken: string) {
    localStorage.setItem(this.tokenKey, token);
    localStorage.setItem(this.refreshTokenKey, refreshToken);
  }

  // 获取 token
  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  // 获取 refresh token
  getRefreshToken(): string | null {
    return localStorage.getItem(this.refreshTokenKey);
  }

  // 清除 token
  clearTokens() {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem(this.refreshTokenKey);
  }
}

export default new AuthService();
```

### 5. 在 React 组件中使用

```tsx
'use client';

import {env} from '@/config/env';

export function Header() {
  return (
    <header>
      <h1>{env.appTitle}</h1>
      <p>当前环境：{env.nodeEnv}</p>
      {env.isDev && <span className="dev-badge">DEV</span>}
    </header>
  );
}
```

### 6. 生产环境验证

```typescript
import {env, validateEnv} from '@/config/env';

// 在生产环境启动时验证环境变量
export function initializeApp() {
  if (env.isProd) {
    try {
      validateEnv();
      console.log('环境变量验证通过');
    } catch (error) {
      console.error('环境变量配置错误:', error);
      throw error;
    }
  }
}
```

## ⚠️ 注意事项

### 1. 环境变量命名规范

- **客户端可用**：以 `NEXT_PUBLIC_` 开头的环境变量可以在浏览器端使用
- **服务端专用**：不以 `NEXT_PUBLIC_` 开头的环境变量只能在服务端使用

```bash
# ✅ 客户端可用
NEXT_PUBLIC_API_BASE_URL=http://localhost:8000

# ❌ 仅服务端可用
DATABASE_URL=postgresql://...
```

### 2. 安全性

```bash
# ✅ 可以暴露在客户端
NEXT_PUBLIC_API_BASE_URL=http://localhost:8000
NEXT_PUBLIC_APP_TITLE=Kratos CMS

# ❌ 不要暴露敏感信息
# NEXT_PUBLIC_SECRET_KEY=xxx  # 错误做法
NEXTAUTH_SECRET=xxx           # 正确：不在客户端使用
```

### 3. 类型安全

项目提供了完整的 TypeScript 类型定义：

```typescript
import {env} from '@/config/env';

// ✅ 类型安全，有智能提示
env.apiBaseUrl;      // string
env.enableMock;      // boolean
env.unknownField;    // ❌ TypeScript 会报错
```

### 4. 环境变量验证

使用 `validateEnv()` 函数验证必需的环境变量：

```typescript
import {validateEnv} from '@/config/env';

try {
  validateEnv();
  console.log('✓ 所有必需的环境变量都已配置');
} catch (error) {
  console.error('✗ 环境变量配置不完整', error);
}
```

## 🚀 快速开始

### 步骤 1：复制示例配置

```bash
cp .env.example .env.local
```

### 步骤 2：修改配置

编辑 `.env.local` 文件，填入适合你本地环境的值：

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8000
NEXT_PUBLIC_APP_TITLE=我的 CMS 系统
NEXT_PUBLIC_DEFAULT_LOCALE=zh-CN
NEXT_PUBLIC_ENABLE_MOCK=false
```

### 步骤 3：启动项目

```bash
# 开发环境
pnpm dev

# 生产环境
pnpm build
pnpm start
```

## 📚 相关文件

- `env.ts` - 环境变量配置和导出
- `env.example.ts` - 使用示例（本文件已迁移到 README）
- `index.ts` - 配置统一导出

## 🔗 参考资料

- [Next.js 环境变量文档](https://nextjs.org/docs/basic-features/environment-variables)
- [TypeScript 类型系统](https://www.typescriptlang.org/docs/handbook/type-system.html)
