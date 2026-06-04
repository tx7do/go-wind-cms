# 偏好设置模块 (Preferences)

## 概述

本模块是 Taro 小程序/H5 CMS 前台的用户偏好设置管理方案，涵盖主题切换、内容偏好、语言切换、功能部件开关等配置项。基于 Zustand + persist + React Context 实现状态管理与 localStorage 持久化。

### 核心特性

- **持久化存储**：偏好设置自动保存到 `localStorage`（键名 `app-preferences`），刷新页面不丢失
- **深度合并更新**：支持传入部分配置进行深度合并，无需提供完整对象
- **主题系统**：支持亮色/暗色/自动三种模式，颜色使用 HSL raw 格式映射 CSS 变量
- **响应式主题**：通过 `useSyncExternalStore` 监听系统主题变化，`auto` 模式下实时跟随
- **Context 隔离**：使用 React Context + Store 工厂函数，每个客户端独立实例
- **Tailwind CSS 集成**：`useThemeConfig` 返回 Tailwind CSS 变量风格的主题配置，由 `ThemeClientProvider` 注入 CSS 变量

---

## 模块结构

```
src/core/preferences/
├── index.ts                         # 统一导出（types + config + hooks + store）
├── types/                           # 类型定义
│   ├── index.ts                     # 类型导出入口
│   ├── preferences-root.ts          # Preferences 总接口 + DeepPartial + InitialOptions
│   ├── app.ts                       # 6 个子模块偏好设置接口
│   ├── theme.ts                     # ThemeModeType 主题模式
│   └── layout.ts                    # SupportedLanguagesType + PageTransitionType + ContentCompactType
├── config/                          # 默认配置
│   ├── index.ts                     # 导出入口
│   └── default.ts                   # defaultPreferences 默认值
├── store/                           # 状态管理
│   └── index.ts                     # createPreferencesStore + PreferencesStoreContext + usePreferencesStore
├── hooks/                           # Hooks
│   ├── index.ts                     # 导出入口
│   ├── usePreferences.ts            # 综合偏好设置 Hook（推荐使用）
│   ├── useThemeConfig.ts            # Tailwind CSS 主题配置 Hook
│   └── useLocale.ts                 # 语言切换 Hook（usePreferencesLocale）
└── utils/
    └── merge.ts                     # mergeDeep 深度合并工具
```

> **注意**：主题渲染由 `src/components/layout/ThemeClientProvider.tsx` 负责（不在本模块内），它通过内联 CSS 变量注入色板，配合 `usePreferences()` 的 `isDark` 状态切换暗色/亮色。

### 导入方式

```typescript
// Hook（组件内推荐）
import { usePreferences, usePreferencesLocale, useThemeConfig } from '@/core/preferences';

// Store（需要精细控制订阅时，必须在 <PreferencesStoreProvider> 内使用）
import { usePreferencesStore, createPreferencesStore, PreferencesStoreContext } from '@/core/preferences';

// 类型
import type { Preferences, DeepPartial, ThemeModeType, AppPreferences } from '@/core/preferences';

// 默认配置
import { defaultPreferences } from '@/core/preferences';

// 深度合并工具
import { mergeDeep, isPlainObject } from '@/core/preferences';
```

---

## 快速开始

### 读取偏好设置

```tsx
import { usePreferences } from '@/core/preferences';

const MyComponent = () => {
  const { theme, app, content, widget, isDark, isMobile } = usePreferences();

  return (
    <div>
      <p>当前主题: {theme.mode}</p>                {/* 'dark' | 'light' | 'auto' */}
      <p>是否暗色: {isDark ? '是' : '否'}</p>      {/* boolean，auto 模式已解析 */}
      <p>语言: {app.locale}</p>                     {/* 'zh-CN' | 'en-US' */}
      <p>紧凑模式: {content.compactMode ? '开' : '关'}</p>
      <p>回到顶部: {widget.backToTop ? '显示' : '隐藏'}</p>
    </div>
  );
};
```

