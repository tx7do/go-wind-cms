import { computed } from 'vue';

import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { defineStore } from 'pinia';

import {
  createSiteServiceClient,
  type siteservicev1_Site_Status as Site_Status,
} from '#/generated/api/admin/service/v1';
import { makeOrderBy, makeQueryString, makeUpdateMask } from '#/utils/query';
import { type Paging, requestClientRequestHandler } from '#/utils/request';

export const useSiteStore = defineStore('site', () => {
  const service = createSiteServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询站点列表
   */
  async function listSite(
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
   * 获取站点
   */
  async function getSite(id: number) {
    return await service.Get({ id });
  }

  /**
   * 创建站点
   */
  async function createSite(values: Record<string, any> = {}) {
    return await service.Create({
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新站点
   */
  async function updateSite(id: number, values: Record<string, any> = {}) {
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
   * 删除站点
   */
  async function deleteSite(id: number) {
    return await service.Delete({ id });
  }

  function $reset() {}

  return {
    $reset,
    listSite,
    getSite,
    createSite,
    updateSite,
    deleteSite,
  };
});

export const siteStatusList = computed(() => [
  {
    value: 'SITE_STATUS_ACTIVE',
    label: $t('enum.site.status.SITE_STATUS_ACTIVE'),
  },
  {
    value: 'SITE_STATUS_INACTIVE',
    label: $t('enum.site.status.SITE_STATUS_INACTIVE'),
  },
  {
    value: 'SITE_STATUS_MAINTENANCE',
    label: $t('enum.site.status.SITE_STATUS_MAINTENANCE'),
  },
]);

export function siteStatusToName(assetType: Site_Status) {
  const values = siteStatusList.value;
  const matchedItem = values.find((item) => item.value === assetType);
  return matchedItem ? matchedItem.label : '';
}

const SITE_STATUS_COLOR_MAP = {
  SITE_STATUS_ACTIVE: '#3b82f6',
  SITE_STATUS_INACTIVE: '#60a5fa',
  SITE_STATUS_MAINTENANCE: '#22c55e',
  // 默认值 - 中性浅灰（无明确类型时的兜底）
  DEFAULT: '#94a3b8',
} as const;

export function siteStatusToColor(assetType: Site_Status) {
  return (
    SITE_STATUS_COLOR_MAP[assetType as keyof typeof SITE_STATUS_COLOR_MAP] ||
    SITE_STATUS_COLOR_MAP.DEFAULT
  );
}
