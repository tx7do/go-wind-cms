import { useMutation, useQueryClient, type UseMutationOptions } from '@tanstack/react-query';
import {
  type notifyservicev1_CreateNotificationRuleRequest,
  type notifyservicev1_DeleteNotificationRuleRequest,
  type notifyservicev1_ListNotificationRuleResponse,
  type notifyservicev1_NotificationRule,
  type notifyservicev1_TestDispatchNotificationRequest,
  type notifyservicev1_TestDispatchNotificationResponse,
  type notifyservicev1_UpdateNotificationRuleRequest,
} from '@/api/generated/admin/service/v1';
import { type PaginationQuery, queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 通知路由规则（事件类型 → 渠道）
// ==============================

const LIST_KEY = 'listNotificationRules';

export function fetchListNotificationRules(
  query: PaginationQuery,
): Promise<notifyservicev1_ListNotificationRuleResponse> {
  return queryClient.fetchQuery({
    queryKey: [LIST_KEY, query],
    queryFn: () => apiClient.notificationRuleService.List(query.toRawParams()),
    retry: 0,
  });
}

export function useCreateNotificationRule(
  options?: UseMutationOptions<
    notifyservicev1_NotificationRule,
    Error,
    notifyservicev1_CreateNotificationRuleRequest
  >,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationRuleService.Create(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useUpdateNotificationRule(
  options?: UseMutationOptions<
    notifyservicev1_NotificationRule,
    Error,
    notifyservicev1_UpdateNotificationRuleRequest
  >,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationRuleService.Update(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useDeleteNotificationRule(
  options?: UseMutationOptions<Record<never, never>, Error, notifyservicev1_DeleteNotificationRuleRequest>,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.notificationRuleService.Delete(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

/**
 * 测试投递：按这条规则当场走一遍完整投递链。
 *
 * 异步规则回的是 SENDING 而不是 SENT —— 结论落在队列那一趟，页面上要按台账状态如实显示，
 * 不能把它翻译成"发送成功"。
 */
export function useTestDispatchNotification(
  options?: UseMutationOptions<
    notifyservicev1_TestDispatchNotificationResponse,
    Error,
    notifyservicev1_TestDispatchNotificationRequest
  >,
) {
  return useMutation({
    mutationFn: (req) => apiClient.notificationRuleService.TestDispatch(req),
    ...options,
  });
}
