import {
  useMutation,
  useQuery,
  type UseQueryOptions,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_ContentModel,
  type contentservicev1_ListContentModelResponse,
  type contentservicev1_ListFieldDefinitionsResponse,
  type contentservicev1_ListFieldDefinitionsRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 内容建模（供分类绑定 / 动态字段表单使用）
// ==============================

export function useListContentModels(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListContentModelResponse, Error>,
) {
  return useQuery({
    queryKey: ['listContentModels', query],
    queryFn: () => apiClient.contentModelService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListContentModels(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listContentModels', params],
    queryFn: () => apiClient.contentModelService.List(params.toRawParams()),
    retry: 0,
  }) as Promise<{ items: contentservicev1_ContentModel[]; total?: any }>;
}

export async function fetchListFieldDefinitions(
  req: contentservicev1_ListFieldDefinitionsRequest,
) {
  return queryClient.fetchQuery({
    queryKey: ['listFieldDefinitions', req],
    queryFn: () => apiClient.contentModelService.ListFieldDefinitions(req),
    retry: 0,
  }) as Promise<contentservicev1_ListFieldDefinitionsResponse>;
}

export function useCreateContentModel(
  options?: UseMutationOptions<contentservicev1_ContentModel, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.contentModelService.Create({
        data: {
          ...values,
          fields: values.fields ?? [],
          translations: values.translations ?? [],
        } as contentservicev1_ContentModel,
      }),
    ...options,
  });
}

export function useUpdateContentModel(
  options?: UseMutationOptions<
    contentservicev1_ContentModel,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }) =>
      apiClient.contentModelService.Update({
        id,
        data: {
          ...values,
          fields: values.fields ?? [],
          translations: values.translations ?? [],
        } as contentservicev1_ContentModel,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteContentModel(
  options?: UseMutationOptions<{}, Error, { id: number }>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.contentModelService.Delete(data),
    ...options,
  });
}

// ==============================
// 字段类型 / 关联实体枚举
// ==============================

const FIELD_TYPE_VALUES = [
  'FIELD_TYPE_TEXT',
  'FIELD_TYPE_NUMBER',
  'FIELD_TYPE_RICHTEXT',
  'FIELD_TYPE_IMAGE',
  'FIELD_TYPE_FILE',
  'FIELD_TYPE_RELATION',
] as const;

export function fieldTypeOptions(t: (key: string, options?: Record<string, any>) => string) {
  return FIELD_TYPE_VALUES.map((value) => ({
    label: t(`fieldTypeMap.${value}`),
    value,
  }));
}

export const RELATION_ENTITY_TYPES = [
  { value: 'post', label: 'Post' },
  { value: 'page', label: 'Page' },
  { value: 'category', label: 'Category' },
  { value: 'tag', label: 'Tag' },
];
