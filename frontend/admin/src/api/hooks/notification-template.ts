import { useMutation, useQueryClient, type UseMutationOptions } from '@tanstack/react-query';
import {
  type notifyservicev1_CreateNotificationTemplateRequest,
  type notifyservicev1_DeleteNotificationTemplateRequest,
  type notifyservicev1_ListNotificationTemplateResponse,
  type notifyservicev1_NotificationTemplate,
  type notifyservicev1_RenderNotificationTemplateRequest,
  type notifyservicev1_RenderNotificationTemplateResponse,
  type notifyservicev1_UpdateNotificationTemplateRequest,
} from '@/api/generated/admin/service/v1';
import { type PaginationQuery, queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 通知模板（可复用的标题/正文占位模板）
// ==============================

const LIST_KEY = 'listNotificationTemplates';

export function fetchListNotificationTemplates(
  query: PaginationQuery,
): Promise<notifyservicev1_ListNotificationTemplateResponse> {
  return queryClient.fetchQuery({
    queryKey: [LIST_KEY, query],
    queryFn: () => apiClient.notificationTemplateService.List(query.toRawParams()),
    retry: 0,
  });
}

export function useCreateNotificationTemplate(
  options?: UseMutationOptions<
    notifyservicev1_NotificationTemplate,
    Error,
    notifyservicev1_CreateNotificationTemplateRequest
  >,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationTemplateService.Create(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useUpdateNotificationTemplate(
  options?: UseMutationOptions<
    notifyservicev1_NotificationTemplate,
    Error,
    notifyservicev1_UpdateNotificationTemplateRequest
  >,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationTemplateService.Update(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useDeleteNotificationTemplate(
  options?: UseMutationOptions<Record<never, never>, Error, notifyservicev1_DeleteNotificationTemplateRequest>,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationTemplateService.Delete(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useRenderNotificationTemplate(
  options?: UseMutationOptions<
    notifyservicev1_RenderNotificationTemplateResponse,
    Error,
    notifyservicev1_RenderNotificationTemplateRequest
  >,
) {
  return useMutation({
    mutationFn: (req) => apiClient.notificationTemplateService.Render(req),
    ...options,
  });
}
