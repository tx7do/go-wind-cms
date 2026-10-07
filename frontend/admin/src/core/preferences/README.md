# 偏好设置模块 (Preferences)

## 概述

本模块提供了完整的用户偏好设置管理方案，涵盖主题切换、布局控制、国际化、快捷键、功能部件等配置项。基于 Zustand + persist 实现状态管理与 localStorage 持久化。

### 核心特性

- **持久化存储**：偏好设置自动保存到 `localStorage`（键名 `app-preferences`），刷新页面不丢失
- **深度合并更新**：支持传入部分配置进行深度合并，无需提供完整对象
- **主题系统**：支持亮色/暗色/自动三种模式 + 15 种内置主题色 + 紧凑模式
- **响应式主题**：`ThemeProvider` 组件自动监听系统主题变化，`auto` 模式下实时跟随
- **AntD 集成**：自动将偏好设置转换为 Ant Design 的 `ThemeConfig` 和 `locale`
- **设置面板**：内置 `PreferencesPanel` Drawer 组件，包含外观/布局/快捷键/通用四个配置页

---

## 模块结构

```
src/core/preferences/
├── index.ts                         # 统一导出
├── types/                           # 类型定义
│   ├── index.ts                     # 类型导出入口
│   ├── preferences-root.ts          # Preferences 总接口 + DeepPartial
│   ├── app.ts                       # 各子模块偏好设置接口（14 个）
│   ├── theme.ts                     # 主题模式 + 内置主题类型
│   └── layout.ts                    # 布局相关枚举类型
├── config/                          # 默认配置
│   ├── default.ts                   # defaultPreferences 默认值
│   └── constants.ts                 # BUILT_IN_THEME_PRESETS 主题预设色板
├── store/                           # 状态管理
│   └── index.ts                     # usePreferencesStore (Zustand + persist)
├── hooks/                           # Hooks
│   ├── usePreferences.ts            # 综合偏好设置 Hook（推荐使用）
│   ├── useThemeConfig.ts            # AntD ThemeConfig 生成 Hook
│   └── useLocale.ts                 # 语言切换 Hook
├── components/                      # UI 组件
│   ├── ThemeProvider/index.tsx      # 全局主题 Provider（ConfigProvider + 水印 + 滤镜）
│   └── PreferencesPanel/            # 偏好设置 Drawer 面板
│       ├── index.tsx                # 面板容器（Drawer + Tab 切换）
│       ├── AppearancePanel.tsx      # 外观设置页（主题色/暗色/圆角等）
│       ├── LayoutPanel.tsx          # 布局设置页（侧边栏/标签页/面包屑等）
│       ├── ShortcutKeyPanel.tsx     # 快捷键设置页
│       └── GeneralPanel.tsx         # 通用设置页（水印/色弱/灰色模式等）
└── utils/
    └── merge.ts                     # mergeDeep 深度合并工具
```

### 导入方式

```typescript
// Hook（组件内推荐）
import { usePreferences, useLocale, useThemeConfig } from '@/core/preferences';

// Store（需要精细控制订阅时）
import { usePreferencesStore } from '@/core/preferences';

// 类型
import type { Preferences, DeepPartial, ThemeModeType } from '@/core/preferences';

// 组件
import { ThemeProvider, PreferencesPanel } from '@/core/preferences';
```

---

## 快速开始

### 读取偏好设置

