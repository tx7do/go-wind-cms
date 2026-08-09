/**
 * 支持的语言类型
 * 与应用使用的语言值保持一致
 *
 * 注意：'ar' 为阿拉伯语（RTL）。其 <html dir="rtl"> 由 @nuxtjs/i18n 依据
 * nuxt.config.ts 中该 locale 的 dir 字段自动注入，无需应用层手写。
 */
export type SupportedLanguagesType = "zh-CN" | "en-US" | "ar";

/**
 * 页面切换动画
 */
export type PageTransitionType = "fade" | "fade-down" | "fade-slide" | "fade-up";
