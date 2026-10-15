import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {useEffect, useState} from 'react';
import {
  type contentservicev1_Page,
  type contentservicev1_ListPageResponse,
  type contentservicev1_PageTranslation,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

// ==============================
// 页面服务 API
// ==============================

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

  return apiClient.pageService.List({
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
  return apiClient.pageService.Get({id});
}

/**
 * 按 slug（或纯数字 id）获取单个页面，附当前语言。
 * 后端对 id=0 且带 slug 的请求按 slug 解析（app BFF Get 内重写 query_by）。
 */
export async function getPageBySlug(slug: string) {
  const locale = currentLocaleLanguageCode();
  const numericId = Number(slug);
  if (Number.isInteger(numericId) && numericId > 0) {
    return apiClient.pageService.Get({id: numericId, locale});
  }
  return apiClient.pageService.Get({id: 0, slug, locale});
}

/**
 * 创建页面
 */
export async function createPage(values: Partial<contentservicev1_Page>) {
  return apiClient.pageService.Create({
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
  return apiClient.pageService.Update({
    id: params.id,
    data: params.values as contentservicev1_Page,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除页面
 */
export async function deletePage(id: number) {
  return apiClient.pageService.Delete({id});
}

// ==============================
// 页面列表 Hook
// ==============================
export function useListPages(
  options?: UseMutationOptions<
    contentservicev1_ListPageResponse,
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
    mutationFn: (params) => listPagesRaw(params),
    ...options,
  });
}

// ==============================================
// 获取页面列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListPages(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}) {
  return queryClient.fetchQuery({
    queryKey: ['listPages', params],
    queryFn: () => listPagesRaw(params),
    retry: 0,
  });
}

// ==============================
// 获取单个页面 Hook
// ==============================
export function useGetPage(
  options?: UseMutationOptions<contentservicev1_Page, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getPage(id),
    ...options,
  });
}

// ==============================================
// 获取单个页面 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchPage(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getPage', id],
    queryFn: () => getPage(id),
    retry: 0,
  });
}

// ==============================================
// 按 slug 获取单个页面 【公开页面渲染用】
// ==============================================
export async function fetchPageBySlug(slug: string) {
  return queryClient.fetchQuery({
    queryKey: ['getPageBySlug', slug, currentLocaleLanguageCode()],
    queryFn: () => getPageBySlug(slug),
    retry: 0,
  });
}

// ==============================
// 页面翻译辅助（对齐 post.getTranslation 约定：匹配当前语言 → 回退第一条）
// ==============================
export function getPageTranslation(page: contentservicev1_Page): contentservicev1_PageTranslation | null {
  if (!page?.translations || page.translations.length === 0) return null;
  const locale = currentLocaleLanguageCode();
  const translation = page.translations?.find(
    (tr: contentservicev1_PageTranslation) => tr.languageCode === locale,
  );
  return translation || page.translations?.[0];
}

export function getPageTitle(page: contentservicev1_Page): string {
  return getPageTranslation(page)?.title || '';
}

// ==============================
// CMS 页面渲染 Hook：slug 拉取已发布页面
// found=false 时调用方回退本地静态内容（公开页遇错静默，不弹登录）
// ==============================
export function useCmsPage(slug: string): {page: contentservicev1_Page | null; loading: boolean} {
  // 结果带 slug 标记：请求归属与当前 slug 不符即视为加载中，避免 effect 体内同步 setState
  const [result, setResult] = useState<{slug: string; page: contentservicev1_Page | null} | null>(null);
  const loading = result?.slug !== slug;

  useEffect(() => {
    let cancelled = false;
    const settle = (page: contentservicev1_Page | null) => {
      if (!cancelled) setResult({slug, page});
    };
    fetchPageBySlug(slug)
      .then((fetched) => settle(fetched ?? null))
      .catch((error) => {
        console.error(`Load cms page "${slug}" failed:`, error);
        settle(null);
      });
    return () => {
      cancelled = true;
    };
  }, [slug]);

  return {page: loading ? null : (result?.page ?? null), loading};
}

// ==============================
// 创建页面 Hook
// ==============================
export function useCreatePage(
  options?: UseMutationOptions<contentservicev1_Page, Error, Partial<contentservicev1_Page>>,
) {
  return useMutation({
    mutationFn: (data) => createPage(data),
    ...options,
  });
}

// ==============================
// 更新页面 Hook
// ==============================
export function useUpdatePage(
  options?: UseMutationOptions<
    contentservicev1_Page,
    Error,
    { id: number; values: Partial<contentservicev1_Page> }
  >,
) {
  return useMutation({
    mutationFn: (params) => updatePage(params),
    ...options,
  });
}

// ==============================
// 删除页面 Hook
// ==============================
export function useDeletePage(options?: UseMutationOptions<Record<string, never>, Error, number>) {
  return useMutation({
    mutationFn: (id) => deletePage(id),
    ...options,
  });
}
