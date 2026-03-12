import React, {useState, useEffect, useCallback, useRef, useMemo} from 'react';
import {Skeleton, Carousel, Button} from 'antd';
import {useTranslations} from 'next-intl';

import {useCategoryStore} from '@/store/slices/category/hooks';
import {
    contentservicev1_Category,
    contentservicev1_ListCategoryResponse,
} from '@/api/generated/app/service/v1';
import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';


import HomeCategoryCard from './HomeCategoryCard';

import styles from './home.module.css';

interface CategoryListSectionProps {
    skeletonCount?: number;
    showCarousel?: boolean;
    pageSize?: number;
    page?: number;
    filter?: Record<string, unknown>;
    orderBy?: string[];
    fieldMask?: string;
    showHeader?: boolean;
}

export default function CategoryListSection({
                                                skeletonCount = 8,
                                                showCarousel = false,
                                                pageSize = 8,
                                                page = 1,
                                                filter = {status: 'CATEGORY_STATUS_ACTIVE'} as Record<string, unknown>,
                                                orderBy = ['-sortOrder', '-postCount'],
                                                fieldMask = 'id,status,sortOrder,icon,code,postCount,directPostCount,parent_id,createdAt,translations.id,translations.categoryId,translations.name,translations.languageCode,translations.description',
                                                showHeader = true
                                            }: CategoryListSectionProps) {
    const t = useTranslations('page.home');
    const categoryStore = useCategoryStore();
    const router = useI18nRouter();

    // 使用 useMemo 稳定对象引用
    const stableFilter = useMemo(() => filter, [filter]);
    const stableOrderBy = useMemo(() => orderBy, [orderBy]);

    const [categories, setCategories] = useState<contentservicev1_Category[]>([]);
    const [loading, setLoading] = useState(false);

    // 用于取消异步操作
    const abortControllerRef = useRef<AbortController | null>(null);

    const loadCategories = useCallback(async () => {
        // 取消之前的请求
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
        }

        // 创建新的 AbortController
        abortControllerRef.current = new AbortController();
        const signal = abortControllerRef.current.signal;

        setLoading(true);
        try {
            const res = await categoryStore.listCategory({
                // @ts-expect-error - 参数类型推断问题
                paging: {page, pageSize},
                formValues: stableFilter,
                fieldMask,
                orderBy: stableOrderBy,
                signal,
            }) as unknown as contentservicev1_ListCategoryResponse;

            if (signal.aborted) return;

            setCategories(res.items || []);
        } catch (error) {
            if (signal.aborted) return;
            console.error('Failed to load categories:', error);
            setCategories([]);
        } finally {
            if (!signal.aborted) {
                setLoading(false);
            }
        }
    }, [categoryStore, stableFilter, fieldMask, stableOrderBy]); // 移除 page, pageSize 依赖

    useEffect(() => {
        loadCategories();

        // 组件卸载时取消请求
        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, []); // 空依赖数组，只在首次渲染时执行

    const handleViewCategoryDetail = (id: number) => {
        router.push(`/category/${id}`);
    };

    const handleViewCategory = () => {
        router.push(`/category`);
    };

    const getIconName = (icon?: string): string => {
        if (!icon) return 'carbon:folder';
        return icon.includes(':') ? icon : `carbon:${icon}`;
    };

    return (
        <section className={`${styles.categoriesSection} scroll-reveal`}>
            {/* Section Header */}
            {showHeader && (
                <div className={styles.sectionHeader}>
                    <h2 className={styles.sectionTitle}>
                        <XIcon name="carbon:folder-details" size={28} style={{color: '#6366f1', marginRight: '8px'}}/>
                        {t('categories')}
                    </h2>
                    <Button type="text" onClick={handleViewCategory}>
                        {t('view_all')} →
                    </Button>
                </div>
            )}

            {/* Loading Skeleton */}
            {loading ? (
                <div className={`${styles.categoriesGrid} desktop-grid`}>
                    {Array.from({length: skeletonCount}).map((_, i) => (
                        <div key={i} className={styles.categoryCardSkeleton}>
                            <Skeleton.Button style={{width: '100%', height: '140px'}}/>
                        </div>
                    ))}
                </div>
            ) : (
                <>
                    {/* Desktop Grid - 始终渲染，由 CSS 控制显示 */}
                    <div className={`${styles.categoriesGrid} desktop-grid`}>
                        {categories.map((category) => (
                            <HomeCategoryCard
                                key={category.id}
                                category={category}
                                onClick={handleViewCategoryDetail}
                            />
                        ))}
                    </div>

                    {/* Mobile Carousel - 始终渲染，由 CSS 控制显示 */}
                    <div className={`${styles.categoriesCarousel} mobile-carousel`}>
                        <Carousel
                            autoplay
                            autoplaySpeed={5000}
                            arrows={false}
                            draggable
                            slidesToShow={1.5}
                            centerMode
                            dots={false}
                            className={styles.carouselContainer}
                        >
                            {categories.map((category) => (
                                <div
                                    key={category.id}
                                    className={styles.carouselItem}
                                    onClick={() => handleViewCategoryDetail(category.id || 0)}
                                >
                                    <HomeCategoryCard
                                        category={category}
                                        onClick={handleViewCategoryDetail}
                                    />
                                </div>
                            ))}
                        </Carousel>
                    </div>
                </>
            )}
        </section>
    );
}
