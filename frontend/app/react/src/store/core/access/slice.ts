import {createSlice, PayloadAction} from '@reduxjs/toolkit';


import type {IAccessState, TokenPayload} from '../../types';
import {StorageManager} from '@/caches/storage-manager';
import {appNamespace} from "@/caches";

const storage = new StorageManager({prefix: appNamespace, storageType: 'localStorage'});

const initialState: IAccessState = {
    permissions: [],
    accessToken: (typeof window !== 'undefined') ? storage.getItem<TokenPayload>('accessToken', null) : null,
    refreshToken: (typeof window !== 'undefined') ? storage.getItem<TokenPayload>('refreshToken', null) : null,
    isAccessChecked: false,
    loginExpired: false,
};

/**
 * 设置 Token 的辅助函数
 */
function setTokenWithExpiry(state: IAccessState, token: TokenPayload | null) {
    state.accessToken = token;
    if (typeof window !== 'undefined') {
        if (token) {
            storage.setItem('accessToken', token, token.expiresAt ? token.expiresAt - Date.now() : undefined);
        } else {
            storage.removeItem('accessToken');
        }
    }
}

/**
 * 检查 Token 是否过期（未使用，但保留）
 */
function isTokenExpired(token: TokenPayload | null): boolean {
    if (!token || !token.expiresAt) {
        return false;
    }
    return token.expiresAt < Date.now();
}

const accessSlice = createSlice({
    name: 'access',
    initialState,
    reducers: {
        /**
         * 设置访问令牌
         */
        setAccessToken(state, action: PayloadAction<TokenPayload | null>) {
            setTokenWithExpiry(state, action.payload);
        },

        /**
         * 设置刷新令牌
         */
        setRefreshToken(state, action: PayloadAction<TokenPayload | null>) {
            state.refreshToken = action.payload;
            if (typeof window !== 'undefined') {
                if (action.payload) {
                    storage.setItem('refreshToken', action.payload, action.payload.expiresAt ? action.payload.expiresAt - Date.now() : undefined);
                } else {
                    storage.removeItem('refreshToken');
                }
            }
        },

        /**
         * 同时设置两个令牌
         */
        setTokens(state, action: PayloadAction<{
            accessToken: TokenPayload | null;
            refreshToken: TokenPayload | null
        }>) {
            const {accessToken, refreshToken} = action.payload;
            setTokenWithExpiry(state, accessToken);
            state.refreshToken = refreshToken;
            if (typeof window !== 'undefined') {
                if (refreshToken) {
                    storage.setItem('refreshToken', refreshToken, refreshToken.expiresAt ? refreshToken.expiresAt - Date.now() : undefined);
                } else {
                    storage.removeItem('refreshToken');
                }
            }
        },

        /**
         * 设置权限列表
         */
        setPermissions(state, action: PayloadAction<string[]>) {
            state.permissions = action.payload;
        },

        /**
         * 添加单个权限
         */
        addPermission(state, action: PayloadAction<string>) {
            if (!state.permissions.includes(action.payload)) {
                state.permissions.push(action.payload);
            }
        },

        /**
         * 移除权限
         */
        removePermission(state, action: PayloadAction<string>) {
            state.permissions = state.permissions.filter(p => p !== action.payload);
        },

        /**
         * 标记权限已检查
         */
        setIsAccessChecked(state, action: PayloadAction<boolean>) {
            state.isAccessChecked = action.payload;
        },

        /**
         * 设置登录过期状态
         */
        setLoginExpired(state, action: PayloadAction<boolean>) {
            state.loginExpired = action.payload;
        },

        /**
         * 清除访问令牌（但保留权限信息）
         */
        clearTokens(state) {
            state.accessToken = null;
            state.refreshToken = null;
            state.loginExpired = false;
            if (typeof window !== 'undefined') {
                storage.removeItem('accessToken');
                storage.removeItem('refreshToken');
            }
        },

        /**
         * 重置所有访问状态
         */
        resetAccess() {
            if (typeof window !== 'undefined') {
                storage.removeItem('accessToken');
                storage.removeItem('refreshToken');
            }
            return initialState;
        },

        /**
         * 从持久化存储恢复状态（用于 SSR 或页面刷新）
         */
        restoreAccess(state, action: PayloadAction<Partial<IAccessState>>) {
            if (typeof window !== 'undefined') {
                if (action.payload.accessToken) {
                    storage.setItem('accessToken', action.payload.accessToken, action.payload.accessToken.expiresAt ? action.payload.accessToken.expiresAt - Date.now() : undefined);
                }
                if (action.payload.refreshToken) {
                    storage.setItem('refreshToken', action.payload.refreshToken, action.payload.refreshToken.expiresAt ? action.payload.refreshToken.expiresAt - Date.now() : undefined);
                }
            }
            return {...initialState, ...action.payload};
        },
    },
});

export const {
    setAccessToken,
    setRefreshToken,
    setTokens,
    setPermissions,
    addPermission,
    removePermission,
    setIsAccessChecked,
    setLoginExpired,
    clearTokens,
    resetAccess,
    restoreAccess,
} = accessSlice.actions;

export default accessSlice.reducer;
