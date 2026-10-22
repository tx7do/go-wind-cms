import {View, Text, Image, ScrollView} from '@tarojs/components';
import {useState} from 'react';
import {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import {fetchListPosts, getPostTitle, getPostSummary, getPostThumbnail} from '@/api/hooks/post';
import {fetchListCategories, getCategoryName} from '@/api/hooks/category';
import {fetchListTags, getTagTranslation} from '@/api/hooks/tag';
import PostCard from '@/components/post/PostCard';
import type {contentservicev1_Post, contentservicev1_Category, contentservicev1_Tag} from '@/api/generated/app/service/v1';

// 分类图标/配色轮换：让网格卡片有节奏感（主题无关的柔和底色）
const CATEGORY_ICONS = ['carbon:document-multiple-01', 'carbon:carbon', 'carbon:globe', 'carbon:cloud-foundry-2', 'carbon:development', 'carbon:data-base', 'carbon:security', 'carbon:analytics'];
const CATEGORY_TINTS = [
    'rgba(46, 150, 255, 0.14)',
    'rgba(34, 211, 238, 0.14)',
    'rgba(139, 92, 246, 0.14)',
    'rgba(255, 159, 75, 0.16)',
];
const CATEGORY_ICON_COLORS = ['#2e96ff', '#22c3d8', '#8b5cf6', '#ff9f4b'];

const LIST_FIELD_MASK = 'id,status,sortOrder,isFeatured,authorName,publishTime,thumbnail,availableLanguages,translations.id,translations.postId,translations.languageCode,translations.title,translations.summary';

/**
 * 首页
 * 布局：搜索入口 → 阅读 hero 入口 → 内容分类网格 → 推荐阅读（横滑） → 最新文章 → 热门标签
 *
 * useDidShow：每次页面显示（包括从子页面/切 tab 返回）时刷新数据。
 */
export default function Home() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    usePageTitle('page.title.home');

    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [featured, setFeatured] = useState<contentservicev1_Post[]>([]);
    const [categories, setCategories] = useState<contentservicev1_Category[]>([]);
    const [tags, setTags] = useState<contentservicev1_Tag[]>([]);
    const [loading, setLoading] = useState(true);

    useDidShow(() => {
        (async () => {
            setLoading(true);
            try {
                const [latest, cats] = await Promise.all([
                    fetchListPosts({
                        paging: {page: 1, pageSize: 5},
                        formValues: {status: 'POST_STATUS_PUBLISHED'},
                        fieldMask: LIST_FIELD_MASK,
                        orderBy: ['-publishTime'],
                    }),
                    fetchListCategories({paging: {page: 1, pageSize: 6}}),
                ]);
                setPosts(latest.items || []);
                setCategories(cats.items || []);
            } catch (error) {
                console.error('Load home feed failed:', error);
                setPosts([]);
                setCategories([]);
            } finally {
                setLoading(false);
            }
        })();

        fetchListPosts({
            paging: {page: 1, pageSize: 6},
            formValues: {status: 'POST_STATUS_PUBLISHED', isFeatured: true},
            fieldMask: LIST_FIELD_MASK,
            orderBy: ['-publishTime'],
        })
            .then((res) => setFeatured(res.items || []))
            .catch((error) => {
                console.error('Load featured posts failed:', error);
                setFeatured([]);
            });

        fetchListTags({paging: {page: 1, pageSize: 12}})
            .then((res) => setTags(res.items || []))
            .catch((error) => {
                console.error('Load tags failed:', error);
                setTags([]);
            });
    });

    function openPost(post: contentservicev1_Post) {
        if (post.id) router.push(`/post/detail?id=${post.id}`);
    }

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[28rpx]'>
                {/* 搜索入口 */}
                <View
                  className='au-anim au-input h-[84rpx] mb-[28rpx]'
                  hoverClass='tap-active'
                  onClick={() => router.push('/search')}
                >
                    <XIcon name='carbon:search' size={24} className='text-textThird' />
                    <Text className='au-input-field text-left text-desc text-textThird leading-[84rpx]'>
                        {t('page.posts.search_placeholder')}
                    </Text>
                </View>

                {/* 阅读 hero 入口 */}
                <View
                  className='au-anim relative overflow-hidden rounded-[32rpx] p-[40rpx] flex items-center justify-between'
                  style={{
                      background: 'linear-gradient(128deg, #57a3ff 0%, #1677ff 52%, #0b55d4 100%)',
                      boxShadow: '0 20rpx 56rpx rgba(12, 110, 255, 0.38)',
                  }}
                  hoverClass='tap-active'
                  onClick={() => router.push('/posts')}
                >
                    {/* 卡内极光装饰 */}
                    <View
                      className='absolute rounded-full'
                      style={{
                          width: '300rpx',
                          height: '300rpx',
                          right: '-90rpx',
                          top: '-120rpx',
                          background: 'radial-gradient(circle, rgba(255,255,255,0.28), transparent 68%)',
                      }}
                    />
                    <View className='relative flex items-center gap-[28rpx]'>
                        <View className='w-[104rpx] h-[104rpx] rounded-[28rpx] bg-white/16 flex items-center justify-center border-[length:1rpx] border-white/25'>
                            <XIcon name='carbon:book' size={48} className='text-white' />
                        </View>
                        <View>
                            <Text className='au-hero-title text-white block'>
                                {t('page.home.browse_posts')}
                            </Text>
                            <Text className='text-desc text-white/85 block mt-[8rpx]'>
                                {t('page.home.hero_description')}
                            </Text>
                        </View>
                    </View>
                    <XIcon name='carbon:chevron-right' size={24} className='relative text-white/85' />
                </View>

                {/* 内容分类 */}
                <View className='au-anim au-anim-1 flex items-center justify-between mt-[44rpx] mb-[22rpx]'>
                    <Text className='au-section-title'>{t('page.home.categories')}</Text>
                    <View
                      className='flex items-center gap-[4rpx] min-h-touch px-[16rpx] -mr-[16rpx]'
                      hoverClass='tap-active'
                      onClick={() => router.push('/categories')}
                    >
                        <Text className='text-desc text-primary font-medium'>
                            {t('page.home.view_all')}
                        </Text>
                        <XIcon name='carbon:chevron-right' size={16} className='text-primary' />
                    </View>
                </View>

                {loading ? (
                    <View className='grid grid-cols-2 gap-[22rpx]'>
                        {[0, 1, 2, 3].map((i) => (
                            <View key={i} className='au-card p-[26rpx]'>
                                <View className='au-skeleton w-[64rpx] h-[64rpx] rounded-[18rpx] mb-[16rpx]' />
                                <View className='au-skeleton w-[70%] h-[30rpx] rounded-[8rpx] mb-[12rpx]' />
                                <View className='au-skeleton w-[45%] h-[22rpx] rounded-[8rpx]' />
                            </View>
                        ))}
                    </View>
                ) : categories.length === 0 ? (
                    <View className='au-card py-[56rpx] flex flex-col items-center'>
                        <View className='w-[96rpx] h-[96rpx] rounded-[28rpx] bg-primary-tint flex items-center justify-center mb-[20rpx]'>
                            <XIcon name='carbon:folder' size={44} className='text-primary' />
                        </View>
                        <Text className='text-body text-textThird'>
                            {t('page.categories.no_categories')}
                        </Text>
                    </View>
                ) : (
                    <View className='grid grid-cols-2 gap-[22rpx]'>
                        {categories.slice(0, 6).map((category, idx) => (
                            <View
                              key={category.id}
                              className='au-card au-card--glow p-[26rpx] relative overflow-hidden flex flex-col'
                              hoverClass='tap-active'
                              onClick={() => category.id && router.push(`/category/detail?id=${category.id}`)}
                            >
                                <View
                                  className='w-[68rpx] h-[68rpx] rounded-[20rpx] flex items-center justify-center mb-[16rpx]'
                                  style={{background: CATEGORY_TINTS[idx % CATEGORY_TINTS.length]}}
                                >
                                    <XIcon
                                      name={category.icon || CATEGORY_ICONS[idx % CATEGORY_ICONS.length]}
                                      size={32}
                                      style={{color: CATEGORY_ICON_COLORS[idx % CATEGORY_ICON_COLORS.length]}}
                                    />
                                </View>
                                <Text className='text-body font-bold text-textMain block' numberOfLines={1}>
                                    {getCategoryName(category) || t('page.home.category_default')}
                                </Text>
                                <View className='flex items-center justify-between mt-auto pt-[18rpx]'>
                                    <Text className='text-tips text-textThird'>
                                        {(category.postCount ?? 0)} {t('page.home.articles_unit')}
                                    </Text>
                                    <XIcon name='carbon:arrow-right' size={16} className='text-primary' />
                                </View>
                            </View>
                        ))}
                    </View>
                )}

                {/* 推荐阅读（横滑） */}
                {featured.length > 0 && (
                    <>
                        <View className='au-anim au-anim-2 flex items-center justify-between mt-[44rpx] mb-[22rpx]'>
                            <Text className='au-section-title'>{t('page.home.featured_reads')}</Text>
                        </View>
                        <View className='-mx-[24rpx] au-anim au-anim-2'>
                            <ScrollView scrollX className='w-full' showScrollbar={false}>
                                <View className='flex gap-[20rpx] px-[24rpx]'>
                                    {featured.map((post) => (
                                        <View
                                          key={post.id}
                                          className='au-card au-card--glow w-[460rpx] shrink-0 overflow-hidden'
                                          hoverClass='tap-active'
                                          onClick={() => openPost(post)}
                                        >
                                            <View className='h-[220rpx]'>
                                                {getPostThumbnail(post) ? (
                                                    <Image
                                                      src={getPostThumbnail(post)}
                                                      mode='aspectFill'
                                                      className='w-full h-full'
                                                      lazyLoad
                                                    />
                                                ) : (
                                                    <View
                                                      className='w-full h-full flex items-center justify-center'
                                                      style={{background: 'linear-gradient(128deg, #57a3ff 0%, #1677ff 52%, #0b55d4 100%)'}}
                                                    >
                                                        <XIcon name='carbon:star' size={44} className='text-white/70' />
                                                    </View>
                                                )}
                                            </View>
                                            <View className='p-[24rpx]'>
                                                <Text className='text-body font-bold text-textMain leading-[1.4]' numberOfLines={2}>
                                                    {getPostTitle(post) || t('page.home.untitled')}
                                                </Text>
                                                {getPostSummary(post) && (
                                                    <Text className='text-tips text-textThird mt-[10rpx] block' numberOfLines={1}>
                                                        {getPostSummary(post)}
                                                    </Text>
                                                )}
                                            </View>
                                        </View>
                                    ))}
                                </View>
                            </ScrollView>
                        </View>
                    </>
                )}

                {/* 最新文章 */}
                <View className='au-anim au-anim-3 flex items-center justify-between mt-[44rpx] mb-[22rpx]'>
                    <Text className='au-section-title'>{t('page.home.latest_posts')}</Text>
                    <View
                      className='flex items-center gap-[4rpx] min-h-touch px-[16rpx] -mr-[16rpx]'
                      hoverClass='tap-active'
                      onClick={() => router.push('/posts')}
                    >
                        <Text className='text-desc text-primary font-medium'>
                            {t('page.home.view_all')}
                        </Text>
                        <XIcon name='carbon:chevron-right' size={16} className='text-primary' />
                    </View>
                </View>
                <View className='flex flex-col gap-[22rpx]'>
                    {posts.map((post) => (
                        <PostCard
                          key={post.id}
                          title={getPostTitle(post) || t('page.home.untitled')}
                          summary={getPostSummary(post)}
                          thumbnail={getPostThumbnail(post)}
                          author={post.authorName}
                          publishTime={post.publishTime}
                          className='au-anim au-anim-4'
                          onClick={() => openPost(post)}
                        />
                    ))}
                </View>

                {/* 热门标签 */}
                {tags.length > 0 && (
                    <>
                        <View className='au-anim au-anim-5 flex items-center justify-between mt-[44rpx] mb-[22rpx]'>
                            <Text className='au-section-title'>{t('page.home.popular_tags')}</Text>
                            <View
                              className='flex items-center gap-[4rpx] min-h-touch px-[16rpx] -mr-[16rpx]'
                              hoverClass='tap-active'
                              onClick={() => router.push('/tags')}
                            >
                                <Text className='text-desc text-primary font-medium'>
                                    {t('page.home.view_all')}
                                </Text>
                                <XIcon name='carbon:chevron-right' size={16} className='text-primary' />
                            </View>
                        </View>
                        <View className='flex flex-wrap gap-[16rpx]'>
                            {tags.map((tag) => (
                                <View
                                  key={tag.id}
                                  className='au-chip'
                                  hoverClass='tap-active'
                                  onClick={() => tag.id && router.push(`/tag/detail?id=${tag.id}`)}
                                >
                                    <XIcon name='carbon:hashtag' size={18} />
                                    <Text>{getTagTranslation(tag)?.name || tag.code || ''}</Text>
                                    <Text className='opacity-70'>({tag.postCount ?? 0})</Text>
                                </View>
                            ))}
                        </View>
                    </>
                )}
            </View>
        </View>
    );
}
