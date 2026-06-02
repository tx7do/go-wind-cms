import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type storageservicev1_DownloadFileResponse,
} from '@/api/generated/app/service/v1';
import { downloadFile, uploadFile } from '@/api/service/file-transfer';
import { queryClient } from '@/core';
import type { AxiosProgressEvent } from 'axios';

// ==============================
// 下载文件 Hook
// ==============================
export function useDownloadFile(
  options?: UseMutationOptions<
    storageservicev1_DownloadFileResponse,
    Error,
    {
      bucketName: string;
      objectName: string;
      preferPresignedUrl?: boolean;
    }
  >,
) {
  return useMutation({
    mutationFn: async ({ bucketName, objectName, preferPresignedUrl = false }) => {
      return downloadFile(bucketName, objectName, preferPresignedUrl);
    },
    ...options,
  });
}

// ==============================================
// 下载文件 【给 Store / 外部调用】不带 Hook 的方法
// ==============================================
export async function fetchDownloadFile(params: {
  bucketName: string;
  objectName: string;
  preferPresignedUrl: boolean;
}) {
  return queryClient.fetchQuery({
    queryKey: ['downloadFile', params],
    queryFn: () => downloadFile(params.bucketName, params.objectName, params.preferPresignedUrl),
    retry: 0,
  });
}

// ==============================
// 上传文件 Hook（支持进度）
// ==============================
export function useUploadFile(
  options?: UseMutationOptions<
    void,
    Error,
    {
      bucketName: string;
      fileDirectory: string;
      file: File;
      method?: 'post' | 'put';
      onUploadProgress?: (progressEvent: AxiosProgressEvent) => void;
    }
  >,
) {
  return useMutation({
    mutationFn: async ({ bucketName, fileDirectory, file, method = 'post', onUploadProgress }) => {
      return uploadFile(bucketName, fileDirectory, file, method, onUploadProgress);
    },
    ...options,
  });
}
