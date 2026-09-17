import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_ListTagResponse,
  type contentservicev1_Tag,
  type contentservicev1_GetTagRequest,
  type contentservicev1_CreateTagRequest,
  type contentservicev1_DeleteTagRequest,
  type translatorservicev1_TranslateRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

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
  options?: UseMutationOptions<contentservicev1_Tag, Error, contentservicev1_CreateTagRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.tagService.Create(data),
    ...options,
  });
}

export function useUpdateTag(
  options?: UseMutationOptions<
    contentservicev1_Tag,
    Error,
    { id: number; data: Record<string, any>; excludeMaskKeys?: string[] }
  >,
) {
  return useMutation({
    mutationFn: ({ id, data, excludeMaskKeys }) =>
      apiClient.tagService.Update({
        id,
        data: data as unknown as contentservicev1_Tag,
        updateMask: makeUpdateMask(
          Object.keys(data ?? {}).filter((k) => !excludeMaskKeys?.includes(k)),
        ),
      }),
    ...options,
  });
}

export function useDeleteTag(
  options?: UseMutationOptions<{}, Error, contentservicev1_DeleteTagRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.tagService.Delete(data),
    ...options,
  });
}

// ==============================
// 翻译
// ==============================

export function useTranslate(
  options?: UseMutationOptions<
    string,
    Error,
    { sourceLanguage: string; targetLanguage: string; content: string }
  >,
) {
  return useMutation({
    mutationFn: async ({ sourceLanguage, targetLanguage, content }) => {
      const req: translatorservicev1_TranslateRequest = {
        sourceLanguage,
        targetLanguage,
        content,
      };
      const resp = await apiClient.translatorService.Translate(req);
      return resp.translatedContent || content;
    },
    ...options,
  });
}
