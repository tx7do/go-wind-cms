import type { CategoryEditProps } from './types';

import { StorageManager } from '@vben-core/shared/cache';

import { defineStore } from 'pinia';

import { useCategoryStore, useLanguageStore } from '#/stores';

const categoryStore = useCategoryStore();
const languageStore = useLanguageStore();

const storageManager = new StorageManager({
  prefix: 'category-draft',
});

/**
 * Generate unique cache key based on category ID, language, and mode
 */
function getCacheKey(
  categoryId: null | number,
  lang: string,
  isCreateMode: boolean,
): string {
  if (isCreateMode) {
    return `create-${lang}`;
  }
  return `edit-${categoryId}-${lang}`;
}

/**
 * Category edit view state interface
 */
interface CategoryEditViewState {
  loading: boolean;
  needTranslate: boolean;
  formData: CategoryEditProps;
  languageOptions: { hasTranslation?: boolean; label: string; value: string }[];
  isCreateMode: boolean;
  categoryId: null | number;
}

/**
 * Category edit view state
 */
export const useCategoryEditViewStore = defineStore('category-edit-view', {
  state: (): CategoryEditViewState => ({
    loading: false,
    needTranslate: false,
    isCreateMode: true,
    categoryId: null,
    formData: {
      name: '',
      slug: '',
      description: '',
      lang: 'zh-CN',
      sortOrder: 0,
      isNav: false,
    },
    languageOptions: [],
  }),

  actions: {
    /**
     * Initialize edit mode
     */
    initEditMode(categoryId: number, initialLang: string) {
      this.isCreateMode = false;
      this.needTranslate = false;
      this.categoryId = categoryId;
      this.formData.lang = initialLang;
    },

    /**
     * Initialize create mode
     */
    initCreateMode(initialLang: string) {
      this.isCreateMode = true;
      this.needTranslate = false;
      this.categoryId = null;
      this.formData = {
        name: '',
        slug: '',
        description: '',
        lang: initialLang,
        sortOrder: 0,
        isNav: false,
      };

      // Try to load draft for this language
      this.loadCategoryDraft();
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
     * Load category data (edit mode only)
     */
    async fetchCategory() {
      if (this.isCreateMode || !this.categoryId) {
        return null;
      }

      this.loading = true;
      try {
        const item = await categoryStore.getCategory(this.categoryId);
        if (!item) {
          throw new Error('Category not found');
        }

        if (!item.translations || item.translations.length === 0) {
          throw new Error('No translations found for category');
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
          throw new Error('No translations found for category');
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
        this.formData.parentId = item.parentId;
        this.formData.icon = item.icon;
        this.formData.isNav = item.isNav;
        this.formData.sortOrder = item.sortOrder;
        this.formData.status = item.status;

        // Try to load draft after fetching backend data
        // Draft will override backend data if exists
        this.loadCategoryDraft();

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
        this.loadCategoryDraft();
      } else {
        // If in edit mode, reload the category for this language
        await this.fetchCategory();
      }
    },

    /**
     * Update form data
     */
    updateFormData(data: Partial<CategoryEditProps>) {
      this.formData = { ...this.formData, ...data };
    },

    /**
     * Save draft data
     */
    saveCategoryDraft() {
      const cacheKey = getCacheKey(
        this.categoryId,
        this.formData.lang,
        this.isCreateMode,
      );
      storageManager.setItem(cacheKey, this.formData);
      console.log(`Draft saved with key: ${cacheKey}`);
    },

    /**
     * Load draft data
     */
    loadCategoryDraft() {
      const cacheKey = getCacheKey(
        this.categoryId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<CategoryEditProps>(cacheKey);
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
    clearCategoryDraft() {
      const cacheKey = getCacheKey(
        this.categoryId,
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
        this.categoryId,
        this.formData.lang,
        this.isCreateMode,
      );
      const draft = storageManager.getItem<CategoryEditProps>(cacheKey);
      return draft !== null && draft !== undefined;
    },

    /**
     * Publish category
     */
    async publishCategory() {
      if (!this.formData.name) {
        return 'Category name is required';
      }
      if (!this.formData.slug) {
        return 'Category slug is required';
      }

      try {
        await (this.isCreateMode
          ? categoryStore.createCategory({
              parentId: this.formData.parentId,
              icon: this.formData.icon,
              isNav: this.formData.isNav,
              sortOrder: this.formData.sortOrder,
              status: this.formData.status,
              translations: [
                {
                  name: this.formData.name,
                  slug: this.formData.slug,
                  description: this.formData.description,
                  languageCode: this.formData.lang,
                },
              ],
            })
          : categoryStore.updateCategory(this.formData.id || 0, {
              parentId: this.formData.parentId,
              icon: this.formData.icon,
              isNav: this.formData.isNav,
              sortOrder: this.formData.sortOrder,
              status: this.formData.status,
              translations: [
                {
                  name: this.formData.name,
                  slug: this.formData.slug,
                  description: this.formData.description,
                  languageCode: this.formData.lang,
                },
              ],
            }));

        // Clear draft after successful publish
        this.clearCategoryDraft();

        return '';
      } catch (error) {
        console.error('Failed to publish category:', error);
        return 'Failed to publish category';
      }
    },

    /**
     * Reset state
     */
    $reset() {
      console.log('resetting category edit view state');
      this.loading = false;
      this.needTranslate = false;
      this.isCreateMode = true;
      this.categoryId = null;
      this.formData = {
        name: '',
        slug: '',
        description: '',
        lang: 'zh-CN',
        sortOrder: 0,
        isNav: false,
      };
      this.languageOptions = [];
    },
  },
});
