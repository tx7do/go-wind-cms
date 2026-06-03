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
    postCount: number;
}

export default function PopularTagsSection() {
    const t = useTranslations('page.tags');
    const router = useI18nRouter();

    const [_tags, setTags] = useState<contentservicev1_Tag[]>([]);
    const [loading, setLoading] = useState(false);
    const [displayTags, setDisplayTags] = useState<TagItem[]>([]);
    const [hoveredId, setHoveredId] = useState<number | null>(null);

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
                .filter(tag => tag.id !== undefined)
                .map((tag, index) => ({
                    id: tag.id!,
                    name: tag.translations?.[0]?.name || t('tag_untitled'),
                    color: tag.color || `hsl(${(index * 67) % 360}, 70%, 55%)`,
                    postCount: tag.postCount || 0,
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
    }, [t]);

    useEffect(() => {
        loadPopularTags();

        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, []);

    const handleViewTag = (tag: TagItem) => {
        router.push(`/tag/detail?id=${tag.id}`);
    };

    return (
        <section className="w-full max-w-300 mx-auto px-8 py-12 max-md:px-4">
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
                    <div className="flex flex-wrap justify-center gap-3">
                        {Array.from({length: 6}).map((_, i) => (
                            <div key={i} className="h-10 w-28">
                                <Skeleton className="h-full w-full rounded-full"/>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="flex flex-wrap justify-center gap-3 max-md:gap-2">
                        {displayTags.map((tag) => {
                            const isHovered = hoveredId === tag.id;
                            return (
                                <button
                                    key={tag.id}
                                    className="group flex cursor-pointer items-center gap-2 rounded-full border-2 px-5 py-2.5 text-sm font-semibold transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lg max-md:px-4 max-md:py-2 max-md:text-xs"
                                    style={isHovered ? {
                                        backgroundColor: tag.color,
                                        color: '#ffffff',
                                        borderColor: tag.color,
                                        boxShadow: `0 8px 24px -8px ${tag.color}80`,
                                    } : {
                                        backgroundColor: `${tag.color}12`,
                                        color: tag.color,
                                        borderColor: `${tag.color}40`,
                                    }}
                                    onMouseEnter={() => setHoveredId(tag.id)}
                                    onMouseLeave={() => setHoveredId(null)}
                                    onClick={() => handleViewTag(tag)}
                                >
                                    {/* 左侧圆点指示器 */}
                                    <span
                                        className="h-2 w-2 rounded-full transition-all duration-300"
                                        style={{backgroundColor: isHovered ? '#fff' : tag.color}}
                                    />
                                    <span className="truncate tracking-tight">
                                        {tag.name}
                                    </span>
                                    {/* 文章数 */}
                                    {tag.postCount > 0 && (
                                        <span
                                            className="rounded-full px-2 py-0.5 text-[10px] font-bold leading-none transition-colors duration-300"
                                            style={isHovered ? {
                                                backgroundColor: 'rgba(255,255,255,0.25)',
                                                color: '#fff',
                                            } : {
                                                backgroundColor: `${tag.color}25`,
                                            }}
                                        >
                                            {tag.postCount}
                                        </span>
                                    )}
                                </button>
                            );
                        })}
                    </div>
                )}
            </div>
        </section>
    );
}
