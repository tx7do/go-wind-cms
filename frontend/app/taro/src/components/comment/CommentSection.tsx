import Taro from '@tarojs/taro';
import {View, Text, Input} from '@tarojs/components';
import React, {useState, useEffect} from 'react';

import {useTranslations} from '@/lib/next-intl-compat';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';
import {AppEmpty} from '@/components/ui';
import {Skeleton} from '@/components/ui/skeleton';
import {fetchListComments} from '@/api/hooks/comment';
import {createComment as createCommentApi} from '@/api/service/comment';
import type {
    commentservicev1_Comment,
    commentservicev1_Comment_ContentType,
    commentservicev1_ListCommentResponse,
} from '@/api/generated/app/service/v1';

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

    const [newComment, setNewComment] = useState<CommentForm>({
        content: '',
        authorName: '',
        authorEmail: '',
    });

    // 计算属性
    const displayComments = comments;
    const hasComments = displayComments.length > 0;
    const showLoadMore = hasMore && !loadingMore;

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

    // 加载态
    if (loading) {
        return (
            <View className='flex flex-col gap-[24rpx] py-[32rpx]'>
                <View className='flex items-center gap-[12rpx]'>
                    <XIcon name='carbon:chat' size={20} className='text-primary' />
                    <Skeleton className='h-[36rpx] w-[200rpx] rounded-[8rpx]' />
                </View>
                {Array.from({length: 3}).map((_, i) => (
                    <View key={i} className='rounded-[16rpx] bg-cardBg p-[24rpx]'>
                        <View className='flex items-center gap-[16rpx] mb-[16rpx]'>
                            <Skeleton className='h-[64rpx] w-[64rpx] rounded-full' />
                            <Skeleton className='h-[28rpx] w-[160rpx] rounded-[8rpx]' />
                        </View>
                        <Skeleton className='h-[24rpx] w-full rounded-[8rpx] mb-[8rpx]' />
                        <Skeleton className='h-[24rpx] w-[70%] rounded-[8rpx]' />
                    </View>
                ))}
            </View>
        );
    }

    return (
        <View className='flex flex-col'>
            {/* 标题栏 - 灰色背景 + 底部蓝线 */}
            <View className='flex items-center gap-[12rpx] bg-pageBg px-[24rpx] py-[16rpx] border-b-[2rpx] border-primary'>
                <XIcon name='carbon:chat' size={20} className='text-primary' />
                <Text className='text-card-title font-bold text-textMain'>
                    {t('comments_count', {count: displayComments.length})}
                </Text>
            </View>

            {/* 评论表单 - 白色卡片 */}
            <View className='bg-cardBg p-[24rpx] mb-[16rpx]'>
                {/* 表单标题 */}
                <View className='flex items-center gap-[8rpx] mb-[16rpx]'>
                    <XIcon name='carbon:edit' size={18} className='text-primary' />
                    <Text className='text-desc font-bold text-textMain'>{t('write_comment')}</Text>
                </View>

                {/* 昵称 + 邮箱 - 灰色背景合并 */}
                <View className='flex gap-[16rpx] mb-[16rpx] bg-pageBg rounded-[12rpx] px-[16rpx] py-[12rpx]'>
                    <Input
                      value={newComment.authorName}
                      onInput={(e) => setNewComment({...newComment, authorName: e.detail.value})}
                      placeholder={t('nickname') + ' *'}
                      className='flex-1 h-[48rpx] text-desc text-textMain bg-transparent'
                      disabled={submitting}
                    />
                    <Input
                      value={newComment.authorEmail}
                      onInput={(e) => setNewComment({...newComment, authorEmail: e.detail.value})}
                      placeholder={t('email') + ' *'}
                      className='flex-1 h-[48rpx] text-desc text-textMain bg-transparent'
                      disabled={submitting}
                    />
                </View>

                {/* 内容编辑器 */}
                <RichTextEditor
                  value={newComment.content}
                  onChange={(content) => setNewComment({...newComment, content})}
                  onSubmit={handleSubmitComment}
                  submitting={submitting}
                  placeholder={t('write_comment')}
                  maxLength={1000}
                  submitLabel={t('submit_comment')}
                />

                {/* 提示 */}
                <View className='flex items-center gap-[6rpx] mt-[12rpx]'>
                    <XIcon name='carbon:information' size={14} className='text-primary' />
                    <Text className='text-tips text-textSec'>{t('fill_form_info')}</Text>
                </View>
            </View>

            {/* 评论列表 */}
            {hasComments ? (
                <View className='flex flex-col gap-[16rpx]'>
                    <CommentTree
                      comments={displayComments}
                      onReply={handleReply}
                      onLoadChildren={loadChildren}
                    />

                    {/* 加载更多 - 居中独立按钮 */}
                    {showLoadMore && (
                        <View className='flex justify-center py-[24rpx]'>
                            <View
                              className={cn(
                                'flex items-center gap-[8rpx] px-[40rpx] py-[20rpx] rounded-[12rpx]',
                                'bg-cardBg text-desc text-textSec',
                              )}
                              onClick={loadMoreComments}
                              hoverClass='tap-active'
                            >
                                <XIcon name='carbon:chevron-down' size={16} />
                                <Text>{t('load_more')}</Text>
                            </View>
                        </View>
                    )}

                    {/* 没有更多 */}
                    {!hasMore && hasComments && (
                        <View className='py-[24rpx] text-center'>
                            <Text className='text-tips text-textSec'>{t('no_more')}</Text>
                        </View>
                    )}
                </View>
            ) : (
                <AppEmpty description={t('no_comments')} inContainer />
            )}
        </View>
    );
};

export default CommentSection;
