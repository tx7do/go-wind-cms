import { useMutation, useQuery, useQueryClient, type UseMutationOptions, type UseQueryOptions } from '@tanstack/react-query';
import {
  type notifyservicev1_CreateNotificationChannelRequest,
  type notifyservicev1_DeleteNotificationChannelRequest,
  type notifyservicev1_ListNotificationChannelResponse,
  type notifyservicev1_NotificationChannel,
  type notifyservicev1_SendTestEmailRequest,
  type notifyservicev1_SendTestEmailResponse,
  type notifyservicev1_UpdateNotificationChannelRequest,
} from '@/api/generated/admin/service/v1';
import { type PaginationQuery, queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 通知渠道（平台级配置）
// ==============================

const LIST_KEY = 'listNotificationChannels';

export function useListNotificationChannels(
  query: PaginationQuery,
  options?: UseQueryOptions<notifyservicev1_ListNotificationChannelResponse, Error>,
) {
  return useQuery({
    queryKey: ['listNotificationChannels', query],
    queryFn: () => apiClient.notificationChannelService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListNotificationChannels(query: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: [LIST_KEY, query],
    queryFn: () => apiClient.notificationChannelService.List(query.toRawParams()),
    retry: 0,
  });
}

export function useCreateNotificationChannel(
  options?: UseMutationOptions<
    notifyservicev1_NotificationChannel,
    Error,
    notifyservicev1_CreateNotificationChannelRequest
  >,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationChannelService.Create(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useUpdateNotificationChannel(
  options?: UseMutationOptions<
    notifyservicev1_NotificationChannel,
    Error,
    notifyservicev1_UpdateNotificationChannelRequest
  >,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationChannelService.Update(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useDeleteNotificationChannel(
  options?: UseMutationOptions<Record<never, never>, Error, notifyservicev1_DeleteNotificationChannelRequest>,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationChannelService.Delete(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useSendTestEmail(
  options?: UseMutationOptions<notifyservicev1_SendTestEmailResponse, Error, notifyservicev1_SendTestEmailRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.notificationChannelService.SendTestEmail(req),
    ...options,
  });
}
