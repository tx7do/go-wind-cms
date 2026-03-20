import React, {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {useCategoryStore} from '@/store/slices/category/hooks';
import type {contentservicev1_Category, contentservicev1_ListCategoryResponse} from '@/api/generated/app/service/v1';
import XIcon from '@/plugins/xicon';

import './index.scss';

interface CategoryFilterProps {
  categories?: contentservicev1_Category[];
  selectedCategory?: number | null;
  treeMode?: boolean;
  parentId?: number | null;
  autoLoad?: boolean;
  onCategoryChange?: (categoryId: number | null) => void;
  onLoaded?: (categories: contentservicev1_Category[]) => void;
}

const CategoryFilter: React.FC<CategoryFilterProps> = ({
                                                         categories: externalCategories,
                                                         selectedCategory = null,
                                                         treeMode = false,
                                                         parentId = null,
                                                         autoLoad = true,
                                                         onCategoryChange,
                                                         onLoaded
                                                       }) => {
  const {t} = useTranslation();
  const categoryStore = useCategoryStore();

  const [internalCategories, setInternalCategories] = useState<contentservicev1_Category[]>([]);
  const [loading, setLoading] = useState(false);
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set());

  const toggleExpanded = (nodeId: number) => {
    setExpandedIds(prev => {
      const next = new Set(prev);
      if (next.has(nodeId)) {
        next.delete(nodeId);
      } else {
        next.add(nodeId);
      }
      return next;
    });
  };

  const loadCategories = async () => {
    setLoading(true);
    try {
    const query: Record<string, any> = { status: 'CATEGORY_STATUS_ACTIVE' };

    if (parentId !== undefined && parentId !== null) {
    query.parentId = parentId;
    }

    // fieldMask and orderBy should use backend-expected snake_case
    const res = await categoryStore.listCategory({
    paging: undefined,
    formValues: query,
    fieldMask: 'id,status,sortOrder,icon,code,postCount,directPostCount,parentId,createdAt,children,translations.id,translations.categoryId,translations.name,translations.languageCode,translations.description',
    orderBy: ['-sortOrder']
    }) as unknown as contentservicev1_ListCategoryResponse;

      const items = res.items || [];
      setInternalCategories(items);
      onLoaded?.(items);
      console.log('[CategoryFilter] Categories loaded:', items.length);
    } catch (error) {
      console.error('[CategoryFilter] Load categories failed:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!autoLoad || externalCategories) {
      return;
    }

    loadCategories();
  }, [autoLoad, externalCategories, parentId]);

  function getCategoryName(category: contentservicev1_Category | null): string {
    if (!category?.id) return '';
    return categoryStore.getCategoryName(category, t);
  }

  const displayCategories = externalCategories || internalCategories;
  const rootCategories = displayCategories.filter(cat => !cat.parentId);

  const handleCategoryChange = (categoryId: number | null) => {
    onCategoryChange?.(categoryId);
  };

  const handleCategoryClick = (nodeId: number) => {
    handleCategoryChange(nodeId);
  };

  // 递归渲染分类项（支持多级）
  const renderCategoryItem = (node: contentservicev1_Category, level: number = 0, parentPath: string = '') => {
    const nodeId = node.id || 0;
    const uniqueKey = `${parentPath}-${nodeId}`; // 使用路径确保 key 唯一
    const isExpanded = expandedIds.has(nodeId);
    const hasChildNodes = node.children && node.children.length > 0;

    console.log('[renderCategoryItem]', {
      nodeId,
      level,
      hasChildNodes,
      isExpanded,
      treeMode,
      childrenCount: node.children?.length
    });

    return (
      <View key={uniqueKey} className='category-item-wrapper' style={{marginLeft: level > 0 ? `${level * 8}px` : undefined}}>
        <View
          className={`category-tab ${selectedCategory === nodeId ? 'active' : ''}`}
          onClick={(e) => {
            e.stopPropagation();
            console.log('[CategoryTab onClick]', { nodeId, hasChildNodes, treeMode, isExpanded });
            if (hasChildNodes && treeMode) {
              // 只有在树形模式且有子分类时，才切换展开/收起
              console.log('[toggleExpanded]', nodeId);
              toggleExpanded(nodeId);
            } else {
              // 否则选择该分类
              console.log('[handleCategoryClick]', nodeId);
              handleCategoryClick(nodeId);
            }
          }}
        >
          <XIcon name={node.icon || 'carbon:folder'} size={24} className='category-icon' />
          <Text>{getCategoryName(node)}</Text>
          {hasChildNodes && treeMode && (
            <XIcon
              name={isExpanded ? 'carbon:chevron-up' : 'carbon:chevron-down'}
              size={16}
              style={{marginLeft: 'auto', flexShrink: 0}}
            />
          )}
        </View>

        {/* 递归渲染子分类 - 只在树形模式且展开时显示 */}
        {hasChildNodes && treeMode && isExpanded && node.children && (
          <View className='category-submenu' style={{display: 'flex', flexDirection: 'column'}}>
            {node.children.map((child) => renderCategoryItem(child, level + 1, uniqueKey))}
          </View>
        )}
      </View>
    );
  };

  if (loading && autoLoad) {
    return <View className='loading'>{t('common.loading')}</View>;
  }

  return (
    <View className='category-filter'>
      <View className='category-tabs'>
        {/* 所有分类按钮 */}
        <View
          className={`category-tab ${selectedCategory === null ? 'active' : ''}`}
          onClick={() => handleCategoryChange(null)}
        >
          <XIcon name={'carbon:folder'} size={24} className='category-icon' />
          <Text>{t('page.posts.all_categories')}</Text>
        </View>

        {/* 树形模式 - 只渲染根分类，子分类递归渲染 */}
        {treeMode && rootCategories.map((node) => renderCategoryItem(node))}

        {/* 平铺模式 - 只显示根分类 */}
        {!treeMode && rootCategories.map((cat) => {
          return (
            <View
              key={cat.id}
              className={`category-tab ${selectedCategory === cat.id ? 'active' : ''}`}
              onClick={() => handleCategoryChange(cat.id || 0)}
            >
              <XIcon name={cat.icon || 'carbon:folder'} size={24} className='category-icon' />
              <Text>{getCategoryName(cat)}</Text>
            </View>
          );
        })}
      </View>
    </View>
  );
};

export default CategoryFilter;
