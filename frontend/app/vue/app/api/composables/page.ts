import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import {
  type contentservicev1_Page,
  type contentservicev1_PageTranslation,
} from '@/api/generated/app/service/v1';
import { type Paging } from '@/core/transport/rest';
import {
  createPage,
  deletePage,
  getPage,
  listPage,
  updatePage,
} from '@/api/service/page';
import { queryClient } from '@/plugins/vue-query';
import { getCurrentLocale } from '@/utils/locale';

// 直接导出 service 层函数
export { createPage, deletePage, getPage, listPage, updatePage };

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