```tsx
import { usePreferences } from '@/core/preferences';

const MyComponent = () => {
  const { theme, app, sidebar, tabbar } = usePreferences();

  return (
    <div>
      <p>当前主题: {theme.mode}</p>          {/* 'dark' | 'light' | 'auto' */}
      <p>布局方式: {app.layout}</p>           {/* 'sidebar-nav' | 'header-nav' | ... */}
      <p>侧边栏折叠: {sidebar.collapsed}</p>  {/* boolean */}
      <p>标签页风格: {tabbar.styleType}</p>    {/* 'chrome' | 'card' | ... */}
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
  const handleChangeColor = () => updateTheme({ colorPrimary: 'hsl(245 82% 67%)' });

  // 同时更新多个分组
  const handleBatch = () => updatePreferences({
    theme: { mode: 'light' },
    sidebar: { collapsed: true },
    app: { compact: true },
  });

  return <Button onClick={handleToggle}>切换主题</Button>;
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
| `app` | `AppPreferences` | 全局应用配置 | `sidebar-nav` 布局，`zh-CN`，前端权限模式 |
| `theme` | `ThemePreferences` | 主题配置 | `dark` 模式，`default` 主题色 |
| `sidebar` | `SidebarPreferences` | 侧边栏配置 | 宽度 224，不折叠 |
| `header` | `HeaderPreferences` | 顶栏配置 | 固定模式 |
| `tabbar` | `TabbarPreferences` | 标签页配置 | `chrome` 风格，启用拖拽/图标/更多 |
| `breadcrumb` | `BreadcrumbPreferences` | 面包屑配置 | 启用，显示首页图标 |
| `navigation` | `NavigationPreferences` | 导航菜单配置 | 手风琴模式，圆润风格 |
| `footer` | `FooterPreferences` | 底栏配置 | 默认隐藏 |
| `logo` | `LogoPreferences` | Logo 配置 | 启用，使用 `/logo.png` |
| `copyright` | `CopyrightPreferences` | 版权配置 | 启用 |
| `transition` | `TransitionPreferences` | 页面动画配置 | `fade-slide`，显示进度条 |
| `shortcutKeys` | `ShortcutKeyPreferences` | 快捷键配置 | 全部启用 |
| `widget` | `WidgetPreferences` | 功能部件配置 | 全部启用 |
| `widget` | `WidgetPreferences` | 功能部件开关 | 全部启用 |

### app — 全局配置

```typescript
interface AppPreferences {
  accessMode: 'frontend' | 'backend';        // 权限模式，默认 'frontend'
  authPageLayout: 'panel-center' | 'panel-left' | 'panel-right'; // 登录页布局
  locale: 'zh-CN' | 'en-US';                 // 界面语言
  layout: 'sidebar-nav' | 'header-nav' | 'mixed-nav' | 'sidebar-mixed-nav' | 'full-content';
  compact: boolean;                           // AntD 紧凑模式
  contentCompact: 'compact' | 'wide';         // 内容区域宽度
  colorGrayMode: boolean;                     // 灰色模式
  colorWeakMode: boolean;                     // 色弱模式
  watermark: boolean;                         // 水印
  enableRefreshToken: boolean;                // 启用 Token 刷新
  enableTenant: boolean;                      // 启用多租户
  enablePreferences: boolean;                 // 显示偏好设置入口
  isMobile: boolean;                          // 移动端模式
  // ...更多字段见 types/app.ts
}
```

### theme — 主题配置

```typescript
interface ThemePreferences {
  mode: 'auto' | 'dark' | 'light';           // 主题模式
  builtinType: BuiltinThemeType;              // 内置主题名
  colorPrimary: string;                       // 主题色（HSL 格式）
  colorSuccess: string;                       // 成功色
  colorWarning: string;                       // 警告色
  colorDestructive: string;                   // 错误色
  radius: string;                             // 圆角大小
  semiDarkHeader: boolean;                    // 半深色顶栏（仅 light 模式）
  semiDarkSidebar: boolean;                   // 半深色侧边栏（仅 light 模式）
}
```

**内置主题色**（`builtinType`）：

| 主题名 | 色值 | 主题名 | 色值 |
|--------|------|--------|------|
| `default` | `hsl(212 100% 45%)` 蓝色 | `green` | `hsl(161 90% 43%)` 绿色 |
| `violet` | `hsl(245 82% 67%)` 紫色 | `deep-green` | `hsl(181 84% 32%)` 深绿 |
| `pink` | `hsl(347 77% 60%)` 粉色 | `deep-blue` | `hsl(211 91% 39%)` 深蓝 |
| `yellow` | `hsl(42 84% 61%)` 黄色 | `orange` | `hsl(18 89% 40%)` 橙色 |
| `sky-blue` | `hsl(231 98% 65%)` 天蓝 | `rose` | `hsl(0 75% 42%)` 玫红 |
| `zinc` | `hsl(240 5% 26%)` 锌色 | `neutral` | `hsl(0 0% 25%)` 中性 |
| `slate` | `hsl(215 25% 27%)` 石板灰 | `gray` | `hsl(217 19% 27%)` 灰色 |
| `custom` | 自定义 | | |

### sidebar — 侧边栏配置

```typescript
interface SidebarPreferences {
  collapsed: boolean;              // 是否折叠
  collapsedShowTitle: boolean;     // 折叠时是否显示标题
  enable: boolean;                 // 是否可见
  expandOnHover: boolean;          // 鼠标悬停自动展开
  extraCollapse: boolean;          // 扩展区域折叠
  hidden: boolean;                 // CSS 级隐藏
  width: number;                   // 宽度（px），默认 224
}
```

### tabbar — 标签页配置

```typescript
interface TabbarPreferences {
  enable: boolean;                 // 是否启用标签页
  styleType: 'brisk' | 'card' | 'chrome' | 'plain';  // 标签页风格
  draggable: boolean;              // 是否可拖拽排序
  showIcon: boolean;               // 显示标签页图标
  showMaximize: boolean;           // 显示最大化按钮
  showMore: boolean;               // 显示更多按钮
  persist: boolean;                // 持久化标签页状态
  keepAlive: boolean;              // 标签页缓存
  height: number;                  // 标签页高度，默认 38
}
```

---

## Hooks API

### `usePreferences()` — 综合偏好设置 Hook（推荐）

在 React 组件中使用的首选方式，提供分组访问和便捷更新方法：

```typescript
const {
  // 分组访问
  preferences,   // Preferences — 完整偏好设置对象
  app,           // AppPreferences
  theme,         // ThemePreferences
  sidebar,       // SidebarPreferences
  tabbar,        // TabbarPreferences
  breadcrumb,    // BreadcrumbPreferences
  header,        // HeaderPreferences
  footer,        // FooterPreferences
  logo,          // LogoPreferences
  navigation,    // NavigationPreferences
  copyright,     // CopyrightPreferences
  transition,    // TransitionPreferences
  shortcutKeys,  // ShortcutKeyPreferences
  widget,        // WidgetPreferences

  // 计算属性
  isDark,        // boolean — 当前是否暗色模式（含 auto 解析）
  isMobile,      // boolean — 当前是否移动端

  // 通用更新
  updatePreferences,  // (overrides: DeepPartial<Preferences>) => void
  resetPreferences,   // () => void
  getPreference,      // <K>(key: K) => Preferences[K]

  // 便捷更新
  updateTheme,   // (overrides: DeepPartial<ThemePreferences>) => void
  updateApp,     // (overrides: DeepPartial<AppPreferences>) => void
  toggleTheme,   // () => void — 明暗切换
  setThemeMode,  // (mode) => void — 设置主题模式
  setLanguage,   // (locale) => void — 设置语言
} = usePreferences();
```

### `usePreferencesStore` — Zustand Store（精细订阅）

适用于需要精确控制重渲染的场景，使用 selector 只订阅需要的字段：

```tsx
import { usePreferencesStore } from '@/core/preferences';

