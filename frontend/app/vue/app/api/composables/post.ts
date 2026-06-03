import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import {
  type contentservicev1_Post,
  type contentservicev1_PostTranslation,
} from '@/api/generated/app/service/v1';
import { type Paging } from '@/core/transport/rest';
import {
  createPost,
  deletePost,
  getPost,
  listPost,
  updatePost,
} from '@/api/service/post';
import { queryClient } from '@/plugins/vue-query';
import { getCurrentLocale } from '@/utils/locale';

// 直接导出 service 层函数
export { createPost, deletePost, getPost, listPost, updatePost };

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

export function getPostThumbnail(post: contentservicev1_Post, fallback = '/placeholder.jpg') {
  const translation = getPostTranslation(post);
  return translation?.thumbnail || fallback;
}

export function getPostContent(post: contentservicev1_Post) {
  const translation = getPostTranslation(post);
  return translation?.content || '';
}
