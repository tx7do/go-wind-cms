import {View} from '@tarojs/components';

import FeaturedPostsSection from '@/components/home/FeaturedPostsSection';
import CategoryListSection from '@/components/home/CategoryListSection';
import LatestPostsSection from '@/components/home/LatestPostsSection';
import PopularTagsSection from '@/components/home/PopularTagsSection';
import FeaturesSection from '@/components/home/FeaturesSection';

/**
 * 首页
 * 布局：精选 → 分类 → 最新文章 → 热门标签 → 平台特性
 * 每个区块通过 SectionContainer 统一 px-[24rpx] 内边距
 */
export default function Home() {
    return (
        <View className='w-full bg-pageBg pb-[48rpx]'>
            {/* 精选文章 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <FeaturedPostsSection />
            </View>

            {/* 分类（横向滑动） */}
            <View className='px-[24rpx] py-[32rpx]'>
                <CategoryListSection />
            </View>

            {/* 最新文章 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <LatestPostsSection />
            </View>

            {/* 热门标签 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <PopularTagsSection />
            </View>

            {/* 平台特性 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <FeaturesSection />
            </View>
        </View>
    );
}
