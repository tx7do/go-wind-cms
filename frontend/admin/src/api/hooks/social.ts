import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type socialservicev1_CreateSocialAccountRequest,
  type socialservicev1_CreateSocialDistributionRequest,
  type socialservicev1_DeleteSocialAccountRequest,
  type socialservicev1_DeleteSocialDistributionRequest,
  type socialservicev1_ListSocialAccountsResponse,
  type socialservicev1_ListSocialDistributionsResponse,
  type socialservicev1_ListSocialDistributionStatsHistoryRequest,
  type socialservicev1_MoveSocialDistributionRequest,
  type socialservicev1_ListSocialInboxRequest,
  type socialservicev1_ListSocialInboxResponse,
  type socialservicev1_MarkSocialInboxReadRequest,
  type socialservicev1_SocialDistributionSummary,
  type socialservicev1_ListSocialDistributionStatsHistoryResponse,
  type socialservicev1_PublishSocialDistributionRequest,
  type socialservicev1_RefreshSocialDistributionStatsRequest,
  type socialservicev1_SocialAccount,
  type socialservicev1_SocialDistribution,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 社交分发：账号 + 内容分发
// ==============================

// ── 社交账号 ──

export function useListSocialAccounts(
  query: PaginationQuery,
  options?: UseQueryOptions<socialservicev1_ListSocialAccountsResponse, Error>,
) {
  return useQuery({
    queryKey: ['listSocialAccounts', query],
    queryFn: () => apiClient.socialService.ListSocialAccounts(query.toRawParams()),
    ...options,
  });
}

export async function fetchListSocialAccounts(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listSocialAccounts', params],
    queryFn: () => apiClient.socialService.ListSocialAccounts(params.toRawParams()),
    retry: 0,
  });
}

export function useCreateSocialAccount(
  options?: UseMutationOptions<
    socialservicev1_SocialAccount,
    Error,
    socialservicev1_CreateSocialAccountRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.socialService.CreateSocialAccount(data),
    ...options,
  });
}

export function useUpdateSocialAccount(
  options?: UseMutationOptions<
    socialservicev1_SocialAccount,
    Error,
    { id: number; values: Record<string, unknown>; maskKeys?: string[] }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values, maskKeys }) =>
      apiClient.socialService.UpdateSocialAccount({
        id,
        data: { ...values } as unknown as socialservicev1_SocialAccount,
        updateMask: makeUpdateMask(maskKeys ?? Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteSocialAccount(
  options?: UseMutationOptions<unknown, Error, socialservicev1_DeleteSocialAccountRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.socialService.DeleteSocialAccount(data),
    ...options,
  });
}

// ── 内容分发 ──

export function useListSocialDistributions(
  query: PaginationQuery,
  options?: UseQueryOptions<socialservicev1_ListSocialDistributionsResponse, Error>,
) {
  return useQuery({
    queryKey: ['listSocialDistributions', query],
    queryFn: () => apiClient.socialService.ListSocialDistributions(query.toRawParams()),
    ...options,
  });
}

export async function fetchListSocialDistributions(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listSocialDistributions', params],
    queryFn: () => apiClient.socialService.ListSocialDistributions(params.toRawParams()),
    retry: 0,
  });
}

export function useCreateSocialDistribution(
  options?: UseMutationOptions<
    socialservicev1_SocialDistribution,
    Error,
    socialservicev1_CreateSocialDistributionRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.socialService.CreateSocialDistribution(data),
    ...options,
  });
}

export function useUpdateSocialDistribution(
  options?: UseMutationOptions<
    socialservicev1_SocialDistribution,
    Error,
    { id: number; values: Record<string, unknown>; maskKeys?: string[] }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values, maskKeys }) =>
      apiClient.socialService.UpdateSocialDistribution({
        id,
        data: { ...values } as unknown as socialservicev1_SocialDistribution,
        updateMask: makeUpdateMask(maskKeys ?? Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteSocialDistribution(
  options?: UseMutationOptions<unknown, Error, socialservicev1_DeleteSocialDistributionRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.socialService.DeleteSocialDistribution(data),
    ...options,
  });
}

export function useMoveSocialDistribution(
  options?: UseMutationOptions<
    socialservicev1_SocialDistribution,
    Error,
    socialservicev1_MoveSocialDistributionRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.socialService.MoveSocialDistribution(data),
    ...options,
  });
}

export function useListSocialInbox(
  req: socialservicev1_ListSocialInboxRequest,
  options?: UseQueryOptions<socialservicev1_ListSocialInboxResponse, Error>,
) {
  return useQuery({
    queryKey: ['listSocialInbox', req],
    queryFn: () => apiClient.socialService.ListSocialInbox(req),
    ...options,
  });
}

export function useMarkSocialInboxRead(
  options?: UseMutationOptions<unknown, Error, socialservicev1_MarkSocialInboxReadRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.socialService.MarkSocialInboxRead(data),
    ...options,
  });
}

export function usePublishSocialDistribution(
  options?: UseMutationOptions<
    socialservicev1_SocialDistribution,
    Error,
    socialservicev1_PublishSocialDistributionRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.socialService.PublishSocialDistribution(data),
    ...options,
  });
}

export function useSocialDistributionSummary(
  options?: Omit<
    UseQueryOptions<socialservicev1_SocialDistributionSummary, Error>,
    'queryKey' | 'queryFn'
  >,
) {
  return useQuery({
    queryKey: ['socialDistributionSummary'],
    queryFn: () => apiClient.socialService.GetSocialDistributionSummary({}),
    ...options,
  });
}

export function useSocialDistributionStatsHistory(
  req: socialservicev1_ListSocialDistributionStatsHistoryRequest,
  options?: Omit<
    UseQueryOptions<socialservicev1_ListSocialDistributionStatsHistoryResponse, Error>,
    'queryKey' | 'queryFn'
  >,
) {
  return useQuery({
    queryKey: ['socialDistributionStatsHistory', req],
    queryFn: () => apiClient.socialService.ListSocialDistributionStatsHistory(req),
    enabled: !!req.id,
    ...options,
  });
}

export function useRefreshSocialDistributionStats(
  options?: UseMutationOptions<
    socialservicev1_SocialDistribution,
    Error,
    socialservicev1_RefreshSocialDistributionStatsRequest
  >,
) {
  return useMutation({
    mutationFn: (data) => apiClient.socialService.RefreshSocialDistributionStats(data),
    ...options,
  });
}
