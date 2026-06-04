import {useEffect} from 'react';
import {View} from '@tarojs/components';

import HeroSection from '@/components/home/HeroSection';
import FeaturedPostsSection from '@/components/home/FeaturedPostsSection';
import CategoryListSection from '@/components/home/CategoryListSection';
import PopularTagsSection from '@/components/home/PopularTagsSection';
import LatestPostsSection from '@/components/home/LatestPostsSection';
import FeaturesSection from '@/components/home/FeaturesSection';
import SectionContainer from '@/components/layout/SectionContainer';

export default function Home() {
    useEffect(() => {
        console.log('[Home Page] Component mounted');
    }, []);

    return (
        <View className='w-full min-h-screen bg-background'>
            {/* Hero Section */}
            <HeroSection />

            {/* Featured Posts */}
            <SectionContainer>
                <FeaturedPostsSection />
            </SectionContainer>

            {/* Categories */}
            <SectionContainer>
                <CategoryListSection />
            </SectionContainer>

            {/* Latest Posts */}
            <SectionContainer>
                <LatestPostsSection />
            </SectionContainer>

            {/* Popular Tags */}
            <SectionContainer>
                <PopularTagsSection />
            </SectionContainer>

            {/* Features — 自带背景和容器 */}
            <FeaturesSection />
        </View>
    );
}
