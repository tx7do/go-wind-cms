<script setup lang="ts">
import { computed, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft, LucideSparkles } from '@vben/icons';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { Editor } from '#/adapter/component/Editor';
import { apiClient, editorTypeOptions, uploadMediaAsset } from '#/api';
import { router } from '#/router';

import { usePageEditViewStore } from './page-edit-view.state';

const pageEditViewStore = usePageEditViewStore();

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
  console.log('handleLanguageChange', newLang);

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
    console.log('Title translation response:', titleResp);
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
    console.log('Content translation response:', contentResp);
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
  console.log('Upload image:', file);

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
  console.log('init');

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
      </div>
    </template>

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
