<script lang="ts" setup>
import type { VxeGridListeners, VxeGridProps } from '#/adapter/vxe-table';

import { h, onMounted, ref, watch } from 'vue';

import { Page, useVbenDrawer, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';
import { isEqual } from '@vben/utils';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { type siteservicev1_Navigation as Navigation } from '#/generated/api/admin/service/v1';
import { $t } from '#/locales';
import {
  navigationLocationList,
  navigationLocationToColor,
  navigationLocationToName,
  useLanguageStore,
  useNavigationStore,
} from '#/stores';

import NavigationDrawer from './navigation-drawer.vue';
import { useNavigationViewStore } from './navigation-view.state';

const navigationStore = useNavigationStore();
const languageStore = useLanguageStore();
const navigationViewStore = useNavigationViewStore();

const languageOptions = ref<{ label: string; value: string }[]>([]);

const formOptions: VbenFormProps = {
  collapsed: true,
  showCollapseButton: true,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.navigation.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'location',
      label: $t('page.navigation.location'),
      componentProps: {
        options: navigationLocationList,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
      },
      rules: 'selectRequired',
    },
    {
      component: 'Select',
      fieldName: 'locale',
      label: $t('page.navigation.locale'),
      componentProps: {
        options: languageOptions,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
      },
    },
  ],
};

const gridOptions: VxeGridProps<Navigation> = {
  toolbarConfig: {
    custom: true,
    export: true,
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
        return await navigationViewStore.fetchNavigationList(
          page.currentPage,
          page.pageSize,
          formValues,
        );
      },
    },
  },

  columns: [
    {
      title: $t('page.navigation.name'),
      field: 'name',
      align: 'left',
      fixed: 'left',
      minWidth: 140,
    },
    {
      title: $t('page.navigation.location'),
      field: 'location',
      minWidth: 140,
      slots: {
        default: 'location',
      },
    },
    {
      title: $t('page.navigation.locale'),
      field: 'locale',
      minWidth: 110,
    },
    {
      title: $t('page.navigation.isActive'),
      field: 'isActive',
      minWidth: 100,
      formatter: ({ cellValue }) =>
        cellValue ? $t('ui.button.yes') : $t('ui.button.no'),
    },
    {
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      minWidth: 160,
    },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      minWidth: 90,
    },
  ],
};

const gridEvents: VxeGridListeners<Navigation> = {
  cellClick: ({ row }) => {
    // console.log(`cell-click: ${row.id}`);
    navigationViewStore.setCurrentNavigationId(
      typeof row.id === 'number' ? row.id : 0,
    );
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions,
  formOptions,
  gridEvents,
});

const [Drawer, drawerApi] = useVbenDrawer({
  connectedComponent: NavigationDrawer,
  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      gridApi.reload();
    }
  },
});

function openDrawer(create: boolean, row?: any) {
  drawerApi.setData({
    create,
    row,
  });
  drawerApi.open();
}

function handleCreate() {
  openDrawer(true);
}

function handleEdit(row: any) {
  openDrawer(false, row);
}

async function handleDelete(row: any) {
  try {
    await navigationStore.deleteNavigation(row.id);
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

watch(
  () => navigationViewStore.needReloadNavigationList,
  (newValues, oldValue) => {
    if (isEqual(newValues, oldValue) || !newValues) {
      return;
    }
    navigationViewStore.needReloadNavigationList = false;
    gridApi.reload();
  },
);

onMounted(async () => {
  try {
    const resp = await languageStore.listLanguage(undefined, {}, undefined, [
      'sortOrder',
    ] as any);
    languageOptions.value =
      resp.items?.map((lang) => ({
        label: lang.nativeName || lang.languageCode || '',
        value: lang.languageCode || '',
      })) || [];
  } catch (error) {
    console.error('Failed to load language list:', error);
  }
});
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('page.navigation.moduleName')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('page.navigation.button.create') }}
        </a-button>
      </template>
      <template #location="{ row }">
        <a-tag :color="navigationLocationToColor(row.location)">
          {{ navigationLocationToName(row.location) }}
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
              moduleName: $t('page.navigation.moduleName'),
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

<style scoped></style>
