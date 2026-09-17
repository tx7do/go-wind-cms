/**
 * 内部消息模块枚举映射常量
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 消息状态 ==========

export const STATUS_COLORS: Record<string, string> = {
  DRAFT: 'default',
  PUBLISHED: 'success',
  SCHEDULED: 'processing',
  REVOKED: 'warning',
  ARCHIVED: 'default',
  DELETED: 'error',
};

export function getStatusMap(t: TFn) {
  return {
    DRAFT: { text: t('statusMap.DRAFT'), color: STATUS_COLORS.DRAFT },
    PUBLISHED: { text: t('statusMap.PUBLISHED'), color: STATUS_COLORS.PUBLISHED },
    SCHEDULED: { text: t('statusMap.SCHEDULED'), color: STATUS_COLORS.SCHEDULED },
    REVOKED: { text: t('statusMap.REVOKED'), color: STATUS_COLORS.REVOKED },
    ARCHIVED: { text: t('statusMap.ARCHIVED'), color: STATUS_COLORS.ARCHIVED },
    DELETED: { text: t('statusMap.DELETED'), color: STATUS_COLORS.DELETED },
  };
}

export function getStatusOptions(t: TFn) {
  return [
    { label: t('statusMap.DRAFT'), value: 'DRAFT' },
    { label: t('statusMap.PUBLISHED'), value: 'PUBLISHED' },
    { label: t('statusMap.SCHEDULED'), value: 'SCHEDULED' },
    { label: t('statusMap.REVOKED'), value: 'REVOKED' },
    { label: t('statusMap.ARCHIVED'), value: 'ARCHIVED' },
    { label: t('statusMap.DELETED'), value: 'DELETED' },
  ];
}

// ========== 消息类型 ==========

export const TYPE_COLORS: Record<string, string> = {
  NOTIFICATION: 'blue',
  PRIVATE: 'green',
  GROUP: 'orange',
};

export function getTypeMap(t: TFn) {
  return {
    NOTIFICATION: { text: t('typeMap.NOTIFICATION'), color: TYPE_COLORS.NOTIFICATION },
    PRIVATE: { text: t('typeMap.PRIVATE'), color: TYPE_COLORS.PRIVATE },
    GROUP: { text: t('typeMap.GROUP'), color: TYPE_COLORS.GROUP },
  };
}

export function getTypeOptions(t: TFn) {
  return [
    { label: t('typeMap.NOTIFICATION'), value: 'NOTIFICATION' },
    { label: t('typeMap.PRIVATE'), value: 'PRIVATE' },
    { label: t('typeMap.GROUP'), value: 'GROUP' },
  ];
}
