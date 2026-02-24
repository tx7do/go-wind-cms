<script setup lang="ts">
import {
  computed,
  getCurrentInstance,
  nextTick,
  onMounted,
  onUnmounted,
  ref,
  watch,
} from 'vue';

import { preferences } from '@vben/preferences';

import Vue3JsonEditor from 'json-editor-vue';

import { isDarkMode } from '#/adapter/component/UEditor/src/utils';

import 'jsoneditor/dist/jsoneditor.min.css';

// 类型定义抽离
type EditorMode = 'code' | 'form' | 'text' | 'tree' | 'view';
interface Props {
  modelValue: string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  options?: {
    mode?: EditorMode;
    modes?: EditorMode[];
    search?: boolean;
  };
}

// Props定义（移除冗余类型断言）
const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  height: 500,
  placeholder: '{}',
  options: () => ({
    mode: 'text',
    modes: ['tree', 'code', 'form', 'text', 'view'],
    search: true,
  }),
});

const emit = defineEmits<{
  (e: 'change', value: string): void;
  (e: 'error', error: Error): void;
  (e: 'ready'): void;
  (e: 'update:modelValue', value: string): void;
}>();

// 响应式数据
const localValue = ref(props.modelValue);
const jsonData = ref<any[] | null | Record<string, any>>(null);
const parseError = ref<string>('');
const isValidJson = ref(true);
const instance = getCurrentInstance();
let observer: MutationObserver | null = null; // 声明观察者变量，便于销毁

// 验证格式化JSON（优化逻辑）
const validateAndFormat = (value: string) => {
  try {
    if (!value?.trim()) {
      parseError.value = '';
      isValidJson.value = true;
      return { parsed: null, formatted: '' };
    }
    const parsed = JSON.parse(String(value));
    const formatted = JSON.stringify(parsed, null, 2);
    parseError.value = '';
    isValidJson.value = true;
    return { parsed, formatted };
  } catch (error) {
    const err = error as Error;
    parseError.value = `JSON解析错误: ${err.message || '未知错误'}`;
    isValidJson.value = false;
    emit('error', err);
    return { parsed: null, formatted: value };
  }
};

// 初始化数据
const initData = () => {
  const { parsed, formatted } = validateAndFormat(props.modelValue);
  localValue.value = formatted || props.placeholder;
  try {
    jsonData.value = parsed || JSON.parse(props.placeholder);
  } catch {
    jsonData.value = {};
  }
};

// 监听外部值变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== localValue.value) {
      const { parsed, formatted } = validateAndFormat(newVal);
      localValue.value = formatted || newVal || props.placeholder;
      try {
        jsonData.value = parsed || JSON.parse(props.placeholder);
      } catch {
        jsonData.value = {};
      }
    }
  },
  { immediate: true, deep: false },
);

// 监听编辑器内部数据变化
watch(
  () => jsonData.value,
  (newVal) => {
    if (newVal === null) return;
    try {
      const newValue = JSON.stringify(newVal, null, 2);
      if (newValue !== localValue.value) {
        localValue.value = newValue;
        emit('update:modelValue', newValue);
        emit('change', newValue);
      }
      parseError.value = '';
      isValidJson.value = true;
    } catch (error) {
      const err = error as Error;
      parseError.value = `JSON序列化错误: ${err.message || '未知错误'}`;
      isValidJson.value = false;
      emit('error', err);
    }
  },
  { deep: true },
);

// 高度计算（优化类型安全）
const editorHeight = computed(() => {
  let baseHeight = 500;
  if (typeof props.height === 'number') {
    baseHeight = props.height;
  } else if (typeof props.height === 'string') {
    const parsedHeight = Number.parseInt(props.height, 10);
    baseHeight = Number.isNaN(parsedHeight) ? 500 : parsedHeight;
  }
  const finalHeight = Math.max(baseHeight - 40, 200);
  return `${finalHeight}px`;
});

// 暗黑模式类名
const editorClass = computed(() => ({
  'json-editor-container': true,
  'json-editor-dark': isDarkMode(),
}));

// 优化刷新逻辑：移除内联样式硬编码，完全依赖CSS
const refreshEditor = () => {};

// 监听主题变化
watch(
  () => preferences.theme.mode,
  () => {
    refreshEditor();
  },
  { immediate: true },
);

// 监听编辑器模式变化
watch(
  () => props.options?.mode,
  () => {
    refreshEditor();
  },
);

