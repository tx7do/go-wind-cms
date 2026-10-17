import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * 通知域路由配置：渠道配置、路由规则（事件 → 渠道）、投递台账、通知模板。
 * 站内信只是本域的一个 INTERNAL 渠道，故这几页与 /internal-message 平级而非挂在其下。
 */
export const notificationRoutes: AppRouteObject[] = [
  {
    name: 'notification',
    path: 'notification', // 相对路径，会自动拼接到父路由 '/'
    meta: {
      title: 'routes:notification',
      icon: 'lucide:bell-ring', // Iconify 格式
      order: 2007,
      keepAlive: true, // 保持组件状态
      authority: ['sys:platform_admin'], // 平台级配置，仅平台管理员
    },
    children: [
      {
        name: 'notification-channels',
        path: 'channels', // 相对路径，最终为 /notification/channels
        element: createLazyRoute(() => import('@/pages/app/notification/channel')),
        meta: {
          title: 'routes:notification-channels',
          icon: 'lucide:mail', // Iconify 格式
          order: 1,
          authority: ['sys:platform_admin'],
        },
      },
      {
        name: 'notification-rules',
        path: 'rules', // 相对路径，最终为 /notification/rules
        element: createLazyRoute(() => import('@/pages/app/notification/rule')),
        meta: {
          title: 'routes:notification-rules',
          icon: 'lucide:git-branch',
          order: 2,
          authority: ['sys:platform_admin'],
        },
      },
      {
        name: 'notification-deliveries',
        path: 'deliveries', // 相对路径，最终为 /notification/deliveries
        element: createLazyRoute(() => import('@/pages/app/notification/delivery')),
        meta: {
          title: 'routes:notification-deliveries',
          icon: 'lucide:send',
          order: 3,
          authority: ['sys:platform_admin'],
        },
      },
      {
        name: 'notification-templates',
        path: 'templates', // 相对路径，最终为 /notification/templates
        element: createLazyRoute(() => import('@/pages/app/notification/template')),
        meta: {
          title: 'routes:notification-templates',
          icon: 'lucide:layout-template',
          order: 4,
          authority: ['sys:platform_admin'],
        },
      },
    ],
  },
];

export default notificationRoutes;
