'use client';

import {useState, useEffect, useMemo} from 'react';
import {useParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {Button} from '@/components/ui/button';
import {ArrowLeft} from 'lucide-react';
import {AppEmpty} from '@/components/ui';

import {useI18nRouter} from "@/i18n/helpers";

import {
    fetchCategory,
    getCategoryName,
    getCategoryDescription,
} from '@/api/hooks/category';
import CategoryList from '@/components/category/CategoryList';
import PostListWithPagination from '@/components/post/PostList';
import {contentservicev1_Category} from "@/api/generated/app/service/v1";

import '../../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用

export default function CategoryDetailPage() {
    const t = useTranslations('page');
    const params = useParams();
    const router = useI18nRouter();

    const [_loading, setLoading] = useState(false);
    const [category, setCategory] = useState<contentservicev1_Category | null>(null);
    const [childCategories, setChildCategories] = useState<contentservicev1_Category[]>([]);

    const categoryId = useMemo(() => {
        const id = params?.id;
        return id ? parseInt(id as string) : null;
    }, [params?.id]);

    // Get parent category ID
    const parentCategoryId = useMemo(() => {
        if (!category?.parentId) return null;
        return category.parentId;
    }, [category?.parentId]);

    async function loadCategory() {
        if (!categoryId) return;

        setLoading(true);
        try {
            const categoryData = (await fetchCategory({
                id: categoryId!,
                fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image'
            })) as contentservicev1_Category;
            setCategory(categoryData);
            // Extract child categories - create new array reference to ensure reactive update
            if (categoryData?.children && categoryData.children.length > 0) {
                setChildCategories([...categoryData.children]);
            } else {
                setChildCategories([]);
            }
        } catch (error) {
            console.error('Load category failed:', error);
        } finally {
            setLoading(false);
        }
    }

    function handleViewChildCategory(id: number) {
        router.push(`/category/${id}`);
    }

    function handleBackToParent() {
        if (parentCategoryId) {
            router.push(`/category/${parentCategoryId}`);
        } else {
            router.push('/category');
        }
    }

    useEffect(() => {
        loadCategory();
    }, [categoryId]);

    if (!categoryId) {
        return <AppEmpty variant="error" description="Invalid category ID"/>;
    }

    return (
        <div className="w-full">
            {/* Hero Section */}
            <section className="w-full flex min-h-[300px] items-center justify-center overflow-hidden border-b border-border bg-gradient-to-br from-primary/10 via-background to-background py-20">
                <div className="w-full max-w-3xl px-8 text-center">
                    <h1 className="mb-4 text-4xl font-bold text-foreground max-md:text-3xl">
                        {getCategoryName(category)}
                    </h1>
                    {getCategoryDescription(category) && (
                        <p className="mb-6 text-lg text-muted-foreground max-md:text-base">
                            {getCategoryDescription(category)}
                        </p>
                    )}
                    <div className="flex items-center justify-center gap-6 text-sm text-muted-foreground">
                        <div className="flex items-center gap-2">
                            <span className="iconfont icon-document"/>
                            <span>{category?.postCount || 0} {t('posts.articles')}</span>
                        </div>
                    </div>
                </div>
            </section>

            {/* Posts Section */}
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                {/* Back to Parent Button */}
                <div className="mb-8">
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={handleBackToParent}
                        className="gap-2"
                    >
                        <ArrowLeft className="h-4 w-4"/>
                        {parentCategoryId ? t('categories.back_to_parent') : t('categories.back_to_list')}
                    </Button>
                </div>

                {/* Sub Categories List */}
                {childCategories.length > 0 && (
                    <div className="mb-8">
                        <CategoryList
                            categories={childCategories}
                            loading={false}
                            showSkeleton={false}
                            onCategoryClick={handleViewChildCategory}
                        />
                    </div>
                )}

                {/* Posts List with Pagination */}
                <PostListWithPagination
                    key={categoryId}
                    initialPageSize={10}
                    pageSizes={[10, 20, 30, 40]}
                    categoryId={categoryId}
                    from="category"
                />
            </div>
        </div>
    );
}
