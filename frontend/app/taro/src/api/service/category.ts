import {
  type contentservicev1_Category,
  type contentservicev1_GetCategoryRequest,
  type contentservicev1_ListCategoryResponse,
  createCategoryServiceClient,
} from '@/api/generated/app/service/v1';
import { type PaginationQuery, requestApi } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

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

// ==============================
// 分类管理 API
// ==============================

/**
 * 获取分类列表
 */
export async function listCategories(query: PaginationQuery) {
  const params = query.toRawParams();
  const locale = currentLocaleLanguageCode();
  const formValues = query.formValues ? {...query.formValues, locale} : {locale};

  return getCategoryService().List({
    ...params,
    query: JSON.stringify(formValues),
  });
}

/**
 * 兼容旧调用方式 - 通过原始参数获取分类列表
 */
export async function listCategoriesRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<contentservicev1_ListCategoryResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return getCategoryService().List({
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
 * 获取单个分类
 */
export async function getCategory(request: contentservicev1_GetCategoryRequest) {
  const locale = currentLocaleLanguageCode();
  return getCategoryService().Get({
    ...request,
    locale,
  });
}

/**
 * 创建分类
 */
export async function createCategory(data: Partial<contentservicev1_Category>) {
  return getCategoryService().Create({
    data: {
      ...data,
      translations: data.translations ?? [],
      availableLanguages: data.availableLanguages ?? [],
      customFields: data.customFields ?? [],
      children: data.children ?? [],
    } as contentservicev1_Category,
  });
}

/**
 * 更新分类
 */
export async function updateCategory(params: {
  id: number;
  values: Partial<contentservicev1_Category>;
}) {
  return getCategoryService().Update({
    id: params.id,
    data: {
      ...params.values,
      translations: params.values.translations ?? [],
      availableLanguages: params.values.availableLanguages ?? [],
      customFields: params.values.customFields ?? [],
      children: params.values.children ?? [],
    } as contentservicev1_Category,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除分类
 */
export async function deleteCategory(id: number) {
  return getCategoryService().Delete({id});
}
