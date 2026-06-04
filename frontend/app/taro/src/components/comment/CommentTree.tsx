import {View, Text, Textarea} from '@tarojs/components';
import React, {useState} from 'react';
import Taro from '@tarojs/taro';
import {useTranslations} from '@/lib/next-intl-compat';
import {Spinner} from '@/components/ui/spinner';
import {XIcon} from '@/plugins/xicon';
import type {commentservicev1_Comment} from '@/api/generated/app/service/v1';
import {formatDate} from '@/utils/date';
import {cn} from '@/lib/utils';

interface CommentTreeProps {
    comments: commentservicev1_Comment[];
    onReply: (comment: commentservicev1_Comment, content: string) => void;
    onLoadChildren: (comment: commentservicev1_Comment) => Promise<void>;
}

const CommentTree: React.FC<CommentTreeProps> = ({
    comments,
    onReply,
    onLoadChildren,
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
            Taro.showToast({title: t('empty_content'), icon: 'none'});
            return;
        }
        if (submitting) return;

        setSubmitting(true);
        try {
            onReply(comment, replyContent.trim());
            cancelReply();
        } catch (error) {
            console.error('Submit reply failed:', error);
            Taro.showToast({title: t('submit_comment_failed'), icon: 'none'});
        } finally {
            setSubmitting(false);
        }
    }

    // 点赞
    function handleLike(comment: commentservicev1_Comment) {
        const commentId = comment.id || 0;
        if (likedComments.has(commentId)) {
            setLikedComments(prev => {
                const next = new Set(prev);
                next.delete(commentId);
                return next;
            });
        } else {
            setLikedComments(prev => new Set(prev).add(commentId));
        }
    }

    // 分享 - 复制链接
    function handleShare(comment: commentservicev1_Comment) {
        const commentUrl = `#comment-${comment.id}`;
        Taro.setClipboardData({
            data: commentUrl,
            success: () => Taro.showToast({title: t('link_copied'), icon: 'success'}),
            fail: () => Taro.showToast({title: t('copy_failed'), icon: 'none'}),
        });
    }

    // 展开/收起子评论
    async function toggleExpand(comment: commentservicev1_Comment) {
        const id = comment.id || 0;
        if (isExpanded(comment)) {
            setExpandedComments(prev => {
                const next = new Set(prev);
                next.delete(id);
                return next;
            });
        } else {
            if (!comment.children && comment.replyCount && comment.replyCount > 0) {
                setLoadingChildren(prev => new Set(prev).add(id));
                try {
                    await onLoadChildren(comment);
                    setExpandedComments(prev => new Set(prev).add(id));
                } catch (error) {
                    console.error('Load children failed:', error);
                } finally {
                    setLoadingChildren(prev => {
                        const next = new Set(prev);
                        next.delete(id);
                        return next;
                    });
                }
            } else {
                setExpandedComments(prev => new Set(prev).add(id));
            }
        }
    }

    if (!comments || comments.length === 0) return null;

    return (
        <View className='flex flex-col gap-[16rpx]'>
            {comments.map((comment) => (
                <View key={comment.id} className='flex flex-col'>
                    {/* 评论主体 - 白色卡片 */}
                    <View className='bg-cardBg p-[24rpx]'>
                        {/* 头部：头像 + 用户名 + 时间 */}
                        <View className='flex items-start gap-[12rpx] mb-[12rpx]'>
                            <View className='flex items-center justify-center w-[64rpx] h-[64rpx] rounded-full bg-primary/10 flex-shrink-0'>
                                <Text className='text-body font-bold text-primary'>
                                    {comment.authorName?.charAt(0) || 'U'}
                                </Text>
                            </View>
                            <View className='flex-1 min-w-0'>
                                <View className='flex items-center gap-[8rpx] mb-[4rpx]'>
                                    <Text className='text-desc font-bold text-textMain'>
                                        @{comment.authorName}
                                    </Text>
                                    {isOwnerReply(comment) && (
                                        <View className='flex items-center gap-[4rpx] rounded-[8rpx] bg-primary px-[8rpx] py-[2rpx]'>
                                            <XIcon name='carbon:badge' size={10} className='text-white' />
                                            <Text className='text-tips font-bold text-white'>{t('owner_reply')}</Text>
                                        </View>
                                    )}
                                </View>
                                <View className='flex items-center gap-[8rpx]'>
                                    {comment.location && !isOwnerReply(comment) && (
                                        <Text className='text-tips text-textSec'>{comment.location}</Text>
                                    )}
                                    <Text className='text-tips text-textSec'>{formatDate(comment.createdAt)}</Text>
                                </View>
                            </View>
                        </View>

                        {/* 评论内容 */}
                        <View className='mb-[16rpx]'>
                            <Text className='text-body text-textMain leading-[1.6]'>{comment.content}</Text>
                        </View>

                        {/* 操作栏 - 无分割线 */}
                        <View className='flex items-center gap-[24rpx]'>
                            {/* 点赞 */}
                            <View
                              className='flex items-center gap-[4rpx]'
                              onClick={() => handleLike(comment)}
                              hoverClass='tap-active'
                            >
                                <XIcon
                                  name={likedComments.has(comment.id || 0) ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'}
                                  size={16}
                                  className={likedComments.has(comment.id || 0) ? 'text-primary' : 'text-textSec'}
                                />
                                <Text className='text-tips text-textSec'>
                                    {comment.likeCount || 0}
                                </Text>
                            </View>

                            {/* 回复 */}
                            <View
                              className='flex items-center gap-[4rpx]'
                              onClick={() => handleReply(comment)}
                              hoverClass='tap-active'
                            >
                                <XIcon name='carbon:chat' size={16} className='text-textSec' />
                                <Text className='text-tips text-textSec'>{t('reply')}</Text>
                            </View>

                            {/* 分享 */}
                            <View
                              className='flex items-center gap-[4rpx]'
                              onClick={() => handleShare(comment)}
                              hoverClass='tap-active'
                            >
                                <XIcon name='carbon:share' size={16} className='text-textSec' />
                                <Text className='text-tips text-textSec'>{t('share')}</Text>
                            </View>

                            {/* 查看回复 */}
                            {comment.replyCount && comment.replyCount > 0 ? (
                                <View
                                  className='flex items-center gap-[4rpx] ml-auto'
                                  onClick={() => toggleExpand(comment)}
                                  hoverClass='tap-active'
                                >
                                    <XIcon
                                      name={isExpanded(comment) ? 'carbon:chevron-down' : 'carbon:chevron-up'}
                                      size={14}
                                      className='text-primary'
                                    />
                                    <Text className='text-tips text-primary font-medium'>
                                        {isExpanded(comment)
                                            ? t('hide_replies')
                                            : t('view_replies', {count: comment.replyCount})
                                        }
                                    </Text>
                                </View>
                            ) : null}
                        </View>

                        {/* 回复表单 */}
                        {replyingCommentId === comment.id && (
                            <View className='mt-[16rpx] pt-[16rpx]'>
                                <Textarea
                                  value={replyContent}
                                  onInput={(e) => setReplyContent(e.detail.value ?? '')}
                                  placeholder={t('write_comment')}
                                  maxlength={1000}
                                  className='w-full min-h-[120rpx] text-body text-textMain bg-pageBg rounded-[12rpx] p-[16rpx]'
                                  disabled={submitting}
                                />
                                <View className='flex items-center justify-between mt-[12rpx]'>
                                    <Text className='text-tips text-textSec'>
                                        {replyContent.length} / 1000
                                    </Text>
                                    <View className='flex gap-[12rpx]'>
                                        <View
                                          className='px-[24rpx] py-[12rpx] rounded-[8rpx] bg-pageBg'
                                          onClick={cancelReply}
                                          hoverClass='tap-active'
                                        >
                                            <Text className='text-desc text-textSec'>{t('cancel')}</Text>
                                        </View>
                                        <View
                                          className={cn(
                                              'px-[24rpx] py-[12rpx] rounded-[8rpx]',
                                              'bg-primary',
                                              submitting && 'opacity-50',
                                          )}
                                          onClick={() => submitReply(comment)}
                                          hoverClass='tap-active'
                                        >
                                            <Text className='text-desc text-white'>
                                                {submitting ? '...' : t('submit_comment')}
                                            </Text>
                                        </View>
                                    </View>
                                </View>
                            </View>
                        )}
                    </View>

                    {/* 递归渲染子评论 */}
                    {hasChildren(comment) && isExpanded(comment) && (
                        <View className='ml-[48rpx] mt-[16rpx] border-l-[2rpx] border-primary/20 pl-[24rpx]'>
                            {/* 加载中 */}
                            {isLoading(comment) && (
                                <View className='flex items-center justify-center gap-[12rpx] py-[32rpx]'>
                                    <Spinner size='sm' className='text-primary' />
                                    <Text className='text-desc text-textSec'>{t('loading')}</Text>
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
