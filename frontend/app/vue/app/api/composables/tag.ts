import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import {
  type contentservicev1_Tag,
  type contentservicev1_TagTranslation,
} from '@/api/generated/app/service/v1';
import { type Paging } from '@/core/transport/rest';
import {
  createTag,
  deleteTag,
  getTag,
  listTag,
  updateTag,
} from '@/api/service/tag';
import { queryClient } from '@/plugins/vue-query';
import { preferencesManager } from '@/core/preferences';

// 直接导出 service 层函数
export { createTag, deleteTag, getTag, listTag, updateTag };

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
      const locale = preferencesManager.getPreferences().app.locale;
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
  const locale = preferencesManager.getPreferences().app.locale;
  return queryClient.fetchQuery({
    queryKey: ['listTag', params],
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
    mutationFn: (id) => getTag(id),
    ...options,
  });
}

export async function fetchTag(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getTag', id],
    queryFn: () => getTag(id),
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

  const locale = preferencesManager.getPreferences().app.locale;
  const translation = tag.translations?.find(
    (t: contentservicev1_TagTranslation) => t.languageCode === locale,
  );
  return translation || tag.translations?.[0];
}
