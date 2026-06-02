import { computed, onUnmounted, watch } from "vue";

import { preferencesManager } from "./preferences";
import { isDarkTheme } from "./update-css-variables";
import type { SupportedLanguagesType, ThemeModeType } from "@/core/preferences/types";

function usePreferences() {
  const preferences = preferencesManager.getPreferences();
  const initialPreferences = preferencesManager.getInitialPreferences();

  const appPreferences = computed(() => preferences.app);
  const logoPreferences = computed(() => preferences.logo);
  const themePreferences = computed(() => preferences.theme);
  const transitionPreferences = computed(() => preferences.transition);
  const copyrightPreferences = computed(() => preferences.copyright);

  /**
   * 判断是否为暗黑模式
   */
  const isDark = computed(() => {
    return isDarkTheme(preferences.theme.mode);
  });

  /**
   * 当前语言（ISO 格式，如 "zh-CN"、"en-US"）
   */
  const locale = computed(() => {
    return preferences.app.locale;
  });

  /**
   * 当前语言的 ISO 代码（供 API 调用使用）
   * 与 @nuxtjs/i18n 的 locale.iso 保持同步
   */
  const localeIso = computed(() => preferences.app.locale);

  /**
   * 应用名
   */
  const appName = computed(() => preferences.app.name);

  /**
   * 是否移动端
   */
  const isMobile = computed(() => preferences.app.isMobile);

  /**
   * 主题模式字符串（"dark" | "light"）
   */
  const theme = computed(() => (isDark.value ? "dark" : "light"));

  /**
   * 圆角值
   */
  const radius = computed(() => preferences.theme.radius);

  /**
   * 动态标题
   */
  const dynamicTitle = computed(() => preferences.app.dynamicTitle);

  /**
   * 页面切换动画是否启用
   */
  const transitionEnable = computed(() => preferences.transition.enable);

  /**
   * 页面切换动画名称
   */
  const transitionName = computed(() => preferences.transition.name);

  /**
   * 设置主题模式
   * @param mode 主题模式（light/dark/auto）
   */
  const setTheme = (mode: ThemeModeType) => {
    preferencesManager.updatePreferences({
      theme: { mode },
    });
  };

  /**
   * 切换主题模式（在浅色和深色之间切换）
   */
  const toggleTheme = () => {
    const newMode = preferences.theme.mode === "dark" ? "light" : "dark";
    setTheme(newMode);
  };

  /**
   * 设置主色
   */
  const setColorPrimary = (color: string) => {
    preferencesManager.updatePreferences({
      theme: { colorPrimary: color, builtinType: "custom" },
    });
  };

  /**
   * 设置语言（同步更新 preferences + @nuxtjs/i18n）
   * preferences 是唯一真相源
   */
  const setLocale = (newLocale: SupportedLanguagesType) => {
    preferencesManager.updatePreferences({
      app: { locale: newLocale },
    });
  };

  /**
   * 偏好设置差异
   */
  const diffPreference = computed(() => {
    return JSON.stringify(initialPreferences) !== JSON.stringify(preferences);
  });

  return {
    appName,
    appPreferences,
    copyrightPreferences,
    diffPreference,
    dynamicTitle,
    isDark,
    isMobile,
    locale,
    localeIso,
    logoPreferences,
    radius,
    setColorPrimary,
    setLocale,
    setTheme,
    theme,
    themePreferences,
    toggleTheme,
    transitionEnable,
    transitionName,
    transitionPreferences,
  };
}

/**
 * 监听语言变化 — 内聚替代 useLanguageChangeEffect + language-change.state.ts
 *
 * @param callback 语言切换时的回调函数
 * @param options 配置项
 * @returns 取消订阅函数
 *
 * @example
 * // 在页面/组件中自动监听语言切换并重新加载数据
 * onLocaleChange(async () => {
 *   await loadPosts()
 * })
 */
function onLocaleChange(
  callback: () => void | Promise<void>,
  options: { immediate?: boolean; autoCleanup?: boolean } = {},
) {
  const { immediate = false, autoCleanup = true } = options;
  const { locale } = useI18n();

  const stop = watch(
    () => locale.value,
    async (newLang, oldLang) => {
      if (newLang !== oldLang) {
        await callback();
      }
    },
  );

  if (immediate) {
    callback();
  }

  if (autoCleanup) {
    onUnmounted(() => stop());
  }

  return stop;
}

export { onLocaleChange, usePreferences };
