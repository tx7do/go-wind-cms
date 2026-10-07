/**
 * antd Tag 预设色名池：交给 antd 后由主题算法推导底色/边框/文字，
 * 暗色模式自动适配。勿在此返回 hex/hsl——自定义色不参与主题适配，
 * 暗色下会出现亮底深字的刺眼组合（历史教训，勿回退）。
 */
const TAG_PRESET_COLORS = [
  'blue',
  'green',
  'gold',
  'red',
  'purple',
  'cyan',
  'magenta',
  'orange',
] as const;

export type TagPresetColor = (typeof TAG_PRESET_COLORS)[number];

/**
 * 生成基于字符串的确定性预设色（同一字符串恒定同色），供 Tag color 使用
 */
export const getPresetColor = (str: string): TagPresetColor => {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  return TAG_PRESET_COLORS[Math.abs(hash) % TAG_PRESET_COLORS.length];
};

/**
 * @deprecated 旧名，等价于 getPresetColor。新代码请用 getPresetColor。
 */
export const getRandomColor = getPresetColor;

/**
 * 根据首字母生成固定头像底色（CSS 色，供 Avatar backgroundColor 使用；
 * Tag 场景请用 getPresetColor）
 *
 * 明度取 28% 是量出来的：本函数的产物一律配 antd/EP 头像的白字，而 hsl 同饱和度下
 * 45% 明度是"两种墨都不合格"的死区（360 个色相里最差 white 2.03 / dark 2.01）；
 * 28% + s60 保证任意色相下白字 ≥4.86:1。
 * @param char
 */
export const getCharColor = (char: string) => {
  let hash = 0;
  for (let i = 0; i < char.length; i++) {
    hash = char.charCodeAt(i) + ((hash << 5) - hash);
  }
  const hue = Math.abs(hash % 360);
  const saturation = 60;
  const lightness = 28;
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
};

// 辅助函数：将十六进制颜色转换为 RGB
export function hexToRgb(hex: string): [number, number, number] {
  const bigint = parseInt(hex.slice(1), 16);
  return [(bigint >> 16) & 255, (bigint >> 8) & 255, bigint & 255];
}

// 辅助函数：将 RGB 转换为十六进制颜色
export function rgbToHex(r: number, g: number, b: number): string {
  return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`;
}

/** 解析 `#rgb`/`#rrggbb`/`hsl(h s% l%)`/`hsl(h, s%, l%)` 为 [r,g,b]（0-255）；无法识别返回 null */
export function parseColor(c: string): [number, number, number] | null {
  const s = c.trim();
  if (s.startsWith('#')) {
    const hex = s.slice(1);
    const full = hex.length === 3 ? hex.split('').map((x) => x + x).join('') : hex.slice(0, 6);
    if (!/^[0-9a-fA-F]{6}$/.test(full)) return null;
    return hexToRgb(`#${full}`);
  }
  const hsl = s.match(/^hsla?\(([^)]+)\)$/i);
  if (hsl) {
    const parts = hsl[1].split(/[\s,/]+/).filter(Boolean);
    if (parts.length < 3) return null;
    const h = parseFloat(parts[0]) / 360;
    const sl = parseFloat(parts[1]) / 100;
    const l = parseFloat(parts[2]) / 100;
    if ([h, sl, l].some((n) => Number.isNaN(n))) return null;
    const f = (n: number) => {
      const k = (n + h * 12) % 12;
      return l - sl * Math.min(l, 1 - l) * Math.max(-1, Math.min(k - 3, 9 - k, 1));
    };
    return [f(0) * 255, f(8) * 255, f(4) * 255].map((v) => Math.round(v)) as [number, number, number];
  }
  return null;
}

/**
 * 与 CSS `color-mix(in srgb, base pct%, other)` 同口径的 JS 版本。
 * 存在的理由：antd 的 token 值（colorLink、Menu.itemSelectedColor 等）会被 antd 再做
 * 一次色彩运算，不能填 `var(--gowind-*)`，所以「语义色作文字」的降档/升档要在 TS 里算一遍。
 */
export function mixColor(base: string, other: string, pct: number): string {
  const b = parseColor(base);
  const o = parseColor(other);
  if (!b || !o) return base;
  const p = Math.max(0, Math.min(1, pct));
  return rgbToHex(...(b.map((v, i) => Math.round(p * v + (1 - p) * o[i])) as [number, number, number]));
}