// 编辑器change事件处理
const handleEditorChange = (value: any[] | Record<string, any>) => {
  jsonData.value = value;
  refreshEditor();
};

// 初始化和销毁逻辑
onMounted(() => {
  initData();
  nextTick(() => {
    emit('ready');
    refreshEditor();

    // 初始化MutationObserver（优化销毁逻辑）
    if (instance?.el) {
      const container = instance.el as HTMLElement;
      const editorEl = container.querySelector('.json-editor-core');
      if (editorEl) {
        observer = new MutationObserver((mutations) => {
          // 仅在样式/类名变化且暗黑模式开启时刷新
          const hasStyleChange = mutations.some(
            (m) =>
              m.type === 'attributes' &&
              ['class', 'style'].includes(m.attributeName || ''),
          );
          if (isDarkMode() && hasStyleChange) {
            refreshEditor();
          }
        });

        observer.observe(editorEl, {
          childList: true,
          subtree: true,
          attributes: true,
          attributeFilter: ['class', 'style'],
        });
      }
    }
  });
});

// 销毁观察者，防止内存泄漏
onUnmounted(() => {
  if (observer) {
    observer.disconnect();
    observer = null;
  }
});
</script>

<template>
  <div :class="editorClass">
    <!-- 错误提示 -->
    <div v-if="parseError" class="error-message">
      {{ parseError }}
    </div>

    <!-- 核心编辑器：移除所有冗余类型断言 -->
    <Vue3JsonEditor
      v-model="jsonData"
      :mode="options.mode"
      :modes="options.modes"
      :disabled="disabled"
      :search="options.search"
      :placeholder="placeholder"
      :style="{ height: editorHeight, width: '100%' }"
      class="json-editor-core"
      @change="handleEditorChange"
    />
  </div>
</template>

<style scoped>
/* 基础容器样式：合并类名，简化结构 */
.json-editor-container {
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  width: 100%;
  transition: all 0.2s ease;
}

/* 暗黑模式变量：统一管理，提升优先级 */
.json-editor-dark {
  --bg-primary: #0f172a !important;
  --bg-secondary: #1e293b !important;
  --text-primary: #ffffff !important;
  --text-secondary: #94a3b8 !important;
  --border-primary: #1e293b !important;
  --border-secondary: #334155 !important;
  --error-bg: #2a1a1a !important;
  --error-text: #ff8888 !important;
  --error-border: #4a2222 !important;

  border-color: var(--border-primary) !important;
  background-color: var(--bg-primary) !important;
}

/* 错误提示样式 */
.error-message {
  padding: 8px 12px;
  margin: 0;
  background-color: #fee;
  border-bottom: 1px solid #fcc;
  color: #c33;
  font-size: 12px;
  line-height: 1.4;
}

.json-editor-dark .error-message {
  background-color: var(--error-bg) !important;
  border-bottom-color: var(--error-border) !important;
  color: var(--error-text) !important;
}

/* 编辑器核心样式 */
:deep(.json-editor-core) {
  flex: 1;
  overflow: hidden;
  width: 100%;
}

/* 暗黑模式全覆盖：补充所有模式的样式 */
.json-editor-dark :deep(.jsoneditor) {
  background-color: var(--bg-primary) !important;
  color: var(--text-primary) !important;
  border: none !important;
  font-family: Monaco, Consolas, 'Courier New', monospace !important;
  font-size: 14px !important;
}

/* 覆盖所有子容器背景 */
.json-editor-dark :deep(.jsoneditor > *) {
  background-color: var(--bg-primary) !important;
}

/* 菜单样式 */
.json-editor-dark :deep(.jsoneditor-menu) {
  background-color: var(--bg-secondary) !important;
  border-bottom: 1px solid var(--border-primary) !important;
  padding: 4px !important;
}

.json-editor-dark :deep(.jsoneditor-menu button) {
  color: var(--text-primary) !important;
  border: none !important;
  border-radius: 4px !important;
  padding: 4px 8px !important;
  margin: 0 2px !important;
  transition: background-color 0.2s ease !important;
}

.json-editor-dark :deep(.jsoneditor-menu button:hover) {
  background-color: var(--border-secondary) !important;
}

/* 树状视图样式 */
.json-editor-dark :deep(.jsoneditor-tree) {
  background-color: var(--bg-primary) !important;
  color: var(--text-primary) !important;
  padding: 8px !important;
}

