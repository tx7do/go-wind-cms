import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import './privacy.scss';

export default function PrivacyPage() {
  const {t} = useTranslation();

  return (
    <View className="info-page">
      {/* Hero Section */}
      <View className="hero">
        <View className="hero-content">
          <Text className="hero-title">{t('page.legal.privacy.title')}</Text>
          <Text className="hero-subtitle">{t('page.legal.privacy.description')}</Text>
        </View>
      </View>

      {/* Content Section */}
      <View className="info-card">
        <View className="list">
          <Text>{t('page.legal.privacy.item_1')}</Text>
          <Text>{t('page.legal.privacy.item_2')}</Text>
          <Text>{t('page.legal.privacy.item_3')}</Text>
        </View>
      </View>
    </View>
  );
}
