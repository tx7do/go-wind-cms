import {View, Text} from '@tarojs/components';
import {useState, useEffect, useCallback, useRef} from 'react';
import {Skeleton} from '@/components/ui/skeleton';
import {useTranslations} from '@/lib/next-intl-compat';
import {XIcon} from '@/plugins/xicon';
import {fetchListTags} from '@/api/hooks/tag';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {contentservicev1_ListTagResponse} from '@/api/generated/app/service/v1';

interface TagItem {
    id: number;
    name: string;
    color: string;
    postCount: number;
}

/** 预定义标签色板，避免 HSL 在小程序端兼容性问题 */
const TAG_COLORS = [
    '#1677ff', '#00b42a', '#ff7d00', '#f53f3f',
    '#722ed1', '#13c2c2', '#eb2f96', '#faad14',
];

/**
 * 首页「热门标签」区块
 * - flex wrap 流式布局
 * - 显示最多 6 个标签
 */
export default function PopularTagsSection() {
    const t = useTranslations('page.tags');
    const router = useI18nRouter();
    const [loading, setLoading] = useState(false);
    const [displayTags, setDisplayTags] = useState<TagItem[]>([]);
    const abortRef = useRef<AbortController | null>(null);

    const loadPopularTags = useCallback(async () => {
        if (abortRef.current) abortRef.current.abort();
        abortRef.current = new AbortController();
        setLoading(true);
        try {
            const res = await fetchListTags({
                paging: {page: 1, pageSize: 6},
                formValues: {status: 'TAG_STATUS_ACTIVE', isFeatured: true},
            }) as unknown as contentservicev1_ListTagResponse;
            const tagItems = (res.items || [])
                .filter(tag => tag.id !== undefined)
                .map((tag, index) => ({
                    id: tag.id!,
                    name: tag.translations?.[0]?.name || t('tag_untitled'),
                    color: tag.color || TAG_COLORS[index % TAG_COLORS.length],
                    postCount: tag.postCount || 0,
                }));
            setDisplayTags(tagItems);
        } catch (error) {
            console.error('Failed to load tags:', error);
            setDisplayTags([]);
        } finally {
            setLoading(false);
        }
    }, [t]);

    useEffect(() => {
        loadPopularTags();
        return () => { abortRef.current?.abort(); };
    }, []);

    return (
        <View className='w-full'>
            {/* 标题行 */}
            <View className='flex items-center justify-between mb-[32rpx]'>
                <View className='flex items-center gap-[8rpx]'>
                    <XIcon name='carbon:fire' size={20} className='text-primary' />
                    <Text className='text-card-title font-bold text-textMain'>
                        {t('popular_tags')}
                    </Text>
                </View>
                <View
                  className='px-[24rpx] py-[12rpx] rounded-full flex items-center justify-center'
                  style={{
                      backgroundColor: 'rgba(22,119,255,0.08)',
                  }}
                  onClick={() => router.push('/tag')}
                  hoverClass='tap-active'
                >
                    <Text className='text-tips font-medium text-primary'>{t('view_all')} →</Text>
                </View>
            </View>

            {loading ? (
                <View className='flex flex-wrap gap-[16rpx]'>
                    {Array.from({length: 6}).map((_, i) => (
                        <View key={i} className='h-[64rpx] w-[160rpx]'>
                            <Skeleton className='h-full w-full rounded-full' />
                        </View>
                    ))}
                </View>
            ) : displayTags.length > 0 ? (
                <View className='flex flex-wrap gap-[16rpx]'>
                    {displayTags.map((tag) => (
                        <View
                          key={tag.id}
                          className='flex items-center gap-[8rpx] rounded-full px-[24rpx] py-[12rpx] min-h-[64rpx]'
                          style={{
                                backgroundColor: `${tag.color}14`,
                                borderWidth: '2rpx',
                                borderStyle: 'solid',
                                borderColor: `${tag.color}40`,
                            }}
                          onClick={() => router.push(`/tag/detail?id=${tag.id}`)}
                          hoverClass='tap-active'
                        >
                            <View
                              className='w-[12rpx] h-[12rpx] rounded-full'
                              style={{backgroundColor: tag.color}}
                            />
                            <Text className='text-desc font-semibold' style={{color: tag.color}}>
                                {tag.name}
                            </Text>
                            {tag.postCount > 0 && (
                                <Text className='text-tips' style={{color: tag.color}}>
                                    {tag.postCount}
                                </Text>
                            )}
                        </View>
                    ))}
                </View>
            ) : null}
        </View>
    );
}
