'use client';

import {useState, useEffect, useMemo} from 'react';
import {useParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
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
import PageHero from '@/components/layout/PageHero';
import {XIcon} from '@/plugins/xicon';
import {contentservicev1_Category} from "@/api/generated/app/service/v1";

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
        // 优先级：有父分类 → 回父分类；有历史 → router.back()；降级 → /category
        if (parentCategoryId) {
            router.push(`/category/${parentCategoryId}`);
        } else if (typeof window !== 'undefined' && window.history.length > 2) {
            router.back();
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
            <PageHero
                title={getCategoryName(category) || ''}
                description={getCategoryDescription(category) || undefined}
                icon={category?.icon || 'carbon:folder'}
                size="lg"
                meta={
                    <>
                        <div className="flex items-center gap-2">
                            <XIcon name="carbon:document" size={16}/>
                            <span>{category?.postCount || 0} {t('posts.articles')}</span>
                        </div>
                        {childCategories.length > 0 && (
                            <div className="flex items-center gap-2">
                                <XIcon name="carbon:folder" size={16}/>
                                <span>{childCategories.length} {t('categories.categories')}</span>
                            </div>
                        )}
                    </>
                }
            />

            {/* Posts Section */}
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                {/* Back Button */}
                <div className="mb-8">
                    <button
                        onClick={handleBackToParent}
                        className="group flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground"
                    >
                        <ArrowLeft className="h-4 w-4 transition-transform duration-300 group-hover:-translate-x-1"/>
                        <span>{parentCategoryId ? t('categories.back_to_parent') : t('categories.back_to_list')}</span>
                    </button>
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
