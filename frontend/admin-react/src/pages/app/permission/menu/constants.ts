/**
 * 菜单模块枚举映射常量
 * 页面和 Drawer 共用
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 菜单类型 ==========

/** 菜单类型颜色映射 */
export const MENU_TYPE_COLORS: Record<string, string> = {
  CATALOG: 'blue',
  MENU: 'green',
  BUTTON: 'orange',
  LINK: 'cyan',
  EMBEDDED: 'purple',
};

/** 获取菜单类型映射（text + color），供 render 使用 */
export function getMenuTypeMap(t: TFn) {
  return {
    CATALOG: { text: t('menuType.CATALOG'), color: MENU_TYPE_COLORS.CATALOG },
    MENU: { text: t('menuType.MENU'), color: MENU_TYPE_COLORS.MENU },
    BUTTON: { text: t('menuType.BUTTON'), color: MENU_TYPE_COLORS.BUTTON },
    LINK: { text: t('menuType.LINK'), color: MENU_TYPE_COLORS.LINK },
    EMBEDDED: { text: t('menuType.EMBEDDED'), color: MENU_TYPE_COLORS.EMBEDDED },
  };
}

/** 菜单类型选项 */
export function getMenuTypeOptions(t: TFn) {
  return [
    { label: t('menuType.CATALOG'), value: 'CATALOG' },
    { label: t('menuType.MENU'), value: 'MENU' },
    { label: t('menuType.BUTTON'), value: 'BUTTON' },
    { label: t('menuType.LINK'), value: 'LINK' },
    { label: t('menuType.EMBEDDED'), value: 'EMBEDDED' },
  ];
}

// ========== 状态映射 ==========

/** 状态颜色映射 */
export const STATUS_COLORS: Record<string, string> = {
  ON: 'success',
  OFF: 'error',
};

/** 获取状态映射 */
export function getStatusMap(t: TFn) {
  return {
    ON: { text: t('statusMap.ON'), color: STATUS_COLORS.ON },
    OFF: { text: t('statusMap.OFF'), color: STATUS_COLORS.OFF },
  };
}

/** 状态选项 */
export function getStatusOptions(t: TFn) {
  return [
    { label: t('statusMap.ON'), value: 'ON' },
    { label: t('statusMap.OFF'), value: 'OFF' },
  ];
}

// ========== 菜单类型判断辅助函数 ==========

/** 是否为目录类型 */
export function isCatalog(type?: string): boolean {
  return type === 'CATALOG';
}

/** 是否为菜单类型 */
export function isMenu(type?: string): boolean {
  return type === 'MENU';
}

/** 是否为按钮类型 */
export function isButton(type?: string): boolean {
  return type === 'BUTTON';
}

/** 是否为链接类型 */
export function isLink(type?: string): boolean {
  return type === 'LINK';
}

/** 是否为内嵌类型 */
export function isEmbedded(type?: string): boolean {
  return type === 'EMBEDDED';
}

// ========== 权限标识解析 ==========

/** 将 authority 字段统一解析为字符串数组 */
export function normalizeAuthority(authority: unknown): string[] {
  if (Array.isArray(authority)) return authority;
  if (typeof authority === 'string') {
    return authority
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean);
  }
  return [];
}

/** 获取 ProTable valueEnum 用的 status 映射 */
export const MENU_TYPE_STATUS: Record<string, string> = {
  CATALOG: 'Processing',
  MENU: 'Success',
  BUTTON: 'Warning',
  LINK: 'Default',
  EMBEDDED: 'Default',
};

export const STATUS_TYPE_STATUS: Record<string, string> = {
  ON: 'Success',
  OFF: 'Error',
};
