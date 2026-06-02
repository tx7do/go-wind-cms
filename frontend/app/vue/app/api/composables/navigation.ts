import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import { type siteservicev1_NavigationItem } from '@/api/generated/app/service/v1';
import { type Paging } from '@/core/transport/rest';
import {
  createNavigation,
  deleteNavigation,
  getNavigation,
  listNavigation,
  updateNavigation,
} from '@/api/service/navigation';
import { queryClient } from '@/plugins/vue-query';
import { preferencesManager } from '@/core/preferences';

// 直接导出 service 层函数
export { createNavigation, deleteNavigation, getNavigation, listNavigation, updateNavigation };

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
      const locale = preferencesManager.getPreferences().app.locale;
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
  const locale = preferencesManager.getPreferences().app.locale;
  return queryClient.fetchQuery({
    queryKey: ['listNavigation', params],
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
