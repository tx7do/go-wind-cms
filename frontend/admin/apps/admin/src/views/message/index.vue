<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { VbenFormProps } from '@vben/common-ui';

import { Page } from '@vben/common-ui';
import { useUserStore } from '@vben/stores';
import { $t } from '@vben/locales';

import dayjs from 'dayjs';
import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  fetchListUserInbox,
  type internal_messageservicev1_InternalMessageRecipient as InternalMessageRecipient,
  internalMessageRecipientStatusColor,
  internalMessageRecipientStatusLabel,
  PaginationQuery,
  useDeleteNotificationFromInbox,
  useMarkNotificationAsRead,
} from '#/api';

const userStore = useUserStore();
const { mutateAsync: markNotificationAsRead } = useMarkNotificationAsRead();
const { mutateAsync: deleteNotificationFromInbox } =
  useDeleteNotificationFromInbox();

const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: false,
  submitOnEnter: true,
  schema: [
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

const gridOptions: VxeGridProps<InternalMessageRecipient> = {
  height: 'auto',
  stripe: true,

  exportConfig: {},
  pagerConfig: {},
  rowConfig: {
    isHover: true,
  },

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

        return await fetchListUserInbox(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues: {
              created_at__gte: startTime,
              created_at__lte: endTime,
            },
          }),
        );
      },
    },
  },

  columns: [
    { title: $t('page.internalMessage.title'), field: 'title' },
    {
      title: $t('page.internalMessage.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 100,
    },
    {
      title: $t('page.internalMessage.readAt'),
      field: 'readAt',
      formatter: 'formatDateTime',
      width: 140,
    },
    {
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      width: 140,
    },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      minWidth: 160,
    },
  ],
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

async function handleMarkAllRead() {
  const gridData = gridApi.grid?.getTableData();
  if (!gridData?.fullData || gridData.fullData.length === 0) {
    return;
  }

  const ids: number[] = [];
  for (const row of gridData.fullData as InternalMessageRecipient[]) {
    if (row.status !== 'READ' && row.id) {
      ids.push(row.id);
    }
  }

  if (ids.length === 0) {
    return;
  }

  try {
    await markNotificationAsRead({
      userId: userStore.userInfo?.id ?? 0,
      recipientIds: ids,
    });
    notification.success({
      message: $t('ui.notification.update_success'),
    });
    await gridApi.reload();
  } catch {
    notification.error({
      message: $t('ui.notification.update_failed'),
    });
  }
}

async function handleMarkRead(row: InternalMessageRecipient) {
  if (row.status === 'READ' || !row.id) {
    return;
  }

  try {
    await markNotificationAsRead({
      userId: userStore.userInfo?.id ?? 0,
      recipientIds: [row.id],
    });
    notification.success({
      message: $t('ui.notification.update_success'),
    });
    await gridApi.reload();
  } catch {
    notification.error({
      message: $t('ui.notification.update_failed'),
    });
  }
}

async function handleDelete(row: InternalMessageRecipient) {
  if (!row.id) {
    return;
  }

  try {
    await deleteNotificationFromInbox({
      userId: userStore.userInfo?.id ?? 0,
      recipientIds: [row.id],
    });
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
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleMarkAllRead">
          {{ $t('page.internalMessage.button.markAllRead') }}
        </a-button>
      </template>
      <template #status="{ row }">
        <a-tag :color="internalMessageRecipientStatusColor(row.status)">
          {{ internalMessageRecipientStatusLabel(row.status) }}
        </a-tag>
      </template>
      <template #action="{ row }">
        <a-button
          type="link"
          size="small"
          :disabled="row.status === 'READ'"
          @click="handleMarkRead(row)"
        >
          {{ $t('page.internalMessage.button.markRead') }}
        </a-button>
        <a-button type="link" size="small" danger @click="handleDelete(row)">
          {{ $t('page.internalMessage.button.delete') }}
        </a-button>
      </template>
    </Grid>
  </Page>
</template>

<style scoped></style>
