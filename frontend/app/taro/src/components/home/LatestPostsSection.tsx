import {View, Text} from '@tarojs/components';
﻿import React from 'react';
import {Button} from '@/components/ui/button';
import {useTranslations} from '@/lib/next-intl-compat';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import PostList from '@/components/post/PostList';

export default function LatestPostsSection() {
    const t = useTranslations('page.home');
    const router = useI18nRouter();
    return (
        <View className='w-full scroll-reveal'>
            <View className='mb-8 flex items-center justify-between'>
                <Text className='flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl'>
                    <XIcon name='carbon:document' size={28} className='mr-2 text-primary' />
                    {t('latest_posts')}
                </Text>
                <Button variant='ghost' onClick={() => router.push('/post')}>
                    {t('view_all')} →
                </Button>
            </View>
            <View className='w-full'>
                <PostList
                  queryParams={{status: 'POST_STATUS_PUBLISHED'}}
                  fieldMask='id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail'
                  orderBy={['-createdAt']}
                  page={1}
                  pageSize={6}
                  columns={3}
                  showSkeleton
                  from='home'
                  showPagination={false}
                />
            </View>
        </View>
    );
}
