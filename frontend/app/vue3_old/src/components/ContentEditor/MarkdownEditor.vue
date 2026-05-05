<script setup lang="ts">
import {computed, nextTick, onMounted, onUnmounted, ref, watch} from 'vue'

import {type EditorProps, MdEditor} from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

import {$t} from '@/locales'
import {useAppStore} from '@/stores'

import {isDarkMode} from './utils'

interface UseEditorConfigProps {
  modelValue: string
  height?: number | string
  disabled?: boolean
  placeholder?: string
  options?: Partial<EditorProps>
  enableExport?: boolean
  uploadImage?: (file: File) => Promise<string>
}

const props = withDefaults(defineProps<UseEditorConfigProps>(), {
  disabled: false,
  height: '100%', // 默认撑满父容器
  placeholder: '',
  options: () => ({}),
  enableExport: false,
  uploadImage: undefined,
})

const appStore = useAppStore()

const emit = defineEmits<{
  (e: 'change', value: string): void;
  (e: 'imageUpload', file: File): void;
  (e: 'ready'): void;
  (e: 'update:modelValue', value: string): void;
}>();

const localValue = ref(props.modelValue);
const isDark = ref(isDarkMode());
const wrapperRef = ref<HTMLDivElement>(); // 修正类型定义
let resizeObserver: null | ResizeObserver = null; // 监听容器尺寸变化

const toolbars = computed(() => {
  const base = [
    // === 文本格式 ===
    'bold',
    'underline',
    'italic',
    'strikeThrough',
    '-',

    // === 文本样式 ===
    'title',
    'sub',
    'sup',
    'alignLeft',
    'alignCenter',
    'alignRight',
    'alignJustify',
    '-',

    // === 列表与引用 ===
    'quote',
    'unorderedList',
    'orderedList',
    'task',
    'indent',
    'outdent',
    '-',

    // === 代码与插入 ===
    'codeRow',
    'code',
    'link',
    'image',
    'table',
    'horizontalRule',
    'emoji',
    'footnote',
    '-',

    // === 图表与公式 ===
    'mermaid',
    'katex',
    '-',

    // === 编辑操作 ===
    'revoke',
    'next',
    'clear',
    'save',
    '=',

    // === 视图与导出 ===
    'pageFullscreen',
    'fullscreen',
    'preview',
    'htmlPreview',
    'catalog',
    'help',
  ];

  // 导出功能（按需）
  if (props.enableExport) {
    base.splice(base.indexOf('save') + 1, 0, 'exportPdf', 'exportHtml', '-');
  }

  return base;
});

