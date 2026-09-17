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
 * @param char
 */
export const getCharColor = (char: string) => {
  let hash = 0;
  for (let i = 0; i < char.length; i++) {
    hash = char.charCodeAt(i) + ((hash << 5) - hash);
  }
  const hue = Math.abs(hash % 360);
  const saturation = 60;
  const lightness = 45;
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
