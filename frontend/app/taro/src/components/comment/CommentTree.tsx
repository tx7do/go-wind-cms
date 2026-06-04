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
        if (replyingCommentId === comment.id) {
            setReplyingCommentId(null);
            setReplyContent('');
        } else {
            setReplyingCommentId(comment.id || 0);
            setReplyContent('');
        }
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

    function handleShare(comment: commentservicev1_Comment) {
        const commentUrl = `#comment-${comment.id}`;
        Taro.setClipboardData({
            data: commentUrl,
            success: () => Taro.showToast({title: t('link_copied'), icon: 'success'}),
            fail: () => Taro.showToast({title: t('copy_failed'), icon: 'none'}),
        });
    }

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
        <View className='flex flex-col'>
            {comments.map((comment, index) => (
                <View key={comment.id}>
                    {/* 评论主体 */}
                    <View className='bg-cardBg px-[32rpx] py-[20rpx]'>
                        {/* 头像 - 独立左侧 */}
                        <View className='flex items-start gap-[12rpx]'>
                            {/* 头像 */}
                            <View className='flex items-center justify-center w-[48rpx] h-[48rpx] rounded-full bg-primary/10 flex-shrink-0'>
                                <Text className='text-tips font-bold text-primary'>
                                    {comment.authorName?.charAt(0)?.toUpperCase() || 'U'}
                                </Text>
                            </View>
                            {/* 右侧内容区：用户名 + 正文 + 操作栏，共享同一条左侧对齐线 */}
                            <View className='flex-1 min-w-0 flex flex-col'>
                                {/* 用户名行 */}
                                <View className='flex items-center gap-[8rpx] flex-wrap'>
                                    <Text className='text-desc font-medium text-textMain'>
                                        {comment.authorName}
                                    </Text>
                                    {isOwnerReply(comment) && (
                                        <View className='flex items-center gap-[2rpx] rounded-[6rpx] bg-primary/15 px-[10rpx] py-[2rpx]'>
                                            <XIcon name='carbon:verified' size={12} className='text-primary' />
                                            <Text className='text-tips font-medium text-primary'>{t('owner_reply')}</Text>
                                        </View>
                                    )}
                                    <Text className='text-tips text-textSec'>{formatDate(comment.createdAt)}</Text>
                                </View>

                                {/* 正文 */}
                                <Text className='text-desc text-textMain leading-[1.7] mt-[6rpx] block'>
                                    {comment.content}
                                </Text>

                                {/* 操作栏 */}
                                <View className='flex items-center gap-[20rpx] mt-[12rpx]'>
                                    {/* 点赞 */}
                                    <View
                                      className='flex items-center gap-[4rpx]'
                                      onClick={() => handleLike(comment)}
                                      hoverClass='tap-active'
                                    >
                                        <XIcon
                                          name={likedComments.has(comment.id || 0) ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'}
                                          size={14}
                                          className={likedComments.has(comment.id || 0) ? 'text-primary' : 'text-textSec'}
                                        />
                                        <Text
                                          className={cn(
                                            'text-tips',
                                            likedComments.has(comment.id || 0) ? 'text-primary' : 'text-textSec',
                                          )}
                                        >
                                            {(comment.likeCount || 0) + (likedComments.has(comment.id || 0) ? 1 : 0)}
                                        </Text>
                                    </View>

                                    {/* 回复 */}
                                    <View
                                      className='flex items-center gap-[4rpx]'
                                      onClick={() => handleReply(comment)}
                                      hoverClass='tap-active'
                                    >
                                        <XIcon name='carbon:chat' size={14} className='text-textSec' />
                                        <Text className='text-tips text-textSec'>{t('reply')}</Text>
                                    </View>

                                    {/* 分享 */}
                                    <View
                                      className='flex items-center gap-[4rpx]'
                                      onClick={() => handleShare(comment)}
                                      hoverClass='tap-active'
                                    >
                                        <XIcon name='carbon:share' size={14} className='text-textSec' />
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
                                              name={isExpanded(comment) ? 'carbon:chevron-up' : 'carbon:chevron-down'}
                                              size={12}
                                              className='text-primary'
                                            />
                                            <Text className='text-tips text-primary'>
                                                {isExpanded(comment)
                                                    ? t('hide_replies')
                                                    : t('view_replies', {count: comment.replyCount})
                                                }
                                            </Text>
                                        </View>
                                    ) : null}
                                </View>
                            </View>
                        </View>

                        {/* 回复表单 - 内嵌在评论下方 */}
                        {replyingCommentId === comment.id && (
                            <View className='mt-[12rpx] ml-[60rpx] rounded-[8rpx] p-[16rpx]'>
                                <View className='rounded-[8rpx] bg-pageBg border border-splitLine/60 p-[12rpx]'>
                                    <Textarea
                                      value={replyContent}
                                      onInput={(e) => setReplyContent(e.detail.value ?? '')}
                                      placeholder={t('write_comment')}
                                      maxlength={1000}
                                      className='w-full min-h-[80rpx] text-tips text-textMain'
                                      disabled={submitting}
                                    />
                                </View>
                                <View className='flex items-center justify-between mt-[8rpx]'>
                                    <Text className='text-tips text-textThird'>
                                        {replyContent.length}/1000
                                    </Text>
                                    <View className='flex items-center gap-[8rpx]'>
                                        <View
                                          className='px-[16rpx] py-[6rpx] rounded-[6rpx]'
                                          onClick={cancelReply}
                                          hoverClass='tap-active'
                                        >
                                            <Text className='text-tips text-textSec'>{t('cancel')}</Text>
                                        </View>
                                        <View
                                          className={cn(
                                            'flex items-center gap-[4rpx] px-[20rpx] py-[6rpx] rounded-[6rpx]',
                                            'bg-primary',
                                            (submitting || !replyContent.trim()) && 'opacity-50',
                                          )}
                                          onClick={() => submitReply(comment)}
                                          hoverClass='tap-active'
                                        >
                                            {submitting && <Spinner size='sm' className='text-white' />}
                                            <Text className='text-tips text-white font-medium'>{t('reply')}</Text>
                                        </View>
                                    </View>
                                </View>
                            </View>
                        )}
                    </View>

                    {/* 子评论区域 */}
                    {hasChildren(comment) && isExpanded(comment) && (
                        <View className='ml-[60rpx] border-l-[2rpx] border-primary/20 pl-[16rpx]'>
                            {/* 加载中 */}
                            {isLoading(comment) && (
                                <View className='flex items-center gap-[8rpx] py-[24rpx]'>
                                    <Spinner size='sm' className='text-primary' />
                                    <Text className='text-tips text-textSec'>{t('loading')}</Text>
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

                    {/* 评论间分割线 */}
                    {index < comments.length - 1 && (
                        <View className='h-[1rpx] bg-splitLine ml-[24rpx]' />
                    )}
                </View>
            ))}
        </View>
    );
};

export default CommentTree;
