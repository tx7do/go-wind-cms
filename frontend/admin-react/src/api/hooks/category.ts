import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_Category,
  type contentservicev1_DeleteCategoryRequest,
  type contentservicev1_GetCategoryRequest,
  type contentservicev1_ListCategoryResponse,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 分类管理
// ==============================

export function useListCategories(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListCategoryResponse, Error>,
) {
  return useQuery({
    queryKey: ['listCategories', query],
    queryFn: () => apiClient.categoryService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListCategories(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listCategories', params],
    queryFn: () => apiClient.categoryService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetCategory(
  req: contentservicev1_GetCategoryRequest,
  options?: UseQueryOptions<contentservicev1_Category, Error>,
) {
  return useQuery({
    queryKey: ['getCategory', req],
    queryFn: () => apiClient.categoryService.Get(req),
    ...options,
  });
}

export function useCreateCategory(
  options?: UseMutationOptions<contentservicev1_Category, Error, { data: Record<string, any> }>,
) {
  return useMutation({
    mutationFn: (data) =>
      apiClient.categoryService.Create(data as { data: contentservicev1_Category }),
    ...options,
  });
}

export function useUpdateCategory(
  options?: UseMutationOptions<
    contentservicev1_Category,
    Error,
    { id: number; data: Record<string, any>; excludeMaskKeys?: string[] }
  >,
) {
  return useMutation({
    mutationFn: ({ id, data, excludeMaskKeys }) =>
      apiClient.categoryService.Update({
        id,
        data: data as unknown as contentservicev1_Category,
        updateMask: makeUpdateMask(
          Object.keys(data ?? {}).filter((k) => !excludeMaskKeys?.includes(k)),
        ),
      }),
    ...options,
  });
}

export function useDeleteCategory(
  options?: UseMutationOptions<{}, Error, contentservicev1_DeleteCategoryRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.categoryService.Delete(data),
    ...options,
  });
}

/**
 * 拉取拍平的分类下拉选项（含子分类，以「父级 / 」前缀区分层级）。
 * 供文章编辑页的分类选择与文章列表的分类筛选共用。
 */
export async function fetchFlattenedCategoryOptions(
  lang: string,
): Promise<{ label: string; value: number }[]> {
  const resp = await fetchListCategories(
    new PaginationQuery({
      paging: { page: 1, pageSize: 500 },
      fieldMask:
        'id,parent_id,children,translations.id,translations.language_code,translations.name',
    }),
  );

  const result: { label: string; value: number }[] = [];
  const walk = (items: any[] | undefined, prefix: string) => {
    items?.forEach((item) => {
      const translation =
        item.translations?.find((tr: any) => tr.languageCode === lang) ||
        item.translations?.[0];
      if (translation?.name) {
        result.push({ label: `${prefix}${translation.name}`, value: item.id });
        if (item.children?.length) {
          walk(item.children, `${prefix}${translation.name} / `);
        }
      }
    });
  };
  walk(resp.items, '');
  return result;
}
