<script setup lang="ts">
import { computed, ref, watch } from 'vue';

import { $t } from '@/locales';

interface Props {
  modelValue: string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  height: '100%',
  placeholder: '',
});

const emit = defineEmits<{
  (e: 'change', value: string): void;
  (e: 'ready'): void;
  (e: 'update:modelValue', value: string): void;
}>();

const localValue = ref(props.modelValue);

// 监听外部值变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== localValue.value) {
      localValue.value = newVal;
    }
  },
  { immediate: true },
);

// 计算编辑器高度
const editorHeight = computed(() => {
  if (typeof props.height === 'number') {
    return `${props.height}px`;
  }
  return props.height;
});
const computedPlaceholder = computed(
  () => props.placeholder || $t('component.editor.html_placeholder'),
);

// 处理输入变化
const handleInput = (event: Event) => {
  const target = event.target as HTMLTextAreaElement;
  localValue.value = target.value;
  emit('update:modelValue', target.value);
  emit('change', target.value);
};

// 组件挂载后触发 ready 事件
defineExpose({
  focus: () => {
    const textarea = document.querySelector(
      '.plain-text-editor-textarea',
    ) as HTMLTextAreaElement | null;
    textarea?.focus();
  },
});
</script>

<template>
  <div class="plain-text-editor-container">
    <textarea
      :value="localValue"
      :disabled="props.disabled"
      :placeholder="computedPlaceholder"
      :style="{ height: editorHeight }"
      class="plain-text-editor-textarea"
      @input="handleInput"
    />
  </div>
</template>

<style scoped lang="less">
.plain-text-editor-container {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.plain-text-editor-textarea {
  width: 100%;
  height: 100%;
  min-height: 300px;
  padding: 12px 16px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
  font-size: 14px;
  line-height: 1.6;
  color: var(--editor-text, var(--color-text-primary));
  background-color: var(--editor-bg, var(--color-surface));
  border: 1px solid var(--editor-border, var(--color-border));
  border-radius: var(--radius-sm);
  resize: vertical;
  transition: all 0.3s ease;

  &::placeholder {
    color: var(--editor-placeholder, var(--color-text-secondary));
  }

  &:hover:not(:disabled) {
    border-color: var(--editor-focus-border, var(--color-brand));
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  }

  &:focus {
    outline: none;
    border-color: var(--editor-focus-border, var(--color-brand));
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.18);
  }

  &:disabled {
    background-color: var(--color-surface);
    cursor: not-allowed;
    opacity: 0.6;
  }

  // 滚动条美化
  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(102, 126, 234, 0.3);
    border-radius: 4px;

    &:hover {
      background: rgba(102, 126, 234, 0.6);
    }
  }
}

// 暗色模式
html.dark {
  .plain-text-editor-textarea {
    &::-webkit-scrollbar-thumb {
      background: rgba(129, 147, 255, 0.3);

      &:hover {
        background: rgba(129, 147, 255, 0.6);
      }
    }
  }
}
</style>

