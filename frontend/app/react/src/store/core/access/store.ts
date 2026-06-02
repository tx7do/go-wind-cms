import {create, type StoreApi, useStore} from 'zustand';
import {createContext, useContext, type Context} from 'react';
import {StorageManager} from '@/core/storage';

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
// 持久化存储（仅客户端使用）
// ==============================

const appNamespace = 'gowind';
const storage = new StorageManager({prefix: appNamespace, storageType: 'localStorage'});

export const ACCESS_TOKEN_KEY = 'access-token';
export const REFRESH_TOKEN_KEY = 'refresh-token';

// ==============================
// 初始状态
// ==============================

function getInitialState(): AccessState {
    return {
        permissions: [],
        accessToken: (typeof window !== 'undefined') ? storage.getItem<TokenPayload>(ACCESS_TOKEN_KEY, null) : null,
        refreshToken: (typeof window !== 'undefined') ? storage.getItem<TokenPayload>(REFRESH_TOKEN_KEY, null) : null,
        isAccessChecked: false,
        loginExpired: false,
    };
}

// ==============================
// 辅助函数
// ==============================

function persistToken(key: string, token: TokenPayload | null) {
    if (typeof window === 'undefined') return;
    if (token) {
        storage.setItem(key, token, token.expiresAt ? token.expiresAt - Date.now() : undefined);
    } else {
        storage.removeItem(key);
    }
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
// Store 工厂函数（每个请求/客户端独立实例）
// ==============================

export function createAccessStore(): StoreApi<AccessStore> {
    return create<AccessStore>((set) => ({
        ...getInitialState(),

        setAccessToken: (token) => {
            persistToken(ACCESS_TOKEN_KEY, token);
            set({accessToken: token});
        },

        setRefreshToken: (token) => {
            persistToken(REFRESH_TOKEN_KEY, token);
            set({refreshToken: token});
        },

        setTokens: ({accessToken, refreshToken}) => {
            persistToken(ACCESS_TOKEN_KEY, accessToken);
            persistToken(REFRESH_TOKEN_KEY, refreshToken);
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
            persistToken(ACCESS_TOKEN_KEY, null);
            persistToken(REFRESH_TOKEN_KEY, null);
            set({accessToken: null, refreshToken: null, loginExpired: false});
        },

        resetAccess: () => {
            persistToken(ACCESS_TOKEN_KEY, null);
            persistToken(REFRESH_TOKEN_KEY, null);
            set({...getInitialState()});
        },

        restoreAccess: (payload) => {
            if (typeof window !== 'undefined') {
                if (payload.accessToken) {
                    persistToken(ACCESS_TOKEN_KEY, payload.accessToken);
                }
                if (payload.refreshToken) {
                    persistToken(REFRESH_TOKEN_KEY, payload.refreshToken);
                }
            }
            set((state) => ({...state, ...payload}));
        },
    }));
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
