import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * 媒体中心路由
 * 权限模型对齐 CMS vben 侧（authority 任一匹配即放行）：
 * 分组要求 platform_admin 或 tenant_manager，子页均要求 platform_admin。
 */
export const mediaRoutes: AppRouteObject[] = [
  {
    name: 'media',
    path: 'media', // 相对路径，会自动拼接到父路由 '/'
    meta: {
      title: 'routes:media',
      icon: 'lucide:image',
      order: 1001,
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        name: 'media-assets',
        path: 'media-assets', // 最终为 /media/media-assets
        element: createLazyRoute(() => import('@/pages/app/media/media_asset')),
        meta: {
          title: 'routes:media-assets',
          icon: 'lucide:image',
          order: 1,
          authority: ['sys:platform_admin'],
        },
      },
    ],
  },
];

export default mediaRoutes;
