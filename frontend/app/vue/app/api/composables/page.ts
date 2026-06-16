import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import {
  type contentservicev1_Page,
  type contentservicev1_PageTranslation,
} from '@/api/generated/app/service/v1';
import { type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';
import { apiClient } from '@/api/client';
import { queryClient } from '@/plugins/vue-query';
import { getCurrentLocale } from '@/utils/locale';

// ==============================
// Service 层方法（使用 apiClient）
// ==============================

/**
 * 查询页面列表
 */
export async function listPage(
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
  return await apiClient.pageService.List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取页面
 */
export async function getPage(id: number) {
  return await apiClient.pageService.Get({ id });
}

/**
 * 创建页面
 */
export async function createPage(values: Record<string, any> = {}) {
  return await apiClient.pageService.Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新页面
 */
export async function updatePage(id: number, values: Record<string, any> = {}) {
  return await apiClient.pageService.Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除页面
 */
export async function deletePage(id: number) {
  return await apiClient.pageService.Delete({ id });
}

/** 列表查询参数 */
export interface ListPageParams {
  paging?: Paging;
  formValues?: null | object;
  fieldMask?: string;
  orderBy?: null | string[];
  isTenantUser?: boolean;
}

// ==============================
// 列表 / 详情查询
// ==============================
export function useListPage(
  options?: UseMutationOptions<any, Error, ListPageParams>,
) {
  return useMutation({
    mutationFn: (params) => {
      const locale = getCurrentLocale();
      return listPage(
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

export async function fetchListPage(params: ListPageParams) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['listPage', params, locale],
    queryFn: () =>
      listPage(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser, locale }
      ),
    retry: 0,
  });
}

export function useGetPage(
  options?: UseMutationOptions<any, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getPage(id),
    ...options,
  });
}

export async function fetchPage(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getPage', id],
    queryFn: () => getPage(id),
    retry: 0,
  });
}

// ==============================
// 创建 / 更新 / 删除
// ==============================
export function useCreatePage(
  options?: UseMutationOptions<{}, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) => createPage(values),
    ...options,
  });
}

export function useUpdatePage(
  options?: UseMutationOptions<{}, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }) => updatePage(id, values),
    ...options,
  });
}

export function useDeletePage(
  options?: UseMutationOptions<{}, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => deletePage(id),
    ...options,
  });
}

// ==============================
// 翻译辅助方法
// ==============================
export function getPageTranslation(page: contentservicev1_Page) {
  if (!page?.translations || page.translations.length === 0) return null;

  const locale = getCurrentLocale();
  const translation = page.translations?.find(
    (t: contentservicev1_PageTranslation) => t.languageCode === locale,
  );
  return translation || page.translations?.[0];
}
