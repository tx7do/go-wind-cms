import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type commentservicev1_ListCommentResponse,
  type commentservicev1_Comment,
  type commentservicev1_GetCommentRequest,
  type commentservicev1_CreateCommentRequest,
  type commentservicev1_DeleteCommentRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 评论管理
// ==============================

export function useListComments(
  query: PaginationQuery,
  options?: UseQueryOptions<commentservicev1_ListCommentResponse, Error>,
) {
  return useQuery({
    queryKey: ['listComments', query],
    queryFn: () => apiClient.commentService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListComments(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listComments', params],
    queryFn: () => apiClient.commentService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetComment(
  req: commentservicev1_GetCommentRequest,
  options?: UseQueryOptions<commentservicev1_Comment, Error>,
) {
  return useQuery({
    queryKey: ['getComment', req],
    queryFn: () => apiClient.commentService.Get(req),
    ...options,
  });
}

export function useCreateComment(
  options?: UseMutationOptions<commentservicev1_Comment, Error, commentservicev1_CreateCommentRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.commentService.Create(data),
    ...options,
  });
}

export function useUpdateComment(
  options?: UseMutationOptions<
    commentservicev1_Comment,
    Error,
    { id: number; values: Record<string, any>; maskKeys?: string[] }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values, maskKeys }) =>
      apiClient.commentService.Update({
        id,
        data: { ...values } as unknown as commentservicev1_Comment,
        updateMask: makeUpdateMask(maskKeys ?? Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteComment(
  options?: UseMutationOptions<{}, Error, commentservicev1_DeleteCommentRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.commentService.Delete(data),
    ...options,
  });
}
