import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type storageservicev1_DownloadFileResponse,
} from '@/api/generated/app/service/v1';
import { apiClient } from '@/api/client';
import { RequestClient, queryClient } from '@/core';
import type { AxiosProgressEvent } from 'axios';

// ==============================
// 文件传输服务 API
// ==============================

/**
 * 下载文件
 * @param bucketName 文件桶名称
 * @param objectName 对象名称
 * @param preferPresignedUrl 是否优先使用预签名URL下载
 */
export async function downloadFile(
  bucketName: string,
  objectName: string,
  preferPresignedUrl: boolean,
) {
  return apiClient.fileTransferService.DownloadFile({
    storageObject: {bucketName, objectName},
    preferPresignedUrl,
  });
}

/**
 * 上传文件到服务器
 * @param bucketName 文件桶名称
 * @param fileDirectory 远端存储文件目录
 * @param fileData 文件数据
 * @param method 上传方法，支持 'post' 和 'put'
 * @param onUploadProgress 上传进度回调函数
 */
export async function uploadFile(
  bucketName: string,
  fileDirectory: string,
  fileData: File,
  method: 'post' | 'put' = 'post',
  onUploadProgress?: (progressEvent: AxiosProgressEvent) => void,
) {
  const storageObject = JSON.stringify({
    bucketName,
    fileDirectory,
  });

  await RequestClient.getInstance().upload(
    'app/v1/file/upload',
    {
      file: fileData,
      storageObject,
      sourceFileName: fileData.name,
      mime: fileData.type,
      size: String(fileData.size),
      method,
    },
    {onUploadProgress},
  );
}

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
