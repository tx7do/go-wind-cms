import {
  createPostServiceClient,
} from '@/api/generated/app/service/v1';
import { requestApi, type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';

let _instance: ReturnType<typeof createPostServiceClient> | null = null;

/**
 * 获取帖子服务单例（延迟初始化）
 */
export function getPostService() {
  if (!_instance) {
    _instance = createPostServiceClient(requestApi);
  }
  return _instance;
}

/**
 * 查询帖子列表
 */
export async function listPost(
  paging?: Paging,
  formValues?: null | object,
  fieldMask?: undefined | string,
  orderBy?: null | string[],
  options?: { isTenantUser?: boolean; locale?: string },
) {
  const merged: Record<string, any> = {
    ...(formValues ?? {}),
    locale: options?.locale,
  };

  const noPaging = paging?.page === undefined && paging?.pageSize === undefined;
  // @ts-ignore proto generated code is error.
  return await getPostService().List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取帖子
 */
export async function getPost(id: number, locale?: string) {
  if (!id) return null;
  return await getPostService().Get({ id, locale });
}

/**
 * 创建帖子
 */
export async function createPost(values: Record<string, any> = {}) {
  return await getPostService().Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新帖子
 */
export async function updatePost(id: number, values: Record<string, any> = {}) {
  return await getPostService().Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除帖子
 */
export async function deletePost(id: number) {
  return await getPostService().Delete({ id });
}
