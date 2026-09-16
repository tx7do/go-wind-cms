<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, h, onActivated, onMounted, ref } from 'vue';

import { Page, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';
import { i18n } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  apiClient,
  editorTypeToColor,
  editorTypeToName,
  enableBoolToColor,
  enableBoolToName,
  fetchFlattenedCategoryOptions,
  fetchListPosts,
  PaginationQuery,
  type contentservicev1_Post as Post,
  postStatusList,
  postStatusToColor,
  postStatusToName,
} from '#/api';
import { $t } from '#/locales';
import { router } from '#/router';

/** 状态 Tab：ALL 之外为 POST_STATUS_* 枚举名 */
const ALL_STATUS = 'ALL';
const activeStatus = ref<string>(ALL_STATUS);
const statusCounts = ref<Record<string, number>>({});
const COUNT_REFRESH_THROTTLE_MS = 3000;
let lastCountRefreshAt = 0;
let countsInFlight = false;
let countsRefreshQueued = false;

/** 分类筛选下拉选项（与编辑页分类选择共用数据源） */
const categoryOptions = ref<{ label: string; value: number }[]>([]);
const categoryOptionsLoading = ref(false);

async function fetchCategoryOptions() {
  if (categoryOptionsLoading.value) return;
  categoryOptionsLoading.value = true;
  try {
    categoryOptions.value = await fetchFlattenedCategoryOptions(
      i18n.global.locale.value,
    );
  } finally {
    categoryOptionsLoading.value = false;
  }
}

const statusTabs = computed(() => [
  { key: ALL_STATUS, label: $t('page.post.statusAll') },
  ...postStatusList.value.map((o) => ({
    key: o.value as string,
    label: o.label,
  })),
]);

const formOptions: VbenFormProps = {
  // 默认展开
  collapsed: false,
  // 控制表单是否显示折叠按钮
  showCollapseButton: false,
  // 按下回车时是否提交表单
  submitOnEnter: true,
  schema: [
    {
      // 主表唯一码字段是 code；proto 的 slug 是 translations 层字段，列表查询会 500
      component: 'Input',
      fieldName: 'code',
      label: $t('page.post.slug'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'categoryIds',
      label: $t('page.post.category'),
      componentProps: {
        options: categoryOptions,
        loading: categoryOptionsLoading,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
        showSearch: true,
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
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
        const query: Record<string, unknown> = { ...(formValues ?? {}) };
        // 状态 Tab 是唯一的状态筛选入口，ALL 表示不过滤
        if (activeStatus.value === ALL_STATUS) {
          delete query.status;
        } else {
          query.status = activeStatus.value;
        }
        // 分类过滤走连接表：字段名必须带 __in 后缀且值为数组，
        // 通用转换器对单值生成 Value，而 post_repo 的连接表过滤只认 Values
        if (query.categoryIds != null) {
          const v = query.categoryIds;
          delete query.categoryIds;
          query.category_ids__in = Array.isArray(v) ? v : [v];
        }
        return await fetchListPosts(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues: query,
            fieldMask:
              'id,status,sort_order,is_featured,author_name,available_languages,created_at,code,editor_type,disallow_comment,in_progress,auto_summary,is_featured,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,',
          }),
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

/**
 * 刷新各状态 Tab 的计数（每次一个 pageSize=1 的计数查询，取 total）
 * @param force 跳过节流强制刷新（如删除后）
 */
async function refreshCounts(force = false) {
  const now = Date.now();
  if (!force && now - lastCountRefreshAt < COUNT_REFRESH_THROTTLE_MS) return;
  if (countsInFlight) {
    countsRefreshQueued = true;
    return;
  }
  countsInFlight = true;
  lastCountRefreshAt = now;
  try {
    const keys = [ALL_STATUS, ...postStatusList.value.map((o) => o.value)];
    const entries = await Promise.all(
      keys.map(async (key) => {
        try {
          const resp = await fetchListPosts(
            new PaginationQuery({
              paging: { page: 1, pageSize: 1 },
              formValues: key === ALL_STATUS ? {} : { status: key },
              fieldMask: 'id',
            }),
          );
          return [key, resp.total] as const;
        } catch {
          return [key, undefined] as const;
        }
      }),
    );
    const next: Record<string, number> = {};
    for (const [key, total] of entries) {
      // uint64 经 JSON 序列化是字符串，需强转数字
      const n = Number(total);
      if (Number.isFinite(n)) {
        next[key] = n;
      }
    }
    statusCounts.value = next;
  } finally {
    countsInFlight = false;
    if (countsRefreshQueued) {
      countsRefreshQueued = false;
      refreshCounts(true);
    }
  }
}

/* 状态 Tab 切换：回到第 1 页重新查询 */
function handleStatusTabChange() {
  gridApi.reload();
}

onMounted(() => {
  refreshCounts();
  fetchCategoryOptions();
});

// keep-alive 缓存页从编辑页返回时刷新计数（节流去重首挂载的重复触发）
onActivated(() => {
  refreshCounts();
});

/* 创建 */
function handleCreate() {
  router.push({
    name: 'CreatePost',
    query: { lang: i18n.global.locale.value },
  });
}

/* 编辑 */
function handleEdit(row: any) {
  router.push({
    name: 'EditPost',
    params: { id: String(row.id) },
    query: { lang: i18n.global.locale.value },
  });
}

/* 删除 */
async function handleDelete(row: any) {
  try {
    await apiClient.postService.Delete({ id: row.id });

    notification.success({
      message: $t('ui.notification.delete_success'),
    });

    await gridApi.reload();
    refreshCounts(true);
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
      <template #toolbar-actions>
        <a-tabs
          v-model:activeKey="activeStatus"
          size="small"
          class="post-status-tabs"
          @change="handleStatusTabChange"
        >
          <a-tab-pane v-for="tabItem in statusTabs" :key="tabItem.key">
            <template #tab>
              <span>
                {{ tabItem.label }}
                <span
                  v-if="statusCounts[tabItem.key] !== undefined"
                  class="status-count"
                  >{{ statusCounts[tabItem.key] }}</span
                >
              </span>
            </template>
          </a-tab-pane>
        </a-tabs>
      </template>
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

<style scoped>
.post-status-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 0;
}

/* 仅作筛选器使用，隐藏空 pane 内容区 */
.post-status-tabs :deep(.ant-tabs-content-holder) {
  display: none;
}

.post-status-tabs :deep(.status-count) {
  margin-left: 4px;
  font-size: 12px;
  opacity: 0.6;
}
</style>
