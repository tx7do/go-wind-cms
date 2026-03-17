import {useEffect} from "react";
import {View} from '@tarojs/components';
import Taro from '@tarojs/taro';

import '../../globals.css';
import styles from './page.scss';

import HeroSection from '@/components/home/HeroSection';
import FeaturedPostsSection from '@/components/home/FeaturedPostsSection';
import CategoryListSection from '@/components/home/CategoryListSection';
import PopularTagsSection from '@/components/home/PopularTagsSection';
import LatestPostsSection from '@/components/home/LatestPostsSection';
import FeaturesSection from '@/components/home/FeaturesSection';

export default function Home() {
  // Intersection Observer for scroll reveal - Taro 中需要使用 createSelectorQuery
  useEffect(() => {
    // TODO: 在 Taro 中实现滚动动画，需要使用 Taro API
    // 由于 Taro 多端适配，建议使用 Taro.createIntersectionObserver
  }, []);

  // Scroll to categories - Taro 中使用 pageScrollTo
  const scrollToCategories = async () => {
    try {
      await Taro.pageScrollTo({
        selector: '.categories-section',
        duration: 300
      });
    } catch (err) {
      console.error('Scroll failed:', err);
    }
  };

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

