import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * 内容管理路由
 * 权限模型对齐 CMS vben 侧（authority 任一匹配即放行）：
 * 分组要求 platform_admin 或 tenant_manager；列表子页 platform_admin，
 * 创建/编辑页 tenant_manager 亦可访问（与 vben 侧一致）。
 */
export const contentRoutes: AppRouteObject[] = [
  {
    name: 'content',
    path: 'content', // 相对路径，会自动拼接到父路由 '/'
    meta: {
      title: 'routes:content',
      icon: 'lucide:file-text',
      order: 1000,
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        name: 'content-categories',
        path: 'categories',
        element: createLazyRoute(() => import('@/pages/app/content/category')),
        meta: {
          title: 'routes:content-categories',
          icon: 'lucide:folder',
          order: 3,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-create-category',
        path: 'categories/create',
        element: createLazyRoute(() => import('@/pages/app/content/category/edit/CategoryEditPage')),
        meta: {
          title: 'routes:content-create-category',
          hideInMenu: true,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-edit-category',
        path: 'categories/edit/:id',
        element: createLazyRoute(() => import('@/pages/app/content/category/edit/CategoryEditPage')),
        meta: {
          title: 'routes:content-edit-category',
          hideInMenu: true,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-models',
        path: 'models',
        element: createLazyRoute(() => import('@/pages/app/content/model')),
        meta: {
          title: 'routes:content-models',
          icon: 'lucide:database',
          order: 9,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-pages',
        path: 'pages',
        element: createLazyRoute(() => import('@/pages/app/content/page')),
        meta: {
          title: 'routes:content-pages',
          icon: 'lucide:layout-template',
          order: 2,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-create-page',
        path: 'pages/create',
        element: createLazyRoute(() => import('@/pages/app/content/page/edit/PageEditPage')),
        meta: {
          title: 'routes:content-create-page',
          hideInMenu: true,
          hideInTab: false,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-edit-page',
        path: 'pages/edit/:id',
        element: createLazyRoute(() => import('@/pages/app/content/page/edit/PageEditPage')),
        meta: {
          title: 'routes:content-edit-page',
          hideInMenu: true,
          hideInTab: false,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-posts',
        path: 'posts',
        element: createLazyRoute(() => import('@/pages/app/content/post')),
        meta: {
          title: 'routes:content-posts',
          icon: 'lucide:newspaper',
          order: 1,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-create-post',
        path: 'posts/create',
        element: createLazyRoute(() => import('@/pages/app/content/post/edit/PostEditPage')),
        meta: {
          title: 'routes:content-create-post',
          hideInMenu: true,
          hideInTab: false,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-edit-post',
        path: 'posts/edit/:id',
        element: createLazyRoute(() => import('@/pages/app/content/post/edit/PostEditPage')),
        meta: {
          title: 'routes:content-edit-post',
          hideInMenu: true,
          hideInTab: false,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-tags',
        path: 'tags',
        element: createLazyRoute(() => import('@/pages/app/content/tag')),
        meta: {
          title: 'routes:content-tags',
          icon: 'lucide:tags',
          order: 4,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-create-tag',
        path: 'tags/create',
        element: createLazyRoute(() => import('@/pages/app/content/tag/edit/TagEditPage')),
        meta: {
          title: 'routes:content-create-tag',
          hideInMenu: true,
          hideInTab: false,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
      {
        name: 'content-edit-tag',
        path: 'tags/edit/:id',
        element: createLazyRoute(() => import('@/pages/app/content/tag/edit/TagEditPage')),
        meta: {
          title: 'routes:content-edit-tag',
          hideInMenu: true,
          hideInTab: false,
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
      },
    ],
  },
];

export default contentRoutes;
