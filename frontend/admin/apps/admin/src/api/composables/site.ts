import type {
  siteservicev1_DeleteSiteRequest,
  siteservicev1_GetSiteRequest,
  siteservicev1_ListSiteResponse,
  siteservicev1_Site,
} from '#/api/generated/admin/service/v1';

import { computed } from 'vue';

import { $t } from '@vben/locales';

import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/vue-query';

import { apiClient } from '#/api/client';
import { queryClient } from '#/plugins/vue-query';
import { makeUpdateMask, type PaginationQuery } from '#/transport/rest';

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
  options?: UseMutationOptions<siteservicev1_Site, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.siteService.Create({
        data: { ...values } as siteservicev1_Site,
      }),
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
        data: { ...values } as siteservicev1_Site,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteSite(
  options?: UseMutationOptions<object, Error, siteservicev1_DeleteSiteRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.siteService.Delete(data),
    ...options,
  });
}

// ==============================
// 站点枚举与工具函数
// ==============================

export const siteStatusList = computed(() => [
  {
    value: 'SITE_STATUS_ACTIVE',
    label: $t('enum.site.status.SITE_STATUS_ACTIVE'),
  },
  {
    value: 'SITE_STATUS_INACTIVE',
    label: $t('enum.site.status.SITE_STATUS_INACTIVE'),
  },
  {
    value: 'SITE_STATUS_MAINTENANCE',
    label: $t('enum.site.status.SITE_STATUS_MAINTENANCE'),
  },
]);

export function siteStatusToName(status: siteservicev1_Site['status']) {
  const values = siteStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const SITE_STATUS_COLOR_MAP = {
  SITE_STATUS_ACTIVE: '#3b82f6',
  SITE_STATUS_INACTIVE: '#60a5fa',
  SITE_STATUS_MAINTENANCE: '#22c55e',
  DEFAULT: '#94a3b8',
} as const;

export function siteStatusToColor(status: siteservicev1_Site['status']) {
  return (
    SITE_STATUS_COLOR_MAP[status as keyof typeof SITE_STATUS_COLOR_MAP] ||
    SITE_STATUS_COLOR_MAP.DEFAULT
  );
}
