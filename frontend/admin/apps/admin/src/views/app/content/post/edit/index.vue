<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue';
import { onBeforeRouteLeave, useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { useTabs } from '@vben/hooks';
import {
  LucideArrowLeft,
  LucideChevronRight,
  LucideSettings,
  LucideSparkles,
} from '@vben/icons';
import { $t } from '@vben/locales';

import { Modal, notification } from 'ant-design-vue';

import { Editor } from '#/adapter/component/Editor';
import {
  apiClient,
  compressImageFile,
  editorTypeOptions,
  postStatusToColor,
  postStatusToName,
  notifyUploadError,
  uploadMediaAsset,
} from '#/api';
import { router } from '#/router';

import { usePostEditViewStore } from './post-edit-view.state';
import PostSettingsPanel from './post-settings-panel.vue';

const postEditViewStore = usePostEditViewStore();

const route = useRoute();
const { closeCurrentTab } = useTabs();

// 右侧设置栏默认展开；折叠后编辑器自动占满剩余宽度
const settingsOpen = ref(true);

/** 本地草稿恢复提示每次进入页面只弹一次 */
let draftPromptShown = false;

/**
 * 存在未发布的本地草稿时询问用户是否恢复（不再静默覆盖服务端数据）
 */
function maybePromptPendingDraft() {
  if (draftPromptShown || !postEditViewStore.pendingDraft) {
    return;
  }
  draftPromptShown = true;
  Modal.confirm({
    title: $t('page.post.validation.draftFound'),
    content: $t('page.post.validation.draftFoundMessage'),
    okText: $t('page.post.validation.draftRestore'),
    cancelText: $t('ui.button.cancel'),
    onOk: () => {
      postEditViewStore.restorePendingDraft();
    },
  });
}

/**
 * 未保存修改防抖自动暂存到本地草稿，意外离开后下次进入可恢复
 */
let draftSaveTimer: ReturnType<typeof setTimeout> | undefined;
watch(
  () => postEditViewStore.formData,
  () => {
    if (!postEditViewStore.hasUnsavedChanges()) {
      return;
    }
    if (draftSaveTimer) {
      clearTimeout(draftSaveTimer);
    }
    draftSaveTimer = setTimeout(() => {
      postEditViewStore.savePostDraft();
    }, 1500);
  },
  { deep: true },
);

/**
 * 离开编辑页前拦截未保存修改（关闭标签页/切换路由）
 */
onBeforeRouteLeave((_to, _from, next) => {
  if (!postEditViewStore.hasUnsavedChanges()) {
    next();
    return;
  }
  Modal.confirm({
    title: $t('page.post.validation.unsavedTitle'),
    content: $t('page.post.validation.unsavedMessage'),
    okText: $t('ui.button.ok'),
    cancelText: $t('ui.button.cancel'),
    onOk: () => next(true),
    onCancel: () => next(false),
  });
});

function handleBeforeUnload(event: BeforeUnloadEvent) {
  if (postEditViewStore.hasUnsavedChanges()) {
    event.preventDefault();
    event.returnValue = '';
  }
}
window.addEventListener('beforeunload', handleBeforeUnload);
onUnmounted(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload);
  if (draftSaveTimer) {
    clearTimeout(draftSaveTimer);
  }
});

const initLanguage = computed(() => {
  return (route.query.lang as string) || 'zh-CN';
});

const isCreateMode = computed(() => {
  return route.name === 'CreatePost';
});

const isEditMode = computed(() => {
  return route.name === 'EditPost';
});

const postId = computed(() => {
  if (isCreateMode.value) {
    return null;
  }
  const id = route.params.id ?? -1;
  return Number(id);
});

/**
 * 监听路由查询参数变化
 * 当用户通过URL直接访问时（例如打开书签），自动更新表单语言
 */
watch(
  () => route.query.lang,
  async (newLang) => {
    if (newLang && postEditViewStore.formData.lang !== newLang) {
      postEditViewStore.formData.lang = newLang as string;
    }
  },
);

/**
 * 处理返回按钮点击
 */
function goBack() {
  closeCurrentTab();
  router.push('/content/posts');
}

/**
 * Handle language switch
 */
async function handleLanguageChange(newLang: string) {
  // Update URL query parameter
  await router.replace({
    path: route.path,
    query: { ...route.query, lang: newLang },
  });

  // Use store's switchLanguage method to handle language change with draft loading
  await postEditViewStore.switchLanguage(newLang);

  // Show notification if translation doesn't exist
  if (postEditViewStore.needTranslate) {
    notification.info({
      message: $t('page.post.validation.translationNotExists'),
    });
  }
}

/**
 * 处理一键翻译（当前已有内容时先确认覆盖）
 */
