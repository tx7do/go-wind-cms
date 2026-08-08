import type { SectionEditProps } from './types';

import { $t } from '@vben/locales';
import { StorageManager } from '@vben-core/shared/cache';

import { defineStore } from 'pinia';

import {
  apiClient,
  fetchListLanguages,
  makeUpdateMask,
  PaginationQuery,
} from '#/api';

const storageManager = new StorageManager({
  prefix: 'section-draft',
});

/**
 * Generate unique cache key based on section ID, language, and mode
 */
function getCacheKey(
  sectionId: null | number,
  lang: string,
  isCreateMode: boolean,
): string {
  if (isCreateMode) {
    return `create-${lang}`;
  }
  return `edit-${sectionId}-${lang}`;
}

/**
 * Section edit view state interface
 */
interface SectionEditViewState {
  loading: boolean;
  needTranslate: boolean;
  formData: SectionEditProps;
  languageOptions: { hasTranslation?: boolean; label: string; value: string }[];
  isCreateMode: boolean;
  sectionId: null | number;
}

/**
 * Section edit view state
 */
export const useSectionEditViewStore = defineStore('section-edit-view', {
  state: (): SectionEditViewState => ({
    loading: false,
    needTranslate: false,
    isCreateMode: true,
    sectionId: null,
    formData: {
      name: '',
      content: '',
      lang: 'zh-CN',
    },
    languageOptions: [],
  }),

  actions: {
    /**
     * Initialize edit mode
     */
    initEditMode(sectionId: number, initialLang: string) {
      this.isCreateMode = false;
      this.needTranslate = false;
      this.sectionId = sectionId;
      this.formData.lang = initialLang;
    },

    /**
     * Initialize create mode
     */
    initCreateMode(initialLang: string) {
      this.isCreateMode = true;
      this.needTranslate = false;
      this.sectionId = null;
      this.formData = {
        name: '',
        content: '',
        lang: initialLang,
      };

      // Try to load draft for this language
      this.loadSectionDraft();
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
     * Load section data (edit mode only)
     */
    async fetchSection() {
      if (this.isCreateMode || !this.sectionId) {
        return null;
      }

      this.loading = true;
      try {
        const item = await apiClient.sectionService.Get({ id: this.sectionId });
        if (!item) {
          throw new Error('Section not found');
        }

        if (!item.translations || item.translations.length === 0) {
          throw new Error('No translations found for section');
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
          throw new Error('No translations found for section');
        }

        // Mark translation status in language options using availableLanguages
        const availableLanguages = item.availableLanguages || [];
        this.languageOptions = this.languageOptions.map((option) => ({
          ...option,
          hasTranslation: availableLanguages.includes(option.value),
        }));

        // Update form data
        this.formData.id = item.id;
        this.formData.pageId = item.pageId;
        this.formData.type = item.type;
        this.formData.name = item.name || '';
        this.formData.sortOrder = item.sortOrder;
        this.formData.config = item.config;
        this.formData.content = langItem.content?.body || '';

        // Try to load draft after fetching backend data
        this.loadSectionDraft();

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
      if (this.isCreateMode) {
        this.loadSectionDraft();
      } else {
        await this.fetchSection();
      }
    },

    /**
     * Update form data
     */
    updateFormData(data: Partial<SectionEditProps>) {
      this.formData = { ...this.formData, ...data };
    },

    /**
     * Save draft data
     */
    saveSectionDraft() {
      const cacheKey = getCacheKey(
        this.sectionId,
        this.formData.lang,
        this.isCreateMode,
      );
      storageManager.setItem(cacheKey, this.formData);
    },

    /**
     * Load draft data
     */
    loadSectionDraft() {
      const cacheKey = getCacheKey(
        this.sectionId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<SectionEditProps>(cacheKey);
      if (draft) {
        this.formData = draft;
        return true;
      }
      return false;
    },

    /**
     * Clear draft data
     */
    clearSectionDraft() {
      const cacheKey = getCacheKey(
        this.sectionId,
        this.formData.lang,
        this.isCreateMode,
      );
      storageManager.removeItem(cacheKey);
    },

    /**
     * Save section
     */
    async saveSection() {
      if (!this.formData.type) {
        return $t('page.section.validation.typeRequired');
      }
      if (!this.formData.pageId) {
        return $t('page.section.validation.pageIdRequired');
      }
      if (!this.formData.content) {
        return $t('page.section.validation.contentRequired');
      }

      try {
        const data = {
          type: this.formData.type as any,
          name: this.formData.name,
          pageId: this.formData.pageId,
          sortOrder: this.formData.sortOrder,
          config: this.formData.config,
          translations: [
            {
              languageCode: this.formData.lang,
              content: { body: this.formData.content },
            },
          ],
        };
        await (this.isCreateMode
          ? apiClient.sectionService.Create({ data: data as any })
          : apiClient.sectionService.Update({
              id: this.formData.id || 0,
              data: data as any,
              updateMask: makeUpdateMask(Object.keys(data).filter((k) => k !== 'translations')),
            }));

        // Clear draft after successful save
        this.clearSectionDraft();

        return '';
      } catch (error) {
        console.error('Failed to save section:', error);
        return $t('page.section.validation.saveFailed');
      }
    },

    /**
     * Reset state
     */
    $reset() {
      this.loading = false;
      this.needTranslate = false;
      this.isCreateMode = true;
      this.sectionId = null;
      this.formData = {
        name: '',
        content: '',
        lang: 'zh-CN',
      };
      this.languageOptions = [];
    },
  },
});
