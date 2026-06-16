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
  type commentservicev1_Comment as Comment,
  commentAuthorTypeList,
  commentAuthorTypeToColor,
  commentAuthorTypeToName,
  commentContentTypeList,
  commentContentTypeToColor,
  commentContentTypeToName,
  commentStatusList,
  commentStatusToColor,
  commentStatusToName,
  fetchListComments,
  PaginationQuery,
} from '#/api';
import { $t } from '#/locales';

import CommentDrawer from './comment-drawer.vue';

const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: false,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'authorName',
      label: $t('page.comment.authorName'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'authorType',
      label: $t('page.comment.authorType'),
      componentProps: {
        options: commentAuthorTypeList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'contentType',
      label: $t('page.comment.contentType'),
      componentProps: {
        options: commentContentTypeList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('page.comment.status'),
      componentProps: {
        options: commentStatusList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
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

const gridOptions: VxeGridProps<Comment> = {
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

        return await fetchListComments(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues: {
              authorName: formValues.authorName,
              authorType: formValues.authorType,
              contentType: formValues.contentType,
              status: formValues.status,
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
      title: $t('page.comment.authorName'),
      field: 'authorName',
      minWidth: 120,
    },
    {
      title: $t('page.comment.content'),
      field: 'content',
      minWidth: 200,
      formatter: ({ cellValue }) => {
        if (!cellValue) return '';
        const text = String(cellValue);
        return text.length > 50 ? `${text.slice(0, 50)}...` : text;
      },
    },
    {
      title: $t('page.comment.contentType'),
      field: 'contentType',
      slots: { default: 'contentType' },
      width: 110,
    },
    {
      title: $t('page.comment.authorType'),
      field: 'authorType',
      slots: { default: 'authorType' },
      width: 110,
    },
    {
      title: $t('page.comment.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 100,
    },
    {
      title: $t('page.comment.likeCount'),
      field: 'likeCount',
      width: 90,
      align: 'right',
      headerAlign: 'right',
      sortable: true,
    },
    {
      title: $t('page.comment.dislikeCount'),
      field: 'dislikeCount',
      width: 90,
      align: 'right',
      headerAlign: 'right',
      sortable: true,
    },
    {
      title: $t('page.comment.replyCount'),
      field: 'replyCount',
      width: 90,
      align: 'right',
      headerAlign: 'right',
      sortable: true,
    },
    {
      title: $t('page.comment.ipAddress'),
      field: 'ipAddress',
      minWidth: 130,
    },
    {
      title: $t('page.comment.location'),
      field: 'location',
      minWidth: 120,
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
  connectedComponent: CommentDrawer,
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
      await apiClient.commentService.Delete({ id: row.id });
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
    <Grid :table-title="$t('menu.engagement.comment')">
      <template #status="{ row }">
        <a-tag :color="commentStatusToColor(row.status)">
          {{ commentStatusToName(row.status) }}
        </a-tag>
      </template>
      <template #contentType="{ row }">
        <a-tag :color="commentContentTypeToColor(row.contentType)">
          {{ commentContentTypeToName(row.contentType) }}
        </a-tag>
      </template>
      <template #authorType="{ row }">
        <a-tag :color="commentAuthorTypeToColor(row.authorType)">
          {{ commentAuthorTypeToName(row.authorType) }}
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
              moduleName: $t('page.comment.moduleName'),
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
