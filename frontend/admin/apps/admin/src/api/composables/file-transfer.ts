import { useMutation, type UseMutationOptions } from '@tanstack/vue-query';
import { message } from 'ant-design-vue';

import { $t } from '@vben/locales';

import { apiClient } from '#/api/client';
import { RequestClient } from '#/transport/rest';

/**
 * 上传失败时给出界面提示；413（请求体超限，常见为 nginx client_max_body_size）
 * 单独提示文件过大，其余提示上传失败。错误对象可能是 response.data（字符串/对象）或原生 Error。
 */
export function notifyUploadError(error: unknown) {
  let text = '';
  if (typeof error === 'string') {
    text = error;
  } else if (error instanceof Error) {
    text = error.message;
  } else {
    try {
      text = JSON.stringify(error ?? {});
    } catch {
      text = '';
    }
  }
  if (/413|Request Entity Too Large/i.test(text)) {
    message.error($t('ui.notification.file_too_large'));
  } else {
    message.error($t('ui.notification.upload_failed'));
  }
}

const UPLOAD_PROGRESS_TOAST_KEY = 'upload-progress-toast';

/**
 * 显示/更新上传进度提示（antd message 按 key 复用同一实例更新内容）。
 * percent < 0 表示进度未知（如分块上传拿不到 total），只显示上传中。
 */
export function showUploadProgress(percent: number) {
  const text =
    $t('ui.notification.uploading') +
    (percent >= 0 ? ` ${Math.min(99, percent)}%` : '');
  message.loading({
    key: UPLOAD_PROGRESS_TOAST_KEY,
    content: text,
    duration: 0,
  });
}

/**
 * 关闭上传进度提示；上传结束（无论成败）都必须调用，避免残留常驻 toast。
 */
export function hideUploadProgress() {
  message.destroy(UPLOAD_PROGRESS_TOAST_KEY);
}

// ==============================
// 上传前图片压缩（仅照片类源，画图/截图类 PNG 原样直传）
// ==============================

/** 参与压缩的源类型：相机照片常见封装。PNG/SVG/GIF 是图表/截图主载体，保持原样。 */
const COMPRESSIBLE_IMAGE_TYPES = new Set(['image/jpeg', 'image/webp']);
/** 长边上限，超过才等比缩小；不足则只重编码不缩放 */
const IMAGE_MAX_LONG_EDGE = 1920;
/** 低于该体积分辨率再高也不值得处理，直传 */
const IMAGE_COMPRESS_MIN_BYTES = 300 * 1024;
/** WebP 重编码质量因子 */
const IMAGE_REENCODE_QUALITY = 0.82;

async function reencodeImage(file: File): Promise<File | null> {
  try {
    const bitmap = await createImageBitmap(file, {
      imageOrientation: 'from-image',
    });
    const longEdge = Math.max(bitmap.width, bitmap.height);
    const scale =
      longEdge > IMAGE_MAX_LONG_EDGE
        ? IMAGE_MAX_LONG_EDGE / longEdge
        : 1;
    const width = Math.round(bitmap.width * scale);
    const height = Math.round(bitmap.height * scale);
    if (width < 1 || height < 1) {
      bitmap.close();
      return null;
    }
    const canvas = document.createElement('canvas');
    canvas.width = width;
    canvas.height = height;
    const ctx = canvas.getContext('2d');
    if (!ctx) {
      bitmap.close();
      return null;
    }
    ctx.drawImage(bitmap, 0, 0, width, height);
    bitmap.close();
    const blob = await new Promise<Blob | null>((resolve) => {
      canvas.toBlob(resolve, 'image/webp', IMAGE_REENCODE_QUALITY);
    });
    if (!blob || blob.size === 0 || blob.size >= file.size) {
      return null;
    }
    const name = `${file.name.replace(/\.[^.]+$/, '')}.webp`;
    return new File([blob], name, { type: 'image/webp' });
  } catch {
    return null;
  }
}

/**
 * 上传前的图片预处理：照片类（jpeg/webp）等比缩小到长边上限内并重编码为 WebP
 * （同时剥离 EXIF 等元数据）；任何一步失败或无收益都原样返回，绝不阻断上传。
 */
