import React from 'react';
import {useTranslation} from 'react-i18next';
import {Text, View} from '@tarojs/components';

import {getCategoryName} from '@/api/hooks/category';
import {contentservicev1_Category} from '@/api/generated/app/service/v1';
import XIcon from '@/plugins/xicon';

import './index.scss';

interface HomeCategoryCardProps {
  category: contentservicev1_Category;
  onClick?: (id: number) => void;
}

const HomeCategoryCard: React.FC<HomeCategoryCardProps> = ({category, onClick}) => {
  const {t} = useTranslation();

  // 计算更新时间（假设 category.createdAt 是 ISO 字符串或时间戳）
  const getDaysAgo = () => {
    if (!category.createdAt) return 0;
    const createdDate = new Date(category.createdAt);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - createdDate.getTime());
    return Math.floor(diffTime / (1000 * 60 * 60 * 24));
  };

  const daysAgo = getDaysAgo();

  const handleClick = () => {
    if (!category?.id) return;
    onClick?.(category.id);
  };

  return (
    <View
      className='home-category-card scroll-reveal-item'
      onClick={handleClick}
    >
      <View className='home-category-card-bg' />
      <View className='home-category-card-content'>
        <View className='home-category-card-header'>
          <View className='home-category-icon'>
            <XIcon name={category.icon || 'carbon:folder'} size={32} className='icon-placeholder' />
          </View>
          <View className='home-category-info'>
            <Text className='category-name'>{getCategoryName(category)}</Text>
            <Text className='home-post-count'>
              {t('page.home.article_count', {count: category.postCount || 0})}
            </Text>
          </View>
        </View>
        <View className='home-category-badge'>
          <XIcon name='carbon:time' size={14} style={{marginRight: '4px'}} />
          <Text>{t('page.home.updated_days_ago', {days: daysAgo})}</Text>
        </View>
      </View>
    </View>
  );
};

export default HomeCategoryCard;
