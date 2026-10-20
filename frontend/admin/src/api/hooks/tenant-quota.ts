import { useMutation, useQuery, type UseMutationOptions, type UseQueryOptions } from '@tanstack/react-query';
import { apiClient } from '@/api/client';
import { type PaginationQuery, queryClient } from '@/core';
import { makeUpdateMask } from '@/core/transport/rest/utils';
import type {
  identityservicev1_TenantQuota as TenantQuota,
  identityservicev1_TenantUsage as TenantUsage,
  identityservicev1_DeleteTenantQuotaRequest,
} from '@/api/generated/admin/service/v1';

/** 非组件上下文取列表（ProTable request） */
export async function fetchListTenantQuotas(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listTenantQuotas', params],
    queryFn: () => apiClient.tenantQuotaService.List(params.toRawParams()),
    retry: 0,
  });
}

/** 当前租户用量与配额（个人中心自助面板） */
export function useMyTenantUsage(
  options?: UseQueryOptions<TenantUsage, Error>,
) {
  return useQuery({
    queryKey: ['myTenantUsage'],
    queryFn: () => apiClient.tenantQuotaService.GetMyTenantUsage({ tenantId: 0 } as any),
    retry: 0,
    ...options,
  });
}

/** 创建配额 */
export function useCreateTenantQuota(
  options?: UseMutationOptions<void, Error, { values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: async ({ values }: { values: Record<string, any> }) => {
      await apiClient.tenantQuotaService.Create({
        data: values as any,
      });
    },
    ...options,
  });
}

/** 更新配额（配额值） */
export function useUpdateTenantQuota(
  options?: UseMutationOptions<void, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: async ({ id, values }: { id: number; values: Record<string, any> }) => {
      await apiClient.tenantQuotaService.Update({
        id,
        data: values as any,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      });
    },
    ...options,
  });
}

export function useDeleteTenantQuota(
  options?: UseMutationOptions<void, Error, identityservicev1_DeleteTenantQuotaRequest>,
) {
  return useMutation({
    mutationFn: async (req) => {
      await apiClient.tenantQuotaService.Delete(req);
    },
    ...options,
  });
}

export type { TenantQuota, TenantUsage };
