<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue';

import { preferences } from '@vben/preferences';

import VueMarkdownEditor from '@kangc/v-md-editor';
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js';
import hljs from 'highlight.js/lib/core';
import javascript from 'highlight.js/lib/languages/javascript';
import Prism from 'prismjs';

import { isDarkMode } from './utils';

// 引入Prism暗黑高亮样式
import 'prismjs/themes/prism-tomorrow.css';
import '@kangc/v-md-editor/lib/style/base-editor.css';
import '@kangc/v-md-editor/lib/theme/style/github.css';
import 'highlight.js/styles/github-dark.css';

interface Props {
  modelValue: string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  options?: {
    hideModeSwitch?: boolean;
    initialEditType?: 'markdown' | 'wysiwyg';
    previewStyle?: 'tab' | 'vertical';
    toolbarItems?: string[][];
  };
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  height: '100%',
  placeholder: 'Please enter markdown content...',
  options: () => ({
    hideModeSwitch: false,
    initialEditType: 'markdown',
    previewStyle: 'vertical',
  }),
});
const emit = defineEmits<{
  (e: 'change', value: string): void;
  (e: 'ready'): void;
  (e: 'update:modelValue', value: string): void;
}>();

hljs.registerLanguage('javascript', javascript);

// 初始化主题时配置Prism（确保代码高亮生效）
VueMarkdownEditor.use(githubTheme, {
  Prism,
  Hljs: hljs,
  // 自定义代码高亮样式前缀（适配暗黑模式）
  codeHighlightExtensionMap: {
    vue: 'html',
    js: 'javascript',
  },
});

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

// 提取暗黑样式应用逻辑为独立函数
const applyDarkModeStyles = () => {
  const editors = document.querySelectorAll('.v-md-editor');
  const textareas = document.querySelectorAll('.v-md-editor-textarea');
  const previews = document.querySelectorAll('.v-md-preview');

  if (isDarkMode()) {
    // 应用暗黑样式
    editors.forEach((el) => {
      const editor = el as HTMLElement;
      editor.style.backgroundColor = '#0f172a';
      editor.style.color = '#ffffff';
      editor.style.borderColor = '#1e293b';
    });

    textareas.forEach((el) => {
      const textarea = el as HTMLElement;
      textarea.style.backgroundColor = '#0f172a';
      textarea.style.color = '#ffffff';
      textarea.style.caretColor = '#60a5fa';
    });

    previews.forEach((el) => {
      const preview = el as HTMLElement;
      preview.style.backgroundColor = '#0f172a';
      preview.style.color = '#ffffff';
    });
  } else {
    // 恢复默认样式
    editors.forEach((el) => {
      const editor = el as HTMLElement;
      editor.style.backgroundColor = '';
      editor.style.color = '';
      editor.style.borderColor = '';
    });

    textareas.forEach((el) => {
      const textarea = el as HTMLElement;
      textarea.style.backgroundColor = '';
      textarea.style.color = '';
      textarea.style.caretColor = '';
    });

    previews.forEach((el) => {
      const preview = el as HTMLElement;
      preview.style.backgroundColor = '';
      preview.style.color = '';
    });
  }
};

watch(
  () => preferences.theme.mode,
  () => {
    nextTick(() => {
      applyDarkModeStyles();
    });
  },
  { immediate: true },
);

const handleChange = (value: string) => {
  localValue.value = value;
  emit('update:modelValue', value);
  emit('change', value);
};

onMounted(() => {
  emit('ready');
  // 初始化时应用样式
  nextTick(() => {
    applyDarkModeStyles();
  });
});
</script>

<template>
  <div class="markdown-editor-wrapper" :class="{ 'is-dark': isDarkMode() }">
    <VueMarkdownEditor
      v-model="localValue"
      :disabled="disabled"
      :height="editorHeight"
      :placeholder="placeholder"
      class="markdown-editor"
      @change="handleChange"
    />
  </div>
</template>

<style scoped>
.markdown-editor-wrapper {
  height: 100%;
  width: 100%;
  /* 提升样式优先级 */
  --v-md-editor-dark-bg: #0f172a;
  --v-md-editor-dark-text: #ffffff;
  --v-md-editor-dark-border: #1e293b;
  --v-md-editor-dark-hover: #334155;
  --v-md-editor-primary: #60a5fa;
  --v-md-editor-active-icon: #60a5fa;
  --v-md-editor-active-border: #60a5fa;
}

.markdown-editor {
  height: 100%;
  width: 100%;
  /* 确保样式继承 */
  color: inherit;
  background-color: inherit;
}

/* 确保深层button样式传递 */
.markdown-editor :deep(button) {
  /* 基础样式 */
}

