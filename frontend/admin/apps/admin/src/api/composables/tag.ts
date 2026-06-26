import type {
  contentservicev1_DeleteTagRequest,
  contentservicev1_GetTagRequest,
  contentservicev1_ListTagResponse,
  contentservicev1_Tag,
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
// 标签管理
// ==============================

export function useListTags(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListTagResponse, Error>,
) {
  return useQuery({
    queryKey: ['listTags', query],
    queryFn: () => apiClient.tagService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListTags(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listTags', params],
    queryFn: () => apiClient.tagService.List(params.toRawParams()),
    staleTime: 0,
    retry: 0,
  });
}

export function useGetTag(
  req: contentservicev1_GetTagRequest,
  options?: UseQueryOptions<contentservicev1_Tag, Error>,
) {
  return useQuery({
    queryKey: ['getTag', req],
    queryFn: () => apiClient.tagService.Get(req),
    ...options,
  });
}

export function useCreateTag(
  options?: UseMutationOptions<
    contentservicev1_Tag,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.tagService.Create({
        data: { ...values } as contentservicev1_Tag,
      }),
    ...options,
  });
}

export function useUpdateTag(
  options?: UseMutationOptions<
    contentservicev1_Tag,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.tagService.Update({
        id,
        data: { ...values } as contentservicev1_Tag,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteTag(
  options?: UseMutationOptions<
    object,
    Error,
    contentservicev1_DeleteTagRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.tagService.Delete(data),
    ...options,
  });
}

// ==============================
// 标签枚举与工具函数
// ==============================

export const tagStatusList = computed(() => [
  {
    value: 'TAG_STATUS_ACTIVE',
    label: $t('enum.tag.status.TAG_STATUS_ACTIVE'),
  },
  {
    value: 'TAG_STATUS_HIDDEN',
    label: $t('enum.tag.status.TAG_STATUS_HIDDEN'),
  },
  {
    value: 'TAG_STATUS_ARCHIVED',
    label: $t('enum.tag.status.TAG_STATUS_ARCHIVED'),
  },
]);

export function tagStatusToName(status: contentservicev1_Tag['status']) {
  const values = tagStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const TAG_STATUS_COLOR_MAP = {
  TAG_STATUS_ACTIVE: '#22c55e',
  TAG_STATUS_HIDDEN: '#f97316',
  TAG_STATUS_ARCHIVED: '#92400e',
  DEFAULT: '#94a3b8',
} as const;

export function tagStatusToColor(status: contentservicev1_Tag['status']) {
  return (
    TAG_STATUS_COLOR_MAP[status as keyof typeof TAG_STATUS_COLOR_MAP] ||
    TAG_STATUS_COLOR_MAP.DEFAULT
  );
}
