import type { BuiltinThemeType } from '../types';

interface BuiltinThemePreset {
  color: string;
  darkPrimaryColor?: string;
  primaryColor?: string;
  type: BuiltinThemeType;
}

const BUILT_IN_THEME_PRESETS: BuiltinThemePreset[] = [
  {
    color: '#006BE6',
    type: 'default',
  },
  {
    color: 'hsl(245 82% 67%)',
    type: 'violet',
  },
  {
    color: 'hsl(347 77% 60%)',
    type: 'pink',
  },
  {
    color: 'hsl(42 84% 61%)',
    type: 'yellow',
  },
  {
    color: 'hsl(231 98% 65%)',
    type: 'sky-blue',
  },
  {
    color: 'hsl(161 90% 43%)',
    type: 'green',
  },
  {
    color: 'hsl(240 5% 26%)',
    darkPrimaryColor: 'hsl(0 0% 98%)',
    primaryColor: 'hsl(240 5.9% 10%)',
    type: 'zinc',
  },

  {
    color: 'hsl(181 84% 32%)',
    type: 'deep-green',
  },

  {
    color: 'hsl(211 91% 39%)',
    type: 'deep-blue',
  },
  {
    color: 'hsl(18 89% 40%)',
    type: 'orange',
  },
  {
    color: 'hsl(0 75% 42%)',
    type: 'rose',
  },

  {
    color: 'hsl(0 0% 25%)',
    darkPrimaryColor: 'hsl(0 0% 98%)',
    primaryColor: 'hsl(240 5.9% 10%)',
    type: 'neutral',
  },
  {
    color: 'hsl(215 25% 27%)',
    darkPrimaryColor: 'hsl(0 0% 98%)',
    primaryColor: 'hsl(240 5.9% 10%)',
    type: 'slate',
  },
  {
    color: 'hsl(217 19% 27%)',
    darkPrimaryColor: 'hsl(0 0% 98%)',
    primaryColor: 'hsl(240 5.9% 10%)',
    type: 'gray',
  },
  {
    color: '',
    type: 'custom',
  },
];

export const COLOR_PRESETS = [...BUILT_IN_THEME_PRESETS].slice(0, 7);

/**
 * 「次要文字」档（docs/design-language.md §2.2）。
 *
 * 浅色侧不直接用 gray-500 #6B7280：它是按白底量的，实测离开白底就掉下 4.5——
 * 白 4.83 / 页面大底 #F5F5F5 4.43 / 卡片灰 #F0F2F5 4.31 / 标签栏 #F0F0F0 4.24。
 * 沿同一色相压暗到 #676E7C 后四处全部达标：5.12 / 4.70 / 4.57 / 4.50（数值由
 * `node .zcode/tmp/hue-calc.mjs` 同口径的 sRGB 相对亮度算出）。
 * 三端共用这两个值（vue-element 的 --el-text-color-secondary、vue-vben 的 --gray-*），
 * 改动须同步 docs/design-language.md §2.2 与三端。
 */
export const SECONDARY_TEXT_LIGHT = '#676E7C';

/** 暗色侧的次要文字 = §2.3 的 #8B949E（L0 6.23 / L1 5.77 / L2 5.34 / L3 4.77 全达标） */
export const SECONDARY_TEXT_DARK = '#8B949E';

export { BUILT_IN_THEME_PRESETS };

export type { BuiltinThemePreset };
