# 移动端跨端统一UI设计规范

> **适用终端**：H5 / 微信 / 支付宝 / 抖音小程序 / ReactNative APP  
> **设计基准**：统一 750 宽度设计稿，一套 className 全端复用，脱离单一平台 UI 限制
>
>

---

## 1 基础设计基准

### 1.1 设计稿与单位换算规则

全局统一设计稿宽度：**750**，所有页面、组件严格基于该尺寸开发，框架自动完成多端单位适配，无需差异化编码。

|运行端|编码写法|编译结果|
|---|---|---|
|全品类小程序|`[XXrpx]`|原生 rpx|
|H5 移动端|`[XXrpx]`|自动换算 px（2rpx = 1px）|
|ReactNative/App|`[XXrpx]`|自动换算 dp|

> 编码强制规范：统一使用 `rpx` 任意值格式，**禁止手写 px/rem**，保障全端视觉一致性。
>
>

### 1.2 全局尺寸常量库

项目统一固定尺寸，所有组件、布局严格复用，禁止自定义非标尺寸。

|分类|规格 (rpx)|业务使用场景|
|---|---|---|
|基础间距|8、16、24、32、48、64|内外边距、元素间隙、图文间距、模块间距|
|圆角分级|8、12、16、24、999|输入框、小按钮、常规卡片、大按钮、圆形头像|
|最小触控热区|88|所有可点击控件最小宽高（移动端强制规范，避免触控失效）|

---

## 2 全端统一字体规范

H5、小程序、RN 共用一套字号体系，无差异化配置，实现全端视觉统一。

|文字层级|字号 (rpx)|使用场景|Tailwind 编码示例|
|---|---|---|---|
|PageTitle|48|页面大标题、弹窗头部标题|`text-[48rpx]`|
|CardTitle|36|卡片标题、分区标题、列表主标题|`text-[36rpx]`|
|BodyText|32|页面正文、列表主体文案、常规内容|`text-[32rpx]`|
|DescText|28|备注说明、时间字段、辅助文案|`text-[28rpx]`|
|TipsText|24|标签小字、提示文案、页脚版权|`text-[24rpx]`|

---

## 3 全局通用色板（Tailwind Theme 配置）

项目统一业务色值，无平台绑定，全端通用，直接写入 Tailwind 配置扩展，禁止自定义非标色值。

```js
// tailwind.config.js -> theme.extend.colors
colors: {
  primary: '#1677ff',    // 主色/主按钮/重点强调
  success: '#00b42a',    // 成功状态/完成提示
  warning: '#ff7d00',    // 警告状态/待处理提示
  danger: '#f53f3f',     // 错误状态/删除/高危操作
  textMain: '#1d2129',   // 一级正文/核心内容
  textSec: '#4e5969',    // 二级正文/次要内容
  textThird: '#86909c',  // 三级辅助文字/说明
  textWeak: '#c9cdd4',   // 占位文字/禁用状态文字
  pageBg: '#f2f3f5',     // 页面背景色
  cardBg: '#ffffff',     // 卡片/模块底色
  splitLine: '#e5e6eb'   // 分割线/边框色
}
```

使用示例：`text-primary bg-pageBg border-splitLine`

---

## 4 通用组件尺寸规范

### 4.1 按钮组件

- 通栏主按钮：`h-[88rpx] rounded-[16rpx] px-[48rpx]`

- 常规按钮：高 `72rpx`，圆角 `16rpx`，左右内边距 `32rpx`

- 小型按钮：高 `56rpx`，圆角 `12rpx`，左右内边距 `24rpx`

> 强制约束：所有可点击元素最小尺寸 ≥ 88rpx × 88rpx，保障移动端触控体验。
>
>

### 4.2 列表 Cell

- 单行标准高度：`h-[88rpx]`

- 左右统一内边距：`px-[24rpx]`

- 主标题样式：`text-[32rpx] text-textMain`

- 副文本样式：`text-[28rpx] text-textThird`

### 4.3 输入框

- 基础样式：`h-[80rpx] rounded-[12rpx] px-[24rpx] bg-pageBg`

- 占位文字样式：`text-textWeak`

### 4.4 卡片容器

- 基础样式：`p-[24rpx] rounded-[16rpx] bg-cardBg`

---

## 5 页面布局通用约束

1. **页面留白规范**：所有页面全局左右留白 `px-[24rpx]`，统一内容边距，杜绝贴边布局。

