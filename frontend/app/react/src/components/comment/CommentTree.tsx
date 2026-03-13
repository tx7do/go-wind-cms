'use client';

import React, {useState} from 'react';
import {Button, Input, Avatar, Spin} from 'antd';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import type {commentservicev1_Comment} from "@/api/generated/app/service/v1";
import {formatDate} from "@/utils/date";

import styles from './CommentTree.module.css';

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
    const [likedComments, setLikedComments] = useState<Set<number>>(new Set());

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
            await onReply(comment, replyContent.trim());
            cancelReply();
        } catch (error) {
            console.error('Submit reply failed:', error);
            alert(t('submit_comment_failed'));
        } finally {
            setSubmitting(false);
        }
    }

    // 处理点赞
    function handleLike(comment: commentservicev1_Comment) {
        const commentId = comment.id || 0;
        
        if (likedComments.has(commentId)) {
            // 取消点赞
            setLikedComments(prev => {
                const next = new Set(prev);
                next.delete(commentId);
                return next;
            });
        } else {
            // 点赞
            setLikedComments(prev => new Set(prev).add(commentId));
            // 这里可以调用 API 更新点赞数
            // comment.likeCount = (comment.likeCount || 0) + 1;
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
        <div className={styles.commentTree}>
            {comments.map((comment) => (
                <div key={comment.id} className={styles.commentNode}>
                    {/* 评论主体 */}
                    <div className={styles.commentItem}>
                        <div className={styles.commentAvatar}>
                            <Avatar size={48}>
                                {comment.authorName?.charAt(0) || 'U'}
                            </Avatar>
                        </div>
                        <div className={styles.commentBody}>
                            <div className={styles.commentHeader}>
                                <div className={styles.authorInfo}>
                                    <strong className={styles.commentAuthor}>
                                        {isOwnerReply(comment) && (
                                            <span className={styles.ownerBadge}>
                                                <XIcon name="carbon:badge" size={13}/>
                                                {t('owner_reply')}
                                            </span>
                                        )}
                                        {comment.authorName}
                                    </strong>
                                    {comment.location && !isOwnerReply(comment) && (
                                        <span className={styles.commentLocation}>{comment.location}</span>
                                    )}
                                </div>
                                <span className={styles.commentDate}>{formatDate(comment.createdAt)}</span>
                            </div>
                            <div className={styles.commentContent}>
                                {comment.content}
                            </div>
                            <div className={styles.commentActions}>
                                <span 
                                    className={`${styles.actionItem} ${likedComments.has(comment.id || 0) ? styles.liked : ''}`}
                                    onClick={() => handleLike(comment)}
                                    title={likedComments.has(comment.id || 0) ? t('unlike') : t('like')}
                                >
                                    <XIcon name="carbon:thumbs-up" size={16}/>
                                    {comment.likeCount || 0}
                                </span>
                                <span
                                    className={styles.actionItem}
                                    onClick={() => handleReply(comment)}
                                    title={t('reply')}
                                >
                                    <XIcon name="carbon:chat" size={16}/>
                                    {t('reply')}
                                </span>
                                <span
                                    className={styles.actionItem}
                                    onClick={() => handleShare(comment)}
                                    title={t('share')}
                                >
                                    <XIcon name="carbon:share" size={16}/>
                                    {t('share')}
                                </span>
                                {/* 查看回复按钮 */}
                                {comment.replyCount && comment.replyCount > 0 && (
                                    <span
                                        className={`${styles.actionItem} ${styles.viewReplies}`}
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
                                )}
                            </div>

                            {/* 回复表单 */}
                            {replyingCommentId === comment.id && (
                                <div className={styles.replyForm}>
                                    <Input.TextArea
                                        value={replyContent}
                                        onChange={(e) => setReplyContent(e.target.value)}
                                        rows={3}
                                        placeholder={t('write_comment')}
                                        size="small"
                                        disabled={submitting}
                                        onPressEnter={(e) => {
                                            if (e.ctrlKey) {
                                                e.preventDefault();
                                                submitReply(comment);
                                            }
                                        }}
                                    />
                                    <div className={styles.replyFormActions}>
                                        <Button
                                            size="small"
                                            type="primary"
                                            onClick={() => submitReply(comment)}
                                            loading={submitting}
                                        >
                                            {t('submit_comment')}
                                        </Button>
                                        <Button
                                            size="small"
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
                        <div className={styles.commentChildren}>
                            {/* 加载中提示 */}
                            {isLoading(comment) && (
                                <div className={styles.loadingChildren}>
                                    <Spin size="small"/>
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
