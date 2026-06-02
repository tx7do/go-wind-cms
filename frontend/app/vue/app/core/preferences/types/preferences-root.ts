import type {
  AppPreferences,
  CopyrightPreferences,
  LogoPreferences,
  ThemePreferences,
  TransitionPreferences,
} from "./app";

/**
 * 偏好设置总接口（CMS 前台）
 */
export interface Preferences {
  /** 全局配置 */
  app: AppPreferences;
  /** 版权配置 */
  copyright: CopyrightPreferences;
  /** logo 配置 */
  logo: LogoPreferences;
  /** 主题配置 */
  theme: ThemePreferences;
  /** 动画配置 */
  transition: TransitionPreferences;
}

export type PreferencesKeys = keyof Preferences;

export type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

/**
 * 偏好设置初始化参数
 */
export interface InitialOptions {
  /** 命名空间 */
  namespace: string;
  /** 覆盖默认偏好设置 */
  overrides?: DeepPartial<Preferences>;
}
