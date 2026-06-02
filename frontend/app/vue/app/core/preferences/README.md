# core/preferences 偏好设置模块

> 面向 admin 业务逻辑开发者的使用指南

## 模块定位

提供 **主题、布局、侧边栏、标签页、面包屑、快捷键、动画** 等全应用级别的偏好配置管理。

核心能力：响应式状态 + localStorage 持久化 + CSS 变量联动 + 偏好设置面板。

业务层通过 `updatePreferences()` 修改值，框架自动处理存储、CSS 更新、语言切换等副作用。

---

## 目录结构

```
core/preferences/
├── index.ts                      # 统一导出（preferences, updatePreferences, resetPreferences 等）
├── preferences.ts                # PreferenceManager 核心类（状态管理、持久化、副作用）
├── use-preferences.ts            # usePreferences() composable（计算属性快捷访问）
├── update-css-variables.ts       # CSS 变量联动（主题色、暗色模式、圆角等）
├── config/
│   ├── default.ts                # 所有偏好项的默认值
│   └── constants.ts              # 内置主题预设色板
├── types/
│   ├── preferences-root.ts       # Preferences 总接口 + DeepPartial + InitialOptions
│   ├── app.ts                    # 各子项类型（AppPreferences, ThemePreferences 等）
│   ├── layout.ts                 # 布局相关枚举类型（LayoutType, ThemeModeType 等）
│   └── theme.ts                  # BuiltinThemeType 内置主题名
└── components/
    └── PreferencesPanel/         # 偏好设置面板 UI 组件
```

---

## 快速开始

### 读取偏好值

```ts
import { preferences } from "@/core/preferences";

// 直接读取（响应式 readonly 对象）
const layout = preferences.app.layout;         // "sidebar-nav"
const isDark = preferences.theme.mode;          // "dark"
const locale = preferences.app.locale;          // "zh-CN"
```

### 修改偏好值

```ts
import { updatePreferences } from "@/core/preferences";

// 支持深层部分更新（无需传完整对象）
updatePreferences({
  theme: { mode: "dark" },
});
```

### 在组件中使用 composable

```vue
<script setup lang="ts">
import { usePreferences } from "@/core/preferences";

const { isDark, layout, isMobile, locale, setTheme, toggleTheme } = usePreferences();
</script>

<template>
  <div>当前布局: {{ layout }}</div>
  <button @click="toggleTheme">切换主题</button>
</template>
```

---

## 核心 API

### `updatePreferences(updates: DeepPartial<Preferences>)`

更新偏好设置。支持深层部分更新，只传需要修改的字段。

```ts
// 修改主题色
updatePreferences({ theme: { colorPrimary: "#ff6600" } });

// 同时修改多个子项
updatePreferences({
  app: { locale: "en-US" },
  sidebar: { collapsed: true },
  tabbar: { enable: false },
});
```

**自动触发的副作用：**
| 修改项 | 副作用 |
|---|---|
| `theme.*` | 更新 CSS 变量、切换 `dark` 类名 |
| `app.locale` | 加载对应语言包、同步 vue-i18n |
| `app.colorGrayMode` | 切换 `grayscale-mode` 类名 |
| `app.colorWeakMode` | 切换 `invert-mode` 类名 |

### `initPreferences(options: InitialOptions)`

初始化偏好设置。在 `bootstrap.ts` 中调用一次，**不要在业务代码中调用**。

```ts
await initPreferences({
  namespace: "gowind",           // localStorage 键前缀
  overrides: {                   // 覆盖默认值
    app: { name: "My App" },
  },
});
```

### `resetPreferences()`

重置为初始值（默认值 + overrides 合并结果），同时清除 localStorage。

### `clearPreferencesCache()`

清除 localStorage 中的偏好缓存，不重置内存中的值。

---

## `usePreferences()` composable

提供一系列计算属性，方便在组件中快捷访问常用偏好值。

### 布局判断