// 只订阅 theme.mode（mode 变化才重渲染）
const themeMode = usePreferencesStore((s) => s.preferences.theme.mode);

// 只订阅 sidebar.collapsed
const collapsed = usePreferencesStore((s) => s.preferences.sidebar.collapsed);

// 更新偏好设置
const setPreferences = usePreferencesStore((s) => s.setPreferences);
setPreferences({ theme: { mode: 'dark' } });

// 非组件环境中使用
const locale = usePreferencesStore.getState().preferences.app.locale;
```

### `useThemeConfig()` — AntD 主题配置

生成 Ant Design 的 `ThemeConfig`，已自动处理暗色算法和紧凑算法：

```tsx
import { useThemeConfig } from '@/core/preferences';
import { ConfigProvider } from 'antd';

const themeConfig = useThemeConfig();
// → { algorithm: [darkAlgorithm, compactAlgorithm], token: { colorPrimary, ... } }

<ConfigProvider theme={themeConfig}>
  <App />
</ConfigProvider>
```

> 通常不需要手动使用此 Hook，`ThemeProvider` 组件已内部集成。

### `useLocale()` — 语言切换

提供语言相关的快捷方法和 Ant Design locale 对象：

```tsx
import { useLocale } from '@/core/preferences';

const {
  locale,          // 'zh-CN' | 'en-US'
  localeName,      // 显示名称（如"简体中文"）
  antdLocale,      // Ant Design locale 对象
  isZhCN,          // boolean
  isEnUS,          // boolean
  setLocale,       // (locale) => void
  toggleLocale,    // () => void — 中英切换
  supportedLocales,// [{ value, label }]
} = useLocale();
```

---

## 组件

### `<ThemeProvider>` — 全局主题 Provider

必须在应用最外层使用，负责：

1. **AntD 主题切换**：根据 `theme.mode` 自动应用 `darkAlgorithm` / `defaultAlgorithm`
2. **系统主题同步**：`mode='auto'` 时监听 `prefers-color-scheme` 变化
3. **CSS 滤镜**：色弱模式（`invert`）和灰色模式（`grayscale`）
4. **水印**：根据 `app.watermark` 显示/隐藏水印
5. **国际化**：根据 `app.locale` 注入 AntD locale
6. **根元素样式**：同步背景色和文字颜色，防止暗色模式白闪

```tsx
import { ThemeProvider } from '@/core/preferences';

