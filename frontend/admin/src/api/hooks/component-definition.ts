import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_ComponentDefinition,
  type contentservicev1_FieldDefinition,
  type contentservicev1_ListComponentDefinitionsResponse,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask } from '@/core/transport/rest';
import { apiClient } from '@/api/client';

// ==============================
// 组件定义（内容建模 Components）
// ==============================

export function useListComponentDefinitions(
  options?: UseQueryOptions<contentservicev1_ListComponentDefinitionsResponse, Error>,
) {
  return useQuery({
    queryKey: ['listComponentDefinitions'],
    // AGENTS 铁律：List 禁传空 sorting[]，排序走 orderBy
    queryFn: () =>
      apiClient.componentDefinitionService.List({
        orderBy: JSON.stringify([{ field: 'created_at', direction: 'DESC' }]),
      } as never),
    ...options,
  });
}

export function fetchComponentFieldDefinitions(componentId: number) {
  return apiClient.componentDefinitionService.ListFieldDefinitions({
    contentModelId: 0,
    componentId,
  });
}

export function useCreateComponentDefinition(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.componentDefinitionService.Create>>,
    Error,
    Partial<contentservicev1_ComponentDefinition>
  >,
) {
  return useMutation({
    mutationFn: (values) => apiClient.componentDefinitionService.Create({ data: values }),
    ...options,
  });
}

export function useUpdateComponentDefinition(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.componentDefinitionService.Update>>,
    Error,
    { id: number; values: Partial<contentservicev1_ComponentDefinition> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }) =>
      apiClient.componentDefinitionService.Update({
        id,
        data: values,
        // code 创建后不可改（后端校验），update_mask 不含 code
        updateMask: makeUpdateMask(
          Object.keys(values ?? {}).filter((k) => k !== 'code'),
        ),
      }),
    ...options,
  });
}

export function useDeleteComponentDefinition(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.componentDefinitionService.Delete>>,
    Error,
    { id: number }
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.componentDefinitionService.Delete(data),
    ...options,
  });
}

export function useSaveComponentFieldDefinitions(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.componentDefinitionService.SaveFieldDefinitions>>,
    Error,
    { componentId: number; fields: Array<Partial<contentservicev1_FieldDefinition>> }
  >,
) {
  return useMutation({
    mutationFn: ({ componentId, fields }) =>
      apiClient.componentDefinitionService.SaveFieldDefinitions({
        componentId,
        fields: fields as contentservicev1_FieldDefinition[],
      }),
    ...options,
  });
}
