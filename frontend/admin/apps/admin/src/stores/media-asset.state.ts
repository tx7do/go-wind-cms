import { computed } from 'vue';

import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { defineStore } from 'pinia';

import {
  createMediaAssetServiceClient,
  type mediaservicev1_MediaAsset_AssetType as MediaAsset_AssetType,
  type mediaservicev1_MediaAsset_ProcessingStatus as MediaAsset_ProcessingStatus,
} from '#/generated/api/admin/service/v1';
import { makeOrderBy, makeQueryString, makeUpdateMask } from '#/utils/query';
import { type Paging, requestClientRequestHandler } from '#/utils/request';

export const useMediaAssetStore = defineStore('media-asset', () => {
  const service = createMediaAssetServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询媒体资源列表
   */
  async function listMediaAsset(
    paging?: Paging,
    formValues?: null | object,
    fieldMask?: null | string,
    orderBy?: null | string[],
  ) {
    const noPaging =
      paging?.page === undefined && paging?.pageSize === undefined;
    return await service.List({
      // @ts-ignore proto generated code is error.
      fieldMask,
      orderBy: makeOrderBy(orderBy),
      query: makeQueryString(formValues, userStore.isTenantUser()),
      page: paging?.page,
      pageSize: paging?.pageSize,
      noPaging,
    });
  }

  /**
   * 获取媒体资源
   */
  async function getMediaAsset(id: number) {
    return await service.Get({ id });
  }

  /**
   * 创建媒体资源
   */
  async function createMediaAsset(values: Record<string, any> = {}) {
    return await service.Create({
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新媒体资源
   */
  async function updateMediaAsset(
    id: number,
    values: Record<string, any> = {},
  ) {
    return await service.Update({
      id,
      data: {
        ...values,
      },
      // @ts-ignore proto generated code is error.
      updateMask: makeUpdateMask(Object.keys(values ?? [])),
    });
  }

  /**
   * 删除媒体资源
   */
  async function deleteMediaAsset(id: number) {
    return await service.Delete({ id });
  }

  function $reset() {}

  return {
    $reset,
    listMediaAsset,
    getMediaAsset,
    createMediaAsset,
    updateMediaAsset,
    deleteMediaAsset,
  };
});

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
  status: MediaAsset_ProcessingStatus,
) {
  const values = mediaAssetProcessingStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const MEDIA_ASSET_PROCESSING_STATUS_COLOR_MAP = {
  // 上传中 - 进度蓝（代表传输中、进行中，中性无情绪）
  PROCESSING_STATUS_UPLOADING: '#3b82f6',
  // 处理中 - 进度黄（代表后台处理、待完成，温和提醒）
  PROCESSING_STATUS_PROCESSING: '#f59e0b',
  // 已完成 - 成功绿（代表处理完成、成功，正向积极）
  PROCESSING_STATUS_COMPLETED: '#22c55e',
  // 处理失败 - 警示红（代表出错、失败，需关注）
  PROCESSING_STATUS_FAILED: '#ef4444',
  // 默认值 - 中性浅灰（无明确状态时的兜底）
  DEFAULT: '#94a3b8',
} as const;
export function mediaAssetProcessingStatusToColor(
  status: MediaAsset_ProcessingStatus,
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

export function mediaAssetAssetTypeToName(assetType: MediaAsset_AssetType) {
  const values = mediaAssetAssetTypeList.value;
  const matchedItem = values.find((item) => item.value === assetType);
  return matchedItem ? matchedItem.label : '';
}

const MEDIA_ASSET_ASSET_TYPE_COLOR_MAP = {
  // 图片 - 创意紫（视觉创意、图像类，贴合设计属性）
  ASSET_TYPE_IMAGE: '#8b5cf6',
  // 视频 - 活力蓝（动态影像、多媒体，贴合视频播放属性）
  ASSET_TYPE_VIDEO: '#3b82f6',
  // 文档 - 专业灰蓝（办公、文本、正式文档，中性专业）
  ASSET_TYPE_DOCUMENT: '#64748b',
  // 音频 - 活力青（声音、音频类，区别于视频蓝且有听觉联想）
  ASSET_TYPE_AUDIO: '#14b8a6',
  // 归档文件 - 沉稳棕（压缩、封存、打包，贴合归档属性）
  ASSET_TYPE_ARCHIVE: '#92400e',
  // 其他类型 - 中性浅紫（兜底分类，温和不抢镜）
  ASSET_TYPE_OTHER: '#a855f7',
  // 默认值 - 中性浅灰（无明确类型时的兜底）
  DEFAULT: '#94a3b8',
} as const;

export function mediaAssetAssetTypeToColor(assetType: MediaAsset_AssetType) {
  return (
    MEDIA_ASSET_ASSET_TYPE_COLOR_MAP[
      assetType as keyof typeof MEDIA_ASSET_ASSET_TYPE_COLOR_MAP
    ] || MEDIA_ASSET_ASSET_TYPE_COLOR_MAP.DEFAULT
  );
}
