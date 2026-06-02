import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import { type Paging } from '@/core/transport/rest';
import {
  createComment,
  deleteComment,
  getComment,
  listComment,
  updateComment,
} from '@/api/service/comment';
import { queryClient } from '@/plugins/vue-query';

// 直接导出 service 层函数
export { createComment, deleteComment, getComment, listComment, updateComment };

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
    mutationFn: (params) =>
      listComment(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser },
      ),
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
