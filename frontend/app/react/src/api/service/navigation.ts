import {
  type siteservicev1_Navigation,
  type siteservicev1_ListNavigationResponse,
  createNavigationServiceClient,
} from '@/api/generated/app/service/v1';
import { type PaginationQuery, requestApi } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

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

// ==============================
// 导航管理 API
// ==============================

/**
 * 获取导航列表
 */
export async function listNavigations(query: PaginationQuery) {
  const params = query.toRawParams();
  const locale = currentLocaleLanguageCode();
  const formValues = query.formValues ? {...query.formValues, locale} : {locale};

  return getNavigationService().List({
    ...params,
    query: JSON.stringify(formValues),
  });
}

/**
 * 兼容旧调用方式 - 通过原始参数获取导航列表
 */
export async function listNavigationsRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<siteservicev1_ListNavigationResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return getNavigationService().List({
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
 * 获取单个导航
 */
export async function getNavigation(id: number) {
  return getNavigationService().Get({id});
}

/**
 * 创建导航
 */
export async function createNavigation(values: Partial<siteservicev1_Navigation>) {
  return getNavigationService().Create({
    data: values as siteservicev1_Navigation,
  });
}

/**
 * 更新导航
 */
export async function updateNavigation(params: {
  id: number;
  values: Partial<siteservicev1_Navigation>;
}) {
  return getNavigationService().Update({
    id: params.id,
    data: params.values as siteservicev1_Navigation,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除导航
 */
export async function deleteNavigation(id: number) {
  return getNavigationService().Delete({id});
}