// 在 App.tsx 中包裹
<ThemeProvider>
  <AppRouter />
</ThemeProvider>
```

### `<PreferencesPanel>` — 偏好设置面板

内置的设置 Drawer，包含四个 Tab 页：

| Tab | 面板 | 配置项 |
|-----|------|--------|
| 外观 | `AppearancePanel` | 主题色、明暗模式、圆角、半深色侧边栏/顶栏 |
| 布局 | `LayoutPanel` | 布局方式、侧边栏、标签页风格、面包屑、导航风格 |
| 快捷键 | `ShortcutKeyPanel` | 全局搜索/偏好设置/锁屏/注销快捷键开关 |
| 通用 | `GeneralPanel` | 水印、色弱模式、灰色模式、默认头像、检查更新 |

```tsx
import { PreferencesPanel } from '@/core/preferences';

const [open, setOpen] = useState(false);

<PreferencesPanel open={open} onClose={() => setOpen(false)} />
```

---

## 持久化机制

偏好设置通过 Zustand 的 `persist` 中间件自动保存到 `localStorage`：

- **存储键名**：`app-preferences`
- **存储内容**：`preferences` 对象（`partialize` 只持久化偏好设置数据，不持久化方法）
- **合并策略**：启动时 `persist` 中间件自动将 localStorage 中的值与 `defaultPreferences` 浅合并
- **深度合并**：运行时通过 `setPreferences()` 更新时使用 `mergeDeep` 进行深度合并

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
  sidebar: { collapsed: true, width: 200 },
  tabbar: { styleType: 'card' },
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

### 场景一：根据主题模式切换逻辑

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

### 场景二：非组件环境读取偏好设置

```typescript
import { usePreferencesStore } from '@/core/preferences';

// 在工具函数、service 层等非组件环境中
const locale = usePreferencesStore.getState().preferences.app.locale;
const isDark = usePreferencesStore.getState().preferences.theme.mode === 'dark';
```

### 场景三：精确订阅避免不必要渲染

```tsx
import { usePreferencesStore } from '@/core/preferences';

