import React, {useRef} from 'react';
import {Button} from '@/components/ui/button';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import PostList from '@/components/post/PostList';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import styles from './home.module.css';

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
        <section className={`${styles.featuredSection} scroll-reveal`}>
            <div className={styles.sectionHeader}>
                <h2 className={styles.sectionTitle}>
                    <XIcon name="carbon:star-filled" size={28} style={{color: '#6366f1', marginRight: '8px'}} />
                    {t('featured_posts')}
                </h2>
                <Button variant="ghost" onClick={() => router.push('/post')}>
                    {t('view_all')} →
                </Button>
            </div>
            <div className={styles.featuredGrid}>
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
            <div className={styles.scrollIndicator}>
                <div className={styles.scrollIndicatorContent}>
                    <span className={styles.scrollIndicatorText}>{t('explore_more_categories')}</span>
                    <div className={styles.scrollArrow} onClick={scrollToCategories}>
                        <XIcon name="carbon:arrow-down" size={24} className={styles.arrowIcon} />
                    </div>
                </div>
            </div>
        </section>
    );
}
