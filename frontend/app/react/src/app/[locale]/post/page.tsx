'use client';

import {useState, useMemo} from 'react';
import {useTranslations} from 'next-intl';

import PageHero from '@/components/layout/PageHero';
import CategoryFilter from '@/components/category/CategoryFilter';
import PostList from '@/components/post/PostList';

export default function PostListPage() {
    const t = useTranslations('page');
    const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(null);

    const handleCategoryChange = (categoryId: number | null) => {
        console.log('[PostListPage] Category changed:', categoryId);
        setSelectedCategoryId(categoryId);
    };

    // 构建查询参数 - 使用 useMemo 确保稳定性
    const queryParams = useMemo(() => {
        if (selectedCategoryId) {
            return {category_ids__in: [selectedCategoryId]};
        }
        return {};
    }, [selectedCategoryId]);

    return (
        <div className="w-full">
            <PageHero
                title={t('posts.posts_list')}
                description={t('posts.explore_latest')}
                icon="carbon:document"
                size="md"
            />

            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                <CategoryFilter
                    selectedCategory={selectedCategoryId}
                    treeMode={true}
                    autoLoad={true}
                    onCategoryChange={handleCategoryChange}
                />

                <PostList
                    key={selectedCategoryId || 'all'}  // 使用 key 强制重新渲染
                    queryParams={queryParams}
                    initialPageSize={12}
                    pageSizes={[12, 24, 36, 48]}
                />
            </div>
        </div>
    );
}