async function handleTranslate() {
  if (postEditViewStore.formData.title || postEditViewStore.formData.content) {
    const confirmed = await new Promise<boolean>((resolve) => {
      Modal.confirm({
        title: $t('page.post.validation.translationExists'),
        okText: $t('ui.button.ok'),
        cancelText: $t('ui.button.cancel'),
        onOk: () => resolve(true),
        onCancel: () => resolve(false),
      });
    });
    if (!confirmed) {
      return;
    }
  }

  try {
    const titleResp = await apiClient.translatorService.Translate({
      sourceLanguage: 'auto',
      targetLanguage: postEditViewStore.formData.lang,
      content: postEditViewStore.formData.title,
    });
    postEditViewStore.formData.title =
      titleResp.translatedContent || postEditViewStore.formData.title;
  } catch (error) {
    console.error('Title translation failed:', error);
    notification.error({
      message: $t('page.post.validation.translateTitleFailed'),
    });
    return;
  }

  try {
    const contentResp = await apiClient.translatorService.Translate({
      sourceLanguage: 'auto',
      targetLanguage: postEditViewStore.formData.lang,
      content: postEditViewStore.formData.content,
    });
    postEditViewStore.formData.content =
      contentResp.translatedContent || postEditViewStore.formData.content;
  } catch (error) {
    console.error('Content translation failed:', error);
    notification.error({
      message: $t('page.post.validation.translateContentFailed'),
    });
  }

  // 摘要非空时一并翻译（失败不影响整体流程）
  if (postEditViewStore.formData.summary) {
    try {
      const summaryResp = await apiClient.translatorService.Translate({
        sourceLanguage: 'auto',
        targetLanguage: postEditViewStore.formData.lang,
        content: postEditViewStore.formData.summary,
      });
      postEditViewStore.formData.summary =
        summaryResp.translatedContent || postEditViewStore.formData.summary;
    } catch (error) {
      console.error('Summary translation failed:', error);
    }
  }
}

/**
 * 处理保存草稿（保存为服务端草稿状态）
 */
async function handleSaveDraft() {
  try {
    const resp = await postEditViewStore.saveServerDraft();
    if (resp && resp !== '') {
      notification.error({
        message: resp,
      });
      return;
    }

    notification.success({
      message: $t('page.post.validation.saveDraftSuccess'),
    });

    // 新建模式下保存草稿后已生成文章，切换到编辑路由，后续保存走更新
    if (isCreateMode.value && postEditViewStore.postId) {
      await router.replace({
        name: 'EditPost',
        params: { id: String(postEditViewStore.postId) },
        query: { ...route.query },
      });
    }
  } catch (error) {
    console.error('Save draft failed:', error);
    notification.error({
      message: $t('page.post.validation.saveDraftFailed'),
    });
  }
}

/**
 * 处理发布文章（设置了未来发布时间时会转为定时发布）
 */
async function handlePublish() {
  const resp = await postEditViewStore.publishPost();
  if (!resp || resp === '') {
    notification.success({
      message:
        postEditViewStore.formData.status === 'POST_STATUS_SCHEDULED'
          ? $t('page.post.validation.scheduleSuccess')
          : $t('page.post.validation.publishSuccess'),
    });

    goBack();
  } else {
    notification.error({
      message: resp,
    });
  }
}

/**
 * 处理图片上传（编辑器正文插图：先压缩再带进度上传）
 */
async function handleUploadImage(
  file: File,
  onProgress?: (percent: number) => void,
): Promise<string> {
  try {
    const compressed = await compressImageFile(file);
    const resp = await uploadMediaAsset(
      {},
      compressed,
      onProgress
        ? (progressEvent: any) => {
            const total = progressEvent?.total ?? 0;
            onProgress(
              total > 0
                ? Math.floor((progressEvent.loaded / total) * 100)
                : -1,
            );
          }
        : undefined,
    );
    return (resp as { objectName?: string }).objectName || '';
  } catch (error) {
    console.error('Image upload failed:', error);
    notifyUploadError(error);
    return '';
  }
}

/**
 * 加载文章数据（仅编辑模式）
 */
async function loadPost() {
  if (!isEditMode.value) {
    return;
  }

  try {
    await postEditViewStore.fetchPost();

    if (postEditViewStore.needTranslate) {
      notification.info({
        message: $t('page.post.validation.translationNotExists'),
      });
    }
  } catch (error) {
    console.error('Failed to load post:', error);
    notification.error({
      message: $t('page.post.validation.loadFailed'),
    });
    throw error;
  }
}

/**
 * 初始化页面数据
 */
async function init() {
  postEditViewStore.fetchCategoryOptions();

  try {
    await postEditViewStore.fetchLanguageList();
  } catch {
    notification.error({
      message: $t('page.post.validation.loadLanguageFailed'),
    });
  }

  if (isCreateMode.value) {
    postEditViewStore.initCreateMode(initLanguage.value);
  } else if (isEditMode.value) {
    postEditViewStore.initEditMode(postId.value || 0, initLanguage.value);
    await loadPost();
  } else {
    console.error('Unknown route name:', route.name);
  }

  maybePromptPendingDraft();
}