### 修改偏好设置

```tsx
import { usePreferences } from '@/core/preferences';

const ThemeSwitcher = () => {
  const { toggleTheme, setThemeMode, updateTheme, updatePreferences } = usePreferences();

  // 快捷方法：明暗切换
  const handleToggle = () => toggleTheme();

  // 指定模式
  const handleSetDark = () => setThemeMode('dark');

  // 更新部分主题配置
  const handleChangeColor = () => updateTheme({ colorPrimary: '142.1 76.2% 36.3%' });

  // 同时更新多个分组
  const handleBatch = () => updatePreferences({
    theme: { mode: 'light' },
    app: { compact: true },
    content: { compactMode: true },
  });

  return <button onClick={handleToggle}>切换主题</button>;
};
```

### 重置为默认值

```tsx
const { resetPreferences } = usePreferences();

// 恢复所有设置到初始值
resetPreferences();
```

---

## 偏好设置配置项

### 完整结构概览

| 分组 | 类型 | 说明 | 默认值要点 |
|------|------|------|-----------|
| `app` | `AppPreferences` | 全局应用配置 | `zh-CN`，不紧凑，每页 10 条 |
| `theme` | `ThemePreferences` | 主题配置 | `auto` 模式，绿色主题色 |
| `content` | `ContentPreferences` | 内容偏好（CMS 专用） | 隐藏敏感内容，显示推荐 |
| `copyright` | `CopyrightPreferences` | 版权配置 | 启用，GoWind |
| `widget` | `WidgetPreferences` | 功能部件开关 | 全部启用 |
| `transition` | `TransitionPreferences` | 页面动画配置 | `fade-slide`，显示进度条 |

### app — 全局配置

```typescript
interface AppPreferences {
  name: string;                              // 应用名称，默认 'GoWind CMS'
  title: string;                             // 浏览器标签标题，默认 'GoWind Content Hub'
  defaultAvatar: string;                     // 默认头像路径
  locale: SupportedLanguagesType;            // 'zh-CN' | 'en-US'，默认 'zh-CN'
  isMobile: boolean;                         // 移动端模式，默认 false
  compact: boolean;                          // 紧凑模式，默认 false
  defaultPageSize: number;                   // 默认每页条数，默认 10
}
```

### theme — 主题配置

```typescript
interface ThemePreferences {
  mode: ThemeModeType;                       // 'auto' | 'dark' | 'light'，默认 'auto'
  colorPrimary: string;                      // 主题色（HSL raw 格式），默认 '142.1 76.2% 36.3%'
  colorSuccess: string;                      // 成功色，默认 '142.1 76.2% 36.3%'
  colorWarning: string;                      // 警告色，默认 '38 92% 50%'
  colorDestructive: string;                  // 错误色，默认 '0 84.2% 60.2%'
  radius: string;                            // 圆角大小，默认 '0.6rem'
}
```

> **HSL raw 格式说明**：颜色值使用 `"H S% L%"` 格式（如 `"142.1 76.2% 36.3%"`），直接映射到 CSS `hsl()` 函数。不带 `hsl()` 包裹，由 `ThemeClientProvider` 在注入 CSS 变量时组合。

### content — 内容偏好

```typescript
interface ContentPreferences {
  hideSensitiveContent: boolean;             // 隐藏敏感内容，默认 true
  compactMode: boolean;                      // 紧凑列表/卡片，默认 false
  showRecommendations: boolean;              // 显示推荐内容，默认 true
}
```

### copyright — 版权配置

```typescript
interface CopyrightPreferences {
  enable: boolean;                           // 是否显示版权信息，默认 true
  companyName: string;                       // 公司名称，默认 'GoWind'
  companySiteLink: string;                   // 公司链接
  date: string;                              // 版权日期，默认 '2026'
  icp: string;                               // 备案号
  icpLink: string;                           // 备案号链接
}
```

### widget — 功能部件开关

