import { useMutation, useQuery, useQueryClient, type UseMutationOptions, type UseQueryOptions } from '@tanstack/react-query';
import {
  type aiservicev1_AiAgentApp as AiAgentApp,
  type aiservicev1_AiAgentRun,
  type aiservicev1_CreateAiAgentAppRequest,
  type aiservicev1_DeleteAiAgentAppRequest,
  type aiservicev1_GetAiAgentRunRequest,
  type aiservicev1_ListAiAgentAppsResponse,
  type aiservicev1_ListAiAgentRunsResponse,
  type aiservicev1_UpdateAiAgentAppRequest,
} from '@/api/generated/admin/service/v1';
import { PaginationQuery, queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// AI Agent 编排（可视化画布）
// 本切片（①）：应用 CRUD 与运行台账读取；试运行 Run 在执行引擎切片接入。
// ==============================

const LIST_KEY = 'listAiAgentApps';
const RUNS_KEY = 'listAiAgentRuns';

export function useListAiAgentApps(
  query: PaginationQuery,
  options?: UseQueryOptions<aiservicev1_ListAiAgentAppsResponse, Error>,
) {
  return useQuery({
    queryKey: [LIST_KEY, query],
    queryFn: () => apiClient.aiAgentService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListAiAgentApps(query: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: [LIST_KEY, params(query)],
    queryFn: () => apiClient.aiAgentService.List(params(query)),
    retry: 0,
  });
}

function params(query: PaginationQuery) {
  return query.toRawParams();
}

export function useCreateAiAgentApp(
  options?: UseMutationOptions<AiAgentApp, Error, aiservicev1_CreateAiAgentAppRequest>,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.aiAgentService.Create(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useUpdateAiAgentApp(
  options?: UseMutationOptions<{}, Error, aiservicev1_UpdateAiAgentAppRequest>,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.aiAgentService.Update(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useDeleteAiAgentApp(
  options?: UseMutationOptions<{}, Error, aiservicev1_DeleteAiAgentAppRequest>,
) {
  const queryClientRef = useQueryClient();
  return useMutation({
    mutationFn: (req) => apiClient.aiAgentService.Delete(req),
    onSuccess: () => queryClientRef.invalidateQueries({ queryKey: [LIST_KEY] }),
    ...options,
  });
}

export function useListAiAgentRuns(
  appId: number | undefined,
  query: PaginationQuery,
  options?: UseQueryOptions<aiservicev1_ListAiAgentRunsResponse, Error>,
) {
  return useQuery({
    queryKey: [RUNS_KEY, appId, query],
    // appId 走 formValues → 通用 eq 过滤（camel→snake 归一由后端机制处理）
    queryFn: () =>
      apiClient.aiAgentService.ListRuns(
        new PaginationQuery({
          paging: { page: query.paging?.page, pageSize: query.paging?.pageSize },
          formValues: appId ? { appId: String(appId) } : undefined,
        }).toRawParams(),
      ),
    enabled: !!appId,
    ...options,
  });
}

export function useGetAiAgentRun(
  options?: UseMutationOptions<aiservicev1_AiAgentRun, Error, aiservicev1_GetAiAgentRunRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.aiAgentService.GetRun(req),
    ...options,
  });
}

export type { AiAgentApp };
