<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';
import { i18n } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { type contentservicev1_Post as Post } from '#/generated/api/admin/service/v1';
import { $t } from '#/locales';
import { router } from '#/router';
import {
  editorTypeToColor,
  editorTypeToName,
  enableBoolToColor,
  enableBoolToName,
  postStatusList,
  postStatusToColor,
  postStatusToName,
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
        options: postStatusList,
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
          'id,status,sort_order,is_featured,visits,likes,comment_count,author_name,available_languages,created_at,code,editor_type,disallow_comment,in_progress,auto_summary,is_featured,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,translations.thumbnail',
        );
      },
    },
  },

  columns: [
    {
      title: $t('page.post.postTitle'),
      field: 'translations.title',
      align: 'left',
      fixed: 'left',
      minWidth: 400,
      slots: { default: 'postTitle' },
    },
    {
      title: $t('page.post.slug'),
      field: 'code',
      align: 'left',
      minWidth: 200,
    },
    {
      title: $t('page.post.editorType'),
      field: 'editorType',
      slots: { default: 'editorType' },
      minWidth: 140,
    },
    {
      title: $t('page.post.authorName'),
      field: 'authorName',
      align: 'left',
      minWidth: 140,
    },
    { title: $t('page.post.visits'), field: 'visits', minWidth: 80 },
    { title: $t('page.post.likes'), field: 'likes', minWidth: 80 },
    {
      title: $t('page.post.commentCount'),
      field: 'commentCount',
      minWidth: 80,
    },
    {
      title: $t('page.post.status'),
      field: 'status',
      slots: { default: 'status' },
      minWidth: 95,
    },
    { title: $t('ui.table.sortOrder'), field: 'sortOrder', width: 70 },
    {
      title: $t('page.post.disallowComment'),
      field: 'disallowComment',
      slots: { default: 'disallowComment' },
      minWidth: 80,
    },
    {
      title: $t('page.post.isFeatured'),
      field: 'isFeatured',
      slots: { default: 'isFeatured' },
      minWidth: 80,
    },
    {
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      minWidth: 140,
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

/* 创建 */
function handleCreate() {
  console.log('创建', i18n.global.locale.value);
  router.push({
    name: 'CreatePost',
    query: { lang: i18n.global.locale.value },
  });
}

/* 编辑 */
function handleEdit(row: any) {
  console.log('编辑', row, i18n.global.locale.value);
  router.push({
    name: 'EditPost',
    params: { id: row.id },
    query: { lang: i18n.global.locale.value },
  });
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

function getPostTitle(row: any) {
  const currentLang = i18n.global.locale.value;
  const translation = row.translations?.find(
    (t: any) => t.languageCode === currentLang,
  );
  return translation?.title || row.translations?.[0]?.title || '';
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
      <template #postTitle="{ row }">
        <span>{{ getPostTitle(row) }}</span>
      </template>
      <template #disallowComment="{ row }">
        <a-tag :color="enableBoolToColor(row.disallowComment)">
          {{ enableBoolToName(row.disallowComment) }}
        </a-tag>
      </template>
      <template #isFeatured="{ row }">
        <a-tag :color="enableBoolToColor(row.isFeatured)">
          {{ enableBoolToName(row.isFeatured) }}
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
