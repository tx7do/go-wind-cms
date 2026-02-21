<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page, type VbenFormProps } from '@vben/common-ui';

import dayjs from 'dayjs';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { type commentservicev1_Comment as Comment } from '#/generated/api/admin/service/v1';
import { $t } from '#/locales';
import {
  commentAuthorTypeList,
  commentAuthorTypeToColor,
  commentAuthorTypeToName,
  commentContentTypeList,
  commentContentTypeToColor,
  commentContentTypeToName,
  commentStatusList,
  commentStatusToColor,
  commentStatusToName,
  useCommentStore,
} from '#/stores';

const commentStore = useCommentStore();

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
      fieldName: 'authorId',
      label: $t('page.comment.authorId'),
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
      component: 'Input',
      fieldName: 'ipAddress',
      label: $t('page.comment.ipAddress'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
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

        return await commentStore.listComment(
          {
            page: page.currentPage,
            pageSize: page.pageSize,
          },
          {
            authorId: formValues.authorId,
            authorType: formValues.authorType,
            contentType: formValues.contentType,
            status: formValues.status,
            ipAddress: formValues.ipAddress,
            created_at__gte: startTime,
            created_at__lte: endTime,
          },
          null,
          ['-created_at'],
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
      title: $t('page.comment.content'),
      field: 'content',
    },
    {
      title: $t('page.comment.contentType'),
      field: 'contentType',
      slots: { default: 'contentType' },
      width: 100,
    },
    {
      title: $t('page.comment.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 100,
    },
    {
      title: $t('page.comment.authorType'),
      field: 'authorType',
      slots: { default: 'authorType' },
      width: 100,
    },
    { title: $t('page.comment.likeCount'), field: 'likeCount' },
    { title: $t('page.comment.dislikeCount'), field: 'dislikeCount' },
    { title: $t('page.comment.replyCount'), field: 'replyCount' },
    {
      title: $t('page.comment.ipAddress'),
      field: 'ipAddress',
      width: 140,
    },
    {
      title: $t('page.comment.location'),
      field: 'location',
    },
  ],
};

const [Grid] = useVbenVxeGrid({ gridOptions, formOptions });
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
    </Grid>
  </Page>
</template>
