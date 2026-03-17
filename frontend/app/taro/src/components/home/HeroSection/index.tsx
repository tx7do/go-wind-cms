import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import './index.scss';

export default function HeroSection() {
  const {t} = useTranslation('page.home');
  const {t: brandT} = useTranslation('authentication.login');
  const router = useI18nRouter();

  console.log('[HeroSection] Rendering...', { t, brandT });

  return (
    <View className='hero'>
      {/* Animated Background Elements */}
      <View className='hero-bg-wrapper'>
        <View className='hero-gradient-bg'/>
        {/* Grid Background */}
        <View className='hero-grid-bg'/>
        {/* Animated Shapes */}
        <View className='hero-animated-shapes'>
          <View className='shape shape1'/>
          <View className='shape shape2'/>
          <View className='shape shape3'/>
        </View>
        {/* Code Snippet Decorations */}
        <View className='hero-code-snippets'>
          <View className='code-snippet snippet1'>
            <Text className='code-line'>func GetPost() {'{'}return post{'}'}</Text>
          </View>
          <View className='code-snippet snippet2'>
            <Text className='code-line'>&lt;template&gt;&lt;div
              class="content-hub"&gt;&lt;/div&gt;&lt;/template&gt;</Text>
          </View>
          <View className='code-snippet snippet3'>
            <Text className='code-line'>{`// Content Hub API\nconst api = 'v1'`}</Text>
          </View>
        </View>
        {/* SVG Waves */}
        <View className='hero-waves'/>
      </View>

      <View className='hero-content'>
        <Text className='hero-title'>{brandT('brand_title')}</Text>
        <Text className='hero-subtitle'>{brandT('brand_subtitle')}</Text>
        <Text className='hero-description'>{t('hero_description')}</Text>
        <View className='hero-actions'>
          <View
            className='btn btn-primary'
            onClick={() => router.push('/post')}
          >
            <Text>{t('browse_posts')}</Text>
          </View>
          <View
            className='btn btn-secondary'
            onClick={() => router.push('/about')}
          >
            <Text>{t('learn_more')}</Text>
          </View>
        </View>
      </View>
    </View>
  );
}
