<script setup lang="ts">
import { computed, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page, useVbenModal } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft } from '@vben/icons';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { Editor } from '#/adapter/component/Editor';
import { router } from '#/router';
import {
  convertToUIEditorType,
  editorTypeOptions,
  useFileTransferStore,
  usePostStore,
} from '#/stores';

import { usePostEditViewStore } from './post-edit-view.state';
import PublishPostModal from './publish-post-modal.vue';

const postStore = usePostStore();
const fileTransferStore = useFileTransferStore();
const postEditViewStore = usePostEditViewStore();

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

async function handleLanguageChange(newLang: string) {
  postEditViewStore.formData.lang = newLang;

  // 更新URL中的 lang 查询参数
  await router.replace({
    path: route.path,
    query: { ...route.query, lang: newLang },
  });

  // 如果是编辑模式，需要重新加载该语言版本的文章
  if (isEditMode.value) {
    try {
      await loadPost();
    } catch {
      notification.error({
        message: $t('page.post.validation.loadFailed'),
      });
    }
  }
}

/**
 * 监听路由查询参数变化
 * 当用户通过URL直接访问时（例如打开书签），自动更新表单语言
 */
watch(
  () => route.query.lang,
  async (newLang) => {
    if (newLang && postEditViewStore.formData.lang !== newLang) {
      postEditViewStore.formData.lang = newLang as string;
      // 如果是编辑模式，重新加载该语言版本的文章
      if (isEditMode.value) {
        try {
          await loadPost();
        } catch {
          notification.error({
            message: $t('page.post.validation.loadFailed'),
          });
        }
      }
    }
  },
);

function goBack() {
  closeCurrentTab();
  router.push('/content/posts');
}

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

async function loadPost() {
  if (!isEditMode.value) {
    return;
  }

  try {
    const item = await postStore.getPost(postId.value || 0);
    if (!item) {
      throw new Error('Post not found');
    }

    if (!item.translations || item.translations.length === 0) {
      throw new Error('No translations found for post');
    }

    let langItem = item.translations?.find(
      (t) => t.languageCode === postEditViewStore.formData.lang,
    );
    if (!langItem) {
      langItem = item.translations?.[0];
      postEditViewStore.formData.lang =
        langItem?.languageCode || postEditViewStore.formData.lang;
    }

    if (!langItem) {
      throw new Error('No translations found for post');
    }

    postEditViewStore.formData.id = item.id;
    postEditViewStore.formData.title = langItem.title || '';
    postEditViewStore.formData.content = langItem.content || '';
    postEditViewStore.formData.lang =
      postEditViewStore.formData.lang || 'zh-CN';
    postEditViewStore.formData.editorType = convertToUIEditorType(
      item.editorType,
    );
  } catch (error) {
    console.error('Failed to load post:', error);
    throw error;
  }
}

async function init() {
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
    try {
      await postEditViewStore.fetchPost();
    } catch {
      notification.error({
        message: $t('page.post.validation.loadFailed'),
      });
    }
  } else {
    console.error('Unknown route name:', route.name);
  }
}

onMounted(async () => {
  await init();
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
            {{ option.label }}
          </a-select-option>
        </a-select>
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
</style>
