import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import './index.scss';

interface Feature {
  icon?: string;
  title: string;
  description: string;
}

export default function FeaturesSection() {
  const {t} = useTranslation('page.home');

  // 特性列表
  const features: Feature[] = [
    {
      icon: '📄',
      title: t('flexible_content_management'),
      description: t('content_management_desc'),
    },
    {
      icon: '☁️',
      title: t('multi_tenant_architecture'),
      description: t('multi_tenant_desc'),
    },
    {
      icon: '🔒',
      title: t('enterprise_security'),
      description: t('security_desc'),
    },
    {
      icon: '📊',
      title: t('advanced_analytics'),
      description: t('analytics_desc'),
    },
    {
      icon: '🔌',
      title: t('api_integration'),
      description: t('api_integration_desc'),
    },
    {
      icon: '👥',
      title: t('real_time_collaboration'),
      description: t('real_time_collaboration_desc'),
    },
  ];

  return (
    <View className='features-section'>
      <View className='section-header'>
        <Text className='section-title'>
          🚀 {t('platform_features')}
        </Text>
      </View>
      <View className='features-grid'>
        {features.map((feature, index) => (
          <View key={index} className='feature-card'>
            <View className='feature-icon'>{feature.icon}</View>
            <Text className='feature-title'>{feature.title}</Text>
            <Text className='feature-description'>{feature.description}</Text>
          </View>
        ))}
      </View>
    </View>
  );
}
