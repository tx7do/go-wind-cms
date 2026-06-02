'use client';

import {useState, useMemo} from 'react';
import {useTranslations} from 'next-intl';

import CategoryFilter from '@/components/category/CategoryFilter';
import PostList from '@/components/post/PostList';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用

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
            {/* Hero Section */}
            <section className="w-full flex min-h-[300px] items-center justify-center border-b border-border bg-gradient-to-br from-primary/10 via-background to-background py-20">
                <div className="w-full max-w-3xl px-8 text-center">
                    <h1 className="mb-4 text-4xl font-bold text-foreground max-md:text-3xl">
                        {t('posts.posts_list')}
                    </h1>
                    <p className="text-lg text-muted-foreground max-md:text-base">
                        {t('posts.explore_latest')}
                    </p>
                </div>
            </section>

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
