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

import { useAppStore } from '@/stores';

import VueJsonEditor from 'json-editor-vue';

import 'jsoneditor/dist/jsoneditor.min.css';

// 类型定义
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

const appStore = useAppStore();

// 响应式数据
const localValue = ref(props.modelValue);
const jsonData = ref<any[] | null | Record<string, any>>(null);
const parseError = ref<string>('');
const isValidJson = ref(true);
const instance = getCurrentInstance();
let observer: MutationObserver | null = null;
let themeObserver: MutationObserver | null = null;

const isDark = ref(false);

const updateIsDark = () => {
  const prefersDark = appStore.mode === 'dark';
  if (typeof document === 'undefined') {
    isDark.value = prefersDark;
    return;
  }
  const root = document.documentElement;
  isDark.value =
    prefersDark ||
    root.classList.contains('dark') ||
    root.classList.contains('theme-dark') ||
    root.classList.contains('json-editor-dark');
};

// computed
// 验证并格式化 JSON
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

  // 🛡️ 确保 jsonData 是对象类型
  if (parsed !== null && typeof parsed === 'object') {
    jsonData.value = parsed;
  } else if (parsed === null) {
    jsonData.value = {};
  } else {
    // 兜底：非对象值包装处理
    jsonData.value = { value: parsed };
  }
};

// 监听外部值变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== localValue.value) {
      const { parsed, formatted } = validateAndFormat(newVal);
      localValue.value = formatted || newVal || props.placeholder;
      console.log('props.modelValue');
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

    if (typeof newVal === 'string') {
      if (newVal !== localValue.value) {
        localValue.value = newVal;
        emit('update:modelValue', newVal);
        emit('change', newVal);
      }
      return;
    }

    // 正常对象/数组：序列化为字符串
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
    const numericHeight = Number(props.height);
    if (!Number.isNaN(numericHeight)) {
      baseHeight = numericHeight;
    } else if (props.height.endsWith('px')) {
      const pxValue = Number(props.height.replace('px', ''));
      if (!Number.isNaN(pxValue)) {
        baseHeight = pxValue;
      }
    } else {
      // 百分比等非数值高度直接返回原字符串
      return props.height;
    }
  }

  const finalHeight = Math.max(baseHeight - 40, 200);
  return `${finalHeight}px`;
});

// 刷新编辑器样式
const refreshEditor = () => {
  nextTick(() => {
    const container = instance?.proxy?.$el as HTMLElement | undefined;
    if (!container) return;
    container.dataset.theme = isDark.value ? 'dark' : 'light';

    // 强制应用暗色模式样式
    if (isDark.value) {
      const jsonEditorEl = container.querySelector('.jsoneditor') as HTMLElement;
      if (jsonEditorEl) {
        jsonEditorEl.style.backgroundColor = 'var(--color-surface)';
        jsonEditorEl.style.color = 'var(--color-text-primary)';
      }

      const menuEl = container.querySelector('.jsoneditor-menu') as HTMLElement;
      if (menuEl) {
        menuEl.style.backgroundColor = '#2a3140';
        menuEl.style.borderBottomColor = 'var(--color-border)';
      }

      // ACE 编辑器
      const aceEditorEl = container.querySelector('.ace-jsoneditor') as HTMLElement;
      if (aceEditorEl) {
        aceEditorEl.style.backgroundColor = 'var(--color-surface)';
      }

      const aceScroller = container.querySelector('.ace_scroller') as HTMLElement;
      if (aceScroller) {
        aceScroller.style.backgroundColor = 'var(--color-surface)';
      }

      // 所有内容区域
      const allDivs = container.querySelectorAll('.jsoneditor > div, .jsoneditor-outer, .jsoneditor-tree, .jsoneditor-code, .jsoneditor-text');
      allDivs.forEach((el: Element) => {
        (el as HTMLElement).style.backgroundColor = 'var(--color-surface)';
        (el as HTMLElement).style.color = 'var(--color-text-primary)';
      });
    }
  });
};

