import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type identityservicev1_User,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';
import { encryptByAES } from '@/utils';

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

/**
 * 更新当前用户资料（昵称/头像/简介等白名单字段）
 */
export async function updateMe(values: Partial<identityservicev1_User>) {
  const me = await getMe();
  if (!me?.id) {
    throw new Error('User not loaded');
  }
  const resp = await apiClient.userProfileService.UpdateUser({
    id: me.id,
    data: {...me, ...values} as identityservicev1_User,
    updateMask: Object.keys(values ?? {}).join(','),
  });
  // 资料变更后使缓存失效，下次读取拿新值
  await queryClient.invalidateQueries({queryKey: ['userProfile']});
  await queryClient.invalidateQueries({queryKey: ['getMe']});
  return resp;
}

/**
 * 修改密码（需验证旧密码）。
 * 旧/新密码均需 AES 加密传输（与登录、react 端约定一致，AES_KEY 来自环境注入）。
 */
export async function changePassword(oldPassword: string, newPassword: string) {
  return apiClient.userProfileService.ChangePassword({
    oldPassword: encryptByAES(oldPassword, AES_KEY),
    newPassword: encryptByAES(newPassword, AES_KEY),
  });
}

// ==============================
// 更新资料 Hook
// ==============================
export function useUpdateMe(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof updateMe>>,
    Error,
    Partial<identityservicev1_User>
  >,
) {
  return useMutation({
    mutationFn: (values) => updateMe(values),
    ...options,
  });
}

// ==============================
// 修改密码 Hook
// ==============================
export function useChangePassword(
  options?: UseMutationOptions<
    Record<string, never>,
    Error,
    { oldPassword: string; newPassword: string }
  >,
) {
  return useMutation({
    mutationFn: (params) => changePassword(params.oldPassword, params.newPassword),
    ...options,
  });
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
