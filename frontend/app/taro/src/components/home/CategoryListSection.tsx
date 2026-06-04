import {View, Text} from '@tarojs/components';
﻿import React, {useState, useEffect, useCallback, useRef, useMemo} from 'react';
import {Skeleton} from '@/components/ui/skeleton';
import {Button} from '@/components/ui/button';
import {useTranslations} from '@/lib/next-intl-compat';

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
        router.push(`/category/detail?id=${id}`);
    };

    const handleViewCategory = () => {
        router.push(`/category`);
    };

    return (
        <View className='categories-section w-full scroll-reveal'>
            {/* Section Header */}
            {showHeader && (
                <View className='mb-8 flex items-center justify-between'>
                    <Text className='flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl'>
                        <XIcon name='carbon:folder-details' size={28} className='mr-2 text-primary' />
                        {t('categories')}
                    </Text>
                    <Button variant='ghost' onClick={handleViewCategory}>
                        {t('view_all')} →
                    </Button>
                </View>
            )}

            {/* Loading Skeleton */}
            {loading ? (
                <View className='grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6'>
                    {Array.from({length: skeletonCount}).map((_, i) => (
                        <View key={i} className='flex min-h-50 flex-col rounded-2xl border border-border bg-card p-6'>
                            <Skeleton className='h-35 w-full' />
                        </View>
                    ))}
                </View>
            ) : (
                <>
                    {/* Desktop Grid */}
                    <View className='grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6 max-md:hidden'>
                        {categories.map((category) => (
                            <HomeCategoryCard
                              key={category.id}
                              category={category}
                              onClick={handleViewCategoryDetail}
                            />
                        ))}
                    </View>

                    {/* Mobile：横向传送带（X-Axis Slider） */}
                    <View className='hidden max-md:flex w-full gap-4 overflow-x-auto no-scrollbar scroll-smooth px-4 pb-4 pt-2 -mx-4'>
                        {categories.map((category) => (
                            <HomeCategoryCard
                              key={category.id}
                              category={category}
                              onClick={handleViewCategoryDetail}
                              mobileCompact
                            />
                        ))}
                    </View>
                </>
            )}
        </View>
    );
}
