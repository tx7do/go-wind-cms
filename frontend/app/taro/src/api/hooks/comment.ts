import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type commentservicev1_Comment,
  type commentservicev1_ListCommentResponse,
} from '@/api/generated/app/service/v1';
import {
  listCommentsRaw,
  getComment,
  createComment,
  updateComment,
  deleteComment,
} from '@/api/service/comment';
import { queryClient } from '@/core';

// ==============================
// 评论列表 Hook
// ==============================
export function useListComments(
  options?: UseMutationOptions<
    commentservicev1_ListCommentResponse,
    Error,
    {
      paging?: { page?: number; pageSize?: number };
      formValues?: object | undefined;
      fieldMask?: string | undefined;
      orderBy?: string[] | undefined;
    }
  >,
) {
  return useMutation({
    mutationFn: (params) => listCommentsRaw(params),
    ...options,
  });
}

// ==============================================
// 获取评论列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListComments(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}) {
  return queryClient.fetchQuery({
    queryKey: ['listComments', params],
    queryFn: () => listCommentsRaw(params),
    retry: 0,
  });
}

// ==============================
// 获取单个评论 Hook
// ==============================
export function useGetComment(
  options?: UseMutationOptions<commentservicev1_Comment, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getComment(id),
    ...options,
  });
}

// ==============================================
// 获取单个评论 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchComment(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getComment', id],
    queryFn: () => getComment(id),
    retry: 0,
  });
}

// ==============================
// 创建评论 Hook
// ==============================
export function useCreateComment(
  options?: UseMutationOptions<commentservicev1_Comment, Error, Partial<commentservicev1_Comment>>,
) {
  return useMutation({
    mutationFn: (data) => createComment(data),
    ...options,
  });
}

// ==============================
// 更新评论 Hook
// ==============================
export function useUpdateComment(
  options?: UseMutationOptions<
    commentservicev1_Comment,
    Error,
    { id: number; values: Partial<commentservicev1_Comment> }
  >,
) {
  return useMutation({
    mutationFn: (params) => updateComment(params),
    ...options,
  });
}

// ==============================
// 删除评论 Hook
// ==============================
export function useDeleteComment(options?: UseMutationOptions<Record<string, never>, Error, number>) {
  return useMutation({
    mutationFn: (id) => deleteComment(id),
    ...options,
  });
}

// ==============================================
// 创建评论 【给外部调用】不带 Hook 的方法
// ==============================================
export async function fetchCreateComment(values: Partial<commentservicev1_Comment>) {
  return queryClient.fetchQuery({
    queryKey: ['createComment', values],
    queryFn: () => createComment(values),
    retry: 0,
  });
}
