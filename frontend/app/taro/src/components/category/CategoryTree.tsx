import {View, Text} from '@tarojs/components';
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

const levelMarginClass: Record<number, string> = {
    0: '',
    1: 'ml-8 max-md:ml-4',
    2: 'ml-16 max-md:ml-8',
    3: 'ml-24 max-md:ml-12',
};

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

    const toggleExpand = (e: React.MouseEvent, category: contentservicev1_Category) => {
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
        <View className='flex flex-col gap-4 max-md:gap-3'>
            {categories.map((category) => {
                const hasChildren = !!(category.children && category.children.length > 0);
                const expanded = isExpanded(category);

                return (
                    <View
                      key={category.id}
                      className={cn('flex flex-col', levelMarginClass[level] || '')}
                    >
                        <View
                          className={cn(
                                'group relative flex cursor-pointer items-center overflow-hidden',
                                'rounded-xl border border-border/60 bg-card/50',
                                'transition-all duration-300',
                                'hover:bg-card hover:border-primary/30 hover:shadow-[0_8px_24px_-8px_hsl(var(--primary)/0.12)]',
                                expanded && 'border-primary/40 bg-card shadow-sm',
                            )}
                          onClick={() => handleViewCategory(category.id || 0)}
                        >
                            {/* 左侧色条指示器 */}
                            <View className={cn(
                                'w-1 self-stretch flex-shrink-0 bg-primary/0 transition-colors duration-300',
                                'group-hover:bg-primary',
                                expanded && 'bg-primary',
                            )}
                            />

                            <View className='flex flex-1 items-center gap-4 p-4 max-md:flex-col max-md:items-start max-md:p-3 max-md:gap-3'>
                                {/* 缩图：rounded-lg + 微间距，像精致狐裹的组件 */}
                                <View className={cn(
                                    'relative h-[90px] w-[130px] flex-shrink-0 overflow-hidden rounded-lg bg-muted',
                                    'max-md:h-[180px] max-md:w-full',
                                )}
                                >
                                    <Image
                                      src={getCategoryThumbnail(category)}
                                      alt={getCategoryName(category, t)}
                                      className='h-full w-full object-cover transition-transform duration-500 group-hover:scale-105'
                                    />
                                    <View className='absolute inset-0 bg-gradient-to-t from-black/30 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100' />
                                </View>
                                {/* 文字区：垂直居中 */}
                                <View className='flex flex-1 flex-col gap-1.5 max-md:w-full'>
                                    <Text className={cn(
                                        'flex items-center gap-2 text-base font-semibold leading-tight text-foreground transition-colors',
                                        'group-hover:text-primary',
                                    )}
                                    >
                                        {getCategoryName(category, t)}
                                    </Text>
                                    <Text className='line-clamp-2 text-sm leading-relaxed text-muted-foreground'>
                                        {getCategoryDescription(category)}
                                    </Text>
                                    <View className='mt-1 flex items-center gap-3 text-xs text-muted-foreground'>
                                        <Text className='flex items-center gap-1'>
                                            <XIcon name='carbon:document' size={14} />
                                            {category.postCount || 0} {t('articles_count')}
                                        </Text>
                                        {hasChildren && (
                                            <Text className='flex items-center gap-1'>
                                                <XIcon name='carbon:folder' size={14} />
                                                {category.children!.length}
                                            </Text>
                                        )}
                                    </View>
                                </View>

                                {/* 右侧箭头：无子分类时显示 */}
                                {!hasChildren && (
                                    <View className='flex-shrink-0 text-muted-foreground/50 transition-all duration-300 group-hover:text-primary group-hover:translate-x-1 max-md:hidden'>
                                        <XIcon name='carbon:chevron-right' size={20} />
                                    </View>
                                )}
                            </View>

                            {/* 展开/收起按钮 */}
                            {hasChildren && (
                                <View
                                  className={cn(
                                        'flex items-center justify-center px-3',
                                        'border-l border-border text-muted-foreground',
                                        'transition-colors duration-200',
                                        'hover:bg-primary/5 hover:text-primary',
                                        'max-md:absolute max-md:top-2 max-md:right-2 max-md:z-10',
                                        'max-md:h-8 max-md:w-8 max-md:rounded-full max-md:border max-md:border-border',
                                        'max-md:bg-card/90 max-md:shadow-sm max-md:border-l',
                                        'max-md:hover:bg-primary max-md:hover:text-white max-md:hover:scale-110',
                                    )}
                                  onClick={(e) => toggleExpand(e, category)}
                                >
                                    <XIcon
                                      name={expanded ? 'carbon:chevron-down' : 'carbon:chevron-right'}
                                      size={20}
                                      className={cn('transition-transform duration-300', expanded && 'rotate-180')}
                                    />
                                </View>
                            )}
                        </View>

                        {/* 递归渲染子分类 */}
                        {hasChildren && expanded && (
                            <View className={cn(
                                'mt-2 overflow-hidden',
                                'border-l-2 border-primary/20 ml-3',
                            )}
                            >
                                <CategoryTree
                                  categories={category.children!}
                                  level={level + 1}
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
