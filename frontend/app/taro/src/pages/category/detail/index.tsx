import {useEffect, useMemo, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import Taro from '@tarojs/taro';

import {AppEmpty} from '@/components/ui';
import {fetchCategory, getCategoryName, getCategoryDescription} from '@/api/hooks/category';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import CategoryList from '@/components/category/CategoryList';
import PostListWithPagination from '@/components/post/PostList';
import {contentservicev1_Category} from "@/api/generated/app/service/v1";export default function CategoryDetailPage() {
  const {t} = useTranslation();

  const router = useI18nRouter();

  const [detail, setDetail] = useState<contentservicev1_Category | null>(null);
  const [childCategories, setChildCategories] = useState<contentservicev1_Category[]>([]);

  const categoryId = useMemo(() => {
    const pages = Taro.getCurrentPages();
    const currentPage = pages[pages.length - 1];
    const id = (currentPage as any).options?.id;
    return id ? parseInt(id) : null;
  }, []);

  // Get parent category ID
  const parentCategoryId = useMemo(() => {
    if (!detail?.parentId) return null;
    return detail.parentId;
  }, [detail?.parentId]);

  async function loadCategory() {
    if (!categoryId) return;

    try {
      const loadedCategory = await fetchCategory({
        id: categoryId,
        fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image'
      });
      setDetail(loadedCategory);
      // Extract child categories from the loaded category
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
    router.push(`/category/${id}`);
  }

  function handleBackToParent() {
    if (parentCategoryId) {
      router.push(`/category/${parentCategoryId}`);
    } else {
      router.back();
    }
  }

  useEffect(() => {
    loadCategory();
  }, [categoryId]);

  if (!categoryId) {
    return <AppEmpty variant='error' description='Invalid category ID' />;
  }

  return (
    <View className='category-detail-page'>
      {/* Hero Section */}
      <View className='hero-section'>
        <View className='hero-content'>
          <Text className='hero-title'>{getCategoryName(detail)}</Text>
          {getCategoryDescription(detail) && (
            <Text className='category-description'>
              {getCategoryDescription(detail)}
            </Text>
          )}
          <View className='category-stats'>
            <View className='stat-item'>
              <XIcon name='carbon:document' size={20} />
              <Text>{detail?.postCount || 0} {t('page.posts.articles')}</Text>
            </View>
          </View>
        </View>
      </View>

      {/* Posts Section */}
      <View className='page-container'>
        {/* Back to Parent Button */}
        <View className='back-button-container'>
          <View className='back-btn' onClick={handleBackToParent}>
            <Text>← </Text>
            <Text>{parentCategoryId ? t('page.categories.back_to_parent') : t('page.categories.back_to_list')}</Text>
          </View>
        </View>

        {/* Sub Categories List */}
        {childCategories.length > 0 && (
          <View className='sub-categories-wrapper'>
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
          from='category'
        />
      </View>
    </View>
  );
}
