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
        <div className="flex flex-col gap-4 max-md:gap-3">
            {categories.map((category) => {
                const hasChildren = !!(category.children && category.children.length > 0);
                const expanded = isExpanded(category);

                return (
                    <div
                        key={category.id}
                        className={cn('flex flex-col', levelMarginClass[level] || '')}
                    >
                        <div
                            className={cn(
                                'group relative flex cursor-pointer items-stretch overflow-hidden',
                                'rounded-xl border border-border bg-card',
                                'transition-all duration-300',
                                'hover:border-primary/50 hover:shadow-md',
                                expanded && 'border-primary/40 shadow-sm',
                            )}
                            onClick={() => handleViewCategory(category.id || 0)}
                        >
                            {/* 左侧色条指示器 */}
                            <div className={cn(
                                'w-1 flex-shrink-0 bg-primary/0 transition-colors duration-300',
                                'group-hover:bg-primary',
                                expanded && 'bg-primary',
                            )}/>

                            <div className="flex flex-1 gap-4 p-4 max-md:flex-col max-md:p-3 max-md:gap-3">
                                <div className={cn(
                                    'relative h-[100px] w-[140px] flex-shrink-0 overflow-hidden rounded-lg bg-muted',
                                    'max-md:h-[180px] max-md:w-full',
                                )}>
                                    <img
                                        src={getCategoryThumbnail(category)}
                                        alt={getCategoryName(category, t)}
                                        className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
                                    />
                                    <div className="absolute inset-0 bg-gradient-to-t from-black/30 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100"/>
                                </div>
                                <div className="flex flex-1 flex-col justify-center gap-1.5 max-md:w-full">
                                    <h3 className={cn(
                                        'flex items-center gap-2 text-base font-semibold leading-tight text-foreground transition-colors',
                                        'group-hover:text-primary',
                                    )}>
                                        {getCategoryName(category, t)}
                                    </h3>
                                    <p className="line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                                        {getCategoryDescription(category)}
                                    </p>
                                    <div className="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
                                        <span className="flex items-center gap-1">
                                            <XIcon name="carbon:document" size={14}/>
                                            {category.postCount || 0} {t('articles_count')}
                                        </span>
                                        {hasChildren && (
                                            <span className="flex items-center gap-1">
                                                <XIcon name="carbon:folder" size={14}/>
                                                {category.children!.length}
                                            </span>
                                        )}
                                    </div>
                                </div>
                            </div>

                            {/* 展开/收起按钮 */}
                            {hasChildren && (
                                <button
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
                                </button>
                            )}
                        </div>

                        {/* 递归渲染子分类 */}
                        {hasChildren && expanded && (
                            <div className={cn(
                                'mt-2 overflow-hidden',
                                'border-l-2 border-primary/20 ml-3',
                            )}>
                                <CategoryTree
                                    categories={category.children!}
                                    level={level + 1}
                                    onCategoryClick={onCategoryClick}
                                />
                            </div>
                        )}
                    </div>
                );
            })}
        </div>
    );
};

export default CategoryTree;
