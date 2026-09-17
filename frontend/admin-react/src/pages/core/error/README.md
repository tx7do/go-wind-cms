# 错误页面组件

本目录包含各种错误和特殊状态的React组件，用于处理常见的HTTP错误和应用状态。

## 组件列表

### Fallback
通用的错误/状态展示组件，支持多种内置状态。

**使用示例：**
```tsx
import { Fallback } from '@/pages/core/error';

// 使用内置状态
<Fallback status="404" />

// 自定义内容
<Fallback 
  status="404"
  title="自定义标题"
  description="自定义描述"
  homePath="/dashboard"
/>

// 使用自定义图片
<Fallback 
  image="/custom-error.png"
  title="出错了"
  description="请稍后重试"
/>
```

### 预定义错误页面

- **Unauthorized (401)** - 未授权页面
- **Forbidden (403)** - 禁止访问页面  
- **NotFound (404)** - 页面未找到
- **InternalError (500)** - 服务器内部错误
- **Offline** - 网络离线页面
- **ComingSoon** - 即将推出页面

**使用示例：**
```tsx
import { NotFound, Unauthorized } from '@/pages/core/error';

// 在路由中使用
<Route path="*" element={<NotFound />} />
<Route path="/unauthorized" element={<Unauthorized />} />
```

## Props

### FallbackProps

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| status | `'401' \| '403' \| '404' \| '500' \| 'coming-soon' \| 'offline'` | `'coming-soon'` | 内置状态类型 |
| title | `string` | `''` | 自定义标题（可选） |
| description | `string` | `''` | 自定义描述（可选） |
| image | `string` | `''` | 自定义图片URL（可选） |
| homePath | `string` | `'/'` | 返回首页的路径 |

## 国际化

组件使用 `react-i18next` 进行国际化，支持的翻译键：

- `common.fallback.unauthorized`
- `common.fallback.forbidden`
- `common.fallback.pageNotFound`
- `common.fallback.internalError`
- `common.fallback.comingSoon`
- `common.fallback.offlineError`
- 以及对应的描述文本...

## 样式

组件使用CSS文件进行样式定义，支持响应式设计：
- 移动端：图片宽度50%
- 平板（≥768px）：图片宽度33.33%
- 桌面（≥1024px）：图片宽度25%

## 图标组件

每个状态都有对应的SVG图标组件，位于 `icons/` 目录：
- Icon401.tsx
- Icon403.tsx
- Icon404.tsx
- Icon500.tsx
- IconComingSoon.tsx
- IconOffline.tsx
