import {View, Text} from '@tarojs/components';
import {useState} from 'react';
import {useDidShow} from '@tarojs/taro';
import {useTranslation} from 'react-i18next';
import XIcon from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers';
import {usePageTitle} from '@/hooks/usePageTitle';
import {fetchListCategories, getCategoryName, getCategoryDescription} from '@/api/hooks/category';
import {AppEmpty} from '@/components/ui';
import type {contentservicev1_Category} from '@/api/generated/app/service/v1';

const CATEGORY_ICONS = ['carbon:document-multiple-01', 'carbon:carbon', 'carbon:globe', 'carbon:cloud-foundry-2', 'carbon:development', 'carbon:data-base', 'carbon:security', 'carbon:analytics'];
const CATEGORY_TINTS = [
    'rgba(46, 150, 255, 0.14)',
    'rgba(34, 211, 238, 0.14)',
    'rgba(139, 92, 246, 0.14)',
    'rgba(255, 159, 75, 0.16)',
];
const CATEGORY_ICON_COLORS = ['#2e96ff', '#22c3d8', '#8b5cf6', '#ff9f4b'];

/** 分类总览：树形分类卡片网格（含文章计数与子分类数） */
export default function CategoryListPage() {
    const {t} = useTranslation();
    const router = useI18nRouter();
    usePageTitle('page.title.categories');

    const [categories, setCategories] = useState<contentservicev1_Category[]>([]);
    const [loading, setLoading] = useState(true);

    useDidShow(() => {
        fetchListCategories({paging: {page: 1, pageSize: 50}})
            .then((res) => setCategories(res.items || []))
            .catch((error) => {
                console.error('Load categories failed:', error);
                setCategories([]);
            })
            .finally(() => setLoading(false));
    });

    return (
        <View className='au-page w-full bg-pageBg pb-[48rpx]'>
            <View className='px-[24rpx] pt-[36rpx]'>
                <View className='au-anim mb-[28rpx]'>
                    <Text className='au-hero-title text-textMain block'>{t('page.categories.categories_list')}</Text>
                    <Text className='text-desc text-textThird block mt-[10rpx]'>
                        {t('page.categories.explore_all')}
                        {categories.length > 0 ? ` · ${categories.length}` : ''}
                    </Text>
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
                    <View className='au-card'>
                        <AppEmpty description={t('page.categories.no_categories')} />
                    </View>
                ) : (
                    <View className='grid grid-cols-2 gap-[22rpx]'>
                        {categories.map((category, idx) => {
                            const childCount = category.children?.length ?? 0;
                            return (
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
                                        {getCategoryName(category) || t('page.categories.category_untitled')}
                                    </Text>
                                    {getCategoryDescription(category) && (
                                        <Text className='text-tips text-textThird block mt-[6rpx]' numberOfLines={2}>
                                            {getCategoryDescription(category)}
                                        </Text>
                                    )}
                                    <View className='flex items-center gap-[14rpx] mt-auto pt-[18rpx]'>
                                        <Text className='text-tips text-primary font-medium'>
                                            {(category.postCount ?? 0)} {t('page.categories.articles_count')}
                                        </Text>
                                        {childCount > 0 && (
                                            <Text className='text-tips text-textThird'>
                                                {childCount} {t('page.categories.sub_categories')}
                                            </Text>
                                        )}
                                    </View>
                                </View>
                            );
                        })}
                    </View>
                )}
            </View>
        </View>
    );
}
