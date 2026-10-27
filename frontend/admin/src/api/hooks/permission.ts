import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type permissionservicev1_Permission,
  type permissionservicev1_ListPermissionResponse,
  type permissionservicev1_GetPermissionRequest,
  type permissionservicev1_CreatePermissionRequest,
  type permissionservicev1_DeletePermissionRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 权限点管理
// ==============================

export function useListPermissions(
  query: PaginationQuery,
  options?: UseQueryOptions<permissionservicev1_ListPermissionResponse, Error>,
) {
  return useQuery({
    queryKey: ['listPermissions', query],
    queryFn: () => apiClient.permissionService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListPermissions(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listPermissions', params],
    queryFn: () => apiClient.permissionService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetPermission(
  req: permissionservicev1_GetPermissionRequest,
  options?: UseQueryOptions<permissionservicev1_Permission, Error>,
) {
  return useQuery({
    queryKey: ['getPermission', req],
    queryFn: () => apiClient.permissionService.Get(req),
    ...options,
  });
}

export function useCreatePermission(
  options?: UseMutationOptions<unknown, Error, permissionservicev1_CreatePermissionRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.permissionService.Create(data),
    ...options,
  });
}

export function useUpdatePermission(
  options?: UseMutationOptions<unknown, Error, { id: number; values: Record<string, unknown> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, unknown> }) =>
      apiClient.permissionService.Update({
        id,
        // 表单 Partial → proto 消息边界：字段集为真子集，运行时安全
          // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 表单/proto 边界
          data: { ...values } as any,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeletePermission(
  options?: UseMutationOptions<unknown, Error, permissionservicev1_DeletePermissionRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.permissionService.Delete(req),
    ...options,
  });
}

// ==============================
// 同步权限点（按菜单树重建业务权限分组与权限点）
// ==============================

export function useSyncPermissions(options?: UseMutationOptions<unknown, Error>) {
  return useMutation({
    mutationFn: () => apiClient.permissionService.SyncPermissions({}),
    ...options,
  });
}
