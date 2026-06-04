import {View} from '@tarojs/components';

import FeaturedPostsSection from '@/components/home/FeaturedPostsSection';
import CategoryListSection from '@/components/home/CategoryListSection';
import LatestPostsSection from '@/components/home/LatestPostsSection';
import PopularTagsSection from '@/components/home/PopularTagsSection';
import FeaturesSection from '@/components/home/FeaturesSection';

/**
 * 首页
 * 布局：最新文章（视觉重心） → 内容分类 → 推荐阅读 → 热门标签 → 平台特性
 */
export default function Home() {
    return (
        <View className='w-full bg-pageBg pb-[104rpx]'>
            {/* 最新文章 - 视觉重心，放最前面 */}
            <View className='px-[24rpx] pt-[32rpx] pb-[32rpx]'>
                <LatestPostsSection />
            </View>

            {/* 内容分类（横向滑动） */}
            <View className='px-[24rpx] pb-[32rpx]'>
                <CategoryListSection />
            </View>

            {/* 推荐阅读（横排紧凑卡片，视觉较轻） */}
            <View className='px-[24rpx] pb-[32rpx]'>
                <FeaturedPostsSection />
            </View>

            {/* 热门标签 */}
            <View className='px-[24rpx] pb-[32rpx]'>
                <PopularTagsSection />
            </View>

            {/* 平台特性 */}
            <View className='px-[24rpx] pb-[24rpx]'>
                <FeaturesSection />
            </View>
        </View>
    );
}
