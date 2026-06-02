import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_Page,
  type contentservicev1_ListPageResponse,
} from '@/api/generated/app/service/v1';
import {
  listPagesRaw,
  getPage,
  createPage,
  updatePage,
  deletePage,
} from '@/api/service/page';
import { queryClient } from '@/core';

// ==============================
// 页面列表 Hook
// ==============================
export function useListPages(
  options?: UseMutationOptions<
    contentservicev1_ListPageResponse,
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
    mutationFn: (params) => listPagesRaw(params),
    ...options,
  });
}

// ==============================================
// 获取页面列表 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchListPages(params: {
  paging?: { page?: number; pageSize?: number };
  formValues?: object | undefined;
  fieldMask?: string | undefined;
  orderBy?: string[] | undefined;
}) {
  return queryClient.fetchQuery({
    queryKey: ['listPages', params],
    queryFn: () => listPagesRaw(params),
    retry: 0,
  });
}

// ==============================
// 获取单个页面 Hook
// ==============================
export function useGetPage(
  options?: UseMutationOptions<contentservicev1_Page, Error, number>,
) {
  return useMutation({
    mutationFn: (id) => getPage(id),
    ...options,
  });
}

// ==============================================
// 获取单个页面 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchPage(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getPage', id],
    queryFn: () => getPage(id),
    retry: 0,
  });
}

// ==============================
// 创建页面 Hook
// ==============================
export function useCreatePage(
  options?: UseMutationOptions<contentservicev1_Page, Error, Partial<contentservicev1_Page>>,
) {
  return useMutation({
    mutationFn: (data) => createPage(data),
    ...options,
  });
}

// ==============================
// 更新页面 Hook
// ==============================
export function useUpdatePage(
  options?: UseMutationOptions<
    contentservicev1_Page,
    Error,
    { id: number; values: Partial<contentservicev1_Page> }
  >,
) {
  return useMutation({
    mutationFn: (params) => updatePage(params),
    ...options,
  });
}

// ==============================
// 删除页面 Hook
// ==============================
export function useDeletePage(options?: UseMutationOptions<Record<string, never>, Error, number>) {
  return useMutation({
    mutationFn: (id) => deletePage(id),
    ...options,
  });
}