init();
</script>

<template>
  <Page
    auto-content-height
    content-class="flex h-full min-h-0 flex-col p-0 overflow-hidden"
  >
    <template #title>
      <div class="flex w-full items-center gap-2">
        <a-button type="text" @click="goBack">
          <template #icon>
            <LucideArrowLeft class="text-align:center" />
          </template>
        </a-button>
        <a-input
          v-model:value="postEditViewStore.formData.title"
          :placeholder="$t('page.post.placeholder.title')"
          size="large"
          class="flex-1"
        />
        <a-tag
          v-if="!isCreateMode && postEditViewStore.formData.status"
          :color="postStatusToColor(postEditViewStore.formData.status)"
        >
          {{ postStatusToName(postEditViewStore.formData.status) }}
        </a-tag>
        <a-select
          :value="postEditViewStore.formData.lang"
          style="width: 200px"
          @change="handleLanguageChange"
        >
          <a-select-option
            v-for="option in postEditViewStore.languageOptions"
            :key="option.value"
            :value="option.value"
          >
            <span>
              {{ option.label }}
              <span
                v-if="option.hasTranslation"
                class="ml-2 text-green-600"
                :title="$t('page.post.placeholder.hasTranslation')"
              >
                ✓
              </span>
              <span
                v-else
                class="ml-2 text-orange-500"
                :title="$t('page.post.placeholder.noTranslation')"
              >
                ○
              </span>
            </span>
          </a-select-option>
        </a-select>
        <a-button
          type="primary"
          class="translate-btn"
          @click="handleTranslate"
        >
          <template #icon>
            <LucideSparkles />
          </template>
          {{ $t('page.post.button.oneClickTranslate') }}
        </a-button>
        <a-select
          v-model:value="postEditViewStore.formData.editorType"
          style="width: 200px"
        >
          <a-select-option
            v-for="option in editorTypeOptions"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </a-select-option>
        </a-select>
      </div>
    </template>

    <!-- 主体：左侧编辑列（弹性伸缩） + 右侧内嵌设置栏 -->
    <div class="flex min-h-0 flex-1 overflow-hidden">
      <div class="flex min-h-0 min-w-0 flex-1 flex-col">
        <div class="post-edit-container min-h-0 flex-1">
          <Editor
            class="absolute inset-0"
            height="100%"
            v-model="postEditViewStore.formData.content"
            :editor-type="postEditViewStore.formData.editorType"
            :placeholder="$t('page.post.placeholder.content')"
            :upload-image="handleUploadImage"
          />
        </div>
      </div>

      <!-- 右侧：内嵌设置栏 + 常驻竖向开关轨 -->
      <div class="flex h-full shrink-0">
        <aside
          class="h-full overflow-hidden transition-all duration-300"
          :class="settingsOpen ? 'w-[340px] border-l border-border' : 'w-0'"
        >
          <div class="flex h-full w-[340px] flex-col">
            <div
              class="flex shrink-0 items-center justify-between border-b border-border px-4 py-2"
            >
              <span class="text-sm font-medium">
                {{ $t('page.post.settings') }}
              </span>
              <a-button size="small" type="text" @click="settingsOpen = false">
                <template #icon>
                  <LucideChevronRight />
                </template>
              </a-button>
            </div>
            <div class="min-h-0 flex-1 overflow-y-auto px-4 py-3">
              <PostSettingsPanel />
            </div>
          </div>
        </aside>

        <button
          class="flex w-8 shrink-0 flex-col items-center gap-2 border-border py-3 text-muted-foreground transition-colors hover:bg-accent/40 hover:text-foreground"
          :class="settingsOpen ? '' : 'border-l'"
          :title="$t('page.post.settings')"
          @click="settingsOpen = !settingsOpen"
        >
          <LucideSettings
            class="size-4"
            :class="settingsOpen ? 'text-primary' : ''"
          />
          <span
            class="text-xs tracking-widest"
            style="writing-mode: vertical-rl"
          >
            {{ $t('page.post.settings') }}
          </span>
        </button>
      </div>
    </div>

    <template #footer>
      <div class="flex w-full">
        <a-space class="ml-auto">
          <a-button type="default" @click="handleSaveDraft">
            {{ $t('page.post.button.saveDraft') }}
          </a-button>
          <a-button type="primary" danger @click="handlePublish">
            {{ $t('page.post.button.publish') }}
          </a-button>
        </a-space>
      </div>
    </template>
  </Page>
</template>

<style scoped>
.post-edit-container {
  position: relative;
  width: 100%;
  height: 100%;
}

.translate-btn {
  min-width: 140px;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.translate-btn:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}
</style>
