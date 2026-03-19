import {useEffect} from "react";
import {View} from '@tarojs/components';

import HeroSection from './components/HeroSection';
import FeaturedPostsSection from './components/FeaturedPostsSection';
import CategoryListSection from './components/CategoryListSection';
import PopularTagsSection from './components/PopularTagsSection';
import LatestPostsSection from './components/LatestPostsSection';
import FeaturesSection from './components/FeaturesSection';

export default function Home() {

  // 获取路由参数
  useEffect(() => {
    console.log('[Home Page] Component mounted');
    // 可以在这里获取 URL 参数中的 locale
  }, []);

  return (
    <View className='page'>
      <HeroSection/>
      <View className='section-container'>
        <FeaturedPostsSection/>
      </View>
      <View className='section-container'>
        <CategoryListSection/>
      </View>
      <View className='section-container'>
        <PopularTagsSection/>
      </View>
      <View className='section-container'>
        <LatestPostsSection/>
      </View>
      <View className='section-container'>
        <FeaturesSection/>
      </View>
    </View>
  );
}