```typescript
interface WidgetPreferences {
  themeToggle: boolean;                      // 主题切换部件，默认 true
  languageToggle: boolean;                   // 语言切换部件，默认 true
  globalSearch: boolean;                     // 全局搜索部件，默认 true
  backToTop: boolean;                        // 回到顶部部件，默认 true
}
```

### transition — 页面动画

```typescript
interface TransitionPreferences {
  enable: boolean;                           // 是否启用页面切换动画，默认 true
  loading: boolean;                          // 页面加载 loading，默认 true
  name: PageTransitionType | string;         // 动画名称，默认 'fade-slide'
  progress: boolean;                         // 页面加载进度条，默认 true
}
```

**PageTransitionType**：`'fade' | 'fade-down' | 'fade-slide' | 'fade-up'`

---

## Store 架构

本模块采用 **Zustand 工厂函数 + React Context** 模式，而非全局单例：

```typescript
// 1. 创建独立 store 实例（工厂函数）
const store = createPreferencesStore();

// 2. 通过 Context 注入
<PreferencesStoreContext.Provider value={store}>
  <App />
</PreferencesStoreContext.Provider>

// 3. 组件内通过 Hook 消费
function MyComponent() {
  const { preferences, setPreferences } = usePreferencesStore();
  // ...
}
```

### `createPreferencesStore()` — Store 工厂函数

创建独立的 Zustand store 实例，内置 `persist` 中间件：

```typescript
import { createPreferencesStore } from '@/core/preferences';

const store = createPreferencesStore();
// store 是标准的 Zustand StoreApi<PreferencesState>
```

### `usePreferencesStore()` — 消费 Hook

**必须在 `<PreferencesStoreContext.Provider>` 内部使用**，否则会抛出错误。

支持 selector 精确订阅：

```tsx
import { usePreferencesStore } from '@/core/preferences';

// 完整 state
const { preferences, setPreferences, resetPreferences, getPreference } = usePreferencesStore();

// selector 精确订阅 — 只有 theme.mode 变化才重渲染
const themeMode = usePreferencesStore((s) => s.preferences.theme.mode);

// selector 订阅 widget
const backToTop = usePreferencesStore((s) => s.preferences.widget.backToTop);
```

### PreferencesState 接口

```typescript
interface PreferencesState {
  preferences: Preferences;
  setPreferences: (overrides: DeepPartial<Preferences>) => void;  // 深度合并更新
  resetPreferences: () => void;                                    // 重置为默认值
  getPreference: <K extends keyof Preferences>(key: K) => Preferences[K];
}
```

---

## Hooks API

### `usePreferences()` — 综合偏好设置 Hook（推荐）

在 React 组件中使用的首选方式，提供分组访问和便捷更新方法：

```typescript
const {
  // 完整偏好设置
  preferences,   // Preferences

  // 分组访问
  app,           // AppPreferences
  theme,         // ThemePreferences
  content,       // ContentPreferences
  copyright,     // CopyrightPreferences
  widget,        // WidgetPreferences
  transition,    // TransitionPreferences

  // 计算属性
  isDark,        // boolean — 当前是否暗色模式（auto 模式已解析系统偏好）
  isMobile,      // boolean — 当前是否移动端

  // 通用更新
  updatePreferences,  // (overrides: DeepPartial<Preferences>) => void
  resetPreferences,   // () => void
  getPreference,      // <K>(key: K) => Preferences[K]

  // 分组便捷更新
  updateApp,     // (overrides: DeepPartial<AppPreferences>) => void
  updateTheme,   // (overrides: DeepPartial<ThemePreferences>) => void
  updateContent, // (overrides: DeepPartial<ContentPreferences>) => void
  updateWidget,  // (overrides: DeepPartial<WidgetPreferences>) => void

  // 快捷方法
  toggleTheme,   // () => void — 明暗切换
  setThemeMode,  // (mode: ThemeModeType) => void
  setLanguage,   // (locale: SupportedLanguagesType) => void
} = usePreferences();
```

