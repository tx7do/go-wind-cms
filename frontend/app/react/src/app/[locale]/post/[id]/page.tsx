'use client';

import {useEffect, useState, useRef, useMemo, useCallback} from 'react';
import {useParams, useSearchParams} from 'next/navigation';
import {useTranslations} from 'next-intl';

import CommentSection from '@/components/comment/CommentSection';
import ContentViewer from '@/components/content/ContentViewer';
import PostList from '@/components/post/PostList';
import BackToTop from '@/components/layout/BackToTop';

import {
    fetchPost,
    getPostTitle,
    getPostContent,
    getPostThumbnail,
} from '@/api/hooks/post';
import {formatDate} from "@/utils";
import {contentservicev1_Post} from "@/api/generated/app/service/v1";
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from "@/i18n/helpers";

import '../../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用

// 常量定义
const HEADING_OFFSET = 150;
const THROTTLE_DELAY = 200;

interface TocItem {
    id: string;
    level: number;
    text: string;
    element: Element;
}

export default function PostDetailPage() {
    const t = useTranslations('page');
    const params = useParams();
    const router = useI18nRouter();
    const searchParams = useSearchParams();

    // 本地 state 管理 post 数据（不再使用 Redux store）
    const [post, setPost] = useState<contentservicev1_Post | null>(null);
    const [postLoading, setPostLoading] = useState(false);

    const [localLoading, setLocalLoading] = useState(true);
    const isLoading = localLoading || (postLoading && !post);
    const [tableOfContents, setTableOfContents] = useState<TocItem[]>([]);
    const [activeHeading, setActiveHeading] = useState('');
    const [isLiked, setIsLiked] = useState(false);
    const [isBookmarked, setIsBookmarked] = useState(false);
    const [isTocExpanded, setIsTocExpanded] = useState(true);

    const contentRef = useRef<HTMLDivElement>(null);

    const postId = useMemo(() => {
        const id = params?.id as string;
        return id ? parseInt(id) : null;
    }, [params?.id]);

    // 计算属性 - 使用纯工具函数
    const displayTitle = useMemo(() => {
        if (!post) return '';
        return getPostTitle(post);
    }, [post]);

    const displayContent = useMemo(() => {
        if (!post) return '';
        return getPostContent(post);
    }, [post]);

    const displayThumbnail = useMemo(() => {
        if (!post) return '';
        return getPostThumbnail(post);
    }, [post]);

    // 相关文章数据
    const relatedPostsQuery = useMemo(() => {
        if (!post?.categoryIds) return null;
        return {
            status: 'POST_STATUS_PUBLISHED',
            id__not: postId,
            category_ids__in: post.categoryIds
        };
    }, [post?.categoryIds, postId]);

    // Load post data
    useEffect(() => {
        async function loadPost() {
            if (!postId) return;

            setLocalLoading(true);
            setPostLoading(true);
            try {
                console.log('[PostDetail] Before fetchPost call:', {
                    postId,
                });

                const fetchedPost = (await fetchPost(postId!)) as contentservicev1_Post;

                console.log('[PostDetail] API returned:', {
                    hasFetchedPost: !!fetchedPost,
                    hasTranslations: !!fetchedPost?.translations?.length,
                    firstTranslationTitle: fetchedPost?.translations?.[0]?.title,
                    firstTranslationContentLength: fetchedPost?.translations?.[0]?.content?.length
                });

                setPost(fetchedPost);

                if (fetchedPost) {
                    // SEO
                    document.title = `${getPostTitle(fetchedPost)} - GoWind Content Hub`;
                }
            } catch (error) {
                console.error('Load post failed:', error);
            } finally {
                setLocalLoading(false);
                setPostLoading(false);
            }
        }

        loadPost();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [postId]); // 只依赖 postId

    // Generate table of contents - 在内容渲染后生成目录
    useEffect(() => {
        if (!displayContent || !contentRef.current) {
            console.log('[TOC] Skip generation:', {hasContent: !!displayContent, hasRef: !!contentRef.current});
            return;
        }

        const timeoutId = setTimeout(() => {
            const contentEl = contentRef.current;
            if (!contentEl) {
                console.log('[TOC] No content element');
                return;
            }

            const headings = contentEl.querySelectorAll('h2, h3');
            const toc: TocItem[] = [];

            console.log('[TOC] Found headings:', headings.length);

            headings.forEach((heading, index) => {
                const level = heading.tagName === 'H2' ? 2 : 3;
                const id = `heading-${index}`;
                if (!heading.id) heading.setAttribute('id', id);

                toc.push({
                    id,
                    level,
                    text: heading.textContent || '',
                    element: heading
                });
            });

            setTableOfContents(toc);
            console.log('[TOC] Generated:', toc.length, 'items');

            // 页面加载时检查 URL hash，自动滚动到对应位置
            if (window.location.hash) {
                const hashId = window.location.hash.slice(1);
                setTimeout(() => {
                    scrollToHeading(hashId);
                }, 300);
            }
        }, 100); // 等待内容完全渲染

        return () => clearTimeout(timeoutId);
    }, [displayContent]); // 依赖计算后的内容

    // 监听内容变化，重新生成目录
    useEffect(() => {
        if (displayContent && tableOfContents.length === 0) {
            generateTableOfContents();
        }
    }, [displayContent, tableOfContents.length]);

    // 生成目录函数
    const generateTableOfContents = useCallback(() => {
        setTimeout(() => {
            if (!contentRef.current) return;

            const contentEl = contentRef.current;
            const headings = contentEl.querySelectorAll('h2, h3');
            const toc: TocItem[] = [];

            console.log('[GenerateTOC] Found headings:', headings.length);

            headings.forEach((heading, index) => {
                const level = heading.tagName === 'H2' ? 2 : 3;
                const id = `heading-${index}`;

                if (!heading.id) {
                    heading.setAttribute('id', id);
                    console.log(`[GenerateTOC] Set ID "${id}" on ${heading.tagName}:`, heading.textContent);
                }

                toc.push({
                    id,
                    level,
                    text: heading.textContent || '',
                    element: heading
                });
            });

            setTableOfContents(toc);
            console.log('[GenerateTOC] Generated TOC:', toc);
        }, 500);
    }, []);

    // 防抖函数
    const throttle = useCallback((fn: (...args: unknown[]) => void, delay: number) => {
        let timer: ReturnType<typeof setTimeout> | null = null;
        return function (this: unknown, ...args: unknown[]) {
            if (!timer) {
                timer = setTimeout(() => {
                    fn.apply(this as (...args: unknown[]) => void, args as unknown[]);
                    timer = null;
                }, delay);
            }
        };
    }, []);

    // 滚动处理函数
    const handleScroll = useCallback(() => {
        if (tableOfContents.length > 0) {
            let currentActive = activeHeading;
            for (const heading of tableOfContents) {
                if (heading.element && heading.element.getBoundingClientRect().top < HEADING_OFFSET) {
                    currentActive = heading.id;
                }
            }
            setActiveHeading(currentActive);
        }
    }, [tableOfContents, activeHeading]);

    const throttledScroll = useMemo(() => {
        return throttle(handleScroll, THROTTLE_DELAY);
    }, [throttle, handleScroll]);

    // Scroll handler
    useEffect(() => {
        window.addEventListener('scroll', throttledScroll);
        return () => window.removeEventListener('scroll', throttledScroll);
    }, [throttledScroll]);

    // Handlers
    const handleBack = () => {
        const from = searchParams.get('from');
        const categoryId = searchParams.get('categoryId');

        if (from === 'category' && categoryId) {
            router.push(`/category/${categoryId}`);
        } else if (from === 'user') {
            router.push('/user');
        } else if (from === 'post') {
            router.push('/post');
        } else if (from === 'home') {
            router.push('/');
        } else if (window.history.length > 2) {
            router.back();
        } else {
            router.push('/post');
        }
    };

    const handleLike = () => {
        setIsLiked(!isLiked);
    };

    const handleBookmark = () => {
        setIsBookmarked(!isBookmarked);
    };

    const handleShare = () => {
        const url = window.location.href;
        if (navigator.share) {
            navigator.share({title: displayTitle, url}).then(() => {
                console.log('Shared successfully');
            }).catch(() => {
                copyToClipboard(url);
            });
        } else {
            copyToClipboard(url);
        }
    };

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text).then(() => {
            console.log('Copied!');
        }).catch(() => {
            console.error('Copy failed');
        });
    };

    const scrollToHeading = (id: string) => {
        const tocItem = tableOfContents.find(item => item.id === id);
        const element = tocItem?.element || document.getElementById(id);

        if (element) {
            console.log('[TOC Scroll] Start',
                'ID:', id,
                '| Element:', element.textContent?.trim()
            );

            element.scrollIntoView({behavior: 'auto', block: 'start'});

            setTimeout(() => {
                const currentScroll = window.pageYOffset;
                const headerOffset = 100;
                const targetPosition = currentScroll - headerOffset;

                console.log('[TOC Scroll] Adjusting:',
                    'from:', currentScroll.toFixed(2),
                    'to:', targetPosition.toFixed(2)
                );

                window.scrollTo({
                    top: Math.max(0, targetPosition),
                    behavior: 'smooth'
                });

                setActiveHeading(id);

                if (window.history.pushState) {
                    const currentState = window.history.state || {};
                    window.history.pushState(currentState, '', `#${id}`);
                } else {
                    window.location.hash = `#${id}`;
                }
            }, 50);
        } else {
            console.error('[ScrollToHeading] Element not found:', id);
        }
    };

    if (isLoading) {
        return (
            <div className="w-full">
                {/* Loading skeleton */}
                <div className="w-full max-w-[1200px] mx-auto px-8 py-8 max-md:px-4">
                    <div className="mb-6 h-9 w-32 animate-pulse rounded bg-muted"/>
                    <article className="overflow-hidden rounded-xl border border-border bg-card shadow-sm">
                        <div className="h-[300px] animate-pulse bg-muted max-md:h-[200px]"/>
                        <div className="flex gap-6 p-8 max-md:flex-col max-md:p-4">
                            <aside className="w-[240px] shrink-0 max-md:hidden">
                                <div className="sticky top-24 space-y-3 rounded-lg border border-border bg-background p-4">
                                    <div className="h-6 w-[200px] animate-pulse rounded bg-muted"/>
                                    <div className="mt-4 h-5 w-[180px] animate-pulse rounded bg-muted"/>
                                    <div className="mt-2 h-5 w-[160px] animate-pulse rounded bg-muted"/>
                                    <div className="mt-2 h-5 w-[140px] animate-pulse rounded bg-muted"/>
                                </div>
                            </aside>
                            <div className="flex-1">
                                <div className="mb-4 h-12 w-[80%] animate-pulse rounded bg-muted"/>
                                <div className="mb-4 h-8 w-[60%] animate-pulse rounded bg-muted"/>
                                <div className="mb-6 flex gap-4">
                                    <div className="h-5 w-[100px] animate-pulse rounded bg-muted"/>
                                    <div className="h-5 w-[100px] animate-pulse rounded bg-muted"/>
                                    <div className="h-5 w-[100px] animate-pulse rounded bg-muted"/>
                                    <div className="h-5 w-[100px] animate-pulse rounded bg-muted"/>
                                </div>
                                <div className="space-y-3">
                                    <div className="h-4 w-full animate-pulse rounded bg-muted"/>
                                    <div className="h-4 w-[90%] animate-pulse rounded bg-muted"/>
                                    <div className="h-4 w-[95%] animate-pulse rounded bg-muted"/>
                                </div>
                            </div>
                        </div>
                    </article>
                </div>
            </div>
        );
    }

    if (!post) {
        return (
            <div className="flex w-full items-center justify-center py-20">
                <div className="text-lg text-muted-foreground">Post not found</div>
            </div>
        );
    }

    return (
        <div className="w-full">
            {/* Back Navigation */}
            <div className="w-full max-w-[1200px] mx-auto px-8 py-6 max-md:px-4">
                <button
                    onClick={handleBack}
                    className="flex cursor-pointer items-center gap-2 rounded-lg border border-border bg-card px-4 py-2 text-sm text-foreground transition-all hover:border-primary hover:bg-primary/5"
                    aria-label={t('post_detail.back')}
                >
                    <XIcon name="carbon:arrow-left"/>
                    <span>{t('post_detail.back')}</span>
                </button>
            </div>

            {/* Post Article */}
            <article className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                {/* Post Thumbnail Banner */}
                {displayThumbnail && (
                    <div className="relative mb-8 h-[300px] overflow-hidden rounded-xl max-md:h-[200px]">
                        <img src={displayThumbnail} alt={displayTitle} className="h-full w-full object-cover"/>
                        <div className="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent"/>
                    </div>
                )}

                <div className={`flex gap-6 max-md:flex-col ${!isTocExpanded ? 'lg:pl-12' : ''}`}>
                    {/* Table of Contents Sidebar */}
                    {tableOfContents.length > 0 && isTocExpanded && (
                        <aside className="w-[240px] shrink-0 max-md:hidden">
                            <div className="sticky top-24 rounded-lg border border-border bg-card p-4">
                                <div className="mb-3 flex items-center justify-between">
                                    <h3 className="flex items-center gap-2 text-sm font-semibold text-foreground">
                                        <XIcon name="carbon:list"/>
                                        <span>{t('post_detail.table_of_contents')}</span>
                                    </h3>
                                    <button
                                        onClick={() => setIsTocExpanded(false)}
                                        className="cursor-pointer rounded p-1 text-muted-foreground transition-colors hover:bg-primary/5 hover:text-primary"
                                    >
                                        <XIcon name="carbon:chevron-left"/>
                                    </button>
                                </div>
                                <nav className="space-y-1">
                                    {tableOfContents.map(item => (
                                        <a
                                            key={item.id}
                                            href={`#${item.id}`}
                                            className={`block truncate rounded px-2 py-1.5 text-sm transition-colors hover:bg-primary/5 hover:text-primary ${
                                                item.level === 3 ? 'pl-6 text-muted-foreground' : 'text-foreground'
                                            } ${activeHeading === item.id ? 'bg-primary/10 font-medium text-primary' : ''}`}
                                            onClick={(e) => {
                                                e.preventDefault();
                                                setTimeout(() => scrollToHeading(item.id), 10);
                                            }}
                                        >
                                            {item.text}
                                        </a>
                                    ))}
                                </nav>
                            </div>
                        </aside>
                    )}

                    {/* TOC Expand Button */}
                    {tableOfContents.length > 0 && !isTocExpanded && (
                        <div className="fixed top-1/2 left-4 z-10 -translate-y-1/2">
                            <button
                                onClick={() => setIsTocExpanded(true)}
                                className="flex cursor-pointer items-center gap-2 rounded-lg border border-border bg-card px-3 py-2 text-sm text-foreground shadow-md transition-all hover:border-primary hover:bg-primary/5"
                            >
                                <XIcon name="carbon:list"/>
                                <span className="max-md:hidden">{t('post_detail.table_of_contents')}</span>
                                <XIcon name="carbon:chevron-right"/>
                            </button>
                        </div>
                    )}

                    <div className="flex-1 min-w-0">
                        {/* Post Header */}
                        <header className="mb-8">
                            <h1 className="mb-4 text-3xl font-bold text-foreground max-md:text-2xl">
                                {displayTitle}
                            </h1>
                            <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
                                <div className="flex items-center gap-1.5">
                                    <XIcon name="carbon:user-avatar"/>
                                    <span>{post.authorName}</span>
                                </div>
                                <div className="flex items-center gap-1.5">
                                    <XIcon name="carbon:calendar"/>
                                    <span>{formatDate(post.createdAt)}</span>
                                </div>
                                <div className="flex items-center gap-1.5">
                                    <XIcon name="carbon:view"/>
                                    <span>{post.visits || 0}</span>
                                </div>
                                <div className="flex items-center gap-1.5">
                                    <XIcon name="carbon:thumbs-up"/>
                                    <span>{post.likes || 0}</span>
                                </div>
                            </div>
                        </header>

                        {/* Post Content */}
                        <div className="mb-8" ref={contentRef}>
                            <ContentViewer content={displayContent} type="markdown"/>
                        </div>

                        {/* Post Actions */}
                        <footer className="mb-8 border-t border-border pt-6">
                            <div className="flex items-center gap-3">
                                <button
                                    onClick={handleLike}
                                    className={`flex h-10 w-10 cursor-pointer items-center justify-center rounded-full border border-border bg-card transition-all hover:border-primary hover:bg-primary/5 ${isLiked ? 'border-primary bg-primary/10 text-primary' : 'text-muted-foreground'}`}
                                    aria-label={t('post_detail.likes')}
                                >
                                    <XIcon name={isLiked ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'}/>
                                </button>
                                <button
                                    onClick={handleBookmark}
                                    className={`flex h-10 w-10 cursor-pointer items-center justify-center rounded-full border border-border bg-card transition-all hover:border-primary hover:bg-primary/5 ${isBookmarked ? 'border-primary bg-primary/10 text-primary' : 'text-muted-foreground'}`}
                                    aria-label={t('post_detail.bookmark')}
                                >
                                    <XIcon name={isBookmarked ? 'carbon:bookmark-filled' : 'carbon:bookmark'}/>
                                </button>
                                <button
                                    onClick={handleShare}
                                    className="flex h-10 w-10 cursor-pointer items-center justify-center rounded-full border border-border bg-card text-muted-foreground transition-all hover:border-primary hover:bg-primary/5 hover:text-primary"
                                    aria-label={t('post_detail.share')}
                                >
                                    <XIcon name="carbon:share"/>
                                </button>
                            </div>
                        </footer>
                    </div>
                </div>
            </article>

            {/* Comments Section */}
            <div className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                <CommentSection
                objectId={postId}
                contentType="CONTENT_TYPE_POST"
                onUpdateComments={() => {
                }}
                />
            </div>

            {/* Related Posts */}
            <section className="w-full max-w-[1200px] mx-auto px-8 py-12 max-md:px-4">
                <div className="mb-6">
                    <h2 className="flex items-center gap-2 text-2xl font-bold text-foreground">
                        <XIcon name="carbon:book"/>
                        <span>{t('post_detail.related_posts')}</span>
                    </h2>
                </div>
                {relatedPostsQuery && (
                    <PostList
                        queryParams={relatedPostsQuery}
                        fieldMask="id,status,sort_order,is_featured,visits,likes,comment_count,author_name,available_languages,created_at,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,translations.thumbnail"
                        orderBy={['-sortOrder']}
                        page={1}
                        pageSize={3}
                        showSkeleton={false}
                    />
                )}
            </section>

            {/* Back to Top Button */}
            <BackToTop scrollThreshold={500}/>
        </div>
    );
}
