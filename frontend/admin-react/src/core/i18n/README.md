# 国际化模块 (i18n)

## 概述

本模块基于 [i18next](https://www.i18next.com/) + [react-i18next](https://react.i18next.com/) 提供完整的国际化解决方案，采用**命名空间分离** + **静态全量预加载**的架构。所有翻译资源在构建时通过 Vite `import.meta.glob` 收集，运行时一次性注入 i18next，无需异步加载。

### 核心特性

- **命名空间分离**：`_core/` 存放全局翻译（common/auth/routes/editor），`_modules/` 按功能模块拆分（共 27 个业务模块）
- **静态预加载**：所有 JSON 通过 `import.meta.glob({ eager: true })` 同步加载，无运行时请求
- **双语支持**：`zh-CN`（简体中文）和 `en-US`（英语）
- **语言检测**：localStorage → 浏览器语言 → 默认语言三级回退
- **偏好设置同步**：语言变更自动同步 i18n ↔ PreferencesStore ↔ 浏览器环境
- **格式化工具**：基于 `Intl` API 的日期/数字/货币/百分比/相对时间格式化
- **错误消息 i18n**：后端错误码自动映射为本地化提示文本

---

## 模块结构

```
src/core/i18n/                       # i18n 核心模块
├── index.ts                         # 统一导出
├── config/
│   └── i18n.ts                      # initI18n() 初始化函数
├── hooks/
│   ├── useI18n.ts                   # 类型安全的 useTranslation 封装
│   └── useLocaleSync.ts             # 语言 ↔ PreferencesStore 双向同步
├── utils/
│   ├── detector.ts                  # 语言检测（浏览器/localStorage/默认）
│   ├── formatter.ts                 # Intl 格式化工具（日期/数字/货币/百分比/相对时间）
│   └── sync.ts                      # 语言同步到浏览器环境（HTML lang/dir/存储/标题）
└── components/
    └── LocaleSwitcher/index.tsx      # 语言切换下拉组件

src/locales/                         # 翻译资源
├── index.ts                         # 资源汇总 + 类型导出
├── zh-CN/
│   ├── index.ts                     # Vite glob 收集 → zhCN 对象
│   ├── _core/                       # 核心命名空间
│   │   ├── common.json              # 通用翻译（按钮/错误码/请求状态等）
│   │   ├── auth.json                # 认证相关（登录/注册/Token 过期等）
│   │   ├── routes.json              # 路由标题（侧边栏/标签页/面包屑）
│   │   └── editor.json              # 富文本编辑器
│   └── _modules/                    # 业务模块命名空间（27 个 JSON）
│       ├── dashboard.json
│       ├── user.json
│       ├── role.json
│       ├── permission.json
│       ├── menu.json
│       ├── ...（共 27 个）
│       └── preferences.json
└── en-US/                           # 与 zh-CN 完全对称
    ├── index.ts
    ├── _core/
    └── _modules/
```

### 导入方式

```typescript
// Hook（组件内推荐）
import { useI18n, useLocaleSync } from '@/core/i18n';

// 工具函数
import { formatDate, formatNumber, detectBrowserLocale } from '@/core/i18n';

// 组件
import { LocaleSwitcher } from '@/core/i18n';

// 类型
import type { SupportedLocale } from '@/locales';
```

---

## 快速开始

### 在组件中使用翻译

```tsx
import { useI18n } from '@/core/i18n';

const UserPage = () => {
  // 指定命名空间为 'user'，后续 t() 调用无需 'user:' 前缀
  const { t } = useI18n('user');

  return (
    <>
      <h1>{t('pageTitle')}</h1>
      <span>{t('username')}</span>
      <span>{t('createdAt', { time: '2026-01-01' })}</span>
    </>
  );
};
```

### 切换语言

```tsx
import { useLocaleSync } from '@/core/i18n';

const LanguageSwitcher = () => {
  const { changeLocale } = useLocaleSync();

  return (
    <button onClick={() => changeLocale('en-US')}>English</button>
    <button onClick={() => changeLocale('zh-CN')}>简体中文</button>
  );
};
```

> `changeLocale()` 会同时更新 PreferencesStore、i18next 实例、浏览器 `lang` 属性和 localStorage。

### 添加新的业务模块翻译

**第一步**：创建翻译文件

```
src/locales/zh-CN/_modules/my-module.json
src/locales/en-US/_modules/my-module.json
```

```json
// zh-CN/_modules/my-module.json
{
  "pageTitle": "我的模块",
  "name": "名称",
  "createSuccess": "创建成功"
}
```

**第二步**：在组件中使用

```tsx
const { t } = useI18n('my-module');
// t('pageTitle') → "我的模块"
```

> 无需额外注册，`_modules/*.json` 会通过 `import.meta.glob` 自动收集为新的命名空间。

---

## 命名空间规范

### 命名空间分类

| 类别 | 目录 | 命名空间 | 说明 |
|------|------|----------|------|
| 核心 | `_core/` | `common` | 通用 UI 文本（按钮、错误码、请求状态） |
| 核心 | `_core/` | `auth` | 登录/注册/Token 相关 |
| 核心 | `_core/` | `routes` | 路由标题（侧边栏、标签页、面包屑） |
| 核心 | `_core/` | `editor` | 富文本编辑器 |
| 业务模块 | `_modules/` | `dashboard`、`user`、`role` 等 | 按功能模块独立 |
| 业务模块 | `_modules/` | `preferences` | 偏好设置面板 |
| 业务模块 | `_modules/` | `menu` | 菜单管理界面 |

### 命名空间使用规则

**1. 必须在 `useTranslation` / `useI18n` 中指定命名空间**

```tsx
// ✅ 正确：指定命名空间后，t() 直接使用键名
const { t } = useI18n('user');
t('username');         // 查找 user 命名空间下的 'username'

// ❌ 错误：未指定命名空间，需使用 'namespace:key' 格式
const { t } = useI18n();
t('user:username');    // 虽然能工作，但不符合项目规范
```

**2. 路由标题使用 `routes` 命名空间**

路由的 `meta.title` 使用 `'routes:xxx'` 格式，由路由/菜单系统自动解析：

```typescript
// src/router/modules/system.tsx
{
  meta: {
    title: 'routes:system',  // → t('system', { ns: 'routes' })
  },
}
```

**3. 菜单管理 UI 文本使用 `menu` 命名空间**

菜单管理功能页面的文本使用 `menu` 命名空间，与路由标题的 `routes` 命名空间分离，避免冲突。

```tsx
const { t } = useI18n('menu');
t('pageTitle');  // 菜单管理页面标题
```

**4. 默认命名空间为 `common`**

```tsx
// 不传参数时默认使用 'common'
const { t } = useI18n();
t('button.ok');  // 查找 common 命名空间
```

### 命名空间与文件对应关系

| 命名空间 | 文件路径 |
|----------|----------|
| `common` | `_core/common.json` |
| `auth` | `_core/auth.json` |
| `routes` | `_core/routes.json` |
| `editor` | `_core/editor.json` |
| `user` | `_modules/user.json` |
| `role` | `_modules/role.json` |
| `preferences` | `_modules/preferences.json` |
| ... | `_modules/<name>.json` |

> 命名空间名 = JSON 文件名（不含 `.json` 后缀）。

---

## Hooks API

### `useI18n<N>(ns?)` — 类型安全翻译 Hook

```typescript
import { useI18n } from '@/core/i18n';

const { t, i18n, changeLocale } = useI18n('user');
```

| 返回值 | 类型 | 说明 |
|--------|------|------|
| `t` | `TFunction` | 翻译函数，在指定命名空间下查找 |
| `i18n` | `i18next.i18n` | i18next 实例 |
| `changeLocale` | `(locale: SupportedLocale) => Promise<void>` | 切换语言 |

**参数**：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `ns` | `string \| string[]` | `'common'` | 命名空间，指定后 `t()` 中无需命名空间前缀 |

**使用多个命名空间**：

```tsx
const { t } = useI18n(['user', 'common']);
t('username');       // 优先查找 user 命名空间
t('button.ok');      // user 中找不到时回退到 common
```

### `useLocaleSync()` — 语言双向同步 Hook

确保 i18next 实例和 PreferencesStore 中的语言设置保持一致：

```typescript
import { useLocaleSync } from '@/core/i18n';

const { changeLocale } = useLocaleSync();
```

| 返回值 | 类型 | 说明 |
|--------|------|------|
| `changeLocale` | `(locale: SupportedLocale) => void` | 切换语言（写入 PreferencesStore，自动触发 i18n 同步） |

**同步规则**：

1. **PreferencesStore → i18n**：当 `preferences.app.locale` 变化时，自动调用 `i18n.changeLanguage()`
2. **i18n → PreferencesStore**：i18n 初始化后，反向同步一次当前语言到 store
3. **浏览器环境同步**：自动更新 `document.documentElement.lang` 和 `localStorage`

---

## 组件

### `<LocaleSwitcher />` — 语言切换下拉

```tsx
import { LocaleSwitcher } from '@/core/i18n';

// 直接使用，自带当前语言状态和切换逻辑
<LocaleSwitcher />
```

渲染为 Ant Design 的 `<Select>` 组件，显示 `🇨🇳 简体中文` / `🇺🇸 English`。

---

## 工具函数

### 语言检测（detector.ts）

```typescript
import { detectBrowserLocale, getFinalLocale, isSupportedLocale } from '@/core/i18n';

// 从浏览器语言列表检测
detectBrowserLocale();    // → 'zh-CN' | 'en-US'

// 综合检测：localStorage('app_locale') > 浏览器语言 > 默认('zh-CN')
getFinalLocale();         // → 'zh-CN'

// 校验
isSupportedLocale('ja');  // → false
```

检测优先级：`localStorage` 中存储的语言 > 浏览器首选语言（支持前缀匹配 `zh` → `zh-CN`）> 默认 `zh-CN`。

### 语言同步（sync.ts）

```typescript
import { syncLocaleToBrowser, setupLocaleSync, syncPageTitle } from '@/core/i18n';

// 同步到浏览器环境（HTML lang 属性 + dir 属性 + localStorage）
syncLocaleToBrowser('en-US');

// 监听 i18n 语言变化事件，自动同步
const unsubscribe = setupLocaleSync((locale) => {
  console.log('语言已切换为:', locale);
});
unsubscribe(); // 取消监听

// 同步页面标题
syncPageTitle(t, 'routes:dashboard');
```

### 格式化工具（formatter.ts）

基于浏览器原生 `Intl` API，自动适配当前语言：

```typescript
import { formatDate, formatNumber, formatCurrency, formatPercent, formatRelativeTime } from '@/core/i18n';

// 日期格式化
formatDate(new Date(), 'zh-CN');                      // → "2026年5月28日"
formatDate(new Date(), 'en-US');                       // → "May 28, 2026"
formatDate(new Date(), 'zh-CN', { timeStyle: 'short' }); // 带时间

// 数字格式化（千分位）
formatNumber(1234567.89, 'zh-CN');                     // → "1,234,567.89"
formatNumber(1234567.89, 'en-US');                     // → "1,234,567.89"

// 货币格式化
formatCurrency(99.9, 'zh-CN', 'CNY');                 // → "¥99.90"
formatCurrency(99.9, 'en-US', 'USD');                 // → "$99.90"

// 百分比格式化
formatPercent(0.856, 'zh-CN');                         // → "86%"
formatPercent(0.856, 'en-US');                         // → "86%"

// 相对时间
formatRelativeTime(new Date(Date.now() - 3000), 'zh-CN');  // → "3秒钟前"
formatRelativeTime(new Date(Date.now() - 3600000), 'en-US'); // → "1 hour ago"
```

---

## 初始化流程

```
bootstrap.ts → _initI18n()
  → usePreferencesStore.getState().preferences.app.locale  // 读取偏好语言
  → initI18n(locale)                                        // 初始化 i18next
      → i18n.use(initReactI18next).init({
          lng: locale,
          resources,           // 所有 JSON 预加载
          ns: allNamespaces,   // 声明所有命名空间
          defaultNS: 'common',
          fallbackLng: 'zh-CN',
        })
  → RequestClient.init(...)                                 // 注入 getLocale 回调
```

**关键点**：
- `resources` 在构建时通过 Vite glob 静态收集，运行时同步注入
- `allNamespaces` 从 `zhCN` 对象的 keys 自动提取
- 开发环境下缺失 key 会输出 `console.warn`

---

## 错误消息国际化

`bootstrap.ts` 中的 `getErrorMsg()` 函数实现后端错误消息的本地化映射，优先级：

```
1. reason 字段 → i18n key "request.reason.{reason}"
2. reason 无翻译 → 使用后端返回的 message
3. 都没有 → HTTP code → i18n key "request.status.{code}"
4. 全部失败 → 兜底文本 "请求发生错误,请稍后重试"
```

```typescript
// 在 RequestClient 中自动调用
// 后端返回 { reason: 'USER_FREEZE', message: '...' }
// → 自动翻译为 "用户被冻结，请联系管理员"
```

---

## 翻译键约定

### 插值语法

使用 `{{var}}` 格式（双花括号），**不是** `#{var}`：

```json
{
  "welcome": "欢迎，{{name}}",
  "deleteConfirmDesc": "确定要删除这个{{moduleName}}吗？"
}
```

```tsx
t('welcome', { name: '张三' });          // → "欢迎，张三"
t('deleteConfirmDesc', { moduleName: '用户' }); // → "确定要删除这个用户吗？"
```

### 路由标题格式

路由 `meta.title` 使用 `'routes:xxx'` 格式，系统自动处理两种格式：

| 格式 | 示例 | 处理方式 |
|------|------|----------|
| 带前缀 | `'routes:dashboard'` | 提取 `routes` 作为命名空间，`dashboard` 作为键 |
| 不带前缀 | `'dashboard'` | 直接作为键，使用默认命名空间 |

### 月份文本

月份文本使用 `months.*` 键前缀：

```tsx
t('months.jan');  // → "一月" / "January"
```

---

## 典型场景

### 场景一：新页面使用多语言

```tsx
// 1. 确保翻译文件存在：src/locales/zh-CN/_modules/my-page.json
// 2. 在组件中使用
import { useI18n } from '@/core/i18n';

const MyPage = () => {
  const { t } = useI18n('my-page');

  return (
    <div>
      <h1>{t('pageTitle')}</h1>
      <ProColumns[{
        title: t('name'),          // 表格列标题
        dataIndex: 'name',
      }]} />
    </div>
  );
};
```

### 场景二：在非组件环境中获取翻译

```typescript
import i18n from 'i18next';

// 在工具函数、service 层等非 React 场景
const msg = i18n.t('request.error.networkError', { ns: 'common' });

// 获取当前语言
const currentLang = i18n.language;  // → 'zh-CN'
```

### 场景三：多命名空间组件

```tsx
// 一个页面需要引用多个模块的翻译
const { t } = useI18n(['user-detail', 'user']);

t('basicInfo');           // 优先查找 user-detail
t('createSuccess');       // user-detail 中没有则回退到 user
```

---

## 常见问题

### Q1: 新增翻译 JSON 后不生效？

翻译文件放在 `_modules/` 或 `_core/` 目录下即可自动被 `import.meta.glob` 收集。如果仍然不生效：
1. 确认 JSON 格式正确（无语法错误）
2. 确认文件名与 `useI18n()` 中指定的命名空间一致
3. 重启开发服务器（Vite glob 在启动时收集）

### Q2: 翻译键存在但页面显示键名而不是翻译内容？

检查是否指定了正确的命名空间：

```tsx
// ❌ 错误：命名空间不匹配
const { t } = useI18n('common');
t('pageTitle');  // common 中没有 'pageTitle'

// ✅ 正确：指定正确的命名空间
const { t } = useI18n('dashboard');
t('pageTitle');  // dashboard 中有 'pageTitle'
```

### Q3: 如何处理复数/性别等复杂语法？

i18next 原生支持，参见 [i18next 复数文档](https://www.i18next.com/translation-function/plurals)。当前项目使用简单的插值即可满足需求。

### Q4: 翻译文件中可以使用嵌套结构吗？

可以。JSON 支持嵌套对象，`t()` 中用 `.` 连接：

```json
{
  "button": {
    "create": "新建",
    "delete": "删除"
  }
}
```

```tsx
t('button.create');  // → "新建"
t('button.delete');  // → "删除"
```

### Q5: 如何获取 Ant Design 的 locale 对象？

```tsx
import { useLocale } from '@/core/preferences';

const { antdLocale } = useLocale();
// 传给 ConfigProvider，或直接用在 DatePicker 等组件
```

> i18n 模块本身不直接提供 AntD locale，由 preferences 模块统一管理。

---

## 注意事项

1. **必须指定命名空间**：使用 `useI18n('moduleName')` 指定命名空间，避免使用 `namespace:key` 前缀格式
2. **路由标题使用 `routes:` 前缀**：`meta.title` 中使用 `'routes:xxx'` 格式，由路由系统自动解析
3. **插值用 `{{var}}`**：不要使用 `#{var}` 或 `${var}`
4. **翻译文件必须双语对称**：`zh-CN` 和 `en-US` 的 `_modules/` 目录下的 JSON 文件名和键名必须一一对应
5. **获取当前语言用 `i18n.language`**：在 React 环境中使用 `i18n.language`，而非 `i18n.locale`
6. **硬编码文本必须提取**：所有用户可见的中文/英文文本必须放入翻译文件，通过 `t()` 渲染
7. **Vite glob 限制**：动态导入必须使用相对路径，不能使用 `@/` 别名

---

## 相关文件

- [i18n 初始化](./config/i18n.ts) — initI18n()
- [useI18n Hook](./hooks/useI18n.ts) — 类型安全翻译 Hook
- [useLocaleSync Hook](./hooks/useLocaleSync.ts) — 语言双向同步
- [语言检测工具](./utils/detector.ts) — detectBrowserLocale / getFinalLocale
- [格式化工具](./utils/formatter.ts) — formatDate / formatNumber / formatCurrency
- [语言同步工具](./utils/sync.ts) — syncLocaleToBrowser / setupLocaleSync
- [翻译资源入口](../../locales/index.ts) — resources 汇总 + 类型导出
- [中文翻译](../../locales/zh-CN/index.ts) — zh-CN glob 收集
- [英文翻译](../../locales/en-US/index.ts) — en-US glob 收集
- [引导入口](../../bootstrap.ts) — 应用启动时初始化 i18n
- [偏好设置模块](../preferences/) — 语言偏好管理与 AntD locale
