import React from 'react';
import {Button} from '@/components/ui/button';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import PostList from '@/components/post/PostList';

export default function LatestPostsSection() {
    const t = useTranslations('page.home');
    const router = useI18nRouter();
    return (
        <section className="w-full max-w-[1200px] mx-auto scroll-reveal px-8 py-12 max-md:px-4">
            <div className="mb-8 flex items-center justify-between">
                <h2 className="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
                    <XIcon name="carbon:document" size={28} className="mr-2 text-primary"/>
                    {t('latest_posts')}
                </h2>
                <Button variant="ghost" onClick={() => router.push('/post')}>
                    {t('view_all')} →
                </Button>
            </div>
            <div className="w-full">
                <PostList
                    queryParams={{status: 'POST_STATUS_PUBLISHED'}}
                    fieldMask="id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail"
                    orderBy={['-createdAt']}
                    page={1}
                    pageSize={6}
                    columns={3}
                    showSkeleton={true}
                    from="home"
                    showPagination={false}
                />
            </div>
        </section>
    );
}
