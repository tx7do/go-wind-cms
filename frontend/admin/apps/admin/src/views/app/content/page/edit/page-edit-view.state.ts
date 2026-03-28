import type { PageEditProps } from './types';

import { $t } from '@vben/locales';
import { StorageManager } from '@vben-core/shared/cache';

import { defineStore } from 'pinia';

import { EditorType } from '#/adapter/component/Editor';
import {
  convertToUIEditorType,
  useLanguageStore,
  usePageStore,
} from '#/stores';

const pageStore = usePageStore();
const languageStore = useLanguageStore();

const storageManager = new StorageManager({
  prefix: 'page-draft',
});

/**
 * Generate unique cache key based on page ID, language, and mode
 */
function getCacheKey(
  pageId: null | number,
  lang: string,
  isCreateMode: boolean,
): string {
  if (isCreateMode) {
    return `create-${lang}`;
  }
  return `edit-${pageId}-${lang}`;
}

/**
 * Page edit view state interface
 */
interface PageEditViewState {
  loading: boolean;
  needTranslate: boolean;
  formData: PageEditProps;
  languageOptions: { hasTranslation?: boolean; label: string; value: string }[];
  isCreateMode: boolean;
  pageId: null | number;
}

/**
 * Page edit view state
 */
export const usePageEditViewStore = defineStore('page-edit-view', {
  state: (): PageEditViewState => ({
    loading: false,
    needTranslate: false,
    isCreateMode: true,
    pageId: null,
    formData: {
      title: '',
      slug: '',
      content: '',
      lang: 'zh-CN',
      editorType: EditorType.MARKDOWN,
    },
    languageOptions: [],
  }),

  actions: {
    /**
     * Initialize edit mode
     */
    initEditMode(pageId: number, initialLang: string) {
      this.isCreateMode = false;
      this.needTranslate = false;
      this.pageId = pageId;
      this.formData.lang = initialLang;
    },

    /**
     * Initialize create mode
     */
    initCreateMode(initialLang: string) {
      this.isCreateMode = true;
      this.needTranslate = false;
      this.pageId = null;
      this.formData = {
        title: '',
        slug: '',
        content: '',
        lang: initialLang,
        editorType: EditorType.MARKDOWN,
      };

      // Try to load draft for this language
      this.loadPageDraft();
    },

    /**
     * Load language list
     */
    async fetchLanguageList() {
      try {
        const resp = await languageStore.listLanguage(
          undefined,
          {},
          undefined,
          ['sortOrder'],
        );
        this.languageOptions =
          resp.items?.map((lang) => ({
            label: lang.nativeName || '',
            value: lang.languageCode || '',
          })) || [];
        return this.languageOptions;
      } catch (error) {
        console.error('Failed to load language list:', error);
        this.languageOptions = [];
        throw error;
      }
    },

    /**
     * Load page data (edit mode only)
     */
    async fetchPage() {
      if (this.isCreateMode || !this.pageId) {
        return null;
      }

      this.loading = true;
      try {
        const item = await pageStore.getPage(this.pageId);
        if (!item) {
          throw new Error('Page not found');
        }

        if (!item.translations || item.translations.length === 0) {
          throw new Error('No translations found for page');
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
          console.log(
            'No translation found for selected language, using first available translation',
            this.formData.lang,
          );
        }

        if (!langItem) {
          throw new Error('No translations found for page');
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
        this.formData.slug = langItem.slug || '';
        this.formData.content = langItem.content || '';
        this.formData.editorType = convertToUIEditorType(item.editorType);
        this.formData.parentId = item.parentId;
        this.formData.type = item.type;
        this.formData.status = item.status;
        this.formData.showInNavigation = item.showInNavigation;
        this.formData.disallowComment = item.disallowComment;
        this.formData.template = item.template;
        this.formData.isCustomTemplate = item.isCustomTemplate;
        this.formData.customHead = item.customHead;
        this.formData.customFoot = item.customFoot;
        this.formData.sortOrder = item.sortOrder;

        // Try to load draft after fetching backend data
        // Draft will override backend data if exists
        this.loadPageDraft();

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
      // If in create mode, try to load draft for this language
      if (this.isCreateMode) {
        this.loadPageDraft();
      } else {
        // If in edit mode, reload the page for this language
        await this.fetchPage();
      }
    },

    /**
     * Update form data
     */
    updateFormData(data: Partial<PageEditProps>) {
      this.formData = { ...this.formData, ...data };
    },

    /**
     * Save draft data
     */
    savePageDraft() {
      const cacheKey = getCacheKey(
        this.pageId,
        this.formData.lang,
        this.isCreateMode,
      );
      storageManager.setItem(cacheKey, this.formData);
      console.log(`Draft saved with key: ${cacheKey}`);
    },

    /**
     * Load draft data
     */
    loadPageDraft() {
      const cacheKey = getCacheKey(
        this.pageId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<PageEditProps>(cacheKey);
      if (draft) {
        console.log(`Draft loaded with key: ${cacheKey}`, draft);
        this.formData = draft;
        return true;
      }
      console.log(`No draft found with key: ${cacheKey}`);
      return false;
    },

    /**
     * Clear draft data
     */
    clearPageDraft() {
      const cacheKey = getCacheKey(
        this.pageId,
        this.formData.lang,
        this.isCreateMode,
      );
      storageManager.removeItem(cacheKey);
      console.log(`Draft cleared with key: ${cacheKey}`);
    },

    /**
     * Check if draft exists
     */
    hasDraft(): boolean {
      const cacheKey = getCacheKey(
        this.pageId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<PageEditProps>(cacheKey);
      return draft !== null && draft !== undefined;
    },

    /**
     * Publish page
     */
    async publishPage() {
      if (!this.formData.title) {
        return $t('page.page.validation.titleRequired');
      }
      if (!this.formData.slug) {
        return $t('page.page.validation.slugRequired');
      }
      if (!this.formData.content) {
        return $t('page.page.validation.contentRequired');
      }

      try {
        await (this.isCreateMode
          ? pageStore.createPage({
              editorType: this.formData.editorType as any,
              parentId: this.formData.parentId,
              type: this.formData.type,
              status: this.formData.status,
              showInNavigation: this.formData.showInNavigation,
              disallowComment: this.formData.disallowComment,
              template: this.formData.template,
              isCustomTemplate: this.formData.isCustomTemplate,
              customHead: this.formData.customHead,
              customFoot: this.formData.customFoot,
              sortOrder: this.formData.sortOrder,
              translations: [
                {
                  title: this.formData.title,
                  slug: this.formData.slug,
                  content: this.formData.content,
                  languageCode: this.formData.lang,
                },
              ],
            })
          : pageStore.updatePage(this.formData.id || 0, {
              editorType: this.formData.editorType as any,
              parentId: this.formData.parentId,
              type: this.formData.type,
              status: this.formData.status,
              showInNavigation: this.formData.showInNavigation,
              disallowComment: this.formData.disallowComment,
              template: this.formData.template,
              isCustomTemplate: this.formData.isCustomTemplate,
              customHead: this.formData.customHead,
              customFoot: this.formData.customFoot,
              sortOrder: this.formData.sortOrder,
              translations: [
                {
                  title: this.formData.title,
                  slug: this.formData.slug,
                  content: this.formData.content,
                  languageCode: this.formData.lang,
                },
              ],
            }));

        // Clear draft after successful publish
        this.clearPageDraft();

        return '';
      } catch (error) {
        console.error('Failed to publish page:', error);
        return $t('page.page.validation.publishFailed');
      }
    },

    /**
     * Reset state
     */
    $reset() {
      console.log('resetting page edit view state');
      this.loading = false;
      this.needTranslate = false;
      this.isCreateMode = true;
      this.pageId = null;
      this.formData = {
        title: '',
        slug: '',
        content: '',
        lang: 'zh-CN',
        editorType: EditorType.MARKDOWN,
      };
      this.languageOptions = [];
    },
  },
});
