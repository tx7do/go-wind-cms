import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const media: RouteRecordRaw[] = [
  {
    path: '/media',
    name: 'MediaManagement',
    component: BasicLayout,
    redirect: '/media/navigations',
    meta: {
      order: 1001,
      icon: 'lucide:image',
      title: $t('menu.media.moduleName'),
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        path: 'media-assets',
        name: 'MediaAssetManagement',
        meta: {
          order: 1,
          icon: 'lucide:image',
          title: $t('menu.media.mediaAsset'),
          authority: ['sys:platform_admin'],
        },
        component: () => import('#/views/app/media/media_asset/index.vue'),
      },

    ],
  },
];

export default media;
