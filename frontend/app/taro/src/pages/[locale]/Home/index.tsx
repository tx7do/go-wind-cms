import React from "react";
import {View} from "@tarojs/components";

import HeroSection from '../../components/home/HeroSection';
import FeaturedPostsSection from '../../components/home/FeaturedPostsSection';
import CategoryListSection from '../../components/home/CategoryListSection';
import PopularTagsSection from '../../components/home/PopularTagsSection';
import LatestPostsSection from '../../components/home/LatestPostsSection';
import FeaturesSection from '../../components/home/FeaturesSection';

export default function Home() {
  return (
    <View className='home-page'>
      <HeroSection/>
      <FeaturedPostsSection/>
      <CategoryListSection/>
      <PopularTagsSection/>
      <LatestPostsSection/>
      <FeaturesSection/>
    </View>
  );
}
