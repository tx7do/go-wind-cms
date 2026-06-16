import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import { type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';
import { apiClient } from '@/api/client';
import { queryClient } from '@/plugins/vue-query';

// ==============================
// Service 层方法（使用 apiClient）
// ==============================

/**
 * 查询评论列表
 */
export async function listComment(
  paging?: Paging,
  formValues?: null | object,
  fieldMask?: undefined | string,
  orderBy?: null | string[],
  options?: { isTenantUser?: boolean },
) {
  const noPaging = paging?.page === undefined && paging?.pageSize === undefined;
  // @ts-ignore proto generated code is error.
  return await apiClient.commentService.List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(formValues, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取评论
 */
export async function getComment(id: number) {
  return await apiClient.commentService.Get({ id });
}

/**
 * 创建评论
 */
export async function createComment(values: Record<string, any> = {}) {
  return await apiClient.commentService.Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新评论
 */
export async function updateComment(id: number, values: Record<string, any> = {}) {
  return await apiClient.commentService.Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除评论
 */
export async function deleteComment(id: number) {
  return await apiClient.commentService.Delete({ id });
}

/** 列表查询参数 */
export interface ListCommentParams {
  paging?: Paging;
  formValues?: null | object;
  fieldMask?: string;
  orderBy?: null | string[];
  isTenantUser?: boolean;
}

// ==============================
// 列表 / 详情查询
// ==============================
export function useListComment(
  options?: UseMutationOptions<any, Error, ListCommentParams>,
) {
  return useMutation({
    mutationFn: (params) => {
      return listComment(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser },
      );
    },
    ...options,
  });
}

export async function fetchListComment(params: ListCommentParams) {
  return queryClient.fetchQuery({
    queryKey: ['listComment', params],
    queryFn: () =>
      listComment(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser },
      ),
    retry: 0,
  });
}

export function useGetComment(
  options?: UseMutationOptions<any, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getComment(id),
    ...options,
  });
}

export async function fetchComment(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getComment', id],
    queryFn: () => getComment(id),
    retry: 0,
  });
}

// ==============================
// 创建 / 更新 / 删除
// ==============================
export function useCreateComment(
  options?: UseMutationOptions<{}, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) => createComment(values),
    ...options,
  });
}

export function useUpdateComment(
  options?: UseMutationOptions<{}, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }) => updateComment(id, values),
    ...options,
  });
}

export function useDeleteComment(
  options?: UseMutationOptions<{}, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => deleteComment(id),
    ...options,
  });
}