// 编辑器配置
const editorProps = computed(() => ({
  previewOnly: false,
  preview: true,
  showCodeRowNumber: true,
  noMermaid: false,
  noKatex: false,
  toolbars: toolbars.value as any,
  // 合并用户自定义配置
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
    // eslint-disable-next-line unicorn/prefer-ternary
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

async function doUploadImage(file: File): Promise<string> {
  console.log('Uploading image:', file);

  if (!file || !props.uploadImage) {
    emit('imageUpload', file!);
    return '';
  }

  try {
    return await props.uploadImage(file);
  } catch (error) {
    console.error('Image upload failed:', error);
    return '';
  }
}

const handleUploadImages = async (
  files: File[],
  callback: (urls: string[]) => void,
) => {
  const uploadPromises = files.map((file) => doUploadImage(file));
  const urls = await Promise.all(uploadPromises);
  callback(urls);
};

const handleSave = (val: string, _html: string) => {
  // 创建 Blob 对象（Markdown 格式）
  const blob = new Blob([val], {type: 'text/markdown;charset=utf-8'});

  // 生成文件名（带时间戳）
  const timestamp = new Date().toISOString().slice(0, 19).replaceAll(':', '-');
  const filename = `document-${timestamp}.md`;

  // 创建下载链接
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;

  // 触发下载
  document.body.append(link);
  link.click();

  // 清理
  link.remove();
  URL.revokeObjectURL(url);
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
  () => appStore.mode,
  (mode) => {
    isDark.value = mode === 'dark'
    nextTick(() => updateEditorHeight())
  },
)

// 监听props.height变化
watch(
  () => props.height,
  () => {
    nextTick(updateEditorHeight);
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
      :placeholder="props.placeholder || $t('component.editor.markdown_placeholder')"
      :preview-only="false"
      :disabled="props.disabled"
      v-bind="editorProps"
      @change="handleChange"
      @on-upload-img="handleUploadImages"
      @on-save="handleSave"
      class="md-editor-inner"
    />
  </div>
</template>

<style scoped lang="less">
.md-editor-wrapper {
  box-sizing: border-box;
  width: 100%;
  height: 100%;
  min-height: 300px;
  overflow: hidden;
  background-color: var(--editor-bg, var(--color-surface));
  border: 1px solid var(--editor-border, var(--color-border));
  border-radius: var(--radius-md);
  transition: all 0.3s ease;

  &:hover {
    border-color: var(--editor-focus-border, var(--color-brand));
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  }

  &:focus-within {
    border-color: var(--editor-focus-border, var(--color-brand));
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.18);
  }
}

.md-editor-wrapper {
  :deep(.md-editor-inner) {
    width: 100% !important;
    height: 100% !important;
    border: none !important;
  }

  :deep(.md-editor),
  :deep(.md-editor-dark),
  :deep(.md-editor-preview),
  :deep(.md-editor-preview-wrapper) {
    background: var(--editor-bg, var(--color-surface)) !important;
    color: var(--editor-text, var(--color-text-primary)) !important;
  }

  :deep(.md-editor .md-editor-content) {
    height: calc(100% - 40px) !important;
    background: var(--editor-bg, var(--color-surface)) !important;
  }

  :deep(.md-editor .md-editor-preview),
  :deep(.md-editor .md-editor-input-wrapper),
  :deep(.md-editor .md-editor-input) {
    width: 100% !important;
    height: 100% !important;
    background: var(--editor-bg, var(--color-surface)) !important;
    color: var(--editor-text, var(--color-text-primary)) !important;
  }

  :deep(.md-editor .md-editor-toolbar) {
    background: var(--color-surface) !important;
    border-color: var(--editor-border, var(--color-border)) !important;
  }

  :deep(.md-editor .md-editor-toolbar-left button),
  :deep(.md-editor .md-editor-toolbar-right button),
  :deep(.md-editor .md-editor-toolbar-item) {
    color: var(--editor-text, var(--color-text-primary)) !important;

    &:hover {
      background: rgba(102, 126, 234, 0.08) !important;
      color: var(--editor-focus-border, var(--color-brand)) !important;
    }
  }

  :deep(.md-editor .cm-editor),
  :deep(.md-editor .cm-scroller),
  :deep(.md-editor .cm-content) {
    background: var(--editor-bg, var(--color-surface)) !important;
    color: var(--editor-text, var(--color-text-primary)) !important;
  }

  :deep(.md-editor .cm-gutters) {
    background: var(--color-surface) !important;
    border-color: var(--editor-border, var(--color-border)) !important;
  }

  // md-editor-v3 的分隔条
  :deep(.md-editor-preview-wrapper) {
    border-left: 1px solid var(--editor-border, var(--color-border)) !important;
  }

  // 编辑区和预览区之间的分隔条（通用选择器）
  :deep(.md-editor > div > div) {
    border-color: var(--editor-border, var(--color-border)) !important;
  }

  :deep(.md-editor.dark) {
    height: 100% !important;
  }

  // 分割条样式优化 - 覆盖所有可能的分割器
  :deep(.md-editor-splitter),
  :deep(.md-editor-resize-bar),
  :deep(.md-editor .divider),
  :deep([class*='split']),
  :deep([class*='divider']),
  :deep([class*='resize']) {
    background-color: var(--editor-border, var(--color-border)) !important;
    opacity: 0.5 !important;
    transition: all 0.3s ease !important;

    &:hover {
      background-color: var(--color-brand) !important;
      opacity: 0.8 !important;
    }
  }

  // 滚动条美化
  :deep(.md-editor ::-webkit-scrollbar),
  :deep(.cm-scroller::-webkit-scrollbar) {
    width: 8px;
    height: 8px;
  }

  :deep(.md-editor ::-webkit-scrollbar-track),
  :deep(.cm-scroller::-webkit-scrollbar-track) {
    background: transparent;
  }

  :deep(.md-editor ::-webkit-scrollbar-thumb),
  :deep(.cm-scroller::-webkit-scrollbar-thumb) {
    background: rgba(102, 126, 234, 0.3);
    border-radius: 4px;

    &:hover {
      background: rgba(102, 126, 234, 0.6);
    }
  }
}

// 暗色模式
html.dark {
  .md-editor-wrapper {
    :deep(.md-editor-preview-wrapper) {
      border-left-color: var(--color-border) !important;
    }

    :deep(.md-editor > div > div) {
      border-color: var(--color-border) !important;
    }

    :deep(.md-editor-splitter),
    :deep(.md-editor-resize-bar),
    :deep(.md-editor .divider),
    :deep([class*='split']),
    :deep([class*='divider']),
    :deep([class*='resize']) {
      background-color: var(--color-border) !important;
      opacity: 0.6 !important;

      &:hover {
        background-color: var(--color-brand) !important;
        opacity: 1 !important;
      }
    }

    :deep(.md-editor ::-webkit-scrollbar-thumb),
    :deep(.cm-scroller::-webkit-scrollbar-thumb) {
      background: rgba(129, 147, 255, 0.3);

      &:hover {
        background: rgba(129, 147, 255, 0.6);
      }
    }
  }
}
</style>
