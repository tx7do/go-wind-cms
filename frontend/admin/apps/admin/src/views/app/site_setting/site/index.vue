<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, useVbenDrawer, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  apiClient,
  fetchListSites,
  PaginationQuery,
  type siteservicev1_Site as Site,
  siteStatusList,
  siteStatusToColor,
  siteStatusToName,
} from '#/api';
import { $t } from '#/locales';

import SiteDrawer from './site-drawer.vue';

const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: false,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.site.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Input',
      fieldName: 'slug',
      label: $t('page.site.slug'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Input',
      fieldName: 'domain',
      label: $t('page.site.domain'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('page.site.status'),
      componentProps: {
        options: siteStatusList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
  ],
};

const gridOptions: VxeGridProps<Site> = {
  toolbarConfig: {
    custom: true,
    export: true,
    // import: true,
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
        return await fetchListSites(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues: {
              name: formValues.name,
              slug: formValues.slug,
              domain: formValues.domain,
              status: formValues.status,
            },
            orderBy: ['-created_at'],
          }),
        );
      },
    },
  },

  columns: [
    {
      title: $t('page.site.name'),
      field: 'name',
      fixed: 'left',
      align: 'left',
      minWidth: 140,
    },
    {
      title: $t('page.site.slug'),
      field: 'slug',
      align: 'left',
      minWidth: 120,
    },
    {
      title: $t('page.site.domain'),
      field: 'domain',
      align: 'left',
      minWidth: 140,
    },
    {
      title: $t('page.site.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 110,
    },
    {
      title: $t('page.site.defaultLocale'),
      field: 'defaultLocale',
      width: 100,
    },
    {
      title: $t('page.site.template'),
      field: 'template',
      minWidth: 120,
    },
    {
      title: $t('page.site.theme'),
      field: 'theme',
      minWidth: 120,
    },
    {
      title: $t('page.site.isDefault'),
      field: 'isDefault',
      width: 100,
      formatter: ({ cellValue }) =>
        cellValue ? $t('ui.button.yes') : $t('ui.button.no'),
    },
    {
      title: $t('page.site.visitCount'),
      field: 'visitCount',
      width: 90,
      align: 'right',
      headerAlign: 'right',
      sortable: true,
    },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      width: 100,
    },
  ],
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

const [Drawer, drawerApi] = useVbenDrawer({
  // 连接抽离的组件
  connectedComponent: SiteDrawer,

  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      // 关闭时，重载表格数据
      gridApi.reload();
    }
  },
});

function openDrawer(create: boolean, row?: any) {
  drawerApi.setData({
    create,
    row,
  });
  drawerApi.open();
}

/* 创建 */
function handleCreate() {
  openDrawer(true);
}

/* 编辑 */
function handleEdit(row: any) {
  openDrawer(false, row);
}

/* 删除 */
function handleDelete(row: any) {
  (async () => {
    try {
      await apiClient.siteService.Delete({ id: row.id });
      notification.success({ message: $t('ui.notification.delete_success') });
      await gridApi.reload();
    } catch {
      notification.error({ message: $t('ui.notification.delete_failed') });
    }
  })();
}
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('menu.siteSetting.site')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('page.site.button.create') }}
        </a-button>
      </template>
      <template #status="{ row }">
        <a-tag :color="siteStatusToColor(row.status)">
          {{ siteStatusToName(row.status) }}
        </a-tag>
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
              moduleName: $t('page.site.moduleName'),
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
