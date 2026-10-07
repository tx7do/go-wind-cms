import { create } from 'zustand';

/**
 * 页面刷新管理器
 * 用于全局触发特定页面的刷新
 */
interface PageRefreshState {
  // 当前激活页面的刷新触发器
  refreshTrigger: number;
  
  // 刷新回调函数集合
  refreshCallbacks: Map<string, () => void | Promise<void>>;
  
  // 设置刷新回调
  setRefreshCallback: (pageKey: string, callback: () => void | Promise<void>) => void;
  
  // 移除刷新回调
  removeRefreshCallback: (pageKey: string) => void;
  
  // 触发当前页面刷新
  triggerRefresh: () => void;
  
  // 增加刷新触发器（强制刷新）
  incrementTrigger: () => void;
}

export const usePageRefreshStore = create<PageRefreshState>((set, get) => ({
  refreshTrigger: 0,
  refreshCallbacks: new Map(),
  
  setRefreshCallback: (pageKey, callback) => {
    set((state) => {
      const newCallbacks = new Map(state.refreshCallbacks);
      newCallbacks.set(pageKey, callback);
      return { refreshCallbacks: newCallbacks };
    });
  },
  
  removeRefreshCallback: (pageKey) => {
    set((state) => {
      const newCallbacks = new Map(state.refreshCallbacks);
      newCallbacks.delete(pageKey);
      return { refreshCallbacks: newCallbacks };
    });
  },
  
  triggerRefresh: () => {
    const { refreshCallbacks } = get();
    
    // 触发所有注册的回调
    refreshCallbacks.forEach(async (callback) => {
      try {
        await callback();
      } catch (error) {
        console.error('[PageRefresh] Error during refresh:', error);
      }
    });
    
    // 增加 refreshTrigger，强制重新渲染 Outlet
    set((state) => ({ refreshTrigger: state.refreshTrigger + 1 }));
  },
  
  incrementTrigger: () => {
    set((state) => ({ refreshTrigger: state.refreshTrigger + 1 }));
  },
}));
