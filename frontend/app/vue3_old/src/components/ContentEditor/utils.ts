import { useAppStore } from '@/stores'

/**
 * 检查当前是否为暗黑模式
 */
export function isDarkMode(): boolean {
  try {
    const appStore = useAppStore()
    return appStore.mode === 'dark'
  } catch {
    // 如果 store 未初始化，从 DOM 或 localStorage 获取
    return (
      document.documentElement.classList.contains('dark') ||
      localStorage.getItem('theme') === 'dark'
    )
  }
}

/**
 * 获取当前主题模式
 */
export function getThemeMode(): 'light' | 'dark' {
  return isDarkMode() ? 'dark' : 'light'
}
