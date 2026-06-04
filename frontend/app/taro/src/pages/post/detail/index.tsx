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
import {contentservicev1_Post} from '@/api/generated/app/service/v1';
import {fetchPost, getPostTitle, getPostContent, getPostThumbnail} from '@/api/hooks/post';

export default function PostDetailPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
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
    const displayContent = useMemo(() => post ? getPostContent(post) : '', [post]);
    const displayThumbnail = useMemo(() => post ? getPostThumbnail(post) : '', [post]);

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

    if (loading) {
        return (
            <View className='w-full bg-pageBg min-h-screen'>
                <View className='bg-cardBg px-[24rpx] py-[32rpx]'>
                    <View className='h-[48rpx] w-[60%] bg-splitLine rounded mb-[24rpx]' />
                    <View className='h-[32rpx] w-[40%] bg-splitLine rounded' />
                </View>
                <View className='px-[24rpx] py-[32rpx]'>
                    <View className='h-[400rpx] bg-splitLine rounded mb-[24rpx]' />
                    <View className='h-[32rpx] w-[80%] bg-splitLine rounded mb-[16rpx]' />
                    <View className='h-[32rpx] w-[60%] bg-splitLine rounded mb-[16rpx]' />
                    <View className='h-[32rpx] w-[70%] bg-splitLine rounded' />
                </View>
            </View>
        );
    }

    if (!post) {
        return (
            <View className='w-full bg-pageBg min-h-screen flex items-center justify-center'>
                <XIcon name='carbon:warning' size={28} className='text-textWeak' />
                <Text className='text-body text-textThird mt-[16rpx]'>{t('page.post_detail.not_found') || 'Post not found'}</Text>
            </View>
        );
    }

    return (
        <View className='w-full bg-pageBg'>
            {/* 返回导航 */}
            <View className='bg-cardBg px-[24rpx]'>
                <View
                  className='flex items-center gap-[8rpx] h-[88rpx] min-h-touch'
                  onClick={() => router.back()}
                  hoverClass='tap-active'
                >
                    <XIcon name='carbon:arrow-left' size={16} className='text-textSec' />
                    <Text className='text-desc text-textSec'>{t('page.post_detail.back')}</Text>
                </View>
            </View>

            {/* 文章封面 */}
            {displayThumbnail && (
                <View className='w-full h-[360rpx]'>
                    <Image src={displayThumbnail} mode='aspectFill' className='w-full h-full' />
                </View>
            )}

            {/* 文章头部 */}
            <View className='bg-cardBg px-[24rpx] py-[32rpx]'>
                <Text className='text-title font-bold text-textMain block mb-[16rpx]'>{displayTitle}</Text>
                <View className='flex flex-wrap gap-[24rpx]'>
                    <View className='flex items-center gap-[8rpx]'>
                        <XIcon name='carbon:user' size={14} className='text-textThird' />
                        <Text className='text-desc text-textSec'>{post.authorName}</Text>
                    </View>
                    <View className='flex items-center gap-[8rpx]'>
                        <XIcon name='carbon:calendar' size={14} className='text-textThird' />
                        <Text className='text-desc text-textSec'>{formatDate(post.createdAt)}</Text>
                    </View>
                    <View className='flex items-center gap-[8rpx]'>
                        <XIcon name='carbon:view' size={14} className='text-textThird' />
                        <Text className='text-desc text-textSec'>{post.visits || 0}</Text>
                    </View>
                    <View className='flex items-center gap-[8rpx]'>
                        <XIcon name='carbon:thumb-up' size={14} className='text-textThird' />
                        <Text className='text-desc text-textSec'>{post.likes || 0}</Text>
                    </View>
                </View>
            </View>

            {/* 文章内容 */}
            <View className='bg-cardBg px-[24rpx] py-[24rpx] mt-[16rpx]'>
                <ContentViewer content={displayContent} type='markdown' />
            </View>

            {/* 操作按钮 */}
            <View className='flex justify-center gap-[48rpx] py-[32rpx] bg-cardBg mt-[16rpx]'>
                <View
                  className='flex flex-col items-center gap-[8rpx] min-w-touch min-h-touch justify-center'
                  onClick={() => setIsLiked(!isLiked)}
                  hoverClass='tap-active'
                >
                    <XIcon name={isLiked ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'} size={24}
                      className={isLiked ? 'text-primary' : 'text-textThird'}
                    />
                    <Text className='text-tips text-textSec'>{t('page.post_detail.likes')}</Text>
                </View>
                <View
                  className='flex flex-col items-center gap-[8rpx] min-w-touch min-h-touch justify-center'
                  onClick={() => setIsBookmarked(!isBookmarked)}
                  hoverClass='tap-active'
                >
                    <XIcon name={isBookmarked ? 'carbon:bookmark-filled' : 'carbon:bookmark'} size={24}
                      className={isBookmarked ? 'text-primary' : 'text-textThird'}
                    />
                    <Text className='text-tips text-textSec'>{t('page.post_detail.bookmark')}</Text>
                </View>
                <View
                  className='flex flex-col items-center gap-[8rpx] min-w-touch min-h-touch justify-center'
                  hoverClass='tap-active'
                >
                    <XIcon name='carbon:share' size={24} className='text-textThird' />
                    <Text className='text-tips text-textSec'>{t('page.post_detail.share')}</Text>
                </View>
            </View>

            {/* 评论区 */}
            <CommentSection objectId={postId} contentType='CONTENT_TYPE_POST' onUpdateComments={() => {}} />

            {/* 相关文章 */}
            {relatedPostsQuery && (
                <View className='px-[24rpx] py-[32rpx]'>
                    <Text className='text-card-title font-bold text-textMain block mb-[24rpx]'>
                        {t('page.post_detail.related_posts')}
                    </Text>
                    <PostList
                      queryParams={relatedPostsQuery}
                      fieldMask='id,status,sort_order,is_featured,visits,likes,comment_count,author_name,available_languages,created_at,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,translations.thumbnail'
                      orderBy={['-sortOrder']}
                      page={1}
                      pageSize={3}
                      columns={1}
                      showSkeleton={false}
                      showPagination={false}
                    />
                </View>
            )}

            <BackToTop scrollThreshold={500} />
        </View>
    );
}
