import {View, Text, Swiper, SwiperItem} from '@tarojs/components';
import {useState, useEffect} from 'react';
import Taro from '@tarojs/taro';

import {fetchListPosts,getPostTitle, getPostThumbnail} from '@/api/hooks/post';

import type {contentservicev1_Post} from '@/api/generated/app/service/v1';
import Image from '@/components/ui/image';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {XIcon} from '@/plugins/xicon';
import {useTranslations} from '@/lib/next-intl-compat';


/**
 * 首页「最新文章」区块
 * - 走马灯轮播，自动切换
 * - 封面图 + 底部叠加标题
 */
export default function LatestPostsSection() {
  const t = useTranslations('page.home');
  const router = useI18nRouter();
  const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
  const [loading, setLoading] = useState(true);

  // 获取文章列表
  useEffect(() => {
    const fetchPosts = async () => {
      try {
        setLoading(true);
        const res = await fetchListPosts({
          paging: {page: 1, pageSize: 5},
          formValues: {status: 'POST_STATUS_PUBLISHED', isFeatured: false},
          fieldMask: 'id,status,sortOrder,isFeatured,visits,likes,commentCount,authorName,availableLanguages,createdAt,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary,translations.thumbnail',
          orderBy: ['-createdAt'],
        });
        setPosts(res.items || []);
      } catch (error) {
        console.error('Failed to fetch posts:', error);
        setPosts([]);
      } finally {
        setLoading(false);
      }
    };
    fetchPosts();
  }, []);

  const handleViewPost = (post: contentservicev1_Post) => {
    router.push(`/post/detail?id=${post.id}&from=home`);
    Taro.pageScrollTo({scrollTop: 0, duration: 300});
  };

  // 加载骨架屏
  if (loading) {
    return (
      <View className='w-full'>
        <View className='flex items-center justify-between mb-[24rpx]'>
          <View className='flex items-center gap-[8rpx]'>
            <XIcon name='carbon:document' size={20} className='text-primary' />
            <Text className='text-card-title font-bold text-textMain'>
              {t('latest_posts')}
            </Text>
          </View>
        </View>
        {/* 骨架屏 */}
        <View className='rounded-[16rpx] bg-cardBg overflow-hidden'>
          <View className='w-full h-[160rpx] bg-pageBg' />
          <View className='p-[20rpx] flex flex-col gap-[12rpx]'>
            <View className='h-[28rpx] w-3/4 bg-pageBg rounded' />
            <View className='h-[24rpx] w-full bg-pageBg rounded' />
            <View className='h-[24rpx] w-1/2 bg-pageBg rounded' />
          </View>
        </View>
      </View>
    );
  }

  // 空状态
  if (posts.length === 0) {
    return null;
  }

  return (
    <View className='w-full'>
      {/* 标题行 */}
      <View className='flex items-center justify-between mb-[32rpx]'>
        <View className='flex items-center gap-[8rpx]'>
          <XIcon name='carbon:document' size={20} className='text-primary' />
          <Text className='text-card-title font-bold text-textMain'>
            {t('latest_posts')}
          </Text>
        </View>
        <View
          className='px-[24rpx] py-[12rpx] rounded-full flex items-center justify-center'
          style={{
            backgroundColor: 'rgba(22,119,255,0.08)',
          }}
          onClick={() => router.push('/post')}
          hoverClass='tap-active'
        >
          <Text className='text-tips font-medium text-primary'>{t('view_all')} →</Text>
        </View>
      </View>

      {/* 走马灯轮播 - 只显示封面图 */}
      <Swiper
        className='rounded-[16rpx] overflow-hidden'
        style={{height: '360rpx'}}
        autoplay
        interval={4000}
        duration={500}
        circular
        indicatorDots
        indicatorColor='rgba(22,119,255,0.2)'
        indicatorActiveColor='rgba(22,119,255,1)'
        onChange={() => {}}
      >
        {posts.map((post) => {
          const thumbnail = getPostThumbnail(post);
          const title = getPostTitle(post);

          return (
            <SwiperItem key={post.id}>
              <View
                className='w-full h-full flex items-center justify-center'
                onClick={() => handleViewPost(post)}
                hoverClass='tap-active'
                style={{
                  background: thumbnail ? 'none' : 'linear-gradient(135deg, rgba(22,119,255,0.06) 0%, rgba(114,46,209,0.06) 100%)',
                }}
              >
                {/* 封面图 */}
                {thumbnail ? (
                  <Image src={thumbnail} mode='aspectFill' className='w-full h-full' />
                ) : (
                  <XIcon name='carbon:document' size={80} className='text-primary' style={{opacity: 0.4}} />
                )}

                {/* 底部标题叠加层 - 增强渐变遮罩，提升文字可读性 */}
                {title && title !== 'title' && (
                  <View
                    className='absolute bottom-0 left-0 right-0 p-[24rpx]'
                    style={{
                      background: 'linear-gradient(to top, rgba(0,0,0,0.9) 0%, rgba(0,0,0,0.7) 40%, rgba(0,0,0,0.3) 70%, rgba(0,0,0,0) 100%)',
                    }}
                  >
                    <Text
                      className='text-body font-bold text-white leading-[1.4]'
                      numberOfLines={2}
                      style={{
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        display: '-webkit-box',
                        WebkitLineClamp: 2,
                        WebkitBoxOrient: 'vertical',
                      }}
                    >
                      {title}
                    </Text>
                  </View>
                )}
              </View>
            </SwiperItem>
          );
        })}
      </Swiper>
    </View>
  );
}
