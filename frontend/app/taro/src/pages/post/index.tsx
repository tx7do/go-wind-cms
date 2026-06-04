import {useState, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import CategoryFilter from '@/components/category/CategoryFilter';
import PostList from '@/components/post/PostList';

export default function PostListPage() {
    const {t} = useTranslation();
    const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(null);

    const queryParams = useMemo(() => {
        if (selectedCategoryId) return {category_ids__in: [selectedCategoryId]};
        return {};
    }, [selectedCategoryId]);

    return (
        <View className='min-h-screen w-full bg-pageBg pb-[160rpx]'>
            {/* 页面标题 */}
            <View className='bg-cardBg px-[24rpx] pt-[40rpx] pb-[24rpx] border-b-[1rpx] border-splitLine'>
                <Text className='text-title font-bold text-textMain'>{t('page.posts.posts_list')}</Text>
            </View>

            {/* 分类筛选 */}
            <View className='px-[24rpx] pt-[24rpx]'>
                <CategoryFilter
                  selectedCategory={selectedCategoryId}
                  treeMode
                  autoLoad
                  onCategoryChange={setSelectedCategoryId}
                />
            </View>

            {/* 文章列表 */}
            <View className='px-[24rpx] pt-[24rpx]'>
                <PostList
                  key={selectedCategoryId || 'all'}
                  queryParams={queryParams}
                  initialPageSize={12}
                  showPagination
                />
            </View>
        </View>
    );
}
