import { useMutation, type UseMutationOptions, useQuery, type UseQueryOptions } from '@tanstack/react-query';
import {
  type online_sessionservicev1_ForceLogoutSessionRequest,
  type online_sessionservicev1_ListOnlineSessionResponse,
} from '@/api/generated/admin/service/v1';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 在线会话（只读列表 + 强制下线；数据来自 Redis 会话元数据）
// ==============================

export type ListOnlineSessionsParams = {
  page: number;
  pageSize: number;
  keyword?: string;
};

export async function fetchListOnlineSessions(params: ListOnlineSessionsParams) {
  return queryClient.fetchQuery({
    queryKey: ['listOnlineSessions', params],
    queryFn: () =>
      apiClient.onlineSessionService.ListOnlineSession({
        keyword: params.keyword,
        page: params.page,
        pageSize: params.pageSize,
      }),
    retry: 0,
  });
}

export function useListMyOnlineSessions(
  currentJti: string | undefined,
  userId: number | undefined,
  options?: UseQueryOptions<online_sessionservicev1_ListOnlineSessionResponse, Error>,
) {
  return useQuery({
    queryKey: ['listMyOnlineSessions', userId],
    queryFn: () =>
      apiClient.onlineSessionService.ListMyOnlineSession({
        userId: userId ?? 0,
        currentJti: currentJti,
      }),
    enabled: !!userId,
    ...options,
  });
}

export function useForceLogoutSession(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.onlineSessionService.ForceLogoutSession>>,
    Error,
    online_sessionservicev1_ForceLogoutSessionRequest
  >,
) {
  return useMutation({
    mutationFn: (req: online_sessionservicev1_ForceLogoutSessionRequest) =>
      apiClient.onlineSessionService.ForceLogoutSession(req),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['listOnlineSessions'] }),
    ...options,
  });
}

export function useRevokeMyOnlineSession(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.onlineSessionService.RevokeMyOnlineSession>>,
    Error,
    { clientType?: string; jti: string }
  >,
) {
  return useMutation({
    mutationFn: ({ clientType, jti }: { clientType?: string; jti: string }) =>
      apiClient.onlineSessionService.RevokeMyOnlineSession({ clientType: clientType as never, jti }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['listMyOnlineSessions'] }),
    ...options,
  });
}
