import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import './terms.scss';

export default function TermsPage() {
  const {t} = useTranslation();

  return (
    <View className="info-page">
      {/* Hero Section */}
      <View className="hero">
        <View className="hero-content">
          <Text className="hero-title">{t('page.legal.terms.title')}</Text>
          <Text className="hero-subtitle">{t('page.legal.terms.description')}</Text>
        </View>
      </View>

      {/* Content Section */}
      <View className="info-card">
        <View className="list">
          <Text>{t('page.legal.terms.item_1')}</Text>
          <Text>{t('page.legal.terms.item_2')}</Text>
          <Text>{t('page.legal.terms.item_3')}</Text>
        </View>
      </View>
    </View>
  );
}
