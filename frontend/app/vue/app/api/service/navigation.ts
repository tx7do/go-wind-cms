import {
  createNavigationServiceClient,
} from '@/api/generated/app/service/v1';
import { requestApi, type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';

let _instance: ReturnType<typeof createNavigationServiceClient> | null = null;

/**
 * 获取导航服务单例（延迟初始化）
 */
export function getNavigationService() {
  if (!_instance) {
    _instance = createNavigationServiceClient(requestApi);
  }
  return _instance;
}

/**
 * 查询导航列表
 */
export async function listNavigation(
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
  return await getNavigationService().List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取导航
 */
export async function getNavigation(id: number) {
  return await getNavigationService().Get({ id });
}

/**
 * 创建导航
 */
export async function createNavigation(values: Record<string, any> = {}) {
  return await getNavigationService().Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新导航
 */
export async function updateNavigation(id: number, values: Record<string, any> = {}) {
  return await getNavigationService().Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除导航
 */
export async function deleteNavigation(id: number) {
  return await getNavigationService().Delete({ id });
}
