import React, {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Input, Textarea} from '@tarojs/components';

import {AppEmpty} from '@/components/ui';
import {useCommentStore} from '@/store/slices/comment/hooks';
import type {
  commentservicev1_Comment,
  commentservicev1_Comment_ContentType, commentservicev1_ListCommentResponse,
} from '@/api/generated/app/service/v1';

import CommentTree from '../CommentTree';

import './index.scss';

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
  const {t} = useTranslation('comment');
  const commentStore = useCommentStore();

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
      const res = await commentStore.listComment({
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
      await commentStore.createComment({
        postId: objectId,
        content: newComment.content,
        authorName: newComment.authorName,
        authorEmail: newComment.authorEmail,
        status: 'STATUS_PENDING',
      });

      alert(t('comment_submitted'));
      setNewComment({content: '', authorName: '', authorEmail: ''});
      setCurrentPage(1);
      await loadComments(true);
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
      await commentStore.createComment({
        postId: objectId,
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
      const res = await commentStore.listComment({
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
    <View className='comments-section'>
      <View className='section-header'>
        <Text className='section-title'>
          💬 {t('comments_count', {count: displayComments.length})}
        </Text>
      </View>

      {/* Comment Form */}
      <View className='comment-form'>
        <View className='form-header'>
          <View className='form-icon'>✏️</View>
          <Text className='form-title'>{t('write_comment')}</Text>
        </View>
        <View className='comment-form-wrapper'>
          <View className='form-row'>
            <View className='form-group'>
              <Input
                value={newComment.authorName}
                onInput={(e) => setNewComment({...newComment, authorName: e.detail.value})}
                placeholder={t('nickname') + ' *'}
                className='input-field'
                disabled={submitting}
              />
            </View>
            <View className='form-group'>
              <Input
                value={newComment.authorEmail}
                onInput={(e) => setNewComment({...newComment, authorEmail: e.detail.value})}
                placeholder={t('email') + ' *'}
                className='input-field'
                disabled={submitting}
              />
            </View>
          </View>
          <View className='form-group'>
            <Textarea
              value={newComment.content}
              onInput={(e) => setNewComment({...newComment, content: e.detail.value})}
              placeholder={t('write_comment')}
              className='textarea-field'
              disabled={submitting}
              maxlength={1000}
            />
            <View className='char-count'>
              {newComment.content.length} / 1000
            </View>
          </View>
          <View className='form-actions'>
            <View
              className={`submit-btn ${submitting ? 'disabled' : ''}`}
              onClick={handleSubmitComment}
            >
              <Text>{submitting ? t('submitting') : t('submit_comment')}</Text>
            </View>
            <View className='form-tip'>
              <Text>ℹ️ {t('fill_form_info')}</Text>
            </View>
          </View>
        </View>
      </View>

      {/* Comments List */}
      {hasComments ? (
        <View className='comments-list'>
          <CommentTree
            comments={displayComments}
            onReply={handleReply}
            onLoadChildren={loadChildren}
          />

          {/* 加载更多按钮 */}
          {showLoadMore && (
            <View className='load-more-container'>
              <View
                className='load-more-btn'
                onClick={loadMoreComments}
              >
                <Text>{loadingMore ? t('loading') : t('load_more')}</Text>
              </View>
            </View>
          )}

          {/* 没有更多提示 */}
          {!hasMore && hasComments && (
            <View className='no-more-text'>
              <Text>{t('no_more')}</Text>
            </View>
          )}
        </View>
      ) : (
        <AppEmpty
          description={t('no_comments')}
          inContainer
        />
      )}
    </View>
  );
};

export default CommentSection;
