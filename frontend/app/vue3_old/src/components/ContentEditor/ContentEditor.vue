<script setup lang="ts">
import type {EditorEmits, EditorProps} from './types';

import {computed, defineAsyncComponent} from 'vue';

import {$t} from '@/locales'

const props = withDefaults(defineProps<EditorProps>(), {
  editorType: 'markdown',
  height: '100%',
  disabled: false,
  placeholder: $t('ui.editor.please_input_content'),
});

const emit = defineEmits<EditorEmits>();

// 懒加载编辑器组件
const TiptapEditor = defineAsyncComponent(() => import('./TiptapEditor.vue'));
const MarkdownEditor = defineAsyncComponent(
  () => import('./MarkdownEditor.vue'),
);
const JsonEditor = defineAsyncComponent(() => import('./JsonEditor.vue'));
const PlainTextEditor = defineAsyncComponent(
  () => import('./PlainTextEditor.vue'),
);
const CodeEditor = defineAsyncComponent(() => import('./CodeEditor.vue'));

// 根据编辑器类型确定使用哪个编辑器
const currentEditorComponent = computed(() => {
  const type = props.editorType;

  switch (type) {
    case 'code': {
      return CodeEditor;
    }

    case 'json': {
      return JsonEditor;
    }

    case 'markdown': {
      return MarkdownEditor;
    }

    case 'text': {
      return PlainTextEditor;
    }

    case 'richtext': {
      return TiptapEditor;
    }

    default: {
      // 其他未知类型，默认使用 Markdown 编辑器
      return MarkdownEditor;
    }
  }
});

const handleUpdate = (value: string) => {
  emit('update:modelValue', value);
};

const handleChange = (value: string) => {
  emit('change', value);
};

const handleReady = () => {
  emit('ready');
};

// 计算当前编辑器的选项
const currentOptions = computed(() => {
  const type = props.editorType;
  switch (type) {
    case 'code': {
      return props.codeOptions;
    }
    case 'json': {
      return props.jsonOptions;
    }
    case 'markdown': {
      return props.markdownOptions;
    }
    default: {
      return undefined;
    }
  }
});
</script>

<template>
  <div class="editor-wrapper">
    <component
      :is="currentEditorComponent"
      :model-value="modelValue"
      :height="height"
      :disabled="disabled"
      :placeholder="placeholder"
      :upload-image="uploadImage"
      :options="currentOptions"
      @update:model-value="handleUpdate"
      @change="handleChange"
      @ready="handleReady"
    />
  </div>
</template>

<style scoped lang="less">
.editor-wrapper {
  width: 100%;
  height: 100%;
  min-height: 300px;
  display: flex;
  flex-direction: column;

  // 主题变量
  --editor-bg: var(--color-surface);
  --editor-text: var(--color-text-primary);
  --editor-border: var(--color-border);
  --editor-focus-border: var(--color-brand);
  --editor-placeholder: var(--color-text-secondary);

  // 子组件继承样式
  :deep(.editor-container),
  :deep(.markdown-editor-wrapper),
  :deep(.code-editor-wrapper),
  :deep(.json-editor-wrapper),
  :deep(.plain-text-editor-container) {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
  }
}

// 暗色模式适配
html.dark {
  .editor-wrapper {
    --editor-bg: var(--color-surface);
    --editor-text: var(--color-text-primary);
    --editor-border: var(--color-border);
    --editor-focus-border: var(--color-brand);
    --editor-placeholder: var(--color-text-secondary);
  }
}
</style>