### `useThemeConfig()` — Tailwind CSS 主题配置

返回 Tailwind CSS 变量风格的主题配置信息（无 AntD 依赖）：

```typescript
interface ThemeConfig {
  isDark: boolean;          // 是否暗色模式
  isCompact: boolean;       // 是否紧凑模式
  colorPrimary: string;     // 主色调（HSL raw）
  radius: string;           // 圆角
}
```

```tsx
import { useThemeConfig } from '@/core/preferences';

const { isDark, colorPrimary, radius } = useThemeConfig();
```

> 主题渲染由 `src/components/layout/ThemeClientProvider.tsx` 负责，它读取 `usePreferences()` 的 `isDark` 注入对应色板的 CSS 变量。

### `usePreferencesLocale()` — 语言切换

提供语言相关的快捷方法，整合了 next-intl URL locale 和 preferences 持久化 locale：

```tsx
import { usePreferencesLocale } from '@/core/preferences';

const {
  locale,           // SupportedLanguagesType — 当前语言（优先 next-intl URL locale）
  isZhCN,           // boolean
  isEnUS,           // boolean
  setLocale,        // (locale: SupportedLanguagesType) => void
  setZhCN,          // () => void
  setEnUS,          // () => void
  toggleLocale,     // () => void — 中英切换
  supportedLocales, // [{ value, label }]
} = usePreferencesLocale();
```

> **命名说明**：导出名是 `usePreferencesLocale` 而非 `useLocale`，避免与 next-intl 的 `useLocale` 冲突。

---

## 持久化机制

偏好设置通过 Zustand 的 `persist` 中间件自动保存到 `localStorage`：

- **存储键名**：`app-preferences`
- **存储内容**：`preferences` 对象（`partialize` 只持久化偏好设置数据，不持久化方法）
- **合并策略**：运行时通过 `setPreferences()` 更新时使用 `mergeDeep` 进行深度合并

```typescript
// 查看当前存储的偏好设置
localStorage.getItem('app-preferences');

// 清除偏好设置（恢复默认）
localStorage.removeItem('app-preferences');
```

---

## 更新机制

### `setPreferences()` — 深度合并更新

支持传入部分配置，自动与现有值深度合并：

```typescript
// 只更新 theme.mode，其他 theme 字段保持不变
setPreferences({ theme: { mode: 'dark' } });

// 同时更新多个分组
setPreferences({
  app: { compact: true },
  content: { compactMode: true },
  widget: { backToTop: false },
});
```

### `mergeDeep()` — 深度合并规则

| 规则 | 说明 |
|------|------|
| 普通对象 | 递归合并 |
| 基本类型 | 直接覆盖 |
| `undefined` | 跳过（不覆盖） |
| `null` | 跳过（不覆盖） |
| React 元素 / DOM 节点 / 函数 | 跳过并警告 |

---

## 典型场景

### 场景一：根据主题模式切换样式

```tsx
import { usePreferences } from '@/core/preferences';

const MyComponent = () => {
  const { isDark, theme } = usePreferences();

  // isDark 已自动处理 auto 模式
  return (
    <div style={{ background: isDark ? '#1a1a1a' : '#ffffff' }}>
      <p>当前模式: {theme.mode}</p>
    </div>
  );
};
```

### 场景二：精确订阅避免不必要渲染

```tsx
import { usePreferencesStore } from '@/core/preferences';

// 只订阅 theme.mode — 只有该字段变化时才重渲染
const ThemeIndicator = () => {
  const themeMode = usePreferencesStore((s) => s.preferences.theme.mode);
  return <span>{themeMode}</span>;
};
```

### 场景三：语言切换

```tsx
import { usePreferencesLocale } from '@/core/preferences';

const LanguageSwitcher = () => {
  const { locale, toggleLocale, isZhCN } = usePreferencesLocale();

  return (
    <button onClick={toggleLocale}>
      {isZhCN ? 'English' : '中文'}
    </button>
  );
};
```

