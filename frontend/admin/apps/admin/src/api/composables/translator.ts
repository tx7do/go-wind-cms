import type { translatorservicev1_TranslateResponse } from '#/api/generated/admin/service/v1';

import { useMutation, type UseMutationOptions } from '@tanstack/vue-query';

import { apiClient } from '#/api/client';

// ==============================
// 翻译服务
// ==============================

export function useTranslate(
  options?: UseMutationOptions<
    translatorservicev1_TranslateResponse,
    Error,
    { content: string; sourceLanguage?: string; targetLanguage: string }
  >,
) {
  return useMutation({
    mutationFn: ({ targetLanguage, content, sourceLanguage = 'auto' }) =>
      apiClient.translatorService.Translate({
        sourceLanguage,
        targetLanguage,
        content,
      }),
    ...options,
  });
}
