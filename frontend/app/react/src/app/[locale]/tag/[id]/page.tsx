'use client';

import {useState, useEffect, useMemo} from 'react';
import {useParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {ArrowLeft} from 'lucide-react';
import {AppEmpty} from '@/components/ui';

import {useI18nRouter} from "@/i18n/helpers";
import PageHero from '@/components/layout/PageHero';
import {XIcon} from '@/plugins/xicon';

import {fetchTag, getTranslation as getTagTranslation} from '@/api/hooks/tag';
import PostListWithPagination from '@/components/post/PostList';
import type {contentservicev1_Tag} from "@/api/generated/app/service/v1";

export default function TagDetailPage() {
    const t = useTranslations('page');
    const params = useParams();
    const router = useI18nRouter();

    const [loading, setLoading] = useState(false);
    const [tag, setTag] = useState<contentservicev1_Tag | null>(null);

    const tagId = useMemo(() => {
        const id = params?.id;
        return id ? parseInt(id as string) : null;
    }, [params?.id]);

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
        // 标签详情页没有 from 参数，router.back() 退回的页面不可预测
        // 按钮文字是「返回列表」，始终导航到标签总览页
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
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                {/* Back Button — ghost 风格，与 Post/Category 详情页统一 */}
                <div className="mb-8">
                    <button
                        onClick={handleBack}
                        className="group flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground"
                    >
                        <ArrowLeft className="h-4 w-4 transition-transform duration-300 group-hover:-translate-x-1"/>
                        <span>{t('categories.back_to_list')}</span>
                    </button>
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
            </div>
        </div>
    );
}
