import {create, type StoreApi, useStore} from 'zustand';
import {createContext, useContext, type Context} from 'react';

// ==============================
// 类型定义
// ==============================

export interface IUser {
    id: number;
    username: string;
    nickname: string;
    realname: string;
    avatar: string;
    roles?: string[];
    homePage: string;
}

interface UserState {
    user: IUser | null;
}

export interface UserStore extends UserState {
    setUser: (user: IUser | null) => void;
    clearUser: () => void;
}

// ==============================
// Store 工厂函数（每个请求/客户端独立实例）
// ==============================

export function createUserStore(): StoreApi<UserStore> {
    return create<UserStore>((set) => ({
        user: null,

        setUser: (user) => set({user}),

        clearUser: () => set({user: null}),
    }));
}

// ==============================
// React Context
// ==============================

export const UserStoreContext: Context<StoreApi<UserStore> | null> = createContext<StoreApi<UserStore> | null>(null);

// ==============================
// 消费 Hook
// ==============================

/**
 * 在组件内访问 user store
 * 必须在 <UserStoreProvider> 内部使用
 */
export function useUserStore<T = UserStore>(selector?: (state: UserStore) => T): T {
    const storeApi = useContext(UserStoreContext);
    if (!storeApi) {
        throw new Error('useUserStore must be used within <UserStoreProvider>');
    }
    return useStore(
        storeApi,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (selector ?? ((state: UserStore) => state)) as (state: UserStore) => T,
    );
}

// ==============================
// 便捷 Hooks
// ==============================

export function useUser(): IUser | null {
    return useUserStore((state) => state.user);
}
