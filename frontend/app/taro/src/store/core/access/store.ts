import {create, type StoreApi, useStore} from 'zustand';
import {persist, createJSONStorage, type StateStorage} from 'zustand/middleware';
import {createContext, useContext, type Context} from 'react';
import {getLocalStorage} from '@/core/storage/storage-adapter';
import {encryptByAES, decryptByAES} from '@/utils/crypto';

/**
 * 加密存储适配器（AUD9-M5）。
 *
 * zustand persist 的 partialize 已将持久化范围限定为 accessToken/refreshToken，
 * 即落盘的 JSON 仅含两个 token 字段。这里在 StateStorage 层对该 JSON 字符串
 * 整体做 AES 加密，使 localStorage / 小程序存储里的值为密文而非明文。
 *
 * key（storage 名）保持明文，以便 zustand 按名读写；仅 value 加密。
 * 解密失败返回 null，触发 zustand 走初始状态（迫使重新登录）。
 */
class EncryptedStateStorage implements StateStorage {
    private readonly backend: Storage;
    constructor(backend: Storage) {
        this.backend = backend;
    }
    getItem(name: string): string | null {
        const raw = this.backend.getItem(name);
        if (raw === null || !AES_KEY) return raw;
        try {
            return decryptByAES(raw, AES_KEY);
        } catch {
            return null;
        }
    }
    setItem(name: string, value: string): void {
        if (!AES_KEY) {
            this.backend.setItem(name, value);
            return;
        }
        try {
            const enc = encryptByAES(value, AES_KEY);
            this.backend.setItem(name, enc);
        } catch {
            this.backend.removeItem(name);
        }
    }
    removeItem(name: string): void {
        this.backend.removeItem(name);
    }
}

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
                name: 'gowind-access-storage', // storage key
                // 注意：zustand persist 默认使用全局 localStorage，在微信小程序等无
                // window 对象的运行时会抛 "window is not defined"。这里注入跨端适配器，
                // 浏览器端用 window.localStorage，小程序端用 Taro 同步存储 API。
                // AUD9-M5: 经 EncryptedStateStorage 包裹，落盘 token 字段整体 AES 加密。
                storage: createJSONStorage(() => new EncryptedStateStorage(getLocalStorage())),
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
