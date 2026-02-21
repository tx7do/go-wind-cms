import { computed } from 'vue';

import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { defineStore } from 'pinia';

import {
  createSiteSettingServiceClient,
  type mediaservicev1_MediaAsset_AssetType as MediaAsset_AssetType,
} from '#/generated/api/admin/service/v1';
import { makeOrderBy, makeQueryString, makeUpdateMask } from '#/utils/query';
import { type Paging, requestClientRequestHandler } from '#/utils/request';

export const useSiteSettingStore = defineStore('site-setting', () => {
  const service = createSiteSettingServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询站点配置列表
   */
  async function listSiteSetting(
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
   * 获取站点配置
   */
  async function getSiteSetting(id: number) {
    return await service.Get({ id });
  }

  /**
   * 创建站点配置
   */
  async function createSiteSetting(values: Record<string, any> = {}) {
    return await service.Create({
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新站点配置
   */
  async function updateSiteSetting(
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
   * 删除站点配置
   */
  async function deleteSiteSetting(id: number) {
    return await service.Delete({ id });
  }

  function $reset() {}

  return {
    $reset,
    listSiteSetting,
    getSiteSetting,
    createSiteSetting,
    updateSiteSetting,
    deleteSiteSetting,
  };
});

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

export function siteSettingTypeToName(assetType: MediaAsset_AssetType) {
  const values = siteSettingTypeList.value;
  const matchedItem = values.find((item) => item.value === assetType);
  return matchedItem ? matchedItem.label : '';
}

const SITE_SETTING_TYPE_COLOR_MAP = {
  // 单行文本 - 基础蓝（通用、基础文本输入，中性属性）
  SETTING_TYPE_TEXT: '#3b82f6',
  // 多行文本 - 浅蓝（与单行文本同源，区分多行特性）
  SETTING_TYPE_TEXTAREA: '#60a5fa',
  // 数字 - 自然绿（数值、计数、量化属性）
  SETTING_TYPE_NUMBER: '#22c55e',
  // 布尔值 - 暖黄（开关、二选一，轻量提醒属性）
  SETTING_TYPE_BOOLEAN: '#f59e0b',
  // URL - 科技青（链接、网络、地址属性）
  SETTING_TYPE_URL: '#14b8a6',
  // 邮箱 - 雅致紫（通讯、账号类属性）
  SETTING_TYPE_EMAIL: '#8b5cf6',
  // 图片 - 暖橙（视觉、媒体、上传类属性）
  SETTING_TYPE_IMAGE: '#f97316',
  // 下拉选择 - 梅子紫（选项、筛选、选择类属性）
  SETTING_TYPE_SELECT: '#a855f7',
  // JSON - 冷蓝（代码、结构化数据、配置类属性）
  SETTING_TYPE_JSON: '#0ea5e9',
  // 默认值 - 中性浅灰（无明确类型时的兜底）
  DEFAULT: '#94a3b8',
} as const;

export function siteSettingTypeToColor(assetType: MediaAsset_AssetType) {
  return (
    SITE_SETTING_TYPE_COLOR_MAP[
      assetType as keyof typeof SITE_SETTING_TYPE_COLOR_MAP
    ] || SITE_SETTING_TYPE_COLOR_MAP.DEFAULT
  );
}
