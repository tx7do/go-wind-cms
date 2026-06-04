import {View, Text} from '@tarojs/components';
﻿import React, {useRef} from 'react';
import {Button} from '@/components/ui/button';
import {useTranslations} from '@/lib/next-intl-compat';

import {XIcon} from '@/plugins/xicon';
import PostList from '@/components/post/PostList';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export default function FeaturedPostsSection() {
    const t = useTranslations('page.home');
    const router = useI18nRouter();

    const scrollToCategories = () => {
        // Taro 兼容：H5 用锚点定位，小程序跳过
        if (typeof window !== 'undefined') {
            const section = document.querySelector('.categories-section');
            if (section) {
                section.scrollIntoView({behavior: 'smooth', block: 'start'});
            }
        }
    };

    return (
        <View className='w-full scroll-reveal'>
            <View className='mb-8 flex items-center justify-between'>
                <Text className='flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl'>
                    <XIcon name='carbon:star-filled' size={28} className='mr-2 text-primary' />
                    {t('featured_posts')}
                </Text>
                <Button variant='ghost' onClick={() => router.push('/post')}>
                    {t('view_all')} →
                </Button>
            </View>
            <View className='w-full'>
                <PostList
                  queryParams={{status: 'POST_STATUS_PUBLISHED', isFeatured: true}}
                  fieldMask='id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail'
                  orderBy={['-sortOrder']}
                  page={1}
                  pageSize={3}
                  columns={3}
                  showSkeleton
                  from='home'
                  showPagination={false}
                />
            </View>
            <View className='mt-8 flex w-full items-center justify-center max-md:hidden'>
                <View className='flex items-center gap-4'>
                    <Text className='text-2xl font-extrabold leading-tight tracking-tight text-primary max-md:text-xl'>
                        {t('explore_more_categories')}
                    </Text>
                    <View
                      className='group flex h-9 w-9 cursor-pointer items-center justify-center rounded-full border border-primary/20 bg-primary/10 transition-all duration-300 hover:border-primary hover:bg-primary hover:shadow-[0_0_20px_hsl(var(--primary)/0.3)]'
                      onClick={scrollToCategories}
                    >
                        <XIcon name='carbon:arrow-down' size={24} className='text-primary transition-colors group-hover:text-white' />
                    </View>
                </View>
            </View>
        </View>
    );
}
