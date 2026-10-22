import {View, Text} from '@tarojs/components';
import {useMemo, useState} from 'react';
import Taro, {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {fetchCategory, getCategoryName, getCategoryDescription} from '@/api/hooks/category';
import {fetchListPosts, getPostTitle, getPostSummary, getPostThumbnail} from '@/api/hooks/post';
import PostCard from '@/components/post/PostCard';
import {AppEmpty} from '@/components/ui';
import type {contentservicev1_Category, contentservicev1_Post} from '@/api/generated/app/service/v1';

const LIST_FIELD_MASK = 'id,status,authorName,publishTime,thumbnail,availableLanguages,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary';

/**
 * 分类详情：分类信息卡 + 子分类 chips + 该分类下文章列表
 */
export default function CategoryDetailPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();

    const categoryId = useMemo(() => {
        const raw = Taro.getCurrentInstance()?.router?.params?.id;
        const id = Number(raw);
        return Number.isInteger(id) && id > 0 ? id : 0;
    }, []);

    const [detail, setDetail] = useState<contentservicev1_Category | null>(null);
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [loading, setLoading] = useState(true);

    useDidShow(() => {
        if (!categoryId) {
            setLoading(false);
            return;
        }
        setLoading(true);
        fetchCategory({id: categoryId})
            .then((loaded) => {
                setDetail(loaded);
                const name = getCategoryName(loaded);
                if (name) Taro.setNavigationBarTitle({title: name});
            })
            .catch((error) => {
                console.error('Load category failed:', error);
                setDetail(null);
            });

        fetchListPosts({
            paging: {page: 1, pageSize: 20},
            formValues: {status: 'POST_STATUS_PUBLISHED', category_ids__in: [categoryId]},
            fieldMask: LIST_FIELD_MASK,
            orderBy: ['-publishTime'],
        })
            .then((res) => setPosts(res.items || []))
            .catch((error) => {
                console.error('Load category posts failed:', error);
                setPosts([]);
            })
            .finally(() => setLoading(false));
    });

    if (!categoryId) {
        return (
            <View className='au-page w-full bg-pageBg'>
                <View className='au-card mx-[24rpx] mt-[28rpx]'>
                    <AppEmpty variant='error' description={t('page.error.not_found')} />
                </View>
            </View>
        );
    }

    const name = detail ? (getCategoryName(detail) || t('page.categories.category_untitled')) : '';
    const children = detail?.children ?? [];

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[28rpx]'>
                {/* 分类信息卡 */}
                <View className='au-anim au-card au-card--glow p-[32rpx] mb-[24rpx]'>
                    <View className='flex items-center gap-[20rpx]'>
                        <View className='w-[88rpx] h-[88rpx] rounded-[24rpx] bg-primary-tint flex items-center justify-center shrink-0'>
                            <XIcon name={detail?.icon || 'carbon:folder'} size={40} className='text-primary' />
                        </View>
                        <View className='flex-1 min-w-0'>
                            <Text className='text-card-title font-bold text-textMain block' numberOfLines={1}>
                                {name}
                            </Text>
                            <Text className='text-tips text-textThird block mt-[8rpx]'>
                                {(detail?.postCount ?? 0)} {t('page.categories.articles_count')}
                            </Text>
                        </View>
                    </View>
                    {detail && getCategoryDescription(detail) && (
                        <Text className='text-desc text-textSec leading-[1.7] block mt-[20rpx]'>
                            {getCategoryDescription(detail)}
                        </Text>
                    )}
                </View>

                {/* 子分类 chips */}
                {children.length > 0 && (
                    <View className='au-anim au-anim-1 flex flex-wrap gap-[14rpx] mb-[24rpx]'>
                        {children.map((child) => (
                            <View
                              key={child.id}
                              className='au-chip'
                              hoverClass='tap-active'
                              onClick={() => child.id && router.push(`/category/detail?id=${child.id}`)}
                            >
                                <XIcon name='carbon:folder' size={18} />
                                <Text numberOfLines={1}>{getCategoryName(child)}</Text>
                                <Text className='opacity-70'>({child.postCount ?? 0})</Text>
                            </View>
                        ))}
                    </View>
                )}

                {/* 文章列表 */}
                <Text className='au-section-title block mb-[22rpx] au-anim au-anim-2'>
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
                        <AppEmpty description={t('page.categories.no_categories')} />
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
