import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import { type siteservicev1_NavigationItem } from '@/api/generated/app/service/v1';
import { type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';
import { apiClient } from '@/api/client';
import { queryClient } from '@/plugins/vue-query';
import { getCurrentLocale } from '@/utils/locale';

// ==============================
// Service 层方法（使用 apiClient）
// ==============================

/**
 * 查询导航列表
 */
export async function listNavigation(
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
  return await apiClient.navigationService.List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取导航
 */
export async function getNavigation(id: number) {
  return await apiClient.navigationService.Get({ id });
}

/**
 * 创建导航
 */
export async function createNavigation(values: Record<string, any> = {}) {
  return await apiClient.navigationService.Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新导航
 */
export async function updateNavigation(id: number, values: Record<string, any> = {}) {
  return await apiClient.navigationService.Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除导航
 */
export async function deleteNavigation(id: number) {
  return await apiClient.navigationService.Delete({ id });
}

/** 列表查询参数 */
export interface ListNavigationParams {
  paging?: Paging;
  formValues?: null | object;
  fieldMask?: string;
  orderBy?: null | string[];
  isTenantUser?: boolean;
}

// ==============================
// 列表 / 详情查询
// ==============================
export function useListNavigation(
  options?: UseMutationOptions<any, Error, ListNavigationParams>,
) {
  return useMutation({
    mutationFn: (params) => {
      const locale = getCurrentLocale();
      return listNavigation(
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

export async function fetchListNavigation(params: ListNavigationParams) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['listNavigation', params, locale],
    queryFn: () =>
      listNavigation(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser, locale }
      ),
    retry: 0,
  });
}

export function useGetNavigation(
  options?: UseMutationOptions<any, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getNavigation(id),
    ...options,
  });
}

export async function fetchNavigation(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getNavigation', id],
    queryFn: () => getNavigation(id),
    retry: 0,
  });
}

// ==============================
// 创建 / 更新 / 删除
// ==============================
export function useCreateNavigation(
  options?: UseMutationOptions<{}, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) => createNavigation(values),
    ...options,
  });
}

export function useUpdateNavigation(
  options?: UseMutationOptions<{}, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }) => updateNavigation(id, values),
    ...options,
  });
}

export function useDeleteNavigation(
  options?: UseMutationOptions<{}, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => deleteNavigation(id),
    ...options,
  });
}

// ==============================
// 工具函数
// ==============================

/**
 * 递归查找导航项
 */
export function findNavItem(
  items: siteservicev1_NavigationItem[],
  key: number,
): siteservicev1_NavigationItem | null {
  for (const item of items) {
    if (item.id === key) return item;
    if (item.children) {
      const found = findNavItem(item.children, key);
      if (found) return found;
    }
  }
  return null;
}
