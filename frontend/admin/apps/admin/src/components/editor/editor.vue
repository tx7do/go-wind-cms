<script setup lang="ts">
import type { EditorEmits, EditorProps } from './types';

import { computed, defineAsyncComponent } from 'vue';

import { EditorType } from './types';

const props = withDefaults(defineProps<EditorProps>(), {
  editorType: EditorType.MARKDOWN,
  height: 500,
  disabled: false,
  placeholder: 'Please enter content...',
});

const emit = defineEmits<EditorEmits>();

// 懒加载编辑器组件
const UEditor = defineAsyncComponent(() => import('./UEditor.vue'));
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
    case EditorType.JSON: {
      return JsonEditor;
    }

    case EditorType.RICH_TEXT: {
      return UEditor;
    }

    case EditorType.MARKDOWN: {
      return MarkdownEditor;
    }

    case EditorType.PLAIN_TEXT: {
      return PlainTextEditor;
    }

    case EditorType.CODE: {
      return CodeEditor;
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
    case EditorType.MARKDOWN:
      return props.markdownOptions;
    case EditorType.JSON:
      return props.jsonOptions;
    case EditorType.CODE:
      return props.codeOptions;
    default:
      return undefined;
  }
});

</script>

<template>
  <div class="editor-container">
    <component
      :is="currentEditorComponent"
      :model-value="modelValue"
      :height="height"
      :disabled="disabled"
      :placeholder="placeholder"
      :config="ueditorConfig"
      :options="currentOptions"
      @update:model-value="handleUpdate"
      @change="handleChange"
      @ready="handleReady"
    />
  </div>
</template>

<style scoped>
.editor-container {
  width: 100%;
  min-height: 300px;
  height: 100%;
}
</style>
