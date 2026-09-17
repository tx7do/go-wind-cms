/**
 * 职位模块枚举映射常量
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 职位状态（ON/OFF） ==========

export const STATUS_COLORS: Record<string, string> = {
  ON: 'success',
  OFF: 'error',
};

export function getStatusMap(t: TFn) {
  return {
    ON: { text: t('statusMap.ON'), color: STATUS_COLORS.ON },
    OFF: { text: t('statusMap.OFF'), color: STATUS_COLORS.OFF },
  };
}

export function getStatusOptions(t: TFn) {
  return [
    { label: t('statusMap.ON'), value: 'ON' },
    { label: t('statusMap.OFF'), value: 'OFF' },
  ];
}

// ========== 职位类型 ==========

export const POSITION_TYPE_COLORS: Record<string, string> = {
  REGULAR: 'blue',
  LEADER: 'purple',
  MANAGER: 'cyan',
  INTERN: 'orange',
  CONTRACT: 'green',
  OTHER: 'default',
};

export function getPositionTypeMap(t: TFn) {
  return {
    REGULAR: { text: t('typeMap.REGULAR'), color: POSITION_TYPE_COLORS.REGULAR },
    LEADER: { text: t('typeMap.LEADER'), color: POSITION_TYPE_COLORS.LEADER },
    MANAGER: { text: t('typeMap.MANAGER'), color: POSITION_TYPE_COLORS.MANAGER },
    INTERN: { text: t('typeMap.INTERN'), color: POSITION_TYPE_COLORS.INTERN },
    CONTRACT: { text: t('typeMap.CONTRACT'), color: POSITION_TYPE_COLORS.CONTRACT },
    OTHER: { text: t('typeMap.OTHER'), color: POSITION_TYPE_COLORS.OTHER },
  };
}

export function getPositionTypeOptions(t: TFn) {
  return [
    { label: t('typeMap.REGULAR'), value: 'REGULAR' },
    { label: t('typeMap.LEADER'), value: 'LEADER' },
    { label: t('typeMap.MANAGER'), value: 'MANAGER' },
    { label: t('typeMap.INTERN'), value: 'INTERN' },
    { label: t('typeMap.CONTRACT'), value: 'CONTRACT' },
    { label: t('typeMap.OTHER'), value: 'OTHER' },
  ];
}

// ========== 组织树构建 ==========

/**
 * 将 API 返回的树形组织数据转为 TreeSelect 格式
 */
export function buildOrgTreeData(items: any[]): any[] {
  const result: any[] = [];

  for (const item of items) {
    const node: any = {
      id: item.id,
      key: item.id,
      value: item.id,
      title: item.name || String(item.id),
      label: item.name || String(item.id),
    };

    if (item.children && item.children.length > 0) {
      const children = buildOrgTreeData(item.children);
      if (children.length > 0) {
        node.children = children;
      }
    }

    result.push(node);
  }

  return result;
}
