import { computed } from 'vue';

import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { defineStore } from 'pinia';

import {
  createPostServiceClient,
  type contentservicev1_EditorType as EditorType,
  type contentservicev1_Post_PostStatus as Post_PostStatus,
} from '#/generated/api/admin/service/v1';
import { makeOrderBy, makeQueryString, makeUpdateMask } from '#/utils/query';
import { type Paging, requestClientRequestHandler } from '#/utils/request';

export const usePostStore = defineStore('post', () => {
  const service = createPostServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询帖子列表
   */
  async function listPost(
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
   * 获取帖子
   */
  async function getPost(id: number) {
    return await service.Get({ id });
  }

  /**
   * 创建帖子
   */
  async function createPost(values: Record<string, any> = {}) {
    return await service.Create({
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新帖子
   */
  async function updatePost(id: number, values: Record<string, any> = {}) {
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
   * 删除帖子
   */
  async function deletePost(id: number) {
    return await service.Delete({ id });
  }

  function $reset() {}

  return {
    $reset,
    listPost,
    getPost,
    createPost,
    updatePost,
    deletePost,
  };
});

export const postStatusList = computed(() => [
  {
    value: 'POST_STATUS_DRAFT',
    label: $t('enum.post.status.POST_STATUS_DRAFT'),
  },
  {
    value: 'POST_STATUS_PUBLISHED',
    label: $t('enum.post.status.POST_STATUS_PUBLISHED'),
  },
  {
    value: 'POST_STATUS_SCHEDULED',
    label: $t('enum.post.status.POST_STATUS_SCHEDULED'),
  },
  {
    value: 'POST_STATUS_TRASHED',
    label: $t('enum.post.status.POST_STATUS_TRASHED'),
  },
]);

export function postStatusToName(status: Post_PostStatus) {
  const values = postStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const POST_STATUS_COLOR_MAP = {
  // 草稿 - 柔和蓝灰色（代表未完成、待编辑，中性无情绪）
  POST_STATUS_DRAFT: '#7c3aed',
  // 已发布 - 清新绿（代表成功、生效，积极正面）
  POST_STATUS_PUBLISHED: '#10b981',
  // 定时发布 - 暖黄色（代表等待、计划，温和的提醒属性）
  POST_STATUS_SCHEDULED: '#f59e0b',
  // 回收站 - 低饱和红（代表删除、归档，警示但不刺眼）
  POST_STATUS_TRASHED: '#ef4444',
  // 默认值 - 中性浅灰（无明确状态时的兜底，保持视觉中立）
  DEFAULT: '#94a3b8',
} as const;

export function postStatusToColor(status: Post_PostStatus) {
  return (
    POST_STATUS_COLOR_MAP[status as keyof typeof POST_STATUS_COLOR_MAP] ||
    POST_STATUS_COLOR_MAP.DEFAULT
  );
}

export const editorTypeList = computed(() => [
  {
    value: 'EDITOR_TYPE_MARKDOWN',
    label: $t('enum.editorType.EDITOR_TYPE_MARKDOWN'),
  },
  {
    value: 'EDITOR_TYPE_RICH_TEXT',
    label: $t('enum.editorType.EDITOR_TYPE_RICH_TEXT'),
  },
  {
    value: 'EDITOR_TYPE_HTML',
    label: $t('enum.editorType.EDITOR_TYPE_HTML'),
  },
  {
    value: 'EDITOR_TYPE_JSON_BLOCK',
    label: $t('enum.editorType.EDITOR_TYPE_JSON_BLOCK'),
  },
  {
    value: 'EDITOR_TYPE_PLAIN_TEXT',
    label: $t('enum.editorType.EDITOR_TYPE_PLAIN_TEXT'),
  },
  {
    value: 'EDITOR_TYPE_CODE',
    label: $t('enum.editorType.EDITOR_TYPE_CODE'),
  },
  {
    value: 'EDITOR_TYPE_WYSIWYG',
    label: $t('enum.editorType.EDITOR_TYPE_WYSIWYG'),
  },
  {
    value: 'EDITOR_TYPE_VISUAL_BUILDER',
    label: $t('enum.editorType.EDITOR_TYPE_VISUAL_BUILDER'),
  },
  {
    value: 'EDITOR_TYPE_SLATE',
    label: $t('enum.editorType.EDITOR_TYPE_SLATE'),
  },
  {
    value: 'EDITOR_TYPE_PROSEMIRROR',
    label: $t('enum.editorType.EDITOR_TYPE_PROSEMIRROR'),
  },
]);

export function editorTypeToName(editorType: EditorType) {
  const values = editorTypeList.value;
  const matchedItem = values.find((item) => item.value === editorType);
  return matchedItem ? matchedItem.label : '';
}

const EDITOR_TYPE_COLOR_MAP = {
  // Markdown编辑器 - 优雅靛蓝（轻量、简洁、文本化）
  EDITOR_TYPE_MARKDOWN: '#4f46e5',
  // 富文本编辑器 - 清新青柠绿（可视化、易用、富内容）
  EDITOR_TYPE_RICH_TEXT: '#14b8a6',
  // HTML编辑器 - 专业暖橙（代码、定制化、网页原生）
  EDITOR_TYPE_HTML: '#f97316',
  // JSON块编辑器 - 沉稳紫玫红（结构化、数据、严谨）
  EDITOR_TYPE_JSON_BLOCK: '#d946ef',
  // 纯文本编辑器 - 中性石板灰（极简、无格式、基础文本）
  EDITOR_TYPE_PLAIN_TEXT: '#64748b',
  // 代码编辑器 - 科技冷蓝（编程、语法高亮、开发者工具）
  EDITOR_TYPE_CODE: '#0ea5e9',
  // WYSIWYG编辑器 - 柔和薄荷绿（所见即所得，与富文本同源但更浅）
  EDITOR_TYPE_WYSIWYG: '#34d399',
  // 可视化构建器 - 活力青蓝（拖拽、可视化、低代码）
  EDITOR_TYPE_VISUAL_BUILDER: '#06b6d4',
  // Slate编辑器 - 雅致紫罗兰（现代富文本框架、轻量定制）
  EDITOR_TYPE_SLATE: '#8b5cf6',
  // ProseMirror编辑器 - 柔和梅子紫（专业编辑器框架、结构化编辑）
  EDITOR_TYPE_PROSEMIRROR: '#a855f7',
  // 默认值 - 中性浅灰（无明确类型时的兜底）
  DEFAULT: '#94a3b8',
} as const;

export function editorTypeToColor(editorType: EditorType) {
  return (
    EDITOR_TYPE_COLOR_MAP[editorType as keyof typeof EDITOR_TYPE_COLOR_MAP] ||
    EDITOR_TYPE_COLOR_MAP.DEFAULT
  );
}
