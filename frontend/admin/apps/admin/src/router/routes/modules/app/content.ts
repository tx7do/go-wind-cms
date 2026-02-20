import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const content: RouteRecordRaw[] = [
  {
    path: '/content',
    name: 'ContentManagement',
    component: BasicLayout,
    redirect: '/content/posts',
    meta: {
      order: 1000,
      icon: 'lucide:file-text',
      title: $t('menu.content.moduleName'),
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        path: 'posts',
        name: 'PostManagement',
        meta: {
          order: 1,
          icon: 'lucide:newspaper',
          title: $t('menu.content.post'),
          authority: ['sys:platform_admin'],
        },
        component: () => import('#/views/app/content/post/index.vue'),
      },

      {
        path: 'pages',
        name: 'PageManagement',
        meta: {
          order: 2,
          icon: 'lucide:file',
          title: $t('menu.content.page'),
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
        component: () => import('#/views/app/content/page/index.vue'),
      },

      {
        path: 'categories',
        name: 'CategoryManagement',
        meta: {
          order: 3,
          icon: 'lucide:folder',
          title: $t('menu.content.category'),
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
        component: () => import('#/views/app/content/category/index.vue'),
      },

      {
        path: 'tags',
        name: 'TagManagement',
        meta: {
          order: 4,
          icon: 'lucide:tags',
          title: $t('menu.content.tag'),
          authority: ['sys:platform_admin', 'sys:tenant_manager'],
        },
        component: () => import('#/views/app/content/tag/index.vue'),
      },
    ],
  },
];

export default content;
