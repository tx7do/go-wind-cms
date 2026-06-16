import Taro from '@tarojs/taro';
import {View, Text, Input, Textarea} from '@tarojs/components';
import React, {useState, useEffect} from 'react';

import {useTranslations} from '@/lib/next-intl-compat';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';
import {AppEmpty} from '@/components/ui';
import {Skeleton} from '@/components/ui/skeleton';
import {Spinner} from '@/components/ui/spinner';
import {fetchListComments, createComment as createCommentApi} from '@/api/hooks/comment';
import type {
    commentservicev1_Comment,
    commentservicev1_Comment_ContentType,
    commentservicev1_ListCommentResponse,
} from '@/api/generated/app/service/v1';

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
    onUpdateComments,
}) => {
    const t = useTranslations('comment');

    // 状态
    const [submitting, setSubmitting] = useState(false);
    const [loading, setLoading] = useState(false);
    const [comments, setComments] = useState<commentservicev1_Comment[]>([]);

    // 分页
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize] = useState(20);
    const [_total, setTotal] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);

    // 表单折叠：默认收起昵称/邮箱，点击评论框后展开
    const [formExpanded, setFormExpanded] = useState(false);

    const [newComment, setNewComment] = useState<CommentForm>({
        content: '',
        authorName: '',
        authorEmail: '',
    });

    // 计算属性
    const displayComments = comments;
    const hasComments = displayComments.length > 0;
    const showLoadMore = hasMore && !loadingMore;

    const textLength = newComment.content.length;
    const isEmpty = textLength === 0;
    const isOverLimit = textLength > 1000;

    // 加载评论
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
                paging: {page: reset ? 1 : currentPage, pageSize},
                formValues: {
                    objectId,
                    contentType,
                    status: 'STATUS_APPROVED',
                },
            }) as unknown as commentservicev1_ListCommentResponse;

            const newComments = res.items || [];
            const newTotal = res.total || 0;

            if (reset) {
                setComments(newComments);
            } else {
                setComments(prev => [...prev, ...newComments]);
            }

            setTotal(newTotal);
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

    function loadMoreComments() {
        if (hasMore && !loadingMore) {
            loadComments();
        }
    }

    // 提交评论
    async function handleSubmitComment() {
        const textContent = newComment.content.trim();
        if (!textContent) {
            Taro.showToast({title: t('empty_content'), icon: 'none'});
            return;
        }
        if (!newComment.authorName.trim()) {
            Taro.showToast({title: t('invalid_nickname'), icon: 'none'});
            return;
        }
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(newComment.authorEmail)) {
            Taro.showToast({title: t('invalid_email'), icon: 'none'});
            return;
        }
        if (!objectId || submitting) return;

        setSubmitting(true);
        try {
            await createCommentApi({
                objectId,
                contentType: contentType ?? undefined,
                content: newComment.content,
                authorName: newComment.authorName,
                authorEmail: newComment.authorEmail,
                status: 'STATUS_PENDING',
            });

            Taro.showToast({title: t('comment_submitted'), icon: 'success'});
            setNewComment({content: '', authorName: '', authorEmail: ''});
            setFormExpanded(false);
            setCurrentPage(1);
            await loadComments(true);
        } catch (error) {
            console.error('Submit comment failed:', error);
            Taro.showToast({title: t('submit_comment_failed'), icon: 'none'});
        } finally {
            setSubmitting(false);
        }
    }

    // 处理回复
    async function handleReply(comment: commentservicev1_Comment, content: string) {
        if (!content.trim()) {
            Taro.showToast({title: t('empty_content'), icon: 'none'});
            return;
        }
        if (!objectId || submitting) return;

        setSubmitting(true);
        try {
            await createCommentApi({
                objectId,
                contentType: contentType ?? undefined,
                content: content.trim(),
                authorName: comment.authorName,
                authorEmail: comment.authorEmail,
                parentId: comment.id,
                replyToId: comment.id,
                status: 'STATUS_PENDING',
            });

            Taro.showToast({title: t('comment_posted'), icon: 'success'});
            setCurrentPage(1);
            await loadComments(true);
        } catch (error) {
            console.error('Submit reply failed:', error);
            Taro.showToast({title: t('submit_comment_failed'), icon: 'none'});
        } finally {
            setSubmitting(false);
        }
    }

    // 加载子评论
    async function loadChildren(parentComment: commentservicev1_Comment) {
        if (!objectId || !contentType) return;

        try {
            const res = await fetchListComments({
                paging: {page: 1, pageSize: 50},
                formValues: {
                    objectId,
                    contentType,
                    parentId: parentComment.id,
                    status: 'STATUS_APPROVED',
                },
            }) as unknown as commentservicev1_ListCommentResponse;

            parentComment.children = res.items || [];
        } catch (error) {
            console.error('Load children failed:', error);
            throw error;
        }
    }

    useEffect(() => {
        if (objectId && contentType) {
            loadComments(true);
        }
    }, [objectId, contentType]);

    // ===== 加载态 =====
    if (loading) {
        return (
            <View className='flex flex-col'>
                {/* 标题骨架 */}
                <View className='flex items-center gap-[8rpx] bg-pageBg px-[32rpx] py-[16rpx]'>
                    <Skeleton className='h-[36rpx] w-[240rpx] rounded-[8rpx]' />
                </View>
                {/* 表单骨架 */}
                <View className='bg-cardBg p-[32rpx]'>
                    <Skeleton className='h-[80rpx] w-full rounded-[12rpx] mb-[16rpx]' />
                    <View className='flex justify-end'>
                        <Skeleton className='h-[56rpx] w-[160rpx] rounded-[8rpx]' />
                    </View>
                </View>
                {/* 评论骨架 */}
                {Array.from({length: 2}).map((_, i) => (
                    <View key={i} className='bg-cardBg p-[32rpx] mt-[12rpx]'>
                        <View className='flex items-center gap-[12rpx] mb-[12rpx]'>
                            <Skeleton className='h-[56rpx] w-[56rpx] rounded-full' />
                            <Skeleton className='h-[28rpx] w-[140rpx] rounded-[8rpx]' />
                        </View>
                        <Skeleton className='h-[24rpx] w-full rounded-[8rpx] mb-[8rpx]' />
                        <Skeleton className='h-[24rpx] w-[60%] rounded-[8rpx]' />
                    </View>
                ))}
            </View>
        );
    }

    return (
        <View className='flex flex-col'>
            {/* 区块标题 */}
            <View className='flex items-center justify-between bg-pageBg px-[32rpx] py-[14rpx]'>
                <View className='flex items-center gap-[8rpx]'>
                    <XIcon name='carbon:chat' size={18} className='text-primary' />
                    <Text className='text-desc font-bold text-textMain'>
                        {t('comments_count', {count: displayComments.length})}
                    </Text>
                </View>
            </View>

            {/* 评论表单 */}
            <View className='bg-cardBg px-[32rpx] py-[20rpx]'>
                {/* 评论输入框 - 点击展开 */}
                {!formExpanded ? (
                    <View
                      className='flex items-center gap-[12rpx] min-h-[72rpx] px-[20rpx] rounded-[12rpx] bg-pageBg border-[1rpx] border-splitLine'
                      onClick={() => setFormExpanded(true)}
                    >
                        <XIcon name='carbon:edit' size={16} className='text-textThird' />
                        <Text className='text-desc text-textThird'>{t('write_comment')}</Text>
                    </View>
                ) : (
                    <View className='flex flex-col gap-[24rpx]'>
                        {/* 昵称 + 邮箱 */}
                        <View className='flex gap-[12rpx]'>
                            <View className='flex-1 rounded-[8rpx] bg-pageBg border border-splitLine/60 px-[16rpx] overflow-hidden'>
                                <Input
                                  value={newComment.authorName}
                                  onInput={(e) => setNewComment({...newComment, authorName: e.detail.value})}
                                  placeholder={t('nickname')}
                                  className='h-[64rpx] text-tips text-textMain'
                                  disabled={submitting}
                                />
                            </View>
                            <View className='flex-1 rounded-[8rpx] bg-pageBg border border-splitLine/60 px-[16rpx] overflow-hidden'>
                                <Input
                                  value={newComment.authorEmail}
                                  onInput={(e) => setNewComment({...newComment, authorEmail: e.detail.value})}
                                  placeholder={t('email')}
                                  className='h-[64rpx] text-tips text-textMain'
                                  disabled={submitting}
                                />
                            </View>
                        </View>

                        {/* 评论内容 */}
                        <View className='w-full rounded-[12rpx] bg-pageBg border border-splitLine/60 p-[16rpx]'>
                            <Textarea
                              value={newComment.content}
                              onInput={(e) => setNewComment({...newComment, content: e.detail.value ?? ''})}
                              placeholder={t('write_comment')}
                              maxlength={1000}
                              className='w-full min-h-[120rpx] text-desc text-textMain'
                              disabled={submitting}
                            />
                        </View>

                        {/* 操作栏 */}
                        <View className='flex items-center justify-between'>
                            <Text className={cn('text-tips', isOverLimit ? 'text-danger' : 'text-textThird')}>
                                {textLength} / 1000
                            </Text>
                            <View className='flex items-center gap-[12rpx]'>
                                <View
                                  className='flex items-center justify-center px-[32rpx] py-[16rpx] rounded-[8rpx] bg-pageBg border border-splitLine/60'
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    setFormExpanded(false);
                                  }}
                                >
                                    <Text className='text-desc text-textSec'>{t('cancel')}</Text>
                                </View>
                                <View
                                  className={cn(
                                    'flex items-center gap-[4rpx] px-[24rpx] py-[10rpx] rounded-[8rpx]',
                                    'bg-primary',
                                    (submitting || isEmpty || isOverLimit) && 'opacity-50',
                                  )}
                                  onClick={() => {
                                    if (!isEmpty && !submitting && !isOverLimit) {
                                        handleSubmitComment();
                                    }
                                  }}
                                  hoverClass='tap-active'
                                >
                                    {submitting && <Spinner size='sm' className='text-white' />}
                                    <Text className='text-tips text-white font-medium'>{t('submit_comment')}</Text>
                                </View>
                            </View>
                        </View>
                    </View>
                )}
            </View>

            {/* 评论列表 */}
            {hasComments ? (
                <View className='flex flex-col'>
                    <CommentTree
                      comments={displayComments}
                      onReply={handleReply}
                      onLoadChildren={loadChildren}
                    />

                    {/* 加载更多 */}
                    {showLoadMore && (
                        <View className='flex justify-center py-[24rpx] pb-[64rpx]'>
                            <View
                              className={cn(
                                'flex items-center gap-[8rpx] px-[32rpx] py-[16rpx] rounded-full',
                                'bg-cardBg border-[1rpx] border-splitLine',
                              )}
                              onClick={loadMoreComments}
                              hoverClass='tap-active'
                            >
                                {loadingMore ? (
                                    <Spinner size='sm' className='text-textThird' />
                                ) : (
                                    <XIcon name='carbon:chevron-down' size={14} className='text-textThird' />
                                )}
                                <Text className='text-tips text-textSec'>{t('load_more')}</Text>
                            </View>
                        </View>
                    )}

                    {/* 没有更多 */}
                    {!hasMore && hasComments && (
                        <View className='py-[24rpx] text-center'>
                            <Text className='text-tips text-textThird'>{t('no_more')}</Text>
                        </View>
                    )}
                </View>
            ) : (
                <View className='bg-cardBg py-[80rpx] flex items-center justify-center'>
                    <AppEmpty description={t('no_comments')} inContainer />
                </View>
            )}
        </View>
    );
};

export default CommentSection;
