import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_DeletePostRequest,
  type contentservicev1_GetPostRequest,
  type contentservicev1_ListPostResponse,
  type contentservicev1_Post,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 帖子管理
// ==============================

export function useListPosts(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListPostResponse, Error>,
) {
  return useQuery({
    queryKey: ['listPosts', query],
    queryFn: () => apiClient.postService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListPosts(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listPosts', params],
    queryFn: () => apiClient.postService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetPost(
  req: contentservicev1_GetPostRequest,
  options?: UseQueryOptions<contentservicev1_Post, Error>,
) {
  return useQuery({
    queryKey: ['getPost', req],
    queryFn: () => apiClient.postService.Get(req),
    ...options,
  });
}

export function useCreatePost(
  options?: UseMutationOptions<contentservicev1_Post, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.postService.Create({ data: { ...values } as contentservicev1_Post }),
    ...options,
  });
}

/**
 * 文章编辑保存掩码键：只允许携带 posts 表真实存在的列。
 * availableLanguages 是 proto 计算字段、tagIds 由仓库层按 TagIds!=nil 分支处理，
 * 带进 updateMask 会让后端生成 SET 不存在的列（如 available_languages=NULL）导致 500。
 */
export const POST_SAVE_MASK_KEYS = [
  'editorType',
  'status',
  'code',
  'isFeatured',
  'sortOrder',
  'disallowComment',
  'categoryIds',
  'thumbnail',
  'publishTime',
  'translations',
] as const;

export function useUpdatePost(
  options?: UseMutationOptions<
    contentservicev1_Post,
    Error,
    { id: number; data: Record<string, any>; maskKeys?: readonly string[] }
  >,
) {
  return useMutation({
    mutationFn: ({ id, data, maskKeys }) =>
      apiClient.postService.Update({
        id,
        data: { ...data } as contentservicev1_Post,
        updateMask: makeUpdateMask([...(maskKeys ?? POST_SAVE_MASK_KEYS)]),
      }),
    ...options,
  });
}

export function useDeletePost(
  options?: UseMutationOptions<{}, Error, contentservicev1_DeletePostRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.postService.Delete(data),
    ...options,
  });
}
