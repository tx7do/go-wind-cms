import {
  type contentservicev1_Tag,
  type contentservicev1_ListTagResponse,
  createTagServiceClient,
} from '@/api/generated/app/service/v1';
import { type PaginationQuery, requestApi } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

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

// ==============================
// 标签管理 API
// ==============================

/**
 * 获取标签列表
 */
export async function listTags(query: PaginationQuery) {
  const params = query.toRawParams();
  const locale = currentLocaleLanguageCode();
  const formValues = query.formValues ? {...query.formValues, locale} : {locale};

  return getTagService().List({
    ...params,
    query: JSON.stringify(formValues),
  });
}

/**
 * 兼容旧调用方式 - 通过原始参数获取标签列表
 */
export async function listTagsRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<contentservicev1_ListTagResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return getTagService().List({
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
 * 获取单个标签
 */
export async function getTag(id: number) {
  return getTagService().Get({id});
}

/**
 * 创建标签
 */
export async function createTag(values: Partial<contentservicev1_Tag>) {
  return getTagService().Create({
    data: values as contentservicev1_Tag,
  });
}

/**
 * 更新标签
 */
export async function updateTag(params: {
  id: number;
  values: Partial<contentservicev1_Tag>;
}) {
  return getTagService().Update({
    id: params.id,
    data: params.values as contentservicev1_Tag,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除标签
 */
export async function deleteTag(id: number) {
  return getTagService().Delete({id});
}
