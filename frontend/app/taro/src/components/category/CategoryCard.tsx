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
                'group flex h-full flex-col overflow-hidden rounded-xl border border-border bg-background',
                /* 空气动力学悬浮：纯 primary 辉光 */
                'transition-all duration-500 ease-out',
                clickable && 'cursor-pointer hover:-translate-y-1.5 hover:border-primary/40 hover:shadow-[0_20px_40px_-8px_hsl(var(--primary)/0.15)]',
                !clickable && 'cursor-default',
            )}
          onClick={handleClick}
        >
            <View className='relative h-[160px] w-full overflow-hidden bg-background'>
                <Image
                  src={getCategoryThumbnail(category)}
                  alt={getCategoryName(category, t)}
                  className={cn(
                        'h-full w-full object-cover transition-transform duration-700 ease-out',
                        clickable && 'group-hover:scale-[1.12]',
                    )}
                />
                {/* hover 风迹渐变蒙层 */}
                <View className={cn(
                    'absolute inset-0 bg-gradient-to-tr from-primary/20 via-transparent to-sky-400/20 transition-opacity duration-500',
                    clickable ? 'opacity-0 group-hover:opacity-100' : 'opacity-0',
                )}
                />
            </View>
            <View className='p-4'>
                <Text className={cn(
                    'mb-2 line-clamp-2 text-base font-bold leading-tight text-foreground transition-colors',
                    clickable && 'group-hover:text-primary',
                )}
                >
                    {getCategoryName(category, t)}
                </Text>
                <Text className='mb-3 line-clamp-2 text-[13px] leading-relaxed text-muted-foreground'>
                    {getCategoryDescription(category)}
                </Text>
                <View className='flex items-center gap-1.5 text-xs font-medium text-muted-foreground'>
                    <Text className='flex items-center'>
                        <XIcon name='carbon:document' size={14} />
                    </Text>
                    <Text>
                        {category.postCount || 0} {t('articles_count')}
                    </Text>
                </View>
            </View>
        </View>
    );
};

export default CategoryCard;
