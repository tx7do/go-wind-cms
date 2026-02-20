import { computed } from 'vue';

import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { defineStore } from 'pinia';

import {
  createTagServiceClient,
  type contentservicev1_Tag_TagStatus as Tag_TagStatus,
} from '#/generated/api/admin/service/v1';
import { makeOrderBy, makeQueryString, makeUpdateMask } from '#/utils/query';
import { type Paging, requestClientRequestHandler } from '#/utils/request';

export const useTagStore = defineStore('tag', () => {
  const service = createTagServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询标签列表
   */
  async function listTag(
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
   * 获取标签
   */
  async function getTag(id: number) {
    return await service.Get({ id });
  }

  /**
   * 创建标签
   */
  async function createTag(values: Record<string, any> = {}) {
    return await service.Create({
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新标签
   */
  async function updateTag(id: number, values: Record<string, any> = {}) {
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
   * 删除标签
   */
  async function deleteTag(id: number) {
    return await service.Delete({ id });
  }

  function $reset() {}

  return {
    $reset,
    listTag,
    getTag,
    createTag,
    updateTag,
    deleteTag,
  };
});

export const tagStatusList = computed(() => [
  {
    value: 'TAG_STATUS_ACTIVE',
    label: $t('enum.tag.status.TAG_STATUS_ACTIVE'),
  },
  {
    value: 'TAG_STATUS_HIDDEN',
    label: $t('enum.tag.status.TAG_STATUS_HIDDEN'),
  },
  {
    value: 'TAG_STATUS_ARCHIVED',
    label: $t('enum.tag.status.TAG_STATUS_ARCHIVED'),
  },
]);

export function tagStatusToName(status: Tag_TagStatus) {
  const values = tagStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const TAG_STATUS_COLOR_MAP = {
  // 活跃 - 自然青绿（代表启用、正常、正向可用，符合核心状态定位）
  TAG_STATUS_ACTIVE: '#22c55e',
  // 隐藏 - 警示暖橙（代表暂不可见、非删除，温和提醒属性）
  TAG_STATUS_HIDDEN: '#f97316',
  // 已归档 - 中性暖棕（代表封存、历史状态，无负面情绪）
  TAG_STATUS_ARCHIVED: '#92400e',
  // 默认值 - 中性浅灰（无明确状态时的兜底，保持视觉中立）
  DEFAULT: '#94a3b8',
} as const;

export function tagStatusToColor(status: Tag_TagStatus) {
  return (
    TAG_STATUS_COLOR_MAP[status as keyof typeof TAG_STATUS_COLOR_MAP] ||
    TAG_STATUS_COLOR_MAP.DEFAULT
  );
}
