'use client';

import {useEffect, useMemo, useState} from 'react';
import {useSearchParams} from 'next/navigation';
import {useTranslations} from 'next-intl';

import type {contentservicev1_Page} from '@/api/generated/app/service/v1';
import {fetchPage, fetchPageBySlug, fetchPagePreview} from '@/api/hooks/page';
import CmsPageView from '@/components/page/CmsPageView';
import SectionContainer from '@/components/layout/SectionContainer';

/**
 * CMS Page 公开渲染页（静态导出，query 传参）：
 *   /[locale]/page/detail?slug=about  或  ?id=5
 * 仅返回已发布页面（app BFF Get 过滤）；携带 preview_token 时放行草稿供后台预览。
 */
export default function PageDetailClient() {
    const t = useTranslations('page');
    const searchParams = useSearchParams();

    const slug = useMemo(() => searchParams.get('slug') || '', [searchParams]);
    const id = useMemo(() => {
        const raw = searchParams.get('id');
        return raw ? parseInt(raw) : 0;
    }, [searchParams]);
    const previewToken = useMemo(() => searchParams.get('preview_token') || '', [searchParams]);

    const [page, setPage] = useState<contentservicev1_Page | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        let cancelled = false;
        async function load() {
            if (!slug && !id) {
                setLoading(false);
                return;
            }
            setLoading(true);
            try {
                const fetched = previewToken
                    ? await fetchPagePreview(id, previewToken)
                    : slug
                      ? await fetchPageBySlug(slug)
                      : await fetchPage(id);
                if (!cancelled) setPage((fetched as contentservicev1_Page | null) ?? null);
            } catch (error) {
                console.error('Load cms page failed:', error);
                if (!cancelled) setPage(null);
            } finally {
                if (!cancelled) setLoading(false);
            }
        }
        load();
        return () => {
            cancelled = true;
        };
    }, [slug, id, previewToken]);

    // 标题随页面翻译更新
    useEffect(() => {
        const tr = page?.translations?.find((x) => x.title);
        if (tr?.title) document.title = `${tr.title} - GoWind Content Hub`;
    }, [page]);

    if (loading) {
        return (
            <div className="mx-auto w-full max-w-[900px] px-8 py-12 max-md:px-4">
                <div className="mb-8 h-9 w-32 animate-pulse rounded bg-muted"/>
                <div className="mb-6 h-10 w-2/3 animate-pulse rounded bg-muted"/>
                <div className="space-y-4">
                    <div className="h-4 w-full animate-pulse rounded bg-muted"/>
                    <div className="h-4 w-[90%] animate-pulse rounded bg-muted"/>
                    <div className="h-4 w-[95%] animate-pulse rounded bg-muted"/>
                </div>
            </div>
        );
    }

    if (!page) {
        return (
            <SectionContainer>
                <div className="flex w-full items-center justify-center py-20">
                    <div className="text-lg text-muted-foreground">
                        {t('post_detail.post_not_found')}
                    </div>
                </div>
            </SectionContainer>
        );
    }

    return <CmsPageView page={page} backLabel={t('post_detail.back')}/>;
}
