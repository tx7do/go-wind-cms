import {create, type StoreApi, useStore} from 'zustand';
import {persist} from 'zustand/middleware';
import {createContext, useContext, type Context} from 'react';

// ==============================
// 类型定义
// ==============================

export interface TokenPayload {
    value: string;
    expiresAt?: number;
}

export interface AccessState {
    permissions: string[];
    accessToken: TokenPayload | null;
    refreshToken: TokenPayload | null;
    isAccessChecked: boolean;
    loginExpired: boolean;
}

// ==============================
// Store 接口
// ==============================

export interface AccessStore extends AccessState {
    setAccessToken: (token: TokenPayload | null) => void;
    setRefreshToken: (token: TokenPayload | null) => void;
    setTokens: (tokens: { accessToken: TokenPayload | null; refreshToken: TokenPayload | null }) => void;
    setPermissions: (permissions: string[]) => void;
    addPermission: (permission: string) => void;
    removePermission: (permission: string) => void;
    setIsAccessChecked: (checked: boolean) => void;
    setLoginExpired: (expired: boolean) => void;
    clearTokens: () => void;
    resetAccess: () => void;
    restoreAccess: (payload: Partial<AccessState>) => void;
}

// ==============================
// Store 工厂函数
// ==============================

export function createAccessStore(): StoreApi<AccessStore> {
    return create<AccessStore>()(
        persist(
            (set) => ({
                // 初始状态
                permissions: [],
                accessToken: null,
                refreshToken: null,
                isAccessChecked: false,
                loginExpired: false,

                setAccessToken: (token) => {
                    console.log('setAccessToken', token);
                    set({accessToken: token});
                },

                setRefreshToken: (token) => {
                    console.log('setRefreshToken', token);
                    set({refreshToken: token});
                },

                setTokens: ({accessToken, refreshToken}) => {
                    console.log('setTokens', {accessToken, refreshToken});
                    set({accessToken, refreshToken});
                },

                setPermissions: (permissions) => set({permissions}),

                addPermission: (permission) => set((state) => {
                    if (state.permissions.includes(permission)) return state;
                    return {permissions: [...state.permissions, permission]};
                }),

                removePermission: (permission) => set((state) => ({
                    permissions: state.permissions.filter(p => p !== permission),
                })),

                setIsAccessChecked: (checked) => set({isAccessChecked: checked}),

                setLoginExpired: (expired) => set({loginExpired: expired}),

                clearTokens: () => {
                    console.log('clearTokens');
                    set({accessToken: null, refreshToken: null, loginExpired: false});
                },

                resetAccess: () => {
                    console.log('resetAccess');
                    set({
                        permissions: [],
                        accessToken: null,
                        refreshToken: null,
                        isAccessChecked: false,
                        loginExpired: false,
                    });
                },

                restoreAccess: (payload) => {
                    console.log('restoreAccess', payload);
                    set((state) => ({...state, ...payload}));
                },
            }),
            {
                name: 'gowind-access-storage', // localStorage key
                partialize: (state) => ({
                    accessToken: state.accessToken,
                    refreshToken: state.refreshToken,
                }),
            }
        )
    );
}

// ==============================
// React Context
// ==============================

export const AccessStoreContext: Context<StoreApi<AccessStore> | null> = createContext<StoreApi<AccessStore> | null>(null);

// ==============================
// 消费 Hook
// ==============================

/**
 * 在组件内访问 access store
 * 必须在 <AccessStoreProvider> 内部使用
 */
export function useAccessStore<T = AccessStore>(selector?: (state: AccessStore) => T): T {
    const storeApi = useContext(AccessStoreContext);
    if (!storeApi) {
        throw new Error('useAccessStore must be used within <AccessStoreProvider>');
    }
    return useStore(
        storeApi,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (selector ?? ((state: AccessStore) => state)) as (state: AccessStore) => T,
    );
}

// ==============================
// 便捷 Hooks
// ==============================

export function useIsLogin(): boolean {
    const access = useAccessStore((state) => state);
    if (!access.accessToken || access.loginExpired) return false;
    if (access.accessToken.expiresAt && access.accessToken.expiresAt < Date.now()) return false;
    return true;
}

export function useAccessToken(): TokenPayload | null {
    return useAccessStore((state) => state.accessToken);
}

export function useRefreshTokenValue(): TokenPayload | null {
    return useAccessStore((state) => state.refreshToken);
}

export function usePermissions(): string[] {
    return useAccessStore((state) => state.permissions);
}
