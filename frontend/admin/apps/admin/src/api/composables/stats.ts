import type {
  statsservicev1_GetContentTrendResponse,
  statsservicev1_GetDashboardOverviewResponse,
  statsservicev1_GetInteractionStatsResponse,
  statsservicev1_GetLoginActivityResponse,
} from '#/api/generated/admin/service/v1';

import { useQuery, type UseQueryOptions } from '@tanstack/vue-query';

import { apiClient } from '#/api/client';

// 趋势图统一窗口:近 30 天(服务端 clamp 上限 90)
const TREND_DAYS = 30;

// ------------------------------
// 1. 分析页概览卡片(总量 + 近 7 天新增)
// ------------------------------
export function useGetDashboardOverview(
  options?: UseQueryOptions<statsservicev1_GetDashboardOverviewResponse, Error>,
) {
  return useQuery({
    queryKey: ['getDashboardOverview'],
    queryFn: () => apiClient.statsService.GetDashboardOverview({}),
    ...options,
  });
}

// ------------------------------
// 2. 内容增长趋势(近 30 天每日新增用户/帖子/评论)
// ------------------------------
export function useGetContentTrend(
  options?: UseQueryOptions<statsservicev1_GetContentTrendResponse, Error>,
) {
  return useQuery({
    queryKey: ['getContentTrend', TREND_DAYS],
    queryFn: () => apiClient.statsService.GetContentTrend({ days: TREND_DAYS }),
    ...options,
  });
}

// ------------------------------
// 3. 互动统计(点赞 TOP 榜 + 类型分布)
// ------------------------------
export function useGetInteractionStats(
  options?: UseQueryOptions<statsservicev1_GetInteractionStatsResponse, Error>,
) {
  return useQuery({
    queryKey: ['getInteractionStats'],
    queryFn: () => apiClient.statsService.GetInteractionStats({ topN: undefined }),
    ...options,
  });
}

// ------------------------------
// 4. 登录活跃趋势(近 30 天每日成功/失败)
// ------------------------------
export function useGetLoginActivity(
  options?: UseQueryOptions<statsservicev1_GetLoginActivityResponse, Error>,
) {
  return useQuery({
    queryKey: ['getLoginActivity', TREND_DAYS],
    queryFn: () => apiClient.statsService.GetLoginActivity({ days: TREND_DAYS }),
    ...options,
  });
}
