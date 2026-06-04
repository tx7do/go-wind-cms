import {View} from '@tarojs/components';
import React from 'react';
import {Skeleton} from '@/components/ui/skeleton';
import {useTranslations} from '@/lib/next-intl-compat';
import {AppEmpty} from '@/components/ui';

import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

import CategoryCard from './CategoryCard';

interface CategoryListProps {
    categories: contentservicev1_Category[];
    loading?: boolean;
    showSkeleton?: boolean;
    columns?: number;
    onCategoryClick?: (id: number) => void;
}

const CategoryList: React.FC<CategoryListProps> = ({
                                                       categories = [],
                                                       loading = false,
                                                       showSkeleton = true,
                                                       columns = 2,
                                                       onCategoryClick
                                                   }) => {
    const t = useTranslations('page.categories');

    const handleCategoryClick = (id: number) => {
        onCategoryClick?.(id);
    };

    return (
        <View className='w-full'>
            {/* Loading Skeleton */}
            {loading && showSkeleton && (
                <View
                  className='grid grid-cols-2 gap-[24rpx]'
                >
                    {Array.from({length: columns * 2}).map((_, index) => (
                        <View key={index} className='overflow-hidden rounded-[16rpx] bg-cardBg'>
                            <Skeleton className='h-[180rpx] w-full' />
                            <View className='p-[24rpx]'>
                                <Skeleton className='h-[32rpx] w-full rounded' />
                                <Skeleton className='h-[24rpx] w-[80%] rounded mt-[12rpx]' />
                                <Skeleton className='h-[24rpx] w-[40%] rounded mt-[12rpx]' />
                            </View>
                        </View>
                    ))}
                </View>
            )}

            {/* Loaded Content */}
            {!loading && categories.length > 0 && (
                <View className='grid grid-cols-2 gap-[24rpx]'>
                    {categories.map((category) => (
                        <CategoryCard
                          key={category.id}
                          category={category}
                          onClick={handleCategoryClick}
                        />
                    ))}
                </View>
            )}

            {/* Empty State */}
            {!loading && categories.length === 0 && (
                <AppEmpty description={t('no_categories')} />
            )}
        </View>
    );
};

export default CategoryList;
