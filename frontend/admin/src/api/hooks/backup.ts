// 备份档案（只读视图：历史列表 + 下载）
// ==============================

import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/api/client';

/** 备份档案条目（对象键/展示名/字节数/上传时间） */
export interface BackupArchive {
  objectKey: string;
  displayName: string;
  sizeBytes: number;
  lastModified?: string;
}

/** 备份历史列表 */
export function useListBackups() {
  return useQuery({
    queryKey: ['listBackups'],
    queryFn: async () => {
      const resp = await apiClient.backupService.ListBackups({});
      const items = (resp.items ?? []) as BackupArchive[];
      return { items, total: items.length };
    },
    staleTime: 30 * 1000,
  });
}

/** 归一化 base64（URL-safe → 标准，补齐 padding） */
function normalizeBase64(s: string): string {
  let str = s.replaceAll(/\s+/g, '');
  str = str.replaceAll('-', '+').replaceAll('_', '/');
  while (str.length % 4 !== 0) str += '=';
  return str;
}

/** base64 → Blob（file bytes 经 protojson 是 base64 串） */
function base64ToBlob(b64: string, type: string): Blob {
  let binary: string;
  try {
    binary = atob(normalizeBase64(b64));
  } catch (error) {
    // base64 非法时兜底空 Blob（产物损坏对用户不可见，必须留痕）
    console.error('base64 decode failed, fallback to empty blob:', error);
    return new Blob([], { type });
  }
  const arr = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    arr[i] = binary.codePointAt(i) ?? 0;
  }
  return new Blob([arr], { type });
}

/** 下载备份档案（bytes 内联，前端解 base64 存盘） */
export async function fetchDownloadBackup(archive: BackupArchive) {
  const resp = await apiClient.backupService.DownloadBackup({
    objectKey: archive.objectKey,
  });
  const contentType = resp.contentType || 'application/gzip';
  // protojson 下 bytes 是 base64 串；生成类型为 Uint8Array 时按字节数组兜底
  const raw: unknown = resp.file;
  let blob: Blob;
  if (typeof raw === 'string') {
    blob = base64ToBlob(raw, contentType);
  } else if (Array.isArray(raw)) {
    blob = new Blob([new Uint8Array(raw as number[]).buffer as ArrayBuffer], { type: contentType });
  } else {
    blob = new Blob([], { type: contentType });
  }
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = archive.displayName || 'backup.json.gz';
  document.body.append(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(url);
}
