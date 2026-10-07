import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * 互动运营路由
 * 权限模型对齐 CMS vben 侧（authority 任一匹配即放行）：
 * 分组要求 platform_admin 或 tenant_manager，子页均要求 platform_admin。
 */
export const engagementRoutes: AppRouteObject[] = [
  {
    name: 'engagement',
    path: 'engagement', // 相对路径，会自动拼接到父路由 '/'
    meta: {
      title: 'routes:engagement',
      icon: 'lucide:message-square',
      order: 1002,
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        name: 'comments',
        path: 'comments', // 最终为 /engagement/comments
        element: createLazyRoute(() => import('@/pages/app/engagement/comment')),
        meta: {
          title: 'routes:comments',
          icon: 'lucide:message-square',
          order: 1,
          authority: ['sys:platform_admin'],
        },
      },
    ],
  },
];

export default engagementRoutes;
