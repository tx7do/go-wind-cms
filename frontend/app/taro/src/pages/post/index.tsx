import {View, Text, ScrollView} from '@tarojs/components';
import {useState} from 'react';
import {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import {fetchListPosts, getPostTitle, getPostSummary, getPostThumbnail} from '@/api/hooks/post';
import {fetchListCategories, getCategoryName} from '@/api/hooks/category';
import PostCard from '@/components/post/PostCard';
import {AppEmpty} from '@/components/ui';
import type {contentservicev1_Post, contentservicev1_Category} from '@/api/generated/app/service/v1';

const LIST_FIELD_MASK = 'id,status,sortOrder,isFeatured,authorName,publishTime,thumbnail,availableLanguages,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary';
const PAGE_SIZE = 10;

/**
 * 文章列表（tabBar「文章」）
 * 布局：搜索入口 → 分类筛选横滑 chips → 文章分页列表（加载更多）
 */
export default function PostListPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    usePageTitle('page.title.posts');

    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [total, setTotal] = useState(0);
    const [categories, setCategories] = useState<contentservicev1_Category[]>([]);
    const [selectedCategory, setSelectedCategory] = useState<number | null>(null);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);

    const loadPage = async (targetPage: number, categoryId: number | null, append: boolean) => {
        append ? setLoadingMore(true) : setLoading(true);
        try {
            const res = await fetchListPosts({
                paging: {page: targetPage, pageSize: PAGE_SIZE},
                formValues: {
                    status: 'POST_STATUS_PUBLISHED',
                    ...(categoryId ? {category_ids__in: [categoryId]} : {}),
                },
                fieldMask: LIST_FIELD_MASK,
                orderBy: ['-publishTime'],
            });
            const items = res.items || [];
            setPosts((prev) => (append ? [...prev, ...items] : items));
            setTotal(Number(res.total ?? items.length));
            setPage(targetPage);
        } catch (error) {
            console.error('Load posts failed:', error);
            if (!append) setPosts([]);
        } finally {
            append ? setLoadingMore(false) : setLoading(false);
        }
    };

    useDidShow(() => {
        // 列表是 tab 常驻页：每次显示回到第一页刷新（新发布的文章要露出）
        loadPage(1, selectedCategory, false);
        // 分类筛选 chips 缺失时补拉（tab 页常驻，正常只拉一次）
        if (categories.length === 0) {
            fetchListCategories({paging: {page: 1, pageSize: 30}})
                .then((res) => setCategories(res.items || []))
                .catch((error) => {
                    console.error('Load categories failed:', error);
                    setCategories([]);
                });
        }
    });

    const switchCategory = (categoryId: number | null) => {
        setSelectedCategory(categoryId);
        loadPage(1, categoryId, false);
    };

    const hasMore = posts.length < total;

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[28rpx]'>
                {/* 搜索入口 */}
                <View
                  className='au-anim au-input h-[84rpx] mb-[24rpx]'
                  hoverClass='tap-active'
                  onClick={() => router.push('/search')}
                >
                    <XIcon name='carbon:search' size={24} className='text-textThird' />
                    <Text className='au-input-field text-left text-desc text-textThird leading-[84rpx]'>
                        {t('page.posts.search_placeholder')}
                    </Text>
                </View>

                {/* 分类筛选 chips */}
                {/* shrink-0 必须加：横滑 ScrollView 溢出时 flex 子项默认可收缩，
                    中文 min-content 宽度 = 单字，整排 chip 会被压成圆形文字竖排 */}
                <View className='-mx-[24rpx] mb-[24rpx]'>
                    <ScrollView scrollX className='w-full' showScrollbar={false}>
                        <View className='flex gap-[14rpx] px-[24rpx]'>
                            <View
                              className={`au-chip shrink-0 ${selectedCategory === null ? 'ring-[length:2rpx] ring-primary' : 'opacity-80'}`}
                              hoverClass='tap-active'
                              onClick={() => switchCategory(null)}
                            >
                                <Text className='whitespace-nowrap'>{t('page.posts.all_categories')}</Text>
                            </View>
                            {categories.map((cat) => (
                                <View
                                  key={cat.id}
                                  className={`au-chip shrink-0 ${selectedCategory === cat.id ? 'ring-[length:2rpx] ring-primary' : 'opacity-80'}`}
                                  hoverClass='tap-active'
                                  onClick={() => cat.id && switchCategory(cat.id)}
                                >
                                    <Text className='whitespace-nowrap' numberOfLines={1}>
                                        {getCategoryName(cat) || t('page.categories.category_untitled')}
                                    </Text>
                                </View>
                            ))}
                        </View>
                    </ScrollView>
                </View>

                {/* 统计 */}
                {!loading && total > 0 && (
                    <Text className='text-tips text-textThird block mb-[20rpx]'>
                        {t('page.posts.found')} {total} {t('page.posts.articles')}
                    </Text>
                )}

                {/* 文章列表 */}
                {loading ? (
                    <View className='flex flex-col gap-[22rpx]'>
                        {[0, 1, 2, 3].map((i) => (
                            <View key={i} className='au-card p-[22rpx] flex gap-[22rpx]'>
                                <View className='flex-1 flex flex-col gap-[14rpx]'>
                                    <View className='au-skeleton w-[90%] h-[32rpx] rounded-[8rpx]' />
                                    <View className='au-skeleton w-[60%] h-[24rpx] rounded-[8rpx]' />
                                    <View className='au-skeleton w-[40%] h-[22rpx] rounded-[8rpx] mt-auto' />
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

                        {/* 加载更多 */}
                        {hasMore && (
                            <View
                              className='au-card h-[92rpx] flex items-center justify-center mt-[8rpx]'
                              hoverClass='tap-active'
                              onClick={() => !loadingMore && loadPage(page + 1, selectedCategory, true)}
                            >
                                <Text className='text-desc text-primary font-medium'>
                                    {loadingMore ? t('common.loading') : t('comment.load_more')}
                                </Text>
                            </View>
                        )}
                        {!hasMore && posts.length > PAGE_SIZE && (
                            <Text className='text-center text-tips text-textWeak py-[20rpx]'>
                                {t('comment.no_more')}
                            </Text>
                        )}
                    </View>
                )}
            </View>
        </View>
    );
}
