import type { Preferences } from "./types";

/**
 * 更新 CSS 变量 — 适配 shadcn-vue
 *
 * shadcn-vue 使用：
 * - `.dark` class on <html> 切换暗色模式
 * - `data-theme` 属性切换主题色预设（通过 main.css 中定义的 [data-theme="xxx"] 选择器）
 * - `--radius` 控制全局圆角
 */
function updateCSSVariables(preferences: Preferences) {
  if (typeof document === "undefined") return;

  const root = document.documentElement;
  if (!root) return;

  const theme = preferences?.theme ?? {};
  const { mode, builtinType, radius } = theme;

  // 1. 暗色模式 — 切换 .dark class
  if (Reflect.has(theme, "mode")) {
    const dark = isDarkTheme(mode);
    root.classList.toggle("dark", dark);
  }

  // 2. 主题色预设 — 切换 data-theme 属性
  // main.css 中定义了每种 [data-theme="xxx"] 对应的 --primary HSL 变量
  if (Reflect.has(theme, "builtinType") && builtinType) {
    root.setAttribute("data-theme", builtinType);
  }

  // 3. 圆角 — 更新 --radius
  if (Reflect.has(theme, "radius")) {
    root.style.setProperty("--radius", `${radius}rem`);
  }
}

/**
 * 判断当前是否为暗色主题
 */
function isDarkTheme(theme: string | undefined): boolean {
  if (!theme) return false;
  if (theme === "dark") return true;
  if (theme === "auto") {
    if (typeof window === "undefined") return false;
    return window.matchMedia("(prefers-color-scheme: dark)").matches;
  }
  return false;
}

export { isDarkTheme, updateCSSVariables };
