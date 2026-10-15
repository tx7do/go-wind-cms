import {
  useMutation,
  type UseMutationOptions,
  useQuery,
} from '@tanstack/react-query';
import {
  type contentservicev1_GetPageVersionRequest,
  type contentservicev1_GetPostVersionRequest,
  type contentservicev1_Page,
  type contentservicev1_PageVersion,
  type contentservicev1_Post,
  type contentservicev1_PostVersion,
  type contentservicev1_RestorePageVersionRequest,
  type contentservicev1_RestorePostVersionRequest,
} from '@/api/generated/admin/service/v1';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 内容版本（帖子/页面快照）：列表/详情/恢复
// ==============================

// ─── 帖子版本 ───

/** 列表查询仅暴露 enabled 开关（抽屉打开才拉取），避免 TanStack v5 要求补全 queryKey 的 options 类型 */
interface ListOptions {
  enabled?: boolean;
}

export function useListPostVersions(
  postId: number | null,
  options?: ListOptions,
) {
  return useQuery({
    queryKey: ['listPostVersions', postId],
    queryFn: () => {
      if (!postId) return Promise.reject(new Error('postId is required'));
      return apiClient.postService.ListVersions({ postId });
    },
    enabled: !!postId && (options?.enabled ?? true),
  });
}

export async function fetchPostVersion(req: contentservicev1_GetPostVersionRequest) {
  return queryClient.fetchQuery({
    queryKey: ['getPostVersion', req],
    queryFn: () => apiClient.postService.GetVersion(req),
    retry: 0,
  });
}

export function useRestorePostVersion(
  options?: UseMutationOptions<contentservicev1_Post, Error, contentservicev1_RestorePostVersionRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.postService.RestoreVersion(req),
    ...options,
  });
}

// ─── 页面版本 ───

export function useListPageVersions(
  pageId: number | null,
  options?: ListOptions,
) {
  return useQuery({
    queryKey: ['listPageVersions', pageId],
    queryFn: () => {
      if (!pageId) return Promise.reject(new Error('pageId is required'));
      return apiClient.pageService.ListVersions({ pageId });
    },
    enabled: !!pageId && (options?.enabled ?? true),
  });
}

export async function fetchPageVersion(req: contentservicev1_GetPageVersionRequest) {
  return queryClient.fetchQuery({
    queryKey: ['getPageVersion', req],
    queryFn: () => apiClient.pageService.GetVersion(req),
    retry: 0,
  });
}

export function useRestorePageVersion(
  options?: UseMutationOptions<contentservicev1_Page, Error, contentservicev1_RestorePageVersionRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.pageService.RestoreVersion(req),
    ...options,
  });
}

// ==============================
// 快照解析辅助：快照 JSON → 指定语言的标题/正文
// ==============================

export interface VersionSnapshot {
  title: string;
  content: string;
}

/**
 * 从版本快照 JSON 提取指定语言的标题/正文。
 * - post：translations[].content 即正文；page：translations[].title（正文在 sections）。
 * - 语言缺失回退第一条；解析失败返回 null（快照损坏时不抛错，由 UI 显示占位）。
 */
export function extractSnapshot(
  version: contentservicev1_PostVersion | contentservicev1_PageVersion | null | undefined,
  lang: string,
  kind: 'post' | 'page',
): VersionSnapshot | null {
  if (!version?.content) return null;
  try {
    const parsed = JSON.parse(version.content) as {
      translations?: { languageCode?: string; title?: string; content?: string }[];
    };
    const translations = parsed.translations ?? [];
    const tr = translations.find((x) => x.languageCode === lang) ?? translations[0];
    if (!tr) return null;
    return {
      title: tr.title ?? '',
      // 页面正文在 sections（区块整体替换恢复），此处只有标题可线性展示
      content: kind === 'post' ? (tr.content ?? '') : '',
    };
  } catch (error) {
    console.error('Failed to parse version snapshot:', error);
    return null;
  }
}

/** 页面快照的区块数量（列表/预览示意用） */
export function snapshotSectionCount(
  version: contentservicev1_PageVersion | null | undefined,
): number | undefined {
  if (!version?.content) return undefined;
  try {
    const parsed = JSON.parse(version.content) as { sections?: unknown[] };
    return Array.isArray(parsed.sections) ? parsed.sections.length : undefined;
  } catch (error) {
    console.error('Failed to parse page snapshot sections:', error);
    return undefined;
  }
}
