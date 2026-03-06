import {defineStore} from 'pinia';
import {ref, watch} from 'vue';
import {i18n} from '@/locales';

/**
 * 语言切换监听 Store
 *
 * 功能:
 * 1. 监听语言切换事件
 * 2. 提供语言切换回调注册机制
 * 3. 当语言切换时，通知所有注册的回调函数
 * 4. 防抖和频率限制，减少服务器压力
 *
 * 使用方式:
 * - 在需要监听语言切换的页面或组件中调用 watchLanguageChange
 * - 传入回调函数，语言切换时会自动调用
 */
export const useLanguageChangeStore = defineStore('language-change', () => {
  // 当前语言
  const currentLanguage = ref<string>(i18n.global.locale.value);

  // 回调函数集合
  const languageChangeCallbacks = ref<Set<() => void | Promise<void>>>(new Set());

  // 是否正在切换语言
  const isSwitching = ref(false);

  // 防抖定时器
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;

  // 防抖延迟 (毫秒)
  const DEBOUNCE_DELAY = 300;

  // 频率限制：记录最近切换的时间
  const recentSwitches = ref<number[]>([]);

  // 频率限制：最小切换间隔 (毫秒)
  const MIN_SWITCH_INTERVAL = 2000;

  // 频率限制：最大切换次数
  const MAX_SWITCHES_IN_WINDOW = 3;

  // 频率限制：时间窗口 (毫秒)
  const SWITCH_WINDOW = 10000;

  /**
   * 注册语言切换回调
   * @param callback - 语言切换时调用的回调函数
   * @returns 取消订阅函数
   */
  function watchLanguageChange(callback: () => void | Promise<void>): () => void {
    languageChangeCallbacks.value.add(callback);

    // 返回取消订阅函数
    return () => {
      languageChangeCallbacks.value.delete(callback);
    };
  }

  /**
   * 检查频率限制
   * @returns 是否允许切换
   */
  function checkRateLimit(): boolean {
    const now = Date.now();

    // 清理过期的记录
    recentSwitches.value = recentSwitches.value.filter(
      time => now - time < SWITCH_WINDOW
    );

    // 检查是否超过限制
    if (recentSwitches.value.length >= MAX_SWITCHES_IN_WINDOW) {
      console.warn('[Language Change] Rate limit exceeded');
      return false;
    }

    // 检查最小间隔
    const lastSwitch = recentSwitches.value[recentSwitches.value.length - 1];
    if (lastSwitch && now - lastSwitch < MIN_SWITCH_INTERVAL) {
      console.warn('[Language Change] Switching too frequently');
      return false;
    }

    return true;
  }

  /**
   * 记录切换时间
   */
  function recordSwitch() {
    recentSwitches.value.push(Date.now());
  }

  /**
   * 手动触发语言切换（供内部使用）
   * 当 i18n 语言改变时，调用此方法通知所有监听者
   */
  async function triggerLanguageChange() {
    const newLanguage = i18n.global.locale.value;

    if (newLanguage === currentLanguage.value) {
      return; // 语言未改变
    }

    // 检查频率限制
    if (!checkRateLimit()) {
      console.warn('[Language Change] Ignoring due to rate limit');
      // 即使被限制，也要更新当前语言，但不触发回调
      currentLanguage.value = newLanguage;
      isSwitching.value = false;
      return;
    }

    // 清除之前的防抖定时器
    if (debounceTimer) {
      clearTimeout(debounceTimer);
      debounceTimer = null;
    }

    console.log('[Language Change] Language switched from', currentLanguage.value, 'to', newLanguage);

    // 使用防抖，避免快速连续触发
    isSwitching.value = true;

    debounceTimer = setTimeout(async () => {
      try {
        // 记录切换时间
        recordSwitch();

        currentLanguage.value = newLanguage;

        // 并行执行所有回调
        const promises = Array.from(languageChangeCallbacks.value).map(callback => {
          try {
            return Promise.resolve(callback());
          } catch (error) {
            console.error('[Language Change] Callback error:', error);
            return Promise.resolve();
          }
        });

        await Promise.all(promises);
        console.log('[Language Change] All callbacks executed successfully');
      } catch (error) {
        console.error('[Language Change] Error executing callbacks:', error);
      } finally {
        isSwitching.value = false;
        debounceTimer = null;
      }
    }, DEBOUNCE_DELAY);
  }

  /**
   * 初始化语言监听
   * 监听 i18n locale 的变化
   */
  function initializeLanguageWatcher() {
    watch(
      () => i18n.global.locale.value,
      async (newLang, oldLang) => {
        if (newLang !== oldLang) {
          await triggerLanguageChange();
        }
      }
    );

    console.log('[Language Change] Initial watcher initialized');
  }

  /**
   * 获取当前语言
   */
  function getCurrentLanguage(): string {
    return currentLanguage.value;
  }

  /**
   * 清除所有回调（用于清理）
   */
  function clearAllCallbacks() {
    languageChangeCallbacks.value.clear();
  }

  return {
    currentLanguage,
    isSwitching,
    watchLanguageChange,
    triggerLanguageChange,
    initializeLanguageWatcher,
    getCurrentLanguage,
    clearAllCallbacks,
  };
});
