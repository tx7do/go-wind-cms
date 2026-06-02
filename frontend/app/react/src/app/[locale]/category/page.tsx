'use client';

import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import {Skeleton} from '@/components/ui/skeleton';
import {AppEmpty} from '@/components/ui';

import {fetchListCategories} from '@/api/hooks/category';
import CategoryTree from '@/components/category/CategoryTree';

import {contentservicev1_Category, contentservicev1_ListCategoryResponse} from "@/api/generated/app/service/v1";
import {useI18nRouter} from "@/i18n/helpers";

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用

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
        <div className="w-full">
            {/* Hero Section */}
            <section className="w-full flex min-h-[300px] items-center justify-center border-b border-border bg-gradient-to-br from-primary/10 via-background to-background py-20">
                <div className="w-full max-w-3xl px-8 text-center">
                    <h1 className="mb-4 text-4xl font-bold text-foreground max-md:text-3xl">
                        {t('categories.categories')}
                    </h1>
                    <p className="text-lg text-muted-foreground max-md:text-base">
                        {t('categories.browse_all')}
                    </p>
                </div>
            </section>

            {/* Content Section */}
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                {/* Loading Skeleton */}
                {loading ? (
                    <div className="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] gap-6">
                        {[...Array(6)].map((_, i) => (
                            <div key={i}>
                                <Skeleton className="h-40 w-full"/>
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
                                image={<span className="i-carbon:folder-blank text-[64px]"/>}
                            />
                        )}
                    </>
                )}
            </div>
        </div>
    );
}
