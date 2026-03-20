import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import './contact.scss';

export default function ContactPage() {
  const {t} = useTranslation();

  return (
    <View className="info-page">
      {/* Hero Section */}
      <View className="hero">
        <View className="hero-content">
          <Text className="hero-title">{t('page.legal.contact.title')}</Text>
          <Text className="hero-subtitle">{t('page.legal.contact.description')}</Text>
        </View>
      </View>

      {/* Content Section */}
      <View className="info-card">
        <View className="list">
          <Text>{t('page.legal.contact.item_1')}</Text>
          <Text>{t('page.legal.contact.item_2')}</Text>
          <Text>{t('page.legal.contact.item_3')}</Text>
        </View>
      </View>
    </View>
  );
}
