import { preferencesManager } from "./preferences";
import type { Preferences } from "./types";

// 偏好设置（带有层级关系）
const preferences: Preferences = preferencesManager.getPreferences.apply(preferencesManager);

// 更新偏好设置
const updatePreferences = preferencesManager.updatePreferences.bind(preferencesManager);

// 重置偏好设置
const resetPreferences = preferencesManager.resetPreferences.bind(preferencesManager);

// 清空偏好设置缓存
const clearPreferencesCache = preferencesManager.clearCache.bind(preferencesManager);

// 初始化偏好设置
const initPreferences = preferencesManager.initPreferences.bind(preferencesManager);

export {
  clearPreferencesCache,
  initPreferences,
  preferences,
  preferencesManager,
  resetPreferences,
  updatePreferences,
};

export * from "./config";
export * from "./types";
export * from "./use-preferences";
