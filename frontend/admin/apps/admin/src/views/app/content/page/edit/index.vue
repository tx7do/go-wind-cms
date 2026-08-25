<script setup lang="ts">
import { computed, ref, watch, watchEffect } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft, LucideSparkles } from '@vben/icons';
import { $t } from '@vben/locales';

import { notification, Select, Input, Textarea } from 'ant-design-vue';

import { Editor } from '#/adapter/component/Editor';
import {
  apiClient,
  editorTypeOptions,
  pageStatusList,
  pageTypeList,
  sectionTypeList,
  uploadMediaAsset,
  fetchListContentModels,
  fetchListFieldDefinitions,
  type contentservicev1_FieldDefinition as FieldDefinition,
  type contentservicev1_ContentModel as ContentModel,
} from '#/api';
import { router } from '#/router';
import { PaginationQuery } from '#/transport/rest';

import { usePageEditViewStore } from './page-edit-view.state';
import RelationFieldSelect from '../../model/RelationFieldSelect.vue';

const pageEditViewStore = usePageEditViewStore();

// 内容模型绑定 + 动态字段表单
const contentModelOptions = ref<{ label: string; value: number }[]>([]);
const fieldDefinitions = ref<FieldDefinition[]>([]);

watchEffect(async () => {
  try {
    const resp = await fetchListContentModels(
      new PaginationQuery({
        paging: { page: 1, pageSize: 100 },
        orderBy: ['sort_order'],
      }),
    );
    contentModelOptions.value = (resp?.items ?? []).map((m: ContentModel) => ({
      label: m.name || m.code || `#${m.id}`,
      value: m.id as number,
    }));
  } catch {
    contentModelOptions.value = [];
  }
});

watch(
  () => pageEditViewStore.formData.contentModelId,
  async (modelId) => {
    if (!modelId) {
      fieldDefinitions.value = [];
      return;
    }
    try {
      const resp = await fetchListFieldDefinitions({ contentModelId: modelId });
      fieldDefinitions.value = (resp?.items ?? []) as FieldDefinition[];
    } catch {
      fieldDefinitions.value = [];
    }
  },
  { immediate: true },
);

// ==============================
// 嵌套区块管理（页面子部件，随页面整体读写）
// ==============================
const sections = computed(() => pageEditViewStore.formData.sections ?? []);

function addSection() {
  const list = pageEditViewStore.formData.sections ?? [];
  list.push({
    type: 'SECTION_TYPE_RICH_TEXT',
    name: '',
    sortOrder: list.length,
    config: {},
    translations: [],
  });
  pageEditViewStore.formData.sections = list;
}

function removeSection(index: number) {
  const list = pageEditViewStore.formData.sections ?? [];
  list.splice(index, 1);
  pageEditViewStore.formData.sections = list;
}

function moveSection(index: number, delta: number) {
  const list = pageEditViewStore.formData.sections ?? [];
  const target = index + delta;
  if (target < 0 || target >= list.length) {
    return;
  }
  const a = list[index];
  const b = list[target];
  if (!a || !b) {
    return;
  }
  list[index] = b;
  list[target] = a;
  pageEditViewStore.formData.sections = list;
}

function ensureSectionTranslations(sectionIndex: number) {
  const list = pageEditViewStore.formData.sections ?? [];
  return list[sectionIndex]?.translations ?? [];
}

function addSectionTranslation(sectionIndex: number, lang: string) {
  const list = pageEditViewStore.formData.sections ?? [];
  if (!list[sectionIndex]) {
    return;
  }
  const trs = list[sectionIndex].translations ?? [];
  if (trs.some((t) => t.languageCode === lang)) {
    return;
  }
  trs.push({ languageCode: lang, content: {} });
  list[sectionIndex].translations = trs;
  pageEditViewStore.formData.sections = list;
}

function removeSectionTranslation(sectionIndex: number, trIndex: number) {
  const list = pageEditViewStore.formData.sections ?? [];
  const trs = list[sectionIndex]?.translations ?? [];
  trs.splice(trIndex, 1);
  if (list[sectionIndex]) {
    list[sectionIndex].translations = trs;
  }
  pageEditViewStore.formData.sections = list;
}

// 区块可用语言选项（复用页面已加载的语言列表，与主译文语言一致）
const sectionLanguageOptions = computed(() =>
  pageEditViewStore.languageOptions.map((o) => ({
    label: o.label,
    value: o.value,
  })),
);

