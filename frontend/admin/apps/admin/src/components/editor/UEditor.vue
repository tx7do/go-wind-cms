<script setup lang="ts">
import { computed, ref, unref, watch } from 'vue';

import { useAccessStore } from '@vben/stores';

import {
  formatSaveShowContent,
  full_toolbars,
  getCurrentLang,
  isDarkMode,
} from '#/adapter/component/UEditor/src/utils';

interface Props {
  modelValue: string;
  contentHeight?: number;
  height?: number | string;
  disabled?: boolean;
  editorId?: string;
  maximumWords?: number;
  placeholder?: string;
  toolbars?: string[][];
  config?: Record<string, any>;
}

const props = withDefaults(defineProps<Props>(), {
  contentHeight: 320,
  height: undefined,
  disabled: false,
  editorId: 'editor',
  maximumWords: 10_000,
  placeholder: 'Please enter content...',
  toolbars: () => full_toolbars,
  config: () => ({}),
});

const emit = defineEmits<{
  (e: 'change', value: string): void;
  (e: 'ready'): void;
  (e: 'update:modelValue', value: string): void;
}>();

const accessStore = useAccessStore();

const baseUrl = import.meta.env.VITE_GLOB_API_URL;
const contentRef = ref(formatSaveShowContent('show', props.modelValue));
const editorInst = ref<any>(null);
const initedRef = ref(false);

const ueditorHomeUrl = ref('/static/UEditorPlus/');
const getUeditorConfigUrl = ref(`${baseUrl}/admin/v1/ueditor`);

const getToken = computed(() => accessStore.accessToken);
const getServerHeaders = computed(() => ({
  Authorization: `Bearer ${getToken.value}`,
}));
const getImageDownloadUrl = computed(
  () => import.meta.env.VITE_GLOB_DOWNLOAD_URL,
);

const theme = computed(() => (isDarkMode() ? 'dark' : 'default'));

const wrapperRef = ref<HTMLDivElement | null>(null);
const observedHeight = ref<null | number>(null);
let resizeObserver: null | ResizeObserver = null;

const getViewportMaxHeight = (el: HTMLDivElement) => {
  const rect = el.getBoundingClientRect();
  return Math.max(0, Math.floor(window.innerHeight - rect.top));
};

const resolvedContentHeight = computed(() => {
  if (observedHeight.value && observedHeight.value > 0) {
    return observedHeight.value;
  }
  if (typeof props.contentHeight === 'number') {
    return props.contentHeight;
  }
  if (typeof props.height === 'number') {
    return props.height;
  }
  return 320;
});

const editorConfig = computed(() => ({
  serverUrl: getUeditorConfigUrl.value,
  serverHeaders: getServerHeaders.value,
  UEDITOR_HOME_URL: ueditorHomeUrl.value,
  UEDITOR_CORS_URL: ueditorHomeUrl.value,
  initialFrameWidth: null,
  initialFrameHeight: resolvedContentHeight.value,
  catchRemoteImageEnable: false,
  autoHeightEnabled: false,
  autoSyncData: false,
  autoFloatEnabled: true,
  theme: theme.value,
  themePath: '/static/UEditorPlus/themes/',
  toolbars: props.toolbars,
  maximumWords: props.maximumWords,
  readonly: props.disabled,
  placeholder: props.placeholder,
  imageActionName: 'uploadImage',
  imageFieldName: 'file',
  imageMaxSize: 2_048_000,
  imageAllowFiles: ['.png', '.jpg', '.jpeg', '.gif', '.bmp'],
  imageCompressEnable: true,
  imageCompressBorder: 1600,
  imageInsertAlign: 'none',
  imageUrlPrefix: getImageDownloadUrl.value,
  scrawlActionName: 'uploadScrawl',
  scrawlFieldName: 'file',
  scrawlMaxSize: 2_048_000,
  scrawlUrlPrefix: getImageDownloadUrl.value,
  scrawlInsertAlign: 'none',
  snapscreenActionName: 'uploadImage',
  snapscreenUrlPrefix: getImageDownloadUrl.value,
  snapscreenInsertAlign: 'none',
  catcherLocalDomain: ['127.0.0.1', 'localhost', 'img.baidu.com'],
  catcherActionName: 'catchImage',
  catcherFieldName: 'file',
  catcherUrlPrefix: getImageDownloadUrl.value,
  catcherMaxSize: 2_048_000,
  catcherAllowFiles: ['.png', '.jpg', '.jpeg', '.gif', '.bmp'],
  videoActionName: 'uploadVideo',
  videoFieldName: 'file',
  videoUrlPrefix: getImageDownloadUrl.value,
  videoMaxSize: 102_400_000,
  videoAllowFiles: [
    '.flv',
    '.swf',
    '.mkv',
    '.avi',
    '.rm',
    '.rmvb',
    '.mpeg',
    '.mpg',
    '.ogg',
    '.ogv',
    '.mov',
    '.wmv',
    '.mp4',
    '.webm',
    '.mp3',
    '.wav',
    '.mid',
  ],
  fileActionName: 'uploadFile',
  fileFieldName: 'file',
  fileUrlPrefix: getImageDownloadUrl.value,
  fileMaxSize: 51_200_000,
  fileAllowFiles: [
    '.png',
    '.jpg',
    '.jpeg',
    '.gif',
    '.bmp',
    '.flv',
    '.swf',
    '.mkv',
    '.avi',
    '.rm',
    '.rmvb',
    '.mpeg',
    '.mpg',
    '.ogg',
    '.ogv',
    '.mov',
    '.wmv',
    '.mp4',
    '.webm',
    '.mp3',
    '.wav',
    '.mid',
    '.rar',
    '.zip',
    '.tar',
    '.gz',
    '.7z',
    '.bz2',
    '.cab',
    '.iso',
    '.doc',
    '.docx',
    '.xls',
    '.xlsx',
    '.ppt',
    '.pptx',
    '.pdf',
    '.txt',
    '.md',
    '.xml',
  ],
  imageManagerActionName: 'listImage',
  imageManagerListSize: 20,
  imageManagerUrlPrefix: getImageDownloadUrl.value,
  imageManagerInsertAlign: 'none',
  imageManagerAllowFiles: ['.png', '.jpg', '.jpeg', '.gif', '.bmp'],
  fileManagerActionName: 'listFile',
  fileManagerUrlPrefix: getImageDownloadUrl.value,
  fileManagerListSize: 20,
  fileManagerAllowFiles: [
    '.png',
    '.jpg',
    '.jpeg',
    '.gif',
    '.bmp',
    '.flv',
    '.swf',
    '.mkv',
    '.avi',
    '.rm',
    '.rmvb',
    '.mpeg',
    '.mpg',
    '.ogg',
    '.ogv',
    '.mov',
    '.wmv',
    '.mp4',
    '.webm',
    '.mp3',
    '.wav',
    '.mid',
    '.rar',
    '.zip',
    '.tar',
    '.gz',
    '.7z',
    '.bz2',
    '.cab',
    '.iso',
    '.doc',
    '.docx',
    '.xls',
    '.xlsx',
    '.ppt',
    '.pptx',
    '.pdf',
    '.txt',
    '.md',
    '.xml',
  ],
  ...props.config,
}));

