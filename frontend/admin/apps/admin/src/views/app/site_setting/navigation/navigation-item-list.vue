<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h, watch } from 'vue';

import { Page, useVbenDrawer, type VbenFormProps } from '@vben/common-ui';
import { IconifyIcon, LucideFilePenLine, LucideTrash2 } from '@vben/icons';
import { isEqual } from '@vben/utils';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  apiClient,
  type siteservicev1_NavigationItem as NavigationItem,
  navigationItemLinkTypeList,
  navigationItemLinkTypeToColor,
  navigationItemLinkTypeToName,
} from '#/api';
import { $t } from '#/locales';

import NavigationItemDrawer from './navigation-item-drawer.vue';
import { useNavigationViewStore } from './navigation-view.state';

const navigationViewStore = useNavigationViewStore();

const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: false,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'title',
      label: $t('page.navigationItem.title'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'linkType',
      label: $t('page.navigationItem.linkType'),
      componentProps: {
        options: navigationItemLinkTypeList,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
      },
    },
  ],
};

const gridOptions: VxeGridProps<NavigationItem> = {
  toolbarConfig: {
    custom: true,
    export: true,
    refresh: true,
    zoom: true,
  },
  height: 'auto',
  exportConfig: {},
  pagerConfig: {},
  rowConfig: {
    isHover: true,
  },
  stripe: true,

  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        console.log(
          'navigation item list query:',
          formValues,
          navigationViewStore.currentNavigationId,
        );

        return await navigationViewStore.fetchItemList(
          navigationViewStore.currentNavigationId,
          page.currentPage,
          page.pageSize,
          formValues,
        );
      },
    },
  },

  columns: [
    {
      title: $t('page.navigationItem.title'),
      field: 'title',
      minWidth: 140,
      fixed: 'left',
      align: 'left',
    },
    {
      title: $t('page.navigationItem.description'),
      field: 'description',
      align: 'left',
      minWidth: 140,
    },
    {
      title: $t('page.navigationItem.icon'),
      field: 'icon',
      slots: {
        default: 'icon',
      },
      minWidth: 110,
    },
    {
      title: $t('page.navigationItem.url'),
      field: 'url',
      align: 'left',
      minWidth: 110,
    },
    {
      title: $t('page.navigationItem.linkType'),
      field: 'linkType',
      slots: {
        default: 'linkType',
      },
      minWidth: 110,
    },
    {
      title: $t('page.navigationItem.isOpenNewTab'),
      field: 'isOpenNewTab',
      minWidth: 100,
      formatter: ({ cellValue }) =>
        cellValue ? $t('ui.button.yes') : $t('ui.button.no'),
    },
    {
      title: $t('page.navigationItem.isInvalid'),
      field: 'isInvalid',
      minWidth: 100,
      formatter: ({ cellValue }) =>
        cellValue ? $t('ui.button.yes') : $t('ui.button.no'),
    },
    {
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      minWidth: 160,
    },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      minWidth: 90,
    },
  ],
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

const [Drawer, drawerApi] = useVbenDrawer({
  connectedComponent: NavigationItemDrawer,
  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      gridApi.reload();
    }
  },
});

function openDrawer(create: boolean, row?: any) {
  console.log('open drawer:', create, row);
  drawerApi.setData({
    create,
    row,
  });
  drawerApi.open();
}

function handleCreate() {
  openDrawer(true);
}

function handleEdit(row: any) {
  openDrawer(false, row);
}

async function handleDelete(row: any) {
  try {
    await apiClient.navigationItemService.Delete({ id: row.id });
    notification.success({
      message: $t('ui.notification.delete_success'),
    });
    await gridApi.reload();
  } catch {
    notification.error({
      message: $t('ui.notification.delete_failed'),
    });
  }
}

function getIcon(row: any) {
  if (!row.icon) return '';
  // 如果已有前缀（如carbon:xxx），直接返回
  if (row.icon.includes(':')) return row.icon;
  // 否则加上carbon:前缀
  return `carbon:${row.icon}`;
}

watch(
  () => navigationViewStore.needReloadNavigationItemList,
  (newValues, oldValue) => {
    if (isEqual(newValues, oldValue) || !newValues) {
      return;
    }
    navigationViewStore.needReloadNavigationItemList = false;
    gridApi.reload();
  },
);
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('page.navigationItem.moduleName')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('page.navigationItem.button.create') }}
        </a-button>
      </template>
      <template #linkType="{ row }">
        <a-tag :color="navigationItemLinkTypeToColor(row.linkType)">
          {{ navigationItemLinkTypeToName(row.linkType) }}
        </a-tag>
      </template>

      <template #icon="{ row }">
        <IconifyIcon v-if="row.icon" :icon="getIcon(row)" class="size-6" />
      </template>
      <template #action="{ row }">
        <a-button
          type="link"
          :icon="h(LucideFilePenLine)"
          @click.stop="handleEdit(row)"
        />
        <a-popconfirm
          :cancel-text="$t('ui.button.cancel')"
          :ok-text="$t('ui.button.ok')"
          :title="
            $t('ui.text.do_you_want_delete', {
              moduleName: $t('page.navigationItem.moduleName'),
            })
          "
          @confirm="handleDelete(row)"
        >
          <a-button danger type="link" :icon="h(LucideTrash2)" />
        </a-popconfirm>
      </template>
    </Grid>
    <Drawer />
  </Page>
</template>

<style scoped></style>
