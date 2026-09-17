/**
 * 登录策略模块枚举映射常量
 * 页面和 Drawer 共用
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 策略类型 ==========

/** 策略类型颜色映射 */
export const POLICY_TYPE_COLORS: Record<string, string> = {
  BLACKLIST: 'error',
  WHITELIST: 'success',
};

/** 获取策略类型映射（text + color） */
export function getPolicyTypeMap(t: TFn) {
  return {
    BLACKLIST: { text: t('policyType.BLACKLIST'), color: POLICY_TYPE_COLORS.BLACKLIST },
    WHITELIST: { text: t('policyType.WHITELIST'), color: POLICY_TYPE_COLORS.WHITELIST },
  };
}

/** 策略类型选项 */
export function getPolicyTypeOptions(t: TFn) {
  return [
    { label: t('policyType.BLACKLIST'), value: 'BLACKLIST' },
    { label: t('policyType.WHITELIST'), value: 'WHITELIST' },
  ];
}

// ========== 登录方式 ==========

/** 登录方式颜色映射 */
export const POLICY_METHOD_COLORS: Record<string, string> = {
  IP: 'blue',
  MAC: 'green',
  REGION: 'cyan',
  TIME: 'purple',
  DEVICE: 'orange',
};

/** 获取登录方式映射（text + color） */
export function getPolicyMethodMap(t: TFn) {
  return {
    IP: { text: t('policyMethod.IP'), color: POLICY_METHOD_COLORS.IP },
    MAC: { text: t('policyMethod.MAC'), color: POLICY_METHOD_COLORS.MAC },
    REGION: { text: t('policyMethod.REGION'), color: POLICY_METHOD_COLORS.REGION },
    TIME: { text: t('policyMethod.TIME'), color: POLICY_METHOD_COLORS.TIME },
    DEVICE: { text: t('policyMethod.DEVICE'), color: POLICY_METHOD_COLORS.DEVICE },
  };
}

/** 登录方式选项 */
export function getPolicyMethodOptions(t: TFn) {
  return [
    { label: t('policyMethod.IP'), value: 'IP' },
    { label: t('policyMethod.MAC'), value: 'MAC' },
    { label: t('policyMethod.REGION'), value: 'REGION' },
    { label: t('policyMethod.TIME'), value: 'TIME' },
    { label: t('policyMethod.DEVICE'), value: 'DEVICE' },
  ];
}
