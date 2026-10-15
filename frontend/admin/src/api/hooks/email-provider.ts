import { useMutation, type UseMutationOptions } from '@tanstack/react-query';
import type {
  notifyservicev1_CreateEmailProviderRequest,
  notifyservicev1_EmailProvider,
  notifyservicev1_TestEmailSendRequest,
  notifyservicev1_TestEmailSendResponse,
  notifyservicev1_UpdateEmailProviderRequest,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';

// ==============================
// 邮件出站 SMTP 配置（平台基础设施；凭据加密落库、接口永不回显）
// ==============================

/** 单配置卡形态：拉最新配置（邮件链取第一个启用项） */
export function useEmailProvidersList(
  options?: UseMutationOptions<{
    items?: notifyservicev1_EmailProvider[];
    total?: number;
  }, Error, void>,
) {
  return useMutation({
    mutationFn: () =>
      apiClient.emailProviderService.List({
        page: 1,
        pageSize: 20,
        offset: 0,
        limit: 20,
        sorting: [],
      }),
    ...options,
  });
}

export function useCreateEmailProvider(
  options?: UseMutationOptions<notifyservicev1_EmailProvider, Error, notifyservicev1_CreateEmailProviderRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.emailProviderService.Create(req),
    ...options,
  });
}

export function useUpdateEmailProvider(
  options?: UseMutationOptions<notifyservicev1_EmailProvider, Error, notifyservicev1_UpdateEmailProviderRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.emailProviderService.Update(req),
    ...options,
  });
}

export function useTestEmailProviderSend(
  options?: UseMutationOptions<notifyservicev1_TestEmailSendResponse, Error, notifyservicev1_TestEmailSendRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.emailProviderService.TestSend(req),
    ...options,
  });
}
