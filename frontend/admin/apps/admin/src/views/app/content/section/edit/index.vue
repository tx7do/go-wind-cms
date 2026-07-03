<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft } from '@vben/icons';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { Editor } from '#/adapter/component/Editor';
import { EditorType } from '#/adapter/component/Editor';
import {
  editorTypeOptions,
  fetchListPages,
  PaginationQuery,
  sectionTypeList,
  uploadMediaAsset,
} from '#/api';
import { router } from '#/router';

import { useSectionEditViewStore } from './section-edit-view.state';

const sectionEditViewStore = useSectionEditViewStore();

const route = useRoute();
const { closeCurrentTab } = useTabs();

// 页面下拉选项
const pageOptions = ref<{ label: string; value: number }[]>([]);

// 内容编辑器类型(本地维护,默认 Markdown)
const contentEditorType = ref<EditorType>(EditorType.MARKDOWN);

const initLanguage = computed(() => {
  return (route.query.lang as string) || 'zh-CN';
});

const isCreateMode = computed(() => {
  return route.name === 'CreateSection';
});

const isEditMode = computed(() => {
  return route.name === 'EditSection';
});

const sectionId = computed(() => {
  if (isCreateMode.value) {
    return null;
  }
  const id = route.params.id ?? -1;
  return Number(id);
});

/**
 * Watch route query parameter changes
 */
watch(
  () => route.query.lang,
  async (newLang) => {
    if (newLang && sectionEditViewStore.formData.lang !== newLang) {
      sectionEditViewStore.formData.lang = newLang as string;
    }
  },
);

/**
 * Handle back button click
 */
function goBack() {
  closeCurrentTab();
  router.push('/content/sections');
}

/**
 * Handle language switch
 */
async function handleLanguageChange(newLang: string) {
  await router.replace({
    path: route.path,
    query: { ...route.query, lang: newLang },
  });

  await sectionEditViewStore.switchLanguage(newLang);

  if (sectionEditViewStore.needTranslate) {
    notification.info({
      message: $t('page.page.validation.translationNotExists'),
    });
  }
}

/**
 * Handle save draft
 */
function handleSaveDraft() {
  try {
    sectionEditViewStore.saveSectionDraft();
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
 * Handle save section
 */
async function handleSave() {
  const resp = await sectionEditViewStore.saveSection();
  if (resp) {
    notification.error({
      message: resp,
    });
  } else {
    notification.success({
      message: $t('page.section.validation.saveSuccess'),
    });
    goBack();
  }
}

/**
 * Handle image upload
 */
async function handleUploadImage(file: File): Promise<string> {
  try {
    const resp: any = await uploadMediaAsset({}, file);
    return resp.objectName || resp.data?.objectName || '';
  } catch (error) {
    console.error('Image upload failed:', error);
    return '';
  }
}

/**
 * Load section data (edit mode only)
 */
async function loadSection() {
  if (!isEditMode.value) {
    return;
  }

  try {
    await sectionEditViewStore.fetchSection();

    if (sectionEditViewStore.needTranslate) {
      notification.info({
        message: $t('page.page.validation.translationNotExists'),
      });
    }
  } catch (error) {
    console.error('Failed to load section:', error);
    notification.error({
      message: $t('page.page.validation.translationNotExists'),
    });
    throw error;
  }
}

/**
 * Load page options for the dropdown
 */
async function loadPageOptions() {
  try {
    const resp = await fetchListPages(
      new PaginationQuery({ paging: { page: 1, pageSize: 500 } }),
    );
    pageOptions.value =
      resp.items?.map((page) => ({
        label:
          page.translations?.[0]?.title || page.slug || String(page.id || ''),
        value: Number(page.id || 0),
      })) || [];
  } catch (error) {
    console.error('Failed to load page list:', error);
    pageOptions.value = [];
  }
}

/**
 * Initialize
 */
async function init() {
  try {
    await sectionEditViewStore.fetchLanguageList();
  } catch {
    notification.error({
      message: $t('page.page.validation.loadLanguageFailed'),
    });
  }

  await loadPageOptions();

  if (isCreateMode.value) {
    sectionEditViewStore.initCreateMode(initLanguage.value);
  } else if (isEditMode.value) {
    sectionEditViewStore.initEditMode(sectionId.value || 0, initLanguage.value);
    await loadSection();
  } else {
    console.error('Unknown route name:', route.name);
  }
}

onMounted(() => {
  init();
});
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
        <span class="text-base font-medium">
          {{
            isCreateMode
              ? $t('page.section.title.create')
              : $t('page.section.title.edit')
          }}
        </span>
      </div>
    </template>

    <div class="section-edit-container min-h-0 flex-1 overflow-auto p-4">
      <a-form layout="vertical">
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item :label="$t('page.section.type')">
              <a-select
                v-model:value="sectionEditViewStore.formData.type"
                :placeholder="$t('ui.placeholder.select')"
              >
                <a-select-option
                  v-for="option in sectionTypeList"
                  :key="option.value"
                  :value="option.value"
                >
                  {{ option.label }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item :label="$t('page.section.pageId')">
              <a-select
                v-model:value="sectionEditViewStore.formData.pageId"
                :placeholder="$t('page.section.placeholder.pageId')"
                show-search
                :filter-option="
                  (input: string, option: any) =>
                    option.label.toLowerCase().includes(input.toLowerCase())
                "
              >
                <a-select-option
                  v-for="option in pageOptions"
                  :key="option.value"
                  :value="option.value"
                >
                  {{ option.label }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item :label="$t('ui.table.sortOrder')">
              <a-input-number
                v-model:value="sectionEditViewStore.formData.sortOrder"
                class="w-full"
                :min="0"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item :label="$t('page.section.name')">
          <a-input
            v-model:value="sectionEditViewStore.formData.name"
            :placeholder="$t('page.section.placeholder.name')"
          />
        </a-form-item>

        <a-form-item>
          <div class="mb-2 flex items-center gap-2">
            <a-select
              :value="sectionEditViewStore.formData.lang"
              style="width: 200px"
              @change="handleLanguageChange"
            >
              <a-select-option
                v-for="option in sectionEditViewStore.languageOptions"
                :key="option.value"
                :value="option.value"
              >
                <span>
                  {{ option.label }}
                  <span
                    v-if="option.hasTranslation"
                    class="ml-2 text-green-600"
                  >
                    ✓
                  </span>
                  <span
                    v-else
                    class="ml-2 text-orange-500"
                  >
                    ○
                  </span>
                </span>
              </a-select-option>
            </a-select>
            <a-select
              v-model:value="contentEditorType"
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
          </div>
          <Editor
            v-model="sectionEditViewStore.formData.content"
            :editor-type="contentEditorType"
            :placeholder="$t('page.section.placeholder.content')"
            :upload-image="handleUploadImage"
          />
        </a-form-item>
      </a-form>
    </div>

    <template #footer>
      <div class="flex w-full">
        <a-space class="ml-auto">
          <a-button type="default" @click="handleSaveDraft">
            {{ $t('page.page.button.saveDraft') }}
          </a-button>
          <a-button type="primary" @click="handleSave">
            {{ $t('page.section.button.save') }}
          </a-button>
        </a-space>
      </div>
    </template>
  </Page>
</template>

<style scoped>
.section-edit-container {
  width: 100%;
  height: 100%;
}
</style>
