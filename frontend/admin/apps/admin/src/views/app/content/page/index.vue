<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';
import { i18n } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  apiClient,
  editorTypeToColor,
  editorTypeToName,
  fetchListPages,
  pageStatusList,
  pageStatusToColor,
  pageStatusToName,
  type contentservicev1_Page as PageType,
  pageTypeToColor,
  pageTypeToName,
  PaginationQuery,
} from '#/api';
import { $t } from '#/locales';
import { router } from '#/router';

const formOptions: VbenFormProps = {
  // Default expanded
  collapsed: false,
  // Control whether the form displays a collapse button
  showCollapseButton: false,
  // Whether to submit the form when Enter is pressed
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'slug',
      label: $t('page.page.slug'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('page.page.status'),
      componentProps: {
        options: pageStatusList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
  ],
};

const gridOptions: VxeGridProps<PageType> = {
  toolbarConfig: {
    custom: true,
    export: true,
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
  treeConfig: {
    transform: true,
    rowField: 'id',
    parentField: 'parentId',
  },

  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        console.log('query:', formValues);

        return await fetchListPages(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues,
          }),
        );
      },
    },
  },

  columns: [
    {
      title: $t('page.page.slug'),
      field: 'slug',
      treeNode: true,
    },
    {
      title: $t('page.page.editorType'),
      field: 'editorType',
      slots: { default: 'editorType' },
    },
    {
      title: $t('page.page.type'),
      field: 'type',
      slots: { default: 'type' },
    },
    {
      title: $t('page.page.showInNavigation'),
      field: 'showInNavigation',
      formatter: ({ cellValue }) =>
        cellValue ? $t('ui.button.yes') : $t('ui.button.no'),
      width: 120,
    },
    {
      title: $t('page.page.disallowComment'),
      field: 'disallowComment',
      formatter: ({ cellValue }) =>
        cellValue ? $t('ui.button.yes') : $t('ui.button.no'),
      width: 100,
    },
    {
      title: $t('page.page.authorName'),
      field: 'authorName',
    },
    {
      title: $t('page.page.visits'),
      field: 'visits',
      width: 80,
    },
    {
      title: $t('page.page.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 100,
    },
    {
      title: $t('ui.table.sortOrder'),
      field: 'sortOrder',
      width: 80,
    },
    {
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      width: 150,
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

/* Create */
function handleCreate() {
  console.log('Create', i18n.global.locale.value);
  router.push({
    name: 'CreatePage',
    query: { lang: i18n.global.locale.value },
  });
}

/* Edit */
function handleEdit(row: any) {
  console.log('Edit', row, i18n.global.locale.value);
  router.push({
    name: 'EditPage',
    params: { id: Number(row.id || 0) },
    query: { lang: i18n.global.locale.value },
  });
}

/* Delete */
async function handleDelete(row: any) {
  console.log('Delete', row);

  try {
    await apiClient.pageService.Delete({ id: row.id });

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
    <Grid :table-title="$t('menu.content.page')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('ui.button.create') }}
        </a-button>
      </template>
      <template #status="{ row }">
        <a-tag :color="pageStatusToColor(row.status)">
          {{ pageStatusToName(row.status) }}
        </a-tag>
      </template>
      <template #type="{ row }">
        <a-tag :color="pageTypeToColor(row.type)">
          {{ pageTypeToName(row.type) }}
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
              moduleName: $t('page.page.moduleName'),
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

<style scoped></style>