const ready = (editorInstance: any) => {
  editorInst.value = editorInstance;
  initedRef.value = true;
  emit('ready');
};

watch(
  () => wrapperRef.value,
  (el, _, onCleanup) => {
    if (!el) {
      return;
    }

    const syncHeight = () => {
      const maxHeight = getViewportMaxHeight(el);
      const baseHeight = Math.floor(el.getBoundingClientRect().height);
      const nextHeight = Math.min(baseHeight, maxHeight);
      if (nextHeight > 0 && nextHeight !== observedHeight.value) {
        observedHeight.value = nextHeight;
        if (editorInst.value) {
          editorInst.value.setOpt({ initialFrameHeight: nextHeight });
        }
      }
    };

    resizeObserver = new ResizeObserver((entries) => {
      const entry = entries[0];
      if (!entry) {
        return;
      }
      const maxHeight = getViewportMaxHeight(el);
      const nextHeight = Math.floor(
        Math.min(entry.contentRect.height, maxHeight),
      );
      if (nextHeight > 0 && nextHeight !== observedHeight.value) {
        observedHeight.value = nextHeight;
        if (editorInst.value) {
          editorInst.value.setOpt({ initialFrameHeight: nextHeight });
        }
      }
    });

    resizeObserver.observe(el);
    window.addEventListener('resize', syncHeight);
    syncHeight();

    onCleanup(() => {
      resizeObserver?.disconnect();
      resizeObserver = null;
      window.removeEventListener('resize', syncHeight);
    });
  },
  { immediate: true },
);

watch(
  () => props.modelValue,
  (newVal) => {
    const formatted = formatSaveShowContent('show', newVal);
    if (formatted !== contentRef.value) {
      contentRef.value = formatted;
    }
  },
);

watch(contentRef, () => {
  const value = formatSaveShowContent('save', contentRef.value);
  emit('update:modelValue', value);
  emit('change', value);
});

watch(
  [() => isDarkMode(), () => initedRef.value],
  ([_val, initialed]) => {
    if (!initialed || !editorInst.value) {
      return;
    }
    editorInst.value.setOpt({ theme: theme.value });
  },
  { immediate: true },
);

watch(
  [() => getCurrentLang, () => initedRef.value],
  ([_val, initialed]) => {
    if (!initialed || !editorInst.value) {
      return;
    }
    editorInst.value.setOpt({ lang: unref(getCurrentLang) });
  },
  { immediate: true },
);
</script>

<template>
  <div ref="wrapperRef" class="ueditor-wrapper">
    <VueUeditorWrap
      :editor-id="editorId"
      v-model="contentRef"
      :config="editorConfig"
      :editor-dependencies="['ueditor.config.js', 'ueditor.all.js']"
      @ready="ready"
    />
  </div>
</template>

<style scoped>
.ueditor-wrapper {
  width: 100%;
  height: 100%;
  max-height: 100%;
  overflow: hidden;
}
</style>
