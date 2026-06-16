import type {
  contentservicev1_DeletePageRequest,
  contentservicev1_GetPageRequest,
  contentservicev1_ListPageResponse,
  contentservicev1_Page,
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
// 页面管理
// ==============================

export function useListPages(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListPageResponse, Error>,
) {
  return useQuery({
    queryKey: ['listPages', query],
    queryFn: () => apiClient.pageService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListPages(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listPages', params],
    queryFn: () => apiClient.pageService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetPage(
  req: contentservicev1_GetPageRequest,
  options?: UseQueryOptions<contentservicev1_Page, Error>,
) {
  return useQuery({
    queryKey: ['getPage', req],
    queryFn: () => apiClient.pageService.Get(req),
    ...options,
  });
}

export function useCreatePage(
  options?: UseMutationOptions<
    contentservicev1_Page,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.pageService.Create({
        data: { ...values } as contentservicev1_Page,
      }),
    ...options,
  });
}

export function useUpdatePage(
  options?: UseMutationOptions<
    contentservicev1_Page,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.pageService.Update({
        id,
        data: { ...values } as contentservicev1_Page,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeletePage(
  options?: UseMutationOptions<
    object,
    Error,
    contentservicev1_DeletePageRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.pageService.Delete(data),
    ...options,
  });
}

// ==============================
// 页面枚举与工具函数
// ==============================

export const pageStatusList = computed(() => [
  {
    value: 'PAGE_STATUS_DRAFT',
    label: $t('enum.page.status.PAGE_STATUS_DRAFT'),
  },
  {
    value: 'PAGE_STATUS_PUBLISHED',
    label: $t('enum.page.status.PAGE_STATUS_PUBLISHED'),
  },
  {
    value: 'PAGE_STATUS_ARCHIVED',
    label: $t('enum.page.status.PAGE_STATUS_ARCHIVED'),
  },
]);

export function pageStatusToName(status: contentservicev1_Page['status']) {
  const values = pageStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const PAGE_STATUS_COLOR_MAP = {
  PAGE_STATUS_DRAFT: '#8b5cf6',
  PAGE_STATUS_PUBLISHED: '#22c55e',
  PAGE_STATUS_ARCHIVED: '#92400e',
  DEFAULT: '#94a3b8',
} as const;

export function pageStatusToColor(status: contentservicev1_Page['status']) {
  return (
    PAGE_STATUS_COLOR_MAP[status as keyof typeof PAGE_STATUS_COLOR_MAP] ||
    PAGE_STATUS_COLOR_MAP.DEFAULT
  );
}

export const pageTypeList = computed(() => [
  {
    value: 'PAGE_TYPE_DEFAULT',
    label: $t('enum.page.type.PAGE_TYPE_DEFAULT'),
  },
  {
    value: 'PAGE_TYPE_HOME',
    label: $t('enum.page.type.PAGE_TYPE_HOME'),
  },
  {
    value: 'PAGE_TYPE_ERROR_404',
    label: $t('enum.page.type.PAGE_TYPE_ERROR_404'),
  },
  {
    value: 'PAGE_TYPE_ERROR_500',
    label: $t('enum.page.type.PAGE_TYPE_ERROR_500'),
  },
  {
    value: 'PAGE_TYPE_CUSTOM',
    label: $t('enum.page.type.PAGE_TYPE_CUSTOM'),
  },
]);

export function pageTypeToName(type: contentservicev1_Page['type']) {
  const values = pageTypeList.value;
  const matchedItem = values.find((item) => item.value === type);
  return matchedItem ? matchedItem.label : '';
}

const PAGE_TYPE_COLOR_MAP = {
  PAGE_TYPE_DEFAULT: '#7c3aed',
  PAGE_TYPE_HOME: '#3b82f6',
  PAGE_TYPE_ERROR_404: '#f97316',
  PAGE_TYPE_ERROR_500: '#ef4444',
  PAGE_TYPE_CUSTOM: '#06b6d4',
  DEFAULT: '#94a3b8',
} as const;

export function pageTypeToColor(type: contentservicev1_Page['type']) {
  return (
    PAGE_TYPE_COLOR_MAP[type as keyof typeof PAGE_TYPE_COLOR_MAP] ||
    PAGE_TYPE_COLOR_MAP.DEFAULT
  );
}
