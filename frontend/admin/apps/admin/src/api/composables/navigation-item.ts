import type {
  siteservicev1_DeleteNavigationItemRequest,
  siteservicev1_GetNavigationItemRequest,
  siteservicev1_ListNavigationItemResponse,
  siteservicev1_NavigationItem,
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
  options?: UseMutationOptions<
    siteservicev1_NavigationItem,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.navigationItemService.Create({
        data: { ...values } as siteservicev1_NavigationItem,
      }),
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
        data: { ...values } as siteservicev1_NavigationItem,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteNavigationItem(
  options?: UseMutationOptions<
    object,
    Error,
    siteservicev1_DeleteNavigationItemRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.navigationItemService.Delete(data),
    ...options,
  });
}

// ==============================
// 导航项枚举与工具函数
// ==============================

export const navigationItemLinkTypeList = computed(() => [
  {
    value: 'LINK_TYPE_CUSTOM',
    label: $t('enum.navigationItem.linkType.LINK_TYPE_CUSTOM'),
  },
  {
    value: 'LINK_TYPE_POST',
    label: $t('enum.navigationItem.linkType.LINK_TYPE_POST'),
  },
  {
    value: 'LINK_TYPE_PAGE',
    label: $t('enum.navigationItem.linkType.LINK_TYPE_PAGE'),
  },
  {
    value: 'LINK_TYPE_CATEGORY',
    label: $t('enum.navigationItem.linkType.LINK_TYPE_CATEGORY'),
  },
  {
    value: 'LINK_TYPE_EXTERNAL',
    label: $t('enum.navigationItem.linkType.LINK_TYPE_EXTERNAL'),
  },
]);

export function navigationItemLinkTypeToName(
  linkType: siteservicev1_NavigationItem['linkType'],
) {
  const values = navigationItemLinkTypeList.value;
  const matchedItem = values.find((item) => item.value === linkType);
  return matchedItem ? matchedItem.label : '';
}

const NAVIGATION_ITEM_LINK_TYPE_COLOR_MAP = {
  LINK_TYPE_CUSTOM: '#6366f1',
  LINK_TYPE_POST: '#059669',
  LINK_TYPE_PAGE: '#d97706',
  LINK_TYPE_CATEGORY: '#dc2626',
  LINK_TYPE_EXTERNAL: '#dc2626',
  DEFAULT: '#94a3b8',
} as const;

export function navigationItemLinkTypeToColor(
  linkType: siteservicev1_NavigationItem['linkType'],
) {
  return (
    NAVIGATION_ITEM_LINK_TYPE_COLOR_MAP[
      linkType as keyof typeof NAVIGATION_ITEM_LINK_TYPE_COLOR_MAP
    ] || NAVIGATION_ITEM_LINK_TYPE_COLOR_MAP.DEFAULT
  );
}
