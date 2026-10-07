import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * 系统管理路由配置
 * 包括菜单管理、API管理、字典管理、文件管理、任务管理、登录策略、语言管理等页面
 */
export const systemRoutes: AppRouteObject[] = [
  {
    name: 'system',
    path: 'system', // 相对路径，会自动拼接到父路由 '/'
    meta: {
      title: 'routes:system',
      icon: 'lucide:settings', // Iconify 格式
      order: 2005,
      keepAlive: true, // 保持组件状态
      // permission: 'sys:platform_admin', // 平台管理员或租户管理员权限（开发阶段暂时注释）
    },
    children: [
      {
        name: 'dict',
        path: 'dict', // 相对路径，最终为 /system/dict
        element: createLazyRoute(() => import('@/pages/app/system/dict')),
        meta: {
          title: 'routes:dict',
          icon: 'lucide:library-big', // Iconify 格式
          order: 3,
          // permission: 'sys:platform_admin', // 仅平台管理员权限（开发阶段暂时注释）
        },
      },
      {
        name: 'files',
        path: 'files', // 相对路径，最终为 /system/files
        element: createLazyRoute(() => import('@/pages/app/system/file')),
        meta: {
          title: 'routes:files',
          icon: 'lucide:file-search', // Iconify 格式
          order: 4,
          // permission: 'sys:platform_admin', // 平台管理员或租户管理员权限（开发阶段暂时注释）
        },
      },
      {
        name: 'tasks',
        path: 'tasks', // 相对路径，最终为 /system/tasks
        element: createLazyRoute(() => import('@/pages/app/system/task')),
        meta: {
          title: 'routes:tasks',
          icon: 'lucide:list-todo', // Iconify 格式
          order: 5,
          // permission: 'sys:platform_admin', // 平台管理员或租户管理员权限（开发阶段暂时注释）
        },
      },
      {
        name: 'login-policies',
        path: 'login-policies', // 相对路径，最终为 /system/login-policies
        element: createLazyRoute(() => import('@/pages/app/system/login-policy')),
        meta: {
          title: 'routes:login-policies',
          icon: 'lucide:shield-x', // Iconify 格式
          order: 6,
          // permission: 'sys:platform_admin', // 仅平台管理员权限（开发阶段暂时注释）
        },
      },
      {
        name: 'languages',
        path: 'languages', // 相对路径，最终为 /system/languages
        element: createLazyRoute(() => import('@/pages/app/system/language')),
        meta: {
          title: 'routes:languages',
          icon: 'lucide:globe', // Iconify 格式
          order: 7,
          // permission: 'sys:platform_admin', // 仅平台管理员权限（开发阶段暂时注释）
        },
      },
                                        ],
  },
];

export default systemRoutes;
