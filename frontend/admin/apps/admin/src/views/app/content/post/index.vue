<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { type contentservicev1_Post as Post } from '#/generated/api/admin/service/v1';
import { $t } from '#/locales';
import { router } from '#/router';
import {
  editorTypeToColor,
  editorTypeToName,
  postStatusToColor,
  postStatusToName,
  statusList,
  usePostStore,
} from '#/stores';

const postStore = usePostStore();

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
      fieldName: 'slug',
      label: $t('page.post.slug'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('page.post.status'),
      componentProps: {
        options: statusList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
  ],
};

const gridOptions: VxeGridProps<Post> = {
  toolbarConfig: {
    custom: true,
    export: true,
    // import: true,
    refresh: true,
    zoom: true,
  },
  exportConfig: {},
  pagerConfig: {},
  rowConfig: {
    isHover: true,
  },
  height: 'auto',
  stripe: true,

  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        console.log('query:', formValues);

        return await postStore.listPost(
          {
            page: page.currentPage,
            pageSize: page.pageSize,
          },
          formValues,
        );
      },
    },
  },

  columns: [
    { title: $t('page.post.slug'), field: 'slug' },
    {
      title: $t('page.post.editorType'),
      field: 'editorType',
      slots: { default: 'editorType' },
    },
    {
      title: $t('page.post.disallowComment'),
      field: 'disallowComment',
    },
    { title: $t('page.post.inProgress'), field: 'inProgress' },
    { title: $t('page.post.autoSummary'), field: 'autoSummary' },
    { title: $t('page.post.isFeatured'), field: 'isFeatured' },
    { title: $t('page.post.visits'), field: 'visits', width: 80 },
    { title: $t('page.post.authorName'), field: 'authorName' },
    { title: $t('page.post.likes'), field: 'likes', width: 80 },
    { title: $t('page.post.commentCount'), field: 'commentCount', width: 80 },
    {
      title: $t('page.post.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 95,
    },
    { title: $t('ui.table.sortOrder'), field: 'sortOrder', width: 70 },
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
      width: 90,
    },
  ],
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

/* 创建 */
function handleCreate() {
  console.log('创建');
  router.push(`/content/posts/create`);
}

/* 编辑 */
function handleEdit(row: any) {
  console.log('编辑', row);
  router.push(`/content/posts/edit/${row.id}`);
}

/* 删除 */
async function handleDelete(row: any) {
  console.log('删除', row);

  try {
    await postStore.deletePost(row.id);

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
    <Grid :table-title="$t('menu.content.post')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('page.post.button.create') }}
        </a-button>
      </template>
      <template #status="{ row }">
        <a-tag :color="postStatusToColor(row.status)">
          {{ postStatusToName(row.status) }}
        </a-tag>
      </template>
      <template #editorType="{ row }">
        <a-tag :color="editorTypeToColor(row.editorType)">
          {{ editorTypeToName(row.editorType) }}
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
              moduleName: $t('page.post.moduleName'),
            })
          "
          @confirm="handleDelete(row)"
        >
          <a-button danger type="link" :icon="h(LucideTrash2)" />
        </a-popconfirm>
      </template>
    </Grid>
  </Page>
</template>
