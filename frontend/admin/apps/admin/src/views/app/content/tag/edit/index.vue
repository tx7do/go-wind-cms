<script setup lang="ts">
import { computed, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft } from '@vben/icons';
import { $t } from '@vben/locales';

import { Col, InputNumber, notification, Row } from 'ant-design-vue';

import { tagStatusList } from '#/api';
import { router } from '#/router';

import { useTagEditViewStore } from './tag-edit-view.state';

const tagEditViewStore = useTagEditViewStore();

const route = useRoute();
const { closeCurrentTab } = useTabs();

const initLanguage = computed(() => {
  return (route.query.lang as string) || 'zh-CN';
});

const isCreateMode = computed(() => {
  return route.name === 'CreateTag';
});

const isEditMode = computed(() => {
  return route.name === 'EditTag';
});

const tagId = computed(() => {
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
    if (newLang && tagEditViewStore.formData.lang !== newLang) {
      tagEditViewStore.formData.lang = newLang as string;
    }
  },
);

/**
 * Handle back button click
 */
function goBack() {
  closeCurrentTab();
  router.push('/content/tags');
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
  await tagEditViewStore.switchLanguage(newLang);

  // Show notification if translation doesn't exist
  if (tagEditViewStore.needTranslate) {
    notification.info({
      message: $t('page.tag.validation.translationNotExists'),
    });
  }
}

/**
 * Handle save draft
 */
function handleSaveDraft() {
  try {
    tagEditViewStore.saveTagDraft();
    notification.success({
      message: $t('page.tag.validation.saveDraftSuccess'),
    });
  } catch (error) {
    console.error('Save draft failed:', error);
    notification.error({
      message: $t('page.tag.validation.saveDraftFailed'),
    });
  }
}

/**
 * Handle publish tag
 */
async function handlePublish() {
  const resp = await tagEditViewStore.publishTag();
  if (resp) {
    notification.error({
      message: resp,
    });
  } else {
    notification.success({
      message: $t('page.tag.validation.publishSuccess'),
    });

    goBack();
  }
}

/**
 * Load tag data (edit mode only)
 */
async function loadTag() {
  if (!isEditMode.value) {
    return;
  }

  try {
    await tagEditViewStore.fetchTag();

    if (tagEditViewStore.needTranslate) {
      notification.info({
        message: $t('page.tag.validation.translationNotExists'),
      });
    }
  } catch (error) {
    console.error('Failed to load tag:', error);
    notification.error({
      message: $t('page.tag.validation.loadFailed'),
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
    await tagEditViewStore.fetchLanguageList();
  } catch {
    notification.error({
      message: $t('page.tag.validation.loadLanguageFailed'),
    });
  }

  if (isCreateMode.value) {
    tagEditViewStore.initCreateMode(initLanguage.value);
  } else if (isEditMode.value) {
    tagEditViewStore.initEditMode(tagId.value || 0, initLanguage.value);
    await loadTag();
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
          v-model:value="tagEditViewStore.formData.name"
          :placeholder="$t('page.tag.placeholder.name')"
          size="large"
          class="flex-1"
        />
        <a-select
          :value="tagEditViewStore.formData.lang"
          style="width: 200px"
          @change="handleLanguageChange"
        >
          <a-select-option
            v-for="option in tagEditViewStore.languageOptions"
            :key="option.value"
            :value="option.value"
          >
            <span>
              {{ option.label }}
              <span
                v-if="option.hasTranslation"
                class="ml-2 text-green-600"
                :title="$t('page.tag.placeholder.hasTranslation')"
              >
                ✓
              </span>
              <span
                v-else
                class="ml-2 text-orange-500"
                :title="$t('page.tag.placeholder.noTranslation')"
              >
                ○
              </span>
            </span>
          </a-select-option>
        </a-select>
      </div>
    </template>

    <div class="tag-edit-container min-h-0 flex-1 overflow-auto p-4">
      <a-form layout="vertical">
        <Row :gutter="16">
          <Col :span="12">
            <a-form-item :label="$t('page.tag.placeholder.name')" required>
              <a-input
                v-model:value="tagEditViewStore.formData.name"
                :placeholder="$t('page.tag.placeholder.name')"
              />
            </a-form-item>
          </Col>
          <Col :span="12">
            <a-form-item :label="$t('page.tag.placeholder.slug')" required>
              <a-input
                v-model:value="tagEditViewStore.formData.slug"
                :placeholder="$t('page.tag.placeholder.slug')"
              />
            </a-form-item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <a-form-item :label="$t('page.tag.color')">
              <a-input
                v-model:value="tagEditViewStore.formData.color"
                :placeholder="$t('page.tag.placeholder.color')"
              />
            </a-form-item>
          </Col>
          <Col :span="12">
            <a-form-item :label="$t('page.tag.icon')">
              <a-input
                v-model:value="tagEditViewStore.formData.icon"
                :placeholder="$t('page.tag.placeholder.icon')"
              />
            </a-form-item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <a-form-item :label="$t('page.tag.group')">
              <a-input
                v-model:value="tagEditViewStore.formData.group"
                :placeholder="$t('page.tag.placeholder.group')"
              />
            </a-form-item>
          </Col>
          <Col :span="12">
            <a-form-item :label="$t('page.tag.status')">
              <a-select
                v-model:value="tagEditViewStore.formData.status"
                :placeholder="$t('ui.placeholder.select')"
                :options="tagStatusList"
              />
            </a-form-item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="8">
            <a-form-item :label="$t('ui.table.sortOrder')">
              <InputNumber
                v-model:value="tagEditViewStore.formData.sortOrder"
                :min="0"
                class="w-full"
              />
            </a-form-item>
          </Col>
          <Col :span="8">
            <a-form-item :label="$t('page.tag.isFeatured')">
              <a-switch
                v-model:checked="tagEditViewStore.formData.isFeatured"
              />
            </a-form-item>
          </Col>
        </Row>

        <a-form-item :label="$t('page.tag.placeholder.description')">
          <a-textarea
            v-model:value="tagEditViewStore.formData.description"
            :placeholder="$t('page.tag.placeholder.description')"
            :rows="4"
          />
        </a-form-item>
      </a-form>
    </div>

    <template #footer>
      <div class="flex w-full">
        <a-space class="ml-auto">
          <a-button type="default" @click="handleSaveDraft">
            {{ $t('page.tag.button.saveDraft') }}
          </a-button>
          <a-button type="primary" danger @click="handlePublish">
            {{ $t('page.tag.button.publish') }}
          </a-button>
        </a-space>
      </div>
    </template>
  </Page>
</template>

<style scoped>
.tag-edit-container {
  width: 100%;
  height: 100%;
}
</style>
