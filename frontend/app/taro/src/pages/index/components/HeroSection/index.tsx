import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';export default function HeroSection() {
  const {t} = useTranslation();
  const router = useI18nRouter();

  console.log('[HeroSection] Rendering...', { t });

  return (
    <View className='hero'>
      {/* 简化背景 - 移除复杂动画 */}
      <View className='hero-bg-wrapper'>
        <View className='hero-gradient-bg' />
        {/* H5 端保留简单的网格背景 */}
        <View className='hero-grid-bg' />
      </View>

      <View className='hero-content'>
        <Text className='hero-title'>{t('authentication.login.brand_title')}</Text>
        <Text className='hero-subtitle'>{t('authentication.login.brand_subtitle')}</Text>
        <Text className='hero-description'>{t('page.home.hero_description')}</Text>
        <View className='hero-actions'>
          <View
            className='btn btn-primary'
            onClick={() => router.push('/post')}
          >
            <Text>{t('page.home.browse_posts')}</Text>
          </View>
          <View
            className='btn btn-secondary'
            onClick={() => router.push('/about')}
          >
            <Text>{t('page.home.learn_more')}</Text>
          </View>
        </View>
      </View>
    </View>
  );
}
