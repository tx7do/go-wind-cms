'use client';

import {useState, useRef} from 'react';
import {useTranslations} from 'next-intl';

import CategoryFilter from '@/components/category/CategoryFilter';
import PostList from '@/components/post/PostList';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './post-list.module.css';

export default function PostListPage() {
    const t = useTranslations('page');
    const postListRef = useRef<HTMLDivElement>(null);
    const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(null);

    const handleCategoryChange = (categoryId: number | null) => {
        setSelectedCategoryId(categoryId);
    };

    return (
        <div className={styles['post-list-page']}>
            {/* Hero Section */}
            <div className={styles['hero-section']}>
                <div className={styles['hero-content']}>
                    <h1>{t('posts.posts_list')}</h1>
                    <p>{t('posts.explore_latest')}</p>
                </div>
            </div>

            <div className={styles['page-container']}>
                <CategoryFilter
                    selectedCategory={selectedCategoryId}
                    treeMode={true}
                    autoLoad={true}
                    onCategoryChange={handleCategoryChange}
                />

                <PostList
                    ref={postListRef}
                    initialPageSize={12}
                    pageSizes={[12, 24, 36, 48]}
                    categoryId={selectedCategoryId ?? undefined}
                />
            </div>
        </div>
    );
}