// 监听主题变化
watch(
  () => appStore.mode,
  () => {
    updateIsDark();
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
const handleEditorChange = (value: any) => {
  if (typeof value === 'string') {
    const rawValue = value;
    localValue.value = rawValue;
    emit('update:modelValue', rawValue);
    emit('change', rawValue);

    const { parsed } = validateAndFormat(rawValue);
    if (parsed !== null && typeof parsed === 'object') {
      jsonData.value = parsed;
    }

    // 刷新样式
    setTimeout(() => refreshEditor(), 50);
    return;
  }

  if (Array.isArray(value) || (value !== null && typeof value === 'object')) {
    setTimeout(() => refreshEditor(), 50);
    return;
  }

  jsonData.value = { value };
  setTimeout(() => refreshEditor(), 50);
};

// 初始化和销毁逻辑
onMounted(() => {
  updateIsDark();
  initData();
  nextTick(() => {
    emit('ready');
    refreshEditor();

    if (typeof document !== 'undefined') {
      const root = document.documentElement;
      themeObserver = new MutationObserver(() => {
        updateIsDark();
        refreshEditor();
      });
      themeObserver.observe(root, {
        attributes: true,
        attributeFilter: ['class'],
      });
    }

    const container = instance?.proxy?.$el as HTMLElement | undefined;
    if (!container) return;
    const editorEl = container.querySelector('.json-editor-core');

    if (editorEl) {
      observer = new MutationObserver((mutations) => {
        const hasStyleChange = mutations.some(
          (m) =>
            m.type === 'attributes' &&
            ['class', 'style'].includes(m.attributeName || ''),
        );
        if (isDark.value && hasStyleChange) {
          refreshEditor();
        }

        // 如果是暗色模式，持续强制应用样式
        if (isDark.value) {
          setTimeout(() => refreshEditor(), 100);
        }
      });

      observer.observe(editorEl, {
        childList: true,
        subtree: true,
        attributes: true,
        attributeFilter: ['class', 'style'],
      });

      // 额外的定时器确保样式应用
      if (isDark.value) {
        const styleInterval = setInterval(() => {
          if (isDark.value) {
            refreshEditor();
          } else {
            clearInterval(styleInterval);
          }
        }, 500);

        // 10秒后停止定时器
        setTimeout(() => clearInterval(styleInterval), 10000);
      }
    }
  });
});

onUnmounted(() => {
  if (themeObserver) {
    themeObserver.disconnect();
    themeObserver = null;
  }
  if (observer) {
    observer.disconnect();
    observer = null;
  }
});
</script>

<template>
  <div class="json-editor-container" :class="{ 'json-editor-dark': isDark }">
    <!-- 错误提示 -->
    <div v-if="parseError" class="error-message">
      {{ parseError }}
    </div>

    <VueJsonEditor
      v-model="jsonData"
      :mode="(options.mode as any)"
      :disabled="disabled"
      :search="options.search"
      :placeholder="placeholder"
      :style="{ height: editorHeight, width: '100%' }"
      class="json-editor-core"
      @change="handleEditorChange"
    />
  </div>
</template>

<style scoped lang="less">
.json-editor-container {
  display: flex;
  flex-direction: column;
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

.json-editor-dark {
  --bg-primary: var(--color-surface);
  --bg-secondary: #2a3140;
  --text-primary: var(--color-text-primary);
  --text-secondary: var(--color-text-secondary);
  --border-primary: var(--color-border);
  --border-secondary: #3a4a60;
  --error-bg: rgba(208, 48, 80, 0.1);
  --error-text: #d03050;
  --error-border: #d03050;

  background-color: var(--bg-primary);
  border-color: var(--border-primary);
}

.error-message {
  padding: 8px 12px;
  margin: 0;
  font-size: 13px;
  line-height: 1.5;
  color: #d03050;
  background-color: rgba(208, 48, 80, 0.1);
  border-bottom: 1px solid rgba(208, 48, 80, 0.2);
}

.json-editor-dark .error-message {
  color: var(--error-text);
  background-color: var(--error-bg);
  border-bottom-color: var(--error-border);
}

.json-editor-container {
  :deep(.json-editor-core) {
    flex: 1;
    width: 100%;
    overflow: hidden;
  }

  // 基础样式
  :deep(.jsoneditor) {
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace !important;
    font-size: 14px !important;
    color: var(--editor-text, var(--color-text-primary)) !important;
    background-color: var(--editor-bg, var(--color-surface)) !important;
    border: none !important;
  }

  :deep(.jsoneditor > *) {
    background-color: var(--editor-bg, var(--color-surface)) !important;
  }

  :deep(.jsoneditor .ace-jsoneditor) {
    background-color: var(--editor-bg, var(--color-surface)) !important;
  }

  :deep(.jsoneditor .ace-jsoneditor .ace_gutter) {
    background-color: var(--color-surface) !important;
    color: var(--color-text-secondary) !important;
  }

  :deep(.jsoneditor .ace-jsoneditor .ace_scroller) {
    background-color: var(--editor-bg, var(--color-surface)) !important;
  }

  :deep(.jsoneditor .ace-jsoneditor .ace_content) {
    background-color: var(--editor-bg, var(--color-surface)) !important;
  }

  // 菜单样式
  :deep(.jsoneditor-menu) {
    padding: 4px !important;
    background-color: var(--color-surface) !important;
    border-bottom: 1px solid var(--editor-border, var(--color-border)) !important;
  }

  :deep(.jsoneditor-menu button) {
    padding: 4px 8px !important;
    margin: 0 2px !important;
    color: var(--editor-text, var(--color-text-primary)) !important;
    border: none !important;
    border-radius: 4px !important;
    transition: all 0.2s ease !important;

    &:hover {
      background-color: rgba(102, 126, 234, 0.08) !important;
      color: var(--editor-focus-border, var(--color-brand)) !important;
    }
  }

  // 树状视图
  :deep(.jsoneditor-tree) {
    padding: 8px !important;
    color: var(--editor-text, var(--color-text-primary)) !important;
    background-color: var(--editor-bg, var(--color-surface)) !important;
  }

  :deep(.jsoneditor-tree .jsoneditor-field) {
    margin-right: 4px !important;
    font-weight: 500 !important;
    color: var(--color-brand) !important;
  }

  :deep(.jsoneditor-tree .jsoneditor-string) {
    color: #10b981 !important;
  }

  :deep(.jsoneditor-tree .jsoneditor-number) {
    color: #f59e0b !important;
  }

  :deep(.jsoneditor-tree .jsoneditor-boolean) {
    color: var(--color-brand) !important;
  }

  :deep(.jsoneditor-tree .jsoneditor-null) {
    color: var(--color-text-secondary) !important;
  }

  // 代码/文本模式
  :deep(.jsoneditor-code) {
    color: var(--editor-text, var(--color-text-primary)) !important;
    background-color: var(--editor-bg, var(--color-surface)) !important;
  }

  :deep(.jsoneditor-code textarea) {
    padding: 8px !important;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace !important;
    color: var(--editor-text, var(--color-text-primary)) !important;
    background-color: var(--editor-bg, var(--color-surface)) !important;
    border: none !important;
    resize: none !important;

    &:focus {
      outline: none !important;
      border: 1px solid var(--editor-focus-border, var(--color-brand)) !important;
    }
  }

  // 表单模式
  :deep(.jsoneditor-form) {
    color: var(--editor-text, var(--color-text-primary)) !important;
    background-color: var(--editor-bg, var(--color-surface)) !important;
  }

  :deep(.jsoneditor-form input),
  :deep(.jsoneditor-form textarea),
  :deep(.jsoneditor-form select) {
    color: var(--editor-text, var(--color-text-primary)) !important;
    background-color: var(--color-surface) !important;
    border: 1px solid var(--editor-border, var(--color-border)) !important;
    border-radius: 4px !important;
    transition: all 0.2s ease !important;

    &:focus {
      border-color: var(--editor-focus-border, var(--color-brand)) !important;
      box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1) !important;
    }
  }

  // 搜索框
  :deep(.jsoneditor-search) {
    padding: 6px 8px !important;
    margin: 0 4px !important;
    color: var(--editor-text, var(--color-text-primary)) !important;
    background-color: var(--color-surface) !important;
    border: 1px solid var(--editor-border, var(--color-border)) !important;
    border-radius: 4px !important;
    transition: all 0.2s ease !important;

    &::placeholder {
      color: var(--editor-placeholder, var(--color-text-secondary)) !important;
      opacity: 1 !important;
    }

    &:focus {
      border-color: var(--editor-focus-border, var(--color-brand)) !important;
      box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1) !important;
    }
  }

  // 滚动条
  :deep(.jsoneditor-tree::-webkit-scrollbar),
  :deep(.jsoneditor-code::-webkit-scrollbar) {
    width: 8px;
    height: 8px;
  }

  :deep(.jsoneditor-tree::-webkit-scrollbar-track),
  :deep(.jsoneditor-code::-webkit-scrollbar-track) {
    background: transparent;
  }

  :deep(.jsoneditor-tree::-webkit-scrollbar-thumb),
  :deep(.jsoneditor-code::-webkit-scrollbar-thumb) {
    background: rgba(102, 126, 234, 0.3);
    border-radius: 4px;

    &:hover {
      background: rgba(102, 126, 234, 0.6);
    }
  }
}

// 暗色模式 - 使用更强的选择器
html.dark {
  .json-editor-container {
    background-color: var(--color-surface) !important;
    border-color: var(--color-border) !important;

    // 所有 jsoneditor 相关元素
    :deep(.jsoneditor),
    :deep(.jsoneditor[class*='jsoneditor']),
    :deep(div.jsoneditor),
    :deep(div[class*='jsoneditor']) {
      background-color: var(--color-surface) !important;
      color: var(--color-text-primary) !important;
    }

    :deep(.jsoneditor *),
    :deep(.jsoneditor > *),
    :deep(.jsoneditor-outer) {
      background-color: var(--color-surface) !important;
    }

    // ACE 编辑器
    :deep(.ace-jsoneditor),
    :deep(div.ace-jsoneditor) {
      background-color: var(--color-surface) !important;
    }

    :deep(.ace-jsoneditor .ace_gutter),
    :deep(.ace_gutter) {
      background-color: #2a3140 !important;
      color: var(--color-text-secondary) !important;
    }

    :deep(.ace-jsoneditor .ace_scroller),
    :deep(.ace_scroller) {
      background-color: var(--color-surface) !important;
    }

    :deep(.ace-jsoneditor .ace_content),
    :deep(.ace_content) {
      background-color: var(--color-surface) !important;
    }

    :deep(.ace_editor),
    :deep(.ace_text-layer) {
      background-color: var(--color-surface) !important;
    }

    // 菜单/工具栏 - 使用多个选择器确保覆盖
    :deep(.jsoneditor-menu),
    :deep(div.jsoneditor-menu),
    :deep(.jsoneditor > .jsoneditor-menu) {
      background-color: #2a3140 !important;
      border-bottom-color: var(--color-border) !important;
    }

    :deep(.jsoneditor-menu button),
    :deep(.jsoneditor-menu > button),
    :deep(div.jsoneditor-menu button) {
      color: var(--color-text-primary) !important;
      background-color: transparent !important;

      &:hover {
        background-color: rgba(129, 147, 255, 0.12) !important;
        color: var(--color-brand) !important;
      }

      &:active,
      &:focus {
        background-color: rgba(129, 147, 255, 0.18) !important;
      }
    }

    // 树状视图
    :deep(.jsoneditor-tree),
    :deep(div.jsoneditor-tree) {
      background-color: var(--color-surface) !important;
      color: var(--color-text-primary) !important;
    }

    // 代码模式
    :deep(.jsoneditor-code),
    :deep(div.jsoneditor-code) {
      background-color: var(--color-surface) !important;
      color: var(--color-text-primary) !important;
    }

    :deep(.jsoneditor-code textarea),
    :deep(.jsoneditor-text),
    :deep(textarea.jsoneditor-text) {
      background-color: var(--color-surface) !important;
      color: var(--color-text-primary) !important;
      border-color: var(--color-border) !important;
    }

    // 表单模式
    :deep(.jsoneditor-form) {
      background-color: var(--color-surface) !important;
      color: var(--color-text-primary) !important;
    }

    :deep(.jsoneditor-form input),
    :deep(.jsoneditor-form textarea),
    :deep(.jsoneditor-form select) {
      background-color: #2a3140 !important;
      color: var(--color-text-primary) !important;
      border-color: var(--color-border) !important;
    }

    // 搜索框
    :deep(.jsoneditor-search),
    :deep(input.jsoneditor-search) {
      background-color: #2a3140 !important;
      color: var(--color-text-primary) !important;
      border-color: var(--color-border) !important;
    }

    // 视图模式
    :deep(.jsoneditor-view) {
      background-color: var(--color-surface) !important;
      color: var(--color-text-primary) !important;
    }

    // 所有面板和容器
    :deep(.jsoneditor-panel),
    :deep(.jsoneditor-statusbar),
    :deep(.jsoneditor-navigation-bar) {
      background-color: #2a3140 !important;
      border-color: var(--color-border) !important;
    }

    // 滚动条
    :deep(.jsoneditor-tree::-webkit-scrollbar-thumb),
    :deep(.jsoneditor-code::-webkit-scrollbar-thumb) {
      background: rgba(129, 147, 255, 0.3) !important;

      &:hover {
        background: rgba(129, 147, 255, 0.6) !important;
      }
    }
  }

  // 针对特定 class 的覆盖
  .json-editor-dark {
    :deep(.jsoneditor),
    :deep(.jsoneditor *) {
      background-color: var(--color-surface) !important;
    }

    :deep(.jsoneditor-menu) {
      background-color: #2a3140 !important;
    }
  }
}

// 全局样式覆盖（最后保底）
:deep(.json-editor-dark .jsoneditor) {
  background-color: var(--color-surface) !important;
}

:deep(.json-editor-dark .jsoneditor-menu) {
  background-color: #2a3140 !important;
}
</style>
