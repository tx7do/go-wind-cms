import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_Page,
  type contentservicev1_ListPageResponse,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

// ==============================
// 页面服务封装（直接使用 apiClient）
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
