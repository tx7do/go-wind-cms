import { useMutation, type UseMutationOptions, useQuery, type UseQueryOptions } from '@tanstack/react-query';
import {
  type monitor_alertservicev1_CreateMonitorAlertRuleRequest,
  type monitor_alertservicev1_EvaluateMonitorAlertsResponse,
  type monitor_alertservicev1_ListMonitorAlertRuleResponse,
  type monitor_alertservicev1_UpdateMonitorAlertRuleRequest,
} from '@/api/generated/admin/service/v1';
import { PaginationQuery, queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 监控告警规则（指标阈值 → 触发通知，联动通知域）
// ==============================

export function useListMonitorAlertRules(
  query: PaginationQuery,
  options?: UseQueryOptions<monitor_alertservicev1_ListMonitorAlertRuleResponse, Error>,
) {
  return useQuery({
    queryKey: ['listMonitorAlertRules', query],
    queryFn: () => apiClient.monitorAlertService.ListMonitorAlertRule(query.toRawParams()),
    ...options,
  });
}

export async function fetchListMonitorAlertRules(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listMonitorAlertRules', params],
    queryFn: () => apiClient.monitorAlertService.ListMonitorAlertRule(params.toRawParams()),
    retry: 0,
  });
}

export function useCreateMonitorAlertRule(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.monitorAlertService.CreateMonitorAlertRule>>,
    Error,
    monitor_alertservicev1_CreateMonitorAlertRuleRequest
  >,
) {
  return useMutation({
    mutationFn: (req: monitor_alertservicev1_CreateMonitorAlertRuleRequest) =>
      apiClient.monitorAlertService.CreateMonitorAlertRule(req),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['listMonitorAlertRules'] }),
    ...options,
  });
}

export function useUpdateMonitorAlertRule(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.monitorAlertService.UpdateMonitorAlertRule>>,
    Error,
    monitor_alertservicev1_UpdateMonitorAlertRuleRequest
  >,
) {
  return useMutation({
    mutationFn: (req: monitor_alertservicev1_UpdateMonitorAlertRuleRequest) =>
      apiClient.monitorAlertService.UpdateMonitorAlertRule(req),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['listMonitorAlertRules'] }),
    ...options,
  });
}

export function useDeleteMonitorAlertRule(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.monitorAlertService.DeleteMonitorAlertRule>>,
    Error,
    { id: number }
  >,
) {
  return useMutation({
    mutationFn: ({ id }: { id: number }) =>
      apiClient.monitorAlertService.DeleteMonitorAlertRule({ id }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['listMonitorAlertRules'] }),
    ...options,
  });
}

/** 手动触发一轮评估（与周期扫描同一内核），返回各启用规则的求值结论。 */
export function useEvaluateMonitorAlerts(
  options?: UseMutationOptions<
    monitor_alertservicev1_EvaluateMonitorAlertsResponse,
    Error,
    void
  >,
) {
  return useMutation({
    mutationFn: () => apiClient.monitorAlertService.EvaluateMonitorAlerts({}),
    ...options,
  });
}
