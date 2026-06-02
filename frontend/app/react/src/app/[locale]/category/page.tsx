'use client';

import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import {Skeleton} from 'antd';
import {AppEmpty} from '@/components/ui';

import {fetchListCategories} from '@/api/hooks/category';
import CategoryTree from '@/components/category/CategoryTree';

import {contentservicev1_Category, contentservicev1_ListCategoryResponse} from "@/api/generated/app/service/v1";
import {useI18nRouter} from "@/i18n/helpers";

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './page.module.css';

export default function CategoryListPage() {
    const t = useTranslations('page');
    const router = useI18nRouter();

    const [loading, setLoading] = useState(false);
    const [categories, setCategories] = useState<contentservicev1_Category[]>([]);

    async function loadCategories() {
        setLoading(true);
        try {
            const res = (await fetchListCategories({
                paging: undefined,
                formValues: {status: 'CATEGORY_STATUS_ACTIVE'},
                fieldMask: 'id,status,sort_order,icon,code,post_count,direct_post_count,parent_id,created_at,children,translations.id,translations.category_id,translations.name,translations.language_code,translations.description,translations.thumbnail,translations.cover_image',
                orderBy: ['-sortOrder']
            })) as unknown as contentservicev1_ListCategoryResponse;
            setCategories(res.items || []);
        } catch (error) {
            console.error('Load categories failed:', error);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        loadCategories();
    }, []);

    const handleCategoryClick = (id: number) => {
        router.push(`/category/${id}`);
    };

    return (
        <div className={styles['category-page']}>
            {/* Hero Section */}
            <div className={styles['hero-section']}>
                <div className={styles['hero-content']}>
                    <h1>{t('categories.categories')}</h1>
                    <p>{t('categories.browse_all')}</p>
                </div>
            </div>

            {/* Content Section */}
            <div className={styles['page-container']}>
                {/* Loading Skeleton */}
                {loading ? (
                    <div className={styles['categories-loading']}>
                        {[...Array(6)].map((_, i) => (
                            <div key={i} className={styles['category-loading-card']}>
                                <Skeleton.Image style={{height: 160}} active/>
                            </div>
                        ))}
                    </div>
                ) : (
                    <>
                        {categories.length > 0 ? (
                            <CategoryTree
                                categories={categories}
                                onCategoryClick={handleCategoryClick}
                            />
                        ) : (
                            <AppEmpty
                                description={t('categories.no_categories')}
                                inContainer
                                image={<span className="i-carbon:folder-blank" style={{fontSize: '64px'}}/>}
                            />
                        )}
                    </>
                )}
            </div>
        </div>
    );
}
