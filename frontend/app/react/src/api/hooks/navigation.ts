import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type siteservicev1_Navigation,
  type siteservicev1_NavigationItem,
  type siteservicev1_ListNavigationResponse,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

// ==============================
// 导航服务 API
// ==============================

/**
 * 兼容旧调用方式 - 通过原始参数获取导航列表
 */
export async function listNavigationsRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<siteservicev1_ListNavigationResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return apiClient.navigationService.List({
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
 * 获取单个导航
 */
export async function getNavigation(id: number) {
  return apiClient.navigationService.Get({id});
}

/**
 * 创建导航
 */
export async function createNavigation(values: Partial<siteservicev1_Navigation>) {
  return apiClient.navigationService.Create({
    data: values as siteservicev1_Navigation,
  });
}

/**
 * 更新导航
 */
export async function updateNavigation(params: {
  id: number;
  values: Partial<siteservicev1_Navigation>;
}) {
  return apiClient.navigationService.Update({
    id: params.id,
    data: params.values as siteservicev1_Navigation,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除导航
 */
export async function deleteNavigation(id: number) {
  return apiClient.navigationService.Delete({id});
}

// ==============================
// 导航列表 Hook
// ==============================
export function useListNavigations(
  options?: UseMutationOptions<
    siteservicev1_ListNavigationResponse,
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
    mutationFn: (params) => listNavigationsRaw(params),
    ...options,
  });
}

// ==============================================
// 获取导航列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListNavigations(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}) {
  return queryClient.fetchQuery({
    queryKey: ['listNavigations', params],
    queryFn: () => listNavigationsRaw(params),
    retry: 0,
  });
}

// ==============================
// 获取单个导航 Hook
// ==============================
export function useGetNavigation(
  options?: UseMutationOptions<siteservicev1_Navigation, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getNavigation(id),
    ...options,
  });
}

// ==============================================
// 获取单个导航 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchNavigation(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getNavigation', id],
    queryFn: () => getNavigation(id),
    retry: 0,
  });
}

// ==============================
// 创建导航 Hook
// ==============================
export function useCreateNavigation(
  options?: UseMutationOptions<siteservicev1_Navigation, Error, Partial<siteservicev1_Navigation>>,
) {
  return useMutation({
    mutationFn: (data) => createNavigation(data),
    ...options,
  });
}

// ==============================
// 更新导航 Hook
// ==============================
export function useUpdateNavigation(
  options?: UseMutationOptions<
    siteservicev1_Navigation,
    Error,
    { id: number; values: Partial<siteservicev1_Navigation> }
  >,
) {
  return useMutation({
    mutationFn: (params) => updateNavigation(params),
    ...options,
  });
}

// ==============================
// 删除导航 Hook
// ==============================
export function useDeleteNavigation(options?: UseMutationOptions<Record<string, never>, Error, number>) {
  return useMutation({
    mutationFn: (id) => deleteNavigation(id),
    ...options,
  });
}

// ==============================
// 辅助函数（纯工具，不依赖 Store）
// ==============================

/**
 * 递归查找导航项
 */
export function findNavItem(items: siteservicev1_NavigationItem[], key: number): siteservicev1_NavigationItem | null {
  for (const item of items) {
    if (item.id === key) return item;
    if (item.children) {
      const found = findNavItem(item.children, key);
      if (found) return found;
    }
  }
  return null;
}
