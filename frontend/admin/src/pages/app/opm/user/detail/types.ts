/**
 * 用户详情页 - 标签页常量
 */
export const TabEnum = {
  BASIC_INFO: 'basicInfo',
  API_AUDIT_LOG: 'apiAuditLog',
  INTERNAL_MESSAGE: 'internalMessage',
} as const;

export type TabEnum = (typeof TabEnum)[keyof typeof TabEnum];
