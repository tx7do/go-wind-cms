import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/vue-query';
import { downloadFile, uploadFile } from '@/api/service/file-transfer';

// 直接导出 service 层函数
export { downloadFile, uploadFile };

// ==============================
// 下载文件（Mutation）
// ==============================
export function useDownloadFile(
  options?: UseMutationOptions<
    void,
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

// ==============================
// 上传文件（Mutation，支持进度）
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
      onUploadProgress?: (progress: any) => void;
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
