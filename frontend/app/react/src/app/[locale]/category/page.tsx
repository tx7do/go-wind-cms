'use client';

import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import {Skeleton} from '@/components/ui/skeleton';
import {AppEmpty} from '@/components/ui';

import {fetchListCategories} from '@/api/hooks/category';
import CategoryTree from '@/components/category/CategoryTree';
import PageHero from '@/components/layout/PageHero';
import SectionContainer from '@/components/layout/SectionContainer';

import {contentservicev1_Category, contentservicev1_ListCategoryResponse} from "@/api/generated/app/service/v1";
import {useI18nRouter} from "@/i18n/helpers";

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
        router.push(`/category/detail?id=${id}`);
    };

    return (
        <div className="w-full">
            <PageHero
                title={t('categories.categories')}
                description={t('categories.browse_all')}
                icon="carbon:folder"
                size="md"
            />

            {/* Content Section */}
            <SectionContainer>
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
            </SectionContainer>
        </div>
    );
}
