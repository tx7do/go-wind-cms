import { usePreferencesStore } from '@/core/preferences';

/**
 * 是否是暗黑模式
 */
export function isDarkMode(): boolean {
  const { theme } = usePreferencesStore.getState().preferences;
  if (theme.mode === 'dark') return true;
  if (theme.mode === 'light') return false;
  // auto 模式：跟随系统
  return window.matchMedia('(prefers-color-scheme: dark)').matches;
}
