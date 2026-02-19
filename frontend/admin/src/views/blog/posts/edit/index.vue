<template>
  <PageWrapper title="写文章" contentFullHeight fixedHeight contentBackground>
    <template #extra>
      <a-space size="middle">
        <a-button :loading="previewSaving" @click="handlePreviewClick">预览</a-button>
        <a-button type="primary" @click="postSettingVisible = true">发布</a-button>
      </a-space>
    </template>

    <a-row :gutter="6">
      <a-col :span="24">
        <div class="mb-4">
          <a-input v-model="postToStage.title" placeholder="请输入文章标题" size="large" />
        </div>
        <div id="editor">
          <MarkDown
            v-model:value="postToStage.originalContent"
            @change="handleChange"

            ref="markDownRef"
            placeholder="这是占位文本"
          />
        </div>
      </a-col>
    </a-row>
  </PageWrapper>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { MarkDown, MarkDownActionType } from '/@/components/Markdown';
  import { PageWrapper } from '/@/components/Page';
  import { Col, Row, Space } from 'ant-design-vue';

  const ASpace = Space;
  const ACol = Col;
  const ARow = Row;

  const markDownRef = ref<Nullable<MarkDownActionType>>(null);
  const valueRef = ref(``);

  let editorHeight = 500;

  let previewSaving = false;
  let postSettingVisible = false;
  let postToStage = {
    title: '',
    originalContent: '',
  };
  let contentChanges = 0;

  function handleChange(v: string) {
    valueRef.value = v;
  }

  function handlePreviewClick() {}
  function handleOpenPreview() {}
  function handleRestoreSavedStatus() {}

  function clearValue() {
    valueRef.value = '';
  }
</script>
