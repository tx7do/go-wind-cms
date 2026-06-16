import type {
  contentservicev1_DeletePostRequest,
  contentservicev1_EditorType,
  contentservicev1_GetPostRequest,
  contentservicev1_ListPostResponse,
  contentservicev1_Post,
} from '#/api/generated/admin/service/v1';

import { computed } from 'vue';

import { $t } from '@vben/locales';

import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/vue-query';

import { EditorType } from '#/adapter/component/Editor';
import { apiClient } from '#/api/client';
import { queryClient } from '#/plugins/vue-query';
import { makeUpdateMask, type PaginationQuery } from '#/transport/rest';

// ==============================
// 帖子管理
// ==============================

export function useListPosts(
  query: PaginationQuery,
  options?: UseQueryOptions<contentservicev1_ListPostResponse, Error>,
) {
  return useQuery({
    queryKey: ['listPosts', query],
    queryFn: () => apiClient.postService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListPosts(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listPosts', params],
    queryFn: () => apiClient.postService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetPost(
  req: contentservicev1_GetPostRequest,
  options?: UseQueryOptions<contentservicev1_Post, Error>,
) {
  return useQuery({
    queryKey: ['getPost', req],
    queryFn: () => apiClient.postService.Get(req),
    ...options,
  });
}

export function useCreatePost(
  options?: UseMutationOptions<
    contentservicev1_Post,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.postService.Create({
        data: { ...values } as contentservicev1_Post,
      }),
    ...options,
  });
}

export function useUpdatePost(
  options?: UseMutationOptions<
    contentservicev1_Post,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.postService.Update({
        id,
        data: { ...values } as contentservicev1_Post,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeletePost(
  options?: UseMutationOptions<
    object,
    Error,
    contentservicev1_DeletePostRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.postService.Delete(data),
    ...options,
  });
}

// ==============================
// 帖子枚举与工具函数
// ==============================

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

export function postStatusToName(status: contentservicev1_Post['status']) {
  const values = postStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const POST_STATUS_COLOR_MAP = {
  POST_STATUS_DRAFT: '#7c3aed',
  POST_STATUS_PUBLISHED: '#10b981',
  POST_STATUS_SCHEDULED: '#f59e0b',
  POST_STATUS_TRASHED: '#ef4444',
  DEFAULT: '#94a3b8',
} as const;

export function postStatusToColor(status: contentservicev1_Post['status']) {
  return (
    POST_STATUS_COLOR_MAP[status as keyof typeof POST_STATUS_COLOR_MAP] ||
    POST_STATUS_COLOR_MAP.DEFAULT
  );
}

// 编辑器类型
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
    value: 'EDITOR_TYPE_VISUAL_BUILDER',
    label: $t('enum.editorType.EDITOR_TYPE_VISUAL_BUILDER'),
  },
]);

export function editorTypeToName(editorType: contentservicev1_EditorType) {
  const values = editorTypeList.value;
  const matchedItem = values.find((item) => item.value === editorType);
  return matchedItem ? matchedItem.label : '';
}

const EDITOR_TYPE_COLOR_MAP = {
  EDITOR_TYPE_MARKDOWN: '#4f46e5',
  EDITOR_TYPE_RICH_TEXT: '#14b8a6',
  EDITOR_TYPE_JSON_BLOCK: '#d946ef',
  EDITOR_TYPE_PLAIN_TEXT: '#64748b',
  EDITOR_TYPE_CODE: '#0ea5e9',
  EDITOR_TYPE_VISUAL_BUILDER: '#06b6d4',
  DEFAULT: '#94a3b8',
} as const;

export function editorTypeToColor(editorType: contentservicev1_EditorType) {
  return (
    EDITOR_TYPE_COLOR_MAP[editorType as keyof typeof EDITOR_TYPE_COLOR_MAP] ||
    EDITOR_TYPE_COLOR_MAP.DEFAULT
  );
}

export function convertToUIEditorType(
  editorType: contentservicev1_EditorType | undefined,
): EditorType {
  switch (editorType) {
    case 'EDITOR_TYPE_CODE': {
      return EditorType.CODE;
    }
    case 'EDITOR_TYPE_JSON_BLOCK': {
      return EditorType.JSON;
    }
    case 'EDITOR_TYPE_MARKDOWN': {
      return EditorType.MARKDOWN;
    }
    case 'EDITOR_TYPE_PLAIN_TEXT': {
      return EditorType.PLAIN_TEXT;
    }
    case 'EDITOR_TYPE_RICH_TEXT': {
      return EditorType.RICH_TEXT;
    }
    case 'EDITOR_TYPE_VISUAL_BUILDER': {
      return EditorType.VISUAL_BUILDER;
    }
  }
  return EditorType.MARKDOWN;
}

export function convertToEditorType(
  editorType: EditorType | undefined,
): contentservicev1_EditorType {
  switch (editorType) {
    case EditorType.CODE: {
      return 'EDITOR_TYPE_CODE';
    }
    case EditorType.JSON: {
      return 'EDITOR_TYPE_JSON_BLOCK';
    }
    case EditorType.MARKDOWN: {
      return 'EDITOR_TYPE_MARKDOWN';
    }
    case EditorType.PLAIN_TEXT: {
      return 'EDITOR_TYPE_PLAIN_TEXT';
    }
    case EditorType.RICH_TEXT: {
      return 'EDITOR_TYPE_RICH_TEXT';
    }
    case EditorType.VISUAL_BUILDER: {
      return 'EDITOR_TYPE_VISUAL_BUILDER';
    }
  }
  return 'EDITOR_TYPE_MARKDOWN';
}

// 编辑器类型选项列表
export const editorTypeOptions = computed(() => [
  {
    label: $t('enum.editorType.EDITOR_TYPE_MARKDOWN'),
    value: EditorType.MARKDOWN,
  },
  {
    label: $t('enum.editorType.EDITOR_TYPE_RICH_TEXT'),
    value: EditorType.RICH_TEXT,
  },
  {
    label: $t('enum.editorType.EDITOR_TYPE_JSON_BLOCK'),
    value: EditorType.JSON,
  },
  { label: $t('enum.editorType.EDITOR_TYPE_CODE'), value: EditorType.CODE },
  {
    label: $t('enum.editorType.EDITOR_TYPE_PLAIN_TEXT'),
    value: EditorType.PLAIN_TEXT,
  },
]);
