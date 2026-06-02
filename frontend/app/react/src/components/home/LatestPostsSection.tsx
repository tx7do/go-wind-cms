import React from 'react';
import {Button} from '@/components/ui/button';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import PostList from '@/components/post/PostList';

import styles from './home.module.css';

export default function LatestPostsSection() {
    const t = useTranslations('page.home');
    const router = useI18nRouter();
    return (
        <section className={`${styles.latestSection} scroll-reveal`}>
            <div className={styles.sectionHeader}>
                <h2 className={styles.sectionTitle}>
                    <XIcon name="carbon:document" size={28} style={{color: '#6366f1', marginRight: '8px'}}/>
                    {t('latest_posts')}
                </h2>
                <Button variant="ghost" onClick={() => router.push('/post')}>
                    {t('view_all')} →
                </Button>
            </div>
            <div className={styles.latestGrid}>
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
