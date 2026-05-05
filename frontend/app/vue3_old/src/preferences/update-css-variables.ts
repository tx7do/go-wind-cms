import type { Preferences } from './types';

/**
 * 更新主题的 CSS 变量以及其他 CSS 变量
 * @param _preferences - 当前偏好设置对象，它的主题值将被用来设置文档的主题。
 */
function updateCSSVariables(_preferences: Preferences) {

}

function isDarkTheme(theme: string) {
  let dark = theme === 'dark';
  if (theme === 'auto') {
    dark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  }
  return dark;
}

export { isDarkTheme, updateCSSVariables };
