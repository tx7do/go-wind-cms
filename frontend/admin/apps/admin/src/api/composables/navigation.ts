import type {
  siteservicev1_DeleteNavigationRequest,
  siteservicev1_GetNavigationRequest,
  siteservicev1_ListNavigationResponse,
  siteservicev1_Navigation,
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
    staleTime: 0,
    retry: 0,
  });
}

export function useGetNavigation(
  req: siteservicev1_GetNavigationRequest,
  options?: UseQueryOptions<siteservicev1_Navigation, Error>,
) {
  return useQuery({
    queryKey: ['getNavigation', req],
    queryFn: () => apiClient.navigationService.Get(req),
    ...options,
  });
}

export function useCreateNavigation(
  options?: UseMutationOptions<
    siteservicev1_Navigation,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.navigationService.Create({
        data: { ...values } as siteservicev1_Navigation,
      }),
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
        data: { ...values } as siteservicev1_Navigation,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteNavigation(
  options?: UseMutationOptions<
    object,
    Error,
    siteservicev1_DeleteNavigationRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.navigationService.Delete(data),
    ...options,
  });
}

// ==============================
// 导航枚举与工具函数
// ==============================

export const navigationLocationList = computed(() => [
  {
    value: 'HEADER',
    label: $t('enum.navigation.location.HEADER'),
  },
  {
    value: 'FOOTER',
    label: $t('enum.navigation.location.FOOTER'),
  },
  {
    value: 'SIDEBAR',
    label: $t('enum.navigation.location.SIDEBAR'),
  },
  {
    value: 'MOBILE',
    label: $t('enum.navigation.location.MOBILE'),
  },
  {
    value: 'TOP_BAR',
    label: $t('enum.navigation.location.TOP_BAR'),
  },
  {
    value: 'OFFCANVAS',
    label: $t('enum.navigation.location.OFFCANVAS'),
  },
]);

export function navigationLocationToName(
  location: siteservicev1_Navigation['location'],
) {
  const values = navigationLocationList.value;
  const matchedItem = values.find((item) => item.value === location);
  return matchedItem ? matchedItem.label : '';
}

const NAVIGATION_LOCATION_COLOR_MAP = {
  HEADER: '#6366f1',
  FOOTER: '#059669',
  SIDEBAR: '#d97706',
  MOBILE: '#dc2626',
  TOP_BAR: '#dc2626',
  OFFCANVAS: '#dc2626',
  DEFAULT: '#94a3b8',
} as const;

export function navigationLocationToColor(
  location: siteservicev1_Navigation['location'],
) {
  return (
    NAVIGATION_LOCATION_COLOR_MAP[
      location as keyof typeof NAVIGATION_LOCATION_COLOR_MAP
    ] || NAVIGATION_LOCATION_COLOR_MAP.DEFAULT
  );
}
