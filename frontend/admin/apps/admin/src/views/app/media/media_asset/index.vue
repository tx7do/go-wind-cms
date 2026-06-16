<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page, type VbenFormProps } from '@vben/common-ui';

import dayjs from 'dayjs';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  fetchListMediaAssets,
  type mediaservicev1_MediaAsset as MediaAsset,
  mediaAssetAssetTypeList,
  mediaAssetAssetTypeToColor,
  mediaAssetAssetTypeToName,
  mediaAssetProcessingStatusList,
  mediaAssetProcessingStatusToColor,
  mediaAssetProcessingStatusToName,
  PaginationQuery,
} from '#/api';
import { $t } from '#/locales';
import { formatBytes } from '#/utils';

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
        console.log('query:', formValues);

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
          console.log(startTime, endTime);
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
    { title: $t('page.mediaAsset.authorId'), field: 'authorId' },
    {
      title: $t('page.mediaAsset.processingStatus'),
      field: 'processingStatus',
      slots: { default: 'processingStatus' },
    },
    {
      title: $t('page.mediaAsset.referenceCount'),
      field: 'referenceCount',
    },
  ],
};

const [Grid] = useVbenVxeGrid({ gridOptions, formOptions });
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('menu.media.mediaAsset')">
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
    </Grid>
  </Page>
</template>
