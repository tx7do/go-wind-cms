import type {
  mediaservicev1_DeleteMediaAssetRequest,
  mediaservicev1_GetMediaAssetRequest,
  mediaservicev1_ListMediaAssetResponse,
  mediaservicev1_MediaAsset,
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
// 媒体资源管理
// ==============================

export function useListMediaAssets(
  query: PaginationQuery,
  options?: UseQueryOptions<mediaservicev1_ListMediaAssetResponse, Error>,
) {
  return useQuery({
    queryKey: ['listMediaAssets', query],
    queryFn: () => apiClient.mediaAssetService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListMediaAssets(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listMediaAssets', params],
    queryFn: () => apiClient.mediaAssetService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetMediaAsset(
  req: mediaservicev1_GetMediaAssetRequest,
  options?: UseQueryOptions<mediaservicev1_MediaAsset, Error>,
) {
  return useQuery({
    queryKey: ['getMediaAsset', req],
    queryFn: () => apiClient.mediaAssetService.Get(req),
    ...options,
  });
}

export function useCreateMediaAsset(
  options?: UseMutationOptions<
    mediaservicev1_MediaAsset,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.mediaAssetService.Create({
        data: { ...values } as mediaservicev1_MediaAsset,
      }),
    ...options,
  });
}

export function useUpdateMediaAsset(
  options?: UseMutationOptions<
    mediaservicev1_MediaAsset,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.mediaAssetService.Update({
        id,
        data: { ...values } as mediaservicev1_MediaAsset,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteMediaAsset(
  options?: UseMutationOptions<
    object,
    Error,
    mediaservicev1_DeleteMediaAssetRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.mediaAssetService.Delete(data),
    ...options,
  });
}

// ==============================
// 媒体资源枚举与工具函数
// ==============================

export const mediaAssetProcessingStatusList = computed(() => [
  {
    value: 'PROCESSING_STATUS_UPLOADING',
    label: $t('enum.mediaAsset.processingStatus.PROCESSING_STATUS_UPLOADING'),
  },
  {
    value: 'PROCESSING_STATUS_PROCESSING',
    label: $t('enum.mediaAsset.processingStatus.PROCESSING_STATUS_PROCESSING'),
  },
  {
    value: 'PROCESSING_STATUS_COMPLETED',
    label: $t('enum.mediaAsset.processingStatus.PROCESSING_STATUS_COMPLETED'),
  },
  {
    value: 'PROCESSING_STATUS_FAILED',
    label: $t('enum.mediaAsset.processingStatus.PROCESSING_STATUS_FAILED'),
  },
]);

export function mediaAssetProcessingStatusToName(
  status: mediaservicev1_MediaAsset['processingStatus'],
) {
  const values = mediaAssetProcessingStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const MEDIA_ASSET_PROCESSING_STATUS_COLOR_MAP = {
  PROCESSING_STATUS_UPLOADING: '#3b82f6',
  PROCESSING_STATUS_PROCESSING: '#f59e0b',
  PROCESSING_STATUS_COMPLETED: '#22c55e',
  PROCESSING_STATUS_FAILED: '#ef4444',
  DEFAULT: '#94a3b8',
} as const;

export function mediaAssetProcessingStatusToColor(
  status: mediaservicev1_MediaAsset['processingStatus'],
) {
  return (
    MEDIA_ASSET_PROCESSING_STATUS_COLOR_MAP[
      status as keyof typeof MEDIA_ASSET_PROCESSING_STATUS_COLOR_MAP
    ] || MEDIA_ASSET_PROCESSING_STATUS_COLOR_MAP.DEFAULT
  );
}

export const mediaAssetAssetTypeList = computed(() => [
  {
    value: 'ASSET_TYPE_IMAGE',
    label: $t('enum.mediaAsset.type.ASSET_TYPE_IMAGE'),
  },
  {
    value: 'ASSET_TYPE_VIDEO',
    label: $t('enum.mediaAsset.type.ASSET_TYPE_VIDEO'),
  },
  {
    value: 'ASSET_TYPE_DOCUMENT',
    label: $t('enum.mediaAsset.type.ASSET_TYPE_DOCUMENT'),
  },
  {
    value: 'ASSET_TYPE_AUDIO',
    label: $t('enum.mediaAsset.type.ASSET_TYPE_AUDIO'),
  },
  {
    value: 'ASSET_TYPE_ARCHIVE',
    label: $t('enum.mediaAsset.type.ASSET_TYPE_ARCHIVE'),
  },
  {
    value: 'ASSET_TYPE_OTHER',
    label: $t('enum.mediaAsset.type.ASSET_TYPE_OTHER'),
  },
]);

export function mediaAssetAssetTypeToName(
  assetType: mediaservicev1_MediaAsset['type'],
) {
  const values = mediaAssetAssetTypeList.value;
  const matchedItem = values.find((item) => item.value === assetType);
  return matchedItem ? matchedItem.label : '';
}

const MEDIA_ASSET_ASSET_TYPE_COLOR_MAP = {
  ASSET_TYPE_IMAGE: '#8b5cf6',
  ASSET_TYPE_VIDEO: '#3b82f6',
  ASSET_TYPE_DOCUMENT: '#64748b',
  ASSET_TYPE_AUDIO: '#14b8a6',
  ASSET_TYPE_ARCHIVE: '#92400e',
  ASSET_TYPE_OTHER: '#a855f7',
  DEFAULT: '#94a3b8',
} as const;

export function mediaAssetAssetTypeToColor(
  assetType: mediaservicev1_MediaAsset['type'],
) {
  return (
    MEDIA_ASSET_ASSET_TYPE_COLOR_MAP[
      assetType as keyof typeof MEDIA_ASSET_ASSET_TYPE_COLOR_MAP
    ] || MEDIA_ASSET_ASSET_TYPE_COLOR_MAP.DEFAULT
  );
}
