import {
  createCommentServiceClient,
} from '@/api/generated/app/service/v1';
import { requestApi, type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';

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
  return await getCommentService().List({
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
  return await getCommentService().Get({ id });
}

/**
 * 创建评论
 */
export async function createComment(values: Record<string, any> = {}) {
  return await getCommentService().Create({
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
  return await getCommentService().Update({
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
  return await getCommentService().Delete({ id });
}
