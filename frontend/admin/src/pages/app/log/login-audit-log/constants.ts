/**
 * 登录审计日志模块枚举映射常量
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 状态 ==========

/** 状态颜色映射 */
export const STATUS_COLORS: Record<string, string> = {
  SUCCESS: 'success',
  FAILED: 'error',
  PARTIAL: 'warning',
  LOCKED: 'default',
};

/** 获取状态映射（text + color） */
export function getStatusMap(t: TFn) {
  return {
    SUCCESS: { text: t('status.SUCCESS'), color: STATUS_COLORS.SUCCESS },
    FAILED: { text: t('status.FAILED'), color: STATUS_COLORS.FAILED },
    PARTIAL: { text: t('status.PARTIAL'), color: STATUS_COLORS.PARTIAL },
    LOCKED: { text: t('status.LOCKED'), color: STATUS_COLORS.LOCKED },
  };
}

/** 获取状态选项列表 */
export function getStatusOptions(t: TFn) {
  return [
    { label: t('status.SUCCESS'), value: 'SUCCESS' },
    { label: t('status.FAILED'), value: 'FAILED' },
    { label: t('status.PARTIAL'), value: 'PARTIAL' },
    { label: t('status.LOCKED'), value: 'LOCKED' },
  ];
}

// ========== 操作类型 ==========

/** 操作类型颜色映射 */
export const ACTION_TYPE_COLORS: Record<string, string> = {
  LOGIN: 'success',
  LOGOUT: 'default',
  SESSION_EXPIRED: 'warning',
  KICKED_OUT: 'error',
  PASSWORD_RESET: 'processing',
};

/** 获取操作类型映射（text + color） */
export function getActionTypeMap(t: TFn) {
  return {
    LOGIN: { text: t('actionType.LOGIN'), color: ACTION_TYPE_COLORS.LOGIN },
    LOGOUT: { text: t('actionType.LOGOUT'), color: ACTION_TYPE_COLORS.LOGOUT },
    SESSION_EXPIRED: { text: t('actionType.SESSION_EXPIRED'), color: ACTION_TYPE_COLORS.SESSION_EXPIRED },
    KICKED_OUT: { text: t('actionType.KICKED_OUT'), color: ACTION_TYPE_COLORS.KICKED_OUT },
    PASSWORD_RESET: { text: t('actionType.PASSWORD_RESET'), color: ACTION_TYPE_COLORS.PASSWORD_RESET },
  };
}

/** 获取操作类型选项列表 */
export function getActionTypeOptions(t: TFn) {
  return [
    { label: t('actionType.LOGIN'), value: 'LOGIN' },
    { label: t('actionType.LOGOUT'), value: 'LOGOUT' },
    { label: t('actionType.SESSION_EXPIRED'), value: 'SESSION_EXPIRED' },
    { label: t('actionType.KICKED_OUT'), value: 'KICKED_OUT' },
    { label: t('actionType.PASSWORD_RESET'), value: 'PASSWORD_RESET' },
  ];
}

// ========== 风险等级 ==========

/** 风险等级颜色映射 */
export const RISK_LEVEL_COLORS: Record<string, string> = {
  LOW: 'success',
  MEDIUM: 'warning',
  HIGH: 'error',
};

/** 获取风险等级映射（text + color） */
export function getRiskLevelMap(t: TFn) {
  return {
    LOW: { text: t('riskLevel.LOW'), color: RISK_LEVEL_COLORS.LOW },
    MEDIUM: { text: t('riskLevel.MEDIUM'), color: RISK_LEVEL_COLORS.MEDIUM },
    HIGH: { text: t('riskLevel.HIGH'), color: RISK_LEVEL_COLORS.HIGH },
  };
}

/** 获取风险等级选项列表 */
export function getRiskLevelOptions(t: TFn) {
  return [
    { label: t('riskLevel.LOW'), value: 'LOW' },
    { label: t('riskLevel.MEDIUM'), value: 'MEDIUM' },
    { label: t('riskLevel.HIGH'), value: 'HIGH' },
  ];
}
