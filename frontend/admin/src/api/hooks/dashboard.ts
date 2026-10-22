import { useQuery, type UseQueryOptions } from '@tanstack/react-query';
import {
  type statsservicev1_GetAiUsageTrendResponse,
  type statsservicev1_GetApiCallTrendResponse,
  type statsservicev1_GetCommentStatusDistributionResponse,
  type statsservicev1_GetContentTrendResponse,
  type statsservicev1_GetDashboardOverviewResponse,
  type statsservicev1_GetInteractionStatsResponse,
  type statsservicev1_GetLoginActivityResponse,
  type statsservicev1_GetMediaAssetTypeDistributionResponse,
  type statsservicev1_GetPostStatusDistributionResponse,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';

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
// 4. 登录活跃趋势(近 30 天每日成功/失败 + 活跃用户)
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

// ------------------------------
// 5. AI 用量趋势(近 30 天每日调用次数/Token 消耗)
// ------------------------------
export function useGetAiUsageTrend(
  options?: UseQueryOptions<statsservicev1_GetAiUsageTrendResponse, Error>,
) {
  return useQuery({
    queryKey: ['getAiUsageTrend', TREND_DAYS],
    queryFn: () => apiClient.statsService.GetAiUsageTrend({ days: TREND_DAYS }),
    ...options,
  });
}

// ------------------------------
// 6. API 调用量趋势(近 30 天每日网关请求)
// ------------------------------
export function useGetApiCallTrend(
  options?: UseQueryOptions<statsservicev1_GetApiCallTrendResponse, Error>,
) {
  return useQuery({
    queryKey: ['getApiCallTrend', TREND_DAYS],
    queryFn: () => apiClient.statsService.GetApiCallTrend({ days: TREND_DAYS }),
    ...options,
  });
}

// ------------------------------
// 7. 帖子状态分布(排除回收站)
// ------------------------------
export function useGetPostStatusDistribution(
  options?: UseQueryOptions<statsservicev1_GetPostStatusDistributionResponse, Error>,
) {
  return useQuery({
    queryKey: ['getPostStatusDistribution'],
    queryFn: () => apiClient.statsService.GetPostStatusDistribution({}),
    ...options,
  });
}

// ------------------------------
// 8. 评论状态分布(NULL 状态按已通过计)
// ------------------------------
export function useGetCommentStatusDistribution(
  options?: UseQueryOptions<statsservicev1_GetCommentStatusDistributionResponse, Error>,
) {
  return useQuery({
    queryKey: ['getCommentStatusDistribution'],
    queryFn: () => apiClient.statsService.GetCommentStatusDistribution({}),
    ...options,
  });
}

// ------------------------------
// 9. 媒体类型分布(NULL 类型按其他计)
// ------------------------------
export function useGetMediaAssetTypeDistribution(
  options?: UseQueryOptions<statsservicev1_GetMediaAssetTypeDistributionResponse, Error>,
) {
  return useQuery({
    queryKey: ['getMediaAssetTypeDistribution'],
    queryFn: () => apiClient.statsService.GetMediaAssetTypeDistribution({}),
    ...options,
  });
}
