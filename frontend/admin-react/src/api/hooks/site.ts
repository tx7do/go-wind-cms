import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type siteservicev1_Site,
  type siteservicev1_ListSiteResponse,
  type siteservicev1_GetSiteRequest,
  type siteservicev1_CreateSiteRequest,
  type siteservicev1_DeleteSiteRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 站点管理
// ==============================

export function useListSites(
  query: PaginationQuery,
  options?: UseQueryOptions<siteservicev1_ListSiteResponse, Error>,
) {
  return useQuery({
    queryKey: ['listSites', query],
    queryFn: () => apiClient.siteService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListSites(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listSites', params],
    queryFn: () => apiClient.siteService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetSite(
  req: siteservicev1_GetSiteRequest,
  options?: UseQueryOptions<siteservicev1_Site, Error>,
) {
  return useQuery({
    queryKey: ['getSite', req],
    queryFn: () => apiClient.siteService.Get(req),
    ...options,
  });
}

export function useCreateSite(
  options?: UseMutationOptions<siteservicev1_Site, Error, siteservicev1_CreateSiteRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.siteService.Create(data),
    ...options,
  });
}

export function useUpdateSite(
  options?: UseMutationOptions<
    siteservicev1_Site,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.siteService.Update({
        id,
        data: { ...values } as unknown as siteservicev1_Site,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteSite(
  options?: UseMutationOptions<{}, Error, siteservicev1_DeleteSiteRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.siteService.Delete(data),
    ...options,
  });
}
