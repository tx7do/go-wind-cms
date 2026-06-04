import {View, Text, ITouchEvent} from '@tarojs/components';
import React, {useState} from 'react';
import {useTranslations} from '@/lib/next-intl-compat';

import {XIcon} from '@/plugins/xicon';
import {
    getCategoryName,
    getCategoryDescription,
    getCategoryThumbnail,
} from '@/api/hooks/category';
import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import {cn} from '@/lib/utils';
import Image from '@/components/ui/image';

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
    const t = useTranslations('page.categories');
    const [expandedCategories, setExpandedCategories] = useState<Set<number>>(new Set());

    const handleViewCategory = (id: number) => {
        onCategoryClick?.(id);
    };

    const toggleExpand = (e: React.MouseEvent | ITouchEvent, category: contentservicev1_Category) => {
        e.stopPropagation();
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
        <View className='flex flex-col gap-[24rpx]'>
            {categories.map((category) => {
                const hasChildren = !!(category.children && category.children.length > 0);
                const expanded = isExpanded(category);

                return (
                    <View
                      key={category.id}
                      className='flex flex-col'
                      style={{marginLeft: level > 0 ? `${level * 24}px` : 0}}
                    >
                        <View
                          className={cn(
                                'flex items-stretch overflow-hidden',
                                'rounded-[16rpx] bg-cardBg',
                                expanded && 'border-[1rpx] border-primary/30',
                            )}
                          onClick={() => handleViewCategory(category.id || 0)}
                          hoverClass='opacity-80'
                        >
                            {/* 左侧色条 */}
                            <View
                              className={cn(
                                    'w-[6rpx] flex-shrink-0',
                                    expanded ? 'bg-primary' : 'bg-transparent',
                                )}
                            />

                            <View className='flex items-center gap-[24rpx] p-[24rpx] flex-1 min-w-0'>
                                {/* 缩略图 */}
                                <View className='w-[180rpx] h-[120rpx] flex-shrink-0 overflow-hidden rounded-[12rpx] bg-splitLine'>
                                    <Image
                                      src={getCategoryThumbnail(category)}
                                      alt={getCategoryName(category, t)}
                                      className='w-full h-full object-cover'
                                    />
                                </View>

                                {/* 文字区 */}
                                <View className='flex flex-1 flex-col gap-[8rpx] min-w-0'>
                                    <Text
                                      className='text-body font-bold text-textMain'
                                      style={{
                                        overflow: 'hidden',
                                        textOverflow: 'ellipsis',
                                        whiteSpace: 'nowrap',
                                      }}
                                    >
                                        {getCategoryName(category, t)}
                                    </Text>
                                    {getCategoryDescription(category) && (
                                        <Text
                                          className='text-tips text-textSec'
                                          style={{
                                            overflow: 'hidden',
                                            textOverflow: 'ellipsis',
                                            display: '-webkit-box',
                                            WebkitLineClamp: 2,
                                            WebkitBoxOrient: 'vertical',
                                          }}
                                        >
                                            {getCategoryDescription(category)}
                                        </Text>
                                    )}
                                    <View className='flex items-center gap-[16rpx] mt-[4rpx]'>
                                        <View className='flex items-center gap-[6rpx]'>
                                            <XIcon name='carbon:document' size={12} className='text-textThird' />
                                            <Text className='text-tips text-textSec'>
                                                {category.postCount || 0} {t('articles_count')}
                                            </Text>
                                        </View>
                                        {hasChildren && (
                                            <View className='flex items-center gap-[6rpx]'>
                                                <XIcon name='carbon:folder' size={12} className='text-textThird' />
                                                <Text className='text-tips text-textSec'>
                                                    {category.children!.length}
                                                </Text>
                                            </View>
                                        )}
                                    </View>
                                </View>

                                {/* 右侧箭头或展开按钮 */}
                                {!hasChildren ? (
                                    <View className='flex-shrink-0'>
                                        <XIcon name='carbon:chevron-right' size={18} className='text-textThird' />
                                    </View>
                                ) : (
                                    <View
                                      className={cn(
                                            'flex items-center justify-center w-[56rpx] h-[56rpx] rounded-full',
                                            'border-[1rpx] border-splitLine',
                                            expanded && 'bg-primary/10 border-primary/30',
                                        )}
                                      onClick={(e) => toggleExpand(e, category)}
                                    >
                                        <XIcon
                                          name={expanded ? 'carbon:chevron-down' : 'carbon:chevron-right'}
                                          size={16}
                                          className={expanded ? 'text-primary' : 'text-textSec'}
                                        />
                                    </View>
                                )}
                            </View>
                        </View>

                        {/* 递归渲染子分类 */}
                        {hasChildren && expanded && (
                            <View className='mt-[16rpx] ml-[32rpx] border-l-[2rpx] border-primary/20 pl-[16rpx]'>
                                <CategoryTree
                                  categories={category.children!}
                                  level={0}
                                  onCategoryClick={onCategoryClick}
                                />
                            </View>
                        )}
                    </View>
                );
            })}
        </View>
    );
};

export default CategoryTree;
