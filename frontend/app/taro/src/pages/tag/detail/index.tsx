import {View, Text} from '@tarojs/components';
import {useMemo, useState} from 'react';
import Taro, {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {fetchTag, getTagTranslation} from '@/api/hooks/tag';
import {fetchListPosts, getPostTitle, getPostSummary, getPostThumbnail} from '@/api/hooks/post';
import PostCard from '@/components/post/PostCard';
import {AppEmpty} from '@/components/ui';
import type {contentservicev1_Tag, contentservicev1_Post} from '@/api/generated/app/service/v1';

const LIST_FIELD_MASK = 'id,status,authorName,publishTime,thumbnail,availableLanguages,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary';

/**
 * 标签详情：标签信息 + 该标签下文章列表
 */
export default function TagDetailPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();

    const tagId = useMemo(() => {
        const raw = Taro.getCurrentInstance()?.router?.params?.id;
        const id = Number(raw);
        return Number.isInteger(id) && id > 0 ? id : 0;
    }, []);

    const [tag, setTag] = useState<contentservicev1_Tag | null>(null);
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [loading, setLoading] = useState(true);

    useDidShow(() => {
        if (!tagId) {
            setLoading(false);
            return;
        }
        setLoading(true);
        fetchTag(tagId)
            .then((loaded) => {
                setTag(loaded);
                const name = getTagTranslation(loaded)?.name;
                if (name) Taro.setNavigationBarTitle({title: name});
            })
            .catch((error) => {
                console.error('Load tag failed:', error);
                setTag(null);
            });

        fetchListPosts({
            paging: {page: 1, pageSize: 30},
            formValues: {status: 'POST_STATUS_PUBLISHED', tag_ids__in: [tagId]},
            fieldMask: LIST_FIELD_MASK,
            orderBy: ['-publishTime'],
        })
            .then((res) => setPosts(res.items || []))
            .catch((error) => {
                console.error('Load tag posts failed:', error);
                setPosts([]);
            })
            .finally(() => setLoading(false));
    });

    if (!tagId) {
        return (
            <View className='au-page w-full bg-pageBg'>
                <View className='au-card mx-[24rpx] mt-[28rpx]'>
                    <AppEmpty variant='error' description={t('page.error.not_found')} />
                </View>
            </View>
        );
    }

    const name = tag ? (getTagTranslation(tag)?.name || tag.code || t('page.tags.tag_untitled')) : '';

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[28rpx]'>
                {/* 标签信息卡 */}
                <View className='au-anim au-card au-card--glow p-[32rpx] mb-[24rpx] flex items-center gap-[20rpx]'>
                    <View
                      className='w-[88rpx] h-[88rpx] rounded-[24rpx] flex items-center justify-center shrink-0'
                      style={{
                          background: tag?.color ? `${tag.color}22` : 'rgba(46, 150, 255, 0.14)',
                      }}
                    >
                        <XIcon name='carbon:hashtag' size={40} style={{color: tag?.color || '#2e96ff'}} />
                    </View>
                    <View className='flex-1 min-w-0'>
                        <Text className='text-card-title font-bold text-textMain block' numberOfLines={1}>
                            {name}
                        </Text>
                        <Text className='text-tips text-textThird block mt-[8rpx]'>
                            {(tag?.postCount ?? 0)} {t('page.categories.articles_count')}
                        </Text>
                    </View>
                </View>

                <Text className='au-section-title block mb-[22rpx] au-anim au-anim-1'>
                    {t('page.categories.all_posts')}
                </Text>
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
                        <AppEmpty description={t('page.posts.no_results')} />
                    </View>
                ) : (
                    <View className='flex flex-col gap-[22rpx]'>
                        {posts.map((post) => (
                            <PostCard
                              key={post.id}
                              title={getPostTitle(post) || t('page.home.untitled')}
                              summary={getPostSummary(post)}
                              thumbnail={getPostThumbnail(post)}
                              author={post.authorName}
                              publishTime={post.publishTime}
                              onClick={() => post.id && router.push(`/post/detail?id=${post.id}`)}
                            />
                        ))}
                    </View>
                )}
            </View>
        </View>
    );
}
