<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, useVbenDrawer, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';

import { notification } from 'ant-design-vue';
import dayjs from 'dayjs';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  apiClient,
  type contentservicev1_ContentModel as ContentModel,
  fetchListContentModels,
  PaginationQuery,
} from '#/api';
import { $t } from '#/locales';

import ContentModelDrawer from './content-model-drawer.vue';

const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: false,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.contentModel.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'RangePicker',
      fieldName: 'createdAt',
      label: $t('ui.table.createdAt'),
      componentProps: {
        showTime: true,
        allowClear: true,
        presets: [
          {
            label: $t('ui.dateRange.today'),
            value: [dayjs().startOf('day'), dayjs().endOf('day')],
          },
          {
            label: $t('ui.dateRange.yesterday'),
            value: [
              dayjs().subtract(1, 'day').startOf('day'),
              dayjs().subtract(1, 'day').endOf('day'),
            ],
          },
          {
            label: $t('ui.dateRange.thisWeek'),
            value: [dayjs().startOf('week'), dayjs().endOf('week')],
          },
          {
            label: $t('ui.dateRange.lastWeek'),
            value: [
              dayjs().subtract(1, 'week').startOf('week'),
              dayjs().subtract(1, 'week').endOf('week'),
            ],
          },
          {
            label: $t('ui.dateRange.thisMonth'),
            value: [dayjs().startOf('month'), dayjs().endOf('month')],
          },
          {
            label: $t('ui.dateRange.lastMonth'),
            value: [
              dayjs().subtract(1, 'month').startOf('month'),
              dayjs().subtract(1, 'month').endOf('month'),
            ],
          },
        ],
      },
    },
  ],
};

const gridOptions: VxeGridProps<ContentModel> = {
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
        let startTime: any;
        let endTime: any;
        if (
          formValues.createdAt !== undefined &&
          formValues.createdAt.length === 2
        ) {
          startTime = dayjs(formValues.createdAt[0]).format(
            'YYYY-MM-DD HH:mm:ss',
          );
          endTime = dayjs(formValues.createdAt[1]).format(
            'YYYY-MM-DD HH:mm:ss',
          );
        }

        return await fetchListContentModels(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues: {
              name: formValues.name,
              created_at__gte: startTime,
              created_at__lte: endTime,
            },
            orderBy: ['-created_at'],
          }),
        );
      },
    },
  },

  columns: [
    {
      title: $t('page.contentModel.name'),
      field: 'name',
      minWidth: 120,
    },
    {
      title: $t('page.contentModel.code'),
      field: 'code',
      minWidth: 120,
    },
    {
      title: $t('page.contentModel.description'),
      field: 'description',
      minWidth: 200,
      formatter: ({ cellValue }) => {
        if (!cellValue) return '';
        const text = String(cellValue);
        return text.length > 50 ? `${text.slice(0, 50)}...` : text;
      },
    },
    {
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      width: 160,
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
  connectedComponent: ContentModelDrawer,
  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      gridApi.reload();
    }
  },
});

function openDrawer(row?: any) {
  drawerApi.setData({ row });
  drawerApi.open();
}

function handleEdit(row: any) {
  openDrawer(row);
}

function handleDelete(row: any) {
  (async () => {
    try {
      await apiClient.contentModelService.Delete({ id: row.id });
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
    <Grid :table-title="$t('page.contentModel.moduleName')">
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
              moduleName: $t('page.contentModel.moduleName'),
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
