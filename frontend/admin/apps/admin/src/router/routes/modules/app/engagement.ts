import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const engagement: RouteRecordRaw[] = [
  {
    path: '/engagement',
    name: 'EngagementManagement',
    component: BasicLayout,
    redirect: '/engagement/navigations',
    meta: {
      order: 1002,
      icon: 'lucide:message-square',
      title: $t('menu.engagement.moduleName'),
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        path: 'comments',
        name: 'CommentManagement',
        meta: {
          order: 1,
          icon: 'lucide:message-square',
          title: $t('menu.engagement.comment'),
          authority: ['sys:platform_admin'],
        },
        component: () => import('#/views/app/engagement/comment/index.vue'),
      },

    ],
  },
];

export default engagement;
