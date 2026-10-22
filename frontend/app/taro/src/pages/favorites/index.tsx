import {View, Text} from '@tarojs/components';
import {useState} from 'react';
import Taro, {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import {listWatchedPosts, unwatchPost} from '@/api/hooks/interaction';
import {getPostTitle, getPostSummary, getPostThumbnail} from '@/api/hooks/post';
import {useIsLogin} from '@/store/core/access/store';
import PostCard from '@/components/post/PostCard';
import {AppEmpty} from '@/components/ui';
import type {contentservicev1_Post} from '@/api/generated/app/service/v1';

/**
 * 我的收藏（tabBar「收藏」）：当前 viewer 收藏的文章列表。
 * 未登录展示登录引导卡；列表支持取消收藏（左滑不便跨端，用卡片右上角按钮）。
 */
export default function FavoritesPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    const isLogin = useIsLogin();
    usePageTitle('page.title.favorites');

    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);

    const loadFavorites = async () => {
        if (!isLogin) return;
        setLoading(true);
        try {
            const res = await listWatchedPosts(1, 50);
            setPosts(res.items || []);
            setTotal(Number(res.total ?? res.items?.length ?? 0));
        } catch (error) {
            console.error('Load favorites failed:', error);
            setPosts([]);
        } finally {
            setLoading(false);
        }
    };

    useDidShow(() => {
        loadFavorites();
    });

    const unwatch = (postId: number) => {
        unwatchPost(postId)
            .then(() => {
                setPosts((prev) => prev.filter((p) => p.id !== postId));
                setTotal((c) => Math.max(0, c - 1));
                Taro.showToast({title: t('page.post_detail.unbookmarked'), icon: 'none'});
            })
            .catch((error) => console.error('Unwatch failed:', error));
    };

    // —— 未登录：登录引导 ——
    if (!isLogin) {
        return (
            <View className='au-page w-full bg-pageBg min-h-screen flex flex-col items-center justify-center px-[48rpx]'>
                <View className='au-card au-card--glow w-full flex flex-col items-center gap-[28rpx] py-[72rpx] px-[40rpx]'>
                    <View className='w-[112rpx] h-[112rpx] rounded-full bg-primary-tint flex items-center justify-center'>
                        <XIcon name='carbon:star' size={52} className='text-primary' />
                    </View>
                    <Text className='text-body text-textThird'>{t('page.user.please_login')}</Text>
                    <View
                      className='btn-hero px-[72rpx] h-[88rpx] rounded-full flex items-center justify-center'
                      hoverClass='tap-active'
                      onClick={() => router.push('/login')}
                    >
                        <Text className='text-body text-white font-bold'>{t('authentication.login.login')}</Text>
                    </View>
                </View>
            </View>
        );
    }

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[36rpx]'>
                <View className='au-anim mb-[28rpx]'>
                    <Text className='au-hero-title text-textMain block'>{t('page.title.favorites')}</Text>
                    <Text className='text-desc text-textThird block mt-[10rpx]'>
                        {total > 0 ? `${total} ${t('page.home.articles_unit')}` : t('page.user.tab_collections')}
                    </Text>
                </View>

                {loading ? (
                    <View className='flex flex-col gap-[22rpx]'>
                        {[0, 1, 2].map((i) => (
                            <View key={i} className='au-card p-[22rpx] flex gap-[22rpx]'>
                                <View className='flex-1 flex flex-col gap-[14rpx]'>
                                    <View className='au-skeleton w-[90%] h-[32rpx] rounded-[8rpx]' />
                                    <View className='au-skeleton w-[50%] h-[24rpx] rounded-[8rpx] mt-auto' />
                                </View>
                                <View className='au-skeleton w-[200rpx] h-[150rpx] rounded-[18rpx]' />
                            </View>
                        ))}
                    </View>
                ) : posts.length === 0 ? (
                    <View className='au-card'>
                        <AppEmpty description={t('page.user.no_collections')} />
                    </View>
                ) : (
                    <View className='flex flex-col gap-[22rpx]'>
                        {posts.map((post) => (
                            <View key={post.id} className='relative'>
                                <PostCard
                                  title={getPostTitle(post) || t('page.home.untitled')}
                                  summary={getPostSummary(post)}
                                  thumbnail={getPostThumbnail(post)}
                                  author={post.authorName}
                                  publishTime={post.publishTime}
                                  onClick={() => post.id && router.push(`/post/detail?id=${post.id}`)}
                                />
                                <View
                                  className='absolute top-[16rpx] right-[16rpx] w-[52rpx] h-[52rpx] rounded-full bg-pageBg flex items-center justify-center z-10'
                                  hoverClass='tap-active'
                                  onClick={() => post.id && unwatch(post.id)}
                                >
                                    <XIcon name='carbon:star-filled' size={24} className='text-primary' />
                                  </View>
                            </View>
                        ))}
                    </View>
                )}
            </View>
        </View>
    );
}
