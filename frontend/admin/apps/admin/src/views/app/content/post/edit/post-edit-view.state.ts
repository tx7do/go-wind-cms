import type { PostEditProps } from './types';

import type { contentservicev1_Post } from '#/api/generated/admin/service/v1';

import { $t } from '@vben/locales';
import { StorageManager } from '@vben-core/shared/cache';

import { defineStore } from 'pinia';

import { EditorType } from '#/adapter/component/Editor';
import {
  apiClient,
  convertToEditorType,
  convertToUIEditorType,
  fetchFlattenedCategoryOptions,
  fetchListLanguages,
  makeUpdateMask,
  PaginationQuery,
} from '#/api';

const storageManager = new StorageManager({
  prefix: 'post-draft',
});

/**
 * Generate unique cache key based on post ID, language, and mode
 */
function getCacheKey(
  postId: null | number,
  lang: string,
  isCreateMode: boolean,
): string {
  if (isCreateMode) {
    return `create-${lang}`;
  }
  return `edit-${postId}-${lang}`;
}

/**
 * 文章编辑视图状态接口
 */
interface PostEditViewState {
  loading: boolean; // 加载状态
  needTranslate: boolean; // 是否需要翻译
  formData: PostEditProps; // 表单数据
  languageOptions: {
    hasTranslation?: boolean;
    label: string;
    value: string;
  }[]; // 语言选项列表（带翻译标记）
  isCreateMode: boolean; // 是否为创建模式
  postId: null | number; // 文章ID（编辑模式下）
  pendingDraft: boolean; // 是否存在待恢复的本地草稿（由页面询问用户处置）
  savedSnapshot: string; // 最近一次与服务端同步的表单快照（用于未保存修改检测）
  categoryOptions: { label: string; value: number }[]; // 分类下拉选项（编辑页必选项）
  categoryOptionsLoading: boolean; // 分类选项加载状态
}

/**
 * 编辑页保存载荷：只允许携带 posts 表真实存在的列（或仓库层显式处理的关联键）。
 * availableLanguages 是 proto 计算字段、tagIds 由仓库层按 TagIds!=nil 分支处理，
 * 带进 updateMask 会让后端生成 SET 不存在的列（如 available_languages=NULL）导致 500。
 */
type PostSaveData = Omit<
  contentservicev1_Post,
  'availableLanguages' | 'customFields' | 'tagIds'
>;

/** 参与字段掩码的保存键（与 savePost 的 data 字段一一对应） */
const POST_SAVE_MASK_KEYS = [
  'editorType',
  'status',
  'code',
  'isFeatured',
  'sortOrder',
  'disallowComment',
  'categoryIds',
  'thumbnail',
  'publishTime',
  'translations',
] as const;

/**
 * 文章编辑视图状态
 */
