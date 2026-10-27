import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
  type mediaservicev1_ListDuplicatedMediaAssetsResponse,
  type mediaservicev1_ListMediaAssetReferencesResponse,
  type mediaservicev1_ListMediaAssetResponse,
  type mediaservicev1_MediaAsset,
  type mediaservicev1_GetMediaAssetRequest,
  type mediaservicev1_CreateMediaAssetRequest,
  type mediaservicev1_DeleteMediaAssetRequest,
} from '@/api/generated/admin/service/v1';
import { makeUpdateMask, type PaginationQuery } from '@/core/transport/rest';
import { RequestClient, queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 媒体资源管理
// ==============================

export function useListMediaAssets(
  query: PaginationQuery,
  options?: UseQueryOptions<mediaservicev1_ListMediaAssetResponse, Error>,
) {
  return useQuery({
    queryKey: ['listMediaAssets', query],
    queryFn: () => apiClient.mediaAssetService.List(query.toRawParams()),
    ...options,
  });
}

export async function fetchListMediaAssets(params: PaginationQuery) {
  return queryClient.fetchQuery({
    queryKey: ['listMediaAssets', params],
    queryFn: () => apiClient.mediaAssetService.List(params.toRawParams()),
    retry: 0,
  });
}

/** 重复资源分组（同租户内同内容哈希 >1 条） */
export async function fetchListDuplicatedMediaAssets(limit = 50) {
  return queryClient.fetchQuery({
    queryKey: ['listDuplicatedMediaAssets', limit],
    queryFn: () =>
      apiClient.mediaAssetService.ListDuplicatedMediaAssets({ limit }),
    retry: 0,
  });
}

/** 重复资源分组（弹窗用 useQuery；limit<=0 时禁用） */
export function useListDuplicatedMediaAssets(
  limit: number,
  options?: UseQueryOptions<mediaservicev1_ListDuplicatedMediaAssetsResponse, Error>,
) {
  return useQuery({
    queryKey: ['listDuplicatedMediaAssets', limit],
    queryFn: () =>
      apiClient.mediaAssetService.ListDuplicatedMediaAssets({ limit }),
    enabled: limit > 0,
    ...options,
  });
}

/** 媒体被内容引用的关系（哪些文章/页面在用） */
export function useListMediaAssetReferences(
  mediaAssetId?: number,
  options?: UseQueryOptions<mediaservicev1_ListMediaAssetReferencesResponse, Error>,
) {
  return useQuery({
    queryKey: ['listMediaAssetReferences', mediaAssetId],
    queryFn: () =>
      apiClient.mediaAssetService.ListMediaAssetReferences({
        mediaAssetId: mediaAssetId ?? 0,
      }),
    enabled: !!mediaAssetId,
    ...options,
  });
}

export function useGetMediaAsset(
  req: mediaservicev1_GetMediaAssetRequest,
  options?: UseQueryOptions<mediaservicev1_MediaAsset, Error>,
) {
  return useQuery({
    queryKey: ['getMediaAsset', req],
    queryFn: () => apiClient.mediaAssetService.Get(req),
    ...options,
  });
}

export function useCreateMediaAsset(
  options?: UseMutationOptions<mediaservicev1_MediaAsset, Error, mediaservicev1_CreateMediaAssetRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.mediaAssetService.Create(data),
    ...options,
  });
}

export function useUpdateMediaAsset(
  options?: UseMutationOptions<
    mediaservicev1_MediaAsset,
    Error,
    { id: number; values: Record<string, unknown>; maskKeys?: string[] }
  >,
) {
  return useMutation({
    mutationFn: ({ id, values, maskKeys }) =>
      apiClient.mediaAssetService.Update({
        id,
        data: { ...values } as unknown as mediaservicev1_MediaAsset,
        updateMask: makeUpdateMask(maskKeys ?? Object.keys(values ?? {})),
      }),
    ...options,
  });
}

export function useDeleteMediaAsset(
  options?: UseMutationOptions<unknown, Error, mediaservicev1_DeleteMediaAssetRequest>,
) {
  return useMutation({
    mutationFn: (data) => apiClient.mediaAssetService.Delete(data),
    ...options,
  });
}

// ==============================
// 媒体资源上传
// ==============================

/**
 * 上传媒体资源文件（后端流式接收，presigned 全层禁用）
 */
export async function uploadMediaAsset(
  data: {
    altText?: string;
    caption?: string;
    fileDirectory?: string;
    title?: string;
  },
  fileData: File,
  onUploadProgress?: (progressEvent: { loaded?: number; total?: number }) => void,
) {
  return RequestClient.getInstance().upload(
    'admin/v1/file/asset/upload',
    {
      file: fileData,
      ...data,
      sourceFileName: fileData.name,
      mime: fileData.type,
      size: fileData.size,
    },
    { onUploadProgress },
  );
}
