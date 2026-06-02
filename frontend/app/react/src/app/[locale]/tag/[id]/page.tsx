'use client';

import {useState, useEffect, useMemo} from 'react';
import {useParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {Button} from '@/components/ui/button';
import {ArrowLeft} from 'lucide-react';
import {AppEmpty} from '@/components/ui';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from "@/i18n/helpers";

import {fetchTag, getTranslation as getTagTranslation} from '@/api/hooks/tag';
import PostListWithPagination from '@/components/post/PostList';
import type {contentservicev1_Tag} from "@/api/generated/app/service/v1";

import '../../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用

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
            {/* Hero Section */}
            <section
                className="w-full flex min-h-[300px] items-center justify-center overflow-hidden border-b border-border py-20"
                style={{
                    background: tag?.color
                        ? `linear-gradient(135deg, ${tag.color} 0%, ${tag.color}dd 50%, ${tag.color}aa 100%)`
                        : 'linear-gradient(135deg, hsl(var(--primary) / 0.1) 0%, hsl(var(--background)) 100%)'
                }}
            >
                <div className="w-full max-w-3xl px-8 text-center">
                    <div
                        className="mb-4 flex items-center justify-center"
                        style={{color: tag?.color || 'hsl(var(--primary))'}}
                    >
                        <XIcon name={`carbon:${tag?.icon || 'tag'}`} size={64}/>
                    </div>
                    <h1 className="mb-4 text-4xl font-bold text-foreground max-md:text-3xl">
                        {getTagTranslation(tag)?.name || t('tags.tag_untitled')}
                    </h1>
                    {getTagTranslation(tag)?.description && (
                        <p className="mb-6 text-lg text-muted-foreground max-md:text-base">
                            {getTagTranslation(tag)?.description}
                        </p>
                    )}
                    <div className="flex items-center justify-center gap-6 text-sm text-muted-foreground">
                        <div className="flex items-center gap-2">
                            <XIcon name="carbon:document" size={20}/>
                            <span>{tag?.postCount || 0} {t('posts.articles')}</span>
                        </div>
                    </div>
                </div>
            </section>

            {/* Posts Section */}
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                {/* Back Button */}
                <div className="mb-8">
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={handleBack}
                        className="gap-2"
                    >
                        <ArrowLeft className="h-4 w-4"/>
                        {t('categories.back_to_list')}
                    </Button>
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
