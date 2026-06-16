import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_Category,
  type contentservicev1_CategoryTranslation,
  type contentservicev1_GetCategoryRequest,
  type contentservicev1_ListCategoryResponse,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';
import placeholderImage from '@/assets/images/placeholder.png';

// ==============================
// 分类服务封装（直接使用 apiClient）
// ==============================

/**
 * 兼容旧调用方式 - 通过原始参数获取分类列表
 */
export async function listCategoriesRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<contentservicev1_ListCategoryResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return apiClient.categoryService.List({
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
 * 获取单个分类
 */
export async function getCategory(request: contentservicev1_GetCategoryRequest) {
  const locale = currentLocaleLanguageCode();
  return apiClient.categoryService.Get({
    ...request,
    locale,
  });
}

/**
 * 创建分类
 */
export async function createCategory(data: Partial<contentservicev1_Category>) {
  return apiClient.categoryService.Create({
    data: {
      ...data,
      translations: data.translations ?? [],
      availableLanguages: data.availableLanguages ?? [],
      customFields: data.customFields ?? [],
      children: data.children ?? [],
    } as contentservicev1_Category,
  });
}

/**
 * 更新分类
 */
export async function updateCategory(params: {
  id: number;
  values: Partial<contentservicev1_Category>;
}) {
  return apiClient.categoryService.Update({
    id: params.id,
    data: {
      ...params.values,
      translations: params.values.translations ?? [],
      availableLanguages: params.values.availableLanguages ?? [],
      customFields: params.values.customFields ?? [],
      children: params.values.children ?? [],
    } as contentservicev1_Category,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除分类
 */
export async function deleteCategory(id: number) {
  return apiClient.categoryService.Delete({id});
}

// ==============================
// 分类列表 Hook（Mutation 形式，兼容旧调用）
// ==============================
export function useListCategories(
  options?: UseMutationOptions<
    contentservicev1_ListCategoryResponse,
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
    mutationFn: (params) => listCategoriesRaw(params),
    ...options,
  });
}

// ==============================================
// 获取分类列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListCategories(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}) {
  return queryClient.fetchQuery({
    queryKey: ['listCategories', params],
    queryFn: () => listCategoriesRaw(params),
    retry: 0,
  });
}

// ==============================
// 获取单个分类 Hook
// ==============================
export function useGetCategory(
  options?: UseMutationOptions<contentservicev1_Category, Error, contentservicev1_GetCategoryRequest>,
) {
  return useMutation({
    mutationFn: (req) => getCategory(req),
    ...options,
  });
}

// ==============================================
// 获取单个分类 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchCategory(params: contentservicev1_GetCategoryRequest & { fieldMask?: string }) {
  return queryClient.fetchQuery({
    queryKey: ['getCategory', params],
    queryFn: () => getCategory(params),
    retry: 0,
  });
}

// ==============================
// 创建分类 Hook
// ==============================
export function useCreateCategory(
  options?: UseMutationOptions<contentservicev1_Category, Error, Partial<contentservicev1_Category>>,
) {
  return useMutation({
    mutationFn: (data) => createCategory(data),
    ...options,
  });
}

// ==============================
// 更新分类 Hook
// ==============================
export function useUpdateCategory(
  options?: UseMutationOptions<
    contentservicev1_Category,
    Error,
    { id: number; values: Partial<contentservicev1_Category> }
  >,
) {
  return useMutation({
    mutationFn: (params) => updateCategory(params),
    ...options,
  });
}

// ==============================
// 删除分类 Hook
// ==============================
export function useDeleteCategory(options?: UseMutationOptions<Record<string, never>, Error, number>) {
  return useMutation({
    mutationFn: (id) => deleteCategory(id),
    ...options,
  });
}

// ==============================
// 辅助函数（纯工具，不依赖 Store）
// ==============================

/**
 * 获取分类的翻译
 */
export function getTranslation(category: contentservicev1_Category | null) {
  if (!category || !category?.translations || category.translations.length === 0) return null;
  const locale = currentLocaleLanguageCode();
  const translation = category.translations?.find(
    (t: contentservicev1_CategoryTranslation) => t.languageCode === locale
  );
  return translation || category.translations?.[0];
}

/**
 * 获取分类名称（支持国际化）
 * @param category 分类对象
 * @param t 可选的翻译函数，如果在组件外部调用且需要国际化支持，请传入 t('page.categories')
 */
export function getCategoryName(category: contentservicev1_Category | null, t?: (key: string) => string) {
  const translation = getTranslation(category);
  return translation?.name || (t ? t('category_untitled') : '未命名分类');
}

/**
 * 获取分类描述
 */
export function getCategoryDescription(category: contentservicev1_Category | null) {
  const translation = getTranslation(category);
  return translation?.description || '';
}

/**
 * 获取分类缩略图
 */
export function getCategoryThumbnail(category: contentservicev1_Category | null) {
  const translation = getTranslation(category);
  return translation?.thumbnail || placeholderImage;
}
