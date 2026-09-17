import {
  useMutation,
  type UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from '@tanstack/react-query';
import {
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
    { id: number; values: Record<string, any>; maskKeys?: string[] }
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
  options?: UseMutationOptions<{}, Error, mediaservicev1_DeleteMediaAssetRequest>,
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
  onUploadProgress?: (progressEvent: any) => void,
) {
  return await RequestClient.getInstance().upload(
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