/* 暗黑模式 - 全局重置（提升优先级） */
.markdown-editor-wrapper.is-dark :deep(.v-md-editor *) {
  --v-md-editor-bg-color: var(--v-md-editor-dark-bg) !important;
  --v-md-editor-text-color: var(--v-md-editor-dark-text) !important;
  color: var(--v-md-editor-dark-text) !important;
}

/* 暗黑模式 - 主容器 */
.markdown-editor-wrapper.is-dark :deep(.v-md-editor) {
  background-color: var(--v-md-editor-dark-bg) !important;
  color: var(--v-md-editor-dark-text) !important;
  border: 1px solid var(--v-md-editor-dark-border) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__main) {
  background-color: var(--v-md-editor-dark-bg) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor-container) {
  background-color: var(--v-md-editor-dark-bg) !important;
}

/* 暗黑模式 - 工具栏 */
.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar) {
  background-color: var(--v-md-editor-dark-border) !important;
  border-bottom: 1px solid var(--v-md-editor-dark-hover) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar button) {
  color: #cbd5e1 !important;
  background-color: transparent !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar button svg) {
  color: inherit !important;
  fill: currentColor !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar button:hover) {
  background-color: var(--v-md-editor-dark-hover) !important;
  color: #cbd5e1 !important;
  cursor: pointer !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar button:focus) {
  background-color: var(--v-md-editor-dark-hover) !important;
  color: #cbd5e1 !important;
  outline: none !important;
  box-shadow: none !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-editor__toolbar button:focus-visible) {
  background-color: var(--v-md-editor-dark-hover) !important;
  color: #cbd5e1 !important;
  outline: 1px solid var(--v-md-editor-primary) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar button:disabled) {
  color: #64748b !important;
  cursor: not-allowed !important;
  opacity: 0.5 !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar button.active) {
  background-color: var(--v-md-editor-dark-border) !important;
  color: var(--v-md-editor-active-icon) !important;
  border-color: var(--v-md-editor-active-border) !important;
  /* 增加内阴影，强化选中视觉效果 */
  box-shadow: inset 0 0 0 1px var(--v-md-editor-active-border) !important;
}

/* 确保所有图标元素都继承选中色 */
.markdown-editor-wrapper.is-dark,
:deep(.v-md-editor__toolbar button.active svg) {
  color: var(--v-md-editor-active-icon) !important;
  fill: var(--v-md-editor-active-icon) !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-editor__toolbar button.active svg path) {
  fill: var(--v-md-editor-active-icon) !important;
  stroke: var(--v-md-editor-active-icon) !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-editor__toolbar button.active svg use) {
  color: var(--v-md-editor-active-icon) !important;
  fill: var(--v-md-editor-active-icon) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar button.active i) {
  color: var(--v-md-editor-active-icon) !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-editor__toolbar button.active span) {
  color: var(--v-md-editor-active-icon) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar .dividing-line) {
  background-color: var(--v-md-editor-dark-hover) !important;
}

/* 暗黑模式 - 工具栏下拉菜单（补充遗漏样式） */
.markdown-editor-wrapper.is-dark :deep(.v-md-editor__toolbar .dropdown-menu) {
  background-color: var(--v-md-editor-dark-border) !important;
  border-color: var(--v-md-editor-dark-hover) !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-editor__toolbar .dropdown-menu li) {
  color: var(--v-md-editor-dark-text) !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-editor__toolbar .dropdown-menu li:hover) {
  background-color: var(--v-md-editor-dark-hover) !important;
}

/* 暗黑模式 - 编辑区 */
.markdown-editor-wrapper.is-dark :deep(.v-md-editor__left) {
  background-color: var(--v-md-editor-dark-bg) !important;
  border-right: 1px solid var(--v-md-editor-dark-border) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor__left-container) {
  background-color: var(--v-md-editor-dark-bg) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor-textarea) {
  background-color: var(--v-md-editor-dark-bg) !important;
  color: var(--v-md-editor-dark-text) !important;
  caret-color: var(--v-md-editor-primary) !important;
  border: none !important;
  font-size: 14px !important;
  line-height: 1.6 !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-editor-textarea::placeholder) {
  color: #94a3b8 !important;
  opacity: 1 !important; /* 确保占位符颜色生效 */
}

/* 暗黑模式 - 编辑区行号（补充遗漏） */
.markdown-editor-wrapper.is-dark :deep(.v-md-editor__line-number) {
  color: #94a3b8 !important;
  background-color: var(--v-md-editor-dark-bg) !important;
  border-right: 1px solid var(--v-md-editor-dark-border) !important;
}

/* 暗黑模式 - 预览区 */
.markdown-editor-wrapper.is-dark :deep(.v-md-editor__right) {
  background-color: var(--v-md-editor-dark-bg) !important;
  border-left: 1px solid var(--v-md-editor-dark-border) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-preview),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview-wrapper),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview__container),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview-html),
