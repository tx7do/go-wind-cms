import { useMutation, type UseMutationOptions } from '@tanstack/react-query';
import { apiClient } from '@/api/client';
import { RequestClient } from '@/core';

// ==============================
// 文件下载 / 上传辅助函数
// ==============================

/**
 * 从MinIO下载文件
 * @param bucketName 文件桶名称
 * @param objectName 对象名称
 * @param preferPresignedUrl 是否优先使用预签名URL下载
 */
async function downloadFile(
  bucketName: string,
  objectName: string,
  preferPresignedUrl: boolean,
) {
  if (preferPresignedUrl) {
    const resp = await apiClient.fileTransferService.DownloadFile({
      storageObject: { bucketName, objectName },
      preferPresignedUrl,
    });

    // 生成类型未导出 downloadUrl 联合分支字段：结构化收窄替代 as any
    const url = (resp as { downloadUrl?: string }).downloadUrl || '';
    if (!url) return;

    console.log('Downloading file transfer...', url);

    const a = document.createElement('a');
    a.href = url;
    a.target = '_blank';
    a.download = objectName || 'download';
    document.body.append(a);
    a.click();
    a.remove();
    return;
  }

  const resp = await apiClient.fileTransferService.DownloadFile({
    storageObject: { bucketName, objectName },
    preferPresignedUrl,
  });

  // 同上：结构化收窄（base64 file / url 二选一分支）
  type DownloadPayload = { contentType?: string; file?: string; data?: string; payload?: string };
  const dl = resp as unknown as DownloadPayload;
  const contentType = dl.contentType || 'application/octet-stream';
  const payload: ArrayBuffer | Blob | string | Uint8Array | undefined =
    dl.file ?? dl.data ?? dl.payload;

  function normalizeBase64(s: string): string {
    let str = s.replaceAll(/\s+/g, '');
    str = str.replaceAll('-', '+').replaceAll('_', '/');
    while (str.length % 4 !== 0) str += '=';
    return str;
  }

  function toBlob(data: unknown, type = contentType): Blob {
    if (!data) return new Blob([], { type });
    if (data instanceof Blob) return data;
    if (data instanceof ArrayBuffer) return new Blob([data], { type });
    if (ArrayBuffer.isView(data)) return new Blob([data as BufferSource], { type });

    if (typeof data === 'string') {
      // 支持 data URI 或纯 base64（处理 URL-safe base64）
      const str = data as string;
    const maybeBase64 = str.includes('base64,') ? str.split('base64,')[1] : str;
      const base64 = normalizeBase64(maybeBase64 ?? '');

      let binary: string;
      try {
        binary = atob(base64);
      } catch (error) {
        // base64 非法时兜底返回空 Blob（下载产物损坏对用户不可见，必须留痕）
        console.error('toBlob: base64 解码失败，返回空 Blob:', error);
        return new Blob([], { type });
      }

      const len = binary.length;
      const arr = new Uint8Array(len);
      for (let i = 0; i < len; i++) {
        arr[i] = (binary.codePointAt(i) ?? 0) & 0xff;
      }
      return new Blob([arr], { type });
    }

    // fallback：残余形态按 JSON 文本收编
    return new Blob([JSON.stringify(data)], { type });
  }

  const blob = toBlob(payload, contentType);
  const objectUrl = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = objectUrl;
  a.download = objectName || 'download';
  document.body.append(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(objectUrl);
}

/**
 * 上传文件到MinIO
 * @param bucketName 文件桶名称
 * @param fileDirectory 远端存储文件目录
 * @param fileData 文件数据
 * @param method 上传方法，支持 'post' 和 'put'
 * @param onUploadProgress 上传进度回调函数
 */
// eslint-disable-next-line max-params -- 上传 API 固有参数面（bucket/目录/文件/方法/进度）
export async function uploadFile(
  bucketName: string,
  fileDirectory: string,
  fileData: File,
  method: 'post' | 'put' = 'post',
  uploadOptions: { onUploadProgress?: (progressEvent: { loaded?: number; total?: number }) => void } = {},
) {
  const { onUploadProgress } = uploadOptions;
  const storageObject = JSON.stringify({
    bucketName,
    fileDirectory,
  });

  return RequestClient.getInstance().upload(
    'admin/v1/file/upload',
    {
      file: fileData,
      storageObject,
      sourceFileName: fileData.name,
      mime: fileData.type,
      size: fileData.size,
      method,
    },
    { onUploadProgress },
  );
}

// -----------------------------------------------------------------------------
// 下载文件 Hook
// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
// 上传文件 Hook（支持进度）
// -----------------------------------------------------------------------------
export function useUploadFile(
  options?: UseMutationOptions<
    void,
    Error,
    {
      bucketName: string;
      fileDirectory: string;
      file: File;
      method?: 'post' | 'put';
      uploadOptions?: { onUploadProgress?: (progress: { loaded?: number; total?: number }) => void };
    }
  >,
) {
  return useMutation({
    mutationFn: async ({ bucketName, fileDirectory, file, method = 'post', uploadOptions }) => {
      await uploadFile(bucketName, fileDirectory, file, method, uploadOptions);
    },
    ...options,
  });
}
