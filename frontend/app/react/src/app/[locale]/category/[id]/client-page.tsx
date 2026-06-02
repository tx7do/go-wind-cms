'use client';

import {useState, useEffect, useMemo} from 'react';
import {useParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {AppEmpty} from '@/components/ui';
import BackButton from '@/components/layout/BackButton';

import {useI18nRouter} from "@/i18n/helpers";

import {
    fetchCategory,
    getCategoryName,
    getCategoryDescription,
} from '@/api/hooks/category';
import CategoryList from '@/components/category/CategoryList';
import PostListWithPagination from '@/components/post/PostList';
import PageHero from '@/components/layout/PageHero';
import SectionContainer from '@/components/layout/SectionContainer';
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
        // 有父分类 → 回父分类；无父分类 → 回分类总览页
        // 分类详情页没有 from 参数，router.back() 退回的页面不可预测，所以直接导航到确定的路由
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
            <SectionContainer>
                {/* Back Button */}
                <div className="mb-8">
                    <BackButton 
                        label={parentCategoryId ? t('categories.back_to_parent') : t('categories.back_to_list')} 
                        onClick={handleBackToParent}
                    />
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
            </SectionContainer>
        </div>
    );
}
