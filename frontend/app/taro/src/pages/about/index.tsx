import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';

import XIcon from '@/plugins/xicon';export default function AboutPage() {
  const {t} = useTranslation();

  const features = [
    {
      icon: 'carbon:document-add',
      title: t('page.about.feature_content'),
      description: t('page.about.feature_content_desc'),
    },
    {
      icon: 'carbon:cloud-upload',
      title: t('page.about.feature_multi_tenant'),
      description: t('page.about.feature_multi_tenant_desc'),
    },
    {
      icon: 'carbon:security',
      title: t('page.about.feature_security'),
      description: t('page.about.feature_security_desc'),
    },
    {
      icon: 'carbon:api',
      title: t('page.about.feature_api'),
      description: t('page.about.feature_api_desc'),
    },
    {
      icon: 'carbon:collaborate',
      title: t('page.about.feature_collaboration'),
      description: t('page.about.feature_collaboration_desc'),
    },
    {
      icon: 'carbon:analytics',
      title: t('page.about.feature_analytics'),
      description: t('page.about.feature_analytics_desc'),
    },
  ];

  const teamMembers = [
    {
      name: t('page.about.team_member_1'),
      role: t('page.about.team_role_1'),
      avatar: '/assets/images/logo.png',
    },
    {
      name: t('page.about.team_member_2'),
      role: t('page.about.team_role_2'),
      avatar: '/assets/images/logo.png',
    },
    {
      name: t('page.about.team_member_3'),
      role: t('page.about.team_role_3'),
      avatar: '/assets/images/logo.png',
    },
  ];

  return (
    <View className='about-page'>
      {/* Hero Section */}
      <View className='hero'>
        <View className='hero-content'>
          <Text className='hero-title'>{t('page.about.title')}</Text>
          <Text className='hero-subtitle'>{t('page.about.subtitle')}</Text>
          <Text className='hero-description'>{t('page.about.description')}</Text>
        </View>
      </View>

      {/* About Section */}
      <View className='about-section'>
        <View className='section-container'>
          <View className='about-content'>
            <h2>{t('page.about.about_us')}</h2>
            <p>{t('page.about.about_us_desc_1')}</p>
            <p>{t('page.about.about_us_desc_2')}</p>
            <p>{t('page.about.about_us_desc_3')}</p>
          </View>
          <View className='about-stats'>
            <View className='stat-card'>
              <Text className='stat-number'>10K+</Text>
              <Text className='stat-label'>{t('page.about.stat_users')}</Text>
            </View>
            <View className='stat-card'>
              <Text className='stat-number'>500+</Text>
              <Text className='stat-label'>{t('page.about.stat_projects')}</Text>
            </View>
            <View className='stat-card'>
              <Text className='stat-number'>99.9%</Text>
              <Text className='stat-label'>{t('page.about.stat_uptime')}</Text>
            </View>
            <View className='stat-card'>
              <Text className='stat-number'>24/7</Text>
              <Text className='stat-label'>{t('page.about.stat_support')}</Text>
            </View>
          </View>
        </View>
      </View>

      {/* Features Section */}
      <View className='features-section'>
        <View className='section-container'>
          <View className='section-header'>
            <h2>{t('page.about.features')}</h2>
            <p>{t('page.about.features_desc')}</p>
          </View>
          <View className='features-grid'>
            {features.map((feature) => (
              <View key={feature.title} className='feature-card'>
                <View className='feature-icon'>
                  <XIcon name={feature.icon} size={36} />
                </View>
                <Text className='feature-title'>{feature.title}</Text>
                <Text className='feature-description'>{feature.description}</Text>
              </View>
            ))}
          </View>
        </View>
      </View>

      {/* Team Section */}
      <View className='team-section'>
        <View className='section-container'>
          <View className='section-header'>
            <h2>{t('page.about.team')}</h2>
            <p>{t('page.about.team_desc')}</p>
          </View>
          <View className='team-grid'>
            {teamMembers.map((member) => (
              <View key={member.name} className='team-card'>
                <View className='team-avatar'>
                  <Image
                    src={member.avatar}
                    style={{width: 120, height: 120}}
                  />
                </View>
                <Text className='team-name'>{member.name}</Text>
                <Text className='team-role'>{member.role}</Text>
              </View>
            ))}
          </View>
        </View>
      </View>

      {/* Values Section */}
      <View className='values-section'>
        <View className='section-container'>
          <View className='section-header'>
            <h2>{t('page.about.values')}</h2>
            <p>{t('page.about.values_desc')}</p>
          </View>
          <View className='values-grid'>
            <View className='value-card'>
              <Text className='value-title'>{t('page.about.value_innovation')}</Text>
              <Text className='value-description'>{t('page.about.value_innovation_desc')}</Text>
            </View>
            <View className='value-card'>
              <Text className='value-title'>{t('page.about.value_reliability')}</Text>
              <Text className='value-description'>{t('page.about.value_reliability_desc')}</Text>
            </View>
            <View className='value-card'>
              <Text className='value-title'>{t('page.about.value_customer')}</Text>
              <Text className='value-description'>{t('page.about.value_customer_desc')}</Text>
            </View>
          </View>
        </View>
      </View>

      {/* CTA Section */}
      <View className='cta-section'>
        <View className='section-container'>
          <h2>{t('page.about.cta_title')}</h2>
          <p>{t('page.about.cta_desc')}</p>
          <View className='cta-buttons'>
            <View className='button-primary'>
              <Text>{t('page.about.cta_explore')}</Text>
            </View>
            <View className='button-default'>
              <Text>{t('page.about.cta_contact')}</Text>
            </View>
          </View>
        </View>
      </View>
    </View>
  );
}
