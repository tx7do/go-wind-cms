import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import {
  type contentservicev1_Tag,
  type contentservicev1_TagTranslation,
} from '@/api/generated/app/service/v1';
import { type Paging, makeOrderBy, makeQueryString, makeUpdateMask } from '@/core/transport/rest';
import { apiClient } from '@/api/client';
import { queryClient } from '@/plugins/vue-query';
import { getCurrentLocale } from '@/utils/locale';

// ==============================
// Service 层方法（使用 apiClient）
// ==============================

/**
 * 查询标签列表
 */
export async function listTag(
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
  return await apiClient.tagService.List({
    fieldMask,
    orderBy: makeOrderBy(orderBy),
    query: makeQueryString(merged, options?.isTenantUser),
    page: paging?.page,
    pageSize: paging?.pageSize,
    noPaging,
  });
}

/**
 * 获取标签
 */
export async function getTag(id: number, locale?: string) {
  return await apiClient.tagService.Get({ id, locale });
}

/**
 * 创建标签
 */
export async function createTag(values: Record<string, any> = {}) {
  return await apiClient.tagService.Create({
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
  });
}

/**
 * 更新标签
 */
export async function updateTag(id: number, values: Record<string, any> = {}) {
  return await apiClient.tagService.Update({
    id,
    // @ts-ignore proto generated code is error.
    data: {
      ...values,
    },
    updateMask: makeUpdateMask(Object.keys(values ?? [])),
  });
}

/**
 * 删除标签
 */
export async function deleteTag(id: number) {
  return await apiClient.tagService.Delete({ id });
}

/** 列表查询参数 */
export interface ListTagParams {
  paging?: Paging;
  formValues?: null | object;
  fieldMask?: string;
  orderBy?: null | string[];
  isTenantUser?: boolean;
}

// ==============================
// 列表 / 详情查询
// ==============================
export function useListTag(
  options?: UseMutationOptions<any, Error, ListTagParams>,
) {
  return useMutation({
    mutationFn: (params) => {
      const locale = getCurrentLocale();
      return listTag(
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

export async function fetchListTag(params: ListTagParams) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['listTag', params, locale],
    queryFn: () =>
      listTag(
        params.paging,
        params.formValues,
        params.fieldMask,
        params.orderBy,
        { isTenantUser: params.isTenantUser, locale }
      ),
    retry: 0,
  });
}

export function useGetTag(
  options?: UseMutationOptions<any, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => {
      const locale = getCurrentLocale();
      return getTag(id, locale);
    },
    ...options,
  });
}

export async function fetchTag(id: number) {
  const locale = getCurrentLocale();
  return queryClient.fetchQuery({
    queryKey: ['getTag', id, locale],
    queryFn: () => getTag(id, locale),
    retry: 0,
  });
}

// ==============================
// 创建 / 更新 / 删除
// ==============================
export function useCreateTag(
  options?: UseMutationOptions<{}, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) => createTag(values),
    ...options,
  });
}

export function useUpdateTag(
  options?: UseMutationOptions<{}, Error, { id: number; values: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: ({ id, values }) => updateTag(id, values),
    ...options,
  });
}

export function useDeleteTag(
  options?: UseMutationOptions<{}, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => deleteTag(id),
    ...options,
  });
}

// ==============================
// 翻译辅助方法
// ==============================
export function getTagTranslation(tag: contentservicev1_Tag) {
  if (!tag?.translations || tag.translations.length === 0) return null;

  const locale = getCurrentLocale();
  const translation = tag.translations?.find(
    (t: contentservicev1_TagTranslation) => t.languageCode === locale,
  );
  return translation || tag.translations?.[0];
}
