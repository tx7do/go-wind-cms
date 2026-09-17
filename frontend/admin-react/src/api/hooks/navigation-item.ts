import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type siteservicev1_ListNavigationItemResponse,
  type siteservicev1_NavigationItem,
  type siteservicev1_GetNavigationItemRequest,
  type siteservicev1_CreateNavigationItemRequest,
  type siteservicev1_DeleteNavigationItemRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 导航项管理
// ==============================

export function useListNavigationItems(
  query: PaginationQuery,
  options?: UseQueryOptions<siteservicev1_ListNavigationItemResponse, Error>,
) {
  return useQuery({
    queryKey: ['listNavigationItems', query],
    queryFn: () => apiClient.navigationItemService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListNavigationItems(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listNavigationItems', params],
    queryFn: () => apiClient.navigationItemService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetNavigationItem(
  req: siteservicev1_GetNavigationItemRequest,
  options?: UseQueryOptions<siteservicev1_NavigationItem, Error>,
) {
  return useQuery({
    queryKey: ['getNavigationItem', req],
    queryFn: () => apiClient.navigationItemService.Get(req),
    ...options,
  });
}

export function useCreateNavigationItem(
  options?: UseMutationOptions<siteservicev1_NavigationItem, Error, siteservicev1_CreateNavigationItemRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.navigationItemService.Create(data),
    ...options,
  });
}

export function useUpdateNavigationItem(
  options?: UseMutationOptions<
    siteservicev1_NavigationItem,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.navigationItemService.Update({
        id,
        data: { ...values } as unknown as siteservicev1_NavigationItem,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteNavigationItem(
  options?: UseMutationOptions<{}, Error, siteservicev1_DeleteNavigationItemRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.navigationItemService.Delete(data),
    ...options,
  });
}
