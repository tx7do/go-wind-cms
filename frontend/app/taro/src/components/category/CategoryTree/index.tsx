import React, {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';

import {useCategoryStore} from '@/store/slices/category/hooks';
import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import './index.scss';

interface CategoryTreeProps {
    categories: contentservicev1_Category[];
    level?: number;
    onCategoryClick?: (id: number) => void;
}

const CategoryTree: React.FC<CategoryTreeProps> = ({
  categories = [],
  level = 0,
  onCategoryClick
}) => {
  const {t} = useTranslation('page.categories');
  const categoryStore = useCategoryStore();
  const [expandedCategories, setExpandedCategories] = useState<Set<number>>(new Set());

  const handleViewCategory = (id: number) => {
    onCategoryClick?.(id);
  };

  const toggleExpand = (category: contentservicev1_Category) => {
    if (category.children && category.children.length > 0) {
      const newExpanded = new Set(expandedCategories);
      if (newExpanded.has(category.id || 0)) {
        newExpanded.delete(category.id || 0);
      } else {
        newExpanded.add(category.id || 0);
      }
      setExpandedCategories(newExpanded);
    }
  };

  const isExpanded = (category: contentservicev1_Category) => {
    return expandedCategories.has(category.id || 0);
  };

  if (!categories || categories.length === 0) return null;

  return (
    <View className='category-tree'>
      {categories.map((category) => (
        <View
          key={category.id}
          className={`category-node level-${level}`}
        >
          <View
            className='category-item'
            onClick={() => handleViewCategory(category.id || 0)}
          >
            <View className='category-info'>
              <View className='category-image'>
                <Image
                  src={categoryStore.getCategoryThumbnail(category)}
                  mode='aspectFill'
                  className='category-img'
                />
                <View className='image-overlay' />
              </View>
              <View className='category-content'>
                <Text className={`category-title ${category.children && category.children.length > 0 ? 'has-children' : ''}`}>
                  {categoryStore.getCategoryName(category, t)}
                </Text>
                <Text className='description'>
                  {categoryStore.getCategoryDescription(category)}
                </Text>
                <View className='category-meta'>
                  <Text className='meta-icon'>📄</Text>
                  <Text className='meta-text'>
                    {category.postCount || 0} {t('articles_count')}
                  </Text>
                </View>
              </View>
            </View>

            {/* 展开/收起按钮 (如果有子分类) */}
            {category.children && category.children.length > 0 && (
              <View
                className={`expand-btn ${isExpanded(category) ? 'expanded' : ''}`}
                onClick={(e) => {
                  e.stopPropagation();
                  toggleExpand(category);
                }}
              >
                <Text>{isExpanded(category) ? '▼' : '▶'}</Text>
              </View>
            )}
          </View>

          {/* 递归渲染子分类 */}
          {category.children && category.children.length > 0 && isExpanded(category) && (
            <View className='slide-container'>
              <CategoryTree
                categories={category.children}
                level={level + 1}
                onCategoryClick={onCategoryClick}
              />
            </View>
          )}
        </View>
      ))}
    </View>
  );
};

export default CategoryTree;
