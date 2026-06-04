import {View, Text} from '@tarojs/components';
import React from 'react';
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
        <View
          className={cn(
                'flex flex-col overflow-hidden rounded-[16rpx] bg-cardBg',
                clickable && 'tap-active',
            )}
          onClick={handleClick}
          hoverClass={clickable ? 'opacity-80' : undefined}
        >
            {/* 封面图 */}
            <View className='relative w-full h-[180rpx] overflow-hidden'>
                <Image
                  src={getCategoryThumbnail(category)}
                  alt={getCategoryName(category, t)}
                  className='w-full h-full object-cover'
                />
                <View className='absolute inset-0 bg-gradient-to-t from-black/20 to-transparent' />
            </View>

            {/* 内容区域 */}
            <View className='flex flex-col p-[24rpx]'>
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
                      className='text-tips text-textSec mt-[8rpx]'
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

                <View className='flex items-center gap-[8rpx] mt-[16rpx]'>
                    <XIcon name='carbon:document' size={12} className='text-textThird' />
                    <Text className='text-tips text-textSec'>
                        {category.postCount || 0} {t('articles_count')}
                    </Text>
                </View>
            </View>
        </View>
    );
};

export default CategoryCard;
