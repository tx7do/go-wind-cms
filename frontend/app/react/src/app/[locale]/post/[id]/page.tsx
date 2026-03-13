'use client';

import {useEffect, useState, useRef, useMemo, useCallback} from 'react';
import {useParams, useRouter, useSearchParams} from 'next/navigation';
import {useTranslations} from 'next-intl';

import CommentSection from '@/components/comment/CommentSection';
import ContentViewer from '@/components/content/ContentViewer';
import PostList from '@/components/post/PostList';

import {usePostStore} from '@/store/slices/post/hooks';
import {formatDate} from "@/utils";
import {contentservicev1_Post} from "@/api/generated/app/service/v1";
import XIcon from '@/plugins/xicon';

import '../../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './post-detail.module.css';

// 常量定义
const SCROLL_THRESHOLD = 500;
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
    const router = useRouter();
    const searchParams = useSearchParams();
    const postStore = usePostStore();

    // 直接使用 store 中的数据，而不是本地 state
    const post = postStore.detail as contentservicev1_Post | null;

    const [localLoading, setLocalLoading] = useState(true);
    const isLoading = localLoading || (postStore.loading && !post);
    const [tableOfContents, setTableOfContents] = useState<TocItem[]>([]);
    const [activeHeading, setActiveHeading] = useState('');
    const [showBackToTop, setShowBackToTop] = useState(false);
    const [isLiked, setIsLiked] = useState(false);
    const [isBookmarked, setIsBookmarked] = useState(false);
    const [isTocExpanded, setIsTocExpanded] = useState(true);

    const relatedPostListRef = useRef<HTMLDivElement>(null);
    const contentRef = useRef<HTMLDivElement>(null);

    const postId = useMemo(() => {
        const id = params?.id as string;
        return id ? parseInt(id) : null;
    }, [params?.id]);

    // 计算属性 - 使用 postStore 提供的工具函数
    const displayTitle = useMemo(() => {
        if (!post) return '';
        return postStore.getPostTitle(post);
    }, [post, postStore]);

    const displayContent = useMemo(() => {
        if (!post) return '';
        return postStore.getPostContent(post);
    }, [post, postStore]);

    const displayThumbnail = useMemo(() => {
        if (!post) return '';
        return postStore.getPostThumbnail(post);
    }, [post, postStore]);

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
            try {
                console.log('[PostDetail] Before getPost call:', {
                    postId,
                    currentDetail: postStore.detail
                });

                const fetchedPost = (await postStore.getPost({
                    // @ts-expect-error - listNavigation 参数类型推断问题
                    id: postId
                })) as contentservicev1_Post;

                console.log('[PostDetail] API returned:', {
                    hasFetchedPost: !!fetchedPost,
                    hasTranslations: !!fetchedPost?.translations?.length,
                    firstTranslationTitle: fetchedPost?.translations?.[0]?.title,
                    firstTranslationContentLength: fetchedPost?.translations?.[0]?.content?.length
                });

                if (fetchedPost) {
                    // SEO
                    document.title = `${postStore.getPostTitle(fetchedPost)} - GoWind Content Hub`;
                }
            } catch (error) {
                console.error('Load post failed:', error);
            } finally {
                setLocalLoading(false);
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
            // 如果还没有生成目录，立即生成一次
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
            console.log('Table of contents regenerated:', toc.length, 'items');
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
        setShowBackToTop(window.scrollY > SCROLL_THRESHOLD);

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
        // TODO: Add toast message
    };

    const handleBookmark = () => {
        setIsBookmarked(!isBookmarked);
        // TODO: Add toast message
    };

    const handleShare = () => {
        const url = window.location.href;
        if (navigator.share) {
            navigator.share({title: post?.title, url}).then(() => {
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
            const offset = 100;
            const elementPosition = element.getBoundingClientRect().top + window.pageYOffset;
            window.scrollTo({
                top: elementPosition - offset,
                behavior: 'smooth'
            });
            setActiveHeading(id);
            window.history.pushState({}, '', `#${id}`);
        }
    };

    const scrollToTop = () => {
        window.scrollTo({top: 0, behavior: 'smooth'});
    };

    if (isLoading) {
        return (
            <div className={styles['post-detail-page']}>
                {/* Loading skeleton */}
                <div className={styles['back-navigation']}>
                    <div className={styles['skeleton-btn']}></div>
                </div>
                <article className={styles['post-article']}>
                    <div className={styles['post-banner']}>
                        <div className={styles['skeleton-image']}></div>
                    </div>
                    <div className={styles['post-wrapper']}>
                        <aside className={styles['toc-sidebar']}>
                            <div className={styles['toc-container']}>
                                <div className={styles['skeleton-line']} style={{width: '200px', height: '24px'}}></div>
                                <div className={styles['skeleton-line']}
                                     style={{width: '180px', height: '20px', marginTop: '16px'}}></div>
                                <div className={styles['skeleton-line']}
                                     style={{width: '160px', height: '20px', marginTop: '8px'}}></div>
                                <div className={styles['skeleton-line']}
                                     style={{width: '140px', height: '20px', marginTop: '8px'}}></div>
                            </div>
                        </aside>
                        <div className={styles['article-content']}>
                            <header className={styles['post-header']}>
                                <div className={styles['skeleton-title']} style={{width: '80%', height: '48px'}}></div>
                                <div className={styles['skeleton-subtitle']}
                                     style={{width: '60%', height: '32px', marginTop: '16px'}}></div>
                                <div className={styles['post-meta']}>
                                    <div className={styles['skeleton-meta']}
                                         style={{width: '100px', height: '20px'}}></div>
                                    <div className={styles['skeleton-meta']}
                                         style={{width: '100px', height: '20px'}}></div>
                                    <div className={styles['skeleton-meta']}
                                         style={{width: '100px', height: '20px'}}></div>
                                    <div className={styles['skeleton-meta']}
                                         style={{width: '100px', height: '20px'}}></div>
                                </div>
                            </header>
                            <div className={styles['post-content']}>
                                <div className={styles['skeleton-paragraph']}></div>
                            </div>
                        </div>
                    </div>
                </article>
            </div>
        );
    }

    if (!post) {
        return (
            <div className={styles['post-detail-page']}>
                <div className={styles['empty-state']}>Post not found</div>
            </div>
        );
    }

    return (
        <div className={styles['post-detail-page']}>
            {/* Back Navigation */}
            <div className={styles['back-navigation']}>
                <button onClick={handleBack} className={styles['back-btn']}>
                    <XIcon name="carbon:arrow-left" />
                    <span>Back</span>
                </button>
            </div>

            {/* Post Article */}
            <article className={styles['post-article']}>
                {/* Post Thumbnail Banner */}
                {displayThumbnail && (
                    <div className={styles['post-banner']}>
                        <img src={displayThumbnail} alt={displayTitle}/>
                        <div className={styles['banner-overlay']}/>
                    </div>
                )}

                <div className={`${styles['post-wrapper']} ${!isTocExpanded ? styles['toc-collapsed'] : ''}`}>
                    {/* Table of Contents Sidebar */}
                    {tableOfContents.length > 0 && isTocExpanded && (
                        <aside className={styles['toc-sidebar']}>
                            <div className={styles['toc-container']}>
                                <div className={styles['toc-header']}>
                                    <h3 className={styles['toc-title']}>
                                        <XIcon name="carbon:list" />
                                        <span>{t('post_detail.table_of_contents')}</span>
                                    </h3>
                                    <button
                                        onClick={() => setIsTocExpanded(false)}
                                        className={styles['toc-collapse-btn']}
                                    >
                                        <XIcon name="carbon:chevron-left" />
                                    </button>
                                </div>
                                <nav className={styles['toc-list']}>
                                    {tableOfContents.map(item => (
                                        <a
                                            key={item.id}
                                            href={`#${item.id}`}
                                            className={`${styles['toc-item']} ${styles[`level-${item.level}`]} ${activeHeading === item.id ? styles['active'] : ''}`}
                                            onClick={(e) => {
                                                e.preventDefault();
                                                scrollToHeading(item.id);
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
                        <div className={styles['toc-expand-trigger']}>
                            <button onClick={() => setIsTocExpanded(true)}>
                                <XIcon name="carbon:list" />
                                <span>{t('post_detail.table_of_contents')}</span>
                                <XIcon name="carbon:chevron-right" />
                            </button>
                        </div>
                    )}

                    <div className={styles['article-content']}>
                        {/* Post Header */}
                        <header className={styles['post-header']}>
                            <h1 className={styles['post-title']}>{displayTitle}</h1>
                            <div className={styles['post-meta']}>
                                <div className={styles['meta-item']}>
                                    <XIcon name="carbon:user-avatar" />
                                    <span>{post.authorName}</span>
                                </div>
                                <div className={styles['meta-item']}>
                                    <XIcon name="carbon:calendar" />
                                    <span>{formatDate(post.createdAt)}</span>
                                </div>
                                <div className={styles['meta-item']}>
                                    <XIcon name="carbon:view" />
                                    <span>{post.visits || 0}</span>
                                </div>
                                <div className={styles['meta-item']}>
                                    <XIcon name="carbon:thumbs-up" />
                                    <span>{post.likes || 0}</span>
                                </div>
                            </div>
                        </header>

                        {/* Post Content */}
                        <div className={styles['post-content']} ref={contentRef}>
                            <ContentViewer content={displayContent} type="markdown"/>
                        </div>

                        {/* Post Actions */}
                        <footer className={styles['post-actions']}>
                            <div className={styles['action-buttons']}>
                                <button
                                    onClick={handleLike}
                                    className={`${styles['action-btn']} ${isLiked ? styles['liked'] : ''}`}
                                    aria-label={t('post_detail.likes')}
                                >
                                    <XIcon name={isLiked ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'} />
                                </button>
                                <button
                                    onClick={handleBookmark}
                                    className={`${styles['action-btn']} ${isBookmarked ? styles['bookmarked'] : ''}`}
                                    aria-label={t('post_detail.bookmark')}
                                >
                                    <XIcon name={isBookmarked ? 'carbon:bookmark-filled' : 'carbon:bookmark'} />
                                </button>
                                <button
                                    onClick={handleShare}
                                    className={styles['action-btn']}
                                    aria-label={t('post_detail.share')}
                                >
                                    <XIcon name="carbon:share" />
                                </button>
                            </div>
                        </footer>
                    </div>
                </div>
            </article>

            {/* Comments Section */}
            <CommentSection
                objectId={postId}
                contentType="CONTENT_TYPE_POST"
                onUpdateComments={() => {
                }}
            />

            {/* Related Posts */}
            <section className={styles['related-section']}>
                <div className={styles['section-header']}>
                    <h2>
                        <XIcon name="carbon:book" />
                        <span>{t('post_detail.related_posts')}</span>
                    </h2>
                </div>
                {relatedPostsQuery && (
                    <PostList
                        ref={relatedPostListRef}
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
            {showBackToTop && (
                <button onClick={scrollToTop} className={styles['back-to-top']}>
                    <XIcon name="carbon:arrow-up" />
                </button>
            )}
        </div>
    );
}
