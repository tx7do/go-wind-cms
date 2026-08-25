/**
 * 通用枚举与工具函数
 * 从 stores/modules/api 迁移而来
 */
import { computed } from 'vue';

import { $t } from '@vben/locales';

import type {
  contentservicev1_Section,
  contentservicev1_SectionType,
} from '#/api/generated/admin/service/v1';

export const enableList = computed(() => [
  { value: 'true', label: $t('enum.enable.true') },
  { value: 'false', label: $t('enum.enable.false') },
]);

export const enableBoolList = computed(() => [
  { value: true, label: $t('enum.enable.true') },
  { value: false, label: $t('enum.enable.false') },
]);

export const successStatusList = computed(() => [
  { value: true, label: $t('enum.successStatus.success') },
  { value: false, label: $t('enum.successStatus.failed') },
]);

export function enableBoolToColor(
  enable: 'false' | 'FALSE' | 'False' | 'true' | 'TRUE' | 'True' | boolean,
) {
  switch (enable) {
    case false:
    case 'false':
    case 'FALSE':
    case 'False': {
      return '#8C8C8C';
    }
    case true:
    case 'true':
    case 'TRUE':
    case 'True': {
      return '#52C41A';
    }
    default: {
      return '#C9CDD4';
    }
  }
}

export function enableBoolToName(
  enable: 'false' | 'FALSE' | 'False' | 'true' | 'TRUE' | 'True' | boolean,
) {
  switch (enable) {
    case true:
    case 'true':
    case 'TRUE':
    case 'True': {
      return $t('enum.enable.true');
    }
    default: {
      return $t('enum.enable.false');
    }
  }
}

export const methodList = [
  { value: 'GET', label: 'GET' },
  { value: 'POST', label: 'POST' },
  { value: 'PUT', label: 'PUT' },
  { value: 'PATCH', label: 'PATCH' },
  { value: 'DELETE', label: 'DELETE' },
];

// HTTP 方法对应的 Tag 颜色类型
export const httpMethodTagTypeMap: Record<
  string,
  'danger' | 'info' | 'primary' | 'success' | 'warning'
> = {
  GET: 'success',
  POST: 'primary',
  PUT: 'warning',
  DELETE: 'danger',
  PATCH: 'info',
  HEAD: 'info',
  OPTIONS: 'info',
};

export const statusList = computed(() => [
  { value: 'ON', label: $t('enum.status.ON') },
  { value: 'OFF', label: $t('enum.status.OFF') },
]);

export function statusToName(status: 'OFF' | 'ON' | undefined) {
  const values = statusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

export function statusToColor(status: 'OFF' | 'ON' | undefined) {
  switch (status) {
    case 'OFF': {
      return '#8C8C8C';
    }
    case 'ON': {
      return '#52C41A';
    }
    default: {
      return '#C9CDD4';
    }
  }
}

// 成功/失败状态
// ==============================

export function successToColor(success: boolean) {
  return success ? 'limegreen' : 'crimson';
}

export function successToName(success: boolean) {
  return success
    ? $t('enum.successStatus.success')
    : $t('enum.successStatus.failed');
}

export function successToNameWithStatusCode(
  success: boolean,
  statusCode: number,
) {
  return success
    ? $t('enum.successStatus.success')
    : ` ${$t('enum.successStatus.failed')} (${statusCode})`;
}

// ==============================
// 页面区块类型枚举与工具函数
// 区块作为页面嵌套子部件管理，类型助手供页面编辑器内联使用。
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
