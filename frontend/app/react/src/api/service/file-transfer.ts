import {
  createFileTransferServiceClient,
} from '@/api/generated/app/service/v1';
import { RequestClient, requestApi } from '@/core';
import type { AxiosProgressEvent } from 'axios';

let _instance: ReturnType<typeof createFileTransferServiceClient> | null = null;

/**
 * 获取文件传输服务单例（延迟初始化）
 */
export function getFileTransferService() {
  if (!_instance) {
    _instance = createFileTransferServiceClient(requestApi);
  }
  return _instance;
}

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
  return getFileTransferService().DownloadFile({
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