.markdown-editor-wrapper.is-dark :deep(.v-md-editor-preview) {
  background-color: var(--v-md-editor-dark-bg) !important;
  color: var(--v-md-editor-dark-text) !important;
  padding: 16px !important;
  line-height: 1.8 !important;
}

/* 暗黑模式 - 文本基础样式（提升优先级） */
.markdown-editor-wrapper.is-dark :deep(.v-md-preview p),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview strong),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview em),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview ul li),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview ol li) {
  color: var(--v-md-editor-dark-text) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-preview em) {
  color: #e8e8ff !important;
}

/* 暗黑模式 - 标题 */
.markdown-editor-wrapper.is-dark :deep(.v-md-preview h1),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview h2),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview h3),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview h4),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview h5),
.markdown-editor-wrapper.is-dark :deep(.v-md-preview h6) {
  color: var(--v-md-editor-dark-text) !important;
  border-bottom-color: var(--v-md-editor-dark-hover) !important;
}

/* 暗黑模式 - 链接 */
.markdown-editor-wrapper.is-dark :deep(.v-md-preview a) {
  color: var(--v-md-editor-primary) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-preview a:hover) {
  color: #93c5fd !important;
}

/* 暗黑模式 - 代码块（适配Prism暗黑样式） */
.markdown-editor-wrapper.is-dark :deep(.v-md-preview code) {
  background-color: var(--v-md-editor-dark-border) !important;
  color: #f1f5f9 !important;
  border-radius: 3px !important;
  padding: 2px 6px !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-preview pre) {
  background-color: var(--v-md-editor-dark-border) !important;
  border: 1px solid var(--v-md-editor-dark-hover) !important;
  color: #f1f5f9 !important;
  padding: 12px !important;
  overflow-x: auto !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-preview pre code) {
  background-color: transparent !important;
  color: inherit !important;
}

/* 暗黑模式 - 引用块 */
.markdown-editor-wrapper.is-dark :deep(.v-md-preview blockquote) {
  background-color: var(--v-md-editor-dark-border) !important;
  border-left: 3px solid var(--v-md-editor-primary) !important;
  color: var(--v-md-editor-dark-text) !important;
  padding: 12px 16px !important;
  margin: 8px 0 !important;
}

/* 暗黑模式 - 表格 */
.markdown-editor-wrapper.is-dark :deep(.v-md-preview table) {
  border-collapse: collapse !important;
  border: 1px solid var(--v-md-editor-dark-hover) !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-preview table th) {
  background-color: var(--v-md-editor-dark-border) !important;
  border: 1px solid var(--v-md-editor-dark-hover) !important;
  color: var(--v-md-editor-dark-text) !important;
  padding: 8px 12px !important;
  font-weight: 600 !important;
}

.markdown-editor-wrapper.is-dark :deep(.v-md-preview table td) {
  border: 1px solid var(--v-md-editor-dark-hover) !important;
  color: var(--v-md-editor-dark-text) !important;
  padding: 8px 12px !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-preview table tr:nth-child(even) td) {
  background-color: var(--v-md-editor-dark-bg) !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-preview table tr:nth-child(odd) td) {
  background-color: var(--v-md-editor-dark-border) !important;
}

/* 暗黑模式 - 分隔线 */
.markdown-editor-wrapper.is-dark :deep(.v-md-preview hr) {
  border-color: var(--v-md-editor-dark-hover) !important;
}

/* 暗黑模式 - 滚动条 */
.markdown-editor-wrapper.is-dark :deep(.v-md-editor__left::-webkit-scrollbar),
.markdown-editor-wrapper.is-dark :deep(.v-md-editor__right::-webkit-scrollbar) {
  width: 8px !important;
}

.markdown-editor-wrapper.is-dark,
:deep(.v-md-editor__left::-webkit-scrollbar-track),
.markdown-editor-wrapper.is-dark
  :deep(.v-md-editor__right::-webkit-scrollbar-track) {
  background-color: var(--v-md-editor-dark-bg) !important;
}

.markdown-editor-wrapper.is-dark
  :deep(.v-md-editor__left::-webkit-scrollbar-thumb),
.markdown-editor-wrapper.is-dark
  :deep(.v-md-editor__right::-webkit-scrollbar-thumb) {
  background-color: var(--v-md-editor-dark-hover) !important;
  border-radius: 4px !important;
}

.markdown-editor-wrapper.is-dark
  :deep(.v-md-editor__left::-webkit-scrollbar-thumb:hover),
.markdown-editor-wrapper.is-dark
  :deep(.v-md-editor__right::-webkit-scrollbar-thumb:hover) {
  background-color: #475569 !important;
}
</style>
