'use client';

import {useState, useEffect, useMemo} from 'react';
import {useParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {Button} from 'antd';
import {ArrowLeftOutlined} from '@ant-design/icons';
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
import styles from './category-detail.module.css';

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
            // If no parent category, go back to category list
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
        <div className={styles['category-detail-page']}>
            {/* Hero Section */}
            <div className={styles['hero-section']}>
                <div className={styles['hero-content']}>
                    <h1>{getCategoryName(category)}</h1>
                    {getCategoryDescription(category) && (
                        <p className={styles['category-description']}>
                            {getCategoryDescription(category)}
                        </p>
                    )}
                    <div className={styles['category-stats']}>
                        <div className={styles['stat-item']}>
                            <span className={`${styles['icon']} iconfont icon-document`}/>
                            <span>{category?.postCount || 0} {t('posts.articles')}</span>
                        </div>
                    </div>
                </div>
            </div>

            {/* Posts Section */}
            <div className={styles['page-container']}>
                {/* Back to Parent Button */}
                <div className={styles['back-button-container']}>
                    <Button
                        size="small"
                        onClick={handleBackToParent}
                        icon={<ArrowLeftOutlined/>}
                    >
                        {parentCategoryId ? t('categories.back_to_parent') : t('categories.back_to_list')}
                    </Button>
                </div>

                {/* Sub Categories List */}
                {childCategories.length > 0 && (
                    <div className={styles['sub-categories-wrapper']}>
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