| 返回值 | 类型 | 说明 |
|---|---|---|
| `layout` | `ComputedRef<string>` | 当前布局名（移动端强制为 `sidebar-nav`） |
| `isMobile` | `ComputedRef<boolean>` | 是否移动端（由 Tailwind `md` 断点自动判断） |
| `isDark` | `ComputedRef<boolean>` | 是否暗色模式 |
| `isSideNav` | `ComputedRef<boolean>` | 是否侧边菜单布局 |
| `isSideMixedNav` | `ComputedRef<boolean>` | 是否侧边混合菜单布局 |
| `isHeaderNav` | `ComputedRef<boolean>` | 是否顶部菜单布局 |
| `isMixedNav` | `ComputedRef<boolean>` | 是否混合导航布局 |
| `isFullContent` | `ComputedRef<boolean>` | 是否全屏内容布局 |
| `isSideMode` | `ComputedRef<boolean>` | 是否包含侧边栏的布局（mixed/side-mixed/side） |
| `sidebarCollapsed` | `ComputedRef<boolean>` | 侧边栏是否折叠 |

### 主题操作

| 返回值 | 签名 | 说明 |
|---|---|---|
| `setTheme` | `(mode: ThemeModeType) => void` | 设置主题模式 (`light`/`dark`/`auto`) |
| `toggleTheme` | `() => void` | 在 light 和 dark 之间切换 |
| `theme` | `ComputedRef<"dark" \| "light">` | 当前实际主题（已解析 auto） |
| `locale` | `ComputedRef<string>` | 当前语言 |

### 功能判断

| 返回值 | 说明 |
|---|---|
| `keepAlive` | 是否开启 keep-alive（需要 tabbar.enable + tabbar.keepAlive） |
| `contentIsMaximize` | 内容区是否最大化（顶栏和侧边栏都隐藏） |
| `authPanelLeft/Right/Center` | 登录页布局方向 |
| `preferencesButtonPosition` | 偏好设置按钮位置 (`{ fixed, header }`) |
| `globalSearchShortcutKey` | 是否启用全局搜索快捷键 |
| `globalLogoutShortcutKey` | 是否启用全局注销快捷键 |
| `getElementPlusLocale` | Element Plus 当前语言包 |

---

## Preferences 完整字段速查

### `app` — 全局配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `accessMode` | `"backend" \| "frontend"` | `"frontend"` | 权限模式 |
| `authPageLayout` | `"panel-center" \| "panel-left" \| "panel-right"` | `"panel-right"` | 登录页布局 |
| `locale` | `"zh-CN" \| "en-US"` | `"zh-CN"` | 当前语言 |
| `layout` | `LayoutType` | `"sidebar-nav"` | 布局方式 |
| `contentCompact` | `"compact" \| "wide"` | `"wide"` | 内容区宽度 |
| `colorGrayMode` | `boolean` | `false` | 灰色模式 |
| `colorWeakMode` | `boolean` | `false` | 色弱模式 |
| `compact` | `boolean` | `false` | 紧凑模式 |
| `enablePreferences` | `boolean` | `true` | 是否显示偏好设置入口 |
| `preferencesButtonPosition` | `"auto" \| "fixed" \| "header"` | `"auto"` | 偏好按钮位置 |
| `enableRefreshToken` | `boolean` | `false` | 是否开启 RefreshToken |
| `enableTenant` | `boolean` | `false` | 是否启用多租户 |
| `enableCheckUpdates` | `boolean` | `false` | 是否开启版本检查 |
| `watermark` | `boolean` | `false` | 是否开启水印 |
| `loginExpiredMode` | `"modal" \| "page"` | `"page"` | 登录过期提示方式 |
| `isMobile` | `boolean` | `false` | 是否移动端（自动检测，勿手动设置） |

### `theme` — 主题配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `mode` | `"light" \| "dark" \| "auto"` | `"dark"` | 主题模式 |
| `builtinType` | `BuiltinThemeType` | `"default"` | 内置主题名 |
| `colorPrimary` | `string` | `"hsl(212 100% 45%)"` | 主题色 |
| `colorSuccess` | `string` | `"hsl(144 57% 58%)"` | 成功色 |
| `colorWarning` | `string` | `"hsl(42 84% 61%)"` | 警告色 |
| `colorDestructive` | `string` | `"hsl(348 100% 61%)"` | 危险色 |
| `radius` | `string` | `"0.5"` | 圆角大小（rem） |
| `semiDarkSidebar` | `boolean` | `false` | 半深色侧边栏（仅浅色模式） |
| `semiDarkHeader` | `boolean` | `false` | 半深色顶栏（仅浅色模式） |

