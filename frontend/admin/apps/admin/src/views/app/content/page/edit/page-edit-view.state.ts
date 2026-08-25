import type { PageEditProps } from './types';

import { $t } from '@vben/locales';
import { StorageManager } from '@vben-core/shared/cache';

import { defineStore } from 'pinia';

import { EditorType } from '#/adapter/component/Editor';
import {
  apiClient,
  convertToUIEditorType,
  fetchListLanguages,
  makeUpdateMask,
  PaginationQuery,
} from '#/api';

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
      type: 'PAGE_TYPE_DEFAULT',
      status: 'PAGE_STATUS_DRAFT',
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
        type: 'PAGE_TYPE_DEFAULT',
        status: 'PAGE_STATUS_DRAFT',
      };

      // Try to load draft for this language
      this.loadPageDraft();
    },

    /**
     * Load language list
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
        const item = await apiClient.pageService.Get({ id: this.pageId });
        if (!item) {
          throw new Error('Page not found');
        }

        // 页面可能尚无任何翻译（新建页、或演示数据中无翻译的页面）。
        // 此时仍加载页面级元数据，并让作者以当前语言撰写首条翻译，
        // 而非直接抛错阻断编辑。
        const hasTranslations =
          !!item.translations && item.translations.length > 0;

        // Find translation for selected language
        let langItem = hasTranslations
          ? item.translations?.find(
              (t) => t.languageCode === this.formData.lang,
            )
          : undefined;

        this.needTranslate = false;

        // If translation not found, use first available translation
        if (!langItem) {
          langItem = hasTranslations ? item.translations?.[0] : undefined;
          this.needTranslate = true;
        }

        // Mark translation status in language options using availableLanguages
        const availableLanguages = item.availableLanguages || [];
        this.languageOptions = this.languageOptions.map((option) => ({
          ...option,
          hasTranslation: availableLanguages.includes(option.value),
        }));

        // Update form data
        this.formData.id = item.id;
        this.formData.title = langItem?.title || '';
        this.formData.slug = langItem?.slug || '';
        this.formData.content = langItem?.content || '';
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
        this.formData.contentModelId = item.contentModelId;
        this.formData.customFields = item.customFields ?? {};

        // 嵌套区块：从页面整体水合（后端按 page_id 返回该页全部 section 及其各语言翻译）。
        this.formData.sections = (item.sections ?? []).map((s) => ({
          id: s.id,
          type: s.type,
          name: s.name,
          sortOrder: s.sortOrder,
          config: s.config,
          translations: (s.translations ?? []).map((tr) => ({
            languageCode: tr.languageCode,
            content: tr.content,
          })),
        }));

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
        this.formData = draft;
        return true;
      }
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
        const data = {
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
          contentModelId: this.formData.contentModelId,
          customFields: this.formData.customFields,
          sections: this.formData.sections,
          translations: [
            {
              title: this.formData.title,
              slug: this.formData.slug,
              content: this.formData.content,
              languageCode: this.formData.lang,
            },
          ],
        } as any;
        await (this.isCreateMode
          ? apiClient.pageService.Create({ data })
          : apiClient.pageService.Update({
              id: this.formData.id || 0,
              data,
              updateMask: makeUpdateMask(
                // translations 与 sections 均为整体替换字段，不纳入 updateMask
                Object.keys(data).filter(
                  (k) => k !== 'translations' && k !== 'sections',
                ),
              ),
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
        type: 'PAGE_TYPE_DEFAULT',
        status: 'PAGE_STATUS_DRAFT',
      };
      this.languageOptions = [];
    },
  },
});
