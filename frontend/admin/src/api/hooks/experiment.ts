import {
  useMutation,
  useQuery,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_CreateExperimentRequest,
  type contentservicev1_DeleteExperimentRequest,
  type contentservicev1_Experiment,
  type contentservicev1_ExperimentVariantStatsResponse,
  type contentservicev1_GetExperimentVariantStatsRequest,
  type contentservicev1_SetExperimentStatusRequest,
  type contentservicev1_UpdateExperimentRequest,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';
import { type PaginationQuery } from '@/core/transport/rest';

// ==============================
// 内容实验（A/B 分流 + 受众定向）
// ==============================

export function useExperimentsList() {
  return useMutation({
    mutationFn: (query: PaginationQuery) =>
      apiClient.experimentService.List(query.toRawParams()),
  });
}

export function useCreateExperiment(
  options?: UseMutationOptions<contentservicev1_Experiment, Error, contentservicev1_CreateExperimentRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.experimentService.Create(req),
    ...options,
  });
}

export function useUpdateExperiment(
  options?: UseMutationOptions<contentservicev1_Experiment, Error, contentservicev1_UpdateExperimentRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.experimentService.Update(req),
    ...options,
  });
}

export function useDeleteExperiment(
  options?: UseMutationOptions<Record<string, never>, Error, contentservicev1_DeleteExperimentRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.experimentService.Delete(req),
    ...options,
  });
}

export function useSetExperimentStatus(
  options?: UseMutationOptions<contentservicev1_Experiment, Error, contentservicev1_SetExperimentStatusRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.experimentService.SetStatus(req),
    ...options,
  });
}

/** 变体参与统计（人数是 uint64 JSON，值为字符串，展示前须 Number()） */
export function useExperimentVariantStats(
  req: contentservicev1_GetExperimentVariantStatsRequest | null,
) {
  return useQuery({
    queryKey: ['experiment-variant-stats', req?.experimentId],
    enabled: Boolean(req?.experimentId),
    queryFn: async () => {
      if (!req?.experimentId) {
        throw new Error('experimentId is required');
      }
      return apiClient.experimentService.GetVariantStats(req);
    },
    select: (resp: contentservicev1_ExperimentVariantStatsResponse) =>
      (resp.items ?? []).map((item) => ({
        variantName: item.variantName ?? '',
        participants: Number(item.participants ?? 0),
        conversions: Number(item.conversions ?? 0),
      })),
  });
}
