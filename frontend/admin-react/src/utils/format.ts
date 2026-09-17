/**
 * 字节数格式化为可读字符串（1024 进制）
 * @param bytes 字节数
 * @param decimals 保留小数位
 */
export function formatBytes(bytes: number, decimals: number = 2): string {
  // 处理0的特殊情况
  if (bytes === 0) return '0 B';

  // 定义单位换算的基数（1024进制）和单位列表
  const k = 1024;
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];

  // 计算最合适的单位索引
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  // 处理超出最大单位（PB）的情况
  const unitIndex = Math.min(i, units.length - 1);

  // 计算对应单位的数值并保留指定小数位
  const value = (bytes / k ** unitIndex).toFixed(decimals);

  // 移除末尾多余的0
  return `${parseFloat(value)} ${units[unitIndex]}`;
}
