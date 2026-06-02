import React, {useState, useEffect, useCallback, useRef} from 'react';
import {Skeleton} from '@/components/ui/skeleton';
import {Button} from '@/components/ui/button';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {fetchListTags} from '@/api/hooks/tag';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {
    contentservicev1_ListTagResponse,
    contentservicev1_Tag
} from '@/api/generated/app/service/v1';

interface TagItem {
    id: number;
    name: string;
    color: string;
}

export default function PopularTagsSection() {
    const t = useTranslations('page.tags');
    const router = useI18nRouter();

    const [_tags, setTags] = useState<contentservicev1_Tag[]>([]);
    const [loading, setLoading] = useState(false);
    const [displayTags, setDisplayTags] = useState<TagItem[]>([]);

    // 用于取消异步操作
    const abortControllerRef = useRef<AbortController | null>(null);

    const loadPopularTags = useCallback(async () => {
        // 取消之前的请求
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
        }

        // 创建新的 AbortController
        abortControllerRef.current = new AbortController();
        const signal = abortControllerRef.current.signal;

        setLoading(true);
        try {
            const res = await fetchListTags({
                paging: {page: 1, pageSize: 6},
                formValues: {status: 'TAG_STATUS_ACTIVE', isFeatured: true},
                fieldMask: undefined,
                orderBy: undefined,
            }) as unknown as contentservicev1_ListTagResponse;

            if (signal.aborted) return;

            const tagItems = res.items || [];
            setTags(tagItems);

            // 生成带颜色的标签
            const taggedItems: TagItem[] = tagItems
                .filter(tag => tag.id !== undefined) // 过滤掉 id 为 undefined 的项
                .map((tag, index) => ({
                    id: tag.id!,
                    name: tag.translations?.[0]?.name || t('tag_untitled'),
                    color: tag.color || `hsl(${index * 60}, 100%, 50%)`,
                }));

            setDisplayTags(taggedItems);
        } catch (error) {
            if (signal.aborted) return;
            console.error('Failed to load tags:', error);
            setTags([]);
            setDisplayTags([]);
        } finally {
            if (!signal.aborted) {
                setLoading(false);
            }
        }
    }, [t]); // 添加 t 依赖保证闭包正确性，useEffect 空依赖不会导致无限循环

    useEffect(() => {
        loadPopularTags();

        // 组件卸载时取消请求
        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, []); // 空依赖数组，只在首次渲染时执行

    const handleViewTag = (tag: TagItem) => {
        router.push(`/tag/${tag.id}`);
    };

    return (
        <section className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
            <div className="mb-8 flex items-center justify-between">
                <h2 className="flex items-center gap-2 text-2xl font-extrabold tracking-tight text-foreground max-md:text-xl">
                    <XIcon name="carbon:fire" size={28} className="mr-2 text-primary"/>
                    {t('popular_tags')}
                </h2>
                <Button variant="ghost" onClick={() => router.push('/tag')}>
                    {t('view_all')} →
                </Button>
            </div>
            <div className="w-full">
                {loading ? (
                    <div className="flex flex-wrap justify-center gap-6">
                        {Array.from({length: 6}).map((_, i) => (
                            <div key={i} className="flex min-w-[90px] items-center rounded-2xl border-[1.5px] border-primary/20 bg-primary/5 px-6 py-4 shadow-sm">
                                <Skeleton className="h-11 w-full"/>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="flex flex-wrap justify-center gap-6 max-md:gap-2">
                        {displayTags.map((tag) => (
                            <div
                                key={tag.id}
                                className="group flex min-w-[90px] cursor-pointer items-center rounded-2xl border-[1.5px] border-primary/20 bg-primary/5 px-6 py-4 font-bold text-primary shadow-sm backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-primary hover:bg-primary hover:text-white hover:shadow-lg max-md:rounded-xl max-md:px-4 max-md:py-2.5 max-md:text-sm"
                                style={{'--tag-color': tag.color} as React.CSSProperties}
                                onClick={() => handleViewTag(tag)}
                            >
                                <span className="truncate text-base font-bold tracking-tight transition-colors group-hover:text-white max-md:text-sm">
                                    {tag.name}
                                </span>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </section>
    );
}
