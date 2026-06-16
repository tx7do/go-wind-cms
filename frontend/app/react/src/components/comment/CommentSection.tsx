'use client';

import React, {useState, useEffect, useRef} from 'react';
import {Button} from '@/components/ui/button';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {AppEmpty} from '@/components/ui';
import {fetchListComments} from '@/api/hooks/comment';
import {createComment as createCommentApi} from '@/api/hooks/comment';
import type {
    commentservicev1_Comment,
    commentservicev1_Comment_ContentType, commentservicev1_ListCommentResponse,
} from '@/api/generated/app/service/v1';

import {cn} from '@/lib/utils';
import CommentTree from './CommentTree';
import RichTextEditor from './RichTextEditor';

interface CommentForm {
    content: string;
    authorName: string;
    authorEmail: string;
}

interface CommentSectionProps {
    objectId: number | null;
    contentType: commentservicev1_Comment_ContentType | null;
    onUpdateComments?: (comments: commentservicev1_Comment[]) => void;
}

const CommentSection: React.FC<CommentSectionProps> = ({
                                                           objectId,
                                                           contentType,
                                                           onUpdateComments
                                                       }) => {
    const t = useTranslations('comment');

    // 状态
    const [submitting, setSubmitting] = useState(false);
    const [loading, setLoading] = useState(false);
    const [comments, setComments] = useState<commentservicev1_Comment[]>([]);

    // 分页相关状态
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize] = useState(20);
    const [_total, setTotal] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);

    const [newComment, setNewComment] = useState<CommentForm>({
        content: '',
        authorName: '',
        authorEmail: '',
    });

    const commentsListRef = useRef<HTMLDivElement>(null);

    // 计算属性
    const displayComments = comments;
    const hasComments = displayComments.length > 0;
    const showLoadMore = hasMore && !loadingMore;

    // 加载评论 (分页)
    async function loadComments(reset = false) {
        if (!objectId || !contentType) return;
        if (loading || (loadingMore && !reset)) return;

        if (reset) {
            setLoading(true);
        } else {
            setLoadingMore(true);
        }

        try {
            const res = await fetchListComments({
                paging: {
                    page: reset ? 1 : currentPage,
                    pageSize: pageSize
                },
                formValues: {
                    objectId: objectId,
                    contentType: contentType,
                    status: 'STATUS_APPROVED'
                }
            }) as unknown as commentservicev1_ListCommentResponse;

            const newComments = res.items || [];
            const newTotal = res.total || 0;

            if (reset) {
                setComments(newComments);
            } else {
                setComments(prev => [...prev, ...newComments]);
            }

            setTotal(newTotal);
            // 判断是否还有更多数据
            const allComments = reset ? newComments : [...comments, ...newComments];
            setHasMore(allComments.length < newTotal);
            setCurrentPage(prev => prev + 1);

            onUpdateComments?.(allComments);
        } catch (error) {
            console.error('Load comments failed:', error);
        } finally {
            setLoading(false);
            setLoadingMore(false);
        }
    }

    // 加载更多
    function loadMoreComments() {
        if (hasMore && !loadingMore) {
            loadComments();
        }
    }

    // 提交评论
    async function handleSubmitComment() {
        // Tiptap 空内容会返回 <p></p>，需要用 textContent 检查
        const tmpDiv = document.createElement('div');
        tmpDiv.innerHTML = newComment.content;
        const textContent = tmpDiv.textContent?.trim() || '';
        if (!textContent) {
            alert(t('empty_content'));
            return;
        }
        if (!newComment.authorName.trim()) {
            alert(t('invalid_nickname'));
            return;
        }
        if (!newComment.authorEmail.trim()) {
            alert(t('invalid_email'));
            return;
        }
        // 邮箱格式验证
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(newComment.authorEmail)) {
            alert(t('invalid_email'));
            return;
        }
        if (!objectId || submitting) return;

        setSubmitting(true);
        try {
            await createCommentApi({
                objectId: objectId,
                contentType: contentType ?? undefined,
                content: newComment.content,
                authorName: newComment.authorName,
                authorEmail: newComment.authorEmail,
                status: 'STATUS_PENDING',
            });

            alert(t('comment_submitted'));
            setNewComment({content: '', authorName: '', authorEmail: ''});
            setCurrentPage(1);
            await loadComments(true);

            // 用户体验优化：提交后滚动到评论列表
            setTimeout(() => {
                commentsListRef.current?.scrollIntoView({behavior: 'smooth'});
            }, 100);
        } catch (error) {
            console.error('Submit comment failed:', error);
            alert(t('submit_comment_failed'));
        } finally {
            setSubmitting(false);
        }
    }

    // 处理回复
    async function handleReply(comment: commentservicev1_Comment, content: string) {
        if (!content.trim()) {
            alert(t('empty_content'));
            return;
        }

        if (!objectId || submitting) return;

        setSubmitting(true);
        try {
            await createCommentApi({
                objectId: objectId,
                contentType: contentType ?? undefined,
                content: content.trim(),
                authorName: comment.authorName,
                authorEmail: comment.authorEmail,
                parentId: comment.id,
                replyToId: comment.id,
                status: 'STATUS_PENDING',
            });

            alert(t('comment_posted'));
            setCurrentPage(1);
            await loadComments(true);
        } catch (error) {
            console.error('Submit reply failed:', error);
            alert(t('submit_comment_failed'));
        } finally {
            setSubmitting(false);
        }
    }

    // 加载子评论
    async function loadChildren(parentComment: commentservicev1_Comment) {
        if (!objectId || !contentType) return;

        try {
            const res = await fetchListComments({
                paging: {
                    page: 1,
                    pageSize: 50
                },
                formValues: {
                    objectId: objectId,
                    contentType: contentType,
                    parentId: parentComment.id,
                    status: 'STATUS_APPROVED'
                }
            }) as unknown as commentservicev1_ListCommentResponse;

            parentComment.children = res.items || [];
        } catch (error) {
            console.error('Load children failed:', error);
            throw error;
        }
    }

    // 生命周期：自动加载评论
    useEffect(() => {
        if (objectId && contentType) {
            loadComments(true);
        }
    }, [objectId, contentType]);

    return (
        <section className={cn(
            'mx-auto mb-10 max-w-300 rounded-2xl border border-border bg-card p-14 shadow-sm backdrop-blur-sm',
            'max-md:rounded-xl max-md:p-10',
            'max-sm:rounded-xl max-sm:p-8',
        )}>
            {/* Section Header */}
            <div className="mb-10 max-md:mb-8">
                <h2 className="flex items-center gap-3.5 text-3xl font-bold tracking-tight text-foreground max-md:text-2xl max-sm:text-xl">
                    <XIcon name="carbon:chat" size={36}/>
                    {t('comments_count', {count: displayComments.length})}
                </h2>
            </div>

            {/* Comment Form */}
            <div className={cn(
                'relative mb-12 overflow-hidden rounded-2xl border border-primary/10 p-12',
                'bg-linear-to-br from-card to-primary/2 shadow-sm',
                'transition-all duration-400 hover:border-primary hover:shadow-md',
                'max-md:rounded-xl max-md:p-8 max-md:mb-9',
                'max-sm:p-7 max-sm:mb-8',
            )}>
                {/* Gradient top bar */}
                <div className="absolute top-0 left-0 right-0 h-1 bg-primary opacity-90"/>

                {/* 标题块：去除绿色背景方框，改为简洁文字 + 绿色图标点绥 */}
                <div className="mb-8 flex items-center gap-2.5 max-md:mb-6">
                    <XIcon name="carbon:edit" size={22} className="text-primary"/>
                    <h3 className="text-lg font-bold tracking-tight text-foreground max-md:text-base">
                        {t('write_comment')}
                    </h3>
                </div>

                <div className="flex flex-col gap-6">
                    <div className="grid grid-cols-2 gap-6 max-md:grid-cols-1 max-md:gap-4">
                        <div className="flex flex-col gap-2">
                            <input
                                value={newComment.authorName}
                                onChange={(e) => setNewComment({...newComment, authorName: e.target.value})}
                                placeholder={t('nickname') + ' *'}
                                className={cn(
                                    'w-full rounded-lg border border-border bg-background px-4 py-3 text-foreground',
                                    'transition-all duration-300',
                                    'hover:border-primary',
                                    'focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15',
                                    'disabled:cursor-not-allowed disabled:opacity-60',
                                )}
                                disabled={submitting}
                            />
                        </div>
                        <div className="flex flex-col gap-2">
                            <input
                                value={newComment.authorEmail}
                                onChange={(e) => setNewComment({...newComment, authorEmail: e.target.value})}
                                placeholder={t('email') + ' *'}
                                className={cn(
                                    'w-full rounded-lg border border-border bg-background px-4 py-3 text-foreground',
                                    'transition-all duration-300',
                                    'hover:border-primary',
                                    'focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15',
                                    'disabled:cursor-not-allowed disabled:opacity-60',
                                )}
                                disabled={submitting}
                                type="email"
                            />
                        </div>
                    </div>
                    <div className="flex flex-col gap-2">
                        <RichTextEditor
                            value={newComment.content}
                            onChange={(content) => setNewComment({...newComment, content})}
                            onSubmit={handleSubmitComment}
                            submitting={submitting}
                            placeholder={t('write_comment')}
                            maxLength={1000}
                            submitLabel={t('submit_comment')}
                        />
                    </div>
                    <div className="flex flex-wrap items-center gap-4 max-md:flex-col max-md:items-stretch">
                        <span className="flex items-center gap-2 rounded-lg border border-primary/10 bg-primary/5 px-3.5 py-2 text-[13px] text-muted-foreground max-md:justify-center max-md:text-xs max-md:px-3 max-md:py-1.5">
                            <XIcon name="carbon:information" size={16} className="text-primary"/>
                            {t('fill_form_info')}
                        </span>
                    </div>
                </div>
            </div>

            {/* Comments List */}
            {hasComments ? (
                <div ref={commentsListRef} className="flex flex-col gap-7">
                    <CommentTree
                        comments={displayComments}
                        onReply={handleReply}
                        onLoadChildren={loadChildren}
                    />

                    {/* 加载更多按钮 */}
                    {showLoadMore && (
                        <div className="mt-8 flex justify-center border-t border-primary/8 pt-6">
                            <Button
                                size="lg"
                                onClick={loadMoreComments}
                                disabled={loadingMore}
                            >
                                {loadingMore ? 'Loading...' : t('load_more')}
                            </Button>
                        </div>
                    )}

                    {/* 没有更多提示 */}
                    {!hasMore && hasComments && (
                        <div className="py-5 text-center text-sm font-medium tracking-wide text-muted-foreground">
                            {t('no_more')}
                        </div>
                    )}
                </div>
            ) : (
                <AppEmpty
                    description={t('no_comments')}
                    inContainer
                />
            )}
        </section>
    );
};

export default CommentSection;
