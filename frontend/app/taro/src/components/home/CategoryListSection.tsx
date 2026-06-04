import {View, Text} from '@tarojs/components';
import {useState, useEffect, useCallback, useRef} from 'react';
import {Skeleton} from '@/components/ui/skeleton';
import {useTranslations} from '@/lib/next-intl-compat';
import {fetchListCategories} from '@/api/hooks/category';
import {contentservicev1_Category, contentservicev1_ListCategoryResponse} from '@/api/generated/app/service/v1';
import {XIcon} from '@/plugins/xicon';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import HomeCategoryCard from './HomeCategoryCard';

/**
 * 首页「分类」区块
 * - 横向滑动卡片布局
 * - 显示最多 8 个分类
 */
export default function CategoryListSection({pageSize = 8, page = 1}: {
    pageSize?: number;
    page?: number;
}) {
    const t = useTranslations('page.home');
    const router = useI18nRouter();
    const [categories, setCategories] = useState<contentservicev1_Category[]>([]);
    const [loading, setLoading] = useState(false);
    const abortRef = useRef<AbortController | null>(null);

    const loadCategories = useCallback(async () => {
        if (abortRef.current) abortRef.current.abort();
        abortRef.current = new AbortController();
        setLoading(true);
        try {
            const res = await fetchListCategories({
                paging: {page, pageSize},
                formValues: {status: 'CATEGORY_STATUS_ACTIVE'},
                fieldMask: 'id,status,sortOrder,icon,code,postCount,directPostCount,parent_id,createdAt,translations.id,translations.categoryId,translations.name,translations.languageCode,translations.description',
                orderBy: ['-sortOrder', '-postCount'],
            }) as unknown as contentservicev1_ListCategoryResponse;
            if (!abortRef.current.signal.aborted) setCategories(res.items || []);
        } catch (error) {
            if (!abortRef.current?.signal?.aborted) {
                console.error('Failed to load categories:', error);
                setCategories([]);
            }
        } finally {
            if (!abortRef.current?.signal?.aborted) setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadCategories();
        return () => { abortRef.current?.abort(); };
    }, []);

    return (
        <View className='w-full'>
            {/* 标题行 */}
            <View className='flex items-center justify-between mb-[32rpx]'>
                <View className='flex items-center gap-[8rpx]'>
                    <XIcon name='carbon:folder-details' size={20} className='text-primary' />
                    <Text className='text-card-title font-bold text-textMain'>
                        {t('categories')}
                    </Text>
                </View>
                <View
                  className='px-[24rpx] py-[12rpx] rounded-full flex items-center justify-center'
                  style={{
                      backgroundColor: 'rgba(22,119,255,0.08)',
                  }}
                  onClick={() => router.push('/category')}
                  hoverClass='tap-active'
                >
                    <Text className='text-tips font-medium text-primary'>{t('view_all')} →</Text>
                </View>
            </View>

            {/* 横向滑动分类列表 - 使用正内边距确保边缘留白 */}
            {loading ? (
                <View className='flex gap-[16rpx] overflow-hidden'>
                    {Array.from({length: 4}).map((_, i) => (
                        <View key={i} className='flex-shrink-0 w-[180rpx] flex flex-col items-center p-[20rpx] rounded bg-cardBg'>
                            <Skeleton className='w-[72rpx] h-[72rpx] rounded mb-[12rpx]' />
                            <Skeleton className='h-[24rpx] w-[100rpx] rounded mb-[8rpx]' />
                            <Skeleton className='h-[20rpx] w-[60rpx] rounded' />
                        </View>
                    ))}
                </View>
            ) : categories.length > 0 ? (
                <View className='flex gap-[16rpx] overflow-x-auto no-scrollbar'>
                    {categories.map((category) => (
                        <HomeCategoryCard
                          key={category.id}
                          category={category}
                          onClick={(id) => router.push(`/category/detail?id=${id}`)}
                          mobileCompact
                        />
                    ))}
                </View>
            ) : null}
        </View>
    );
}
