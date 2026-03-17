import {useState, useEffect, useCallback} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {useCategoryStore} from '@/store/slices/category/hooks';
import type {contentservicev1_Category, contentservicev1_ListCategoryResponse} from '@/api/generated/app/service/v1';

import HomeCategoryCard from '../HomeCategoryCard';

import './index.scss';

interface CategoryListSectionProps {
  skeletonCount?: number;
  pageSize?: number;
  page?: number;
  filter?: Record<string, unknown>;
  orderBy?: string[];
  fieldMask?: string;
  showHeader?: boolean;
}

export default function CategoryListSection({
                                              skeletonCount = 8,
                                              pageSize = 8,
                                              page = 1,
                                              filter = {status: 'CATEGORY_STATUS_ACTIVE'} as Record<string, unknown>,
                                              orderBy = ['-sortOrder', '-postCount'],
                                              fieldMask = 'id,status,sortOrder,icon,code,postCount,directPostCount,parent_id,createdAt,translations.id,translations.categoryId,translations.name,translations.languageCode,translations.description',
                                              showHeader = true
                                            }: CategoryListSectionProps) {
  const {t} = useTranslation('page.home');
  const categoryStore = useCategoryStore();

  const [categories, setCategories] = useState<contentservicev1_Category[]>([]);
  const [loading, setLoading] = useState(false);

  const loadCategories = useCallback(async () => {
    setLoading(true);
    try {
      const res = await categoryStore.listCategory({
        paging: {page, pageSize},
        formValues: filter,
        fieldMask,
        orderBy: orderBy,
      }) as unknown as contentservicev1_ListCategoryResponse;

      setCategories(res.items || []);
    } catch (error) {
      console.error('Failed to load categories:', error);
      setCategories([]);
    } finally {
      setLoading(false);
    }
  }, [categoryStore, filter, fieldMask, orderBy, page, pageSize]);

  useEffect(() => {
    loadCategories();
  }, [loadCategories]);

  const router = {push: (path: string) => console.log('Navigate to:', path)}; // TODO: 使用真实 router

  const handleViewCategoryDetail = (id: number) => {
    router.push(`/category/${id}`);
  };

  const handleViewCategory = () => {
    router.push(`/category`);
  };

  return (
    <View className='categories-section scroll-reveal'>
      {/* Section Header */}
      {showHeader && (
        <View className='section-header'>
          <Text className='section-title'>
            📁 {t('categories')}
          </Text>
          <View
            className='view-all-btn'
            onClick={handleViewCategory}
          >
            <Text>{t('view_all')} →</Text>
          </View>
        </View>
      )}

      {/* Loading Skeleton */}
      {loading ? (
        <View className='categories-grid desktop-grid'>
          {Array.from({length: skeletonCount}).map((_, i) => (
            <View key={i} className='category-card-skeleton'>
              <View className='skeleton-placeholder'/>
            </View>
          ))}
        </View>
      ) : (
        <View className='categories-grid desktop-grid'>
          {categories.map((category) => (
            <HomeCategoryCard
              key={category.id}
              category={category}
              onClick={handleViewCategoryDetail}
            />
          ))}
        </View>
      )}
    </View>
  );
}
