/**
 * 媒体资源模块常量（类型/处理状态枚举展示）
 */

type TFn = (key: string, options?: Record<string, any>) => string;

const ASSET_TYPE_COLOR_MAP: Record<string, string> = {
  ASSET_TYPE_IMAGE: '#8b5cf6',
  ASSET_TYPE_VIDEO: '#3b82f6',
  ASSET_TYPE_DOCUMENT: '#64748b',
  ASSET_TYPE_AUDIO: '#14b8a6',
  ASSET_TYPE_ARCHIVE: '#92400e',
  ASSET_TYPE_OTHER: '#a855f7',
  DEFAULT: '#94a3b8',
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
  PROCESSING_STATUS_UPLOADING: '#3b82f6',
  PROCESSING_STATUS_PROCESSING: '#f59e0b',
  PROCESSING_STATUS_COMPLETED: '#22c55e',
  PROCESSING_STATUS_FAILED: '#ef4444',
  DEFAULT: '#94a3b8',
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
