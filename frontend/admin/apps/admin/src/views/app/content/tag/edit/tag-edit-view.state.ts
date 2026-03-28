import type { TagEditProps } from './types';

import { $t } from '@vben/locales';
import { StorageManager } from '@vben-core/shared/cache';

import { defineStore } from 'pinia';

import { useLanguageStore, useTagStore } from '#/stores';

const tagStore = useTagStore();
const languageStore = useLanguageStore();

const storageManager = new StorageManager({
  prefix: 'tag-draft',
});

/**
 * Generate unique cache key based on tag ID, language, and mode
 */
function getCacheKey(
  tagId: null | number,
  lang: string,
  isCreateMode: boolean,
): string {
  if (isCreateMode) {
    return `create-${lang}`;
  }
  return `edit-${tagId}-${lang}`;
}

/**
 * Tag edit view state interface
 */
interface TagEditViewState {
  loading: boolean;
  needTranslate: boolean;
  formData: TagEditProps;
  languageOptions: { hasTranslation?: boolean; label: string; value: string }[];
  isCreateMode: boolean;
  tagId: null | number;
}

/**
 * Tag edit view state
 */
export const useTagEditViewStore = defineStore('tag-edit-view', {
  state: (): TagEditViewState => ({
    loading: false,
    needTranslate: false,
    isCreateMode: true,
    tagId: null,
    formData: {
      name: '',
      slug: '',
      description: '',
      lang: 'zh-CN',
      sortOrder: 0,
      isFeatured: false,
    },
    languageOptions: [],
  }),

  actions: {
    /**
     * Initialize edit mode
     */
    initEditMode(tagId: number, initialLang: string) {
      this.isCreateMode = false;
      this.needTranslate = false;
      this.tagId = tagId;
      this.formData.lang = initialLang;
    },

    /**
     * Initialize create mode
     */
    initCreateMode(initialLang: string) {
      this.isCreateMode = true;
      this.needTranslate = false;
      this.tagId = null;
      this.formData = {
        name: '',
        slug: '',
        description: '',
        lang: initialLang,
        sortOrder: 0,
        isFeatured: false,
      };

      // Try to load draft for this language
      this.loadTagDraft();
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
     * Load tag data (edit mode only)
     */
    async fetchTag() {
      if (this.isCreateMode || !this.tagId) {
        return null;
      }

      this.loading = true;
      try {
        const item = await tagStore.getTag(this.tagId);
        if (!item) {
          throw new Error('Tag not found');
        }

        if (!item.translations || item.translations.length === 0) {
          throw new Error('No translations found for tag');
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
          throw new Error('No translations found for tag');
        }

        // Mark translation status in language options using availableLanguages
        const availableLanguages = item.availableLanguages || [];
        this.languageOptions = this.languageOptions.map((option) => ({
          ...option,
          hasTranslation: availableLanguages.includes(option.value),
        }));

        // Update form data
        this.formData.id = item.id;
        this.formData.name = langItem.name || '';
        this.formData.slug = langItem.slug || '';
        this.formData.description = langItem.description || '';
        this.formData.coverImage = langItem.coverImage;
        this.formData.template = langItem.template;
        this.formData.color = item.color;
        this.formData.icon = item.icon;
        this.formData.group = item.group;
        this.formData.sortOrder = item.sortOrder;
        this.formData.isFeatured = item.isFeatured;
        this.formData.status = item.status;

        // Try to load draft after fetching backend data
        // Draft will override backend data if exists
        this.loadTagDraft();

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
        this.loadTagDraft();
      } else {
        // If in edit mode, reload the tag for this language
        await this.fetchTag();
      }
    },

    /**
     * Update form data
     */
    updateFormData(data: Partial<TagEditProps>) {
      this.formData = { ...this.formData, ...data };
    },

    /**
     * Save draft data
     */
    saveTagDraft() {
      const cacheKey = getCacheKey(
        this.tagId,
        this.formData.lang,
        this.isCreateMode,
      );
      storageManager.setItem(cacheKey, this.formData);
      console.log(`Draft saved with key: ${cacheKey}`);
    },

    /**
     * Load draft data
     */
    loadTagDraft() {
      const cacheKey = getCacheKey(
        this.tagId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<TagEditProps>(cacheKey);
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
    clearTagDraft() {
      const cacheKey = getCacheKey(
        this.tagId,
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
        this.tagId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<TagEditProps>(cacheKey);
      return draft !== null && draft !== undefined;
    },

    /**
     * Publish tag
     */
    async publishTag() {
      if (!this.formData.name) {
        return $t('page.tag.validation.nameRequired');
      }
      if (!this.formData.slug) {
        return $t('page.tag.validation.slugRequired');
      }

      try {
        await (this.isCreateMode
          ? tagStore.createTag({
              color: this.formData.color,
              icon: this.formData.icon,
              group: this.formData.group,
              sortOrder: this.formData.sortOrder,
              isFeatured: this.formData.isFeatured,
              status: this.formData.status,
              translations: [
                {
                  name: this.formData.name,
                  slug: this.formData.slug,
                  description: this.formData.description,
                  coverImage: this.formData.coverImage,
                  template: this.formData.template,
                  languageCode: this.formData.lang,
                },
              ],
            })
          : tagStore.updateTag(this.formData.id || 0, {
              color: this.formData.color,
              icon: this.formData.icon,
              group: this.formData.group,
              sortOrder: this.formData.sortOrder,
              isFeatured: this.formData.isFeatured,
              status: this.formData.status,
              translations: [
                {
                  name: this.formData.name,
                  slug: this.formData.slug,
                  description: this.formData.description,
                  coverImage: this.formData.coverImage,
                  template: this.formData.template,
                  languageCode: this.formData.lang,
                },
              ],
            }));

        // Clear draft after successful publish
        this.clearTagDraft();

        return '';
      } catch (error) {
        console.error('Failed to publish tag:', error);
        return $t('page.tag.validation.publishFailed');
      }
    },

    /**
     * Reset state
     */
    $reset() {
      console.log('resetting tag edit view state');
      this.loading = false;
      this.needTranslate = false;
      this.isCreateMode = true;
      this.tagId = null;
      this.formData = {
        name: '',
        slug: '',
        description: '',
        lang: 'zh-CN',
        sortOrder: 0,
        isFeatured: false,
      };
      this.languageOptions = [];
    },
  },
});
