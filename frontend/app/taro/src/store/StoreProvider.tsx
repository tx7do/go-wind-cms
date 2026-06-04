import {useMemo, type ReactNode} from 'react';

import {
    createAccessStore,
    AccessStoreContext,
} from '@/store/core/access/store';
import {
    createUserStore,
    UserStoreContext,
} from '@/store/core/user/store';
import {
    createLoadingStore,
    LoadingStoreContext,
} from '@/store/core/loading/store';
import {
    createPreferencesStore,
    PreferencesStoreContext,
} from '@/core/preferences/store';
import {RequestClient} from '@/core/transport/rest/request-client';
import {env} from '@/config/env';
import {refreshToken as apiRefreshToken} from '@/api/service/auth';
import type {IUser} from '@/store/core/user/store';
import {fetchUserProfile} from '@/api/hooks/user-profile';
import Taro from "@tarojs/taro";

/**
 * 应用 Store Provider（聚合）
 *
 * Taro + Zustand 官方推荐模式：
 * 每个客户端创建独立的 store 实例，通过 React Context 注入。
 * 避免全局单例导致的数据泄漏问题。
 *
 * 使用 useMemo 确保 store 实例在组件生命周期内稳定，
 * Provider 重新渲染时不会创建新实例。
 *
 * 同时在客户端引导阶段初始化 RequestClient 单例，
 * 通过 storeApi.getState() 桥接 Context-based stores 到 axios 拦截器 callbacks。
 */
export default function StoreProvider({children}: { children: ReactNode }) {
    const accessStore = useMemo(() => createAccessStore(), []);
    const userStore = useMemo(() => createUserStore(), []);
    const loadingStore = useMemo(() => createLoadingStore(), []);
    const preferencesStore = useMemo(() => createPreferencesStore(), []);

    // 初始化 RequestClient（仅客户端执行一次）
    useMemo(() => {
        if (typeof window === 'undefined') return;

        RequestClient.init(env.apiBaseUrl, {
            getToken: () => {
                const token = accessStore.getState().accessToken;
                return token?.value ?? null;
            },

            getLocale: () => {
                return preferencesStore.getState().preferences.app.locale;
            },

            refreshToken: async () => {
                const state = accessStore.getState();
                const refreshTokenValue = state.refreshToken?.value ?? '';
                if (!refreshTokenValue) return '';

                try {
                    const resp = await apiRefreshToken(refreshTokenValue);
                    const newToken = resp.access_token || '';

                    accessStore.getState().setAccessToken({
                        value: newToken,
                        expiresAt: resp.expires_in && resp.expires_in != 0
                            ? Date.now() + resp.expires_in * 1000
                            : undefined,
                    });

                    // 同时更新 refresh token（如果返回了新的）
                    if (resp.refresh_token) {
                        accessStore.getState().setRefreshToken({
                            value: resp.refresh_token,
                            expiresAt: resp.refresh_expires_in && resp.refresh_expires_in != 0
                                ? Date.now() + resp.refresh_expires_in * 1000
                                : undefined,
                        });
                    }

                    return newToken;
                } catch {
                    // 刷新失败，执行重新认证
                    accessStore.getState().setLoginExpired(true);
                    return '';
                }
            },

            onReAuthenticate: async (redirect?: boolean) => {
                accessStore.getState().clearTokens();
                userStore.getState().clearUser();

                if (redirect && typeof Taro !== 'undefined') {
                    Taro.redirectTo({url: '/pages/user/index'});
                }
            },

            onError: (message: string) => {
                // 避免在控制台刷屏，仅在开发环境打印
                if (env.isDev) {
                    console.error('[RequestClient Error]', message);
                }
            },
        });

        // 同步预加载用户信息（如果有 token 且未过期）
        const state = accessStore.getState();
        if (state.accessToken?.value && !state.loginExpired) {
            const expiresAt = state.accessToken.expiresAt;
            const isValid = !expiresAt || expiresAt > Date.now();
            if (isValid) {
                fetchUserProfile()
                    .then((user) => {
                        userStore.getState().setUser(user as unknown as IUser);
                    })
                    .catch(() => {
                        // 静默处理（token 可能已过期，后续请求会触发刷新或重新认证）
                    });
            }
        }
    }, [accessStore, userStore, preferencesStore]);

    return (
        <PreferencesStoreContext.Provider value={preferencesStore}>
            <AccessStoreContext.Provider value={accessStore}>
                <UserStoreContext.Provider value={userStore}>
                    <LoadingStoreContext.Provider value={loadingStore}>
                        {children}
                    </LoadingStoreContext.Provider>
                </UserStoreContext.Provider>
            </AccessStoreContext.Provider>
        </PreferencesStoreContext.Provider>
    );
}
