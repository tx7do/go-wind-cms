<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue';

import { $t } from '@vben/locales';
import { preferences } from '@vben/preferences';

import { type EditorProps, MdEditor } from 'md-editor-v3';

import { isDarkMode } from './utils';

import 'md-editor-v3/lib/style.css';

interface Props {
  modelValue: string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  options?: Partial<EditorProps>;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  height: '100%', // 默认撑满父容器
  placeholder: $t('ui.editor.please_input_content'),
  options: () => ({}),
});

const emit = defineEmits<{
  (e: 'change', value: string): void;
  (e: 'ready'): void;
  (e: 'update:modelValue', value: string): void;
}>();

const localValue = ref(props.modelValue);
const isDark = ref(isDarkMode());
const wrapperRef = ref<HTMLDivElement>(null); // 修正类型定义
let resizeObserver: null | ResizeObserver = null; // 监听容器尺寸变化

// 编辑器配置
const editorProps = computed<EditorProps>(() => ({
  previewOnly: false,
  preview: true,
  toolbars: [
    'bold',
    'underline',
    'italic',
    '-',
    'strikeThrough',
    'title',
    'sub',
    'sup',
    '-',
    'quote',
    'unorderedList',
    'orderedList',
    'task',
    '-',
    'codeRow',
    'code',
    'link',
    'image',
    'table',
    'mermaid',
    'katex',
    '-',
    'revoke',
    'next',
    'save',
    '=',
    'pageFullscreen',
    'fullscreen',
    'preview',
    'htmlPreview',
    'catalog',
  ],
  ...props.options,
}));

// 响应式高度（传给编辑器）
const editorHeight = ref<number>(600); // 初始兜底高度

// 计算编辑器主题
const theme = computed(() => (isDark.value ? 'dark' : 'light'));

// 核心：计算编辑器应有的高度（撑满可视区域）
const updateEditorHeight = () => {
  if (!wrapperRef.value) return;

  // 获取父容器的实际可用高度（撑满父容器）
  const containerRect = wrapperRef.value.getBoundingClientRect();
  const containerHeight = containerRect.height;

  // 如果传入height是100%，则使用容器高度；否则使用传入值
  if (props.height === '100%' || props.height === '100vh') {
    // 特殊处理：如果是100vh则使用视口高度
    const finalHeight =
      props.height === '100vh' ? window.innerHeight : containerHeight;
    // 避免高度为0的异常情况
    editorHeight.value = finalHeight > 0 ? finalHeight : 600; // 兜底高度600px
  } else if (typeof props.height === 'number') {
    editorHeight.value = props.height;
  } else if (typeof props.height === 'string') {
    // 处理带px的字符串（如"800px"）
    const pxMatch = props.height.match(/^(\d+)px$/);
    if (pxMatch) {
      editorHeight.value = Number(pxMatch[1]);
    } else {
      editorHeight.value = 600; // 兜底
    }
  }
};

const handleChange = (value: string) => {
  localValue.value = value;
  emit('update:modelValue', value);
  emit('change', value);
};

// 监听外部值变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== localValue.value) {
      localValue.value = newVal;
    }
  },
);

// 监听主题变化
watch(
  () => preferences.theme.mode,
  (mode) => {
    isDark.value = mode === 'dark';
    nextTick(() => updateEditorHeight()); // 主题切换后重新计算高度
  },
);

onMounted(() => {
  emit('ready');
  nextTick(() => {
    updateEditorHeight(); // 初始化高度

    // 监听容器尺寸变化（窗口缩放/父容器变化时自动调整）
    if (wrapperRef.value) {
      resizeObserver = new ResizeObserver(() => {
        updateEditorHeight();
      });
      resizeObserver.observe(wrapperRef.value);

      // 监听窗口大小变化（兜底）
      window.addEventListener('resize', updateEditorHeight);
    }
  });
});

// 组件卸载时清理监听
onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect();
    resizeObserver = null;
  }
  window.removeEventListener('resize', updateEditorHeight);
});

// 监听props.height变化
watch(
  () => props.height,
  () => {
    nextTick(updateEditorHeight);
  },
);
</script>

<template>
  <div
    ref="wrapperRef"
    class="md-editor-wrapper"
    :style="{ height: props.height === '100vh' ? '100vh' : '100%' }"
  >
    <MdEditor
      v-model="localValue"
      :theme="theme"
      :placeholder="placeholder"
      :preview-only="false"
      :disabled="disabled"
      v-bind="editorProps"
      @change="handleChange"
      class="md-editor-inner"
    />
  </div>
</template>

<style scoped>
.md-editor-wrapper {
  width: 100%;
  min-height: 1px;
  box-sizing: border-box;
  overflow: hidden;
}

:deep(.md-editor-inner) {
  width: 100% !important;
  height: 100% !important;
}

:deep(.m-md-editor) {
  width: 100% !important;
  height: 100% !important;
  min-height: unset !important;
}

:deep(.m-md-editor .m-md-content) {
  height: calc(100% - 40px) !important;
}

:deep(.m-md-editor .m-md-preview),
:deep(.m-md-editor .m-md-edit-area) {
  height: 100% !important;
  width: 100% !important;
}

:deep(.m-md-editor.dark) {
  height: 100% !important;
}
</style>
