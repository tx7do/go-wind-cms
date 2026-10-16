import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type flowservicev1_CreateFlowRuleRequest,
  type flowservicev1_ListFlowRuleResponse,
  type flowservicev1_UpdateFlowRuleRequest,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';

// ==============================
// 自动化 Flow 规则管理
// ==============================

/** 内容事件族（与 webhook 事件一致） */
export const FLOW_EVENT_OPTIONS = [
  'post.created',
  'post.updated',
  'post.published',
  'post.deleted',
  'page.created',
  'page.updated',
  'page.published',
  'page.deleted',
].map((v) => ({ value: v, label: v }));

/** 实体状态条件候选 */
export const FLOW_STATUS_CONDITION_OPTIONS = [
  { value: '', label: '-' },
  { value: 'POST_STATUS_PUBLISHED', label: 'POST_STATUS_PUBLISHED' },
  { value: 'POST_STATUS_DRAFT', label: 'POST_STATUS_DRAFT' },
  { value: 'POST_STATUS_SCHEDULED', label: 'POST_STATUS_SCHEDULED' },
  { value: 'POST_STATUS_ARCHIVED', label: 'POST_STATUS_ARCHIVED' },
  { value: 'PAGE_STATUS_PUBLISHED', label: 'PAGE_STATUS_PUBLISHED' },
  { value: 'PAGE_STATUS_DRAFT', label: 'PAGE_STATUS_DRAFT' },
];

export function useListFlowRules(
  options?: UseQueryOptions<flowservicev1_ListFlowRuleResponse, Error>,
) {
  return useQuery({
    queryKey: ['listFlowRules'],
    // AGENTS 铁律：List 禁传空 sorting[]，排序走 orderBy
    queryFn: () =>
      apiClient.flowRuleService.List({
        orderBy: JSON.stringify([{ field: 'created_at', direction: 'DESC' }]),
      } as never),
    ...options,
  });
}

export function useCreateFlowRule(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.flowRuleService.Create>>,
    Error,
    flowservicev1_CreateFlowRuleRequest
  >,
) {
  return useMutation({
    mutationFn: (req) => apiClient.flowRuleService.Create(req),
    ...options,
  });
}

export function useUpdateFlowRule(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.flowRuleService.Update>>,
    Error,
    flowservicev1_UpdateFlowRuleRequest
  >,
) {
  return useMutation({
    mutationFn: (req) => apiClient.flowRuleService.Update(req),
    ...options,
  });
}

export function useDeleteFlowRule(
  options?: UseMutationOptions<
    Awaited<ReturnType<typeof apiClient.flowRuleService.Delete>>,
    Error,
    { id: number }
  >,
) {
  return useMutation({
    mutationFn: (req) => apiClient.flowRuleService.Delete(req),
    ...options,
  });
}
