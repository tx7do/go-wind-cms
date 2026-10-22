import {View, Text} from '@tarojs/components';
import {useState} from 'react';
import {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import {fetchListTags, getTagTranslation} from '@/api/hooks/tag';
import {AppEmpty} from '@/components/ui';
import type {contentservicev1_Tag} from '@/api/generated/app/service/v1';

/** 标签总览：热门标签按文章计数排序，胶囊流式布局 */
export default function TagListPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    usePageTitle('page.title.tags');

    const [tags, setTags] = useState<contentservicev1_Tag[]>([]);
    const [loading, setLoading] = useState(true);

    useDidShow(() => {
        fetchListTags({
            paging: {page: 1, pageSize: 60},
            orderBy: ['-postCount'],
        })
            .then((res) => setTags(res.items || []))
            .catch((error) => {
                console.error('Load tags failed:', error);
                setTags([]);
            })
            .finally(() => setLoading(false));
    });

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[36rpx]'>
                <View className='au-anim mb-[28rpx]'>
                    <Text className='au-hero-title text-textMain block'>{t('page.tags.tags_list')}</Text>
                    <Text className='text-desc text-textThird block mt-[10rpx]'>
                        {t('page.tags.total_tags')}
                        {tags.length > 0 ? ` · ${tags.length}` : ''}
                    </Text>
                </View>

                {loading ? (
                    <View className='flex flex-wrap gap-[16rpx]'>
                        {['w-[140rpx]', 'w-[180rpx]', 'w-[120rpx]', 'w-[200rpx]', 'w-[150rpx]', 'w-[170rpx]', 'w-[130rpx]', 'w-[190rpx]'].map((w, i) => (
                            <View key={i} className={`au-skeleton ${w} h-[52rpx] rounded-full`} />
                        ))}
                    </View>
                ) : tags.length === 0 ? (
                    <View className='au-card'>
                        <AppEmpty description={t('page.tags.no_tags')} />
                    </View>
                ) : (
                    <View className='flex flex-wrap gap-[16rpx]'>
                        {tags.map((tag) => (
                            <View
                              key={tag.id}
                              className='au-card au-card--glow px-[26rpx] h-[76rpx] flex items-center gap-[10rpx] rounded-full'
                              hoverClass='tap-active'
                              onClick={() => tag.id && router.push(`/tag/detail?id=${tag.id}`)}
                            >
                                {tag.color && (
                                    <View className='w-[14rpx] h-[14rpx] rounded-full' style={{background: tag.color}} />
                                )}
                                <Text className='text-desc font-medium text-textMain' numberOfLines={1}>
                                    {getTagTranslation(tag)?.name || tag.code}
                                </Text>
                                <Text className='text-tips text-textThird'>
                                    {tag.postCount ?? 0}
                                </Text>
                            </View>
                        ))}
                    </View>
                )}
            </View>
        </View>
    );
}
