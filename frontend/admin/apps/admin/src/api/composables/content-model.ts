import type {
  contentservicev1_ContentModel,
  contentservicev1_DeleteContentModelRequest,
  contentservicev1_FieldDefinition_Type,
  contentservicev1_GetContentModelRequest,
  contentservicev1_ListContentModelResponse,
  contentservicev1_ListFieldDefinitionsRequest,
  contentservicev1_ListFieldDefinitionsResponse,
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
// 内容模型管理
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
    staleTime: 0,
    retry: 0,
  });
}

export function useGetContentModel(
  req: contentservicev1_GetContentModelRequest,
  options?: UseQueryOptions<contentservicev1_ContentModel, Error>,
) {
  return useQuery({
    queryKey: ['getContentModel', req],
    queryFn: () => apiClient.contentModelService.Get(req),
    ...options,
  });
}

export function useCreateContentModel(
  options?: UseMutationOptions<
    contentservicev1_ContentModel,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.contentModelService.Create({
        data: { ...values, fields: values.fields ?? [], translations: values.translations ?? [] } as contentservicev1_ContentModel,
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
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.contentModelService.Update({
        id,
        data: { ...values, fields: values.fields ?? [], translations: values.translations ?? [] } as contentservicev1_ContentModel,
        updateMask: makeUpdateMask(Object.keys(values)),
      }),
    ...options,
  });
}

export function useDeleteContentModel(
  options?: UseMutationOptions<
    object,
    Error,
    contentservicev1_DeleteContentModelRequest
  >,
) {
  return useMutation({
    mutationFn: (req) => apiClient.contentModelService.Delete(req),
    ...options,
  });
}

// ==============================
// 字段定义查询（供内容编辑器拉动态表单 schema）
// ==============================

export function useListFieldDefinitions(
  req: contentservicev1_ListFieldDefinitionsRequest,
  options?: UseQueryOptions<contentservicev1_ListFieldDefinitionsResponse, Error>,
) {
  return useQuery({
    queryKey: ['listFieldDefinitions', req],
    queryFn: () => apiClient.contentModelService.ListFieldDefinitions(req),
    ...options,
  });
}

export async function fetchListFieldDefinitions(
  req: contentservicev1_ListFieldDefinitionsRequest,
) {
  return queryClient.fetchQuery({
    queryKey: ['listFieldDefinitions', req],
    queryFn: () => apiClient.contentModelService.ListFieldDefinitions(req),
    staleTime: 0,
    retry: 0,
  }) as Promise<contentservicev1_ListFieldDefinitionsResponse>;
}

// ==============================
// 字段类型枚举
// ==============================

export const fieldDefinitionTypeList = computed(() => [
  { value: 'FIELD_TYPE_TEXT', label: $t('enum.fieldDefinition.type.text') },
  { value: 'FIELD_TYPE_NUMBER', label: $t('enum.fieldDefinition.type.number') },
  { value: 'FIELD_TYPE_RICHTEXT', label: $t('enum.fieldDefinition.type.richtext') },
  { value: 'FIELD_TYPE_IMAGE', label: $t('enum.fieldDefinition.type.image') },
  { value: 'FIELD_TYPE_FILE', label: $t('enum.fieldDefinition.type.file') },
  { value: 'FIELD_TYPE_RELATION', label: $t('enum.fieldDefinition.type.relation') },
]);

export function fieldDefinitionTypeToName(
  type?: contentservicev1_FieldDefinition_Type,
) {
  const values = fieldDefinitionTypeList.value;
  const matchedItem = values.find((item) => item.value === type);
  return matchedItem ? matchedItem.label : '';
}

export function fieldDefinitionTypeToColor(
  type?: contentservicev1_FieldDefinition_Type,
) {
  switch (type) {
    case 'FIELD_TYPE_TEXT': {
      return 'blue';
    }
    case 'FIELD_TYPE_NUMBER': {
      return 'cyan';
    }
    case 'FIELD_TYPE_RICHTEXT': {
      return 'purple';
    }
    case 'FIELD_TYPE_IMAGE':
    case 'FIELD_TYPE_FILE': {
      return 'orange';
    }
    case 'FIELD_TYPE_RELATION': {
      return 'red';
    }
    default: {
      return 'gray';
    }
  }
}

// relation 字段允许引用的实体类型白名单
export const relationEntityTypeList = [
  { value: 'post', label: 'Post' },
  { value: 'page', label: 'Page' },
  { value: 'category', label: 'Category' },
  { value: 'tag', label: 'Tag' },
];