// 只订阅 sidebar.collapsed — 只有该字段变化时才重渲染
const Sidebar = () => {
  const collapsed = usePreferencesStore((s) => s.preferences.sidebar.collapsed);
  return <aside>{collapsed ? '已折叠' : '展开'}</aside>;
};
```

### 场景四：在 MainLayout 中使用各分组配置

```tsx
import { usePreferences } from '@/core/preferences';

const MainLayout = () => {
  const { sidebar, header, tabbar, breadcrumb, widget } = usePreferences();

  return (
    <Layout>
      {header.enable && <Header mode={header.mode} />}
      {sidebar.enable && <Sider collapsed={sidebar.collapsed} width={sidebar.width} />}
      {tabbar.enable && <TabsBar styleType={tabbar.styleType} />}
      {breadcrumb.enable && <Breadcrumbs />}
      {widget.refresh && <RefreshButton />}
    </Layout>
  );
};
```

---

## 常见问题

### Q1: 修改偏好设置后页面没有更新？

确认使用的是 `usePreferences()` 或 `usePreferencesStore()` Hook（React 组件内），而非直接操作 `localStorage`。Zustand 的响应式机制依赖 Hook 订阅。

### Q2: 如何在非组件环境（service、工具函数）中读取偏好设置？

使用 Zustand 的 `getState()` 静态方法：

```typescript
import { usePreferencesStore } from '@/core/preferences';

const themeMode = usePreferencesStore.getState().preferences.theme.mode;
```

### Q3: 刷新后偏好设置丢失？

偏好设置通过 `localStorage` 持久化（键名 `app-preferences`）。检查：
1. 是否清除了浏览器缓存
2. 是否在隐私/无痕模式下使用
3. localStorage 是否被其他代码清除

### Q4: 如何添加新的偏好设置字段？

1. 在 `types/app.ts` 中对应的接口添加字段
2. 在 `config/default.ts` 的 `defaultPreferences` 中添加默认值
3. TypeScript 会自动提示所有使用 `setPreferences` 的地方需要更新

### Q5: `auto` 主题模式如何工作？

当 `theme.mode` 设为 `'auto'` 时：
- `ThemeProvider` 会监听系统 `prefers-color-scheme: dark` 媒体查询
- 系统切换明暗时自动生效
- `usePreferences()` 的 `isDark` 计算属性已自动处理 `auto` 模式

---

## 注意事项

1. **使用 `DeepPartial` 更新**：`setPreferences` 接受 `DeepPartial<Preferences>`，只需传入要修改的字段
2. **不要直接修改 store**：始终通过 `setPreferences()` 更新，确保 `mergeDeep` 正确执行和持久化触发
3. **精确订阅**：在性能敏感组件中使用 `usePreferencesStore((s) => s.preferences.xxx)` 选择器，避免整个 preferences 对象变化时重渲染
4. **`ThemeProvider` 位置**：必须在最外层（`App.tsx`），在 `ConfigProvider` 之上的层级
5. **主题色格式**：颜色值使用 HSL 格式（如 `hsl(212 100% 45%)`），与 CSS 变量兼容
6. **偏好设置面板**：通过 `app.enablePreferences` 控制是否显示设置入口，设为 `false` 可隐藏

---

## 相关文件

- [Store](./store/index.ts) — usePreferencesStore 状态管理
- [默认配置](./config/default.ts) — defaultPreferences 默认值
- [类型定义](./types/) — 所有偏好设置接口
- [ThemeProvider](./components/ThemeProvider/index.tsx) — 全局主题 Provider
- [PreferencesPanel](./components/PreferencesPanel/index.tsx) — 设置面板
- [深度合并工具](./utils/merge.ts) — mergeDeep
- [i18n 语言同步](../i18n/hooks/useLocaleSync.ts) — 语言偏好与 i18n 同步
- [引导入口](../../bootstrap.ts) — 初始化时读取 locale
