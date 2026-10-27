import {View, Text, Image, Input, Textarea} from '@tarojs/components';
import {useMemo, useRef, useState} from 'react';
import Taro, {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import {fetchPost, fetchListPosts, getPostTitle, getPostSummary, getPostContent, getPostThumbnail} from '@/api/hooks/post';
import {
  getInteractionStatus,
  getCounts,
  likePost,
  unlikePost,
  watchPost,
  unwatchPost,
  recordView,
  extractCount,
} from '@/api/hooks/interaction';
import {fetchListComments, fetchCreateComment} from '@/api/hooks/comment';
import {useIsLogin} from '@/store/core/access/store';
import ContentViewer from '@/components/content/ContentViewer';
import CustomFieldsSection from '@/components/post/CustomFieldsSection';
import PostCard from '@/components/post/PostCard';
import {AppEmpty, Button} from '@/components/ui';
import type {
  contentservicev1_Post,
  contentservicev1_EditorType,
  commentservicev1_Comment,
  interactionservicev1_TargetType,
  interactionservicev1_CounterMetric,
} from '@/api/generated/app/service/v1';

const RELATED_FIELD_MASK = 'id,status,authorName,publishTime,thumbnail,availableLanguages,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary';

/**
 * 文章详情（阅读器）
 * 结构：封面 → 标题/摘要 → 作者/时间/浏览 → 正文（按 editorType 分流）→ 标签 →
 *       相关阅读 → 评论区 → 底部操作条（点赞 / 收藏 / 分享）
 */
export default function PostDetailPage() {
  const {t} = useTranslation();
  const router = useI18nRouter();
  const isLogin = useIsLogin();

  const postId = useMemo(() => {
    const raw = Taro.getCurrentInstance()?.router?.params?.id;
    const id = Number(raw);
    return Number.isInteger(id) && id > 0 ? id : 0;
  }, []);

  const [post, setPost] = useState<contentservicev1_Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [viewCount, setViewCount] = useState<number | null>(null);
  const [liked, setLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(0);
  const [watched, setWatched] = useState(false);
  const [related, setRelated] = useState<contentservicev1_Post[]>([]);
  const [comments, setComments] = useState<commentservicev1_Comment[]>([]);
  const [commentTotal, setCommentTotal] = useState(0);
  const [commentText, setCommentText] = useState('');
  const [guestName, setGuestName] = useState('');
  const [guestEmail, setGuestEmail] = useState('');
  const [submittingComment, setSubmittingComment] = useState(false);

  const loadComments = async () => {
    if (!postId) return;
    try {
      const res = await fetchListComments({
        paging: {page: 1, pageSize: 20},
        formValues: {objectId: postId, contentType: 'CONTENT_TYPE_POST', status: 'STATUS_APPROVED'},
        orderBy: ['-createdAt'],
      });
      setComments(res.items || []);
      setCommentTotal(Number(res.total ?? res.items?.length ?? 0));
    } catch (error) {
      console.error('Load comments failed:', error);
    }
  };

  // PV 上报去重：StrictMode/重复 useDidShow 不重复计数
  const viewReportedRef = useRef<number | null>(null);

  const loadAll = async () => {
    if (!postId) {
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      const loaded = await fetchPost(postId);
      setPost(loaded);

      // 浏览量上报（匿名可报），返回累加后的浏览量直接展示
      if (viewReportedRef.current !== postId) {
        viewReportedRef.current = postId;
        recordView('TARGET_TYPE_POST', postId)
          .then((resp) => setViewCount(resp.viewCount ?? null))
          .catch((error) => console.error('Record view failed:', error));
      }

      // 相关阅读：同分类（首个分类）下其余文章
      const categoryId = loaded.categoryIds?.[0];
      if (categoryId) {
        fetchListPosts({
          paging: {page: 1, pageSize: 6},
          formValues: {status: 'POST_STATUS_PUBLISHED', category_ids__in: [categoryId]},
          fieldMask: RELATED_FIELD_MASK,
          orderBy: ['-publishTime'],
        })
          .then((res) => setRelated((res.items || []).filter((p) => p.id !== postId).slice(0, 4)))
          .catch((error) => console.error('Load related posts failed:', error));
      }

      // 评论区：已通过审核的顶层评论，按时间倒序
      loadComments();

      // 计数（浏览/点赞/收藏）任何人都可读
      const metrics: interactionservicev1_CounterMetric[] = ['COUNTER_METRIC_LIKE', 'COUNTER_METRIC_VIEW', 'COUNTER_METRIC_WATCH'];
      getCounts('TARGET_TYPE_POST' as interactionservicev1_TargetType, [postId], metrics)
        .then((resp) => {
          setLikeCount(extractCount(resp, postId, 'COUNTER_METRIC_LIKE'));
          const views = extractCount(resp, postId, 'COUNTER_METRIC_VIEW');
          setViewCount((prev) => prev ?? (views || null));
        })
        .catch((error) => console.error('Load counts failed:', error));

      // 点赞/收藏态要求 viewer 身份，游客跳过（401 会触发全局重认证跳登录，不能让游客踩到）
      if (isLogin) {
        getInteractionStatus('TARGET_TYPE_POST' as interactionservicev1_TargetType, [postId])
          .then((resp) => {
            const status = resp.statuses?.[String(postId)];
            setLiked(!!status?.liked);
            setWatched(!!status?.watched);
          })
          .catch((error) => console.error('Load interaction status failed:', error));
      }
    } catch (error) {
      console.error('Load post failed:', error);
      setPost(null);
    } finally {
      setLoading(false);
    }
  };

  useDidShow(() => {
    loadAll();
  });

  usePageTitle(post ? (getPostTitle(post) || t('page.title.post_detail')) : t('page.title.post_detail'));

  const requireLogin = (action: () => void) => {
    if (!isLogin) {
      Taro.showToast({title: t('page.user.please_login'), icon: 'none'});
      router.push('/login');
      return;
    }
    action();
  };

  const toggleLike = () =>
    requireLogin(() => {
      if (!postId) return;
      const next = !liked;
      setLiked(next);
      setLikeCount((c) => Math.max(0, c + (next ? 1 : -1)));
      (next ? likePost : unlikePost)('TARGET_TYPE_POST' as interactionservicev1_TargetType, postId)
        .catch((error) => {
          console.error('Toggle like failed:', error);
          setLiked(!next);
          setLikeCount((c) => Math.max(0, c + (next ? -1 : 1)));
        });
    });

  const toggleWatch = () =>
    requireLogin(() => {
      if (!postId) return;
      const next = !watched;
      setWatched(next);
      (next ? watchPost : unwatchPost)(postId)
        .catch((error) => {
          console.error('Toggle watch failed:', error);
          setWatched(!next);
        });
    });

  const shareLink = () => {
    const url = typeof window !== 'undefined'
      ? `${window.location.origin}/posts/detail?id=${postId}`
      : `/posts/detail?id=${postId}`;
    Taro.setClipboardData({data: url})
      .then(() => Taro.showToast({title: t('page.post_detail.link_copied'), icon: 'none'}))
      .catch((error) => {
        console.error('Copy link failed:', error);
        Taro.showToast({title: t('ui.fallback.http.networkError'), icon: 'none'});
      });
  };

  const submitComment = async () => {
    const content = commentText.trim();
    if (!content || !postId || submittingComment) return;

    // 游客留言需留昵称与合法邮箱（对齐 react 端 CommentSection 校验）
    if (!isLogin) {
      const name = guestName.trim();
      const email = guestEmail.trim();
      if (!name) {
        Taro.showToast({title: t('comment.enter_nickname'), icon: 'none'});
        return;
      }
      if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        Taro.showToast({title: t('comment.enter_email'), icon: 'none'});
        return;
      }
    }

    setSubmittingComment(true);
    try {
      await fetchCreateComment({
        objectId: postId,
        contentType: 'CONTENT_TYPE_POST',
        content,
        status: 'STATUS_PENDING',
        ...(!isLogin
          ? {authorName: guestName.trim(), authorEmail: guestEmail.trim(), authorType: 'AUTHOR_TYPE_GUEST'}
          : {}),
      } as Partial<commentservicev1_Comment>);
      setCommentText('');
      Taro.showToast({title: t('comment.comment_submitted'), icon: 'success'});
      await loadComments();
    } catch (error) {
      console.error('Submit comment failed:', error);
      Taro.showToast({title: t('comment.submit_comment_failed'), icon: 'none'});
    } finally {
      setSubmittingComment(false);
    }
  };

  if (loading) {
    return (
      <View className='au-page w-full bg-pageBg pb-[200rpx]'>
        <View className='px-[24rpx] pt-[28rpx]'>
          <View className='au-card p-[32rpx] flex flex-col gap-[20rpx]'>
            <View className='au-skeleton w-[70%] h-[44rpx] rounded-[10rpx]' />
            <View className='au-skeleton w-[40%] h-[26rpx] rounded-[8rpx]' />
            <View className='au-skeleton w-full h-[300rpx] rounded-[16rpx] mt-[10rpx]' />
            <View className='au-skeleton w-full h-[26rpx] rounded-[8rpx]' />
            <View className='au-skeleton w-[92%] h-[26rpx] rounded-[8rpx]' />
            <View className='au-skeleton w-[96%] h-[26rpx] rounded-[8rpx]' />
          </View>
        </View>
      </View>
    );
  }

  if (!post) {
    return (
      <View className='au-page w-full bg-pageBg'>
        <View className='au-card mx-[24rpx] mt-[28rpx]'>
          <AppEmpty variant='error' description={t('page.post_detail.post_not_found')} />
          <View className='pb-[40rpx] flex justify-center'>
            <Button variant='hero' size='sm' className='rounded-full px-[48rpx]' onClick={() => router.back()}>
              {t('page.post_detail.back')}
            </Button>
          </View>
        </View>
      </View>
    );
  }

  const title = getPostTitle(post) || t('page.post_detail.untitled');
  const summary = getPostSummary(post);
  const content = getPostContent(post);
  const cover = post.thumbnail || '';
  // 编辑器类型 → ContentViewer：富文本(html) / 其余按 markdown；纯文本单独分流
  const editorType = (post.editorType || 'EDITOR_TYPE_MARKDOWN') as contentservicev1_EditorType;
  const contentType = editorType === 'EDITOR_TYPE_RICH_TEXT'
    ? 'html'
    : editorType === 'EDITOR_TYPE_PLAIN_TEXT' ? 'text' : 'markdown';

  return (
    <View className='au-page w-full bg-pageBg pb-[200rpx]'>
      <View className='px-[24rpx] pt-[28rpx]'>
        {/* 封面 */}
        {cover && (
          <View className='au-anim rounded-[28rpx] overflow-hidden mb-[28rpx] shadow-[0_16rpx_48rpx_rgba(2,8,26,0.12)]'>
            <Image src={cover} mode='widthFix' className='w-full' lazyLoad />
          </View>
        )}

        {/* 标题与元信息 */}
        <View className='au-anim au-anim-1 mb-[28rpx]'>
          <Text className='text-title font-bold text-textMain leading-[1.35] block'>{title}</Text>
          {summary && (
            <Text className='text-desc text-textSec leading-[1.7] block mt-[16rpx]'>{summary}</Text>
          )}
          <View className='flex items-center gap-[26rpx] mt-[20rpx]'>
            {post.authorName && (
              <View className='flex items-center gap-[8rpx]'>
                <View className='w-[44rpx] h-[44rpx] rounded-full bg-primary-tint flex items-center justify-center'>
                  <Text className='text-tips font-bold text-primary'>
                    {post.authorName.charAt(0).toUpperCase()}
                  </Text>
                </View>
                <Text className='text-tips text-textSec'>{post.authorName}</Text>
              </View>
            )}
            {post.publishTime && (
              <View className='flex items-center gap-[6rpx]'>
                <XIcon name='carbon:calendar' size={18} className='text-textWeak' />
                <Text className='text-tips text-textThird'>{String(post.publishTime).slice(0, 10)}</Text>
              </View>
            )}
            {viewCount !== null && (
              <View className='flex items-center gap-[6rpx]'>
                <XIcon name='carbon:view' size={18} className='text-textWeak' />
                <Text className='text-tips text-textThird'>{viewCount}</Text>
              </View>
            )}
          </View>
        </View>

        {/* 正文 */}
        <View className='au-card au-card--glow p-[32rpx] au-anim au-anim-2 mb-[28rpx]'>
          {content ? (
            <ContentViewer content={content} type={contentType} />
          ) : (
            <AppEmpty description={t('common.noData')} />
          )}
        </View>

        {/* 自定义字段（内容建模前台读侧） */}
        <CustomFieldsSection postId={postId} customFields={post.customFields} />

        {/* 相关阅读 */}
        {related.length > 0 && (
          <View className='au-anim au-anim-3 mb-[28rpx]'>
            <Text className='au-section-title block mb-[22rpx]'>{t('page.post_detail.related_posts')}</Text>
            <View className='flex flex-col gap-[22rpx]'>
              {related.map((rel) => (
                <PostCard
                  key={rel.id}
                  title={getPostTitle(rel) || t('page.home.untitled')}
                  thumbnail={getPostThumbnail(rel)}
                  publishTime={rel.publishTime}
                  onClick={() => rel.id && router.push(`/post/detail?id=${rel.id}`)}
                />
              ))}
            </View>
          </View>
        )}

        {/* 评论区 */}
        <View className='au-anim au-anim-4'>
          <Text className='au-section-title block mb-[22rpx]'>
            {t('comment.comments')} ({commentTotal})
          </Text>

          {/* 评论输入（游客留言需昵称+邮箱，进入待审核；登录用户直接提交） */}
          <View className='au-card p-[28rpx] mb-[22rpx]'>
            {!isLogin && (
              <View className='flex gap-[16rpx] mb-[16rpx]'>
                <View className='au-input flex-1 h-[76rpx]'>
                  <Input
                    className='au-input-field'
                    value={guestName}
                    placeholder={t('comment.nickname')}
                    placeholderClass='text-textWeak'
                    onInput={(e) => setGuestName(e.detail.value)}
                  />
                </View>
                <View className='au-input flex-1 h-[76rpx]'>
                  <Input
                    className='au-input-field'
                    value={guestEmail}
                    placeholder={t('comment.email')}
                    placeholderClass='text-textWeak'
                    onInput={(e) => setGuestEmail(e.detail.value)}
                  />
                </View>
              </View>
            )}
            <Textarea
              className='au-input-field min-h-[120rpx] w-full'
              value={commentText}
              placeholder={t('comment.comment_content')}
              placeholderClass='text-textWeak'
              maxlength={500}
              onInput={(e) => setCommentText(e.detail.value)}
            />
            <View className='flex justify-end mt-[16rpx]'>
              <Button
                variant='hero'
                size='sm'
                className='rounded-full px-[44rpx]'
                loading={submittingComment}
                disabled={!commentText.trim()}
                onClick={submitComment}
              >
                {t('comment.submit_comment')}
              </Button>
            </View>
          </View>

          {comments.length === 0 ? (
            <View className='au-card'>
              <AppEmpty description={t('comment.no_comments')} />
            </View>
          ) : (
            <View className='au-card overflow-hidden'>
              {comments.map((comment, index) => (
                <View
                  key={comment.id ?? index}
                  className={`p-[28rpx] ${index < comments.length - 1 ? 'border-b-[length:1rpx] border-splitLine' : ''}`}
                >
                  <View className='flex items-center gap-[14rpx] mb-[12rpx]'>
                    <View className='w-[52rpx] h-[52rpx] rounded-full bg-primary-tint flex items-center justify-center'>
                      <Text className='text-tips font-bold text-primary'>
                        {(comment.authorName || 'A').charAt(0).toUpperCase()}
                      </Text>
                    </View>
                    <Text className='text-desc font-medium text-textMain'>
                      {comment.authorName || t('comment.author')}
                    </Text>
                    {comment.createdAt && (
                      <Text className='text-tips text-textWeak ml-auto'>
                        {String(comment.createdAt).slice(0, 10)}
                      </Text>
                    )}
                  </View>
                  <Text className='text-desc text-textMain leading-[1.7]'>{comment.content}</Text>
                </View>
              ))}
            </View>
          )}
        </View>
      </View>

      {/* 底部操作条：点赞 / 收藏 / 分享 */}
      <View className='au-bottom-bar au-glass flex items-center gap-[20rpx]'>
        <View
          className='flex-1 h-[88rpx] rounded-full flex items-center justify-center gap-[10rpx] border-[length:1rpx] border-splitLine bg-cardBg'
          hoverClass='tap-active'
          onClick={toggleLike}
        >
          <XIcon
            name={liked ? 'carbon:favorite-filled' : 'carbon:favorite'}
            size={30}
            className={liked ? 'text-[#ff4d4f]' : 'text-textSec'}
          />
          <Text className={`text-desc ${liked ? 'text-[#ff4d4f] font-bold' : 'text-textSec'}`}>
            {likeCount > 0 ? likeCount : t('comment.like')}
          </Text>
        </View>
        <View
          className='flex-1 h-[88rpx] rounded-full flex items-center justify-center gap-[10rpx] border-[length:1rpx] border-splitLine bg-cardBg'
          hoverClass='tap-active'
          onClick={toggleWatch}
        >
          <XIcon
            name={watched ? 'carbon:star-filled' : 'carbon:star'}
            size={30}
            className={watched ? 'text-primary' : 'text-textSec'}
          />
          <Text className={`text-desc ${watched ? 'text-primary font-bold' : 'text-textSec'}`}>
            {watched ? t('page.post_detail.bookmarked') : t('page.post_detail.bookmark')}
          </Text>
        </View>
        <View
          className='flex-1 h-[88rpx] rounded-full flex items-center justify-center gap-[10rpx] btn-hero'
          hoverClass='tap-active'
          onClick={shareLink}
        >
          <XIcon name='carbon:share' size={28} className='text-white' />
          <Text className='text-desc text-white font-bold'>{t('page.post_detail.share')}</Text>
        </View>
      </View>
    </View>
  );
}
