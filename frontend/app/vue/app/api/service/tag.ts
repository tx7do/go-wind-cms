import {
  createTagServiceClient,
} from '@/api/generated/app/service/v1';
import { requestApi, type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';

let _instance: ReturnType<typeof createTagServiceClient> | null = null;

/**
 * 获取标签服务单例（延迟初始化）
 */
export function getTagService() {
  if (!_instance) {
    _instance = createTagServiceClient(requestApi);
  }
  return _instance;
}

/**
 * 查询标签列表
 */
export async function listTag(
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
  return await getTagService().List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取标签
 */
export async function getTag(id: number, locale?: string) {
  return await getTagService().Get({ id, locale });
}

/**
 * 创建标签
 */
export async function createTag(values: Record<string, any> = {}) {
  return await getTagService().Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新标签
 */
export async function updateTag(id: number, values: Record<string, any> = {}) {
  return await getTagService().Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除标签
 */
export async function deleteTag(id: number) {
  return await getTagService().Delete({ id });
}
