import { Navigate } from 'react-router-dom';

import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * 概览（仪表盘）路由：目录 + 分析/工作台叶子，与上游 go-wind-admin 及菜单种子的
 * Dashboard(CATALOG) → Analytics/Workspace(MENU) 结构对齐。
 */
export const dashboardRoutes: AppRouteObject[] = [
  {
    name: 'dashboard',
    path: 'dashboard', // 相对路径，会自动拼接到父路由 '/'
    meta: {
      title: 'routes:dashboard',
      icon: 'lucide:layout-dashboard', // Iconify 格式
      order: 1,
    },
    children: [
      // /dashboard 直入重定向到分析页（react-router 不允许嵌套绝对路径，用 index 路由承担）
      {
        index: true,
        element: <Navigate to="/dashboard/analytics" replace />,
        meta: { title: 'routes:dashboard', hideInMenu: true, hideInTab: true },
      },
      {
        name: 'analytics',
        path: 'analytics', // /dashboard/analytics
        element: createLazyRoute(() => import('@/pages/app/dashboard/analytics')),
        meta: {
          title: 'routes:dashboard-analytics',
          icon: 'lucide:area-chart',
          order: 1,
          affixTab: true, // 与上游及种子菜单一致：分析页固定标签
        },
      },
      {
        name: 'workspace',
        path: 'workspace', // /dashboard/workspace
        element: createLazyRoute(() => import('@/pages/app/dashboard/workspace')),
        meta: {
          title: 'routes:dashboard-workspace',
          icon: 'lucide:armchair',
          order: 2,
        },
      },
    ],
  },
];

export default dashboardRoutes;
