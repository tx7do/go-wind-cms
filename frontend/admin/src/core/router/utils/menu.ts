import type { ProLayoutProps } from '@ant-design/pro-components';
import type { AppRouteObject } from '@/core/router/types';

type MenuRoute = AppRouteObject;

/**
 * 递归排序菜单树（按 order 字段）
 */
const sortMenuTree = (menuList: any[]): any[] => {
  return menuList
    .map((menu) => ({
      ...menu,
      // 递归排序子菜单
      children: menu.children?.length ? sortMenuTree(menu.children) : undefined,
    }))
    .sort((a, b) => {
      // order 越小越靠前，undefined 排最后
      const orderA = a.order ?? Number.MAX_SAFE_INTEGER;
      const orderB = b.order ?? Number.MAX_SAFE_INTEGER;
      return orderA - orderB;
    });
};

export const transformRoutesToMenu = (
  routes: MenuRoute[],
  permissions: string[],
  parentPath: string = '',
): NonNullable<ProLayoutProps['route']>['routes'] => {
  const menus = routes
    .filter((route) => {
      // 过滤隐藏菜单
      if (route.meta?.hideInMenu) return false;

      const meta = route.meta;
      // 权限过滤与路由闸（generate-routes-frontend.hasPermission）同语义：
      // 无权限要求=保留；持有任一权限码=保留；两者皆否但声明了
      // menuVisibleWithForbidden 的路由保留进菜单（点击后由路由侧渲染 403），
      // 对齐 vue-vben / vue-element 的"菜单可见但禁止访问"
      if (!meta?.authority?.length) return true;
      if (meta.authority.some((code: string) => permissions.includes(code))) return true;
      return meta?.menuVisibleWithForbidden === true;
    })
    .map((route) => {
      // 处理路径：将相对路径转换为绝对路径
      const fullPath = route.path?.startsWith('/')
        ? route.path
        : `${parentPath}/${route.path}`.replace(/\/+/g, '/');

      const menuItem: any = {
        path: fullPath, // 使用完整路径作为 key
        name: route.label || route.meta?.title,
        icon: route.meta?.icon,
        order: route.meta?.order, // 提取排序字段
      };

      if (route.children) {
        menuItem.children = transformRoutesToMenu(route.children, permissions, fullPath);
      }

      return menuItem;
    })
    .filter(Boolean);

  // 排序菜单树
  return sortMenuTree(menus);
};
