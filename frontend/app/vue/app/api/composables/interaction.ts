import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationOptions,
  type UseQueryOptions,
} from '@tanstack/vue-query';
import { computed, toValue, type MaybeRefOrGetter } from 'vue';
import {
  type interactionservicev1_LikeResponse,
  type interactionservicev1_WatchResponse,
  type interactionservicev1_GetInteractionStatusResponse,
  type interactionservicev1_GetCountsResponse,
  type interactionservicev1_TargetType,
  type interactionservicev1_CounterMetric,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';

// ==============================
// 交互服务 API（点赞 / 收藏 / 交互状态 / 计数）
// ==============================
//
// 说明：所有写操作（Like/Unlike/Watch/Unwatch）的服务端鉴权由后端 InteractionService
// 强制，viewer 身份取自鉴权上下文，客户端无需也不可传 user_id。
// 点赞计数统一存于 interaction_counter 表（post/comment 表的散落计数列已移除），
// 列表/详情渲染计数改读 GetCounts；点赞/取消后由 LikeResponse.likeCount 乐观更新。

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
 * 批量查询指定目标的计数（如点赞数）。供列表/详情渲染计数展示。
 */
export async function getCounts(
  targetType: interactionservicev1_TargetType,
  targetIds: number[],
  metrics: interactionservicev1_CounterMetric[],
) {
  return apiClient.interactionService.GetCounts({
    targetType,
    targetIds,
    metrics,
  });
}

// ==============================
// 交互状态查询 Hook
// ==============================

/**
 * useInteractionStatus —— 批量查询当前 viewer 对指定目标的 {liked, watched} 状态。
 * 仅当 targetIds 非空时启用查询。用于渲染列表页中点赞/收藏按钮的初始态。
 *
 * 说明：Vue 端 targetIds 接受 MaybeRefOrGetter，以响应列表数据变化（如分页切换）
 * 自动重新发起查询。queryKey 与 enabled 均以 computed 形式注册，保证响应式追踪。
 */
export function useInteractionStatus(
  targetType: interactionservicev1_TargetType,
  targetIds: MaybeRefOrGetter<number[]>,
  options?: Omit<
    UseQueryOptions<interactionservicev1_GetInteractionStatusResponse, Error>,
    'queryKey' | 'queryFn' | 'enabled'
  >,
) {
  return useQuery({
    queryKey: computed(() => [
      'interaction-status',
      targetType,
      toValue(targetIds),
    ]),
    queryFn: () => {
      const ids = toValue(targetIds);
      return getInteractionStatus(targetType, ids);
    },
    enabled: computed(() => toValue(targetIds).length > 0),
    staleTime: 0,
    ...options,
  });
}

// ==============================
// 计数查询 Hook
// ==============================

/**
 * useGetCounts —— 批量查询指定目标的计数（如点赞数）。
 * 仅当 targetIds 非空且 metrics 非空时启用查询。用于列表/详情渲染计数展示。
 * 响应中未出现的 (target, metric) 组合按 0 兜底。
 *
 * 说明：Vue 端 targetIds / metrics 接受 MaybeRefOrGetter，以响应列表数据变化
 * 自动重新发起查询。queryKey 与 enabled 均以 computed 形式注册，保证响应式追踪。
 */
export function useGetCounts(
  targetType: interactionservicev1_TargetType,
  targetIds: MaybeRefOrGetter<number[]>,
  metrics: MaybeRefOrGetter<interactionservicev1_CounterMetric[]>,
  options?: Omit<
    UseQueryOptions<interactionservicev1_GetCountsResponse, Error>,
    'queryKey' | 'queryFn' | 'enabled'
  >,
) {
  return useQuery({
    queryKey: computed(() => [
      'interaction-counts',
      targetType,
      toValue(targetIds),
      toValue(metrics),
    ]),
    queryFn: () => {
      const ids = toValue(targetIds);
      const m = toValue(metrics);
      return getCounts(targetType, ids, m);
    },
    enabled: computed(
      () => toValue(targetIds).length > 0 && toValue(metrics).length > 0,
    ),
    staleTime: 0,
    ...options,
  });
}

/**
 * extractCount —— 从 GetCounts 响应中提取单个 (targetId, metric) 的计数。
 * 未出现则返回 0。
 */
export function extractCount(
  resp: interactionservicev1_GetCountsResponse | undefined,
  targetId: number,
  metric: interactionservicev1_CounterMetric,
): number {
  if (!resp?.counts) return 0;
  const cm = resp.counts[String(targetId)];
  if (!cm?.counts) return 0;
  for (const mc of cm.counts) {
    if (mc.metric === metric && mc.count !== undefined) {
      return mc.count;
    }
  }
  return 0;
}

// ==============================
// 点赞 / 收藏 写操作 Hooks
// ==============================
//
// 写操作成功后主动 invalidate 对应的 interaction-status / interaction-counts 查询，
// 使列表中相关按钮态与计数立即刷新。

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
      qc.invalidateQueries({queryKey: ['interaction-counts', vars.targetType]});
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
      qc.invalidateQueries({queryKey: ['interaction-counts', vars.targetType]});
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
      // post 列表/详情的收藏按钮态与收藏计数
      qc.invalidateQueries({queryKey: ['interaction-status', 'TARGET_TYPE_POST' as interactionservicev1_TargetType]});
      qc.invalidateQueries({queryKey: ['interaction-counts', 'TARGET_TYPE_POST' as interactionservicev1_TargetType]});
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
      // 收藏列表需刷新
      qc.invalidateQueries({queryKey: ['watched-posts']});
      // post 列表/详情的收藏按钮态与收藏计数
      qc.invalidateQueries({queryKey: ['interaction-status', 'TARGET_TYPE_POST' as interactionservicev1_TargetType]});
      qc.invalidateQueries({queryKey: ['interaction-counts', 'TARGET_TYPE_POST' as interactionservicev1_TargetType]});
    },
    ...options,
  });
}
