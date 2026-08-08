import type { CategoryOption, PostEditProps } from './types';

import { $t } from '@vben/locales';
import { StorageManager } from '@vben-core/shared/cache';

import { defineStore } from 'pinia';

import { EditorType } from '#/adapter/component/Editor';
import {
  apiClient,
  convertToEditorType,
  convertToUIEditorType,
  fetchListCategories,
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
  languageOptions: { hasTranslation?: boolean; label: string; value: string }[]; // 语言选项列表（带翻译标记）
  categoryOptions: CategoryOption[]; // 分类选项列表（按当前语言展示名称）
  isCreateMode: boolean; // 是否为创建模式
  postId: null | number; // 文章ID（编辑模式下）
}

/**
 * 文章编辑视图状态
 */
export const usePostEditViewStore = defineStore('post-edit-view', {
  state: (): PostEditViewState => ({
    loading: false,
    needTranslate: false,
    isCreateMode: true,
    postId: null,
    formData: {
      title: '',
      content: '',
      lang: 'zh-CN',
      editorType: EditorType.MARKDOWN,
    },
    languageOptions: [],
    categoryOptions: [],
  }),

  actions: {
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
        editorType: EditorType.MARKDOWN,
      };

      // Try to load draft for this language
      this.loadPostDraft();
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
          resp.items?.map((lang) => ({
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
     * 加载分类列表（按当前编辑语言展示分类名称）
     */
    async fetchCategoryList() {
      try {
        const resp = await fetchListCategories(new PaginationQuery());
        this.categoryOptions =
          resp.items
            ?.filter((category) => category.id !== undefined)
            .map((category) => {
              const id = category.id as number;
              // 优先匹配当前语言的翻译，回退到首个可用翻译
              const translations = category.translations || [];
              const langItem =
                translations.find(
                  (t) => t.languageCode === this.formData.lang,
                ) || translations[0];
              return {
                label: langItem?.name || `#${id}`,
                value: id,
              };
            }) || [];
        return this.categoryOptions;
      } catch (error) {
        console.error('获取分类列表失败:', error);
        this.categoryOptions = [];
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
        this.formData.categoryIds = [...(item.categoryIds || [])];

        // 分类名称依赖当前语言，切换语言/加载文章后刷新分类选项
        await this.fetchCategoryList();

        // Try to load draft after fetching backend data
        // Draft will override backend data if exists
        this.loadPostDraft();

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
      // 分类名称依赖当前语言，切换语言后刷新分类选项
      await this.fetchCategoryList();
      // If in create mode, try to load draft for this language
      if (this.isCreateMode) {
        // this.loadPostDraft();
      } else {
        // If in edit mode, reload the article for this language
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
     * 发布文章
     */
    async publishPost() {
      if (!this.formData.title) {
        return $t('page.post.validation.titleRequired');
      }
      if (!this.formData.content) {
        return $t('page.post.validation.contentRequired');
      }

      try {
        const data = {
          editorType: convertToEditorType(this.formData.editorType),
          categoryIds: this.formData.categoryIds || [],
          translations: [
            {
              title: this.formData.title,
              content: this.formData.content,
              languageCode: this.formData.lang,
            },
          ],
        };
        await (this.isCreateMode
          ? apiClient.postService.Create({ data })
          : apiClient.postService.Update({
              id: this.formData.id || 0,
              data,
              updateMask: makeUpdateMask(Object.keys(data).filter((k) => k !== 'translations')),
            }));

        // 发布成功后清除草稿
        this.clearPostDraft();

        return '';
      } catch (error) {
        console.error('Failed to publish post:', error);
        return $t('page.post.validation.publishFailed');
      }
    },

    /**
     * 重置状态
     */
    $reset() {
      this.loading = false;
      this.needTranslate = false;
      this.isCreateMode = true;
      this.postId = null;
      this.formData = {
        title: '',
        content: '',
        lang: 'zh-CN',
        editorType: EditorType.MARKDOWN,
      };
      this.languageOptions = [];
      this.categoryOptions = [];
    },
  },
});
