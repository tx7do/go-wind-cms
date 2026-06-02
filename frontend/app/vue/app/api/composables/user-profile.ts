import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from '@tanstack/vue-query';
import {
  bindEmail,
  bindPhone,
  changePassword,
  deleteAvatar,
  getMe,
  uploadAvatarBase64,
  uploadAvatarUrl,
  updateUser,
  verifyEmail,
  verifyPhone,
} from '@/api/service/user-profile';
import { queryClient } from '@/plugins/vue-query';

// 直接导出 service 层函数，供非 Vue 上下文使用
export {
  getMe,
  updateUser,
  changePassword,
  uploadAvatarBase64,
  uploadAvatarUrl,
  deleteAvatar,
  bindPhone,
  bindEmail,
  verifyPhone,
  verifyEmail,
};

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
// 头像管理（Mutation）
// ==============================
export function useUploadAvatarBase64(
  options?: UseMutationOptions<{}, Error, string>,
) {
  return useMutation({
    mutationFn: (imageBase64) => uploadAvatarBase64(imageBase64),
    ...options,
  });
}

export function useUploadAvatarUrl(
  options?: UseMutationOptions<{}, Error, string>,
) {
  return useMutation({
    mutationFn: (imageUrl) => uploadAvatarUrl(imageUrl),
    ...options,
  });
}

export function useDeleteAvatar(options?: UseMutationOptions<{}, Error, void>) {
  return useMutation({
    mutationFn: () => deleteAvatar(),
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
