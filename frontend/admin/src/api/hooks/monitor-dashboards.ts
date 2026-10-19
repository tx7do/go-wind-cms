import { useQuery, type UseQueryOptions } from '@tanstack/react-query';
import {
  type server_monitorservicev1_ServerMonitorInfo,
  type redis_cacheservicev1_RedisCacheMonitorInfo,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';

// ==============================
// 监控仪表盘（只读聚合视图）
// ==============================

/** 服务监控：Go 运行时 / 数据库连接池 / 主机信息（10 秒自动刷新由页面传 refetchInterval）。 */
export function useServerMonitorInfo(
  options?: Omit<
    UseQueryOptions<server_monitorservicev1_ServerMonitorInfo, Error>,
    'queryKey' | 'queryFn'
  >,
) {
  return useQuery({
    queryKey: ['serverMonitorInfo'],
    queryFn: () => apiClient.serverMonitorService.Get({}),
    ...options,
  });
}

/** Redis 缓存监控：INFO 分区 / DBSIZE / SLOWLOG。 */
export function useRedisCacheMonitorInfo(
  options?: Omit<
    UseQueryOptions<redis_cacheservicev1_RedisCacheMonitorInfo, Error>,
    'queryKey' | 'queryFn'
  >,
) {
  return useQuery({
    queryKey: ['redisCacheMonitorInfo'],
    queryFn: () => apiClient.redisCacheMonitorService.Get({}),
    ...options,
  });
}
