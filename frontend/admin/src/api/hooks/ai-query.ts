import { useMutation, type UseMutationOptions } from '@tanstack/react-query';
import { apiClient } from '@/api/client';
import type {
  aiservicev1_AiQueryRow,
  aiservicev1_AskAiQueryResponse,
} from '@/api/generated/admin/service/v1';

export interface AskAiQueryHistoryItem {
  question: string;
  sql: string;
  resultSummary: string;
}

export interface AskAiQueryParams {
  question: string;
  lang?: string;
  withAnswer?: boolean;
  history?: AskAiQueryHistoryItem[];
}

// useAskAiQuery — 智能问数：NL → 只读 SQL → 结构化结果（+可选自然语言结论）。
// 平台管理员专属（后端对租户返回 403）。每次调用计入 AI_TOKENS 用量。
export function useAskAiQuery(
  options?: UseMutationOptions<aiservicev1_AskAiQueryResponse, Error, AskAiQueryParams>,
) {
  return useMutation({
    mutationFn: (req: AskAiQueryParams) =>
      apiClient.aiQueryService.Ask({
        question: req.question,
        lang: req.lang ?? 'zh-CN',
        withAnswer: req.withAnswer ?? true,
        history: req.history,
      }),
    ...options,
  });
}

export type AiQueryRow = aiservicev1_AiQueryRow;
