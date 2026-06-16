import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_Tag,
  type contentservicev1_TagTranslation,
  type contentservicev1_ListTagResponse,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';
import { currentLocaleLanguageCode } from '@/i18n';

// ==============================
// 标签服务封装（直接使用 apiClient）
// ==============================

/**
 * 兼容旧调用方式 - 通过原始参数获取标签列表
 */
export async function listTagsRaw(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}): Promise<contentservicev1_ListTagResponse> {
  const locale = currentLocaleLanguageCode();
  const formValues = {...(params.formValues || {}), locale};
  const noPaging =
    params.paging?.page === undefined && params.paging?.pageSize === undefined;

  return apiClient.tagService.List({
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
 * 获取单个标签
 */
export async function getTag(id: number) {
  const locale = currentLocaleLanguageCode();
  return apiClient.tagService.Get({id, locale});
}

/**
 * 创建标签
 */
export async function createTag(values: Partial<contentservicev1_Tag>) {
  return apiClient.tagService.Create({
    data: values as contentservicev1_Tag,
  });
}

/**
 * 更新标签
 */
export async function updateTag(params: {
  id: number;
  values: Partial<contentservicev1_Tag>;
}) {
  return apiClient.tagService.Update({
    id: params.id,
    data: params.values as contentservicev1_Tag,
    updateMask: Object.keys(params.values ?? {}).join(','),
  });
}

/**
 * 删除标签
 */
export async function deleteTag(id: number) {
  return apiClient.tagService.Delete({id});
}

// ==============================
// 标签列表 Hook
// ==============================
export function useListTags(
  options?: UseMutationOptions<
    contentservicev1_ListTagResponse,
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
    mutationFn: (params) => listTagsRaw(params),
    ...options,
  });
}

// ==============================================
// 获取标签列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListTags(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}) {
  return queryClient.fetchQuery({
    queryKey: ['listTags', params],
    queryFn: () => listTagsRaw(params),
    retry: 0,
  });
}

// ==============================
// 获取单个标签 Hook
// ==============================
export function useGetTag(
  options?: UseMutationOptions<contentservicev1_Tag, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getTag(id),
    ...options,
  });
}

// ==============================================
// 获取单个标签 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchTag(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getTag', id],
    queryFn: () => getTag(id),
    retry: 0,
  });
}

// ==============================
// 创建标签 Hook
// ==============================
export function useCreateTag(
  options?: UseMutationOptions<contentservicev1_Tag, Error, Partial<contentservicev1_Tag>>,
) {
  return useMutation({
    mutationFn: (data) => createTag(data),
    ...options,
  });
}

// ==============================
// 更新标签 Hook
// ==============================
export function useUpdateTag(
  options?: UseMutationOptions<
    contentservicev1_Tag,
    Error,
    { id: number; values: Partial<contentservicev1_Tag> }
  >,
) {
  return useMutation({
    mutationFn: (params) => updateTag(params),
    ...options,
  });
}

// ==============================
// 删除标签 Hook
// ==============================
export function useDeleteTag(options?: UseMutationOptions<Record<string, never>, Error, number>) {
  return useMutation({
    mutationFn: (id) => deleteTag(id),
    ...options,
  });
}

// ==============================
// 辅助函数（纯工具，不依赖 Store）
// ==============================

/**
 * 获取标签的翻译
 */
export function getTagTranslation(tag: contentservicev1_Tag | null) {
  if (!tag || !tag?.translations || tag.translations.length === 0) return null;

  const locale = currentLocaleLanguageCode();
  const translation = tag.translations?.find(
    (t: contentservicev1_TagTranslation) => t.languageCode === locale
  );
  return translation || tag.translations?.[0];
}
