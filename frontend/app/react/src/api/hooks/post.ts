import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_Post,
  type contentservicev1_PostTranslation,
  type contentservicev1_ListPostResponse,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

// ==============================
// 文章服务 API
// ==============================

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

  return apiClient.postService.List({
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
  return apiClient.postService.Get({id});
}

/**
 * 创建文章
 */
export async function createPost(values: Partial<contentservicev1_Post>) {
  return apiClient.postService.Create({
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
  return apiClient.postService.Update({
    id: params.id,
    data: params.values as contentservicev1_Post,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除文章
 */
export async function deletePost(id: number) {
  return apiClient.postService.Delete({id});
}

// ==============================
// 文章列表 Hook
// ==============================
export function useListPosts(
  options?: UseMutationOptions<
    contentservicev1_ListPostResponse,
    Error,
    {
      paging?: { page?: number; pageSize?: number };
      formValues?: object | undefined;
      fieldMask?: string | undefined;
      orderBy?: string[] | undefined;
    }
  >,
) {
  return useMutation({
    mutationFn: (params) => listPostsRaw(params),
    ...options,
  });
}

// ==============================================
// 获取文章列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListPosts(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}) {
  return queryClient.fetchQuery({
    queryKey: ['listPosts', params],
    queryFn: () => listPostsRaw(params),
    retry: 0,
  });
}

// ==============================
// 获取单个文章 Hook
// ==============================
export function useGetPost(
  options?: UseMutationOptions<contentservicev1_Post, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getPost(id),
    ...options,
  });
}

// ==============================================
// 获取单个文章 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchPost(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getPost', id],
    queryFn: () => getPost(id),
    retry: 0,
  });
}

// ==============================
// 创建文章 Hook
// ==============================
export function useCreatePost(
  options?: UseMutationOptions<contentservicev1_Post, Error, Partial<contentservicev1_Post>>,
) {
  return useMutation({
    mutationFn: (data) => createPost(data),
    ...options,
  });
}

// ==============================
// 更新文章 Hook
// ==============================
export function useUpdatePost(
  options?: UseMutationOptions<
    contentservicev1_Post,
    Error,
    { id: number; values: Partial<contentservicev1_Post> }
  >,
) {
  return useMutation({
    mutationFn: (params) => updatePost(params),
    ...options,
  });
}

// ==============================
// 删除文章 Hook
// ==============================
export function useDeletePost(options?: UseMutationOptions<Record<string, never>, Error, number>) {
  return useMutation({
    mutationFn: (id) => deletePost(id),
    ...options,
  });
}

// ==============================
// 辅助函数（纯工具，不依赖 Store）
// ==============================

/**
 * 获取帖子的翻译
 */
export function getTranslation(post: contentservicev1_Post): contentservicev1_PostTranslation | null {
  if (!post?.translations || post.translations.length === 0) return null;

  const locale = currentLocaleLanguageCode();
  const translation = post.translations?.find((t: contentservicev1_PostTranslation) => t.languageCode === locale);
  return translation || post.translations?.[0];
}

/**
 * 获取帖子标题
 */
export function getPostTitle(post: contentservicev1_Post): string {
  const translation = getTranslation(post);
  return translation?.title || '';
}

/**
 * 获取帖子摘要
 */
export function getPostSummary(post: contentservicev1_Post): string {
  const translation = getTranslation(post);
  return translation?.summary || '';
}

/**
 * 获取帖子缩略图
 */
export function getPostThumbnail(post: contentservicev1_Post): string {
  const translation = getTranslation(post);
  return translation?.thumbnail || '/placeholder.png';
}

/**
 * 获取帖子内容
 */
export function getPostContent(post: contentservicev1_Post): string {
  const translation = getTranslation(post);
  return translation?.content || '';
}
