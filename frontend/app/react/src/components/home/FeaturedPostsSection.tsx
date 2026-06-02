import React, {useRef} from 'react';
import {Button} from '@/components/ui/button';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import PostList from '@/components/post/PostList';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export default function FeaturedPostsSection() {
    const t = useTranslations('page.home');
    const router = useI18nRouter();

    const scrollToCategories = () => {
        const section = document.querySelector('.categories-section');
        if (section) {
            section.scrollIntoView({behavior: 'smooth', block: 'start'});
        }
    };

    return (
        <section className="w-full max-w-[1200px] mx-auto scroll-reveal px-8 py-12 max-md:px-4">
            <div className="mb-8 flex items-center justify-between">
                <h2 className="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
                    <XIcon name="carbon:star-filled" size={28} className="mr-2 text-primary"/>
                    {t('featured_posts')}
                </h2>
                <Button variant="ghost" onClick={() => router.push('/post')}>
                    {t('view_all')} →
                </Button>
            </div>
            <div className="w-full">
                <PostList
                    queryParams={{status: 'POST_STATUS_PUBLISHED', isFeatured: true}}
                    fieldMask="id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail"
                    orderBy={['-sortOrder']}
                    page={1}
                    pageSize={3}
                    columns={3}
                    showSkeleton={true}
                    from="home"
                    showPagination={false}
                />
            </div>
            <div className="mt-8 flex w-full items-center justify-center max-md:hidden">
                <div className="flex items-center gap-4">
                    <span className="text-2xl font-extrabold leading-tight tracking-tight text-primary max-md:text-xl">
                        {t('explore_more_categories')}
                    </span>
                    <div
                        className="group flex h-9 w-9 cursor-pointer items-center justify-center rounded-full border border-primary/20 bg-primary/10 transition-all duration-300 hover:border-primary hover:bg-primary hover:shadow-[0_0_20px_hsl(var(--primary)/0.3)]"
                        onClick={scrollToCategories}
                    >
                        <XIcon name="carbon:arrow-down" size={24} className="text-primary transition-colors group-hover:text-white"/>
                    </div>
                </div>
            </div>
        </section>
    );
}
