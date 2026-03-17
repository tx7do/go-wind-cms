import {useEffect} from "react";
import {View} from '@tarojs/components';

import styles from './page.scss';

import HeroSection from '@/components/home/HeroSection';
import FeaturedPostsSection from '@/components/home/FeaturedPostsSection';
import CategoryListSection from '@/components/home/CategoryListSection';
import PopularTagsSection from '@/components/home/PopularTagsSection';
import LatestPostsSection from '@/components/home/LatestPostsSection';
import FeaturesSection from '@/components/home/FeaturesSection';

export default function Home() {

  // 获取路由参数
  useEffect(() => {
    console.log('[Home Page] Component mounted');
    // 可以在这里获取 URL 参数中的 locale
  }, []);

  return (
    <View className={styles.page}>
      <HeroSection/>
      <FeaturedPostsSection/>
      <CategoryListSection/>
      <PopularTagsSection/>
      <LatestPostsSection/>
      <FeaturesSection/>
    </View>
  );
}

