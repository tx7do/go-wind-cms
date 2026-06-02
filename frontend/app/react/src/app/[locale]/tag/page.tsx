'use client';

import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import {Skeleton} from '@/components/ui/skeleton';
import {Button} from '@/components/ui/button';
import {AppEmpty} from '@/components/ui';

import {useI18nRouter} from "@/i18n/helpers";
import PageHero from '@/components/layout/PageHero';
import {XIcon} from '@/plugins/xicon';

import {fetchListTags, getTranslation as getTagTranslation} from '@/api/hooks/tag';
import {contentservicev1_ListTagResponse, contentservicev1_Tag} from '@/api/generated/app/service/v1';

export default function TagListPage() {
    const t = useTranslations('page');
    const router = useI18nRouter();

    const [loading, setLoading] = useState(false);
    const [tags, setTags] = useState<contentservicev1_Tag[]>([]);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);

    async function loadTags() {
        setLoading(true);
        try {
            const res = await fetchListTags({
                paging: {
                    page: page,
                    pageSize: pageSize,
                },
                formValues: {
                    status: 'TAG_STATUS_ACTIVE'
                },
                fieldMask: undefined,
                orderBy: undefined,
            }) as unknown as contentservicev1_ListTagResponse;
            setTags(res.items || []);
            setTotal(res.total || 0);
        } catch (error) {
            console.error('Load tags failed:', error);
        } finally {
            setLoading(false);
        }
    }

    function handleTagClick(id: number) {
        router.push(`/tag/${id}`);
    }

    function handlePageChange(newPage: number) {
        setPage(newPage);
        loadTags();
    }

    function handlePageSizeChange(newSize: number) {
        setPageSize(newSize);
        setPage(1);
        loadTags();
    }

    useEffect(() => {
        loadTags();
    }, [page, pageSize]);

    return (
        <div className="w-full">
            <PageHero
                title={t('tags.tags_list')}
                description={t('tags.explore_all')}
                icon="carbon:tag"
                size="md"
                meta={
                    <div className="flex items-center gap-2">
                        <span>{total || tags.length} {t('tags.total_tags')}</span>
                    </div>
                }
            />

            {/* Tags Grid */}
            <div className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                {/* Loading Skeleton */}
                {loading ? (
                    <div className="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6">
                        {Array.from({length: 12}).map((_, i) => (
                            <div key={i}>
                                <Skeleton className="h-[120px] w-full"/>
                            </div>
                        ))}
                    </div>
                ) : (
                    <>
                        {(tags.length > 0 || total > 0) && (
                            <div className="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-6">
                                {tags.map((tag) => (
                                    <div
                                        key={tag.id}
                                        className="group flex cursor-pointer flex-col items-center gap-3 rounded-xl border bg-card p-6 transition-all hover:-translate-y-1 hover:shadow-lg"
                                        style={{
                                            borderColor: tag.color || 'hsl(var(--border))',
                                            background: `linear-gradient(135deg, ${tag.color}10 0%, hsl(var(--card)) 100%)`
                                        }}
                                        onClick={() => handleTagClick(tag.id || 0)}
                                    >
                                        <div
                                            className="flex h-12 w-12 items-center justify-center rounded-lg transition-transform group-hover:scale-110"
                                            style={{color: tag.color || 'hsl(var(--primary))'}}
                                        >
                                            <XIcon name={`carbon:${tag.icon || 'tag'}`} size={48}/>
                                        </div>
                                        <div className="flex flex-1 flex-col items-center gap-2 text-center">
                                            <h3 className="font-semibold text-foreground">
                                                {getTagTranslation(tag)?.name || t('tags.tag_untitled')}
                                            </h3>
                                            <p className="line-clamp-2 text-sm text-muted-foreground">
                                                {getTagTranslation(tag)?.description || ''}
                                            </p>
                                            <div className="mt-2 flex items-center gap-1 text-xs text-muted-foreground">
                                                <XIcon name="carbon:document" size={14}/>
                                                <span>{tag.postCount || 0} {t('posts.articles')}</span>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                                {tags.length === 0 && total > 0 && (
                                    <AppEmpty description={t('tags.no_tags_in_page')}/>
                                )}
                            </div>
                        )}

                        {!loading && tags.length === 0 && total === 0 && (
                            <AppEmpty inContainer description={t('tags.no_tags')}/>
                        )}

                        {/* Simple pagination */}
                        {total > pageSize && (
                            <div className="mt-8 flex items-center justify-center gap-2">
                                <Button
                                    variant="outline"
                                    size="sm"
                                    disabled={page <= 1}
                                    onClick={() => handlePageChange(page - 1)}
                                >
                                    {t('previous') || 'Previous'}
                                </Button>
                                <span className="text-sm text-muted-foreground">
                                    {page} / {Math.ceil(total / pageSize)}
                                </span>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    disabled={page >= Math.ceil(total / pageSize)}
                                    onClick={() => handlePageChange(page + 1)}
                                >
                                    {t('next') || 'Next'}
                                </Button>
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
}
