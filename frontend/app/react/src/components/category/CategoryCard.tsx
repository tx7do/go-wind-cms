'use client';

import React from 'react';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {
    getCategoryName,
    getCategoryDescription,
    getCategoryThumbnail,
} from '@/api/hooks/category';
import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import {cn} from '@/lib/utils';

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
    const t = useTranslations('page.categories');

    const handleClick = () => {
        if (!category?.id || !clickable) return;
        onClick?.(category.id);
    };

    if (!category) return null;

    return (
        <div
            className={cn(
                'group flex h-full flex-col overflow-hidden rounded-xl border border-border bg-background',
                'transition-all duration-300',
                clickable && 'cursor-pointer hover:-translate-y-1.5 hover:border-primary hover:shadow-lg',
                !clickable && 'cursor-default',
            )}
            onClick={handleClick}
        >
            <div className="relative h-[160px] w-full overflow-hidden bg-background">
                <img
                    src={getCategoryThumbnail(category)}
                    alt={getCategoryName(category, t)}
                    className={cn(
                        'h-full w-full object-cover transition-transform duration-500',
                        clickable && 'group-hover:scale-110',
                    )}
                />
                <div className={cn(
                    'absolute inset-0 bg-black/15 transition-opacity duration-300',
                    clickable ? 'opacity-0 group-hover:opacity-50' : 'opacity-0',
                )}/>
            </div>
            <div className="p-4">
                <h3 className={cn(
                    'mb-2 line-clamp-2 text-base font-bold leading-tight text-foreground transition-colors',
                    clickable && 'group-hover:text-primary',
                )}>
                    {getCategoryName(category, t)}
                </h3>
                <p className="mb-3 line-clamp-2 text-[13px] leading-relaxed text-muted-foreground">
                    {getCategoryDescription(category)}
                </p>
                <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
                    <span className="flex items-center opacity-80">
                        <XIcon name="carbon:document" size={14}/>
                    </span>
                    <span>
                        {category.postCount || 0} {t('articles_count')}
                    </span>
                </div>
            </div>
        </div>
    );
};

export default CategoryCard;
