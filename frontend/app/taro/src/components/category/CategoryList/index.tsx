import React from 'react';
import {useTranslation} from 'react-i18next';
import {View} from '@tarojs/components';

import {AppEmpty} from '@/components/ui';

import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import CategoryCard from '../CategoryCard';

import './index.scss';


interface CategoryListProps {
  categories: contentservicev1_Category[];
  loading?: boolean;
  showSkeleton?: boolean;
  columns?: number;
  gap?: number;
  onCategoryClick?: (id: number) => void;
}

const CategoryList: React.FC<CategoryListProps> = ({
                                                     categories = [],
                                                     loading = false,
                                                     showSkeleton = true,
                                                     columns = 3,
                                                     gap = 20,
                                                     onCategoryClick
                                                   }) => {
  const {t} = useTranslation();

  const handleCategoryClick = (id: number) => {
    onCategoryClick?.(id);
  };

  return (
    <View className='category-list-container'>
      {/* Loading Skeleton */}
      {loading && showSkeleton && (
        <View
          className='category-grid'
          style={{
            gridTemplateColumns: `repeat(${columns}, 1fr)`,
            gap: `${gap}px`
          }}
        >
          {Array.from({length: columns}).map((_, index) => (
            <View key={index} className='category-card-skeleton'>
              <View className='skeleton-image-placeholder' />
              <View className='skeleton-content'>
                <View className='skeleton-title-placeholder' />
                <View className='skeleton-paragraph-placeholder' />
                <View className='skeleton-btn-small' />
              </View>
            </View>
          ))}
        </View>
      )}

      {/* Loaded Content */}
      {!loading && categories.length > 0 && (
        <View
          className='category-grid'
          style={{
            gridTemplateColumns: `repeat(auto-fill, minmax(280px, 1fr))`,
            gap: `${gap}px`
          }}
        >
          {categories.map((category) => (
            <CategoryCard
              key={category.id}
              category={category}
              onClick={handleCategoryClick}
            />
          ))}
        </View>
      )}

      {/* Empty State */}
      {!loading && categories.length === 0 && (
        <AppEmpty description={t('page.categories.no_categories')} />
      )}
    </View>
  );
};

export default CategoryList;
