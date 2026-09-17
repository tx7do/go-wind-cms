/**
 * 主题模式
 * auto 自动
 * dark 暗色
 * light 亮色
 */
export type ThemeModeType = "auto" | "dark" | "light";

/**
 * 内置主题名称
 * custom 自定义
 * deep-blue 深蓝
 * deep-green 深绿
 * default 默认
 * gray 灰色
 * green 绿色
 * neutral 中性色
 * orange 橙色
 * pink 粉色
 * red 红色
 * rose 玫红色
 * sky-blue 天蓝色
 * slate 石板灰色
 * stone 石灰色
 * violet 紫罗兰色
 * yellow 黄色
 * zinc 锌色
 * (Record<never, never> & string) 自定义主题
 */
export type BuiltinThemeType =
    | "custom"
    | "deep-blue"
    | "deep-green"
    | "default"
    | "gray"
    | "green"
    | "neutral"
    | "orange"
    | "pink"
    | "red"
    | "rose"
    | "sky-blue"
    | "slate"
    | "stone"
    | "violet"
    | "yellow"
    | "zinc"
    | (Record<never, never> & string);
