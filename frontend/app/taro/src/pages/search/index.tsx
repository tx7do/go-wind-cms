import {View, Text, Input} from '@tarojs/components';
import {useRef, useState} from 'react';
import {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import {fetchSearchPosts, fetchPost, getPostSummary, getPostThumbnail} from '@/api/hooks/post';
import PostCard from '@/components/post/PostCard';
import {AppEmpty} from '@/components/ui';
import type {contentservicev1_Post, contentservicev1_SearchPostHit} from '@/api/generated/app/service/v1';

/**
 * 搜索页：输入即时搜索（防抖 400ms）→ 命中标题列表 → 点击进入详情。
 * SearchPosts 仅返回最小字段集（postId/title），点卡片时用 fetchPost 补全卡片信息。
 */
export default function SearchPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    usePageTitle('page.title.posts');

    const [keyword, setKeyword] = useState('');
    const [hits, setHits] = useState<contentservicev1_SearchPostHit[]>([]);
    const [total, setTotal] = useState(0);
    const [details, setDetails] = useState<Record<number, contentservicev1_Post>>({});
    const [loading, setLoading] = useState(false);
    const [searched, setSearched] = useState(false);

    const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const runSearch = async (query: string) => {
        const q = query.trim();
        if (!q) {
            setHits([]);
            setTotal(0);
            setSearched(false);
            return;
        }
        setLoading(true);
        try {
            const res = await fetchSearchPosts({query: q, page: 0, pageSize: 20});
            const items = res.items || [];
            setHits(items);
            setTotal(Number(res.total ?? items.length));
            setSearched(true);

            // 拉取命中文章的封面/摘要，卡片信息完整（失败不阻塞列表）
            const need = items.filter((it) => it.postId && !details[it.postId]).slice(0, 20);
            await Promise.all(
                need.map((it) =>
                    fetchPost(it.postId as number)
                        .then((post) => {
                            if (post?.id) {
                                setDetails((prev) => ({...prev, [post.id as number]: post}));
                            }
                        })
                        .catch((error) => console.error('Load search hit detail failed:', error)),
                ),
            );
        } catch (error) {
            console.error('Search posts failed:', error);
            setHits([]);
            setTotal(0);
            setSearched(true);
        } finally {
            setLoading(false);
        }
    };

    useDidShow(() => {
        // 详情页返回时保留关键字与结果；首次进入聚焦输入
    });

    const onInput = (value: string) => {
        setKeyword(value);
        if (timerRef.current) clearTimeout(timerRef.current);
        timerRef.current = setTimeout(() => runSearch(value), 400);
    };

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[28rpx]'>
                {/* 搜索框 */}
                <View className='au-anim au-input mb-[24rpx]'>
                    <XIcon name='carbon:search' size={26} className='text-primary shrink-0' />
                    <Input
                      className='au-input-field'
                      value={keyword}
                      focus
                      confirmType='search'
                      placeholder={t('page.posts.search_placeholder')}
                      placeholderClass='text-textWeak'
                      onInput={(e) => onInput(e.detail.value)}
                      onConfirm={() => runSearch(keyword)}
                    />
                    {keyword && (
                        <View
                          className='w-[52rpx] h-[52rpx] flex items-center justify-center rounded-full bg-splitLine shrink-0'
                          hoverClass='tap-active'
                          onClick={() => {
                            setKeyword('');
                            setHits([]);
                            setTotal(0);
                            setSearched(false);
                          }}
                        >
                            <XIcon name='carbon:close' size={20} className='text-textSec' />
                        </View>
                    )}
                </View>

                {/* 结果统计 */}
                {searched && !loading && (
                    <Text className='text-tips text-textThird block mb-[20rpx]'>
                        {t('page.posts.search_results_count', {count: total})}
                    </Text>
                )}

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
                ) : hits.length === 0 ? (
                    searched ? (
                        <View className='au-card'>
                            <AppEmpty
                              description={t('page.posts.no_search_results', {query: keyword.trim()})}
                            />
                        </View>
                    ) : (
                        <View className='au-card'>
                            <AppEmpty description={t('page.posts.search')} />
                        </View>
                    )
                ) : (
                    <View className='flex flex-col gap-[22rpx]'>
                        {hits.map((hit) => {
                            const detail = hit.postId ? details[hit.postId] : undefined;
                            return (
                                <PostCard
                                  key={hit.postId}
                                  title={hit.title || t('page.home.untitled')}
                                  summary={detail ? getPostSummary(detail) : undefined}
                                  thumbnail={detail ? getPostThumbnail(detail) : undefined}
                                  author={detail?.authorName}
                                  publishTime={detail?.publishTime}
                                  onClick={() => hit.postId && router.push(`/post/detail?id=${hit.postId}`)}
                                />
                              );
                        })}
                    </View>
                )}
            </View>
        </View>
    );
}
