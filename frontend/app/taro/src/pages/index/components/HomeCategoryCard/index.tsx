import React from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {useCategoryStore} from '@/store/slices/category/hooks';
import {contentservicev1_Category} from '@/api/generated/app/service/v1';
import XIcon from '@/plugins/xicon';

import './index.scss';

interface HomeCategoryCardProps {
  category: contentservicev1_Category;
  onClick?: (id: number) => void;
}

const HomeCategoryCard: React.FC<HomeCategoryCardProps> = ({category, onClick}) => {
  const {t} = useTranslation();
  const categoryStore = useCategoryStore();

  const handleClick = () => {
    if (!category?.id) return;
    onClick?.(category.id);
  };

  return (
    <View
      className='home-category-card scroll-reveal-item'
      onClick={handleClick}
    >
      <View className='home-category-card-bg'/>
      <View className='home-category-card-content'>
        <View className='home-category-card-header'>
          <View className='home-category-icon'>
            <XIcon name={category.icon || 'carbon:folder'} size={32} className='icon-placeholder' />
          </View>
          <View className='home-category-info'>
            <Text className='category-name'>{categoryStore.getCategoryName(category)}</Text>
            <Text className='home-post-count'>
              {t('page.home.article_count', {count: category.postCount || 0})}
            </Text>
          </View>
        </View>
        <View className='home-category-badge'>
          <Text>⏰ {t('page.home.updated_days_ago', {days: 3})}</Text>
        </View>
      </View>
    </View>
  );
};

export default HomeCategoryCard;
