import React, {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Textarea} from '@tarojs/components';

import type {commentservicev1_Comment} from "@/api/generated/app/service/v1";
import {formatDate} from "@/utils/date";
import XIcon from '@/plugins/xicon';

import './index.scss';

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
    const {t} = useTranslation();
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
            alert(t('comment.empty_content'));
            return;
        }

        if (submitting) return;

        setSubmitting(true);
        try {
            await onReply(comment, replyContent.trim());
            cancelReply();
        } catch (error) {
            console.error('Submit reply failed:', error);
            alert(t('comment.submit_comment_failed'));
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
        }
    }

    function handleShare(comment: commentservicev1_Comment) {
        // TODO: Taro 中分享功能需要使用特定 API
        console.log('Share comment:', comment.id);
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
        <View className='comment-tree'>
            {comments.map((comment) => (
                <View key={comment.id} className='comment-node'>
                    {/* 评论主体 */}
                    <View className='comment-item'>
                        <View className='comment-avatar'>
                            <View className='avatar-placeholder'>
                                {comment.authorName?.charAt(0) || 'U'}
                            </View>
                        </View>
                        <View className='comment-body'>
                    <View className='comment-header'>
                        <View className='author-info'>
                            <Text className='comment-author'>
                                <Text className='author-at'>@</Text>
                                {isOwnerReply(comment) && (
                                    <Text className='owner-badge'>
                                        <XIcon name='carbon:user' size={12} />
                                        <Text> {t('comment.owner_reply')}</Text>
                                    </Text>
                                )}
                                {comment.authorName}
                            </Text>
                            {comment.location && !isOwnerReply(comment) && (
                                <Text className='comment-location'>
                                    <Text className='location-dot'>•</Text>
                                    {comment.location}
                                </Text>
                            )}
                        </View>
                        <Text className='comment-date'>{formatDate(comment.createdAt)}</Text>
                    </View>
                            <View className='comment-content'>
                                {comment.content}
                            </View>
                            <View className='comment-actions'>
                                <View
                                    className={`action-item ${likedComments.has(comment.id || 0) ? 'liked' : ''}`}
                                    onClick={() => handleLike(comment)}
                                >
                                    <XIcon name={likedComments.has(comment.id || 0) ? 'carbon:thumb-up' : 'carbon:thumb-up-outline'} size={16} />
                                    <Text>{comment.likeCount || 0}</Text>
                                </View>
                                <View
                                    className='action-item'
                                    onClick={() => handleReply(comment)}
                                >
                                    <XIcon name='carbon:chat' size={16} />
                                    <Text>{t('comment.reply')}</Text>
                                </View>
                                <View
                                    className='action-item'
                                    onClick={() => handleShare(comment)}
                                >
                                    <XIcon name='carbon:share' size={16} />
                                    <Text>{t('comment.share')}</Text>
                                </View>

                                {/* 查看回复按钮 */}
                                {comment.replyCount && comment.replyCount > 0 ? (
                                    <View
                                        className={`action-item view-replies ${isExpanded(comment) ? 'expanded' : ''}`}
                                        onClick={() => toggleExpand(comment)}
                                    >
                                        <XIcon name={isExpanded(comment) ? 'carbon:chevron-up' : 'carbon:chevron-down'} size={16} />
                                        <Text> {isExpanded(comment)
                                            ? t('comment.hide_replies')
                                            : t('comment.view_replies', {count: comment.replyCount})
                                        }
                                        </Text>
                                    </View>
                                ) : null}
                            </View>

                            {/* 回复表单 */}
                            {replyingCommentId === comment.id && (
                                <View className='reply-form'>
                                    <Textarea
                                        value={replyContent}
                                        onInput={(e) => setReplyContent(e.detail.value)}
                                        placeholder={t('comment.write_comment')}
                                        className='reply-textarea'
                                        disabled={submitting}
                                        maxlength={1000}
                                    />
                                    <View className='char-count'>
                                        {replyContent.length} / 1000
                                    </View>
                                    <View className='reply-form-actions'>
                                        <View
                                            className='submit-btn-small'
                                            onClick={() => submitReply(comment)}
                                        >
                                            <Text>{submitting ? t('comment.loading') : t('comment.submit_comment')}</Text>
                                        </View>
                                        <View
                                            className='cancel-btn-small'
                                            onClick={cancelReply}
                                        >
                                            <Text>{t('comment.cancel')}</Text>
                                        </View>
                                    </View>
                                </View>
                            )}
                        </View>
                    </View>

                    {/* 递归渲染子评论 */}
                    {hasChildren(comment) && isExpanded(comment) && (
                        <View className='comment-children'>
                            {/* 加载中提示 */}
                            {isLoading(comment) && (
                                <View className='loading-children'>
                                    <XIcon name='carbon:time' size={16} className='loading-icon' />
                                    <Text> {t('comment.loading')}</Text>
                                </View>
                            )}
                            {/* 子评论列表 */}
                            {!isLoading(comment) && comment.children && (
                                <CommentTree
                                    comments={comment.children}
                                    onReply={onReply}
                                    onLoadChildren={onLoadChildren}
                                />
                            )}
                        </View>
                    )}
                </View>
            ))}
        </View>
    );
};

export default CommentTree;
