import {
  type contentservicev1_Page,
  type contentservicev1_ListPageResponse,
  createPageServiceClient,
} from '@/api/generated/app/service/v1';
import { type PaginationQuery, requestApi } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

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

// ==============================
// 页面管理 API
// ==============================

/**
 * 获取页面列表
 */
export async function listPages(query: PaginationQuery) {
  const params = query.toRawParams();
  const locale = currentLocaleLanguageCode();
  const formValues = query.formValues ? {...query.formValues, locale} : {locale};

  return getPageService().List({
    ...params,
    query: JSON.stringify(formValues),
  });
}

/**
 * 兼容旧调用方式 - 通过原始参数获取页面列表
 */
export async function listPagesRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<contentservicev1_ListPageResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return getPageService().List({
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
 * 获取单个页面
 */
export async function getPage(id: number) {
  return getPageService().Get({id});
}

/**
 * 创建页面
 */
export async function createPage(values: Partial<contentservicev1_Page>) {
  return getPageService().Create({
    data: values as contentservicev1_Page,
  });
}

/**
 * 更新页面
 */
export async function updatePage(params: {
  id: number;
  values: Partial<contentservicev1_Page>;
}) {
  return getPageService().Update({
    id: params.id,
    data: params.values as contentservicev1_Page,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除页面
 */
export async function deletePage(id: number) {
  return getPageService().Delete({id});
}
