/**
 * 主题模式
 * auto 自动（跟随系统）
 * dark 暗色
 * light 亮色
 */
export type ThemeModeType = "auto" | "dark" | "light";

/**
 * 内置主题名称
 * custom 自定义
 * default 默认（蓝色）
 * green 绿色
 * orange 橙色
 * pink 粉色
 * red 红色
 * sky-blue 天蓝色
 * violet 紫罗兰色
 * yellow 黄色
 */
export type BuiltinThemeType =
  | "custom"
  | "default"
  | "green"
  | "orange"
  | "pink"
  | "red"
  | "sky-blue"
  | "violet"
  | "yellow"
  | (Record<never, never> & string);