const route = useRoute();
const { closeCurrentTab } = useTabs();

const initLanguage = computed(() => {
  return (route.query.lang as string) || 'zh-CN';
});

const isCreateMode = computed(() => {
  return route.name === 'CreatePage';
});

const isEditMode = computed(() => {
  return route.name === 'EditPage';
});

const pageId = computed(() => {
  if (isCreateMode.value) {
    return null;
  }
  const id = route.params.id ?? -1;
  return Number(id);
});

/**
 * Watch route query parameter changes
 * Automatically update form language when user directly accesses via URL (e.g., bookmark)
 */
watch(
  () => route.query.lang,
  async (newLang) => {
    if (newLang && pageEditViewStore.formData.lang !== newLang) {
      pageEditViewStore.formData.lang = newLang as string;
    }
  },
);

/**
 * Handle back button click
 */
function goBack() {
  closeCurrentTab();
  router.push('/content/pages');
}

/**
 * Handle language switch
 */
async function handleLanguageChange(newLang: string) {
  // Update URL query parameter
  await router.replace({
    path: route.path,
    query: { ...route.query, lang: newLang },
  });

  // Use store's switchLanguage method to handle language change with draft loading
  await pageEditViewStore.switchLanguage(newLang);

  // Show notification if translation doesn't exist
  if (pageEditViewStore.needTranslate) {
    notification.info({
      message: $t('page.page.validation.translationNotExists'),
    });
  }
}

/**
 * Handle one-click translation
 */
async function handleTranslate() {
  try {
    const titleResp = await apiClient.translatorService.Translate({
      sourceLanguage: 'auto',
      targetLanguage: pageEditViewStore.formData.lang,
      content: pageEditViewStore.formData.title,
    });
    pageEditViewStore.formData.title =
      titleResp.translatedContent || pageEditViewStore.formData.title;
  } catch (error) {
    console.error('Title translation failed:', error);
    notification.error({
      message: $t('page.page.validation.translateTitleFailed'),
    });
    return;
  }

  try {
    const contentResp = await apiClient.translatorService.Translate({
      sourceLanguage: 'auto',
      targetLanguage: pageEditViewStore.formData.lang,
      content: pageEditViewStore.formData.content,
    });
    pageEditViewStore.formData.content =
      contentResp.translatedContent || pageEditViewStore.formData.content;
  } catch (error) {
    console.error('Content translation failed:', error);
    notification.error({
      message: $t('page.page.validation.translateContentFailed'),
    });
  }
}

/**
 * Handle save draft
 */
function handleSaveDraft() {
  try {
    pageEditViewStore.savePageDraft();
    notification.success({
      message: $t('page.page.validation.saveDraftSuccess'),
    });
  } catch (error) {
    console.error('Save draft failed:', error);
    notification.error({
      message: $t('page.page.validation.saveDraftFailed'),
    });
  }
}

/**
 * Handle publish page
 */
async function handlePublish() {
  const resp = await pageEditViewStore.publishPage();
  if (resp) {
    notification.error({
      message: resp,
    });
  } else {
    notification.success({
      message: $t('page.page.validation.publishSuccess'),
    });

    goBack();
  }
}

/**
 * Handle image upload
 */
async function handleUploadImage(file: File): Promise<string> {
  try {
    const resp = await uploadMediaAsset({}, file);
    return resp.objectName || '';
  } catch (error) {
    console.error('Image upload failed:', error);
    return '';
  }
}

/**
 * Load page data (edit mode only)
 */
async function loadPage() {
  if (!isEditMode.value) {
    return;
  }

  try {
    await pageEditViewStore.fetchPage();

    if (pageEditViewStore.needTranslate) {
      notification.info({
        message: $t('page.page.validation.translationNotExists'),
      });
    }
  } catch (error) {
    console.error('Failed to load page:', error);
    notification.error({
      message: $t('page.page.validation.translationNotExists'),
    });
    throw error;
  }
}

/**
 * Initialize page data
 */
async function init() {
  try {
    await pageEditViewStore.fetchLanguageList();
  } catch {
    notification.error({
      message: $t('page.page.validation.loadLanguageFailed'),
    });
  }

  if (isCreateMode.value) {
    pageEditViewStore.initCreateMode(initLanguage.value);
  } else if (isEditMode.value) {
    pageEditViewStore.initEditMode(pageId.value || 0, initLanguage.value);
    await loadPage();
  } else {
    console.error('Unknown route name:', route.name);
  }
}

