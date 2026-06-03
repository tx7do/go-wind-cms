import {
  type contentservicev1_Post,
  type contentservicev1_ListPostResponse,
  createPostServiceClient,
} from '@/api/generated/app/service/v1';
import { type PaginationQuery, requestApi } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

let _instance: ReturnType<typeof createPostServiceClient> | null = null;

/**
 * 获取文章服务单例（延迟初始化）
 */
export function getPostService() {
  if (!_instance) {
    _instance = createPostServiceClient(requestApi);
  }
  return _instance;
}

// ==============================
// 文章管理 API
// ==============================

/**
 * 获取文章列表
 */
export async function listPosts(query: PaginationQuery) {
  const params = query.toRawParams();
  const locale = currentLocaleLanguageCode();
  const formValues = query.formValues ? {...query.formValues, locale} : {locale};

  return getPostService().List({
    ...params,
    query: JSON.stringify(formValues),
  });
}

/**
 * 兼容旧调用方式 - 通过原始参数获取文章列表
 */
export async function listPostsRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<contentservicev1_ListPostResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return getPostService().List({
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
 * 获取单个文章
 */
export async function getPost(id: number) {
  return getPostService().Get({id});
}

/**
 * 创建文章
 */
export async function createPost(values: Partial<contentservicev1_Post>) {
  return getPostService().Create({
    data: values as contentservicev1_Post,
  });
}

/**
 * 更新文章
 */
export async function updatePost(params: {
  id: number;
  values: Partial<contentservicev1_Post>;
}) {
  return getPostService().Update({
    id: params.id,
    data: params.values as contentservicev1_Post,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除文章
 */
export async function deletePost(id: number) {
  return getPostService().Delete({id});
}
