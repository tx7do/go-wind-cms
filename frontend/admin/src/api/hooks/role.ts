import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type permissionservicev1_CreateRoleRequest,
  type permissionservicev1_DeleteRoleRequest,
  type permissionservicev1_GetRoleRequest,
  type permissionservicev1_ListRoleResponse,
  type permissionservicev1_Role,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery, queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 角色管理
// ==============================

export function useListRoles(
  query: PaginationQuery,
  options?: UseQueryOptions<permissionservicev1_ListRoleResponse, Error>,
) {
  return useQuery({
    queryKey: ['listRoles', query],
    queryFn: () => apiClient.roleService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListRoles(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listRoles', params],
    queryFn: () => apiClient.roleService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetRole(
  req: permissionservicev1_GetRoleRequest,
  options?: UseQueryOptions<permissionservicev1_Role, Error>,
) {
  return useQuery({
    queryKey: ['getRole', req],
    queryFn: () => apiClient.roleService.Get(req),
    ...options,
  });
}

export function useCreateRole(
  options?: UseMutationOptions<unknown, Error, permissionservicev1_CreateRoleRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.roleService.Create(data),
    ...options,
  });
}

export function useUpdateRole(
  options?: UseMutationOptions<unknown, Error, { id: number; values: Record<string, unknown> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, unknown> }) =>
      apiClient.roleService.Update({
        id,
        // 表单 Partial → proto 消息边界：字段集为真子集，运行时安全
          // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 表单/proto 边界
          data: { ...values } as any,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteRole(
  options?: UseMutationOptions<unknown, Error, permissionservicev1_DeleteRoleRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.roleService.Delete(req),
    ...options,
  });
}
