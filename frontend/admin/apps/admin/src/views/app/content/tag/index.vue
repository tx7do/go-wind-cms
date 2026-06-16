<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, type VbenFormProps } from '@vben/common-ui';
import { IconifyIcon, LucideFilePenLine, LucideTrash2 } from '@vben/icons';
import { i18n } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  apiClient,
  fetchListTags,
  PaginationQuery,
  type contentservicev1_Tag as Tag,
  tagStatusList,
  tagStatusToColor,
  tagStatusToName,
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
      component: 'Select',
      fieldName: 'status',
      label: $t('page.tag.status'),
      componentProps: {
        options: tagStatusList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
  ],
};

const gridOptions: VxeGridProps<Tag> = {
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
        console.log('query:', formValues);

        return await fetchListTags(
          new PaginationQuery({
            paging: {
              page: page.currentPage,
              pageSize: page.pageSize,
            },
            formValues,
          }),
        );
      },
    },
  },

  columns: [
    {
      title: $t('page.tag.name'),
      field: 'translations.name',
      minWidth: 100,
      fixed: 'left',
      align: 'left',
      slots: { default: 'tagName' },
    },
    {
      title: $t('page.tag.description'),
      field: 'translations.description',
      minWidth: 250,
      align: 'left',
      slots: { default: 'tagDescription' },
    },
    {
      title: $t('page.tag.color'),
      field: 'color',
      slots: { default: 'color' },
      minWidth: 100,
    },
    {
      title: $t('page.tag.icon'),
      field: 'icon',
      minWidth: 80,
      slots: { default: 'icon' },
    },
    {
      title: $t('page.tag.group'),
      field: 'group',
      minWidth: 120,
    },
    {
      title: $t('page.tag.isFeatured'),
      field: 'isFeatured',
      formatter: ({ cellValue }) =>
        cellValue ? $t('ui.button.yes') : $t('ui.button.no'),
      minWidth: 100,
    },
    {
      title: $t('page.tag.postCount'),
      field: 'postCount',
      minWidth: 100,
    },
    {
      title: $t('page.tag.status'),
      field: 'status',
      slots: { default: 'status' },
      minWidth: 100,
    },
    {
      title: $t('ui.table.sortOrder'),
      field: 'sortOrder',
      minWidth: 80,
    },
    {
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      minWidth: 150,
    },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      minWidth: 100,
    },
  ],
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

/* Create */
function handleCreate() {
  console.log('Create', i18n.global.locale.value);
  router.push({
    name: 'CreateTag',
    query: { lang: i18n.global.locale.value },
  });
}

/* Edit */
function handleEdit(row: any) {
  console.log('Edit', row.id, i18n.global.locale.value);
  router.push({
    name: 'EditTag',
    params: { id: String(row.id) },
    query: { lang: i18n.global.locale.value },
  });
}

/* Delete */
async function handleDelete(row: any) {
  console.log('Delete', row);

  try {
    await apiClient.tagService.Delete({ id: row.id });

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

function getTranslation(row: any) {
  const currentLang = i18n.global.locale.value;
  const translation = row.translations?.find(
    (t: any) => t.languageCode === currentLang,
  );
  return translation || row.translations?.[0] || undefined;
}

function getTagName(row: any) {
  const translation = getTranslation(row);
  return translation?.name || row.translations?.[0]?.name || '';
}

function getTagDescription(row: any) {
  const translation = getTranslation(row);
  return translation?.description || row.translations?.[0]?.description || '';
}

function getIcon(row: any) {
  if (!row.icon) return '';
  // 如果已有前缀（如carbon:xxx），直接返回
  if (row.icon.includes(':')) return row.icon;
  // 否则加上carbon:前缀
  return `carbon:${row.icon}`;
}
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('menu.content.tag')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('ui.button.create') }}
        </a-button>
      </template>
      <template #color="{ row }">
        <a-tag v-if="row.color" :color="row.color">
          {{ row.color }}
        </a-tag>
      </template>
      <template #status="{ row }">
        <a-tag :color="tagStatusToColor(row.status)">
          {{ tagStatusToName(row.status) }}
        </a-tag>
      </template>
      <template #tagName="{ row }">
        {{ getTagName(row) }}
      </template>
      <template #tagDescription="{ row }">
        {{ getTagDescription(row) }}
      </template>
      <template #icon="{ row }">
        <IconifyIcon v-if="row.icon" :icon="getIcon(row)" class="size-6" />
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
              moduleName: $t('page.tag.moduleName'),
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
