<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  apiClient,
  fetchListSections,
  type contentservicev1_Section as SectionType,
  sectionTypeList,
  sectionTypeToColor,
  sectionTypeToName,
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
      fieldName: 'name',
      label: $t('page.section.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'type',
      label: $t('page.section.type'),
      componentProps: {
        options: sectionTypeList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
  ],
};

const gridOptions: VxeGridProps<SectionType> = {
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

  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        return await fetchListSections(
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
      title: $t('page.section.name'),
      field: 'name',
      minWidth: 160,
    },
    {
      title: $t('page.section.type'),
      field: 'type',
      slots: { default: 'type' },
      width: 140,
    },
    {
      title: $t('page.section.pageId'),
      field: 'pageId',
      width: 120,
    },
    {
      title: $t('ui.table.sortOrder'),
      field: 'sortOrder',
      width: 90,
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
  router.push({ name: 'CreateSection' });
}

/* Edit */
function handleEdit(row: any) {
  router.push({
    name: 'EditSection',
    params: { id: Number(row.id || 0) },
  });
}

/* Delete */
async function handleDelete(row: any) {
  try {
    await apiClient.sectionService.Delete({ id: row.id });

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
    <Grid :table-title="$t('menu.content.section')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('ui.button.create') }}
        </a-button>
      </template>
      <template #type="{ row }">
        <a-tag :color="sectionTypeToColor(row.type)">
          {{ sectionTypeToName(row.type) }}
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
              moduleName: $t('page.section.moduleName'),
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
