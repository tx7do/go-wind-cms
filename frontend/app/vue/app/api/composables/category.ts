import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import {
  type contentservicev1_Category,
  type contentservicev1_CategoryTranslation,
} from '@/api/generated/app/service/v1';
import { type Paging } from '@/core/transport/rest';
import {
  createCategory,
  deleteCategory,
  getCategory,
  listCategory,
  updateCategory,
} from '@/api/service/category';
import { queryClient } from '@/plugins/vue-query';
import { getCurrentLocale } from '@/utils/locale';

// 直接导出 service 层函数
export { createCategory, deleteCategory, getCategory, listCategory, updateCategory };

/** 列表查询参数 */
export interface ListCategoryParams {
  paging?: Paging;
  formValues?: null | object;
  fieldMask?: string;
  orderBy?: null | string[];
  isTenantUser?: boolean;
}

// ==============================
// 列表 / 详情查询
// ==============================
export function useListCategory(
  options?: UseMutationOptions<any, Error, ListCategoryParams>,
) {
  return useMutation({
    mutationFn: (params) => {
      const locale = getCurrentLocale();
      return listCategory(
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

export async function fetchListCategory(params: ListCategoryParams) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['listCategory', params, locale],
    queryFn: () =>
      listCategory(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser, locale }
      ),
    retry: 0,
  });
}

export function useGetCategory(
  options?: UseMutationOptions<any, Error, { id: number; fieldMask?: string }>,
) {
  return useMutation({
    mutationFn: ({ id, fieldMask }) => {
      const locale = getCurrentLocale();
      return getCategory(id, fieldMask, locale);
    },
    ...options,
  });
}

export async function fetchCategory(id: number, fieldMask?: string) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['getCategory', id, fieldMask, locale],
    queryFn: () => getCategory(id, fieldMask, locale),
    retry: 0,
  });
}

// ==============================
// 创建 / 更新 / 删除
// ==============================
export function useCreateCategory(
  options?: UseMutationOptions<{}, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) => createCategory(values),
    ...options,
  });
}

export function useUpdateCategory(
  options?: UseMutationOptions<{}, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }) => updateCategory(id, values),
    ...options,
  });
}

export function useDeleteCategory(
  options?: UseMutationOptions<{}, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => deleteCategory(id),
    ...options,
  });
}

// ==============================
// 翻译辅助方法
// ==============================
export function getCategoryTranslation(category: contentservicev1_Category) {
  if (!category?.translations || category.translations.length === 0) return null;

  const locale = getCurrentLocale();
  const translation = category.translations?.find(
    (t: contentservicev1_CategoryTranslation) => t.languageCode === locale,
  );
  return translation || category.translations?.[0];
}

export function getCategoryName(category: contentservicev1_Category, fallback = '') {
  const translation = getCategoryTranslation(category);
  return translation?.name || fallback;
}

export function getCategoryDescription(category: contentservicev1_Category) {
  const translation = getCategoryTranslation(category);
  return translation?.description || '';
}

export function getCategoryThumbnail(category: contentservicev1_Category, fallback = '/placeholder.jpg') {
  const translation = getCategoryTranslation(category);
  return translation?.thumbnail || fallback;
}
