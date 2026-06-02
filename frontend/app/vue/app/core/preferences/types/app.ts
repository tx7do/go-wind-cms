import type { BuiltinThemeType, ThemeModeType } from "./theme";
import type { PageTransitionType, SupportedLanguagesType } from "./layout";

/**
 * 应用偏好设置（CMS 前台）
 */
export interface AppPreferences {
  /** 应用名 */
  name: string;
  /** 应用版本 */
  version: string;
  /** 支持的语言 */
  locale: SupportedLanguagesType;
  /** 是否移动端 */
  isMobile: boolean;
  /** 应用默认头像 */
  defaultAvatar: string;
  /** 开启动态标题 */
  dynamicTitle: boolean;
}

/**
 * 主题偏好设置（CMS 前台）
 */
export interface ThemePreferences {
  /** 内置主题名 */
  builtinType: BuiltinThemeType;
  /** 主题色 */
  colorPrimary: string;
  /** 当前主题模式 */
  mode: ThemeModeType;
  /** 圆角 */
  radius: string;
}

/**
 * Logo 偏好设置
 */
export interface LogoPreferences {
  /** logo 是否可见 */
  enable: boolean;
  /** logo 地址 */
  source: string;
}

/**
 * 版权偏好设置
 */
export interface CopyrightPreferences {
  /** 版权是否可见 */
  enable: boolean;
  /** 版权公司名 */
  companyName: string;
  /** 版权公司名链接 */
  companySiteLink: string;
  /** 版权日期 */
  date: string;
  /** 备案号 */
  icp: string;
  /** 备案号链接 */
  icpLink: string;
}

/**
 * 页面切换动画偏好设置
 */
export interface TransitionPreferences {
  /** 页面切换动画是否启用 */
  enable: boolean;
  /** 是否开启页面加载 loading */
  loading: boolean;
  /** 页面切换动画 */
  name: PageTransitionType | string;
  /** 是否开启页面加载进度动画 */
  progress: boolean;
}
