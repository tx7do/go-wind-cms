import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type siteservicev1_ListNavigationResponse,
  type siteservicev1_Navigation,
  type siteservicev1_CreateNavigationRequest,
  type siteservicev1_DeleteNavigationRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 导航管理
// ==============================

export function useListNavigations(
  query: PaginationQuery,
  options?: UseQueryOptions<siteservicev1_ListNavigationResponse, Error>,
) {
  return useQuery({
    queryKey: ['listNavigations', query],
    queryFn: () => apiClient.navigationService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListNavigations(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listNavigations', params],
    queryFn: () => apiClient.navigationService.List(params.toRawParams()),
    retry: 0,
  });
}


export function useCreateNavigation(
  options?: UseMutationOptions<siteservicev1_Navigation, Error, siteservicev1_CreateNavigationRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.navigationService.Create(data),
    ...options,
  });
}

export function useUpdateNavigation(
  options?: UseMutationOptions<
    siteservicev1_Navigation,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.navigationService.Update({
        id,
        data: { ...values } as unknown as siteservicev1_Navigation,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteNavigation(
  options?: UseMutationOptions<{}, Error, siteservicev1_DeleteNavigationRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.navigationService.Delete(data),
    ...options,
  });
}
