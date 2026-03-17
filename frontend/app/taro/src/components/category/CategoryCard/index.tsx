import React from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';

import {useCategoryStore} from '@/store/slices/category/hooks';
import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import './index.scss';

interface CategoryCardProps {
  category: contentservicev1_Category | null;
  clickable?: boolean;
  onClick?: (id: number) => void;
}

const CategoryCard: React.FC<CategoryCardProps> = ({
                                                     category,
                                                     clickable = true,
                                                     onClick
                                                   }) => {
  const {t} = useTranslation('page.categories');
  const categoryStore = useCategoryStore();

  const handleClick = () => {
    if (!category?.id || !clickable) return;
    onClick?.(category.id);
  };

  if (!category) return null;

  return (
    <View
      className={`category-card ${clickable ? 'clickable' : ''}`}
      onClick={handleClick}
    >
      <View className='category-card-image'>
        <Image
          src={categoryStore.getCategoryThumbnail(category)}
          mode='aspectFill'
          className='category-image-img'
        />
        <View className='image-overlay'/>
      </View>
      <View className='category-card-content'>
        <Text className='category-title'>{categoryStore.getCategoryName(category, t)}</Text>
        <Text className='category-description'>{categoryStore.getCategoryDescription(category)}</Text>
        <View className='category-card-meta'>
          <Text className='meta-icon'>
            📄
          </Text>
          <Text className='meta-text'>
            {category.postCount || 0} {t('articles_count')}
          </Text>
        </View>
      </View>
    </View>
  );
};

export default CategoryCard;
