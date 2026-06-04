import {View, Text} from '@tarojs/components';
﻿import React from 'react';
import {XIcon} from '@/plugins/xicon';
import {getCategoryName} from '@/api/hooks/category';
import {contentservicev1_Category} from '@/api/generated/app/service/v1';
import {useTranslations} from '@/lib/next-intl-compat';

import {cn} from '@/lib/utils';

interface HomeCategoryCardProps {
    category: contentservicev1_Category;
    onClick?: (id: number) => void;
    mobileCompact?: boolean;
}

const HomeCategoryCard: React.FC<HomeCategoryCardProps> = ({category, onClick, mobileCompact = false}) => {
    const t = useTranslations('page.home');

    const handleClick = () => {
        if (!category?.id) return;
        onClick?.(category.id);
    };

    const getIconName = (icon?: string): string => {
        if (!icon) return 'carbon:folder';
        return icon.includes(':') ? icon : `carbon:${icon}`;
    };

    return (
        <View
          className={cn(
                'group scroll-reveal-item relative cursor-pointer flex-col rounded-2xl',
                /* 移动端紧凑模式：固定宽度 + 垂直居中布局 */
                mobileCompact && 'w-[170px] shrink-0 flex items-center text-center p-4',
                /* PC端：保持原有满宽布局 */
                !mobileCompact && 'min-h-50 h-full flex p-5',
                /* 亮色：纯白卡片 + 柔和边框；暗色：深空黑 + slate边框 */
                'border border-slate-100 bg-white',
                'dark:border-slate-800/40 dark:bg-slate-900/50',
                /* 风行速度感：温和浮起 + 品牌绿辉光 */
                'transition-all duration-300 ease-out',
                'hover:-translate-y-1 hover:border-primary/30 hover:bg-slate-50',
                'dark:hover:bg-slate-900',
                'hover:shadow-[0_12px_30px_rgba(34,197,94,0.06)]',
            )}
          onClick={handleClick}
          style={{willChange: 'transform, box-shadow'}}
        >
            {/* hover 风迹渐变蒙层 */}
            <View className='pointer-events-none absolute inset-0 bg-linear-to-tr from-primary/5 via-transparent to-sky-400/5 opacity-0 transition-opacity duration-500 group-hover:opacity-100 rounded-2xl' />

            {mobileCompact ? (
                /* 移动端紧凑模式：图标在上，文字在下，垂直居中 */
                <View className='relative z-1 flex h-full w-full flex-col items-center justify-between gap-3'>
                    <View className={cn(
                        'flex h-14 w-14 shrink-0 items-center justify-center',
                        'rounded-xl bg-emerald-50 text-emerald-600 border border-emerald-100/50',
                        'dark:bg-primary/10 dark:text-primary dark:border-primary/20',
                        'transition-all duration-500 ease-out',
                        'group-hover:scale-110 group-hover:bg-emerald-500 group-hover:text-white group-hover:border-emerald-500',
                        'dark:group-hover:bg-primary dark:group-hover:text-white dark:group-hover:border-primary',
                        'group-hover:shadow-[0_8px_24px_-4px_hsl(var(--primary)/0.4)]',
                    )}
                    >
                        <XIcon name={getIconName(category.icon)} size={32} />
                    </View>
                    <View className='flex w-full min-w-0 flex-1 flex-col items-center gap-1'>
                        <Text className='line-clamp-2 text-sm font-bold leading-tight tracking-wide text-slate-900 dark:text-foreground transition-colors'>
                            {getCategoryName(category)}
                        </Text>
                        <Text className='text-xs text-slate-500 dark:text-muted-foreground'>
                            <Text className='mr-1 font-medium text-emerald-600 dark:text-primary'>{category.postCount || 0}</Text>
                            {t('articles_unit')}
                        </Text>
                    </View>
                    <View className={cn(
                        'mt-auto inline-flex w-fit items-center gap-1.5 rounded-md px-2.5 py-1',
                        'bg-slate-50 border border-slate-100',
                        'dark:bg-slate-900/40 dark:border-slate-800/40',
                        'text-xs text-slate-500 dark:text-muted-foreground transition-colors duration-300',
                        'group-hover:text-slate-700 group-hover:border-slate-200',
                        'dark:group-hover:text-foreground/70 dark:group-hover:border-slate-700/50',
                    )}
                    >
                        <XIcon name='carbon:time' size={14} />
                        <Text>{t('updated_days_ago', {days: 3})}</Text>
                    </View>
                </View>
            ) : (
                /* PC端：图标在左，文字在右 */
                <View className='relative z-1 flex h-full flex-col justify-between'>
                    <View className='mb-5 flex gap-5'>
                        <View className={cn(
                            'mb-4 flex h-17.5 w-17.5 shrink-0 items-center justify-center',
                            'rounded-xl bg-emerald-50 text-emerald-600 border border-emerald-100/50 text-3xl',
                            'dark:bg-primary/10 dark:text-primary dark:border-primary/20',
                            'transition-all duration-500 ease-out',
                            'group-hover:scale-110 group-hover:bg-emerald-500 group-hover:text-white group-hover:border-emerald-500',
                            'dark:group-hover:bg-primary dark:group-hover:text-white dark:group-hover:border-primary',
                            'group-hover:shadow-[0_8px_24px_-4px_hsl(var(--primary)/0.4)]',
                        )}
                        >
                            <XIcon name={getIconName(category.icon)} size={48} />
                        </View>
                        <View className='flex-1 w-full min-w-0'>
                            <Text className='mb-2 text-lg font-extrabold leading-tight tracking-wide text-slate-900 dark:text-foreground transition-colors'>
                                {getCategoryName(category)}
                            </Text>
                            <Text className='text-xs text-slate-500 dark:text-muted-foreground mt-1'>
                                <Text className='mr-1 font-medium text-emerald-600 dark:text-primary'>{category.postCount || 0}</Text>
                                {t('articles_unit')}
                            </Text>
                        </View>
                    </View>
                    <View className={cn(
                        'mt-auto w-fit inline-flex items-center gap-1.5 rounded-md px-2.5 py-1',
                        'bg-slate-50 border border-slate-100',
                        'dark:bg-slate-900/40 dark:border-slate-800/40',
                        'text-xs text-slate-500 dark:text-muted-foreground transition-colors duration-300',
                        'group-hover:text-slate-700 group-hover:border-slate-200',
                        'dark:group-hover:text-foreground/70 dark:group-hover:border-slate-700/50',
                    )}
                    >
                        <XIcon name='carbon:time' size={14} />
                        <Text>{t('updated_days_ago', {days: 3})}</Text>
                    </View>
                </View>
            )}
        </View>
    );
};

export default HomeCategoryCard;