2. **大屏适配规范**：PC/平板大屏 H5，外层容器添加 `max-w-[750rpx] mx-auto`，内容居中展示、两侧留白，避免布局拉伸变形。

3. **全端安全区适配**：全局注入安全区 CSS 变量，适配刘海屏、底部小黑条等异形屏幕。

```css
/* 全局 app.css 通用安全区变量 */
:root {
  --safe-top: env(safe-area-inset-top, 0rpx);
  --safe-bottom: env(safe-area-inset-bottom, 0rpx);
}
```

页面根容器统一添加：`pt-[var(--safe-top)] pb-[var(--safe-bottom)]`

---

## 6 交互统一规范

### 6.1 全局点击反馈

全局自定义按压工具类，实现 H5、小程序双端统一点击效果。

```css
@layer utilities {
  .tap-active { opacity: 0.8; }
}
```

- H5 端：使用 `active:tap-active`

- 小程序端：使用 `hover-class="tap-active"`

### 6.2 底部导航规范

- Tab 栏标准高度：`h-[98rpx]`

- 页面底部统一预留内边距 `pb-[98rpx]`，避免内容被导航栏遮挡

### 6.3 分割线规范

全局统一分割线写法：`h-[1rpx] bg-splitLine`

---

## 7 Tailwind4 项目核心配置

适配 Taro4 多端差异化，解决 H5、小程序样式冲突、失效问题。

```js
/** tailwind.config.js */
export default {
  content: ['./src/**/*.{tsx,jsx}'],
  corePlugins: {
    // 小程序关闭preflight避免wxss冲突，H5开启标准化样式
    preflight: process.env.TARO_ENV !== 'weapp'
  },
  theme: {
    extend: {
      colors: {
        // 粘贴本文档3.章节全局色板配置
      }
    }
  }
}
```

---

## 8 Taro4 多端编码强制规范

1. **单位统一规范**：所有尺寸使用 `[XXrpx]` 任意值语法，禁止手写 px/rem，由框架自动完成全端转换。

2. **类名语法差异化**：


    - H5 端：使用原生 Tailwind 语法，如 `w-1/2`

    - 小程序端：weapp-tailwindcss 自动将 `/` 转为 `_`，无需手动修改

3. **插件环境隔离**：仅小程序端启用 `weapp-tailwindcss` 转换插件，H5 端禁用，防止类名错误转译导致样式失效。

4. **禁止动态拼接类名**：不使用`bg-${color}-500` 动态拼接，Tailwind 无法扫描生成样式，统一使用任意值语法 `bg-[#xxx]`。

---

## 9 小程序 Hero区（头部Banner）规范

### 9.1 核心使用原则

Hero区非全局刚需，**功能页禁用，营销页按需启用**，贴合小程序轻量化设计理念：

- 工具/功能型页面（列表、表单、详情、个人中心）：**禁止使用 Hero区**，避免挤占有效内容可视区

- 营销/展示/活动型页面：**推荐使用 Hero区**，用于品牌展示、活动引流、功能引导

### 9.2 标准尺寸规范（全端通用）

- 宽度：`w-full` 通栏铺满屏幕

- 常规Banner高度：**320rpx**（默认标准，清爽轻量化）

- 大营销活动Banner高度：**360rpx**（高视觉冲击力，项目二选一统一）

- 圆角：底部圆角 `rounded-b-[24rpx]`，顶部跟随系统导航无圆角

- 安全适配：顶部承接状态栏安全区 `pt-[var(--safe-top)]`

### 9.3 标准编码模板

```html
<!-- 通用标准 Hero Banner 组件 -->
<View className="w-full h-[320rpx] rounded-b-[24rpx] pt-[var(--safe-top)] overflow-hidden">
  <Image className="w-full h-full object-cover" src="xxx" mode="widthFix" />
</View>
```

### 9.4 小程序专属禁忌

1. 禁止高度超过 400rpx，避免首屏无有效业务内容，造成用户流失

2. 禁止内容侵入右上角系统胶囊操作区域，预留安全操作空间

3. 禁止复杂动效、视频轮播、多层叠加动画，避免首屏加载卡顿、性能损耗

4. 所有功能型业务页面，强制禁用 Hero 头部横幅

### 9.5 双端差异化说明

- H5端：自由度更高，可自定义更高尺寸、轮播、视频背景、复杂动效

- 小程序端：严格限制尺寸与动效，优先保障加载速度和轻量化体验

**最终原则**：功能优先无Hero，营销展示必规范Hero，全项目视觉、体验统一。