init();
</script>

<template>
  <Page
    auto-content-height
    content-class="flex h-full min-h-0 flex-col p-0 overflow-hidden"
  >
    <template #title>
      <div class="flex w-full items-center gap-2">
        <a-button type="text" @click="goBack">
          <template #icon>
            <LucideArrowLeft class="text-align:center" />
          </template>
        </a-button>
        <a-input
          v-model:value="pageEditViewStore.formData.title"
          :placeholder="$t('page.page.placeholder.title')"
          size="large"
          class="flex-1"
        />
        <a-input
          v-model:value="pageEditViewStore.formData.slug"
          :placeholder="$t('page.page.placeholder.slug')"
          size="large"
          style="width: 200px"
        />
        <a-select
          :value="pageEditViewStore.formData.lang"
          style="width: 200px"
          @change="handleLanguageChange"
        >
          <a-select-option
            v-for="option in pageEditViewStore.languageOptions"
            :key="option.value"
            :value="option.value"
          >
            <span>
              {{ option.label }}
              <span
                v-if="option.hasTranslation"
                class="ml-2 text-green-600"
                :title="$t('page.page.placeholder.hasTranslation')"
              >
                ✓
              </span>
              <span
                v-else
                class="ml-2 text-orange-500"
                :title="$t('page.page.placeholder.noTranslation')"
              >
                ○
              </span>
            </span>
          </a-select-option>
        </a-select>
        <a-button
          v-show="pageEditViewStore.needTranslate"
          type="primary"
          class="translate-btn"
          @click="handleTranslate"
        >
          <template #icon>
            <LucideSparkles />
          </template>
          {{ $t('page.page.button.oneClickTranslate') }}
        </a-button>
        <a-select
          v-model:value="pageEditViewStore.formData.editorType"
          style="width: 200px"
        >
          <a-select-option
            v-for="option in editorTypeOptions"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </a-select-option>
        </a-select>
        <a-select
          v-model:value="pageEditViewStore.formData.status"
          style="width: 200px"
        >
          <a-select-option
            v-for="option in pageStatusList"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </a-select-option>
        </a-select>
        <a-select
          v-model:value="pageEditViewStore.formData.type"
          style="width: 200px"
        >
          <a-select-option
            v-for="option in pageTypeList"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </a-select-option>
        </a-select>
      </div>
    </template>

    <!-- 内容模型绑定 + 动态字段表单 -->
    <div v-if="contentModelOptions.length > 0 || fieldDefinitions.length > 0" class="px-4 py-2 border-b border-splitLine">
      <div class="mb-2">
        <label class="text-xs text-textSec block mb-1">{{ $t('page.contentModel.moduleName') }}</label>
        <Select
          v-model:value="pageEditViewStore.formData.contentModelId"
          :options="contentModelOptions"
          allow-clear
          size="small"
          style="width: 100%"
        />
      </div>
      <a-form-item
        v-for="field in fieldDefinitions"
        :key="field.id"
        :label="field.label || field.name"
        class="mb-2"
      >
        <Input
          v-if="field.type === 'FIELD_TYPE_TEXT' || field.type === 'FIELD_TYPE_NUMBER' || field.type === 'FIELD_TYPE_IMAGE' || field.type === 'FIELD_TYPE_FILE'"
          v-model:value="(pageEditViewStore.formData.customFields ?? (pageEditViewStore.formData.customFields = {}))[field.name ?? '']"
          size="small"
        />
        <Textarea
          v-else-if="field.type === 'FIELD_TYPE_RICHTEXT'"
          v-model:value="(pageEditViewStore.formData.customFields ?? (pageEditViewStore.formData.customFields = {}))[field.name ?? '']"
          :rows="2"
          size="small"
        />
        <RelationFieldSelect
          v-else-if="field.type === 'FIELD_TYPE_RELATION'"
          :target-entity-type="field.relationConfig?.targetEntityType"
          :model-value="(pageEditViewStore.formData.customFields ?? (pageEditViewStore.formData.customFields = {}))[field.name ?? '']"
          @update:model-value="(pageEditViewStore.formData.customFields ?? (pageEditViewStore.formData.customFields = {}))[field.name ?? ''] = $event"
          :size="'small'"
        />
      </a-form-item>
    </div>

    <!-- 嵌套区块管理（页面子部件，随页面整体读写） -->
    <div class="px-4 py-2 border-b border-splitLine">
      <div class="flex items-center justify-between mb-2">
        <label class="text-xs text-textSec block">{{ $t('page.section.moduleName') }}</label>
        <a-button size="small" type="primary" @click="addSection">{{ $t('page.section.button.create') }}</a-button>
      </div>

      <div v-if="sections.length === 0" class="text-xs text-textDis">
        {{ $t('page.section.placeholder.content') }}
      </div>

      <div
        v-for="(section, sIdx) in sections"
        :key="sIdx"
        class="border border-splitLine rounded mb-2 p-2"
      >
        <div class="flex items-center gap-2 mb-2">
          <a-button size="small" @click="moveSection(sIdx, -1)">↑</a-button>
          <a-button size="small" @click="moveSection(sIdx, 1)">↓</a-button>
          <span class="text-xs text-textSec">#{{ sIdx }}</span>
          <div class="flex-1" />
          <a-button size="small" danger @click="removeSection(sIdx)">×</a-button>
        </div>

        <div class="flex gap-2 mb-2">
          <div class="flex-1">
            <label class="text-xs text-textSec block mb-1">{{ $t('page.section.type') }}</label>
            <Select
              v-model:value="section.type"
              :options="sectionTypeList"
              size="small"
              style="width: 100%"
            />
          </div>
          <div class="flex-1">
            <label class="text-xs text-textSec block mb-1">{{ $t('page.section.name') }}</label>
            <Input v-model:value="section.name" size="small" :placeholder="$t('page.section.placeholder.name')" />
          </div>
          <div style="width: 90px">
            <label class="text-xs text-textSec block mb-1">{{ $t('page.section.sortOrder') }}</label>
            <Input v-model:value="section.sortOrder" size="small" />
          </div>
        </div>

        <div class="mb-2">
          <label class="text-xs text-textSec block mb-1">{{ $t('page.section.config') }}</label>
          <Textarea
            :value="section.config ? JSON.stringify(section.config) : ''"
            @update:value="(v: string) => { try { section.config = v ? JSON.parse(v) : {} } catch { /* ignore malformed */ } }"
            :rows="2"
            size="small"
          />
        </div>

        <div class="border-t border-splitLine pt-2">
          <div class="flex items-center justify-between mb-2">
            <label class="text-xs text-textSec block">{{ $t('page.section.content') }}</label>
            <div class="flex gap-1">
              <Select
                :value="undefined"
                :options="sectionLanguageOptions"
                size="small"
                style="width: 120px"
                :placeholder="$t('page.section.placeholder.content')"
                @update:value="(lang: string) => { if (lang) addSectionTranslation(sIdx, lang); }"
              />
            </div>
          </div>

          <div
            v-for="(tr, trIdx) in ensureSectionTranslations(sIdx)"
            :key="trIdx"
            class="border border-splitLine rounded mb-2 p-2"
          >
            <div class="flex items-center justify-between mb-2">
              <span class="text-xs text-textSec">{{ tr.languageCode }}</span>
              <a-button size="small" danger @click="removeSectionTranslation(sIdx, trIdx)">×</a-button>
            </div>
            <Textarea
              :value="tr.content ? JSON.stringify(tr.content) : ''"
              @update:value="(v: string) => { try { tr.content = v ? JSON.parse(v) : {} } catch { /* ignore malformed */ } }"
              :rows="3"
              size="small"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="page-edit-container min-h-0 flex-1">
      <Editor
        class="h-full"
        height="100%"
        v-model="pageEditViewStore.formData.content"
        :editor-type="pageEditViewStore.formData.editorType"
        :placeholder="$t('page.page.placeholder.content')"
        :upload-image="handleUploadImage"
      />
    </div>

    <template #footer>
      <div class="flex w-full">
        <a-space class="ml-auto">
          <a-button type="default" @click="handleSaveDraft">
            {{ $t('page.page.button.saveDraft') }}
          </a-button>
          <a-button type="primary" danger @click="handlePublish">
            {{ $t('page.page.button.publish') }}
          </a-button>
        </a-space>
      </div>
    </template>
  </Page>
</template>

<style scoped>
.page-edit-container {
  width: 100%;
  height: 100%;
}

.translate-btn {
  min-width: 140px;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.translate-btn:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}
</style>
