import {
  type permissionservicev1_CreateMenuRequest,
  type permissionservicev1_DeleteMenuRequest,
  type permissionservicev1_GetMenuRequest,
  type permissionservicev1_ListMenuResponse,
  type permissionservicev1_Menu,
} from '@/api/generated/admin/service/v1';
import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import { makeUpdateMask, type PaginationQuery, queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 菜单管理
// ==============================

export function useListMenus(
  query: PaginationQuery,
  options?: UseQueryOptions<permissionservicev1_ListMenuResponse, Error>,
) {
  return useQuery({
    queryKey: ['listMenus', query],
    queryFn: () => apiClient.menuService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListMenus(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listMenus', params],
    queryFn: () => apiClient.menuService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetMenu(
  req: permissionservicev1_GetMenuRequest,
  options?: UseQueryOptions<permissionservicev1_Menu, Error>,
) {
  return useQuery({
    queryKey: ['getMenu', req],
    queryFn: () => apiClient.menuService.Get(req),
    ...options,
  });
}

export function useCreateMenu(
  options?: UseMutationOptions<unknown, Error, permissionservicev1_CreateMenuRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.menuService.Create(data),
    ...options,
  });
}

export function useUpdateMenu(
  options?: UseMutationOptions<unknown, Error, { id: number; values: Record<string, unknown> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, unknown> }) =>
      apiClient.menuService.Update({
        id,
        // 表单 Partial → proto 消息边界：字段集为真子集，运行时安全
          // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 表单/proto 边界
          data: { ...values } as any,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteMenu(
  options?: UseMutationOptions<unknown, Error, permissionservicev1_DeleteMenuRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.menuService.Delete(data),
    ...options,
  });
}

// ==============================
// 同步菜单（将前端静态路由推送到后端）
// ==============================

