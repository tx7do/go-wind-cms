/**
 * 媒体资源模块常量（类型/处理状态枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const ASSET_TYPE_COLOR_MAP: Record<string, string> = {
  ASSET_TYPE_IMAGE: 'purple',
  ASSET_TYPE_VIDEO: 'blue',
  ASSET_TYPE_DOCUMENT: 'default',
  ASSET_TYPE_AUDIO: 'cyan',
  ASSET_TYPE_ARCHIVE: 'gold',
  ASSET_TYPE_OTHER: 'magenta',
  DEFAULT: 'default',
};

const ASSET_TYPE_VALUES = [
  'ASSET_TYPE_IMAGE',
  'ASSET_TYPE_VIDEO',
  'ASSET_TYPE_DOCUMENT',
  'ASSET_TYPE_AUDIO',
  'ASSET_TYPE_ARCHIVE',
  'ASSET_TYPE_OTHER',
] as const;

const PROCESSING_STATUS_COLOR_MAP: Record<string, string> = {
  PROCESSING_STATUS_UPLOADING: 'blue',
  PROCESSING_STATUS_PROCESSING: 'gold',
  PROCESSING_STATUS_COMPLETED: 'green',
  PROCESSING_STATUS_FAILED: 'red',
  DEFAULT: 'default',
};

const PROCESSING_STATUS_VALUES = [
  'PROCESSING_STATUS_UPLOADING',
  'PROCESSING_STATUS_PROCESSING',
  'PROCESSING_STATUS_COMPLETED',
  'PROCESSING_STATUS_FAILED',
] as const;

export function getAssetTypeColor(type: string): string {
  return ASSET_TYPE_COLOR_MAP[type] || ASSET_TYPE_COLOR_MAP.DEFAULT;
}

export function getAssetTypeLabel(t: TFn, type: string): string {
  // i18next 未开 returnObjects，取对象需走点路径（对齐 dict 模块约定）
  return t(`typeMap.${type}`, { defaultValue: type });
}

export function assetTypeOptions(t: TFn) {
  return ASSET_TYPE_VALUES.map((value) => ({
    label: t(`typeMap.${value}`),
    value,
  }));
}

export function getProcessingStatusColor(status: string): string {
  return PROCESSING_STATUS_COLOR_MAP[status] || PROCESSING_STATUS_COLOR_MAP.DEFAULT;
}

export function getProcessingStatusLabel(t: TFn, status: string): string {
  return t(`processingStatusMap.${status}`, { defaultValue: status });
}

export function processingStatusOptions(t: TFn) {
  return PROCESSING_STATUS_VALUES.map((value) => ({
    label: t(`processingStatusMap.${value}`),
    value,
  }));
}
