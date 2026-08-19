import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from '@tanstack/vue-query';
import { makeUpdateMask, omit } from '@/core/transport/rest';
import { apiClient } from '@/api/client';
import { queryClient } from '@/plugins/vue-query';

// ==============================
// Service 层方法（使用 apiClient）
// ==============================

/**
 * 获取当前用户
 */
export async function getMe() {
  try {
    return await apiClient.userProfileService.GetUser({});
  } catch (error) {
    console.error('getMe failed:', error);
    return null;
  }
}

/**
 * 更新当前用户
 */
export async function updateUser(id: number, values: Record<string, any> = {}) {
  const password = values.password ?? null;
  const cleaned = omit(values, 'password');
  const updateMask = makeUpdateMask(Object.keys(cleaned ?? []));
  return await apiClient.userProfileService.UpdateUser({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...cleaned,
    },
    password,
    // @ts-ignore proto generated code is error.
    updateMask,
  });
}

/**
 * 修改用户密码
 */
export async function changePassword(oldPassword: string, newPassword: string) {
  return await apiClient.userProfileService.ChangePassword({
    oldPassword,
    newPassword,
  });
}

/**
 * 绑定手机号
 */
export async function bindPhone(phone: string, code: string) {
  return await apiClient.userProfileService.BindContact({
    phone: { phone, code },
  });
}

/**
 * 绑定邮箱
 */
export async function bindEmail(email: string, verificationCode: string) {
  return await apiClient.userProfileService.BindContact({
    email: { email, verificationCode },
  });
}

/**
 * 验证手机号
 */
export async function verifyPhone(phone: string, code: string, verificationId?: string) {
  return await apiClient.userProfileService.VerifyContact({
    phone: { phone, code },
    verificationId,
  });
}

/**
 * 验证邮箱
 */
export async function verifyEmail(email: string, code: string, verificationId?: string) {
  return await apiClient.userProfileService.VerifyContact({
    email: { email, code },
    verificationId,
  });
}

// ==============================
// 获取当前用户（Query）
// ==============================
export function useGetMe(options?: UseQueryOptions<Awaited<ReturnType<typeof getMe>>, Error>) {
  return useQuery({
    queryKey: ['getMe'],
    queryFn: () => getMe(),
    ...options,
  });
}

/**
 * 获取当前用户【给 Store / 外部调用】不带 Hook 的方法
 */
export async function fetchMe() {
  return queryClient.fetchQuery({
    queryKey: ['getMe'],
    queryFn: () => getMe(),
    retry: 0,
  });
}

// ==============================
// 更新用户资料（Mutation）
// ==============================
export function useUpdateUser(
  options?: UseMutationOptions<
    {},
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }) => updateUser(id, values),
    ...options,
  });
}

/**
 * 更新用户资料【给 Store / 外部调用】不带 Hook 的方法
 */
export async function fetchUpdateUser(id: number, values: Record<string, any> = {}) {
  return queryClient.fetchQuery({
    queryKey: ['updateUser', id, values],
    queryFn: () => updateUser(id, values),
    retry: 0,
  });
}

// ==============================
// 修改密码（Mutation）
// ==============================
export function useChangePassword(
  options?: UseMutationOptions<
    {},
    Error,
    { oldPassword: string; newPassword: string }
  >,
) {
  return useMutation({
    mutationFn: ({ oldPassword, newPassword }) => changePassword(oldPassword, newPassword),
    ...options,
  });
}

// ==============================
// 绑定联系方式（Mutation）
// ==============================
export function useBindPhone(
  options?: UseMutationOptions<{}, Error, { phone: string; code: string }>,
) {
  return useMutation({
    mutationFn: ({ phone, code }) => bindPhone(phone, code),
    ...options,
  });
}

export function useBindEmail(
  options?: UseMutationOptions<{}, Error, { email: string; verificationCode: string }>,
) {
  return useMutation({
    mutationFn: ({ email, verificationCode }) => bindEmail(email, verificationCode),
    ...options,
  });
}

// ==============================
// 验证联系方式（Mutation）
// ==============================
export function useVerifyPhone(
  options?: UseMutationOptions<{}, Error, { phone: string; code: string; verificationId?: string }>,
) {
  return useMutation({
    mutationFn: ({ phone, code, verificationId }) => verifyPhone(phone, code, verificationId),
    ...options,
  });
}

export function useVerifyEmail(
  options?: UseMutationOptions<{}, Error, { email: string; code: string; verificationId?: string }>,
) {
  return useMutation({
    mutationFn: ({ email, code, verificationId }) => verifyEmail(email, code, verificationId),
    ...options,
  });
}