export const usePostEditViewStore = defineStore('post-edit-view', {
  state: (): PostEditViewState => ({
    loading: false,
    needTranslate: false,
    isCreateMode: true,
    postId: null,
    pendingDraft: false,
    savedSnapshot: '',
    categoryOptions: [],
    categoryOptionsLoading: false,
    formData: {
      title: '',
      content: '',
      lang: 'zh-CN',
      editorType: EditorType.RICH_TEXT,
      seo: {},
    },
    languageOptions: [],
  }),

  actions: {
    /**
     * 加载分类下拉选项（分类树拍平，子分类用“父级 / ”前缀区分层级）
     */
    async fetchCategoryOptions() {
      if (this.categoryOptionsLoading) {
        return;
      }
      this.categoryOptionsLoading = true;
      try {
        this.categoryOptions = await fetchFlattenedCategoryOptions(
          this.formData.lang,
        );
      } finally {
        this.categoryOptionsLoading = false;
      }
    },

    /**
     * 初始化编辑模式
     */
    initEditMode(postId: number, initialLang: string) {
      this.isCreateMode = false;
      this.needTranslate = false;
      this.postId = postId;
      this.formData.lang = initialLang;
    },

    /**
     * Initialize create mode
     */
    initCreateMode(initialLang: string) {
      this.isCreateMode = true;
      this.needTranslate = false;
      this.postId = null;
      this.formData = {
        title: '',
        content: '',
        lang: initialLang,
        editorType: EditorType.RICH_TEXT,
        seo: {},
      };

      // 是否有可恢复的本地草稿由页面询问用户，不再静默覆盖
      this.pendingDraft = this.hasDraft();
      this.refreshSavedSnapshot();
    },

    /**
     * 加载语言列表
     */
    async fetchLanguageList() {
      try {
        const resp = await fetchListLanguages(
          new PaginationQuery({ orderBy: ['sortOrder'] }),
        );
        this.languageOptions =
          resp.items
            ?.filter((lang) => lang.isEnabled !== false)
            .map((lang) => ({
              label: lang.nativeName || '',
              value: lang.languageCode || '',
            })) || [];
        return this.languageOptions;
      } catch (error) {
        console.error('获取语言列表失败:', error);
        this.languageOptions = [];
        throw error;
      }
    },

    /**
     * Load post data (edit mode only)
     */
    async fetchPost() {
      if (this.isCreateMode || !this.postId) {
        return null;
      }

      this.loading = true;
      try {
        const item = await apiClient.postService.Get({ id: this.postId });
        if (!item) {
          throw new Error('Post not found');
        }

        if (!item.translations || item.translations.length === 0) {
          throw new Error('No translations found for post');
        }

        // Find translation for selected language
        let langItem = item.translations?.find(
          (t) => t.languageCode === this.formData.lang,
        );

        this.needTranslate = false;

        // If translation not found, use first available translation
        if (!langItem) {
          langItem = item.translations?.[0];
          this.needTranslate = true;
        }

        if (!langItem) {
          throw new Error('No translations found for post');
        }

        // Mark translation status in language options using availableLanguages
        const availableLanguages = item.availableLanguages || [];
        this.languageOptions = this.languageOptions.map((option) => ({
          ...option,
          hasTranslation: availableLanguages.includes(option.value),
        }));

        // Update form data
        this.formData.id = item.id;
        this.formData.title = langItem.title || '';
        this.formData.content = langItem.content || '';
        this.formData.editorType = convertToUIEditorType(item.editorType);
        this.formData.status = item.status ?? undefined;

        // Post 级字段
        this.formData.code = item.code || '';
        this.formData.categoryIds = item.categoryIds
          ? [...item.categoryIds]
          : [];
        this.formData.isFeatured = item.isFeatured ?? false;
        this.formData.sortOrder = item.sortOrder ?? 0;
        this.formData.disallowComment = item.disallowComment ?? false;
        this.formData.publishTime = item.publishTime ?? undefined;
        // 缩略图帖子级，全语言共用；旧数据无帖子级值时回退翻译级
        this.formData.thumbnail = item.thumbnail || '';

        // 翻译级字段（当前语言）；seo 保持对象形态，设置面板的输入框直接绑定其属性
        this.formData.summary = langItem.summary || '';
        this.formData.slug = langItem.slug || '';
        this.formData.seo = { ...(langItem.seo ?? {}) };

        // 是否有可恢复的本地草稿由页面询问用户，不再静默覆盖
        this.pendingDraft = this.hasDraft();
        this.refreshSavedSnapshot();

        return item;
      } finally {
        this.loading = false;
      }
    },

    /**
     * Switch language
     */
    async switchLanguage(languageCode: string) {
      this.formData.lang = languageCode;
      // 分类名称依赖当前语言，清缓存后按新语言重建选项
      this.categoryOptions = [];
      await this.fetchCategoryOptions();
      // If in create mode, keep the draft as-is; edit mode reloads the post
      if (!this.isCreateMode) {
        await this.fetchPost();
      }
    },

    /**
     * 更新表单数据
     */
    updateFormData(data: Partial<PostEditProps>) {
      this.formData = { ...this.formData, ...data };
    },

    /**
     * Save draft data
     */
    savePostDraft() {
      const cacheKey = getCacheKey(
        this.postId,
        this.formData.lang,
        this.isCreateMode,
      );
      storageManager.setItem(cacheKey, this.formData);
    },

    /**
     * Load draft data
     */
    loadPostDraft() {
      const cacheKey = getCacheKey(
        this.postId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<PostEditProps>(cacheKey);
      if (draft) {
        this.formData = draft;
        // 早期草稿可能没有 seo 字段，补齐为对象避免绑定空引用
        if (!this.formData.seo) {
          this.formData.seo = {};
        }
        return true;
      }
      return false;
    },

    /**
     * Clear draft data
     */
    clearPostDraft() {
      const cacheKey = getCacheKey(
        this.postId,
        this.formData.lang,
        this.isCreateMode,
      );
      storageManager.removeItem(cacheKey);
    },

    /**
     * Check if draft exists
     */
    hasDraft(): boolean {
      const cacheKey = getCacheKey(
        this.postId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<PostEditProps>(cacheKey);
      return draft !== null && draft !== undefined;
    },

    /**
     * 恢复待处置的本地草稿（覆盖当前表单）
     */
    restorePendingDraft(): boolean {
      this.pendingDraft = false;
      return this.loadPostDraft();
    },

    /**
     * 丢弃待处置的本地草稿（保留当前表单内容）
     */
    discardPendingDraft() {
      this.pendingDraft = false;
      this.clearPostDraft();
    },

    /**
     * 表单中需要参与未保存检测的字段快照
     */
    snapshotPayload() {
      const f = this.formData;
      return {
        title: f.title,
        content: f.content,
        editorType: f.editorType,
        code: f.code ?? '',
        categoryIds: [...(f.categoryIds ?? [])].sort(),
        isFeatured: f.isFeatured ?? false,
        sortOrder: f.sortOrder ?? 0,
        disallowComment: f.disallowComment ?? false,
        publishTime: f.publishTime,
        summary: f.summary ?? '',
        slug: f.slug ?? '',
        thumbnail: f.thumbnail ?? '',
        seo: f.seo ?? {},
      };
    },

    /**
     * 记录当前表单为“已同步”基线
     */
    refreshSavedSnapshot() {
      this.savedSnapshot = JSON.stringify(this.snapshotPayload());
    },

    /**
     * 表单相对最近一次同步是否有未保存修改
     */
    hasUnsavedChanges(): boolean {
      if (!this.savedSnapshot) {
        return false;
      }
      return JSON.stringify(this.snapshotPayload()) !== this.savedSnapshot;
    },

    /**
     * 保存文章（status 决定保存为草稿还是发布）。
     * 新建模式下保存成功后会自动切换为编辑模式（后续保存走 Update）。
     * 返回空字符串表示成功，否则返回错误提示文案。
     */
    async savePost(
      status: NonNullable<PostEditProps['status']>,
    ): Promise<string> {
      if (!this.formData.title) {
        return $t('page.post.validation.titleRequired');
      }
      if (!this.formData.content) {
        return $t('page.post.validation.contentRequired');
      }
      if (!this.formData.categoryIds || this.formData.categoryIds.length === 0) {
        return $t('page.post.validation.categoryRequired');
      }

      const wasCreateMode = this.isCreateMode;
      // 全部为空的 seo 不提交，避免用空对象覆盖服务端已有值
      const seo = this.formData.seo ?? {};
      const hasSeo = Object.values(seo).some(
        (v) => typeof v === 'string' && v.length > 0,
      );
      const data: PostSaveData = {
        editorType: convertToEditorType(this.formData.editorType),
        status,
        code: this.formData.code || '',
        isFeatured: this.formData.isFeatured ?? false,
        sortOrder: this.formData.sortOrder ?? 0,
        disallowComment: this.formData.disallowComment ?? false,
        categoryIds: this.formData.categoryIds ?? [],
        // 缩略图为帖子级字段，全语言共用一张
        thumbnail: this.formData.thumbnail || '',
        translations: [
          {
            title: this.formData.title,
            content: this.formData.content,
            languageCode: this.formData.lang,
            summary: this.formData.summary || '',
            slug: this.formData.slug || '',
            seo: hasSeo ? seo : undefined,
          },
        ],
      };
      // 始终携带 publishTime 进掩码：未设置时提交 null，由后端置空该列（支持清空）。
      // 生成类型的 timestamp 为 string，无 null 形态，故绕过类型标注赋值。
      (data as { publishTime?: string | null }).publishTime =
        this.formData.publishTime ?? null;

      try {
        if (wasCreateMode) {
          const created = await apiClient.postService.Create({
            data: data as contentservicev1_Post,
          });

          // 先按创建模式的 key 清理本地草稿，再切换为编辑模式
          this.clearPostDraft();

          if (created?.id) {
            this.postId = Number(created.id);
            this.formData.id = Number(created.id);
            this.isCreateMode = false;
          }
        } else {
          await apiClient.postService.Update({
            id: this.formData.id || 0,
            data: data as contentservicev1_Post,
            updateMask: makeUpdateMask([...POST_SAVE_MASK_KEYS]),
          });
          this.clearPostDraft();
        }

        this.formData.status = status;
        this.refreshSavedSnapshot();

        return '';
      } catch (error) {
        console.error('Failed to save post:', error);
        return wasCreateMode
          ? $t('page.post.validation.publishFailed')
          : $t('page.post.validation.saveDraftFailed');
      }
    },

    /**
     * 保存为服务端草稿
     */
    async saveServerDraft(): Promise<string> {
      return this.savePost('POST_STATUS_DRAFT');
    },

    /**
     * 发布文章。设置了未来发布时间的未发布文章会转为定时发布（SCHEDULED），
     * 到达时间后由后端调度器自动置为已发布。
     */
    async publishPost(): Promise<string> {
      const publishTime = this.formData.publishTime;
      if (publishTime) {
        const ts = new Date(publishTime).getTime();
        const isFuture = !Number.isNaN(ts) && ts > Date.now() + 60_000;
        const isScheduled = this.formData.status === 'POST_STATUS_SCHEDULED';
        const isPublished = this.formData.status === 'POST_STATUS_PUBLISHED';

        if (isFuture && !isScheduled && !isPublished) {
          return this.savePost('POST_STATUS_SCHEDULED');
        }
      }
      return this.savePost('POST_STATUS_PUBLISHED');
    },

    /**
     * 重置状态
     */
    $reset() {
      this.loading = false;
      this.needTranslate = false;
      this.isCreateMode = true;
      this.postId = null;
      this.pendingDraft = false;
      this.savedSnapshot = '';
      this.categoryOptions = [];
      this.categoryOptionsLoading = false;
      this.formData = {
        title: '',
        content: '',
        lang: 'zh-CN',
        editorType: EditorType.RICH_TEXT,
        seo: {},
      };
      this.languageOptions = [];
    },
  },
});
