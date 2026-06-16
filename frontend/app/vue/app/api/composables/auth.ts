import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import type {
  authenticationservicev1_LoginRequest,
  authenticationservicev1_LoginResponse,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/plugins/vue-query';

// ==============================
// Service 层方法（使用 apiClient）
// ==============================

/**
 * 登录
 */
export async function login(request: authenticationservicev1_LoginRequest): Promise<authenticationservicev1_LoginResponse> {
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
export async function refreshToken(refreshTokenValue: string) {
  return apiClient.authenticationService.RefreshToken({
    grant_type: 'password',
    refresh_token: refreshTokenValue ?? '',
  });
}

// ==============================
// 登录（Mutation）
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

/**
 * 登录【给 Store / 外部调用】不带 Hook 的方法
 */
export async function fetchLogin(request: authenticationservicev1_LoginRequest) {
  return queryClient.fetchQuery({
    queryKey: ['login', request],
    queryFn: () => login(request),
    retry: 0,
  });
}

// ==============================
// 登出（Mutation）
// ==============================
export function useLogout(options?: UseMutationOptions<{}, Error, void>) {
  return useMutation({
    mutationFn: () => logout(),
    ...options,
  });
}

/**
 * 登出【给 Store / 外部调用】不带 Hook 的方法
 */
export async function fetchLogout() {
  return queryClient.fetchQuery({
    queryKey: ['logout'],
    queryFn: () => logout(),
    retry: 0,
  });
}

// ==============================
// 刷新 Token（Mutation）
// ==============================
export function useRefreshToken(
  options?: UseMutationOptions<
    authenticationservicev1_LoginResponse,
    Error,
    string
  >,
) {
  return useMutation({
    mutationFn: (token) => refreshToken(token),
    ...options,
  });
}

/**
 * 刷新 Token【给 Store / 外部调用】不带 Hook 的方法
 */
export async function fetchRefreshToken(refreshTokenValue: string) {
  return queryClient.fetchQuery({
    queryKey: ['refreshToken', refreshTokenValue],
    queryFn: () => refreshToken(refreshTokenValue),
    retry: 0,
  });
}
