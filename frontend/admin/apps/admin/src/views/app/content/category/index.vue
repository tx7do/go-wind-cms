<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, type VbenFormProps } from '@vben/common-ui';
import { IconifyIcon, LucideFilePenLine, LucideTrash2 } from '@vben/icons';
import { i18n } from '@vben/locales';

import { Image, notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { type contentservicev1_Category as Category } from '#/generated/api/admin/service/v1';
import { $t } from '#/locales';
import { router } from '#/router';
import {
  categoryStatusList,
  categoryStatusToColor,
  categoryStatusToName,
  useCategoryStore,
} from '#/stores';

const categoryStore = useCategoryStore();

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
      label: $t('page.category.status'),
      componentProps: {
        options: categoryStatusList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
  ],
};

const gridOptions: VxeGridProps<Category> = {
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
  treeConfig: {
    rowField: 'id',
    childrenField: 'children',
    expandAll: true,
  },

  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        console.log('query:', formValues);

        return await categoryStore.listCategory(
          {
            page: page.currentPage,
            pageSize: page.pageSize,
          },
          formValues,
          'id,status,sort_order,is_nav,icon,code,post_count,direct_post_count,available_languages,parent_id,children,created_by,created_at,translations,translations.id,translations.language_code,translations.name,translations.slug,translations.description,translations.thumbnail,translations.cover_image',
        );
      },
    },
  },

  columns: [
    {
      title: $t('page.category.name'),
      field: 'translations',
      treeNode: true,
      fixed: 'left',
      slots: { default: 'categoryName' },
      minWidth: 150,
    },
    {
      title: $t('page.category.icon'),
      field: 'icon',
      slots: { default: 'icon' },
      minWidth: 80,
    },
    {
      title: $t('page.category.thumbnail'),
      field: 'thumbnail',
      slots: { default: 'thumbnail' },
      minWidth: 80,
    },
    {
      title: $t('page.category.coverImage'),
      field: 'coverImage',
      slots: { default: 'coverImage' },
      minWidth: 80,
    },
    {
      title: $t('page.category.description'),
      field: 'description',
      slots: { default: 'description' },
      minWidth: 200,
    },
    {
      title: $t('page.category.isNav'),
      field: 'isNav',
      minWidth: 100,
      formatter: ({ cellValue }) =>
        cellValue ? $t('ui.button.yes') : $t('ui.button.no'),
    },
    {
      title: $t('page.category.postCount'),
      field: 'postCount',
      minWidth: 100,
    },
    {
      title: $t('page.category.directPostCount'),
      field: 'directPostCount',
      minWidth: 120,
    },
    {
      title: $t('page.category.status'),
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
    name: 'CreateCategory',
    query: { lang: i18n.global.locale.value },
  });
}

/* Edit */
function handleEdit(row: any) {
  console.log('Edit', row, i18n.global.locale.value);
  router.push({
    name: 'EditCategory',
    params: { id: row.id },
    query: { lang: i18n.global.locale.value },
  });
}

/* Delete */
async function handleDelete(row: any) {
  console.log('Delete', row);

  try {
    await categoryStore.deleteCategory(row.id);

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

function getCategoryTranslation(row: any) {
  const currentLang = i18n.global.locale.value;
  const translation = row.translations?.find(
    (t: any) => t.languageCode === currentLang,
  );
  return translation || row.translations?.[0] || undefined;
}

function getCategoryName(row: any) {
  const translation = getCategoryTranslation(row);
  return translation?.name || row.translations?.[0]?.name || '';
}

function getCategoryDescription(row: any) {
  const translation = getCategoryTranslation(row);
  return translation?.description || row.translations?.[0]?.description || '';
}

function getCategoryThumbnail(row: any) {
  const translation = getCategoryTranslation(row);
  return translation?.thumbnail || row.translations?.[0]?.thumbnail || '';
}

function getCategoryCoverImage(row: any) {
  const translation = getCategoryTranslation(row);
  return translation?.coverImage || row.translations?.[0]?.coverImage || '';
}

const expandAll = () => {
  gridApi.grid?.setAllTreeExpand(true);
};

const collapseAll = () => {
  gridApi.grid?.setAllTreeExpand(false);
};
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('menu.content.category')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('ui.button.create') }}
        </a-button>
        <a-button class="mr-2" @click="expandAll">
          {{ $t('ui.tree.expand_all') }}
        </a-button>
        <a-button class="mr-2" @click="collapseAll">
          {{ $t('ui.tree.collapse_all') }}
        </a-button>
      </template>
      <template #status="{ row }">
        <a-tag :color="categoryStatusToColor(row.status)">
          {{ categoryStatusToName(row.status) }}
        </a-tag>
      </template>
      <template #categoryName="{ row }">
        {{ getCategoryName(row) }}
      </template>
      <template #description="{ row }">
        {{ getCategoryDescription(row) }}
      </template>
      <template #icon="{ row }">
        <IconifyIcon v-if="row.icon" :icon="row.icon" class="size-6" />
      </template>
      <template #thumbnail="{ row }">
        <Image
          v-if="getCategoryThumbnail(row)"
          :src="getCategoryThumbnail(row)"
          width="50"
        />
      </template>
      <template #coverImage="{ row }">
        <Image
          v-if="getCategoryCoverImage(row)"
          :src="getCategoryCoverImage(row)"
          width="50"
        />
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
              moduleName: $t('page.category.moduleName'),
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
