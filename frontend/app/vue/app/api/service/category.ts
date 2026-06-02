import {
  createCategoryServiceClient,
} from '@/api/generated/app/service/v1';
import { requestApi, type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';

let _instance: ReturnType<typeof createCategoryServiceClient> | null = null;

/**
 * 获取分类服务单例（延迟初始化）
 */
export function getCategoryService() {
  if (!_instance) {
    _instance = createCategoryServiceClient(requestApi);
  }
  return _instance;
}

/**
 * 查询分类列表
 */
export async function listCategory(
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
  return await getCategoryService().List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取分类
 */
export async function getCategory(id: number, fieldMask?: undefined | string, locale?: string) {
  return await getCategoryService().Get({
    id,
    locale,
    viewMask: fieldMask,
  });
}

/**
 * 创建分类
 */
export async function createCategory(values: Record<string, any> = {}) {
  return await getCategoryService().Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新分类
 */
export async function updateCategory(id: number, values: Record<string, any> = {}) {
  return await getCategoryService().Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除分类
 */
export async function deleteCategory(id: number) {
  return await getCategoryService().Delete({ id });
}
