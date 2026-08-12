import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationOptions,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type interactionservicev1_LikeResponse,
  type interactionservicev1_WatchResponse,
  type interactionservicev1_GetInteractionStatusResponse,
  type interactionservicev1_TargetType,
  type contentservicev1_ListPostResponse,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { queryClient } from '@/core';

// ==============================
// 交互服务 API（点赞 / 收藏）
// ==============================
//
// 说明：所有写操作（Like/Unlike/Watch/Unwatch）的服务端鉴权由后端 InteractionService
// 强制，viewer 身份取自鉴权上下文，客户端无需也不可传 user_id。
// 计数（likes）由后端在 ledger 写入的同一事务内递增/递减并回读，前端直读返回值即可。

/**
 * 点赞（post 或 comment）。幂等。
 */
export async function likePost(targetType: interactionservicev1_TargetType, targetId: number) {
  return apiClient.interactionService.Like({
    targetType,
    targetId,
  });
}

/**
 * 取消点赞。幂等。
 */
export async function unlikePost(targetType: interactionservicev1_TargetType, targetId: number) {
  return apiClient.interactionService.Unlike({
    targetType,
    targetId,
  });
}

/**
 * 收藏 post。幂等。
 */
export async function watchPost(postId: number) {
  return apiClient.interactionService.Watch({
    postId,
  });
}

/**
 * 取消收藏 post。幂等。
 */
export async function unwatchPost(postId: number) {
  return apiClient.interactionService.Unwatch({
    postId,
  });
}

/**
 * 批量查询当前 viewer 对指定目标的交互状态（供前端渲染按钮态）。
 */
export async function getInteractionStatus(
  targetType: interactionservicev1_TargetType,
  targetIds: number[],
) {
  return apiClient.interactionService.GetInteractionStatus({
    targetType,
    targetIds,
  });
}

/**
 * 列出当前 viewer 收藏的 post。
 */
export async function listWatchedPosts(page?: number, pageSize?: number) {
  return apiClient.interactionService.ListWatchedPosts({
    page,
    pageSize,
    sorting: undefined,
  });
}

// ==============================
// 交互状态查询 Hook
// ==============================

/**
 * useInteractionStatus —— 批量查询当前 viewer 对指定目标的 {liked, watched} 状态。
 * 仅当 targetIds 非空时启用查询。用于渲染列表页中点赞/收藏按钮的初始态。
 */
export function useInteractionStatus(
  targetType: interactionservicev1_TargetType,
  targetIds: number[],
  options?: Omit<
    UseQueryOptions<interactionservicev1_GetInteractionStatusResponse, Error>,
    'queryKey' | 'queryFn' | 'enabled'
  >,
) {
  return useQuery({
    queryKey: ['interaction-status', targetType, targetIds],
    queryFn: () => getInteractionStatus(targetType, targetIds),
    enabled: targetIds.length > 0,
    staleTime: 0,
    ...options,
  });
}

// ==============================
// 点赞 / 收藏 写操作 Hooks
// ==============================
//
// 写操作成功后主动 invalidate 对应的 interaction-status 查询，使列表中相关按钮态
// 立即刷新。点赞同时 invalidate post 详情（likes 计数变化），收藏同时 invalidate
// watched-posts 列表。

export function useLike(
  options?: UseMutationOptions<
    interactionservicev1_LikeResponse,
    Error,
    { targetType: interactionservicev1_TargetType; targetId: number }
  >,
) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({targetType, targetId}) => likePost(targetType, targetId),
    onSettled: (_d, _e, vars) => {
      qc.invalidateQueries({queryKey: ['interaction-status', vars.targetType]});
    },
    ...options,
  });
}

export function useUnlike(
  options?: UseMutationOptions<
    interactionservicev1_LikeResponse,
    Error,
    { targetType: interactionservicev1_TargetType; targetId: number }
  >,
) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({targetType, targetId}) => unlikePost(targetType, targetId),
    onSettled: (_d, _e, vars) => {
      qc.invalidateQueries({queryKey: ['interaction-status', vars.targetType]});
    },
    ...options,
  });
}

export function useWatch(
  options?: UseMutationOptions<
    interactionservicev1_WatchResponse,
    Error,
    { postId: number }
  >,
) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({postId}) => watchPost(postId),
    onSettled: () => {
      // 收藏列表需刷新
      qc.invalidateQueries({queryKey: ['watched-posts']});
      // post 列表/详情的收藏按钮态
      qc.invalidateQueries({queryKey: ['interaction-status', 'TARGET_TYPE_POST' as interactionservicev1_TargetType]});
    },
    ...options,
  });
}

export function useUnwatch(
  options?: UseMutationOptions<
    interactionservicev1_WatchResponse,
    Error,
    { postId: number }
  >,
) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({postId}) => unwatchPost(postId),
    onSettled: () => {
      qc.invalidateQueries({queryKey: ['watched-posts']});
      qc.invalidateQueries({queryKey: ['interaction-status', 'TARGET_TYPE_POST' as interactionservicev1_TargetType]});
    },
    ...options,
  });
}

// ==============================
// 收藏列表 Hook
// ==============================

/**
 * useWatchedPosts —— 列出当前 viewer 收藏的 post。
 */
export function useWatchedPosts(
  page: number,
  pageSize: number,
  options?: Omit<
    UseQueryOptions<contentservicev1_ListPostResponse, Error>,
    'queryKey' | 'queryFn'
  >,
) {
  return useQuery({
    queryKey: ['watched-posts', page, pageSize],
    queryFn: () => listWatchedPosts(page, pageSize),
    ...options,
  });
}

// 防止 queryClient 未使用导入被 tree-shake 报错（与 comment.ts 一致）
void queryClient;
