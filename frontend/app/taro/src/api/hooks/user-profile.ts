import {
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type identityservicev1_User,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';

// ==============================
// 用户资料服务封装（直接使用 apiClient）
// ==============================

/**
 * 获取当前登录用户信息
 */
export async function getMe(): Promise<identityservicev1_User | null> {
  try {
    return await apiClient.userProfileService.GetUser({});
  } catch (error) {
    console.error('getMe failed:', error);
    return null;
  }
}

// ==============================
// 获取当前用户信息 Hook
// ==============================
export function useGetUserProfile(options?: UseQueryOptions<identityservicev1_User | null, Error>) {
  return useQuery({
    queryKey: ['getMe'],
    queryFn: () => getMe(),
    ...options,
  });
}

// ==============================================
// 获取当前用户信息 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchUserProfile() {
  return queryClient.fetchQuery({
    queryKey: ['userProfile'],
    queryFn: () => getMe(),
    retry: 0,
  });
}
