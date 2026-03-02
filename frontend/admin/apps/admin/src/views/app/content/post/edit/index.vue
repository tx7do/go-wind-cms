<script setup lang="ts">
import { computed, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page, useVbenModal } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft, LucideSparkles } from '@vben/icons';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { Editor } from '#/adapter/component/Editor';
import { router } from '#/router';
import {
  editorTypeOptions,
  useFileTransferStore,
  useTranslatorStore,
} from '#/stores';

import { usePostEditViewStore } from './post-edit-view.state';
import PublishPostModal from './publish-post-modal.vue';

const fileTransferStore = useFileTransferStore();
const postEditViewStore = usePostEditViewStore();
const translatorStore = useTranslatorStore();

const route = useRoute();
const { closeCurrentTab } = useTabs();

const initLanguage = computed(() => {
  return (route.query.lang as string) || 'zh-CN';
});

const isCreateMode = computed(() => {
  return route.name === 'CreatePost';
});

const isEditMode = computed(() => {
  return route.name === 'EditPost';
});

const postId = computed(() => {
  if (isCreateMode.value) {
    return null;
  }
  const id = route.params.id ?? -1;
  return Number(id);
});

const [Modal] = useVbenModal({
  // 连接抽离的组件
  connectedComponent: PublishPostModal,
});

/**
 * 监听路由查询参数变化
 * 当用户通过URL直接访问时（例如打开书签），自动更新表单语言
 */
watch(
  () => route.query.lang,
  async (newLang) => {
    if (newLang && postEditViewStore.formData.lang !== newLang) {
      postEditViewStore.formData.lang = newLang as string;
    }
  },
);

/**
 * 处理返回按钮点击
 */
function goBack() {
  closeCurrentTab();
  router.push('/content/posts');
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
  await postEditViewStore.switchLanguage(newLang);

  // Show notification if translation doesn't exist
  if (postEditViewStore.needTranslate) {
    notification.info({
      message: $t('page.post.validation.translationNotExists'),
    });
  }
}

/**
 * 处理一键翻译
 */
async function handleTranslate() {
  try {
    const titleResp = await translatorStore.translate(
      postEditViewStore.formData.lang,
      postEditViewStore.formData.title,
    );
    console.log('Title translation response:', titleResp);
    postEditViewStore.formData.title =
      titleResp.translatedContent || postEditViewStore.formData.title;
  } catch (error) {
    console.error('Title translation failed:', error);
    notification.error({
      message: $t('page.post.validation.translateTitleFailed'),
    });
    return;
  }

  try {
    const contentResp = await translatorStore.translate(
      postEditViewStore.formData.lang,
      postEditViewStore.formData.content,
    );
    console.log('Content translation response:', contentResp);
    postEditViewStore.formData.content =
      contentResp.translatedContent || postEditViewStore.formData.content;
  } catch (error) {
    console.error('Content translation failed:', error);
    notification.error({
      message: $t('page.post.validation.translateContentFailed'),
    });
  }
}

/**
 * 处理保存草稿
 */
function handleSaveDraft() {
  try {
    postEditViewStore.savePostDraft();
    notification.success({
      message: $t('page.post.validation.saveDraftSuccess'),
    });
  } catch (error) {
    console.error('Save draft failed:', error);
    notification.error({
      message: $t('page.post.validation.saveDraftFailed'),
    });
  }
}

/**
 * 处理发布文章
 */
async function handlePublish() {
  const resp = await postEditViewStore.publishPost();
  if (!resp || resp === '') {
    notification.success({
      message: $t('page.post.validation.publishSuccess'),
    });

    goBack();
  } else {
    notification.error({
      message: resp,
    });
  }
}

/**
 * 处理图片上传
 */
async function handleUploadImage(file: File): Promise<string> {
  console.log('Upload image:', file);

  try {
    const resp = await fileTransferStore.uploadMediaAsset({}, file);
    return resp.objectName || '';
  } catch (error) {
    console.error('Image upload failed:', error);
    return '';
  }
}

/**
 * 加载文章数据（仅编辑模式）
 */
async function loadPost() {
  if (!isEditMode.value) {
    return;
  }

  try {
    await postEditViewStore.fetchPost();

    if (postEditViewStore.needTranslate) {
      notification.info({
        message: $t('page.post.validation.translationNotExists'),
      });
    }
  } catch (error) {
    console.error('Failed to load post:', error);
    notification.error({
      message: $t('page.post.validation.translationNotExists'),
    });
    throw error;
  }
}

/**
 * 初始化页面数据
 */
async function init() {
  console.log('init');

  try {
    await postEditViewStore.fetchLanguageList();
  } catch {
    notification.error({
      message: $t('page.post.validation.loadLanguageFailed'),
    });
  }

  if (isCreateMode.value) {
    postEditViewStore.initCreateMode(initLanguage.value);
  } else if (isEditMode.value) {
    postEditViewStore.initEditMode(postId.value || 0, initLanguage.value);
    await loadPost();
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
          v-model:value="postEditViewStore.formData.title"
          :placeholder="$t('page.post.placeholder.title')"
          size="large"
          class="flex-1"
        />
        <a-select
          :value="postEditViewStore.formData.lang"
          style="width: 200px"
          @change="handleLanguageChange"
        >
          <a-select-option
            v-for="option in postEditViewStore.languageOptions"
            :key="option.value"
            :value="option.value"
          >
            <span>
              {{ option.label }}
              <span
                v-if="option.hasTranslation"
                class="ml-2 text-green-600"
                :title="$t('page.post.placeholder.hasTranslation')"
              >
                ✓
              </span>
              <span
                v-else
                class="ml-2 text-orange-500"
                :title="$t('page.post.placeholder.noTranslation')"
              >
                ○
              </span>
            </span>
          </a-select-option>
        </a-select>
        <a-button
          v-show="postEditViewStore.needTranslate"
          type="primary"
          class="translate-btn"
          @click="handleTranslate"
        >
          <template #icon>
            <LucideSparkles />
          </template>
          {{ $t('page.post.button.oneClickTranslate') }}
        </a-button>
        <a-select
          v-model:value="postEditViewStore.formData.editorType"
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

    <div class="post-edit-container min-h-0 flex-1">
      <Editor
        class="h-full"
        height="100%"
        v-model="postEditViewStore.formData.content"
        :editor-type="postEditViewStore.formData.editorType"
        :placeholder="$t('page.post.placeholder.content')"
        :upload-image="handleUploadImage"
      />
    </div>

    <template #footer>
      <div class="flex w-full">
        <a-space class="ml-auto">
          <a-button type="default" @click="handleSaveDraft">
            {{ $t('page.post.button.saveDraft') }}
          </a-button>
          <a-button type="primary" danger @click="handlePublish">
            {{ $t('page.post.button.publish') }}
          </a-button>
        </a-space>
      </div>
    </template>
    <Modal />
  </Page>
</template>

<style scoped>
.post-edit-container {
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
