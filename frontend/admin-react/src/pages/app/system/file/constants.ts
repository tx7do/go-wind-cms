/**
 * 文件模块枚举映射常量
 * 页面和 Drawer 共用
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 存储提供商 ==========

/** 提供商颜色映射 */
export const PROVIDER_COLORS: Record<string, string> = {
  MINIO: 'blue',
  ALIYUN: 'orange',
  TENCENT: 'green',
  AWS: 'cyan',
  AZURE: 'purple',
  BAIDU: 'default',
  QINIU: 'default',
  GOOGLE: 'default',
  HUAWEI: 'default',
  LOCAL: 'default',
};

/** 获取提供商映射（text + color） */
export function getProviderMap(t: TFn) {
  return {
    MINIO: { text: t('providerMap.MINIO'), color: PROVIDER_COLORS.MINIO },
    ALIYUN: { text: t('providerMap.ALIYUN'), color: PROVIDER_COLORS.ALIYUN },
    TENCENT: { text: t('providerMap.TENCENT'), color: PROVIDER_COLORS.TENCENT },
    AWS: { text: t('providerMap.AWS'), color: PROVIDER_COLORS.AWS },
    AZURE: { text: t('providerMap.AZURE'), color: PROVIDER_COLORS.AZURE },
    BAIDU: { text: t('providerMap.BAIDU'), color: PROVIDER_COLORS.BAIDU },
    QINIU: { text: t('providerMap.QINIU'), color: PROVIDER_COLORS.QINIU },
    GOOGLE: { text: t('providerMap.GOOGLE'), color: PROVIDER_COLORS.GOOGLE },
    HUAWEI: { text: t('providerMap.HUAWEI'), color: PROVIDER_COLORS.HUAWEI },
    LOCAL: { text: t('providerMap.LOCAL'), color: PROVIDER_COLORS.LOCAL },
  };
}

// ========== 文件大小格式化 ==========

/** 将字节数格式化为可读字符串；兼容 protojson 把 uint64 序列化成字符串的运行时形态 */
export function formatFileSize(bytes: number | string | undefined | null): string {
  const n = Number(bytes);
  if (!Number.isFinite(n) || n <= 0) return '-';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let size = n;
  let unitIndex = 0;
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }
  return `${size.toFixed(unitIndex === 0 ? 0 : 2)} ${units[unitIndex]}`;
}
