import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * 社交分发路由
 * 分组要求 platform_admin 或 tenant_manager，子页同权（社交账号是租户级资产）。
 */
export const socialRoutes: AppRouteObject[] = [
  {
    name: 'social',
    path: 'social',
    meta: {
      title: 'routes:social',
      icon: 'lucide:share-2',
      order: 1003,
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        name: 'social-accounts',
        path: 'accounts',
        element: createLazyRoute(() => import('@/pages/app/social/account')),
        meta: {
          title: 'routes:social-accounts',
          icon: 'lucide:link-2',
          order: 1,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'social-inbox',
        path: 'inbox',
        element: createLazyRoute(() => import('@/pages/app/social/inbox')),
        meta: {
          title: 'routes:social-inbox',
          icon: 'lucide:inbox',
          order: 3,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'social-distributions',
        path: 'distributions',
        element: createLazyRoute(() => import('@/pages/app/social/distribution')),
        meta: {
          title: 'routes:social-distributions',
          icon: 'lucide:send',
          order: 2,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
    ],
  },
];

export default socialRoutes;
