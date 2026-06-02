'use client';

import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import {Skeleton} from '@/components/ui/skeleton';
import {Button} from '@/components/ui/button';
import {AppEmpty} from '@/components/ui';

import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from "@/i18n/helpers";

import {fetchListTags, getTranslation as getTagTranslation} from '@/api/hooks/tag';
import {contentservicev1_ListTagResponse, contentservicev1_Tag} from '@/api/generated/app/service/v1';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './tag-list.module.css';

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
        <div className={styles['tag-list-page']}>
            {/* Hero Section */}
            <div className={styles['hero-section']}>
                <div className={styles['hero-content']}>
                    <h1>{t('tags.tags_list')}</h1>
                    <p>{t('tags.explore_all')}</p>
                    <div className={styles['tag-stats']}>
                        <div className={styles['stat-item']}>
                            <XIcon name="carbon:tag" size={20}/>
                            <span>{tags.length} {t('tags.total_tags')}</span>
                        </div>
                    </div>
                </div>
            </div>

            {/* Tags Grid */}
            <div className={styles['page-container']}>
                {/* Loading Skeleton */}
                {loading ? (
                    <div className={styles['tags-grid']}>
                        {Array.from({length: 12}).map((_, i) => (
                            <div key={i} className={styles['tag-card-skeleton']}>
                                <Skeleton className="w-full" style={{height: 120}}/>
                            </div>
                        ))}
                    </div>
                ) : (
                    <>
                        {(tags.length > 0 || total > 0) && (
                            <div className={styles['tags-grid']}>
                                {tags.map((tag) => (
                                    <div
                                        key={tag.id}
                                        className={styles['tag-card']}
                                        style={{
                                            borderColor: tag.color || 'var(--color-border)',
                                            background: `linear-gradient(135deg, ${tag.color}10 0%, var(--color-surface) 100%)`
                                        }}
                                        onClick={() => handleTagClick(tag.id || 0)}
                                    >
                                        <div
                                            className={styles['tag-icon']}
                                            style={{color: tag.color || 'var(--color-brand)'}}
                                        >
                                            <XIcon name={`carbon:${tag.icon || 'tag'}`} size={48}/>
                                        </div>
                                        <div className={styles['tag-content']}>
                                            <h3>{getTagTranslation(tag)?.name || t('tags.tag_untitled')}</h3>
                                            <p className={styles['tag-description']}>
                                                {getTagTranslation(tag)?.description || ''}
                                            </p>
                                            <div className={styles['tag-meta']}>
                                                <span className={styles['meta-icon']}>
                                                    <XIcon name="carbon:document" size={14}/>
                                                </span>
                                                <span className={styles['meta-text']}>
                                                    {tag.postCount || 0} {t('posts.articles')}
                                                </span>
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
                            <div className="flex items-center justify-center gap-2 mt-8">
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
