import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import './disclaimer.scss';

export default function DisclaimerPage() {
  const {t} = useTranslation();

  return (
    <View className="info-page">
      {/* Hero Section */}
      <View className="hero">
        <View className="hero-content">
          <Text className="hero-title">{t('page.legal.disclaimer.title')}</Text>
          <Text className="hero-subtitle">{t('page.legal.disclaimer.description')}</Text>
        </View>
      </View>

      {/* Content Section */}
      <View className="info-card">
        <View className="list">
          <Text>{t('page.legal.disclaimer.item_1')}</Text>
          <Text>{t('page.legal.disclaimer.item_2')}</Text>
          <Text>{t('page.legal.disclaimer.item_3')}</Text>
        </View>
      </View>
    </View>
  );
}
