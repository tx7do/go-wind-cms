import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type siteservicev1_ListSiteSettingResponse,
  type siteservicev1_SiteSetting,
  type siteservicev1_GetSiteSettingRequest,
  type siteservicev1_CreateSiteSettingRequest,
  type siteservicev1_DeleteSiteSettingRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 站点配置管理
// ==============================

export function useListSiteSettings(
  query: PaginationQuery,
  options?: UseQueryOptions<siteservicev1_ListSiteSettingResponse, Error>,
) {
  return useQuery({
    queryKey: ['listSiteSettings', query],
    queryFn: () => apiClient.siteSettingService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListSiteSettings(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listSiteSettings', params],
    queryFn: () => apiClient.siteSettingService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetSiteSetting(
  req: siteservicev1_GetSiteSettingRequest,
  options?: UseQueryOptions<siteservicev1_SiteSetting, Error>,
) {
  return useQuery({
    queryKey: ['getSiteSetting', req],
    queryFn: () => apiClient.siteSettingService.Get(req),
    ...options,
  });
}

export function useCreateSiteSetting(
  options?: UseMutationOptions<siteservicev1_SiteSetting, Error, siteservicev1_CreateSiteSettingRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.siteSettingService.Create(data),
    ...options,
  });
}

export function useUpdateSiteSetting(
  options?: UseMutationOptions<
    siteservicev1_SiteSetting,
    Error,
    { id: number; values: Record<string, any>; maskKeys?: string[] }
  >,
) {
  return useMutation({
    mutationFn: ({
      id,
      values,
      maskKeys,
    }: {
      id: number;
      values: Record<string, any>;
      maskKeys?: string[];
    }) =>
      apiClient.siteSettingService.Update({
        id,
        data: { ...values } as unknown as siteservicev1_SiteSetting,
        updateMask: makeUpdateMask(maskKeys ?? Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteSiteSetting(
  options?: UseMutationOptions<{}, Error, siteservicev1_DeleteSiteSettingRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.siteSettingService.Delete(data),
    ...options,
  });
}
