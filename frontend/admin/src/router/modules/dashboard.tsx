import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

export const dashboardRoutes: AppRouteObject[] = [
  {
    name: 'dashboard',
    path: 'dashboard', // 相对路径，会自动拼接到父路由 '/'
    element: createLazyRoute(() => import('@/pages/app/dashboard')),
    meta: {
      title: 'routes:dashboard',
      icon: 'lucide:layout-dashboard', // Iconify 格式
      order: 1,
      hideInMenu: false,
    },
  },
  {
    name: 'workspace',
    path: 'workspace', // 相对路径，最终为 /workspace
    element: createLazyRoute(() => import('@/pages/app/dashboard/workspace')),
    meta: {
      title: 'routes:workspace',
      icon: 'lucide:rocket', // Iconify 格式
      order: 2,
      hideInMenu: false,
    },
  },
];

export default dashboardRoutes;
