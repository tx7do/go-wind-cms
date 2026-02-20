import { computed } from 'vue';

import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { defineStore } from 'pinia';

import {
  type contentservicev1_Category_CategoryStatus as Category_CategoryStatus,
  createCategoryServiceClient,
} from '#/generated/api/admin/service/v1';
import { makeOrderBy, makeQueryString, makeUpdateMask } from '#/utils/query';
import { type Paging, requestClientRequestHandler } from '#/utils/request';

export const useCategoryStore = defineStore('category', () => {
  const service = createCategoryServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询分类列表
   */
  async function listCategory(
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
   * 获取分类
   */
  async function getCategory(id: number) {
    return await service.Get({ id });
  }

  /**
   * 创建分类
   */
  async function createCategory(values: Record<string, any> = {}) {
    return await service.Create({
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新分类
   */
  async function updateCategory(id: number, values: Record<string, any> = {}) {
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
   * 删除分类
   */
  async function deleteCategory(id: number) {
    return await service.Delete({ id });
  }

  function $reset() {}

  return {
    $reset,
    listCategory,
    getCategory,
    createCategory,
    updateCategory,
    deleteCategory,
  };
});

export const categoryStatusList = computed(() => [
  {
    value: 'CATEGORY_STATUS_ACTIVE',
    label: $t('enum.category.status.CATEGORY_STATUS_ACTIVE'),
  },
  {
    value: 'CATEGORY_STATUS_HIDDEN',
    label: $t('enum.category.status.CATEGORY_STATUS_HIDDEN'),
  },
  {
    value: 'CATEGORY_STATUS_ARCHIVED',
    label: $t('enum.category.status.CATEGORY_STATUS_ARCHIVED'),
  },
]);

export function categoryStatusToName(status: Category_CategoryStatus) {
  const values = categoryStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const CATEGORY_STATUS_COLOR_MAP = {
  // 活跃 - 自然青绿（代表启用、正常、正向可用，符合核心状态定位）
  CATEGORY_STATUS_ACTIVE: '#22c55e',
  // 隐藏 - 警示暖橙（代表暂不可见、非删除，温和提醒属性）
  CATEGORY_STATUS_HIDDEN: '#f97316',
  // 已归档 - 中性暖棕（代表封存、历史状态，无负面情绪）
  CATEGORY_STATUS_ARCHIVED: '#92400e',
  // 默认值 - 中性浅灰（无明确状态时的兜底，保持视觉中立）
  DEFAULT: '#94a3b8',
} as const;

export function categoryStatusToColor(status: Category_CategoryStatus) {
  return (
    CATEGORY_STATUS_COLOR_MAP[
      status as keyof typeof CATEGORY_STATUS_COLOR_MAP
    ] || CATEGORY_STATUS_COLOR_MAP.DEFAULT
  );
}
