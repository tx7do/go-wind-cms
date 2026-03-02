<script setup lang="ts">
import type { PublishProps } from './types';

import { computed, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page, useVbenModal } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft } from '@vben/icons';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { Editor, EditorType } from '#/adapter/component/Editor';
import { router } from '#/router';
import {
  convertToEditorType,
  convertToUIEditorType,
  editorTypeOptions,
  useFileTransferStore,
  useLanguageStore,
  usePostStore,
} from '#/stores';

import PublishPostModal from './publish-post-modal.vue';

const postStore = usePostStore();
const languageStore = useLanguageStore();
const fileTransferStore = useFileTransferStore();

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

// 表单数据
const formData = ref<PublishProps>({
  title: '',
  content: '',
  lang: initLanguage.value,
  editorType: EditorType.MARKDOWN,
});

const languageOptions = ref<{ label: string; value: string }[]>([]);

async function reloadLanguages() {
  try {
    const resp = await languageStore.listLanguage(undefined, {}, undefined, [
      'id',
    ]);
    languageOptions.value =
      resp.items?.map((lang) => ({
        label: lang.nativeName || '',
        value: lang.languageCode || '',
      })) || [];
  } catch (error) {
    console.error('Failed to load languages:', error);
  }
}

const [Modal, modalApi] = useVbenModal({
  // 连接抽离的组件
  connectedComponent: PublishPostModal,
});

/**
 * 当用户切换语言时，同时更新URL查询参数
 */
function handleLanguageChange(newLang: string) {
  formData.value.lang = newLang;

  // 更新URL中的 lang 查询参数
  router.replace({
    path: route.path,
    query: { ...route.query, lang: newLang },
  });

  // 如果是编辑模式，需要重新加载该语言版本的文章
  if (isEditMode.value) {
    loadPost();
  }
}

/**
 * 监听路由查询参数变化
 * 当用户通过URL直接访问时（例如打开书签），自动更新表单语言
 */
watch(
  () => route.query.lang,
  (newLang) => {
    if (newLang && formData.value.lang !== newLang) {
      formData.value.lang = newLang as string;
      // 如果是编辑模式，重新加载该语言版本的文章
      if (isEditMode.value) {
        loadPost();
      }
    }
  },
);

/* 打开模态窗口 */
function openModal() {
  modalApi.setData(formData.value);
  modalApi.open();
}

function goBack() {
  closeCurrentTab();
  router.push('/content/posts');
}

function handleSaveDraft() {
  console.log('Save post:', formData.value);
  // TODO: 调用保存接口
}

async function handlePublish() {
  if (!formData.value.title) {
    notification.error({
      message: $t('page.post.validation.titleRequired'),
    });
    return;
  }
  if (!formData.value.content) {
    notification.error({
      message: $t('page.post.validation.contentRequired'),
    });
    return;
  }

  try {
    await (isCreateMode.value
      ? postStore.createPost({
          editorType: convertToEditorType(formData.value.editorType),
          translations: [
            {
              title: formData.value.title,
              content: formData.value.content,
              languageCode: formData.value.lang,
            },
          ],
        })
      : postStore.updatePost(formData.value.id || 0, {
          editorType: convertToEditorType(formData.value.editorType),
          translations: [
            {
              title: formData.value.title,
              content: formData.value.content,
              languageCode: formData.value.lang,
            },
          ],
        }));

    notification.success({
      message: $t('page.post.validation.publishSuccess'),
    });

    goBack();
  } catch (error) {
    console.error('Failed to publish post:', error);
    notification.error({
      message: $t('page.post.validation.publishFailed'),
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
      console.error('Post not found');
      return;
    }

    if (!item.translations || item.translations.length === 0) {
      console.error('No translations found for post');
      return;
    }

    let langItem = item.translations?.find(
      (t) => t.languageCode === formData.value.lang,
    );
    if (!langItem) {
      langItem = item.translations?.[0];
      formData.value.lang = langItem?.languageCode || formData.value.lang;
    }

    if (!langItem) {
      console.error('No translations found for post');
      return;
    }

    formData.value.id = item.id;
    formData.value.title = langItem.title || '';
    formData.value.content = langItem.content || '';
    formData.value.lang = formData.value.lang || 'zh-CN';
    formData.value.editorType = convertToUIEditorType(item.editorType);
  } catch (error) {
    console.error('Failed to load post:', error);
  }
}

onMounted(async () => {
  await reloadLanguages();

  if (isCreateMode.value) {
    formData.value.title = $t('page.post.placeholder.untitled');
  } else {
    await loadPost();
  }
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
          v-model:value="formData.title"
          :placeholder="$t('page.post.placeholder.title')"
          size="large"
          class="flex-1"
        />
        <a-select
          :value="formData.lang"
          style="width: 200px"
          @change="handleLanguageChange"
        >
          <a-select-option
            v-for="option in languageOptions"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </a-select-option>
        </a-select>
        <a-select v-model:value="formData.editorType" style="width: 200px">
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
        v-model="formData.content"
        :editor-type="formData.editorType"
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