### `sidebar` — 侧边栏配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `enable` | `boolean` | `true` | 是否可见 |
| `collapsed` | `boolean` | `false` | 是否折叠 |
| `collapsedShowTitle` | `boolean` | `false` | 折叠时是否显示标题 |
| `expandOnHover` | `boolean` | `false` | 鼠标悬停自动展开 |
| `extraCollapse` | `boolean` | `true` | 扩展区域是否折叠 |
| `hidden` | `boolean` | `false` | CSS 隐藏 |
| `width` | `number` | `224` | 侧边栏宽度（px） |

### `tabbar` — 标签页配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `enable` | `boolean` | `true` | 是否开启标签页 |
| `keepAlive` | `boolean` | `true` | 是否缓存页面 |
| `persist` | `boolean` | `true` | 是否持久化标签 |
| `draggable` | `boolean` | `true` | 是否可拖拽 |
| `showIcon` | `boolean` | `true` | 显示图标 |
| `showMaximize` | `boolean` | `true` | 显示最大化按钮 |
| `showMore` | `boolean` | `true` | 显示更多按钮 |
| `styleType` | `"brisk" \| "card" \| "chrome" \| "plain"` | `"chrome"` | 标签页风格 |
| `height` | `number` | `38` | 标签页高度 |

### `header` — 顶栏配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `enable` | `boolean` | `true` | 是否启用 |
| `hidden` | `boolean` | `false` | CSS 隐藏 |
| `mode` | `"auto" \| "auto-scroll" \| "fixed" \| "static"` | `"fixed"` | 顶栏模式 |

### `breadcrumb` — 面包屑配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `enable` | `boolean` | `true` | 是否启用 |
| `hideOnlyOne` | `boolean` | `false` | 只有一个时隐藏 |
| `showHome` | `boolean` | `true` | 显示首页图标 |
| `showIcon` | `boolean` | `true` | 显示图标 |
| `styleType` | `"background" \| "normal"` | `"normal"` | 面包屑风格 |

### `navigation` — 导航配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `accordion` | `boolean` | `true` | 手风琴模式 |
| `split` | `boolean` | `true` | 菜单切割（仅 mixed-nav） |
| `styleType` | `"plain" \| "rounded"` | `"rounded"` | 导航风格 |

### `transition` — 动画配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `enable` | `boolean` | `true` | 是否启用切换动画 |
| `name` | `"fade" \| "fade-down" \| "fade-slide" \| "fade-up" \| string` | `"fade-slide"` | 动画名称 |
| `progress` | `boolean` | `true` | 页面加载进度条 |
| `loading` | `boolean` | `true` | 页面加载 loading |

### `shortcutKeys` — 快捷键配置

| 字段 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `enable` | `boolean` | `true` | 是否启用快捷键 |
| `globalSearch` | `boolean` | `true` | 全局搜索 |
| `globalLogout` | `boolean` | `true` | 全局注销 |
| `globalLockScreen` | `boolean` | `true` | 全局锁屏 |
| `globalPreferences` | `boolean` | `true` | 全局偏好设置 |

### `widget` — 功能部件

| 字段 | 类型 | 默认值 |
|---|---|---|
| `fullscreen` | `boolean` | `true` |
| `globalSearch` | `boolean` | `true` |
| `languageToggle` | `boolean` | `true` |
| `lockScreen` | `boolean` | `true` |
| `notification` | `boolean` | `true` |
| `refresh` | `boolean` | `true` |
| `sidebarToggle` | `boolean` | `true` |
| `themeToggle` | `boolean` | `true` |

### `logo` / `copyright` / `footer`

| 子项 | 常用字段 | 说明 |
|---|---|---|
| `logo` | `enable`, `source` | Logo 可见性和图片地址 |
| `copyright` | `enable`, `companyName`, `companySiteLink`, `date`, `icp`, `icpLink` | 版权信息 |
| `footer` | `enable`, `fixed` | 底栏可见性和固定 |

