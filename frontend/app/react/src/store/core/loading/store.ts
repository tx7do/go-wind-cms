import {create, type StoreApi, useStore} from 'zustand';
import {createContext, useContext, type Context, useRef} from 'react';

// ==============================
// 类型定义
// ==============================

export interface LoadingState {
    isLoading: boolean;
    error: boolean | null;
}

export interface LoadingStore extends LoadingState {
    start: () => void;
    finish: () => void;
    error_: () => void;
}

// ==============================
// Store 工厂函数（每个请求/客户端独立实例）
// ==============================

export function createLoadingStore(): StoreApi<LoadingStore> {
    return create<LoadingStore>((set) => ({
        isLoading: false,
        error: null,

        start: () => set({isLoading: true, error: null}),

        finish: () => set({isLoading: false}),

        error_: () => set({isLoading: false, error: true}),
    }));
}

// ==============================
// React Context
// ==============================

export const LoadingStoreContext: Context<StoreApi<LoadingStore> | null> = createContext<StoreApi<LoadingStore> | null>(null);

// ==============================
// 消费 Hook
// ==============================

/**
 * 在组件内访问 loading store
 * 必须在 <LoadingStoreProvider> 内部使用
 */
export function useLoadingStore<T = LoadingStore>(selector?: (state: LoadingStore) => T): T {
    const storeApi = useContext(LoadingStoreContext);
    if (!storeApi) {
        throw new Error('useLoadingStore must be used within <LoadingStoreProvider>');
    }
    return useStore(
        storeApi,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (selector ?? ((state: LoadingStore) => state)) as (state: LoadingStore) => T,
    );
}

// ==============================
// 便捷 Hook（兼容旧 API）
// ==============================

/**
 * 提供 loading 状态和操作方法
 * 使用 useRef 暴露稳定的方法引用，避免不必要的 rerender
 */
export function useLoading() {
    const store = useLoadingStore();
    return {
        isLoading: store.isLoading,
        start: store.start,
        finish: store.finish,
        error: store.error_,
    };
}

// ==============================
// 非 Hook 调用支持
// ==============================

/**
 * 获取当前 LoadingStore 实例的便捷方法
 *
 * 注意：必须运行在 React 树内（依赖 Context）
 * 在组件中推荐使用 useLoadingStore() Hook
 *
 * 使用示例（在 hook 内部）：
 *   const loadingApi = useLoadingApi();
 *   loadingApi.start();
 */
export function useLoadingApi(): StoreApi<LoadingStore> {
    const api = useContext(LoadingStoreContext);
    if (!api) {
        throw new Error('useLoadingApi must be used within <LoadingStoreProvider>');
    }
    return api;
}

/**
 * 兼容旧 API 的 Hook 形式
 * 在 hook 内调用，返回 start/finish/error 函数
 */
export function useLoadingActions() {
    const api = useLoadingApi();
    return useRef({
        startLoading: () => api.getState().start(),
        finishLoading: () => api.getState().finish(),
        loadingError: () => api.getState().error_(),
    }).current;
}
