import type {
  commentservicev1_Comment,
  commentservicev1_DeleteCommentRequest,
  commentservicev1_GetCommentRequest,
  commentservicev1_ListCommentResponse,
} from '#/api/generated/admin/service/v1';

import { computed } from 'vue';

import { $t } from '@vben/locales';

import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/vue-query';

import { apiClient } from '#/api/client';
import { queryClient } from '#/plugins/vue-query';
import { makeUpdateMask, type PaginationQuery } from '#/transport/rest';

// ==============================
// 评论管理
// ==============================

export function useListComments(
  query: PaginationQuery,
  options?: UseQueryOptions<commentservicev1_ListCommentResponse, Error>,
) {
  return useQuery({
    queryKey: ['listComments', query],
    queryFn: () => apiClient.commentService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListComments(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listComments', params],
    queryFn: () => apiClient.commentService.List(params.toRawParams()),
    retry: 0,
  });
}

export function useGetComment(
  req: commentservicev1_GetCommentRequest,
  options?: UseQueryOptions<commentservicev1_Comment, Error>,
) {
  return useQuery({
    queryKey: ['getComment', req],
    queryFn: () => apiClient.commentService.Get(req),
    ...options,
  });
}

export function useCreateComment(
  options?: UseMutationOptions<
    commentservicev1_Comment,
    Error,
    Record<string, any>
  >,
) {
  return useMutation({
    mutationFn: (values) =>
      apiClient.commentService.Create({
        data: { ...values } as commentservicev1_Comment,
      }),
    ...options,
  });
}

export function useUpdateComment(
  options?: UseMutationOptions<
    commentservicev1_Comment,
    Error,
    { id: number; values: Record<string, any> }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.commentService.Update({
        id,
        data: { ...values } as commentservicev1_Comment,
        updateMask: makeUpdateMask(Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteComment(
  options?: UseMutationOptions<
    object,
    Error,
    commentservicev1_DeleteCommentRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.commentService.Delete(data),
    ...options,
  });
}

// ==============================
// 评论枚举与工具函数
// ==============================

export const commentStatusList = computed(() => [
  { value: 'STATUS_PENDING', label: $t('enum.comment.status.STATUS_PENDING') },
  {
    value: 'STATUS_APPROVED',
    label: $t('enum.comment.status.STATUS_APPROVED'),
  },
  {
    value: 'STATUS_REJECTED',
    label: $t('enum.comment.status.STATUS_REJECTED'),
  },
  {
    value: 'STATUS_SPAM',
    label: $t('enum.comment.status.STATUS_SPAM'),
  },
]);

export function commentStatusToName(
  status: commentservicev1_Comment['status'],
) {
  const values = commentStatusList.value;
  const matchedItem = values.find((item) => item.value === status);
  return matchedItem ? matchedItem.label : '';
}

const COMMENT_STATUS_COLOR_MAP = {
  STATUS_PENDING: '#60a5fa',
  STATUS_APPROVED: '#22c55e',
  STATUS_REJECTED: '#f97316',
  STATUS_SPAM: '#ef4444',
  DEFAULT: '#94a3b8',
} as const;

export function commentStatusToColor(
  status: commentservicev1_Comment['status'],
) {
  return (
    COMMENT_STATUS_COLOR_MAP[status as keyof typeof COMMENT_STATUS_COLOR_MAP] ||
    COMMENT_STATUS_COLOR_MAP.DEFAULT
  );
}

export const commentContentTypeList = computed(() => [
  {
    value: 'CONTENT_TYPE_POST',
    label: $t('enum.comment.contentType.CONTENT_TYPE_POST'),
  },
  {
    value: 'CONTENT_TYPE_PAGE',
    label: $t('enum.comment.contentType.CONTENT_TYPE_PAGE'),
  },
  {
    value: 'CONTENT_TYPE_PRODUCT',
    label: $t('enum.comment.contentType.CONTENT_TYPE_PRODUCT'),
  },
]);

export function commentContentTypeToName(
  contentType: commentservicev1_Comment['contentType'],
) {
  const values = commentContentTypeList.value;
  const matchedItem = values.find((item) => item.value === contentType);
  return matchedItem ? matchedItem.label : '';
}

const COMMENT_CONTENT_TYPE_COLOR_THEME = {
  light: {
    CONTENT_TYPE_POST: '#2563eb',
    CONTENT_TYPE_PAGE: '#7c3aed',
    CONTENT_TYPE_PRODUCT: '#ea580c',
  },
  dark: {
    CONTENT_TYPE_POST: '#3b82f6',
    CONTENT_TYPE_PAGE: '#8b5cf6',
    CONTENT_TYPE_PRODUCT: '#fdba74',
  },
} as const;

export function commentContentTypeToColor(
  contentType: commentservicev1_Comment['contentType'],
  theme: 'dark' | 'light' = 'light',
): string {
  const colorMap = COMMENT_CONTENT_TYPE_COLOR_THEME[theme];
  return (
    colorMap[contentType as keyof typeof colorMap] || colorMap.CONTENT_TYPE_POST
  );
}

export const commentAuthorTypeList = computed(() => [
  {
    value: 'AUTHOR_TYPE_GUEST',
    label: $t('enum.comment.authorType.AUTHOR_TYPE_GUEST'),
  },
  {
    value: 'AUTHOR_TYPE_USER',
    label: $t('enum.comment.authorType.AUTHOR_TYPE_USER'),
  },
  {
    value: 'AUTHOR_TYPE_ADMIN',
    label: $t('enum.comment.authorType.AUTHOR_TYPE_ADMIN'),
  },
  {
    value: 'AUTHOR_TYPE_MODERATOR',
    label: $t('enum.comment.authorType.AUTHOR_TYPE_MODERATOR'),
  },
]);

export function commentAuthorTypeToName(
  authorType: commentservicev1_Comment['authorType'],
) {
  const values = commentAuthorTypeList.value;
  const matchedItem = values.find((item) => item.value === authorType);
  return matchedItem ? matchedItem.label : '';
}

const COMMENT_AUTHOR_TYPE_COLOR_MAP = {
  AUTHOR_TYPE_GUEST: '#64748b',
  AUTHOR_TYPE_USER: '#10b981',
  AUTHOR_TYPE_ADMIN: '#3b82f6',
  AUTHOR_TYPE_MODERATOR: '#f97316',
  DEFAULT: '#94a3b8',
} as const;

export function commentAuthorTypeToColor(
  authorType: commentservicev1_Comment['authorType'],
) {
  return (
    COMMENT_AUTHOR_TYPE_COLOR_MAP[
      authorType as keyof typeof COMMENT_AUTHOR_TYPE_COLOR_MAP
    ] || COMMENT_AUTHOR_TYPE_COLOR_MAP.DEFAULT
  );
}
