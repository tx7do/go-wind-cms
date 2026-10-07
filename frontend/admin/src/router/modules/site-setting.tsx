import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * 站点配置路由
 * 权限模型对齐 CMS vben 侧（authority 任一匹配即放行）：
 * 分组要求 platform_admin 或 tenant_manager，子页均要求 platform_admin。
 */
export const siteSettingRoutes: AppRouteObject[] = [
  {
    name: 'site-setting',
    path: 'site-setting', // 相对路径，会自动拼接到父路由 '/'
    meta: {
      title: 'routes:site-setting',
      icon: 'lucide:settings-2',
      order: 1003,
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        name: 'site-settings',
        path: 'site-settings', // 最终为 /site-setting/site-settings
        element: createLazyRoute(() => import('@/pages/app/site_setting/setting')),
        meta: {
          title: 'routes:site-settings',
          icon: 'lucide:settings',
          order: 1,
          authority: ['sys:platform_admin'],
        },
      },
      {
        name: 'navigations',
        path: 'navigations', // 最终为 /site-setting/navigations
        element: createLazyRoute(() => import('@/pages/app/site_setting/navigation')),
        meta: {
          title: 'routes:navigations',
          icon: 'lucide:menu',
          order: 2,
          authority: ['sys:platform_admin'],
        },
      },
      {
        name: 'sites',
        path: 'sites', // 相对路径，最终为 /site-setting/sites
        element: createLazyRoute(() => import('@/pages/app/site_setting/site')),
        meta: {
          title: 'routes:sites',
          icon: 'lucide:globe',
          order: 3,
          authority: ['sys:platform_admin'],
        },
      },
    ],
  },
];

export default siteSettingRoutes;