export async function compressImageFile(file: File): Promise<File> {
  if (!COMPRESSIBLE_IMAGE_TYPES.has(file.type)) {
    return file;
  }
  if (file.size < IMAGE_COMPRESS_MIN_BYTES) {
    return file;
  }
  const out = await reencodeImage(file);
  return out ?? file;
}

/**
 * 从MinIO下载文件
 * @param fileId 文件元数据 ID（由列表接口返回，后端据此做归属校验）
 * @param preferPresignedUrl 是否优先使用预签名URL下载
 */
export async function downloadFile(
  fileId: number,
  preferPresignedUrl: boolean,
) {
  if (preferPresignedUrl) {
    const resp = await apiClient.fileTransferService.DownloadFile({
      fileId,
      preferPresignedUrl,
    });

    const url = (resp as any).downloadUrl || '';
    if (!url) return;

    // 不打印预签名 URL（含签名凭证），避免泄露到控制台

    const a = document.createElement('a');
    a.href = url;
    a.target = '_blank';
    a.download = 'download';
    document.body.append(a);
    a.click();
    a.remove();
    return;
  }

  const resp = await apiClient.fileTransferService.DownloadFile({
    fileId,
    preferPresignedUrl,
  });

  const contentType = (resp as any).contentType || 'application/octet-stream';
  const payload: ArrayBuffer | Blob | string | Uint8Array | undefined =
    (resp as any).file ?? (resp as any).data ?? (resp as any).payload ?? resp;

  function normalizeBase64(s: string): string {
    let str = s.replaceAll(/\s+/g, '');
    str = str.replaceAll('-', '+').replaceAll('_', '/');
    while (str.length % 4 !== 0) str += '=';
    return str;
  }

  function toBlob(data: any, type = contentType): Blob {
    if (!data) return new Blob([], { type });
    if (data instanceof Blob) return data;
    if (data instanceof ArrayBuffer) return new Blob([data], { type });
    if (ArrayBuffer.isView(data))
      return new Blob([data as BufferSource], { type });

    if (typeof data === 'string') {
      // 支持 data URI 或纯 base64（处理 URL-safe base64）
      const maybeBase64 = data.includes('base64,')
        ? data.split('base64,')[1]
        : data;
      const base64 = normalizeBase64(maybeBase64 ?? '');

      let binary: string;
      try {
        binary = atob(base64);
      } catch {
        // 如果仍然失败，返回空 Blob（也可以改为抛错或走异步 fetch fallback）
        return new Blob([], { type });
      }

      const len = binary.length;
      const arr = new Uint8Array(len);
      for (let i = 0; i < len; i++) {
        arr[i] = (binary.codePointAt(i) ?? 0) & 0xff;
      }
      return new Blob([arr], { type });
    }

    // fallback
    return new Blob([data], { type });
  }

  const blob = toBlob(payload, contentType);
  const objectUrl = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = objectUrl;
  a.download = 'download';
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
export async function uploadFile(
  bucketName: string,
  fileDirectory: string,
  fileData: File,
  method: 'post' | 'put' = 'post',
  onUploadProgress?: (progressEvent: any) => void,
) {
  const storageObject = JSON.stringify({
    bucketName,
    fileDirectory,
  });

  await RequestClient.getInstance().upload(
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

/**
 * 上传媒体资源文件
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

// -----------------------------------------------------------------------------
// 下载文件 Hook
// -----------------------------------------------------------------------------
export function useDownloadFile(
  options?: UseMutationOptions<
    void,
    Error,
    {
      fileId: number;
      preferPresignedUrl?: boolean;
    }
  >,
) {
  return useMutation({
    mutationFn: async ({
      fileId,
      preferPresignedUrl = false,
    }) => {
      return downloadFile(fileId, preferPresignedUrl);
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
      file: File;
      fileDirectory: string;
      method?: 'post' | 'put';
      onUploadProgress?: (progressEvent: any) => void;
    }
  >,
) {
  return useMutation({
    mutationFn: async ({
      bucketName,
      fileDirectory,
      file,
      method = 'post',
      onUploadProgress,
    }) => {
      return uploadFile(
        bucketName,
        fileDirectory,
        file,
        method,
        onUploadProgress,
      );
    },
    ...options,
  });
}
