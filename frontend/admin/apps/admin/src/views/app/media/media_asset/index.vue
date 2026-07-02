<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, useVbenDrawer, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';

import { notification, Upload } from 'ant-design-vue';
import dayjs from 'dayjs';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  apiClient,
  fetchListMediaAssets,
  type mediaservicev1_MediaAsset as MediaAsset,
  mediaAssetAssetTypeList,
  mediaAssetAssetTypeToColor,
  mediaAssetAssetTypeToName,
  mediaAssetProcessingStatusList,
  mediaAssetProcessingStatusToColor,
  mediaAssetProcessingStatusToName,
  PaginationQuery,
  uploadMediaAsset,
} from '#/api';
import { $t } from '#/locales';
import { formatBytes } from '#/utils';

import MediaAssetDrawer from './media-asset-drawer.vue';

const formOptions: VbenFormProps = {
  // 默认展开
  collapsed: false,
  // 控制表单是否显示折叠按钮
  showCollapseButton: false,
  // 按下回车时是否提交表单
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'filename',
      label: $t('page.mediaAsset.filename'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'type',
      label: $t('page.mediaAsset.type'),
      componentProps: {
        options: mediaAssetAssetTypeList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'processingStatus',
      label: $t('page.mediaAsset.processingStatus'),
      componentProps: {
        options: mediaAssetProcessingStatusList,
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

const gridOptions: VxeGridProps<MediaAsset> = {
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

        return await fetchListMediaAssets(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues: {
              type: formValues.type,
              filename: formValues.filename,
              processingStatus: formValues.processingStatus,
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
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      width: 140,
    },
    {
      title: $t('page.mediaAsset.filename'),
      field: 'filename',
    },
    {
      title: $t('page.mediaAsset.type'),
      field: 'type',
      slots: { default: 'type' },
    },
    {
      title: $t('page.mediaAsset.size'),
      field: 'size',
      slots: { default: 'size' },
    },
    { title: $t('page.mediaAsset.storagePath'), field: 'storagePath' },
    { title: $t('page.mediaAsset.title'), field: 'title' },
    { title: $t('page.mediaAsset.caption'), field: 'caption' },
    {
      title: $t('page.mediaAsset.processingStatus'),
      field: 'processingStatus',
      slots: { default: 'processingStatus' },
    },
    {
      title: $t('page.mediaAsset.referenceCount'),
      field: 'referenceCount',
    },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      width: 110,
    },
  ],
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

const [Drawer, drawerApi] = useVbenDrawer({
  // 连接抽离的组件
  connectedComponent: MediaAssetDrawer,

  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      // 关闭时，重载表格数据
      gridApi.reload();
    }
  },
});

/* 打开模态窗口 */
function openDrawer(create: boolean, row?: any) {
  drawerApi.setData({
    create,
    row,
  });

  drawerApi.open();
}

/* 编辑 */
function handleEdit(row: any) {
  openDrawer(false, row);
}

/* 删除 */
async function handleDelete(row: any) {
  try {
    await apiClient.mediaAssetService.Delete({ id: row.id });

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

/* 上传媒体资源 */
async function handleUploadMediaAsset(options: any) {
  const { file, onSuccess, onError } = options;

  try {
    await uploadMediaAsset(
      {
        fileDirectory: 'media',
        title: file.name,
      },
      file,
    );

    onSuccess?.({}, file);

    await gridApi.reload();

    notification.success({
      message: $t('ui.notification.upload_success'),
    });
  } catch (error) {
    console.error('上传媒体资源失败', error);

    try {
      onError?.(error, file);
    } catch {}

    notification.error({
      message: $t('ui.notification.upload_failed'),
    });
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('menu.media.mediaAsset')">
      <template #toolbar-tools>
        <Upload
          :multiple="false"
          :custom-request="handleUploadMediaAsset"
          :show-upload-list="false"
        >
          <a-button class="mr-2" type="primary">
            {{ $t('ui.button.upload') }}
          </a-button>
        </Upload>
      </template>
      <template #type="{ row }">
        <a-tag :color="mediaAssetAssetTypeToColor(row.type)">
          {{ mediaAssetAssetTypeToName(row.type) }}
        </a-tag>
      </template>
      <template #processingStatus="{ row }">
        <a-tag :color="mediaAssetProcessingStatusToColor(row.processingStatus)">
          {{ mediaAssetProcessingStatusToName(row.processingStatus) }}
        </a-tag>
      </template>
      <template #size="{ row }">
        {{ formatBytes(row.size) }}
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
              moduleName: $t('page.mediaAsset.moduleName'),
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
