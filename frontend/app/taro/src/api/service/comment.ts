import {
  type commentservicev1_Comment,
  type commentservicev1_ListCommentResponse,
  createCommentServiceClient,
} from '@/api/generated/app/service/v1';
import { type PaginationQuery, requestApi } from '@/core';

let _instance: ReturnType<typeof createCommentServiceClient> | null = null;

/**
 * 获取评论服务单例（延迟初始化）
 */
export function getCommentService() {
  if (!_instance) {
    _instance = createCommentServiceClient(requestApi);
  }
  return _instance;
}

// ==============================
// 评论管理 API
// ==============================

/**
 * 获取评论列表
 */
export async function listComments(query: PaginationQuery) {
  const params = query.toRawParams();
  return getCommentService().List(params);
}

/**
 * 兼容旧调用方式 - 通过原始参数获取评论列表
 */
export async function listCommentsRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<commentservicev1_ListCommentResponse> {
  const formValues = {...(params.formValues || {})};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return getCommentService().List({
    fieldMask: params.fieldMask,
    orderBy: params.orderBy ? JSON.stringify(params.orderBy) : undefined,
    sorting: Array.isArray(params.orderBy)
      ? params.orderBy.map((o) => ({field: o, direction: 'ASC'}))
      : undefined,
    query: JSON.stringify(formValues),
    page: params.paging?.page,
    pageSize: params.paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取单个评论
 */
export async function getComment(id: number) {
  return getCommentService().Get({id});
}

/**
 * 创建评论
 */
export async function createComment(values: Partial<commentservicev1_Comment>) {
  return getCommentService().Create({
    data: values as commentservicev1_Comment,
  });
}

/**
 * 更新评论
 */
export async function updateComment(params: {
  id: number;
  values: Partial<commentservicev1_Comment>;
}) {
  return getCommentService().Update({
    id: params.id,
    data: params.values as commentservicev1_Comment,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除评论
 */
export async function deleteComment(id: number) {
  return getCommentService().Delete({id});
}
