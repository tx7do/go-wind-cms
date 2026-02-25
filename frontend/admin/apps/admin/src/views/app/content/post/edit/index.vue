<script setup lang="ts">
import type { Props } from './types';

import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { Page, useVbenModal } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import { LucideArrowLeft } from '@vben/icons';
import { $t } from '@vben/locales';

import { Editor, EditorType } from '#/components/editor';
import { useFileTransferStore, useLanguageStore, usePostStore } from "#/stores";

import PublishPostModal from './publish-post-modal.vue';

const postStore = usePostStore();
const languageStore = useLanguageStore();
const fileTransferStore = useFileTransferStore();

const route = useRoute();
const { closeCurrentTab } = useTabs();

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

// 编辑器类型列表
const editorTypeOptions = [
  {
    label: $t('enum.editorType.EDITOR_TYPE_MARKDOWN'),
    value: EditorType.MARKDOWN,
  },
  {
    label: $t('enum.editorType.EDITOR_TYPE_RICH_TEXT'),
    value: EditorType.RICH_TEXT,
  },
  {
    label: $t('enum.editorType.EDITOR_TYPE_JSON_BLOCK'),
    value: EditorType.JSON,
  },
  {
    label: $t('enum.editorType.EDITOR_TYPE_PLAIN_TEXT'),
    value: EditorType.PLAIN_TEXT,
  },
  { label: $t('enum.editorType.EDITOR_TYPE_CODE'), value: EditorType.CODE },
];

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

// 表单数据
const formData = ref<Props>({
  title: '',
  content: '',
  lang: 'zh-CN',
  editorType: EditorType.MARKDOWN,
});

const [Modal, modalApi] = useVbenModal({
  // 连接抽离的组件
  connectedComponent: PublishPostModal,
});

/* 打开模态窗口 */
function openModal() {
  modalApi.setData(formData.value);
  modalApi.open();
}

function goBack() {
  closeCurrentTab();
  // router.push('/content/posts');
}

function handleSave() {
  console.log('Save post:', formData.value);
  // TODO: 调用保存接口
}

function handlePublish() {
  openModal();
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

onMounted(async () => {
  await reloadLanguages();

  if (isCreateMode.value) {
    formData.value.title = $t('page.post.placeholder.untitled');
  } else {
    const item = await postStore.getPost(postId.value || 0);
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
        <a-select v-model:value="formData.lang" style="width: 200px">
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
          <a-button type="default" @click="handleSave">
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
