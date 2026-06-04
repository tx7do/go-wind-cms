import {View} from '@tarojs/components';
import {useState} from 'react';
import {useDidShow} from '@tarojs/taro';

import FeaturedPostsSection from '@/components/home/FeaturedPostsSection';
import CategoryListSection from '@/components/home/CategoryListSection';
import LatestPostsSection from '@/components/home/LatestPostsSection';
import PopularTagsSection from '@/components/home/PopularTagsSection';
import FeaturesSection from '@/components/home/FeaturesSection';

/**
 * 首页
 * 布局：最新文章（视觉重心） → 内容分类 → 推荐阅读 → 热门标签 → 平台特性
 *
 * useDidShow：每次页面显示（包括从子页面返回）时递增 refreshKey，
 * 各 Section 以 refreshKey 为 key 重新挂载并刷新数据。
 */
export default function Home() {
    const [refreshKey, setRefreshKey] = useState(0);

    // 页面每次显示时触发（首次进入 + 从其他页面返回）
    useDidShow(() => {
        setRefreshKey(prev => prev + 1);
    });

    return (
        <View className='w-full bg-pageBg pb-[104rpx]'>
            {/* 最新文章 - 视觉重心，放最前面 */}
            <View className='px-[24rpx] pt-[32rpx] pb-[32rpx]'>
                <LatestPostsSection key={`latest-${refreshKey}`} />
            </View>

            {/* 内容分类（横向滑动） */}
            <View className='px-[24rpx] pb-[32rpx]'>
                <CategoryListSection key={`cat-${refreshKey}`} />
            </View>

            {/* 推荐阅读（横排紧凑卡片，视觉较轻） */}
            <View className='px-[24rpx] pb-[32rpx]'>
                <FeaturedPostsSection key={`feat-${refreshKey}`} />
            </View>

            {/* 热门标签 */}
            <View className='px-[24rpx] pb-[32rpx]'>
                <PopularTagsSection key={`tags-${refreshKey}`} />
            </View>

            {/* 平台特性 */}
            <View className='px-[24rpx] pb-[24rpx]'>
                <FeaturesSection />
            </View>
        </View>
    );
}
