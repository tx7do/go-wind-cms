'use client';

import {useState, useEffect} from 'react';
import {useSearchParams} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {Skeleton} from '@/components/ui/skeleton';
import {AppEmpty} from '@/components/ui';
import {SearchX, FileText, ChevronRight} from 'lucide-react';

import {fetchSearchPosts} from '@/api/hooks/post';
import PageHero from '@/components/layout/PageHero';
import SectionContainer from '@/components/layout/SectionContainer';

import {useI18nRouter} from '@/i18n/helpers';

interface SearchHit {
    postId?: number;
    language?: string;
    title?: string;
    summary?: string;
    highlightFragments?: string[];
}

interface SearchResponse {
    items?: SearchHit[];
    total?: number;
}

/**
 * 受控高亮渲染：把含 <mark>…</mark> 的片段拆成文本/高亮交替的 React 节点。
 * 不走 innerHTML——片段内容来自我们自己的 ES 高亮标签，但正文本身是用户内容，
 * 按标记拆分渲染可避免任何注入面。
 */
function HighlightedText({fragments, fallback}: {fragments?: string[]; fallback: string}) {
    // 优先用命中片段；无片段时按纯文本展示
    const text = fragments?.[0] ?? fallback;
    if (!text) return null;
    const parts = text.split(/(<mark>.*?<\/mark>)/g).filter(Boolean);
    return (
        <>
            {parts.map((part, i) => {
                if (part.startsWith('<mark>') && part.endsWith('</mark>')) {
                    return (
                        <mark key={i} className="rounded-sm bg-primary/15 px-0.5 text-primary">
                            {part.slice(6, -7)}
                        </mark>
                    );
                }
                return <span key={i}>{part}</span>;
            })}
        </>
    );
}

export default function SearchPage() {
    const t = useTranslations('page');
    const router = useI18nRouter();
    const searchParams = useSearchParams();
    const query = searchParams.get('q')?.trim() ?? '';

    const [loading, setLoading] = useState(false);
    const [hits, setHits] = useState<SearchHit[]>([]);
    const [total, setTotal] = useState(0);

    useEffect(() => {
        if (!query) {
            // 宏任务清态：规避 set-state-in-effect 规则
            const clear = setTimeout(() => {
                setHits([]);
                setTotal(0);
            }, 0);
            return () => clearTimeout(clear);
        }
        let cancelled = false;
        setLoading(true);
        fetchSearchPosts({query, pageSize: 20})
            .then((res) => {
                if (cancelled) return;
                const data = res as unknown as SearchResponse;
                setHits(data?.items ?? []);
                setTotal(data?.total ?? 0);
            })
            .catch((err) => {
                console.error('Search failed:', err);
                if (cancelled) return;
                setHits([]);
                setTotal(0);
            })
            .finally(() => {
                if (!cancelled) setLoading(false);
            });
        return () => {
            cancelled = true;
        };
    }, [query]);

    const openPost = (postId?: number) => {
        if (!postId) return;
        router.push(`/post/${postId}`);
        if (typeof window !== 'undefined') {
            window.scrollTo({top: 0, behavior: 'smooth'});
        }
    };

    return (
        <div className="w-full">
            <PageHero
                title={t('posts.search')}
                description={query ? `"${query}"` : undefined}
                icon="lucide:search"
                size="md"
            />

            <SectionContainer>
                {/* 加载中 */}
                {loading ? (
                    <div className="grid gap-3">
                        {[...Array(5)].map((_, i) => (
                            <Skeleton key={i} className="h-16 w-full"/>
                        ))}
                    </div>
                ) : !query ? (
                    <AppEmpty
                        description={t('posts.search_placeholder')}
                        inContainer
                        image={<span className="i-carbon:search text-[64px]"/>}
                    />
                ) : hits.length === 0 ? (
                    <AppEmpty
                        description={t('posts.no_search_results', {query})}
                        inContainer
                        image={<SearchX className="h-16 w-16 text-muted-foreground/40"/>}
                    />
                ) : (
                    <div>
                        <p className="mb-6 text-sm text-muted-foreground">
                            {t('posts.search_results_count', {count: total})}
                        </p>
                        <div className="grid gap-3">
                            {hits.map((hit, idx) => (
                                <article
                                    key={hit.postId ?? idx}
                                    onClick={() => openPost(hit.postId)}
                                    className="group flex cursor-pointer items-center gap-3 rounded-xl border border-border bg-card p-5 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
                                >
                                    <FileText className="h-5 w-5 shrink-0 text-muted-foreground transition-colors group-hover:text-primary"/>
                                    <div className="min-w-0 flex-1">
                                        <h3 className="line-clamp-2 text-base font-semibold text-foreground transition-colors group-hover:text-primary">
                                            <HighlightedText
                                                fragments={hit.highlightFragments}
                                                fallback={hit.title ?? ''}
                                            />
                                        </h3>
                                        {(hit.summary || hit.highlightFragments?.length) && (
                                            <p className="mt-1 line-clamp-1 text-sm text-muted-foreground">
                                                <HighlightedText
                                                    fragments={
                                                        hit.highlightFragments && hit.highlightFragments.length > 1
                                                            ? hit.highlightFragments.slice(1)
                                                            : undefined
                                                    }
                                                    fallback={hit.summary ?? ''}
                                                />
                                            </p>
                                        )}
                                    </div>
                                    <ChevronRight className="h-[18px] w-[18px] shrink-0 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100"/>
                                </article>
                            ))}
                        </div>
                    </div>
                )}
            </SectionContainer>
        </div>
    );
}
