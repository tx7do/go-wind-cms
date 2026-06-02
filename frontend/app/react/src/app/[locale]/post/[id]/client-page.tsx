'use client';

import {useEffect, useState, useMemo, useRef} from 'react';
import {useParams, useSearchParams} from 'next/navigation';
import {useTranslations} from 'next-intl';

import CommentSection from '@/components/comment/CommentSection';
import ContentViewer from '@/components/content/ContentViewer';
import PostList from '@/components/post/PostList';
import PostMetaBar from '@/components/post/PostMetaBar';
import PostFloatingActions from '@/components/post/PostFloatingActions';
import TableOfContents from '@/components/post/TableOfContents';
import BackToTop from '@/components/layout/BackToTop';
import BackButton from '@/components/layout/BackButton';
import SectionContainer from '@/components/layout/SectionContainer';

import {
    fetchPost,
    getPostTitle,
    getPostContent,
    getPostThumbnail,
} from '@/api/hooks/post';
import {contentservicev1_Post} from "@/api/generated/app/service/v1";
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from "@/i18n/helpers";

export default function PostDetailPage() {
    const t = useTranslations('page');
    const params = useParams();
    const router = useI18nRouter();
    const searchParams = useSearchParams();

    const [post, setPost] = useState<contentservicev1_Post | null>(null);
    const [postLoading, setPostLoading] = useState(false);
    const [localLoading, setLocalLoading] = useState(true);
    const isLoading = localLoading || (postLoading && !post);
    const [isLiked, setIsLiked] = useState(false);
    const [isBookmarked, setIsBookmarked] = useState(false);
    const [isTocExpanded, setIsTocExpanded] = useState(true);

    const contentRef = useRef<HTMLDivElement>(null);

    const postId = useMemo(() => {
        const id = params?.id as string;
        return id ? parseInt(id) : null;
    }, [params?.id]);

    const displayTitle = useMemo(() => post ? getPostTitle(post) : '', [post]);
    const displayContent = useMemo(() => post ? getPostContent(post) : '', [post]);
    const displayThumbnail = useMemo(() => post ? getPostThumbnail(post) : '', [post]);

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
                const fetchedPost = (await fetchPost(postId!)) as contentservicev1_Post;
                setPost(fetchedPost);
                if (fetchedPost) {
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
    }, [postId]);

    // Handlers
    const handleBack = () => {
        const from = searchParams.get('from');
        const categoryId = searchParams.get('categoryId');

        if (from === 'category' && categoryId) {
            router.push(`/category/${categoryId}`);
            return;
        }
        if (from === 'tag') { router.push('/tag'); return; }
        if (from === 'user') { router.push('/user'); return; }
        if (from === 'home') { router.push('/'); return; }
        if (from === 'post-list') { router.push('/post'); return; }

        if (typeof window !== 'undefined' && window.history.length > 2) {
            router.back();
            return;
        }
        router.push('/post');
    };

    const handleShare = () => {
        const url = window.location.href;
        if (navigator.share) {
            navigator.share({title: displayTitle, url}).catch(() => {
                navigator.clipboard.writeText(url);
            });
        } else {
            navigator.clipboard.writeText(url);
        }
    };

    if (isLoading) {
        return (
            <div className="w-full">
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
            <SectionContainer noPadding className="!py-6">
                <BackButton label={t('post_detail.back')} onClick={handleBack}/>
            </SectionContainer>

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
                    {/* Table of Contents */}
                    <TableOfContents
                        contentRef={contentRef}
                        contentKey={displayContent}
                        isExpanded={isTocExpanded}
                        onCollapse={() => setIsTocExpanded(false)}
                        onExpand={() => setIsTocExpanded(true)}
                        title={t('post_detail.table_of_contents')}
                    />

                    <div className="flex-1 min-w-0">
                        {/* Post Header */}
                        <header className="mb-8">
                            <h1 className="mb-5 text-3xl font-bold leading-tight text-foreground max-md:text-2xl">
                                {displayTitle}
                            </h1>
                            <PostMetaBar
                                authorName={post.authorName}
                                createdAt={post.createdAt}
                                visits={post.visits}
                                likes={post.likes}
                            />
                        </header>

                        {/* Post Content */}
                        <div className="mb-8" ref={contentRef}>
                            <ContentViewer content={displayContent} type="markdown"/>
                        </div>

                        {/* Floating Actions */}
                        <PostFloatingActions
                            isLiked={isLiked}
                            isBookmarked={isBookmarked}
                            onLike={() => setIsLiked(!isLiked)}
                            onBookmark={() => setIsBookmarked(!isBookmarked)}
                            onShare={handleShare}
                            labels={{
                                likes: t('post_detail.likes'),
                                bookmark: t('post_detail.bookmark'),
                                share: t('post_detail.share'),
                            }}
                        />
                    </div>
                </div>
            </article>

            {/* Comments Section */}
            <SectionContainer noPadding>
                <CommentSection
                    objectId={postId}
                    contentType="CONTENT_TYPE_POST"
                    onUpdateComments={() => {}}
                />
            </SectionContainer>

            {/* Related Posts */}
            <SectionContainer>
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
            </SectionContainer>

            <BackToTop scrollThreshold={500}/>
        </div>
    );
}
