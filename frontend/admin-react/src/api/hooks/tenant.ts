import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type identityservicev1_CreateTenantRequest,
  type identityservicev1_CreateTenantWithAdminUserRequest,
  type identityservicev1_DeleteTenantRequest,
  type identityservicev1_GetTenantRequest,
  type identityservicev1_ListTenantResponse,
  type identityservicev1_Tenant,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';

// ==============================
// 获取租户列表
// ==============================
export function useListTenants(
  query: PaginationQuery,
  options?: UseQueryOptions<identityservicev1_ListTenantResponse, Error>,
) {
  return useQuery({
    queryKey: ['listTenants', query],
    queryFn: () => apiClient.tenantService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListTenants(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listTenants', params],
    queryFn: () => apiClient.tenantService.List(params.toRawParams()),
    retry: 0,
  });
}

// ==============================
// 获取单个租户
// ==============================
export function useGetTenant(
  req: identityservicev1_GetTenantRequest,
  options?: UseQueryOptions<identityservicev1_Tenant, Error>,
) {
  return useQuery({
    queryKey: ['getTenant', req],
    queryFn: () => apiClient.tenantService.Get(req),
    ...options,
  });
}

// ==============================
// 创建租户
// ==============================
export function useCreateTenant(
  options?: UseMutationOptions<{}, Error, identityservicev1_CreateTenantRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.tenantService.Create(data),
    ...options,
  });
}

// ==============================
// 更新租户
// ==============================
export function useUpdateTenant(
  options?: UseMutationOptions<{}, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.tenantService.Update({
        id,
        data: {
          ...values,
        },
        updateMask: makeUpdateMask(Object.keys(values ?? [])),
      }),
    ...options,
  });
}

// ==============================
// 删除租户
// ==============================
export function useDeleteTenant(
  options?: UseMutationOptions<{}, Error, identityservicev1_DeleteTenantRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.tenantService.Delete(req),
    ...options,
  });
}


export function useCreateTenantWithAdminUser(
  options?: UseMutationOptions<{}, Error, identityservicev1_CreateTenantWithAdminUserRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.tenantService.CreateTenantWithAdminUser(req),
    ...options,
  });
}


