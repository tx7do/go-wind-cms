import React, {useState, useEffect, useCallback, useRef, useMemo} from 'react';
import {Skeleton} from '@/components/ui/skeleton';
import {Button} from '@/components/ui/button';
import {
    Carousel,
    CarouselContent,
    CarouselItem,
} from '@/components/ui/carousel';
import {useTranslations} from 'next-intl';

import {fetchListCategories} from '@/api/hooks/category';
import {
    contentservicev1_Category,
    contentservicev1_ListCategoryResponse,
} from '@/api/generated/app/service/v1';
import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';


import HomeCategoryCard from './HomeCategoryCard';

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
            const res = await fetchListCategories({
                paging: {page, pageSize},
                formValues: stableFilter,
                fieldMask,
                orderBy: stableOrderBy,
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
    }, [stableFilter, fieldMask, stableOrderBy]); // 移除 page, pageSize 依赖

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

    return (
        <section className="categories-section w-full max-w-300 mx-auto scroll-reveal px-8 py-12 max-md:px-4">
            {/* Section Header */}
            {showHeader && (
                <div className="mb-8 flex items-center justify-between">
                    <h2 className="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
                        <XIcon name="carbon:folder-details" size={28} className="mr-2 text-primary"/>
                        {t('categories')}
                    </h2>
                    <Button variant="ghost" onClick={handleViewCategory}>
                        {t('view_all')} →
                    </Button>
                </div>
            )}

            {/* Loading Skeleton */}
            {loading ? (
                <div className="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6">
                    {Array.from({length: skeletonCount}).map((_, i) => (
                        <div key={i} className="flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6">
                            <Skeleton className="h-35 w-full"/>
                        </div>
                    ))}
                </div>
            ) : (
                <>
                    {/* Desktop Grid */}
                    <div className="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6 max-md:hidden">
                        {categories.map((category) => (
                            <HomeCategoryCard
                                key={category.id}
                                category={category}
                                onClick={handleViewCategoryDetail}
                            />
                        ))}
                    </div>

                    {/* Mobile Carousel */}
                    <div className="hidden max-md:block w-full px-2 pb-8">
                        <Carousel
                            opts={{
                                align: 'center',
                                loop: false,
                            }}
                            className="w-full py-6"
                        >
                            <CarouselContent>
                                {categories.map((category) => (
                                    <CarouselItem
                                        key={category.id}
                                        className="flex items-stretch justify-center px-2 py-3"
                                        onClick={() => handleViewCategoryDetail(category.id || 0)}
                                    >
                                        <HomeCategoryCard
                                            category={category}
                                            onClick={handleViewCategoryDetail}
                                        />
                                    </CarouselItem>
                                ))}
                            </CarouselContent>
                        </Carousel>
                    </div>
                </>
            )}
        </section>
    );
}
