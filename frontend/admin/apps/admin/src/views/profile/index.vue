<script lang="ts" setup>
import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

import BaseSettingPage from './base-setting-page.vue';
import EditPasswordPage from './edit-password-page.vue';
import LoginRecordsPage from './login-records-page.vue';

// securitySettings/accountBind/notification 三个原占位 tab 已移除：
// 均为无 API 的纯静态假数据（安全设置/账号绑定后端为 501 未实现，
// 通知开关无持久化），保留只会误导用户。待后端能力就绪再恢复。
const settingList = [
  {
    key: '1',
    name: $t('page.user.profile.tab.basicSettings'),
    component: BaseSettingPage,
  },
  {
    key: '2',
    name: $t('page.user.profile.tab.editPassword'),
    component: EditPasswordPage,
  },
  {
    key: '3',
    name: $t('page.user.profile.tab.loginRecords'),
    component: LoginRecordsPage,
  },
];
</script>

<template>
  <Page auto-content-height>
    <a-card>
      <a-tabs
        tab-position="left"
        :tab-bar-style="{ width: '220px' }"
        class="edge-tabs"
      >
        <template v-for="item in settingList" :key="item.key">
          <a-tab-pane :tab="item.name">
            <component :is="item.component" />
          </a-tab-pane>
        </template>
      </a-tabs>
    </a-card>
  </Page>
</template>

<style lang="less">
.edge-tabs {
  margin: 0;
}

/* 使用 Vue scoped 的深度选择器，覆盖 Antd 的默认内边距/外边距 */
.edge-tabs,
::v-deep(.ant-tabs-content, .ant-tabs-content-holder, .ant-tabs-tabpane) {
  padding: 0 !important;
  margin: 0 !important;
  box-sizing: border-box;
}
</style>
