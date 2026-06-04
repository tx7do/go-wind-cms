import {useEffect, useState, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';
import Taro from '@tarojs/taro';
import CommentSection from '@/components/comment/CommentSection';
import ContentViewer from '@/components/content/ContentViewer';
import PostList from '@/components/post/PostList';
import BackToTop from '@/components/layout/BackToTop';
import XIcon from '@/plugins/xicon';
import {formatDate} from '@/utils';
import {useI18nRouter} from '@/i18n/helpers';
import {usePreferences} from '@/core/preferences';
import {contentservicev1_Post} from '@/api/generated/app/service/v1';
import {fetchPost, getPostTitle, getPostContent, getPostThumbnail, getPostSummary} from '@/api/hooks/post';
import {Skeleton} from '@/components/ui/skeleton';

export default function PostDetailPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    const {isDark} = usePreferences();
    const [post, setPost] = useState<contentservicev1_Post | null>(null);
    const [loading, setLoading] = useState(true);
    const [isLiked, setIsLiked] = useState(false);
    const [isBookmarked, setIsBookmarked] = useState(false);

    const postId = useMemo(() => {
        let id: string | null = null;
        if (typeof Taro.getCurrentInstance === 'function') {
            const instance = Taro.getCurrentInstance();
            const routeId = instance?.router?.params?.id;
            id = typeof routeId === 'string' ? routeId : null;
        } else {
            const pages = Taro.getCurrentPages();
            const currentPage = pages[pages.length - 1];
            const pageId = (currentPage as any).options?.id;
            id = typeof pageId === 'string' ? pageId : null;
        }
        return id ? parseInt(id) : null;
    }, []);

    const displayTitle = useMemo(() => post ? getPostTitle(post) : '', [post]);

    // 文章加载完成后，用文章标题作为导航栏标题
    useEffect(() => {
        if (displayTitle) {
            Taro.setNavigationBarTitle({title: displayTitle});
        }
    }, [displayTitle]);

    const displayContent = useMemo(() => post ? getPostContent(post) : '', [post]);
    const displayThumbnail = useMemo(() => post ? getPostThumbnail(post) : '', [post]);
    const displaySummary = useMemo(() => post ? getPostSummary(post) : '', [post]);

    const relatedPostsQuery = useMemo(() => {
        if (!post?.categoryIds) return null;
        return {status: 'POST_STATUS_PUBLISHED', id__not: postId, category_ids__in: post.categoryIds};
    }, [post?.categoryIds, postId]);

    useEffect(() => {
        async function loadPost() {
            if (!postId) return;
            setLoading(true);
            try {
                const fetchedPost = await fetchPost(postId);
                if (fetchedPost) setPost(fetchedPost);
            } catch (error) {
                console.error('Load post failed:', error);
            } finally {
                setLoading(false);
            }
        }
        loadPost();
    }, [postId]);

    // 分享功能
    function handleShare() {
        Taro.setClipboardData({
            data: window?.location?.href || '',
        }).then(() => {
            Taro.showToast({title: t('page.post_detail.link_copied'), icon: 'success'});
        }).catch(() => {
            Taro.showToast({title: t('page.post_detail.copy_failed'), icon: 'none'});
        });
    }

    // ===== 加载态 =====
    if (loading) {
        return (
            <View className='min-h-screen w-full bg-pageBg'>
                {/* 头部骨架 */}
                <View className='bg-cardBg'>
                    <View className='h-[360rpx] bg-splitLine/30' />
                    <View className='px-[32rpx] py-[32rpx]'>
                        <Skeleton className='h-[44rpx] w-[80%] rounded-[8rpx] mb-[16rpx]' />
                        <Skeleton className='h-[32rpx] w-[60%] rounded-[8rpx] mb-[24rpx]' />
                        <View className='flex gap-[24rpx]'>
                            <Skeleton className='h-[24rpx] w-[120rpx] rounded-[8rpx]' />
                            <Skeleton className='h-[24rpx] w-[160rpx] rounded-[8rpx]' />
                            <Skeleton className='h-[24rpx] w-[100rpx] rounded-[8rpx]' />
                        </View>
                    </View>
                </View>
                {/* 内容骨架 */}
                <View className='bg-cardBg mt-[16rpx] px-[32rpx] py-[32rpx]'>
                    <Skeleton className='h-[32rpx] w-full rounded-[8rpx] mb-[16rpx]' />
                    <Skeleton className='h-[32rpx] w-[90%] rounded-[8rpx] mb-[16rpx]' />
                    <Skeleton className='h-[32rpx] w-[75%] rounded-[8rpx] mb-[32rpx]' />
                    <Skeleton className='h-[32rpx] w-full rounded-[8rpx] mb-[16rpx]' />
                    <Skeleton className='h-[32rpx] w-[60%] rounded-[8rpx]' />
                </View>
            </View>
        );
    }

    // ===== 空态 =====
    if (!post) {
        return (
            <View className='min-h-screen w-full bg-pageBg flex flex-col items-center justify-center'>
                <View className='w-[120rpx] h-[120rpx] rounded-full bg-pageBg flex items-center justify-center mb-[24rpx]'>
                    <XIcon name='carbon:document-unknown' size={48} className='text-textWeak' />
                </View>
                <Text className='text-body text-textSec mb-[8rpx]'>{t('page.post_detail.post_not_found')}</Text>
                <View
                  className='flex items-center gap-[8rpx] px-[32rpx] py-[16rpx] rounded-full bg-primary/10 mt-[24rpx]'
                  onClick={() => router.reLaunch('/')}
                  hoverClass='tap-active'
                >
                    <XIcon name='carbon:home' size={16} className='text-primary' />
                    <Text className='text-desc text-primary'>{t('page.error.go_home')}</Text>
                </View>
            </View>
        );
    }

    // ===== 正常渲染 =====
    return (
        <View className='min-h-screen w-full bg-pageBg pb-[400rpx]'>
            {/* 封面图 */}
            {displayThumbnail && (
                <View className='w-full h-[400rpx] overflow-hidden'>
                    <Image src={displayThumbnail} mode='aspectFill' className='w-full h-full' />
                </View>
            )}

            {/* 文章头部 - 独立白色卡片 */}
            <View className={`bg-cardBg px-[32rpx] pt-[32rpx] pb-[28rpx] ${displayThumbnail ? '' : 'pt-[40rpx]'}`}>
                {/* 标题 - 全页最大字号 */}
                <Text className='text-title font-bold text-textMain block mb-[16rpx] leading-[1.35]'>
                    {displayTitle}
                </Text>

                {/* 摘要 */}
                {displaySummary && (
                    <Text className='text-desc text-textSec block mb-[24rpx] leading-[1.6]'>
                        {displaySummary}
                    </Text>
                )}

                {/* 作者信息行 */}
                <View className='flex items-center gap-[16rpx] mb-[20rpx]'>
                    <View className='w-[52rpx] h-[52rpx] rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0'>
                        <Text className='text-desc font-bold text-primary'>
                            {post.authorName?.charAt(0) || 'A'}
                        </Text>
                    </View>
                    <View className='flex-1 min-w-0'>
                        <Text className='text-desc font-medium text-textMain block'>{post.authorName}</Text>
                        <Text className='text-tips text-textSec'>{formatDate(post.createdAt)}</Text>
                    </View>
                </View>

                {/* 统计标签 */}
                <View className='flex items-center gap-[16rpx] flex-wrap'>
                    <View className='flex items-center gap-[4rpx] px-[14rpx] py-[4rpx] rounded-full bg-pageBg'>
                        <XIcon name='carbon:view' size={12} className='text-textThird' />
                        <Text className='text-tips text-textSec'>{post.visits || 0} {t('page.post_detail.views')}</Text>
                    </View>
                    <View className='flex items-center gap-[4rpx] px-[14rpx] py-[4rpx] rounded-full bg-pageBg'>
                        <XIcon name='carbon:thumbs-up' size={12} className='text-textThird' />
                        <Text className='text-tips text-textSec'>{post.likes || 0} {t('page.post_detail.likes')}</Text>
                    </View>
                </View>
            </View>

            {/* 文章正文 - 独立白色卡片，确保内容有左右安全边距 */}
            <View className='bg-cardBg px-[32rpx] py-[32rpx] mt-[16rpx]'>
                <ContentViewer content={displayContent} type='markdown' />
            </View>

            {/* 标签 */}
            {post.tagIds && post.tagIds.length > 0 && (
                <View className='bg-cardBg px-[32rpx] py-[20rpx] mt-[16rpx]'>
                    <View className='flex items-center gap-[8rpx] flex-wrap'>
                        <XIcon name='carbon:tag' size={14} className='text-textThird' />
                        {post.tagIds.map((tagId) => (
                            <View
                              key={tagId}
                              className='px-[16rpx] py-[6rpx] rounded-full bg-pageBg'
                              onClick={() => router.push(`/tag/detail?id=${tagId}`)}
                              hoverClass='tap-active'
                            >
                                <Text className='text-tips text-textSec'>#{tagId}</Text>
                            </View>
                        ))}
                    </View>
                </View>
            )}

            {/* 评论区 */}
            <View className='mt-[32rpx] mb-[32rpx]'>
                <CommentSection objectId={postId} contentType='CONTENT_TYPE_POST' onUpdateComments={() => {}} />
            </View>

            {/* 相关文章 */}
            {relatedPostsQuery && (
                <View className='mt-[32rpx]'>
                    {/* 区块标题 */}
                    <View className='flex items-center gap-[8rpx] bg-pageBg px-[32rpx] py-[14rpx]'>
                        <XIcon name='carbon:document' size={16} className='text-primary' />
                        <Text className='text-desc font-bold text-textMain'>
                            {t('page.post_detail.related_posts')}
                        </Text>
                    </View>
                    <View className='px-[32rpx] py-[16rpx]'>
                        <PostList
                          queryParams={relatedPostsQuery}
                          fieldMask='id,status,sort_order,is_featured,visits,likes,comment_count,author_name,available_languages,created_at,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,translations.thumbnail'
                          orderBy={['-sortOrder']}
                          page={1}
                          pageSize={4}
                          compact
                          showSkeleton={false}
                          showPagination={false}
                        />
                    </View>
                </View>
            )}

            {/* 返回按钮 */}
            <View className='mt-[32rpx] px-[32rpx] pb-[48rpx]'>
                <View
                  className='flex items-center justify-center gap-[8rpx] py-[20rpx] rounded-[12rpx] bg-cardBg border-[1rpx] border-splitLine'
                  onClick={() => router.back()}
                  hoverClass='tap-active'
                >
                    <XIcon name='carbon:arrow-left' size={16} className='text-textSec' />
                    <Text className='text-desc text-textSec'>{t('page.post_detail.back')}</Text>
                </View>
            </View>

            {/* 浮动底部操作栏 - 毛玻璃遮罩 */}
            <View
              className='fixed bottom-0 left-0 right-0 z-50 flex items-center justify-around'
              style={{
                backgroundColor: isDark ? 'rgba(35, 35, 38, 0.88)' : 'rgba(255, 255, 255, 0.88)',
                backdropFilter: 'blur(16px)',
                WebkitBackdropFilter: 'blur(16px)',
                borderTop: isDark ? '1rpx solid rgba(255, 255, 255, 0.08)' : '1rpx solid rgba(0, 0, 0, 0.06)',
                paddingTop: '16rpx',
                paddingBottom: 'calc(32rpx + env(safe-area-inset-bottom))',
                paddingLeft: '16rpx',
                paddingRight: '16rpx',
              }}
            >
                <View
                  className='flex items-center gap-[6rpx] px-[32rpx] py-[12rpx] rounded-full'
                  onClick={() => setIsLiked(!isLiked)}
                  hoverClass='tap-active'
                >
                    <XIcon name={isLiked ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'} size={20}
                      className={isLiked ? 'text-primary' : 'text-textSec'}
                    />
                    <Text className={isLiked ? 'text-tips text-primary' : 'text-tips text-textSec'}>
                        {t('page.post_detail.likes')}
                    </Text>
                </View>
                <View
                  className='flex items-center gap-[6rpx] px-[32rpx] py-[12rpx] rounded-full'
                  onClick={() => setIsBookmarked(!isBookmarked)}
                  hoverClass='tap-active'
                >
                    <XIcon name={isBookmarked ? 'carbon:bookmark-filled' : 'carbon:bookmark'} size={20}
                      className={isBookmarked ? 'text-primary' : 'text-textSec'}
                    />
                    <Text className={isBookmarked ? 'text-tips text-primary' : 'text-tips text-textSec'}>
                        {t('page.post_detail.bookmark')}
                    </Text>
                </View>
                <View
                  className='flex items-center gap-[6rpx] px-[32rpx] py-[12rpx] rounded-full'
                  onClick={handleShare}
                  hoverClass='tap-active'
                >
                    <XIcon name='carbon:share' size={20} className='text-textSec' />
                    <Text className='text-tips text-textSec'>{t('page.post_detail.share')}</Text>
                </View>
            </View>

            <BackToTop scrollThreshold={500} />
        </View>
    );
}
