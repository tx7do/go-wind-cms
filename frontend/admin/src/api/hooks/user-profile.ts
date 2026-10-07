import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type identityservicev1_User,
  type identityservicev1_ChangePasswordRequest,
  type identityservicev1_BindContactRequest,
  type identityservicev1_VerifyContactRequest,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';
import { makeUpdateMask, queryClient } from '@/core';
import { encryptPassword } from '@/utils';

/**
 * 获取当前用户
 */
async function getMe(): Promise<identityservicev1_User | null> {
  return apiClient.userProfileService.GetUser({});
}

export function useGetUserProfile(options?: UseQueryOptions<identityservicev1_User | null, Error>) {
  return useQuery({
    queryKey: ['getMe'],
    queryFn: () => getMe(),
    ...options,
  });
}

// ==============================================
// 获取用户资料 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchUserProfile() {
  return queryClient.fetchQuery({
    queryKey: ['userProfile'],
    queryFn: () => getMe(),
    retry: 0,
  });
}

export function useUpdateUserProfile(
  options?: UseMutationOptions<{}, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.userProfileService.UpdateUser({
        id,
        data: { ...values } as any,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useChangePassword(
  options?: UseMutationOptions<{}, Error, identityservicev1_ChangePasswordRequest>,
) {
  return useMutation({
    // 后端 NeedDecrypt 要求 AES 密文传输（与登录同规），明文会被当密文解密导致校验必败
    mutationFn: (data) =>
      apiClient.userProfileService.ChangePassword({
        oldPassword: encryptPassword(data.oldPassword ?? ''),
        newPassword: encryptPassword(data.newPassword ?? ''),
      }),
    ...options,
  });
}

export function useBindContact(
  options?: UseMutationOptions<{}, Error, identityservicev1_BindContactRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.userProfileService.BindContact(data),
    ...options,
  });
}

export function useVerifyContact(
  options?: UseMutationOptions<{}, Error, identityservicev1_VerifyContactRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.userProfileService.VerifyContact(data),
    ...options,
  });
}
