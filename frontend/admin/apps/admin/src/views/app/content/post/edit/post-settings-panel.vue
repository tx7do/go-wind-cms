<script setup lang="ts">
import type { Dayjs } from 'dayjs';

import { computed, ref } from 'vue';

import { $t } from '@vben/locales';

import dayjs from 'dayjs';

import { usePostEditViewStore } from './post-edit-view.state';
import PostImageField from './post-image-field.vue';

const postEditViewStore = usePostEditViewStore();

// 默认展开基础与媒体组，SEO 组默认收起
const activeGroups = ref<string[]>(['base', 'media']);

const SENTENCE_END_REGEX = /[。！？；.!?;]/g;

/**
 * 与后端 summary.GenerateSummaryByRule 对齐的摘要生成：
 * 剥离HTML/Markdown语法 → 合并空白 → 按句子边界截断（100字）→ 补省略号
 */
function generateSummary(content: string, maxLength = 100): string {
  let text = content.replace(/<[^>]+>/g, ' ');
  text = text.replace(/```[\s\S]*?```/g, ' ');
  text = text.replace(/`[^`]*`/g, ' ');
  text = text.replace(/!\[[^\]]*\]\([^)]*\)/g, ' ');
  text = text.replace(/\[([^\]]*)\]\([^)]*\)/g, '$1');
  text = text.replace(/^[ \t]{0,3}#{1,6}[ \t]+/gm, '');
  text = text.replace(/^[ \t]*>[ \t]?/gm, '');
  text = text.replace(/^[ \t]*([-*_][ \t]*){3,}$/gm, ' ');
  text = text.replace(/^[ \t]*[-+*][ \t]+/gm, '');
  text = text.replace(/^[ \t]*\d+\.[ \t]+/gm, '');
  text = text.replace(/\|/g, ' ');
  text = text.replace(/[*_~]+/g, '');
  text = text.replace(/\s+/g, ' ').trim();

  if (!text) {
    return '暂无摘要';
  }

  const runes = Array.from(text);
  if (runes.length <= maxLength) {
    return text;
  }

  const truncated = runes.slice(0, maxLength).join('');
  const ends = [
    ...truncated.matchAll(new RegExp(SENTENCE_END_REGEX.source, 'g')),
  ];
  if (ends.length === 0) {
    return `${truncated}...`;
  }
  const lastEnd = ends[ends.length - 1]!;
  return `${truncated.slice(0, (lastEnd.index ?? -1) + 1)}...`;
}

function handleGenerateSummary() {
  postEditViewStore.formData.summary = generateSummary(
    postEditViewStore.formData.content || '',
  );
}

const publishTimeValue = computed<Dayjs | null>({
  get: () =>
    postEditViewStore.formData.publishTime
      ? dayjs(postEditViewStore.formData.publishTime)
      : null,
  set: (value) => {
    postEditViewStore.formData.publishTime = value
      ? value.toISOString()
      : undefined;
  },
});
</script>

<template>
  <div class="flex flex-col gap-3">
    <!-- 分类（必选项，常驻侧栏顶部，不随手风琴收起） -->
    <a-form layout="vertical" class="mb-0">
      <a-form-item class="mb-2" required>
        <template #label>
          {{ $t('page.post.category') }}
        </template>
        <a-select
          v-model:value="postEditViewStore.formData.categoryIds"
          mode="multiple"
          :options="postEditViewStore.categoryOptions"
          :loading="postEditViewStore.categoryOptionsLoading"
          :placeholder="$t('ui.placeholder.select')"
          allow-clear
        />
      </a-form-item>
    </a-form>

    <a-collapse
      v-model:active-key="activeGroups"
      :bordered="false"
      expand-icon-position="end"
      class="settings-collapse"
      ghost
    >
    <!-- 基础设置 -->
    <a-collapse-panel
      key="base"
      :header="$t('page.post.settingsGroup.base')"
    >
      <a-form layout="vertical">
        <a-form-item :label="$t('page.post.code')">
          <a-input
            v-model:value="postEditViewStore.formData.code"
            :placeholder="$t('page.post.placeholder.code')"
            allow-clear
          />
        </a-form-item>

        <a-form-item :label="$t('page.post.slug')">
          <a-input
            v-model:value="postEditViewStore.formData.slug"
            :placeholder="$t('page.post.placeholder.slug')"
            allow-clear
          />
        </a-form-item>

        <!-- 注意：antd 4.2.6 中 label prop 会覆盖 #label slot，此处必须省略 label prop -->
        <a-form-item>
          <template #label>
            <div class="flex w-full items-center justify-between">
              <span>{{ $t('page.post.summary') }}</span>
              <a-button
                size="small"
                type="link"
                class="px-1"
                :disabled="!postEditViewStore.formData.content"
                @click="handleGenerateSummary"
              >
                {{ $t('page.post.autoSummary') }}
              </a-button>
            </div>
          </template>
          <a-textarea
            v-model:value="postEditViewStore.formData.summary"
            :rows="3"
            :placeholder="$t('page.post.placeholder.summary')"
          />
        </a-form-item>
      </a-form>
    </a-collapse-panel>

    <!-- 媒体与高级 -->
    <a-collapse-panel
      key="media"
      :header="$t('page.post.settingsGroup.media')"
    >
      <a-form layout="vertical">
        <a-form-item>
          <template #label>
            {{ $t('page.post.thumbnail') }}
            <span class="text-xs text-gray-400">
              {{ $t('page.post.thumbnailSharedHint') }}
            </span>
          </template>
          <PostImageField v-model:value="postEditViewStore.formData.thumbnail" />
        </a-form-item>

        <a-form-item :label="$t('page.post.publishTime')">
          <a-date-picker
            v-model:value="publishTimeValue"
            class="w-full"
            show-time
            :placeholder="$t('ui.placeholder.select')"
          />
        </a-form-item>

        <a-form-item :label="$t('page.post.sortOrder')">
          <a-input-number
            v-model:value="postEditViewStore.formData.sortOrder"
            class="w-full"
            :min="0"
            :precision="0"
          />
        </a-form-item>

        <a-form-item :label="$t('page.post.isFeatured')">
          <a-switch v-model:checked="postEditViewStore.formData.isFeatured" />
        </a-form-item>

        <a-form-item :label="$t('page.post.disallowComment')">
          <a-switch v-model:checked="postEditViewStore.formData.disallowComment" />
        </a-form-item>
      </a-form>
    </a-collapse-panel>

    <!-- SEO 设置 -->
    <a-collapse-panel key="seo" :header="$t('page.post.seo.title')">
      <!-- v-if 护栏：旧本地草稿/异常路径下 seo 可能为空，避免渲染期空引用 -->
      <a-form v-if="postEditViewStore.formData.seo" layout="vertical">
        <a-form-item :label="$t('page.post.seo.seoTitle')">
          <a-input
            v-model:value="postEditViewStore.formData.seo.seoTitle"
            :placeholder="$t('ui.placeholder.input')"
            allow-clear
          />
        </a-form-item>

        <a-form-item :label="$t('page.post.seo.metaKeywords')">
          <a-input
            v-model:value="postEditViewStore.formData.seo.metaKeywords"
            :placeholder="$t('ui.placeholder.input')"
            allow-clear
          />
        </a-form-item>

        <a-form-item :label="$t('page.post.seo.metaDescription')">
          <a-textarea
            v-model:value="postEditViewStore.formData.seo.metaDescription"
            :rows="3"
            :placeholder="$t('ui.placeholder.input')"
          />
        </a-form-item>

        <a-form-item :label="$t('page.post.seo.ogTitle')">
          <a-input
            v-model:value="postEditViewStore.formData.seo.ogTitle"
            :placeholder="$t('ui.placeholder.input')"
            allow-clear
          />
        </a-form-item>

        <a-form-item :label="$t('page.post.seo.ogDescription')">
          <a-textarea
            v-model:value="postEditViewStore.formData.seo.ogDescription"
            :rows="2"
            :placeholder="$t('ui.placeholder.input')"
          />
        </a-form-item>

        <a-form-item :label="$t('page.post.seo.ogImage')">
          <PostImageField v-model:value="postEditViewStore.formData.seo.ogImage" />
        </a-form-item>

        <a-form-item :label="$t('page.post.seo.canonicalUrl')">
          <a-input
            v-model:value="postEditViewStore.formData.seo.canonicalUrl"
            :placeholder="$t('ui.placeholder.input')"
            allow-clear
          />
        </a-form-item>
      </a-form>
    </a-collapse-panel>
    </a-collapse>
  </div>
</template>

<style scoped>
.settings-collapse :deep(.ant-collapse-header) {
  font-weight: 500;
  padding-left: 0 !important;
  padding-right: 0 !important;
}

.settings-collapse :deep(.ant-collapse-content-box) {
  padding-left: 0 !important;
  padding-right: 0 !important;
}
</style>
