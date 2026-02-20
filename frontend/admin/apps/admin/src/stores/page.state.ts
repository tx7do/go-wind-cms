import { computed } from 'vue';

import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { defineStore } from 'pinia';

import {
  createPageServiceClient,
  type contentservicev1_Page_PageStatus as Page_PageStatus,
  type contentservicev1_Page_PageType as Page_PageType,
} from '#/generated/api/admin/service/v1';
import { makeOrderBy, makeQueryString, makeUpdateMask } from '#/utils/query';
import { type Paging, requestClientRequestHandler } from '#/utils/request';

export const usePageStore = defineStore('page', () => {
  const service = createPageServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询页面列表
   */
  async function listPage(
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
   * 获取页面
   */
  async function getPage(id: number) {
    return await service.Get({ id });
  }

  /**
   * 创建页面
   */
  async function createPage(values: Record<string, any> = {}) {
    return await service.Create({
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新页面
   */
  async function updatePage(id: number, values: Record<string, any> = {}) {
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
   * 删除页面
   */
  async function deletePage(id: number) {
    return await service.Delete({ id });
  }

  function $reset() {}

  return {
    $reset,
    listPage,
    getPage,
    createPage,
    updatePage,
    deletePage,
  };
});

export const pageStatusList = computed(() => [
  {
    value: 'POST_STATUS_DRAFT',
    label: $t('enum.page.status.POST_STATUS_DRAFT'),
  },
  {
    value: 'POST_STATUS_PUBLISHED',
    label: $t('enum.page.status.POST_STATUS_PUBLISHED'),
  },
  {
    value: 'PAGE_STATUS_ARCHIVED',
    label: $t('enum.page.status.PAGE_STATUS_ARCHIVED'),
  },
]);

export function pageStatusToName(status: Page_PageStatus) {
  const values = pageStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const PAGE_STATUS_COLOR_MAP = {
  // 草稿 - 雅致紫（代表创作中、未完成，中性且有设计感）
  POST_STATUS_DRAFT: '#8b5cf6',
  // 已发布 - 自然青绿（代表成功、生效，积极正面且不刺眼）
  POST_STATUS_PUBLISHED: '#22c55e',
  // 已归档 - 中性暖棕（代表归档、封存，无负面情绪，仅体现状态变更）
  POST_STATUS_ARCHIVED: '#92400e',
  // 默认值 - 中性浅灰（无明确状态时的兜底，保持视觉中立）
  DEFAULT: '#94a3b8',
} as const;

export function pageStatusToColor(status: Page_PageStatus) {
  return (
    PAGE_STATUS_COLOR_MAP[status as keyof typeof PAGE_STATUS_COLOR_MAP] ||
    PAGE_STATUS_COLOR_MAP.DEFAULT
  );
}

export const pageTypeList = computed(() => [
  {
    value: 'PAGE_TYPE_DEFAULT',
    label: $t('enum.page.type.PAGE_TYPE_DEFAULT'),
  },
  {
    value: 'PAGE_TYPE_HOME',
    label: $t('enum.page.type.PAGE_TYPE_HOME'),
  },
  {
    value: 'PAGE_TYPE_ERROR_404',
    label: $t('enum.page.type.PAGE_TYPE_ERROR_404'),
  },
  {
    value: 'PAGE_TYPE_ERROR_500',
    label: $t('enum.page.type.PAGE_TYPE_ERROR_500'),
  },
  {
    value: 'PAGE_TYPE_CUSTOM',
    label: $t('enum.page.type.PAGE_TYPE_CUSTOM'),
  },
]);

export function pageTypeToName(type: Page_PageType) {
  const values = pageTypeList.value;
  const matchedItem = values.find((item) => item.value === type);
  return matchedItem ? matchedItem.label : '';
}

const PAGE_TYPE_COLOR_MAP = {
  // 默认页面 - 雅致紫（通用、基础、无特定属性的页面）
  PAGE_TYPE_DEFAULT: '#7c3aed',
  // 首页 - 品牌蓝（核心、门户、正向的主页面）
  PAGE_TYPE_HOME: '#3b82f6',
  // 404错误页 - 警示橙（资源未找到，轻度警示）
  PAGE_TYPE_ERROR_404: '#f97316',
  // 500错误页 - 警示红（服务器错误，重度警示）
  PAGE_TYPE_ERROR_500: '#ef4444',
  // 自定义页面 - 创意青蓝（个性化、定制化、灵活配置）
  PAGE_TYPE_CUSTOM: '#06b6d4',
  // 默认值 - 中性浅灰（无明确页面类型时的兜底）
  DEFAULT: '#94a3b8',
} as const;

export function pageTypeToColor(type: Page_PageType) {
  return (
    PAGE_TYPE_COLOR_MAP[type as keyof typeof PAGE_TYPE_COLOR_MAP] ||
    PAGE_TYPE_COLOR_MAP.DEFAULT
  );
}
