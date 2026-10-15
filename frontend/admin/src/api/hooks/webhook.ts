import {
  useMutation,
  type UseMutationOptions,
  useQuery,
} from '@tanstack/react-query';
import {
  type webhookservicev1_CreateWebhookRequest as CreateWebhookRequest,
  type webhookservicev1_DeleteWebhookRequest as DeleteWebhookRequest,
  type webhookservicev1_GetWebhookRequest as GetWebhookRequest,
  type webhookservicev1_ListWebhookResponse,
  type webhookservicev1_TestWebhookPingRequest as TestWebhookPingRequest,
  type webhookservicev1_UpdateWebhookRequest as UpdateWebhookRequest,
  type webhookservicev1_Webhook as Webhook,
  type webhookservicev1_WebhookDelivery as WebhookDelivery,
} from '@/api/generated/admin/service/v1';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// Webhook 订阅管理（内容事件外发）
// ==============================

export function useListWebhooks(
  options?: UseMutationOptions<webhookservicev1_ListWebhookResponse, Error, void>,
) {
  return useMutation({
    // 订阅量级小：一次拉全量（pageSize 封顶 200），前端本地分页
    mutationFn: () =>
      apiClient.webhookService.List({
        page: 1,
        pageSize: 200,
        offset: 0,
        limit: 200,
        sorting: [],
      } as never),
    ...options,
  });
}

export async function fetchWebhook(req: GetWebhookRequest) {
  return queryClient.fetchQuery({
    queryKey: ['getWebhook', req],
    queryFn: () => apiClient.webhookService.Get(req),
    retry: 0,
  });
}

export function useCreateWebhook(
  options?: UseMutationOptions<Webhook, Error, CreateWebhookRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.webhookService.Create(req),
    ...options,
  });
}

export function useUpdateWebhook(
  options?: UseMutationOptions<Webhook, Error, UpdateWebhookRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.webhookService.Update(req),
    ...options,
  });
}

export function useDeleteWebhook(
  options?: UseMutationOptions<Record<string, never>, Error, DeleteWebhookRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.webhookService.Delete(req),
    ...options,
  });
}

export function useTestWebhookPing(
  options?: UseMutationOptions<WebhookDelivery, Error, TestWebhookPingRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.webhookService.TestPing(req),
    ...options,
  });
}

export function useListWebhookDeliveries(
  params: { webhookId: number | null; limit?: number },
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ['listWebhookDeliveries', params.webhookId, params.limit ?? 50],
    queryFn: () => {
      if (!params.webhookId) {
        return Promise.reject(new Error('webhookId is required'));
      }
      return apiClient.webhookService.ListDeliveries({
        webhookId: params.webhookId,
        limit: params.limit ?? 50,
      });
    },
    enabled: !!params.webhookId && (options?.enabled ?? true),
  });
}

/** 事件类型清单（admin 表单选项与 core/task 包常量对齐） */
export const WEBHOOK_EVENT_OPTIONS = [
  { value: 'post.created', label: 'post.created' },
  { value: 'post.updated', label: 'post.updated' },
  { value: 'post.published', label: 'post.published' },
  { value: 'post.deleted', label: 'post.deleted' },
  { value: 'page.created', label: 'page.created' },
  { value: 'page.updated', label: 'page.updated' },
  { value: 'page.published', label: 'page.published' },
  { value: 'page.deleted', label: 'page.deleted' },
] as const;
