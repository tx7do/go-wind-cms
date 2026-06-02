'use client';

import React from 'react';
import {Skeleton} from '@/components/ui/skeleton';
import {useTranslations} from 'next-intl';
import {AppEmpty} from '@/components/ui';

import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import CategoryCard from './CategoryCard';


interface CategoryListProps {
    categories: contentservicev1_Category[];
    loading?: boolean;
    showSkeleton?: boolean;
    columns?: number;
    gap?: number;
    onCategoryClick?: (id: number) => void;
}

const CategoryList: React.FC<CategoryListProps> = ({
                                                       categories = [],
                                                       loading = false,
                                                       showSkeleton = true,
                                                       columns = 3,
                                                       gap = 20,
                                                       onCategoryClick
                                                   }) => {
    const t = useTranslations('page.categories');

    const handleCategoryClick = (id: number) => {
        onCategoryClick?.(id);
    };

    return (
        <div className="w-full">
            {/* Loading Skeleton */}
            {loading && showSkeleton && (
                <div
                    className="grid max-md:grid-cols-1"
                    style={{
                        gridTemplateColumns: `repeat(${columns}, 1fr)`,
                        gap: `${gap}px`
                    }}
                >
                    {Array.from({length: columns}).map((_, index) => (
                        <div key={index} className="overflow-hidden rounded-xl border border-border bg-background">
                            <Skeleton className="h-40 w-full"/>
                            <div className="p-4">
                                <Skeleton className="h-4 w-full"/>
                                <Skeleton className="h-4 w-3/4"/>
                                <Skeleton className="h-6 w-[60px]"/>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Loaded Content */}
            {!loading && categories.length > 0 && (
                <div
                    className="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] max-md:grid-cols-1"
                    style={{gap: `${gap}px`}}
                >
                    {categories.map((category) => (
                        <CategoryCard
                            key={category.id}
                            category={category}
                            onClick={handleCategoryClick}
                        />
                    ))}
                </div>
            )}

            {/* Empty State */}
            {!loading && categories.length === 0 && (
                <AppEmpty description={t('no_categories')}/>
            )}
        </div>
    );
};

export default CategoryList;
