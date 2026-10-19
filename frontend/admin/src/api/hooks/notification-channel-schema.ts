import { useQuery, type UseQueryOptions } from '@tanstack/react-query';
import {
  type notifyservicev1_ChannelSchema,
  type notifyservicev1_ListChannelSchemasResponse,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';

// ==============================
// 扩展渠道 schema（配置表单动态渲染的数据源）
// ==============================

export function useListChannelSchemas(
  options?: UseQueryOptions<notifyservicev1_ListChannelSchemasResponse, Error>,
) {
  return useQuery({
    queryKey: ['listChannelSchemas'],
    queryFn: () => apiClient.notificationChannelService.ListChannelSchemas({}),
    staleTime: 5 * 60 * 1000,
    ...options,
  });
}

export type { notifyservicev1_ChannelSchema as ChannelSchema };
