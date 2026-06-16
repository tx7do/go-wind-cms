import type {
  siteservicev1_DeleteSiteSettingRequest,
  siteservicev1_GetSiteSettingRequest,
  siteservicev1_ListSiteSettingResponse,
  siteservicev1_SiteSetting,
  siteservicev1_SiteSetting_SettingType,
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
// 站点配置管理
// ==============================

export function useListSiteSettings(
  query: PaginationQuery,
  options?: UseQueryOptions<siteservicev1_ListSiteSettingResponse, Error>,
) {
  return useQuery({
    queryKey: ['listSiteSettings', query],
    queryFn: () => apiClient.siteSettingService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListSiteSettings(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listSiteSettings', params],
    queryFn: () => apiClient.siteSettingService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetSiteSetting(
  req: siteservicev1_GetSiteSettingRequest,
  options?: UseQueryOptions<siteservicev1_SiteSetting, Error>,
) {
  return useQuery({
    queryKey: ['getSiteSetting', req],
    queryFn: () => apiClient.siteSettingService.Get(req),
    ...options,
  });
}

export function useCreateSiteSetting(
  options?: UseMutationOptions<
    siteservicev1_SiteSetting,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.siteSettingService.Create({
        data: { ...values } as siteservicev1_SiteSetting,
      }),
    ...options,
  });
}

export function useUpdateSiteSetting(
  options?: UseMutationOptions<
    siteservicev1_SiteSetting,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.siteSettingService.Update({
        id,
        data: { ...values } as siteservicev1_SiteSetting,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteSiteSetting(
  options?: UseMutationOptions<
    object,
    Error,
    siteservicev1_DeleteSiteSettingRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.siteSettingService.Delete(data),
    ...options,
  });
}

// ==============================
// 站点配置枚举与工具函数
// ==============================

export const siteSettingTypeList = computed(() => [
  {
    value: 'SETTING_TYPE_TEXT',
    label: $t('enum.siteSetting.type.SETTING_TYPE_TEXT'),
  },
  {
    value: 'SETTING_TYPE_TEXTAREA',
    label: $t('enum.siteSetting.type.SETTING_TYPE_TEXTAREA'),
  },
  {
    value: 'SETTING_TYPE_NUMBER',
    label: $t('enum.siteSetting.type.SETTING_TYPE_NUMBER'),
  },
  {
    value: 'SETTING_TYPE_BOOLEAN',
    label: $t('enum.siteSetting.type.SETTING_TYPE_BOOLEAN'),
  },
  {
    value: 'SETTING_TYPE_URL',
    label: $t('enum.siteSetting.type.SETTING_TYPE_URL'),
  },
  {
    value: 'SETTING_TYPE_EMAIL',
    label: $t('enum.siteSetting.type.SETTING_TYPE_EMAIL'),
  },
  {
    value: 'SETTING_TYPE_IMAGE',
    label: $t('enum.siteSetting.type.SETTING_TYPE_IMAGE'),
  },
  {
    value: 'SETTING_TYPE_SELECT',
    label: $t('enum.siteSetting.type.SETTING_TYPE_SELECT'),
  },
  {
    value: 'SETTING_TYPE_JSON',
    label: $t('enum.siteSetting.type.SETTING_TYPE_JSON'),
  },
]);

export function siteSettingTypeToName(
  settingType: siteservicev1_SiteSetting_SettingType,
) {
  const values = siteSettingTypeList.value;
  const matchedItem = values.find((item) => item.value === settingType);
  return matchedItem ? matchedItem.label : '';
}

const SITE_SETTING_TYPE_COLOR_MAP = {
  SETTING_TYPE_TEXT: '#3b82f6',
  SETTING_TYPE_TEXTAREA: '#60a5fa',
  SETTING_TYPE_NUMBER: '#22c55e',
  SETTING_TYPE_BOOLEAN: '#f59e0b',
  SETTING_TYPE_URL: '#14b8a6',
  SETTING_TYPE_EMAIL: '#8b5cf6',
  SETTING_TYPE_IMAGE: '#f97316',
  SETTING_TYPE_SELECT: '#a855f7',
  SETTING_TYPE_JSON: '#0ea5e9',
  DEFAULT: '#94a3b8',
} as const;

export function siteSettingTypeToColor(
  settingType: siteservicev1_SiteSetting_SettingType,
) {
  return (
    SITE_SETTING_TYPE_COLOR_MAP[
      settingType as keyof typeof SITE_SETTING_TYPE_COLOR_MAP
    ] || SITE_SETTING_TYPE_COLOR_MAP.DEFAULT
  );
}
