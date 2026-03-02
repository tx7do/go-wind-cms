import type { PostEditProps } from './types';

import { $t } from '@vben/locales';
import { StorageManager } from '@vben-core/shared/cache';

import { defineStore } from 'pinia';

import { EditorType } from '#/adapter/component/Editor';
import {
  convertToEditorType,
  convertToUIEditorType,
  useLanguageStore,
  usePostStore,
} from '#/stores';

const postStore = usePostStore();
const languageStore = useLanguageStore();

const storageManager = new StorageManager({
  prefix: 'post-draft',
});

const storageKeyMessage = 'post-edit-draft';

/**
 * 文章编辑视图状态接口
 */
interface PostEditViewState {
  loading: boolean; // 加载状态
  formData: PostEditProps; // 表单数据
  languageOptions: { label: string; value: string }[]; // 语言选项列表
  isCreateMode: boolean; // 是否为创建模式
  postId: null | number; // 文章ID（编辑模式下）
}

/**
 * 文章编辑视图状态
 */
export const usePostEditViewStore = defineStore('post-edit-view', {
  state: (): PostEditViewState => ({
    loading: false,
    formData: {
      title: '',
      content: '',
      lang: 'zh-CN',
      editorType: EditorType.MARKDOWN,
    },
    languageOptions: [],
    isCreateMode: true,
    postId: null,
  }),

  actions: {
    /**
     * 初始化编辑模式
     */
    initEditMode(postId: number, initialLang: string) {
      this.isCreateMode = false;
      this.postId = postId;
      this.formData.lang = initialLang;
    },

    /**
     * 初始化创建模式
     */
    initCreateMode(initialLang: string) {
      this.isCreateMode = true;
      this.postId = null;
      this.formData = {
        title: '',
        content: '',
        lang: initialLang,
        editorType: EditorType.MARKDOWN,
      };
    },

    /**
     * 加载语言列表
     */
    async fetchLanguageList() {
      try {
        const resp = await languageStore.listLanguage(
          undefined,
          {},
          undefined,
          ['id'],
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
     * 加载文章数据（编辑模式）
     */
    async fetchPost() {
      if (this.isCreateMode || !this.postId) {
        return null;
      }

      this.loading = true;
      try {
        const item = await postStore.getPost(this.postId);
        if (!item) {
          throw new Error('Post not found');
        }

        if (!item.translations || item.translations.length === 0) {
          throw new Error('No translations found for post');
        }

        // 查找对应语言的翻译
        let langItem = item.translations?.find(
          (t) => t.languageCode === this.formData.lang,
        );

        // 如果没找到对应语言，使用第一个翻译
        if (!langItem) {
          langItem = item.translations?.[0];
          this.formData.lang = langItem?.languageCode || this.formData.lang;
        }

        if (!langItem) {
          throw new Error('No translations found for post');
        }

        // 更新表单数据
        this.formData.id = item.id;
        this.formData.title = langItem.title || '';
        this.formData.content = langItem.content || '';
        this.formData.editorType = convertToUIEditorType(item.editorType);

        return item;
      } finally {
        this.loading = false;
      }
    },

    /**
     * 切换语言
     */
    async switchLanguage(languageCode: string) {
      this.formData.lang = languageCode;
      // 如果是编辑模式，重新加载该语言版本的文章
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
     * 保存草稿数据
     */
    savePostDraft() {
      storageManager.setItem(storageKeyMessage, this.formData);
    },

    /**
     * 加载草稿数据
     */
    loadPostDraft() {
      const draft = storageManager.getItem<PostEditProps>(storageKeyMessage);
      if (draft) {
        this.formData = draft;
      }
    },

    /**
     * 清除草稿数据
     */
    clearPostDraft() {
      storageManager.removeItem(storageKeyMessage);
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
        await (this.isCreateMode
          ? postStore.createPost({
              editorType: convertToEditorType(this.formData.editorType),
              translations: [
                {
                  title: this.formData.title,
                  content: this.formData.content,
                  languageCode: this.formData.lang,
                },
              ],
            })
          : postStore.updatePost(this.formData.id || 0, {
              editorType: convertToEditorType(this.formData.editorType),
              translations: [
                {
                  title: this.formData.title,
                  content: this.formData.content,
                  languageCode: this.formData.lang,
                },
              ],
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
      this.formData = {
        title: '',
        content: '',
        lang: 'zh-CN',
        editorType: EditorType.MARKDOWN,
      };
      this.languageOptions = [];
      this.isCreateMode = true;
      this.postId = null;
    },
  },
});
