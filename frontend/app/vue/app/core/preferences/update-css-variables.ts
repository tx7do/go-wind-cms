import type { Preferences } from "./types";

/**
 * 更新 CSS 变量 — 适配 shadcn-vue
 *
 * 暗色模式策略（三层）：
 * 1. CSS 后备：@media (prefers-color-scheme: dark) 在 JS 执行前自动跟随系统
 * 2. JS 精确控制：通过 .dark / .light class 覆盖 @media
 *    - mode='auto'  → 移除 .dark 和 .light，由 @media 接管
 *    - mode='dark'  → 添加 .dark，移除 .light
 *    - mode='light' → 添加 .light，移除 .dark
 * 3. data-theme 属性切换主题色预设
 * 4. --radius 控制全局圆角
 */
function updateCSSVariables(preferences: Preferences) {
  if (typeof document === "undefined") return;

  const root = document.documentElement;
  if (!root) return;

  const theme = preferences?.theme ?? {};
  const { mode, builtinType, radius } = theme;

  // 1. 暗色模式 — 通过 .dark / .light class 精确控制
  if (Reflect.has(theme, "mode")) {
    if (mode === "dark") {
      root.classList.add("dark");
      root.classList.remove("light");
    } else if (mode === "light") {
      root.classList.add("light");
      root.classList.remove("dark");
    } else {
      // auto：移除两个 class，让 CSS @media (prefers-color-scheme: dark) 接管
      root.classList.remove("dark");
      root.classList.remove("light");
    }
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