### 场景四：使用 Tailwind 主题配置

```tsx
import { useThemeConfig } from '@/core/preferences';

const StyledCard = () => {
  const { isDark, colorPrimary, radius } = useThemeConfig();

  return (
    <div style={{
      borderRadius: radius,
      borderColor: `hsl(${colorPrimary})`,
    }}>
      卡片内容
    </div>
  );
};
```

### 场景五：控制功能部件显示

```tsx
import { usePreferences } from '@/core/preferences';

const Layout = () => {
  const { widget } = usePreferences();

  return (
    <>
      {widget.backToTop && <BackToTopButton />}
      {widget.themeToggle && <ThemeToggleButton />}
      {widget.globalSearch && <SearchButton />}
    </>
  );
};
```

---

## 常见问题

### Q1: 修改偏好设置后页面没有更新？

确认使用的是 `usePreferences()` 或 `usePreferencesStore()` Hook（React 组件内），而非直接操作 `localStorage`。Zustand 的响应式机制依赖 Hook 订阅。

### Q2: `usePreferencesStore` 报错 "must be used within PreferencesStoreProvider"？

`usePreferencesStore` 通过 React Context 获取 store 实例，必须在应用根节点用 `PreferencesStoreContext.Provider` 包裹并传入 `createPreferencesStore()` 创建的 store。

### Q3: 刷新后偏好设置丢失？

偏好设置通过 `localStorage` 持久化（键名 `app-preferences`）。检查：
1. 是否清除了浏览器缓存
2. 是否在隐私/无痕模式下使用
3. localStorage 是否被其他代码清除

### Q4: 如何添加新的偏好设置字段？

1. 在 `types/app.ts` 中对应的接口添加字段
2. 在 `config/default.ts` 的 `defaultPreferences` 中添加默认值
3. TypeScript 会自动提示所有使用 `setPreferences` 的地方需要更新
4. 如果需要新的分组，在 `types/preferences-root.ts` 的 `Preferences` 接口中添加

### Q5: `auto` 主题模式如何工作？

当 `theme.mode` 设为 `'auto'` 时：
- `usePreferences()` 通过 `useSyncExternalStore` 监听系统 `prefers-color-scheme: dark` 媒体查询
- H5 端：使用 `window.matchMedia` 监听变化
- 小程序端：使用 `Taro.onThemeChange` 监听变化
- `isDark` 计算属性已自动处理 `auto` 模式
- `ThemeClientProvider` 读取 `isDark` 注入对应色板的 CSS 变量

---

## 注意事项

1. **Context 依赖**：所有 hooks 必须在 `PreferencesStoreContext.Provider` 内使用
2. **使用 `DeepPartial` 更新**：`setPreferences` 接受 `DeepPartial<Preferences>`，只需传入要修改的字段
3. **不要直接修改 store**：始终通过 `setPreferences()` 更新，确保 `mergeDeep` 正确执行和持久化触发
4. **精确订阅**：在性能敏感组件中使用 `usePreferencesStore((s) => s.preferences.xxx)` 选择器，避免整个 preferences 对象变化时重渲染
5. **主题色格式**：颜色值使用 HSL raw 格式（如 `"142.1 76.2% 36.3%"`），由 CSS `hsl()` 函数组合
6. **主题渲染分离**：本模块只管理偏好数据，CSS 变量注入由 `src/components/layout/ThemeClientProvider.tsx` 负责

---

## 相关文件

- [Store](./store/index.ts) — createPreferencesStore + PreferencesStoreContext + usePreferencesStore
- [默认配置](./config/default.ts) — defaultPreferences 默认值
- [类型定义](./types/) — 所有偏好设置接口
- [深度合并工具](./utils/merge.ts) — mergeDeep + isPlainObject
- [ThemeClientProvider](../../components/layout/ThemeClientProvider.tsx) — 主题 CSS 变量注入（外部模块）
