<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { LucideArrowLeft } from '@vben/icons';
import { $t } from '@vben/locales';

import { Editor, EditorType } from '#/components/editor';
import { router } from '#/router';

const route = useRoute();

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
  { label: 'Markdown Editor', value: EditorType.MARKDOWN },
  { label: 'Rich Text Editor (UEditor)', value: EditorType.RICH_TEXT },
  { label: 'JSON Editor', value: EditorType.JSON },
  { label: 'Plain Text Editor', value: EditorType.PLAIN_TEXT },
  { label: 'Code Editor', value: EditorType.CODE },
];

// 表单数据
const formData = ref({
  title: '',
  content: '',
  editorType: EditorType.MARKDOWN,
  excerpt: '',
  slug: '',
  status: 'POST_STATUS_DRAFT',
});

function goBack() {
  router.push('/content/posts');
}

function handleSave() {
  console.log('Save post:', formData.value);
  // TODO: 调用保存接口
}

function handlePublish() {
  formData.value.status = 'POST_STATUS_PUBLISHED';
  handleSave();
}
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
        v-model="formData.content"
        :editor-type="formData.editorType"
        height="100%"
        :placeholder="$t('page.post.placeholder.content')"
        class="h-full"
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
  </Page>
</template>

<style scoped>
.post-edit-container {
  width: 100%;
  height: 100%;
}
</style>
