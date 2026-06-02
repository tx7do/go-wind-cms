import { markRaw, reactive, readonly, watch } from "vue";

import { breakpointsTailwind, useBreakpoints, useDebounceFn } from "@vueuse/core";

import { defaultPreferences } from "./config/default";
import type { DeepPartial, InitialOptions, Preferences } from "./types";
import type { SupportedLanguagesType } from "./types";
import { updateCSSVariables } from "./update-css-variables";
import { StorageManager } from "@/core/storage";
import { merge } from "@/utils/merge";
import { isMacOs } from "@/utils/inference";

const STORAGE_KEY = "preferences";
const STORAGE_KEY_LOCALE = `${STORAGE_KEY}-locale`;
const STORAGE_KEY_THEME = `${STORAGE_KEY}-theme`;

class PreferenceManager {
  private cache: null | StorageManager = null;
  private initialPreferences: Preferences = defaultPreferences;
  private isInitialized: boolean = false;
  private readonly savePreferences: (preference: Preferences) => void;
  private state: Preferences = reactive<Preferences>({
    ...this.loadPreferences(),
  });

  constructor() {
    this.cache = new StorageManager();

    // 避免频繁的操作缓存
    this.savePreferences = useDebounceFn(
      (preference: Preferences) => this._savePreferences(preference),
      150,
    );
  }

  /**
   * 保存偏好设置
   */
  private _savePreferences(preference: Preferences) {
    this.cache?.setItem(STORAGE_KEY, preference);
    this.cache?.setItem(STORAGE_KEY_LOCALE, preference.app.locale);
    this.cache?.setItem(STORAGE_KEY_THEME, preference.theme.mode);
  }

  /**
   * 处理更新的键值
   */
  private handleUpdates(updates: DeepPartial<Preferences>) {
    const themeUpdates = updates.theme || {};

    if (themeUpdates && Object.keys(themeUpdates).length > 0) {
      updateCSSVariables(this.state);
    }

    // 语言切换：同步到 @nuxtjs/i18n
    const appUpdates = updates.app || {};
    if (Reflect.has(appUpdates, "locale") && appUpdates.locale) {
      this.syncLocale(appUpdates.locale as SupportedLanguagesType);
    }
  }

  /**
   * 同步 locale 到 @nuxtjs/i18n + 存储
   * preferences 是唯一真相源，i18n 跟随它。
   */
  private syncLocale(locale: SupportedLanguagesType) {
    this.cache?.setItem(STORAGE_KEY_LOCALE, locale);
    // 在 Nuxt 上下文中同步 i18n locale
    // preferences.ts 可能不在 setup() 内调用，所以安全获取
    if (typeof window !== "undefined") {
      try {
        const nuxtApp = (window as any).__NUXT__;
        if (nuxtApp?.locale?.setLocale) {
          // 映射 iso → code: zh-CN → zh, en-US → en
          const code = locale.split("-")[0];
          nuxtApp.locale.setLocale(code);
        }
      } catch {
        // 静默处理——plugin 初始化时 i18n 可能尚未就绪
      }
    }
  }

  private initPlatform() {
    if (typeof document === "undefined") return;
    const dom = document.documentElement;
    dom.dataset.platform = isMacOs() ? "macOs" : "window";
  }

  /**
   * 从缓存中加载偏好设置
   */
  private loadCachedPreferences() {
    return this.cache?.getItem<Preferences>(STORAGE_KEY);
  }

  /**
   * 加载偏好设置
   */
  private loadPreferences(): Preferences {
    return this.loadCachedPreferences() || { ...defaultPreferences };
  }

  /**
   * 监听状态和系统偏好设置的变化
   */
  private setupWatcher() {
    if (this.isInitialized) {
      return;
    }

    // 监听断点，判断是否移动端
    const breakpoints = useBreakpoints(breakpointsTailwind);
    const isMobile = breakpoints.smaller("md");
    watch(
      () => isMobile.value,
      (val) => {
        this.updatePreferences({
          app: { isMobile: val },
        });
      },
      { immediate: true },
    );

    // 监听系统主题偏好设置变化
    if (typeof window !== "undefined") {
      window
        .matchMedia("(prefers-color-scheme: dark)")
        .addEventListener("change", ({ matches: isDark }) => {
          // 仅在 auto 模式下跟随系统
          if (this.state.theme.mode === "auto") {
            updateCSSVariables(this.state);
          }
        });
    }
  }

  clearCache() {
    [STORAGE_KEY, STORAGE_KEY_LOCALE, STORAGE_KEY_THEME].forEach((key) => {
      this.cache?.removeItem(key);
    });
  }

  public getInitialPreferences() {
    return this.initialPreferences;
  }

  public getPreferences() {
    return readonly(this.state);
  }

  /**
   * 覆盖偏好设置
   */
  public async initPreferences({ namespace, overrides }: InitialOptions) {
    if (this.isInitialized) {
      return;
    }

    // 初始化存储管理器
    this.cache = new StorageManager({ prefix: namespace });

    // 合并初始偏好设置
    this.initialPreferences = merge({}, overrides, defaultPreferences);

    // 加载并合并当前存储的偏好设置
    const mergedPreference = merge(
      {},
      this.loadCachedPreferences() || {},
      this.initialPreferences,
    );

    // 从 @nuxtjs/i18n 恢复 locale（如果 i18n 已初始化且与存储不同）
    // 这样 nuxt-i18n 的 URL 路由前缀可以覆盖存储值
    this.restoreLocaleFromI18n(mergedPreference);

    // 更新偏好设置
    this.updatePreferences(mergedPreference);

    this.setupWatcher();
    this.initPlatform();

    this.isInitialized = true;
  }

  /**
   * 从 @nuxtjs/i18n 的当前 locale 恢复 preferences
   * 优先级：URL 路由 > localStorage > 默认值
   */
  private restoreLocaleFromI18n(prefs: Preferences) {
    if (typeof window === "undefined") return;
    try {
      const nuxtApp = (window as any).__NUXT__;
      if (nuxtApp?.locale?.locale) {
        const code = nuxtApp.locale.locale.value as string; // e.g. "zh-CN" or "en-US"
        prefs.app.locale = code as SupportedLanguagesType;
      }
    } catch {
      // i18n 尚未就绪
    }
  }

  /**
   * 重置偏好设置
   */
  resetPreferences() {
    Object.assign(this.state, this.initialPreferences);
    this.savePreferences(this.state);
    [STORAGE_KEY, STORAGE_KEY_THEME, STORAGE_KEY_LOCALE].forEach((key) => {
      this.cache?.removeItem(key);
    });
    this.updatePreferences(this.state);
  }

  /**
   * 更新偏好设置
   */
  public updatePreferences(updates: DeepPartial<Preferences>) {
    const mergedState = merge({}, updates, markRaw(this.state));
    Object.assign(this.state, mergedState);
    this.handleUpdates(updates);
    this.savePreferences(this.state);
  }
}

const preferencesManager = new PreferenceManager();
export { PreferenceManager, preferencesManager };
