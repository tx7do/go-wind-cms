import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import {
  type contentservicev1_Post,
  type contentservicev1_PostTranslation,
} from '@/api/generated/app/service/v1';
import { type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';
import { apiClient } from '@/api/client';
import { queryClient } from '@/plugins/vue-query';
import { getCurrentLocale } from '@/utils/locale';

// ==============================
// Service 层方法（使用 apiClient）
// ==============================

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
  return await apiClient.postService.List({
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
  return await apiClient.postService.Get({ id, locale });
}

/**
 * 创建帖子
 */
export async function createPost(values: Record<string, any> = {}) {
  return await apiClient.postService.Create({
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
  return await apiClient.postService.Update({
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
  return await apiClient.postService.Delete({ id });
}

/** 列表查询参数 */
export interface ListPostParams {
  paging?: Paging;
  formValues?: null | object;
  fieldMask?: string;
  orderBy?: null | string[];
  isTenantUser?: boolean;
}

// ==============================
// 列表 / 详情查询
// ==============================
export function useListPost(
  options?: UseMutationOptions<any, Error, ListPostParams>,
) {
  return useMutation({
    mutationFn: (params) => {
      const locale = getCurrentLocale();
      return listPost(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser, locale }
      );
    },
    ...options,
  });
}

export async function fetchListPost(params: ListPostParams) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['listPost', params, locale],
    queryFn: () =>
      listPost(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser, locale }
      ),
    retry: 0,
  });
}

export function useGetPost(
  options?: UseMutationOptions<any, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => {
      const locale = getCurrentLocale();
      return getPost(id, locale);
    },
    ...options,
  });
}

export async function fetchPost(id: number) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['getPost', id, locale],
    queryFn: () => getPost(id, locale),
    retry: 0,
  });
}

// ==============================
// 创建 / 更新 / 删除
// ==============================
export function useCreatePost(
  options?: UseMutationOptions<{}, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) => createPost(values),
    ...options,
  });
}

export function useUpdatePost(
  options?: UseMutationOptions<{}, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }) => updatePost(id, values),
    ...options,
  });
}

export function useDeletePost(
  options?: UseMutationOptions<{}, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => deletePost(id),
    ...options,
  });
}

// ==============================
// 翻译辅助方法
// ==============================

/**
 * 获取帖子的翻译
 */
export function getPostTranslation(post: contentservicev1_Post) {
  if (!post?.translations || post.translations.length === 0) return null;

  const locale = getCurrentLocale();
  // 优先查找当前语言的翻译
  const translation = post.translations?.find(
    (t: contentservicev1_PostTranslation) => t.languageCode === locale,
  );
  // 如果找不到，回退到第一个翻译
  return translation || post.translations?.[0];
}

export function getPostTitle(post: contentservicev1_Post, fallback = '') {
  const translation = getPostTranslation(post);
  return translation?.title || fallback;
}

export function getPostSummary(post: contentservicev1_Post) {
  const translation = getPostTranslation(post);
  return translation?.summary || '';
}

export function getPostThumbnail(post: contentservicev1_Post, fallback = '/placeholder.png') {
  const translation = getPostTranslation(post);
  return translation?.thumbnail || fallback;
}

export function getPostContent(post: contentservicev1_Post) {
  const translation = getPostTranslation(post);
  return translation?.content || '';
}