---

## 内置主题色板

| 主题名 | 色值 |
|---|---|
| `default` | hsl(212 100% 45%) |
| `violet` | hsl(245 82% 67%) |
| `pink` | hsl(347 77% 60%) |
| `yellow` | hsl(42 84% 61%) |
| `sky-blue` | hsl(231 98% 65%) |
| `green` | hsl(161 90% 43%) |
| `zinc` | hsl(240 5% 26%) |
| `deep-green` | hsl(181 84% 32%) |
| `deep-blue` | hsl(211 91% 39%) |
| `orange` | hsl(18 89% 40%) |
| `rose` | hsl(0 75% 42%) |
| `neutral` | hsl(0 0% 25%) |
| `slate` | hsl(215 25% 27%) |
| `gray` | hsl(217 19% 27%) |
| `red` | hsl(0 72% 51%) |
| `stone` | hsl(25 5% 45%) |
| `custom` | 自定义（由 `colorPrimary` 决定） |

---

## 数据流

```
用户操作 → updatePreferences({ theme: { mode: "dark" } })
    │
    ├─→ 合并到 reactive state
    │
    ├─→ handleUpdates() 副作用
    │     ├─ theme.* → updateCSSVariables() → 更新 CSS 变量 + html.dark
    │     ├─ app.locale → loadLocaleMessages() → 加载语言包
    │     └─ app.colorGrayMode/WeakMode → 更新 DOM class
    │
    └─→ savePreferences() (防抖 150ms)
          └─→ StorageManager.setItem() → localStorage
```

---

## 常见场景

### 场景 1：业务代码中切换布局

```ts
import { updatePreferences } from "@/core/preferences";

updatePreferences({ app: { layout: "header-nav" } });
```

### 场景 2：根据业务状态折叠侧边栏

```ts
import { updatePreferences } from "@/core/preferences";

// 折叠
updatePreferences({ sidebar: { collapsed: true } });

// 展开
updatePreferences({ sidebar: { collapsed: false } });
```

### 场景 3：在组件中判断当前布局

```vue
<script setup lang="ts">
import { usePreferences } from "@/core/preferences";

const { isSideNav, isHeaderNav, isMobile } = usePreferences();
</script>

<template>
  <div v-if="isSideNav && !isMobile">侧边栏布局</div>
  <div v-else-if="isHeaderNav">顶部布局</div>
</template>
```

### 场景 4：切换主题色

```ts
import { updatePreferences } from "@/core/preferences";

// 使用内置主题
updatePreferences({ theme: { builtinType: "violet", colorPrimary: "hsl(245 82% 67%)" } });

// 使用自定义颜色
updatePreferences({ theme: { builtinType: "custom", colorPrimary: "#ff6600" } });
```

### 场景 5：重置所有偏好

```ts
import { resetPreferences } from "@/core/preferences";

resetPreferences(); // 恢复为默认值 + overrides，清除 localStorage
```

---

## 偏好设置面板

框架内置了完整的偏好设置面板 UI，位于 `components/PreferencesPanel/`：

- **LayoutPanel.vue** — 布局方式、内容宽度
- **AppearancePanel.vue** — 主题模式、主题色、圆角、灰色/色弱模式
- **GeneralPanel.vue** — 语言、水印、动画、版本检查
- **ShortcutKeyPanel.vue** — 快捷键开关

在布局组件中通过 `<PreferencesPanel />` 挂载即可使用。

---

## 注意事项

1. **不要直接修改 `preferences` 对象** — 它是 `readonly` 的，始终通过 `updatePreferences()` 修改
2. **`isMobile` 由框架自动管理** — 基于 Tailwind `md` 断点检测，无需手动设置
3. **`initPreferences()` 只调用一次** — 在 `bootstrap.ts` 中已完成，业务代码不要重复调用
4. **颜色值支持多种格式** — `colorPrimary` 等字段接受 HSL、RGB、Hex 格式，内部会统一转换
5. **持久化有防抖** — 连续多次 `updatePreferences()` 不会频繁写入 localStorage（150ms 防抖）
