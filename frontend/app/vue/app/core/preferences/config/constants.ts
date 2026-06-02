import type { BuiltinThemeType } from "../types";

interface BuiltinThemePreset {
  /** 预设类型标识，对应 main.css 中的 [data-theme="xxx"] */
  type: BuiltinThemeType;
  /** 色块颜色（用于 PreferencesPanel 中展示，HSL 格式） */
  color: string;
  /** 暗色模式下展示的颜色（可选） */
  darkPrimaryColor?: string;
  /** 浅色模式下展示的颜色（可选） */
  primaryColor?: string;
}

const BUILT_IN_THEME_PRESETS: BuiltinThemePreset[] = [
  { color: "hsl(220 100% 55%)", type: "default" },
  { color: "hsl(245 82% 67%)", type: "violet" },
  { color: "hsl(347 77% 60%)", type: "pink" },
  { color: "hsl(42 84% 61%)", type: "yellow" },
  { color: "hsl(200 90% 50%)", type: "sky-blue" },
  { color: "hsl(161 90% 43%)", type: "green" },
  { color: "hsl(18 89% 40%)", type: "orange" },
  { color: "hsl(0 72% 51%)", type: "red" },
  { color: "", type: "custom" },
];

export const COLOR_PRESETS = [...BUILT_IN_THEME_PRESETS].slice(0, 8);

export { BUILT_IN_THEME_PRESETS };

export type { BuiltinThemePreset };
