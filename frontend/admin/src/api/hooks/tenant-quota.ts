import { useMutation, useQuery, type UseMutationOptions, type UseQueryOptions } from '@tanstack/react-query';
import { apiClient } from '@/api/client';
import { type PaginationQuery, queryClient } from '@/core';
import { makeUpdateMask } from '@/core/transport/rest/utils';
import type {
  identityservicev1_TenantQuota as TenantQuota,
  identityservicev1_TenantUsage as TenantUsage,
  identityservicev1_DeleteTenantQuotaRequest,
  identityservicev1_QuotaPlan,
  identityservicev1_ApplyQuotaPlanResponse,
  identityservicev1_TenantQuota_QuotaType,
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
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- proto 可选字段 {tenantId:0} 非完整消息
        queryFn: () => apiClient.tenantQuotaService.GetMyTenantUsage({ tenantId: 0 } as any),
    retry: 0,
    ...options,
  });
}

/** 创建配额 */
export function useCreateTenantQuota(
  options?: UseMutationOptions<void, Error, { values: Record<string, unknown> }>,
) {
  return useMutation({
    mutationFn: async ({ values }: { values: Record<string, unknown> }) => {
      await apiClient.tenantQuotaService.Create({
        // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 表单/proto 边界
        data: values as any,
      });
    },
    ...options,
  });
}

/** 更新配额（配额值） */
export function useUpdateTenantQuota(
  options?: UseMutationOptions<void, Error, { id: number; values: Record<string, unknown> }>,
) {
  return useMutation({
    mutationFn: async ({ id, values }: { id: number; values: Record<string, unknown> }) => {
      await apiClient.tenantQuotaService.Update({
        id,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 表单/proto 边界
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

// ── 配额套餐模板 ──

export function useListQuotaPlans() {
  return useQuery({
    queryKey: ['quotaPlans'],
    queryFn: () => apiClient.tenantQuotaService.ListQuotaPlans({}),
    staleTime: 30_000,
  });
}

export function useCreateQuotaPlan(
  options?: UseMutationOptions<
    identityservicev1_QuotaPlan,
    Error,
    {
      data: {
        name: string;
        description?: string;
        items: { quotaType: identityservicev1_TenantQuota_QuotaType; limit: number }[];
      };
    }
  >,
) {
  return useMutation({
    mutationFn: (req) => apiClient.tenantQuotaService.CreateQuotaPlan(req),
    ...options,
  });
}

// 删除套餐：返回值是生成层的 wellKnownEmpty（本地类型未导出），用 Awaited 推断
export function useDeleteQuotaPlan(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.tenantQuotaService.DeleteQuotaPlan>>,
    Error,
    { id: number }
  >,
) {
  return useMutation({
    mutationFn: (req) =>
      apiClient.tenantQuotaService.DeleteQuotaPlan({ id: req.id }),
    ...options,
  });
}

export function useApplyQuotaPlan(
  options?: UseMutationOptions<
    identityservicev1_ApplyQuotaPlanResponse,
    Error,
    { planId: number; tenantId: number }
  >,
) {
  return useMutation({
    mutationFn: (req) =>
      apiClient.tenantQuotaService.ApplyQuotaPlan({
        planId: req.planId,
        tenantId: req.tenantId,
      }),
    ...options,
  });
}
