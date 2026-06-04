/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{ts,tsx,js,jsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      /* ========== 全局通用色板（规范 §3） ========== */
      colors: {
        // 主色系
        primary: '#1677ff',       // 主色/主按钮/重点强调
        'primary-foreground': '#ffffff', // 主色上的文字色
        success: '#00b42a',       // 成功状态/完成提示
        warning: '#ff7d00',       // 警告状态/待处理提示
        danger: '#f53f3f',        // 错误状态/删除/高危操作
        
        // 文字色（通过 CSS 变量支持主题切换）
        textMain: 'var(--color-text-main)',    // 一级正文/核心内容
        textSec: 'var(--color-text-sec)',      // 二级正文/次要内容
        textThird: 'var(--color-text-third)', // 三级辅助文字/说明
        textWeak: 'var(--color-text-weak)',   // 占位文字/禁用状态文字
        
        // 背景色（通过 CSS 变量支持主题切换）
        pageBg: 'var(--color-page-bg)',       // 页面背景色
        cardBg: 'var(--color-card-bg)',       // 卡片/模块底色
        splitLine: 'var(--color-split-line)', // 分割线/边框色
        
        // shadcn/ui 兼容色值（映射到自定义变量）
        border: 'var(--color-split-line)',    // 边框色（兼容 border-border）
        background: 'var(--color-page-bg)',   // 背景色（兼容 bg-background）
        foreground: 'var(--color-text-main)', // 前景色（兼容 text-foreground）
        muted: 'var(--color-split-line)',     // 弱化背景（兼容 bg-muted）
        'muted-foreground': 'var(--color-text-third)', // 弱化文字（兼容 text-muted-foreground）
        popover: 'var(--color-card-bg)',      // 弹出层背景（兼容 bg-popover）
        input: 'var(--color-split-line)',     // 输入框边框（兼容 border-input）
        ring: 'var(--color-text-main)',       // 聚焦环（兼容 ring-ring）
      },
      /* ========== 圆角分级（规范 §1.2） ========== */
      borderRadius: {
        'sm': '8rpx',    // 输入框
        'md': '12rpx',   // 小按钮
        DEFAULT: '16rpx', // 常规卡片、按钮
        'lg': '24rpx',   // 大按钮
        'full': '999rpx', // 圆形头像
      },
      /* ========== 间距常量（规范 §1.2） ========== */
      spacing: {
        'sp-8': '8rpx',
        'sp-16': '16rpx',
        'sp-24': '24rpx',
        'sp-32': '32rpx',
        'sp-48': '48rpx',
        'sp-64': '64rpx',
      },
      /* ========== 字号体系（规范 §2） ========== */
      fontSize: {
        'title': '48rpx',   // PageTitle 页面大标题
        'card-title': '36rpx', // CardTitle 卡片标题
        'body': '32rpx',    // BodyText 正文
        'desc': '28rpx',    // DescText 辅助文案
        'tips': '24rpx',    // TipsText 标签小字
      },
      /* ========== 最小触控热区（规范 §1.2） ========== */
      minHeight: {
        'touch': '88rpx',
      },
      minWidth: {
        'touch': '88rpx',
      },
    },
  },
  // 小程序不支持的关键字选择器：*, ::before, ::after, ::backdrop, :hover, :focus 等
  corePlugins: {
    preflight: false,
    space: false,
    divideWidth: false,
    divideColor: false,
    divideStyle: false,
    ringWidth: false,
    ringColor: false,
    ringOffsetWidth: false,
    ringOffsetColor: false,
    appearance: false,
    cursor: false,
    placeholderColor: false,
    outline: false,
  },
}
