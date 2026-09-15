<script setup lang="ts">
import { computed, ref, watch } from 'vue';

import { $t } from '@vben/locales';

import { Image as AImage, Spin as ASpin } from 'ant-design-vue';

import { compressImageFile, notifyUploadError, uploadMediaAsset } from '#/api';

const props = withDefaults(
  defineProps<{
    height?: number;
    value?: string;
    width?: number;
  }>(),
  { height: 80, value: '', width: 128 },
);

const emit = defineEmits<{
  (e: 'update:value', url: string): void;
}>();

const uploading = ref(false);
// 上传进度百分比（-1 表示未知/不可得）
const uploadPercent = ref(-1);
// 上传期间的本地即时预览（blob URL），成功后切换为远端地址
const localPreview = ref('');

const displayUrl = computed(() => localPreview.value || props.value || '');

watch(localPreview, (_val, old) => {
  if (old) {
    URL.revokeObjectURL(old);
  }
});

async function handleUpload(options: any) {
  const file = options?.file as File | undefined;
  if (!file) {
    return;
  }

  localPreview.value = URL.createObjectURL(file);
  uploading.value = true;
  uploadPercent.value = -1;
  try {
    const compressed = await compressImageFile(file);
    const resp = await uploadMediaAsset({}, compressed, (progressEvent: any) => {
      const total = progressEvent?.total ?? 0;
      uploadPercent.value =
        total > 0 ? Math.floor((progressEvent.loaded / total) * 100) : -1;
    });
    const url = (resp as { objectName?: string }).objectName || '';
    if (!url) {
      notifyUploadError(new Error('empty upload response'));
      return;
    }
    emit('update:value', url);
  } catch (error) {
    console.error('Upload image failed:', error);
    notifyUploadError(error);
  } finally {
    uploading.value = false;
    uploadPercent.value = -1;
    localPreview.value = '';
  }
}

function handleRemove() {
  emit('update:value', '');
}
</script>

<template>
  <div class="flex items-start gap-3">
    <!-- 预览卡片：有图可点击放大，无图显示上传占位 -->
    <div
      class="relative shrink-0 overflow-hidden rounded border border-border"
      :style="{ width: `${width}px`, height: `${height}px` }"
    >
      <AImage
        v-if="displayUrl"
        :src="displayUrl"
        class="h-full w-full object-cover"
        wrapper-class-name="h-full w-full"
        :root-class-name="'h-full w-full'"
      />
      <a-upload
        v-else
        accept="image/*"
        :show-upload-list="false"
        :custom-request="handleUpload"
      >
        <div
          class="flex h-full w-full flex-col items-center justify-center gap-1 border-dashed bg-accent/20 text-muted-foreground transition-colors hover:bg-accent/40"
          :style="{ width: `${width - 2}px`, height: `${height - 2}px` }"
        >
          <span class="text-xl leading-none">＋</span>
          <span class="text-xs">{{ $t('ui.button.upload') }}</span>
        </div>
      </a-upload>

      <div
        v-if="uploading"
        class="absolute inset-0 flex flex-col items-center justify-center gap-2 bg-black/40"
      >
        <ASpin size="small" />
        <span class="text-xs text-white">
          {{ $t('ui.notification.uploading') }}{{ uploadPercent >= 0 ? ` ${uploadPercent}%` : '' }}
        </span>
      </div>
    </div>

    <!-- 操作按钮仅在已有图片时出现；无图时占位卡本身即上传入口 -->
    <div v-if="value" class="flex flex-col gap-2">
      <a-upload
        accept="image/*"
        :show-upload-list="false"
        :custom-request="handleUpload"
      >
        <a-button size="small" :loading="uploading">
          {{ $t('page.post.button.replaceImage') }}
        </a-button>
      </a-upload>
      <a-button size="small" danger @click="handleRemove">
        {{ $t('page.post.button.removeImage') }}
      </a-button>
    </div>
  </div>
</template>
