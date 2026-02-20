import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const site_setting: RouteRecordRaw[] = [
  {
    path: '/site-setting',
    name: 'SiteSettingManagement',
    component: BasicLayout,
    redirect: '/site-setting/site-settings',
    meta: {
      order: 1003,
      icon: 'lucide:settings-2',
      title: $t('menu.siteSetting.moduleName'),
      keepAlive: true,
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        path: 'site-settings',
        name: 'SiteSetting',
        meta: {
          order: 1,
          icon: 'lucide:settings',
          title: $t('menu.siteSetting.siteSetting'),
          authority: ['sys:platform_admin'],
        },
        component: () =>
          import('#/views/app/site_setting/site_setting/index.vue'),
      },

      {
        path: 'navigations',
        name: 'NavigationManagement',
        meta: {
          order: 2,
          icon: 'lucide:menu',
          title: $t('menu.siteSetting.navigation'),
          authority: ['sys:platform_admin'],
        },
        component: () => import('#/views/app/site_setting/navigation/index.vue'),
      },
    ],
  },
];

export default site_setting;
