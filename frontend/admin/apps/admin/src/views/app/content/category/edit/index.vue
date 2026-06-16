<script setup lang="ts">
import { computed, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft } from '@vben/icons';
import { $t } from '@vben/locales';

import { Col, InputNumber, notification, Row } from 'ant-design-vue';

import { categoryStatusList } from '#/api';
import { router } from '#/router';

import { useCategoryEditViewStore } from './category-edit-view.state';

const categoryEditViewStore = useCategoryEditViewStore();

const route = useRoute();
const { closeCurrentTab } = useTabs();

const initLanguage = computed(() => {
  return (route.query.lang as string) || 'zh-CN';
});

const isCreateMode = computed(() => {
  return route.name === 'CreateCategory';
});

const isEditMode = computed(() => {
  return route.name === 'EditCategory';
});

const categoryId = computed(() => {
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
    if (newLang && categoryEditViewStore.formData.lang !== newLang) {
      categoryEditViewStore.formData.lang = newLang as string;
    }
  },
);

/**
 * Handle back button click
 */
function goBack() {
  closeCurrentTab();
  router.push('/content/categories');
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
  await categoryEditViewStore.switchLanguage(newLang);

  // Show notification if translation doesn't exist
  if (categoryEditViewStore.needTranslate) {
    notification.info({
      message: $t('page.category.validation.translationNotExists'),
    });
  }
}

/**
 * Handle save draft
 */
function handleSaveDraft() {
  try {
    categoryEditViewStore.saveCategoryDraft();
    notification.success({
      message: $t('page.category.validation.saveDraftSuccess'),
    });
  } catch (error) {
    console.error('Save draft failed:', error);
    notification.error({
      message: $t('page.category.validation.saveDraftFailed'),
    });
  }
}

/**
 * Handle publish category
 */
async function handlePublish() {
  const resp = await categoryEditViewStore.publishCategory();
  if (resp) {
    notification.error({
      message: resp,
    });
  } else {
    notification.success({
      message: $t('page.category.validation.publishSuccess'),
    });

    goBack();
  }
}

/**
 * Load category data (edit mode only)
 */
async function loadCategory() {
  if (!isEditMode.value) {
    return;
  }

  try {
    await categoryEditViewStore.fetchCategory();

    if (categoryEditViewStore.needTranslate) {
      notification.info({
        message: $t('page.category.validation.translationNotExists'),
      });
    }
  } catch (error) {
    console.error('Failed to load category:', error);
    notification.error({
      message: $t('page.category.validation.translationNotExists'),
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
    await categoryEditViewStore.fetchLanguageList();
  } catch {
    notification.error({
      message: $t('page.category.validation.loadLanguageFailed'),
    });
  }

  if (isCreateMode.value) {
    categoryEditViewStore.initCreateMode(initLanguage.value);
  } else if (isEditMode.value) {
    categoryEditViewStore.initEditMode(
      categoryId.value || 0,
      initLanguage.value,
    );
    await loadCategory();
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
          v-model:value="categoryEditViewStore.formData.name"
          :placeholder="$t('page.category.placeholder.name')"
          size="large"
          class="flex-1"
        />
        <a-select
          :value="categoryEditViewStore.formData.lang"
          style="width: 200px"
          @change="handleLanguageChange"
        >
          <a-select-option
            v-for="option in categoryEditViewStore.languageOptions"
            :key="option.value"
            :value="option.value"
          >
            <span>
              {{ option.label }}
              <span
                v-if="option.hasTranslation"
                class="ml-2 text-green-600"
                :title="$t('page.category.placeholder.hasTranslation')"
              >
                ✓
              </span>
              <span
                v-else
                class="ml-2 text-orange-500"
                :title="$t('page.category.placeholder.noTranslation')"
              >
                ○
              </span>
            </span>
          </a-select-option>
        </a-select>
      </div>
    </template>

    <div class="category-edit-container min-h-0 flex-1 overflow-auto p-4">
      <a-form layout="vertical">
        <Row :gutter="16">
          <Col :span="12">
            <a-form-item :label="$t('page.category.placeholder.name')" required>
              <a-input
                v-model:value="categoryEditViewStore.formData.name"
                :placeholder="$t('page.category.placeholder.name')"
              />
            </a-form-item>
          </Col>
          <Col :span="12">
            <a-form-item :label="$t('page.category.placeholder.slug')" required>
              <a-input
                v-model:value="categoryEditViewStore.formData.slug"
                :placeholder="$t('page.category.placeholder.slug')"
              />
            </a-form-item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <a-form-item :label="$t('page.category.icon')">
              <a-input
                v-model:value="categoryEditViewStore.formData.icon"
                :placeholder="$t('page.category.placeholder.icon')"
              />
            </a-form-item>
          </Col>
          <Col :span="12">
            <a-form-item :label="$t('page.category.status')">
              <a-select
                v-model:value="categoryEditViewStore.formData.status"
                :placeholder="$t('ui.placeholder.select')"
                :options="categoryStatusList"
              />
            </a-form-item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="8">
            <a-form-item :label="$t('ui.table.sortOrder')">
              <InputNumber
                v-model:value="categoryEditViewStore.formData.sortOrder"
                :min="0"
                class="w-full"
              />
            </a-form-item>
          </Col>
          <Col :span="8">
            <a-form-item :label="$t('page.category.isNav')">
              <a-switch
                v-model:checked="categoryEditViewStore.formData.isNav"
              />
            </a-form-item>
          </Col>
        </Row>

        <a-form-item :label="$t('page.category.placeholder.description')">
          <a-textarea
            v-model:value="categoryEditViewStore.formData.description"
            :placeholder="$t('page.category.placeholder.description')"
            :rows="4"
          />
        </a-form-item>
      </a-form>
    </div>

    <template #footer>
      <div class="flex w-full">
        <a-space class="ml-auto">
          <a-button type="default" @click="handleSaveDraft">
            {{ $t('page.category.button.saveDraft') }}
          </a-button>
          <a-button type="primary" danger @click="handlePublish">
            {{ $t('page.category.button.publish') }}
          </a-button>
        </a-space>
      </div>
    </template>
  </Page>
</template>

<style scoped>
.category-edit-container {
  width: 100%;
  height: 100%;
}
</style>
