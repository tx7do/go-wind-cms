import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type authenticationservicev1_LoginRequest,
  type authenticationservicev1_LoginResponse,
  type authenticationservicev1_GrantType,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';
import { encryptByAES } from '@/utils';
import { useAccessStore } from '@/store/core/access/store';
import { useUserStore, type IUser } from '@/store/core/user/store';
import { useI18nRouter } from '@/i18n/helpers';
import { fetchUserProfile } from '@/api/hooks/user-profile';

// ==============================
// 认证服务 API
// ==============================

/**
 * 登录
 */
export async function login(request: authenticationservicev1_LoginRequest) {
  return apiClient.authenticationService.Login(request);
}

/**
 * 登出
 */
export async function logout() {
  return apiClient.authenticationService.Logout({});
}

/**
 * 刷新 Token
 */
export async function refreshToken(refreshToken: string) {
  return apiClient.authenticationService.RefreshToken({
    grant_type: 'refresh_token',
    refresh_token: refreshToken ?? '',
  });
}

// ==============================
// 登录 Hook
// ==============================
export function useLogin(
  options?: UseMutationOptions<
    authenticationservicev1_LoginResponse,
    Error,
    authenticationservicev1_LoginRequest
  >,
) {
  return useMutation({
    mutationFn: (req) => login(req),
    ...options,
  });
}

// ==============================================
// 登录 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchLogin(request: authenticationservicev1_LoginRequest) {
  return queryClient.fetchQuery({
    queryKey: ['login', request],
    queryFn: () => login(request),
    retry: 0,
  });
}

// ==============================
// 登出 Hook
// ==============================
export function useLogout(options?: UseMutationOptions<Record<string, never>, Error, void>) {
  return useMutation({
    mutationFn: () => logout(),
    ...options,
  });
}

// ==============================================
// 登出 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchLogout() {
  return queryClient.fetchQuery({
    queryKey: ['logout'],
    queryFn: () => logout(),
    retry: 0,
  });
}

// ==============================
// 刷新 Token Hook
// ==============================
export function useRefreshToken(
  options?: UseMutationOptions<authenticationservicev1_LoginResponse, Error, string>,
) {
  return useMutation({
    mutationFn: (token) => refreshToken(token),
    ...options,
  });
}

// ==============================================
// 刷新 Token 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchRefreshToken(refreshTokenValue: string) {
  return queryClient.fetchQuery({
    queryKey: ['refreshToken', refreshTokenValue],
    queryFn: () => refreshToken(refreshTokenValue),
    retry: 0,
  });
}

// ==============================
// 认证编排 Hook（不依赖 Redux slice）
// ==============================

export interface LoginParams {
    username?: string;
    email?: string;
    mobile?: string;
    password: string;
    grant_type?: string;
}

export const DEFAULT_HOME_PATH = '/';

/**
 * 加密密码（AES）
 */
function encryptPassword(password: string): string {
    const key = process.env.NEXT_PUBLIC_AES_KEY || '';
    return encryptByAES(password, key);
}

/**
 * 认证编排 Hook
 * 提供 login / logout / refreshToken / reauthenticate 方法
 * 不再使用 Redux slice，仅协调 accessStore、userStore 和路由
 */
export function useAuth() {
    const accessStore = useAccessStore();
    const userStore = useUserStore();
    const router = useI18nRouter();

    /**
     * 异步处理用户登录操作并获取 accessToken
     */
    async function login(
        params: LoginParams,
        onSuccess?: () => Promise<void> | void,
    ) {
        let userInfo: IUser | null = null;

        // 调用登录 API
        const response = await fetchLogin({
            username: params.username,
            email: params.email,
            mobile: params.mobile,
            password: encryptPassword(params.password),
            grant_type: (params.grant_type || 'password') as authenticationservicev1_GrantType,
        });
        if (response.access_token) {
            accessStore.setAccessToken({
                value: response.access_token,
                expiresAt: response.expires_in && response.expires_in != 0
                    ? Date.now() + response.expires_in * 1000
                    : undefined,
            });
            accessStore.setRefreshToken({
                value: response.refresh_token || '',
                expiresAt: response.refresh_expires_in && response.refresh_expires_in != 0
                    ? Date.now() + response.refresh_expires_in * 1000
                    : undefined,
            });
            accessStore.setLoginExpired(false);

            userInfo = await fetchUserProfile() as unknown as IUser;
            userStore.setUser(userInfo);

            onSuccess
                ? await onSuccess?.()
                : router.push(userInfo.homePage || DEFAULT_HOME_PATH);
        }

        return {
            userInfo,
        };
    }

    async function logout(redirect: boolean = true) {
        try {
            await fetchLogout();
        } catch {
            // 不做任何处理
        } finally {
            accessStore.clearTokens();
            accessStore.setLoginExpired(false);

            // 回登录页带上当前路由地址
            router.replace(redirect ? '/login?redirect=' + encodeURIComponent(window.location.pathname) : '/login');
        }
    }

    /**
     * 重新认证
     */
    async function reauthenticate() {
        console.warn('Access token or refresh token is invalid or expired.');
        accessStore.setAccessToken(null);
        accessStore.setLoginExpired(true);
    }

    /**
     * 刷新访问令牌
     */
    async function refreshTokenAction() {
        const resp = await fetchRefreshToken(
            accessStore.refreshToken ? accessStore.refreshToken.value : '',
        );
        const newToken = resp.access_token || '';

        accessStore.setAccessToken({
            value: newToken,
            expiresAt: resp.expires_in && resp.expires_in != 0
                ? Date.now() + resp.expires_in * 1000
                : undefined,
        });

        return newToken;
    }

    return {
        login,
        logout,
        refreshToken: refreshTokenAction,
        reauthenticate,
        fetchUserProfile,
    };
}
