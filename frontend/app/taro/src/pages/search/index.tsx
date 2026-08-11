import {useState, useEffect} from 'react';
import Taro from '@tarojs/taro';
import {View, Text, Input} from '@tarojs/components';
import {useTranslation} from 'react-i18next';
import {useI18nRouter} from '@/i18n/helpers';
import {fetchSearchPosts} from '@/api/hooks/post';
import {Skeleton} from '@/components/ui/skeleton';
import {usePageTitle} from '@/hooks/usePageTitle';

interface SearchHit {
    postId?: number;
    language?: string;
    title?: string;
}
interface SearchResponse {
    items?: SearchHit[];
    total?: number;
}

export default function SearchPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    usePageTitle('page.title.posts');

    // 从入口页可能带入初始关键词（如首页搜索入口）
    const initialQuery = (() => {
        try {
            const instance = Taro.getCurrentInstance();
            return (instance?.router?.params?.q as string)?.trim() ?? '';
        } catch {
            return '';
        }
    })();

    const [keyword, setKeyword] = useState(initialQuery);
    const [loading, setLoading] = useState(false);
    const [hits, setHits] = useState<SearchHit[]>([]);
    const [total, setTotal] = useState(0);
    const [hasSearched, setHasSearched] = useState(false);

    async function runSearch(q: string) {
        const query = q.trim();
        if (!query) {
            setHits([]);
            setTotal(0);
            setHasSearched(false);
            return;
        }
        setLoading(true);
        setHasSearched(true);
        try {
            const res = (await fetchSearchPosts({query, pageSize: 20})) as unknown as SearchResponse;
            setHits(res?.items ?? []);
            setTotal(res?.total ?? 0);
        } catch (err) {
            console.error('Search failed:', err);
            setHits([]);
            setTotal(0);
        } finally {
            setLoading(false);
        }
    }

    // 入口带关键词时自动搜索
    useEffect(() => {
        if (initialQuery) runSearch(initialQuery);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [initialQuery]);

    const openPost = (postId?: number) => {
        if (!postId) return;
        router.push(`/post/${postId}`);
    };

    return (
        <View className='min-h-screen w-full bg-pageBg pb-[160rpx]'>
            {/* 搜索栏 */}
            <View className='flex items-center gap-[16rpx] bg-cardBg px-[24rpx] py-[20rpx] border-b-[1rpx] border-splitLine'>
                <Input
                    className='flex-1 h-[64rpx] px-[24rpx] text-[28rpx] text-textMain rounded-[12rpx] bg-pageBg'
                    value={keyword}
                    placeholder={t('page.posts.search_placeholder')}
                    onInput={(e) => setKeyword(e.detail.value)}
                    onConfirm={() => runSearch(keyword)}
                    confirmType='search'
                />
                <Text
                    className='text-[28rpx] text-primary shrink-0 px-[8rpx]'
                    onClick={() => runSearch(keyword)}
                >
                    {t('page.posts.search')}
                </Text>
            </View>

            {/* 内容区 */}
            <View className='px-[24rpx] pt-[24rpx]'>
                {loading ? (
                    <View className='flex flex-col gap-[16rpx]'>
                        {[...Array(5)].map((_, i) => (
                            <Skeleton key={i} className='h-[120rpx] w-full rounded-[12rpx]'/>
                        ))}
                    </View>
                ) : !hasSearched ? (
                    <View className='flex flex-col items-center py-[120rpx]'>
                        <Text className='text-[28rpx] text-textThird'>
                            {t('page.posts.search_placeholder')}
                        </Text>
                    </View>
                ) : hits.length === 0 ? (
                    <View className='flex flex-col items-center py-[120rpx]'>
                        <Text className='text-[28rpx] text-textThird'>
                            {t('page.posts.no_search_results', {query: initialQuery || keyword})}
                        </Text>
                    </View>
                ) : (
                    <View className='flex flex-col gap-[16rpx]'>
                        <Text className='text-[24rpx] text-textSec'>
                            {t('page.posts.search_results_count', {count: total})}
                        </Text>
                        {hits.map((hit, idx) => (
                            <View
                                key={hit.postId ?? idx}
                                className='flex items-center gap-[16rpx] bg-cardBg rounded-[12rpx] px-[24rpx] py-[28rpx]'
                                onClick={() => openPost(hit.postId)}
                            >
                                <Text className='text-[30rpx] font-medium text-textMain flex-1 line-clamp-2'>
                                    {hit.title}
                                </Text>
                                <Text className='text-[24rpx] text-textThird shrink-0'>›</Text>
                            </View>
                        ))}
                    </View>
                )}
            </View>
        </View>
    );
}
