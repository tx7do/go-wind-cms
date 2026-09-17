import type { PageContainerProps } from '../types';

/**
 * 检查页面权限
 * @param permission - 所需权限码（字符串或数组）
 * @param userPermissions - 用户拥有的权限码列表
 * @returns 是否有权限访问
 */
export const checkPagePermission = (
  permission: PageContainerProps['permission'],
  userPermissions: string[],
): boolean => {
  // 无权限要求 = 放行
  if (!permission) return true;

  // 转换为数组
  const required = Array.isArray(permission) ? permission : [permission];

  // 满足任一即可
  return required.some((code) => userPermissions.includes(code));
};

/**
 * 检查按钮级权限（用于 extra/footer 等操作区）
 */
export const checkButtonPermission = (
  permission: string | undefined,
  userPermissions: string[],
): boolean => {
  if (!permission) return true;
  return userPermissions.includes(permission);
};
