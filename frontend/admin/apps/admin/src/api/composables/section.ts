import type {
  contentservicev1_DeleteSectionRequest,
  contentservicev1_GetSectionRequest,
  contentservicev1_ListSectionResponse,
  contentservicev1_Section,
  contentservicev1_SectionType,
} from '#/api/generated/admin/service/v1';

import { computed } from 'vue';

import { $t } from '@vben/locales';

import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/vue-query';

import { apiClient } from '#/api/client';
import { queryClient } from '#/plugins/vue-query';
import { makeUpdateMask, type PaginationQuery } from '#/transport/rest';

// ==============================
// 页面区块管理
// ==============================

export function useListSections(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListSectionResponse, Error>,
) {
  return useQuery({
    queryKey: ['listSections', query],
    queryFn: () => apiClient.sectionService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListSections(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listSections', params],
    queryFn: () => apiClient.sectionService.List(params.toRawParams()),
    staleTime: 0,
    retry: 0,
  });
}

export function useGetSection(
  req: contentservicev1_GetSectionRequest,
  options?: UseQueryOptions<contentservicev1_Section, Error>,
) {
  return useQuery({
    queryKey: ['getSection', req],
    queryFn: () => apiClient.sectionService.Get(req),
    ...options,
  });
}

export function useCreateSection(
  options?: UseMutationOptions<
    contentservicev1_Section,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.sectionService.Create({
        data: { ...values } as contentservicev1_Section,
      }),
    ...options,
  });
}

export function useUpdateSection(
  options?: UseMutationOptions<
    contentservicev1_Section,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.sectionService.Update({
        id,
        data: { ...values } as contentservicev1_Section,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteSection(
  options?: UseMutationOptions<
    object,
    Error,
    contentservicev1_DeleteSectionRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.sectionService.Delete(data),
    ...options,
  });
}

// ==============================
// 区块类型枚举与工具函数
// ==============================

export const sectionTypeList = computed(() => [
  {
    value: 'SECTION_TYPE_RICH_TEXT',
    label: $t('enum.section.type.SECTION_TYPE_RICH_TEXT'),
  },
  {
    value: 'SECTION_TYPE_MARKDOWN',
    label: $t('enum.section.type.SECTION_TYPE_MARKDOWN'),
  },
  {
    value: 'SECTION_TYPE_TITLE',
    label: $t('enum.section.type.SECTION_TYPE_TITLE'),
  },
  {
    value: 'SECTION_TYPE_IMAGE',
    label: $t('enum.section.type.SECTION_TYPE_IMAGE'),
  },
  {
    value: 'SECTION_TYPE_GALLERY',
    label: $t('enum.section.type.SECTION_TYPE_GALLERY'),
  },
  {
    value: 'SECTION_TYPE_VIDEO',
    label: $t('enum.section.type.SECTION_TYPE_VIDEO'),
  },
  {
    value: 'SECTION_TYPE_BUTTON',
    label: $t('enum.section.type.SECTION_TYPE_BUTTON'),
  },
  {
    value: 'SECTION_TYPE_DIVIDER',
    label: $t('enum.section.type.SECTION_TYPE_DIVIDER'),
  },
  {
    value: 'SECTION_TYPE_SPACER',
    label: $t('enum.section.type.SECTION_TYPE_SPACER'),
  },
  {
    value: 'SECTION_TYPE_CODE',
    label: $t('enum.section.type.SECTION_TYPE_CODE'),
  },
  {
    value: 'SECTION_TYPE_HTML',
    label: $t('enum.section.type.SECTION_TYPE_HTML'),
  },
  {
    value: 'SECTION_TYPE_FORM',
    label: $t('enum.section.type.SECTION_TYPE_FORM'),
  },
  {
    value: 'SECTION_TYPE_CAROUSEL',
    label: $t('enum.section.type.SECTION_TYPE_CAROUSEL'),
  },
  {
    value: 'SECTION_TYPE_CUSTOM',
    label: $t('enum.section.type.SECTION_TYPE_CUSTOM'),
  },
]);

export function sectionTypeToName(type: contentservicev1_Section['type']) {
  const values = sectionTypeList.value;
  const matchedItem = values.find((item) => item.value === type);
  return matchedItem ? matchedItem.label : '';
}

const SECTION_TYPE_COLOR_MAP = {
  SECTION_TYPE_RICH_TEXT: '#3b82f6',
  SECTION_TYPE_MARKDOWN: '#0ea5e9',
  SECTION_TYPE_TITLE: '#6366f1',
  SECTION_TYPE_IMAGE: '#8b5cf6',
  SECTION_TYPE_GALLERY: '#a855f7',
  SECTION_TYPE_VIDEO: '#ec4899',
  SECTION_TYPE_BUTTON: '#f43f5e',
  SECTION_TYPE_DIVIDER: '#64748b',
  SECTION_TYPE_SPACER: '#94a3b8',
  SECTION_TYPE_CODE: '#14b8a6',
  SECTION_TYPE_HTML: '#f97316',
  SECTION_TYPE_FORM: '#eab308',
  SECTION_TYPE_CAROUSEL: '#d946ef',
  SECTION_TYPE_CUSTOM: '#06b6d4',
  DEFAULT: '#94a3b8',
} as const;

export function sectionTypeToColor(type: contentservicev1_SectionType) {
  return (
    SECTION_TYPE_COLOR_MAP[type as keyof typeof SECTION_TYPE_COLOR_MAP] ||
    SECTION_TYPE_COLOR_MAP.DEFAULT
  );
}
