import {
  type authenticationservicev1_LoginRequest,
  type authenticationservicev1_LoginResponse,
  createAuthenticationServiceClient,
} from '@/api/generated/app/service/v1';
import { requestApi } from '@/core';

let _instance: ReturnType<typeof createAuthenticationServiceClient> | null = null;

/**
 * 获取认证服务单例（延迟初始化）
 */
export function getAuthService() {
  if (!_instance) {
    _instance = createAuthenticationServiceClient(requestApi);
  }
  return _instance;
}

/**
 * 登录
 */
export async function login(request: authenticationservicev1_LoginRequest) {
  return getAuthService().Login(request);
}

/**
 * 登出
 */
export async function logout() {
  return getAuthService().Logout({});
}

/**
 * 刷新 Token
 */
export async function refreshToken(refreshToken: string) {
  return getAuthService().RefreshToken({
    grant_type: 'refresh_token',
    refresh_token: refreshToken ?? '',
  });
}
