import {
  createPageServiceClient,
} from '@/api/generated/app/service/v1';
import { requestApi, type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';

let _instance: ReturnType<typeof createPageServiceClient> | null = null;

/**
 * 获取页面服务单例（延迟初始化）
 */
export function getPageService() {
  if (!_instance) {
    _instance = createPageServiceClient(requestApi);
  }
  return _instance;
}

/**
 * 查询页面列表
 */
export async function listPage(
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
  return await getPageService().List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取页面
 */
export async function getPage(id: number) {
  return await getPageService().Get({ id });
}

/**
 * 创建页面
 */
export async function createPage(values: Record<string, any> = {}) {
  return await getPageService().Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新页面
 */
export async function updatePage(id: number, values: Record<string, any> = {}) {
  return await getPageService().Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除页面
 */
export async function deletePage(id: number) {
  return await getPageService().Delete({ id });
}
