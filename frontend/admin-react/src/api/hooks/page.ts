import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_DeletePageRequest,
  type contentservicev1_GetPageRequest,
  type contentservicev1_ListPageResponse,
  type contentservicev1_Page,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 页面管理
// ==============================

export function useListPages(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListPageResponse, Error>,
) {
  return useQuery({
    queryKey: ['listPages', query],
    queryFn: () => apiClient.pageService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListPages(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listPages', params],
    queryFn: () => apiClient.pageService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetPage(
  req: contentservicev1_GetPageRequest,
  options?: UseQueryOptions<contentservicev1_Page, Error>,
) {
  return useQuery({
    queryKey: ['getPage', req],
    queryFn: () => apiClient.pageService.Get(req),
    ...options,
  });
}

export function useCreatePage(
  options?: UseMutationOptions<contentservicev1_Page, Error, Record<string, any>>,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.pageService.Create({ data: { ...values } as contentservicev1_Page }),
    ...options,
  });
}

export function useUpdatePage(
  options?: UseMutationOptions<
    contentservicev1_Page,
    Error,
    { id: number; data: Record<string, any>; excludeMaskKeys?: string[] }
  >,
) {
  return useMutation({
    mutationFn: ({ id, data, excludeMaskKeys }) =>
      apiClient.pageService.Update({
        id,
        data: { ...data } as contentservicev1_Page,
        updateMask: makeUpdateMask(
          Object.keys(data ?? {}).filter((k) => !excludeMaskKeys?.includes(k)),
        ),
      }),
    ...options,
  });
}

export function useDeletePage(
  options?: UseMutationOptions<{}, Error, contentservicev1_DeletePageRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.pageService.Delete(data),
    ...options,
  });
}
