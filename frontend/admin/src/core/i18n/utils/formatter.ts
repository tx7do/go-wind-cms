import type { SupportedLocale } from '@/locales';

/**
 * 日期格式化工具
 * @param date 日期对象/时间戳
 * @param locale 目标语言
 * @param options Intl.DateTimeFormat 配置
 */
export const formatDate = (
  date: Date | number | string,
  locale: SupportedLocale,
  options?: Intl.DateTimeFormatOptions,
): string => {
  const targetDate = typeof date === 'string' || typeof date === 'number' ? new Date(date) : date;

  // 默认格式化配置
  const defaultOptions: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  };

  return new Intl.DateTimeFormat(locale, { ...defaultOptions, ...options }).format(targetDate);
};

/**
 * 数字格式化工具（支持千分位、小数位）
 * @param num 数字
 * @param locale 目标语言
 * @param options Intl.NumberFormat 配置
 */
export const formatNumber = (
  num: number,
  locale: SupportedLocale,
  options?: Intl.NumberFormatOptions,
): string => {
  const defaultOptions: Intl.NumberFormatOptions = {
    maximumFractionDigits: 2,
  };

  return new Intl.NumberFormat(locale, { ...defaultOptions, ...options }).format(num);
};

/**
 * 货币格式化工具
 * @param amount 金额
 * @param locale 目标语言
 * @param currency 货币代码（CNY/USD等）
 * @param options Intl.NumberFormat 配置
 */
export const formatCurrency = (
  amount: number,
  locale: SupportedLocale,
  currency: string = 'CNY',
  options?: Intl.NumberFormatOptions,
): string => {
  const defaultOptions: Intl.NumberFormatOptions = {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  };

  return new Intl.NumberFormat(locale, { ...defaultOptions, ...options }).format(amount);
};

/**
 * 百分比格式化工具
 * @param value 小数（0.85 → 85%）
 * @param locale 目标语言
 * @param options Intl.NumberFormat 配置
 */
export const formatPercent = (
  value: number,
  locale: SupportedLocale,
  options?: Intl.NumberFormatOptions,
): string => {
  const defaultOptions: Intl.NumberFormatOptions = {
    style: 'percent',
    maximumFractionDigits: 2,
  };

  return new Intl.NumberFormat(locale, { ...defaultOptions, ...options }).format(value);
};

/**
 * 相对时间格式化（如“3分钟前”“昨天”）
 * @param date 目标日期
 * @param locale 目标语言
 * @param options Intl.RelativeTimeFormat 配置
 */
export const formatRelativeTime = (
  date: Date | number,
  locale: SupportedLocale,
  options?: Intl.RelativeTimeFormatOptions,
): string => {
  const targetDate = typeof date === 'number' ? new Date(date) : date;
  const now = Date.now();
  const diffMs = targetDate.getTime() - now;

  // 时间单位换算
  const SECOND = 1000;
  const MINUTE = 60 * SECOND;
  const HOUR = 60 * MINUTE;
  const DAY = 24 * HOUR;

  let value: number;
  let unit: Intl.RelativeTimeFormatUnit;

  const absDiff = Math.abs(diffMs);
  if (absDiff < MINUTE) {
    value = Math.round(diffMs / SECOND);
    unit = 'second';
  } else if (absDiff < HOUR) {
    value = Math.round(diffMs / MINUTE);
    unit = 'minute';
  } else if (absDiff < DAY) {
    value = Math.round(diffMs / HOUR);
    unit = 'hour';
  } else {
    value = Math.round(diffMs / DAY);
    unit = 'day';
  }

  return new Intl.RelativeTimeFormat(locale, { numeric: 'auto', ...options }).format(value, unit);
};
