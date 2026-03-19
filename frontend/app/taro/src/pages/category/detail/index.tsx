import {useEffect, useMemo, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import Taro from '@tarojs/taro';
import {AppEmpty} from '@/components/ui';
import {useCategoryStore} from '@/store/slices/category/hooks';
import XIcon from '@/plugins/xicon';

import CategoryList from '@/components/category/CategoryList';
import PostListWithPagination from '@/components/post/PostList';
import {contentservicev1_Category} from "@/api/generated/app/service/v1";

import styles from './category-detail.scss';

export default function CategoryDetailPage() {
  const {t} = useTranslation('page');
  const categoryStore = useCategoryStore();
  const [childCategories, setChildCategories] = useState<contentservicev1_Category[]>([]);

  const categoryId = useMemo(() => {
    const pages = Taro.getCurrentPages();
    const currentPage = pages[pages.length - 1];
    const id = (currentPage as any).options?.id;
    return id ? parseInt(id) : null;
  }, []);

  // Get parent category ID
  const parentCategoryId = useMemo(() => {
    if (!categoryStore.detail?.parentId) return null;
    return categoryStore.detail.parentId;
  }, [categoryStore.detail?.parentId]);

  async function loadCategory() {
    if (!categoryId) return;

    try {
      await categoryStore.getCategory({
        id: categoryId,
        fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image'
      });
      // Extract child categories from the loaded category
      const loadedCategory = categoryStore.detail;
      if (loadedCategory?.children && loadedCategory.children.length > 0) {
        setChildCategories([...loadedCategory.children]);
      } else {
        setChildCategories([]);
      }
    } catch (error) {
      console.error('Load category failed:', error);
    }
  }

  function handleViewChildCategory(id: number) {
    Taro.navigateTo({url: `/pages/[locale]/category/${id}/page`});
  }

  function handleBackToParent() {
    if (parentCategoryId) {
      Taro.navigateTo({url: `/pages/[locale]/category/${parentCategoryId}/page`});
    } else {
      Taro.navigateBack();
    }
  }

  useEffect(() => {
    loadCategory();
  }, [categoryId]);

  if (!categoryId) {
    return <AppEmpty variant="error" description="Invalid category ID"/>;
  }

  return (
    <View className={styles['category-detail-page']}>
      {/* Hero Section */}
      <View className={styles['hero-section']}>
        <View className={styles['hero-content']}>
          <Text>{categoryStore.getCategoryName(categoryStore.detail)}</Text>
          {categoryStore.getCategoryDescription(categoryStore.detail) && (
            <Text className={styles['category-description']}>
              {categoryStore.getCategoryDescription(categoryStore.detail)}
            </Text>
          )}
          <View className={styles['category-stats']}>
            <View className={styles['stat-item']}>
              <XIcon name='carbon:document' size={20} />
              <Text>{categoryStore.detail?.postCount || 0} {t('page.posts.articles')}</Text>
            </View>
          </View>
        </View>
      </View>

      {/* Posts Section */}
      <View className={styles['page-container']}>
        {/* Back to Parent Button */}
        <View className={styles['back-button-container']}>
          <View className={styles['back-btn']} onClick={handleBackToParent}>
            <Text>← </Text>
            <Text>{parentCategoryId ? t('categories.back_to_parent') : t('categories.back_to_list')}</Text>
          </View>
        </View>

        {/* Sub Categories List */}
        {childCategories.length > 0 && (
          <View className={styles['sub-categories-wrapper']}>
            <CategoryList
              categories={childCategories}
              loading={false}
              showSkeleton={false}
              onCategoryClick={handleViewChildCategory}
            />
          </View>
        )}

        {/* Posts List with Pagination */}
        <PostListWithPagination
          key={categoryId}
          initialPageSize={10}
          pageSizes={[10, 20, 30, 40]}
          categoryId={categoryId}
          from="category"
        />
      </View>
    </View>
  );
}
