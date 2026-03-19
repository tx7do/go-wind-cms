import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';

import './index.scss';

interface Feature {
  icon: string;
  title: string;
  description: string;
}

export default function FeaturesSection() {
  const {t} = useTranslation();

  // 特性列表
  const features: Feature[] = [
    {
      icon: 'carbon:document',
      title: t('page.home.flexible_content_management'),
      description: t('page.home.content_management_desc'),
    },
    {
      icon: 'carbon:cloud',
      title: t('page.home.multi_tenant_architecture'),
      description: t('page.home.multi_tenant_desc'),
    },
    {
      icon: 'carbon:security',
      title: t('page.home.enterprise_security'),
      description: t('page.home.security_desc'),
    },
    {
      icon: 'carbon:analytics',
      title: t('page.home.advanced_analytics'),
      description: t('page.home.analytics_desc'),
    },
    {
      icon: 'carbon:plug',
      title: t('page.home.api_integration'),
      description: t('page.home.api_integration_desc'),
    },
    {
      icon: 'carbon:group',
      title: t('page.home.real_time_collaboration'),
      description: t('page.home.real_time_collaboration_desc'),
    },
  ];

  return (
    <View className='features-section'>
      <View className='section-header'>
        <Text className='section-title'>
          <XIcon name='carbon:rocket' size={24} /> {t('page.home.platform_features')}
        </Text>
      </View>
      <View className='features-grid'>
        {features.map((feature, index) => (
          <View key={index} className='feature-card'>
            <View className='feature-icon'>
              <XIcon name={feature.icon} size={32} />
            </View>
            <Text className='feature-title'>{feature.title}</Text>
            <Text className='feature-description'>{feature.description}</Text>
          </View>
        ))}
      </View>
    </View>
  );
}
