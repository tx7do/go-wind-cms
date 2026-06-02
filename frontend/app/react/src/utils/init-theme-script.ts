/**
 * 主题初始化脚本（内联到 HTML <head> 中）
 *
 * 在页面加载时立即执行，从 zustand persist 存储（localStorage key: "app-preferences"）
 * 读取用户主题偏好，并应用到 <html> 上，防止主题闪烁（FOUC）。
 *
 * zustand persist 的存储格式：
 * { "state": { "preferences": { "theme": { "mode": "dark" | "light" | "auto" }, ... } } }
 */
export default function initThemeScript() {
    const STORAGE_KEY = 'app-preferences';

    let themeMode = 'auto'; // 默认跟随系统

    try {
        const raw = localStorage.getItem(STORAGE_KEY);
        if (raw) {
            const parsed = JSON.parse(raw);
            const mode = parsed?.state?.preferences?.theme?.mode;
            if (mode === 'dark' || mode === 'light' || mode === 'auto') {
                themeMode = mode;
            }
        }
    } catch {
        // 忽略解析错误，使用默认值
    }

    let themeToUse: string;

    if (themeMode === 'auto') {
        const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        themeToUse = isDark ? 'dark' : 'light';
    } else {
        themeToUse = themeMode;
    }

    document.documentElement.classList.add(themeToUse);
}
