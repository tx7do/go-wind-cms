'use client';

import React, {useState} from 'react';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {
    getCategoryName,
    getCategoryDescription,
    getCategoryThumbnail,
} from '@/api/hooks/category';
import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import {cn} from '@/lib/utils';

interface CategoryTreeProps {
    categories: contentservicev1_Category[];
    level?: number;
    onCategoryClick?: (id: number) => void;
}

const levelMarginClass: Record<number, string> = {
    0: '',
    1: 'ml-10 max-md:ml-4',
    2: 'ml-20 max-md:ml-8',
    3: 'ml-30 max-md:ml-12',
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
        <div className="flex flex-col gap-6 max-md:gap-4">
            {categories.map((category) => (
                <div
                    key={category.id}
                    className={cn('flex flex-col', levelMarginClass[level] || '')}
                >
                    <div
                        className={cn(
                            'group flex cursor-pointer items-center justify-between overflow-hidden',
                            'rounded-2xl border border-border bg-card shadow-sm',
                            'transition-all duration-400',
                            'hover:-translate-y-1 hover:border-primary hover:shadow-lg',
                        )}
                        onClick={() => handleViewCategory(category.id || 0)}
                    >
                        <div className="flex flex-1 gap-5 p-5 max-md:flex-col max-md:p-4 max-md:gap-3">
                            <div className={cn(
                                'relative h-[120px] w-[160px] flex-shrink-0 overflow-hidden rounded-xl bg-background',
                                'max-md:h-[200px] max-md:w-full max-md:rounded-t-xl max-md:rounded-b-none',
                            )}>
                                <img
                                    src={getCategoryThumbnail(category)}
                                    alt={getCategoryName(category, t)}
                                    className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
                                />
                                <div className="absolute inset-0 bg-black/15 opacity-0 transition-opacity duration-300 group-hover:opacity-50"/>
                            </div>
                            <div className="flex flex-1 flex-col gap-2 max-md:w-full">
                                <h3 className={cn(
                                    'flex items-center gap-2 text-lg font-bold leading-tight text-foreground transition-colors',
                                    'group-hover:text-primary',
                                    'max-md:text-[17px]',
                                )}>
                                    {getCategoryName(category, t)}
                                </h3>
                                <p className="line-clamp-2 text-sm leading-relaxed text-muted-foreground max-md:text-[13px]">
                                    {getCategoryDescription(category)}
                                </p>
                                <div className="flex items-center gap-2 text-[13px] font-medium text-muted-foreground max-md:text-xs">
                                    <span className="flex items-center opacity-80">
                                        <XIcon name="carbon:document" size={16}/>
                                    </span>
                                    <span>
                                        {category.postCount || 0} {t('articles_count')}
                                    </span>
                                </div>
                            </div>
                        </div>

                        {/* 展开/收起按钮 (如果有子分类) */}
                        {category.children && category.children.length > 0 && (
                            <button
                                className={cn(
                                    'mr-3 flex items-center justify-center rounded-lg border border-border p-3',
                                    'text-muted-foreground transition-all duration-300',
                                    'hover:bg-primary/10 hover:border-primary hover:text-primary',
                                    'max-md:absolute max-md:top-3 max-md:right-3 max-md:z-10',
                                    'max-md:m-0 max-md:h-9 max-md:w-9 max-md:rounded-full max-md:p-0',
                                    'max-md:bg-card max-md:shadow-lg max-md:border-border',
                                    'max-md:hover:bg-primary max-md:hover:text-white max-md:hover:scale-110',
                                )}
                                onClick={(e) => toggleExpand(e, category)}
                            >
                                <XIcon
                                    name={isExpanded(category) ? 'carbon:chevron-down' : 'carbon:chevron-right'}
                                    size={20}
                                />
                            </button>
                        )}
                    </div>

                    {/* 递归渲染子分类 */}
                    {category.children && category.children.length > 0 && isExpanded(category) && (
                        <div className="overflow-hidden">
                            <CategoryTree
                                categories={category.children}
                                level={level + 1}
                                onCategoryClick={onCategoryClick}
                            />
                        </div>
                    )}
                </div>
            ))}
        </div>
    );
};

export default CategoryTree;