.json-editor-dark :deep(.jsoneditor-tree .jsoneditor-field) {
  color: #93c5fd !important;
  font-weight: 500 !important;
  margin-right: 4px !important;
}

.json-editor-dark :deep(.jsoneditor-tree .jsoneditor-string) {
  color: #a3e635 !important;
}

.json-editor-dark :deep(.jsoneditor-tree .jsoneditor-number) {
  color: #f87171 !important;
}

.json-editor-dark :deep(.jsoneditor-tree .jsoneditor-boolean) {
  color: #60a5fa !important;
}

.json-editor-dark :deep(.jsoneditor-tree .jsoneditor-null) {
  color: #94a3b8 !important;
}

/* 代码/文本模式样式（核心修复） */
.json-editor-dark :deep(.jsoneditor-code) {
  background-color: var(--bg-primary) !important;
  color: var(--text-primary) !important;
}

.json-editor-dark :deep(.jsoneditor-code textarea) {
  background-color: var(--bg-primary) !important;
  color: var(--text-primary) !important;
  border: none !important;
  font-family: Monaco, Consolas, 'Courier New', monospace !important;
  padding: 8px !important;
  resize: none !important;
}

.json-editor-dark :deep(.jsoneditor-code textarea:focus) {
  outline: none !important;
  border: 1px solid #60a5fa !important;
  border-radius: 2px !important;
}

/* 表单模式样式 */
.json-editor-dark :deep(.jsoneditor-form) {
  background-color: var(--bg-primary) !important;
  color: var(--text-primary) !important;
}

.json-editor-dark :deep(.jsoneditor-form input),
.json-editor-dark :deep(.jsoneditor-form textarea),
.json-editor-dark :deep(.jsoneditor-form select) {
  background-color: var(--bg-secondary) !important;
  color: var(--text-primary) !important;
  border: 1px solid var(--border-secondary) !important;
  border-radius: 4px !important;
}

/* 搜索框样式 */
.json-editor-dark :deep(.jsoneditor-search) {
  background-color: var(--bg-secondary) !important;
  border: 1px solid var(--border-primary) !important;
  color: var(--text-primary) !important;
  border-radius: 4px !important;
  padding: 4px 8px !important;
  margin: 0 4px !important;
}

.json-editor-dark :deep(.jsoneditor-search::placeholder) {
  color: var(--text-secondary) !important;
  opacity: 1 !important;
}

/* 输入框样式 */
.json-editor-dark :deep(.jsoneditor-text-input) {
  background-color: var(--bg-secondary) !important;
  color: var(--text-primary) !important;
  border: 1px solid var(--border-secondary) !important;
  border-radius: 2px !important;
  padding: 2px 4px !important;
}

.json-editor-dark :deep(.jsoneditor-text-input:focus) {
  border-color: #60a5fa !important;
  outline: none !important;
  background-color: #1a2436 !important;
}

/* 禁用状态 */
:deep(.jsoneditor-disabled) {
  opacity: 0.8 !important;
  cursor: not-allowed !important;
  background-color: #1a2436 !important;
}

/* 滚动条样式 */
.json-editor-dark :deep(.jsoneditor-tree::-webkit-scrollbar),
.json-editor-dark :deep(.jsoneditor-code::-webkit-scrollbar) {
  width: 8px !important;
  height: 8px !important;
}

.json-editor-dark :deep(.jsoneditor-tree::-webkit-scrollbar-track),
.json-editor-dark :deep(.jsoneditor-code::-webkit-scrollbar-track) {
  background: var(--bg-primary) !important;
}

.json-editor-dark :deep(.jsoneditor-tree::-webkit-scrollbar-thumb),
.json-editor-dark :deep(.jsoneditor-code::-webkit-scrollbar-thumb) {
  background: var(--border-secondary) !important;
  border-radius: 4px !important;
}

.json-editor-dark :deep(.jsoneditor-tree::-webkit-scrollbar-thumb:hover),
.json-editor-dark :deep(.jsoneditor-code::-webkit-scrollbar-thumb:hover) {
  background: #475569 !important;
}

/* 亮色模式基础样式 */
:deep(.jsoneditor) {
  font-family: Monaco, Consolas, 'Courier New', monospace !important;
  font-size: 14px !important;
}

:deep(.jsoneditor-menu) {
  padding: 4px !important;
}

:deep(.jsoneditor-menu button) {
  border-radius: 4px !important;
  padding: 4px 8px !important;
}
</style>
