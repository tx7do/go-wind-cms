import type {
  contentservicev1_Category,
  contentservicev1_DeleteCategoryRequest,
  contentservicev1_GetCategoryRequest,
  contentservicev1_ListCategoryResponse,
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
// 分类管理
// ==============================

export function useListCategories(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListCategoryResponse, Error>,
) {
  return useQuery({
    queryKey: ['listCategories', query],
    queryFn: () => apiClient.categoryService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListCategories(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listCategories', params],
    queryFn: () => apiClient.categoryService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetCategory(
  req: contentservicev1_GetCategoryRequest,
  options?: UseQueryOptions<contentservicev1_Category, Error>,
) {
  return useQuery({
    queryKey: ['getCategory', req],
    queryFn: () => apiClient.categoryService.Get(req),
    ...options,
  });
}

export function useCreateCategory(
  options?: UseMutationOptions<
    contentservicev1_Category,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.categoryService.Create({
        data: { ...values } as contentservicev1_Category,
      }),
    ...options,
  });
}

export function useUpdateCategory(
  options?: UseMutationOptions<
    contentservicev1_Category,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.categoryService.Update({
        id,
        data: { ...values } as contentservicev1_Category,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteCategory(
  options?: UseMutationOptions<
    object,
    Error,
    contentservicev1_DeleteCategoryRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.categoryService.Delete(data),
    ...options,
  });
}

// ==============================
// 分类枚举与工具函数
// ==============================

export const categoryStatusList = computed(() => [
  {
    value: 'CATEGORY_STATUS_ACTIVE',
    label: $t('enum.category.status.CATEGORY_STATUS_ACTIVE'),
  },
  {
    value: 'CATEGORY_STATUS_HIDDEN',
    label: $t('enum.category.status.CATEGORY_STATUS_HIDDEN'),
  },
  {
    value: 'CATEGORY_STATUS_ARCHIVED',
    label: $t('enum.category.status.CATEGORY_STATUS_ARCHIVED'),
  },
]);

export function categoryStatusToName(
  status: contentservicev1_Category['status'],
) {
  const values = categoryStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const CATEGORY_STATUS_COLOR_MAP = {
  CATEGORY_STATUS_ACTIVE: '#22c55e',
  CATEGORY_STATUS_HIDDEN: '#f97316',
  CATEGORY_STATUS_ARCHIVED: '#92400e',
  DEFAULT: '#94a3b8',
} as const;

export function categoryStatusToColor(
  status: contentservicev1_Category['status'],
) {
  return (
    CATEGORY_STATUS_COLOR_MAP[
      status as keyof typeof CATEGORY_STATUS_COLOR_MAP
    ] || CATEGORY_STATUS_COLOR_MAP.DEFAULT
  );
}
