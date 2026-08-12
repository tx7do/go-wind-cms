'use client';

import React, {useState} from 'react';
import {Button} from '@/components/ui/button';
import {Avatar, AvatarFallback} from '@/components/ui/avatar';
import {Spinner} from '@/components/ui/spinner';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import type {commentservicev1_Comment} from "@/api/generated/app/service/v1";
import {formatDate} from "@/utils/date";

import {cn} from '@/lib/utils';
import {useInteractionStatus, useLike, useUnlike} from '@/api/hooks/interaction';

interface CommentTreeProps {
    comments: commentservicev1_Comment[];
    onReply: (comment: commentservicev1_Comment, content: string) => void;
    onLoadChildren: (comment: commentservicev1_Comment) => Promise<void>;
}

const CommentTree: React.FC<CommentTreeProps> = ({
                                                     comments,
                                                     onReply,
                                                     onLoadChildren
                                                 }) => {
    const t = useTranslations('comment');
    const [replyingCommentId, setReplyingCommentId] = useState<number | null>(null);
    const [replyContent, setReplyContent] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [expandedComments, setExpandedComments] = useState<Set<number>>(new Set());
    const [loadingChildren, setLoadingChildren] = useState<Set<number>>(new Set());

    // 点赞状态：通过 interaction 子系统批量查询当前 viewer 对这些评论的点赞态。
    // targetIds 取当前渲染的顶层评论 id（children 由递归子组件自行查询）。
    const commentIds = comments
        .map(c => c.id)
        .filter((id): id is number => typeof id === 'number');
    const statusQuery = useInteractionStatus('TARGET_TYPE_COMMENT', commentIds);
    const likeMutation = useLike();
    const unlikeMutation = useUnlike();

    function isLiked(commentId?: number): boolean {
        if (typeof commentId !== 'number') return false;
        const map = statusQuery.data?.statuses;
        if (!map) return false;
        return map[String(commentId)]?.liked === true;
    }

    function hasChildren(comment: commentservicev1_Comment): boolean {
        return !!comment.children || (comment.replyCount !== undefined && comment.replyCount > 0);
    }

    function isOwnerReply(comment: commentservicev1_Comment): boolean {
        return comment.authorType === 'AUTHOR_TYPE_ADMIN' && !!comment.replyToId;
    }

    function isExpanded(comment: commentservicev1_Comment): boolean {
        return expandedComments.has(comment.id || 0);
    }

    function isLoading(comment: commentservicev1_Comment): boolean {
        return loadingChildren.has(comment.id || 0);
    }

    function handleReply(comment: commentservicev1_Comment) {
        setReplyingCommentId(comment.id || 0);
        setReplyContent('');
    }

    function cancelReply() {
        setReplyingCommentId(null);
        setReplyContent('');
    }

    async function submitReply(comment: commentservicev1_Comment) {
        if (!replyContent.trim()) {
            alert(t('empty_content'));
            return;
        }

        if (submitting) return;

        setSubmitting(true);
        try {
            onReply(comment, replyContent.trim());
            cancelReply();
        } catch (error) {
            console.error('Submit reply failed:', error);
            alert(t('submit_comment_failed'));
        } finally {
            setSubmitting(false);
        }
    }

    // 处理点赞：调用 interaction 子系统，服务端在同一事务内更新 ledger 与计数缓存。
    // 幂等由后端保证；前端按当前态决定 Like/Unlike。计数直读 comment.likeCount（缓存）。
    async function handleLike(comment: commentservicev1_Comment) {
        const commentId = comment.id;
        if (typeof commentId !== 'number') return;
        const liked = isLiked(commentId);
        try {
            if (liked) {
                await unlikeMutation.mutateAsync({targetType: 'TARGET_TYPE_COMMENT', targetId: commentId});
            } else {
                await likeMutation.mutateAsync({targetType: 'TARGET_TYPE_COMMENT', targetId: commentId});
            }
        } catch (e) {
            console.error('toggle like failed:', e);
        }
    }

    function handleShare(comment: commentservicev1_Comment) {
        // 获取当前 URL（不包含 hash）
        const baseUrl = typeof window !== 'undefined' ? window.location.href.split('#')[0] : '';
        const commentUrl = `${baseUrl}#comment-${comment.id}`;

        navigator.clipboard.writeText(commentUrl).then(() => {
            // message.success(t('link_copied'));
        }).catch(() => {
            // message.error(t('copy_failed'));
        });
    }

    // 切换展开/收起子评论
    async function toggleExpand(comment: commentservicev1_Comment) {
        if (isExpanded(comment)) {
            // 收起
            setExpandedComments(prev => {
                const next = new Set(prev);
                next.delete(comment.id || 0);
                return next;
            });
        } else {
            // 展开
            if (!comment.children && comment.replyCount && comment.replyCount > 0) {
                // 需要动态加载子评论
                setLoadingChildren(prev => new Set(prev).add(comment.id || 0));
                try {
                    await onLoadChildren(comment);
                    setExpandedComments(prev => new Set(prev).add(comment.id || 0));
                } catch (error) {
                    console.error('Load children failed:', error);
                } finally {
                    setLoadingChildren(prev => {
                        const next = new Set(prev);
                        next.delete(comment.id || 0);
                        return next;
                    });
                }
            } else {
                // 已经有子评论数据，直接展开
                setExpandedComments(prev => new Set(prev).add(comment.id || 0));
            }
        }
    }

    if (!comments || comments.length === 0) return null;

    return (
        <div className="flex flex-col gap-7">
            {comments.map((comment) => (
                <div key={comment.id} className="flex flex-col">
                    {/* 评论主体 */}
                    <div className={cn(
                        'group relative flex gap-5 overflow-hidden rounded-2xl border border-border bg-linear-to-br from-card to-primary/2 p-7',
                        'shadow-sm transition-all duration-400',
                        'hover:border-primary hover:shadow-md hover:translate-x-1 hover:-translate-y-0.5',
                        'max-md:p-6 max-md:gap-4',
                        'max-sm:p-4 max-sm:gap-3',
                    )}>
                        {/* Left accent bar on hover */}
                        <div className="absolute top-0 start-0 h-full w-1 bg-primary opacity-0 transition-opacity duration-300 group-hover:opacity-100"/>

                        <div className="shrink-0">
                            <Avatar className="h-12 w-12 ring-2 ring-border/30 shadow-md">
                                <AvatarFallback className="bg-primary text-primary-foreground font-semibold">
                                    {comment.authorName?.charAt(0) || 'U'}
                                </AvatarFallback>
                            </Avatar>
                        </div>
                        <div className="min-w-0 flex-1">
                            {/* Header */}
                            <div className={cn(
                                'mb-4 flex items-center justify-between border-b border-primary/8 pb-3',
                                'max-md:flex-col max-md:items-start max-md:gap-2 max-md:mb-3 max-md:pb-2.5',
                                'max-sm:gap-1.5 max-sm:mb-2 max-sm:pb-1.5',
                            )}>
                                <div className="flex flex-wrap items-center gap-3 max-md:flex-col max-md:items-start max-md:gap-1.5">
                                    <strong className="flex items-center gap-2 text-[17px] font-semibold tracking-tight text-foreground max-md:text-base max-sm:text-[15px]">
                                        <span className="text-sm font-bold text-primary">@</span>
                                        {isOwnerReply(comment) && (
                                            <span className="inline-flex items-center gap-1 rounded-lg bg-primary px-2.5 py-0.5 text-[11px] font-bold text-primary-foreground shadow-sm max-sm:text-[10px] max-sm:px-1.5 max-sm:py-0.5">
                                                <XIcon name="carbon:badge" size={13}/>
                                                {t('owner_reply')}
                                            </span>
                                        )}
                                        {comment.authorName}
                                    </strong>
                                    {comment.location && !isOwnerReply(comment) && (
                                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                                            {comment.location}
                                        </span>
                                    )}
                                </div>
                                <span className="flex items-center gap-1.5 text-[13px] text-muted-foreground max-md:text-xs max-sm:text-[11px]">
                                    {formatDate(comment.createdAt)}
                                </span>
                            </div>

                            {/* Content */}
                            <div className="mb-4 ps-1 text-[15px] leading-relaxed text-foreground max-md:text-sm max-sm:text-[13px]">
                                {comment.content}
                            </div>

                            {/* Actions */}
                            <div className={cn(
                                'flex gap-5 border-t border-primary/5 pt-3',
                                'max-md:flex-wrap max-md:gap-4',
                                'max-sm:gap-2.5',
                            )}>
                                <span
                                    className={cn(
                                        'flex cursor-pointer items-center gap-1.5 text-[13px] text-muted-foreground transition-all duration-200 select-none hover:text-primary hover:-translate-y-0.5',
                                        'max-md:text-xs max-sm:text-[11px]',
                                        isLiked(comment.id) && 'text-primary font-semibold',
                                    )}
                                    onClick={() => handleLike(comment)}
                                    title={isLiked(comment.id) ? t('unlike') : t('like')}
                                >
                                    <XIcon name="carbon:thumbs-up" size={16}/>
                                    {comment.likeCount || 0}
                                </span>
                                <span
                                    className="flex cursor-pointer items-center gap-1.5 text-[13px] text-muted-foreground transition-all duration-200 select-none hover:text-primary hover:-translate-y-0.5 max-md:text-xs max-sm:text-[11px]"
                                    onClick={() => handleReply(comment)}
                                    title={t('reply')}
                                >
                                    <XIcon name="carbon:chat" size={16}/>
                                    {t('reply')}
                                </span>
                                <span
                                    className="flex cursor-pointer items-center gap-1.5 text-[13px] text-muted-foreground transition-all duration-200 select-none hover:text-primary hover:-translate-y-0.5 max-md:text-xs max-sm:text-[11px]"
                                    onClick={() => handleShare(comment)}
                                    title={t('share')}
                                >
                                    <XIcon name="carbon:share" size={16}/>
                                    {t('share')}
                                </span>
                                {/* 查看回复按钮 */}
                                {comment.replyCount && comment.replyCount > 0 ? (
                                    <span
                                        className="flex cursor-pointer items-center gap-1.5 text-[13px] font-medium text-primary transition-all duration-200 select-none hover:text-primary max-md:text-xs max-sm:text-[11px]"
                                        onClick={() => toggleExpand(comment)}
                                        title={isExpanded(comment) ? t('hide_replies') : t('view_replies', {count: comment.replyCount})}
                                    >
                                        <XIcon
                                            name={isExpanded(comment) ? 'carbon:chevron-up' : 'carbon:chevron-down'}
                                            size={18}
                                        />
                                        {isExpanded(comment)
                                            ? t('hide_replies')
                                            : t('view_replies', {count: comment.replyCount})
                                        }
                                    </span>
                                ) : null}
                            </div>

                            {/* 回复表单 */}
                            {replyingCommentId === comment.id && (
                                <div className="mt-4 border-t border-primary/8 pt-4 max-md:mt-3 max-md:pt-3">
                                    <textarea
                                        value={replyContent}
                                        onChange={(e) => setReplyContent(e.target.value)}
                                        rows={3}
                                        placeholder={t('write_comment')}
                                        className={cn(
                                            'w-full min-h-20 resize-y rounded-lg border border-border bg-muted px-4 py-3 text-sm text-foreground',
                                            'transition-all duration-300',
                                            'hover:border-primary',
                                            'focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15',
                                            'disabled:cursor-not-allowed disabled:opacity-60 disabled:bg-background',
                                        )}
                                        disabled={submitting}
                                        onKeyDown={(e) => {
                                            if (e.ctrlKey || e.metaKey) {
                                                e.preventDefault();
                                                submitReply(comment);
                                            }
                                        }}
                                    />
                                    <div className="mt-2 mb-3 text-end text-xs text-muted-foreground">
                                        {replyContent.length} / 1000
                                    </div>
                                    <div className="flex justify-end gap-3 max-md:flex-col max-md:gap-2">
                                        <Button
                                            size="sm"
                                            onClick={() => submitReply(comment)}
                                            disabled={submitting}
                                        >
                                            {submitting ? '...' : t('submit_comment')}
                                        </Button>
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            onClick={cancelReply}
                                            disabled={submitting}
                                        >
                                            {t('cancel')}
                                        </Button>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* 递归渲染子评论 */}
                    {hasChildren(comment) && isExpanded(comment) && (
                        <div className={cn(
                            'relative mt-5 ps-16',
                            'max-md:ps-12 max-md:mt-4',
                            'max-sm:ps-9 max-sm:mt-3',
                        )}>
                            {/* Vertical line */}
                            <div className="absolute top-0 start-8 h-full w-0.5 bg-linear-to-b from-primary/20 to-primary/5 max-md:start-6 max-sm:start-4.5 max-sm:w-[1.5px]"/>

                            {/* 加载中提示 */}
                            {isLoading(comment) && (
                                <div className="flex items-center justify-center gap-3 py-5 text-sm text-muted-foreground">
                                    <Spinner size="sm" className="text-primary"/>
                                    <span>{t('loading')}</span>
                                </div>
                            )}
                            {/* 子评论列表 */}
                            {!isLoading(comment) && comment.children && (
                                <CommentTree
                                    comments={comment.children}
                                    onReply={onReply}
                                    onLoadChildren={onLoadChildren}
                                />
                            )}
                        </div>
                    )}
                </div>
            ))}
        </div>
    );
};

export default CommentTree;
