import {onUnmounted} from 'vue';
import {useLanguageChangeStore} from '@/stores/modules/app/language-change.state';

/**
 * 语言切换监听 Composable
 *
 * 用于在页面或组件中监听语言切换事件，并自动重新加载数据
 *
 * @param loadData - 语言切换时需要重新加载数据的函数
 * @param options - 配置选项
 * @returns 取消订阅函数
 *
 * @example
 * // 在页面中使用
 * import {useLanguageChangeEffect} from '@/hooks';
 *
 * // 基本用法
 * useLanguageChangeEffect(async () => {
 *   await loadPosts();
 *   await loadCategories();
 * });
 *
 * @example
 * // 带清理功能
 * const unsubscribe = useLanguageChangeEffect(async () => {
 *   await loadData();
 * });
 *
 * // 在组件卸载时调用
 * onUnmounted(() => {
 *   unsubscribe();
 * });
 */
export function useLanguageChangeEffect(
  loadData: () => void | Promise<void>,
  options: {
    // 是否立即加载一次
    immediate?: boolean;
    // 是否在组件卸载时自动取消订阅
    autoCleanup?: boolean;
  } = {}
) {
  const {immediate = true, autoCleanup = true} = options;

  const languageChangeStore = useLanguageChangeStore();

  // 注册语言切换回调
  const unsubscribe = languageChangeStore.watchLanguageChange(async () => {
    console.log('[Language Change Effect] Reloading data...');
    await loadData();
  });

  // 立即执行一次
  if (immediate) {
    loadData();
  }

  // 自动清理
  if (autoCleanup) {
    onUnmounted(() => {
      unsubscribe();
      console.log('[Language Change Effect] Unsubscribed from language changes');
    });
  }

  // 返回取消订阅函数
  return unsubscribe;
}

/**
 * 语言切换监听 Hook (用于页面级别)
 *
 * 与 useLanguageChangeEffect 类似，但专门用于页面级别，
 * 会自动在 onMounted 时初始化监听器
 *
 * @param loadData - 语言切换时需要重新加载数据的函数
 *
 * @example
 * // 在页面中使用
 * import {useLanguageChangeListener} from '@/hooks';
 *
 * useLanguageChangeListener(async () => {
 *   await fetchPageData();
 * });
 */
export function useLanguageChangeListener(
  loadData: () => void | Promise<void>
) {
  const languageChangeStore = useLanguageChangeStore();

  // 注册回调
  const unsubscribe = languageChangeStore.watchLanguageChange(async () => {
    console.log('[Language Change Listener] Reloading page data...');
    await loadData();
  });

  // 组件卸载时自动清理
  onUnmounted(() => {
    unsubscribe();
    console.log('[Language Change Listener] Unsubscribed from language changes');
  });

  return unsubscribe;
}
