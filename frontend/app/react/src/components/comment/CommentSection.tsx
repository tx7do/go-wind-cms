'use client';

import React, {useState, useEffect, useRef} from 'react';
import {Button} from '@/components/ui/button';
import {Spinner} from '@/components/ui/spinner';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';
import {AppEmpty} from '@/components/ui';
import {fetchListComments} from '@/api/hooks/comment';
import {createComment as createCommentApi} from '@/api/service/comment';
import type {
    commentservicev1_Comment,
    commentservicev1_Comment_ContentType, commentservicev1_ListCommentResponse,
} from '@/api/generated/app/service/v1';

import styles from './CommentSection.module.css';
import CommentTree from './CommentTree';

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
        if (!newComment.content.trim()) {
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
        <section className={styles.commentsSection}>
            <div className={styles.sectionHeader}>
                <h2>
                    <XIcon name="carbon:chat" size={36}/>
                    {t('comments_count', {count: displayComments.length})}
                </h2>
            </div>

            {/* Comment Form */}
            <div className={styles.commentForm}>
                <div className={styles.formHeader}>
                    <div className={styles.formIcon}>
                        <XIcon name="carbon:edit" size={24}/>
                    </div>
                    <h3>{t('write_comment')}</h3>
                </div>
                <div className={styles.commentFormWrapper}>
                    <div className={styles.formRow}>
                        <div className={styles.formGroup}>
                            <input
                                value={newComment.authorName}
                                onChange={(e) => setNewComment({...newComment, authorName: e.target.value})}
                                placeholder={t('nickname') + ' *'}
                                className={styles.inputField}
                                disabled={submitting}
                            />
                        </div>
                        <div className={styles.formGroup}>
                            <input
                                value={newComment.authorEmail}
                                onChange={(e) => setNewComment({...newComment, authorEmail: e.target.value})}
                                placeholder={t('email') + ' *'}
                                className={styles.inputField}
                                disabled={submitting}
                                type="email"
                            />
                        </div>
                    </div>
                    <div className={styles.formGroup}>
                        <textarea
                            value={newComment.content}
                            onChange={(e) => setNewComment({...newComment, content: e.target.value})}
                            placeholder={t('write_comment')}
                            className={styles.textareaField}
                            disabled={submitting}
                            rows={6}
                            onKeyDown={(e) => {
                                // Ctrl/Cmd + Enter 提交
                                if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
                                    e.preventDefault();
                                    handleSubmitComment();
                                }
                            }}
                        />
                        <div className={styles.charCount}>
                            {newComment.content.length} / 1000
                        </div>
                    </div>
                    <div className={styles.formActions}>
                        <Button
                            size="lg"
                            onClick={handleSubmitComment}
                            disabled={submitting}
                            className="gap-2"
                        >
                            {submitting ? <Spinner size="sm"/> : <XIcon name="carbon:send-alt" size={18}/>}
                            {t('submit_comment')}
                        </Button>
                        <span className={styles.formTip}>
                            <XIcon name="carbon:information" size={16}/>
                            {t('fill_form_info')}
                        </span>
                        <span className={styles.formShortcut}>
                            <XIcon name="carbon:keyboard" size={16}/>
                            Ctrl + Enter {t('submit_comment')}
                        </span>
                    </div>
                </div>
            </div>

            {/* Comments List */}
            {hasComments ? (
                <div ref={commentsListRef} className={styles.commentsList}>
                    <CommentTree
                        comments={displayComments}
                        onReply={handleReply}
                        onLoadChildren={loadChildren}
                    />

                    {/* 加载更多按钮 */}
                    {showLoadMore && (
                        <div className={styles.loadMoreContainer}>
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
                        <div className={styles.noMoreText}>
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
