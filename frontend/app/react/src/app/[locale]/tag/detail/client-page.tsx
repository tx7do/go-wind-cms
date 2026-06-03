'use client';

import {useState, useEffect, useMemo} from 'react';
import {useSearchParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {AppEmpty} from '@/components/ui';
import BackButton from '@/components/layout/BackButton';

import {useI18nRouter} from "@/i18n/helpers";
import PageHero from '@/components/layout/PageHero';
import SectionContainer from '@/components/layout/SectionContainer';
import {XIcon} from '@/plugins/xicon';

import {fetchTag, getTranslation as getTagTranslation} from '@/api/hooks/tag';
import PostListWithPagination from '@/components/post/PostList';
import type {contentservicev1_Tag} from "@/api/generated/app/service/v1";

export default function TagDetailPage() {
    const t = useTranslations('page');
    const router = useI18nRouter();
    const searchParams = useSearchParams();

    const [loading, setLoading] = useState(false);
    const [tag, setTag] = useState<contentservicev1_Tag | null>(null);

    const tagId = useMemo(() => {
        const id = searchParams.get('id');
        return id ? parseInt(id) : null;
    }, [searchParams]);

    async function loadTag() {
        if (!tagId) return;

        setLoading(true);
        try {
            const tagData = await fetchTag(tagId!) as unknown as contentservicev1_Tag;
            setTag(tagData);
        } catch (error) {
            console.error('Load tag failed:', error);
        } finally {
            setLoading(false);
        }
    }

    function handleBack() {
        router.push('/tag');
    }

    useEffect(() => {
        loadTag();
    }, [tagId]);

    if (!tagId) {
        return <AppEmpty variant="error" description="Invalid tag ID"/>;
    }

    return (
        <div className="w-full">
            <PageHero
                title={getTagTranslation(tag)?.name || t('tags.tag_untitled')}
                description={getTagTranslation(tag)?.description || undefined}
                icon={`carbon:${tag?.icon || 'tag'}`}
                iconSize={56}
                size="lg"
                accentColor={tag?.color || undefined}
                meta={
                    <div className="flex items-center gap-2">
                        <XIcon name="carbon:document" size={16}/>
                        <span>{tag?.postCount || 0} {t('posts.articles')}</span>
                    </div>
                }
            />

            {/* Posts Section */}
            <SectionContainer>
                {/* Back Button — ghost 风格，与 Post/Category 详情页统一 */}
                <div className="mb-8">
                    <BackButton 
                        label={t('categories.back_to_list')} 
                        onClick={handleBack}
                    />
                </div>

                {/* Posts List with Pagination */}
                {tagId && (
                    <PostListWithPagination
                        key={tagId}
                        initialPageSize={10}
                        pageSizes={[10, 20, 30, 40]}
                        tagId={tagId}
                        from="tag"
                    />
                )}
            </SectionContainer>
        </div>
    );
}
