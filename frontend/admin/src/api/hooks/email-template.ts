import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import type {
  contentservicev1_CreateEmailTemplateRequest,
  contentservicev1_DeleteEmailTemplateRequest,
  contentservicev1_EmailTemplate,
  contentservicev1_UpdateEmailTemplateRequest,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';
import { type PaginationQuery } from '@/core';

// ==============================
// 邮件模板（系统通知出站的可视化模板）
// ==============================

export function useEmailTemplatesList() {
  return useMutation({
    mutationFn: (query: PaginationQuery) =>
      apiClient.emailTemplateService.List(query.toRawParams()),
  });
}

export function useCreateEmailTemplate(
  options?: UseMutationOptions<contentservicev1_EmailTemplate, Error, contentservicev1_CreateEmailTemplateRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.emailTemplateService.Create(req),
    ...options,
  });
}

export function useUpdateEmailTemplate(
  options?: UseMutationOptions<contentservicev1_EmailTemplate, Error, contentservicev1_UpdateEmailTemplateRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.emailTemplateService.Update(req),
    ...options,
  });
}

export function useDeleteEmailTemplate(
  options?: UseMutationOptions<Record<string, never>, Error, contentservicev1_DeleteEmailTemplateRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.emailTemplateService.Delete(req),
    ...options,
  });
}
