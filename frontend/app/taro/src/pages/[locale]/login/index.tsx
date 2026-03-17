import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';

import {useI18nRouter} from '@/i18n/helpers';

import ControlPanel from '@/components/layout/ControlPanel';

import AccountLoginPage from './components/AccountLoginPage';
import EmailLoginPage from './components/EmailLoginPage';
import PhoneLoginPage from './components/PhoneLoginPage';
import OtherLoginPage from './components/OtherLoginPage';

import './index.scss';

export default function LoginPage() {
  const {t} = useTranslation('authentication');
  const [activeTab, setActiveTab] = useState<'account' | 'email' | 'phone' | 'other'>('account');

  const router = useI18nRouter();

  return (
    <View className='login-page'>
      {/* 顶部控制按钮 */}
      <ControlPanel/>

      {/* 左侧品牌区 */}
      <View className='login-left'>
        <View className='brand'>
          <Image src='/logo.png' className='brand-logo'/>
          <Text className='brand-title'>{t('login.brand_title')}</Text>
          <Text className='brand-subtitle'>{t('login.brand_subtitle')}</Text>
        </View>

        <View className='features-list'>
          <View className='feature-item'>
            <Text>✓</Text>
            <Text>{t('login.feature_projects')}</Text>
          </View>
          <View className='feature-item'>
            <Text>✓</Text>
            <Text>{t('login.feature_isolation')}</Text>
          </View>
          <View className='feature-item'>
            <Text>✓</Text>
            <Text>{t('login.feature_permissions')}</Text>
          </View>
          <View className='feature-item'>
            <Text>✓</Text>
            <Text>{t('login.feature_analytics')}</Text>
          </View>
        </View>
      </View>

      {/* 右侧登录卡片 */}
      <View className='login-right'>
        <View className='login-card'>
          <View className='card-header'>
            <Text className='card-title'>{t('login.title')}</Text>
            <Text className='card-subtitle'>{t('login.login_with')}</Text>
          </View>

          {/* Tab 切换 */}
          <View className='login-tabs'>
            <button
              className={`tab ${activeTab === 'account' ? 'active' : ''}`}
              onClick={() => setActiveTab('account')}
            >
              {t('login.tab_account')}
            </button>
            <button
              className={`tab ${activeTab === 'email' ? 'active' : ''}`}
              onClick={() => setActiveTab('email')}
            >
              {t('login.tab_email')}
            </button>
            <button
              className={`tab ${activeTab === 'phone' ? 'active' : ''}`}
              onClick={() => setActiveTab('phone')}
            >
              {t('login.tab_phone')}
            </button>
            <button
              className={`tab ${activeTab === 'other' ? 'active' : ''}`}
              onClick={() => setActiveTab('other')}
            >
              {t('login.tab_other')}
            </button>
          </View>

          {/* 登录表单内容 */}
          <View className='login-content'>
            {activeTab === 'account' && <AccountLoginPage/>}
            {activeTab === 'email' && <EmailLoginPage/>}
            {activeTab === 'phone' && <PhoneLoginPage/>}
            {activeTab === 'other' && <OtherLoginPage/>}
          </View>

          {/* 注册链接 */}
          <View className='register-section'>
            <Text>
              {t('login.no_account')}
              <button className='text-btn' onClick={() => router.push('/register')}>
                {t('login.register_now')}
              </button>
            </Text>
          </View>

          {/* 返回首页 */}
          <View className='back-home'>
            <button className='text-btn' onClick={() => router.push('/')}>
              ← {t('login.back_home')}
            </button>
          </View>

          {/* 服务条款 */}
          <View className='terms'>
            <Text>
              {t('login.terms_prefix')}
              <button className='text-btn' onClick={() => router.push('/terms')}>
                {t('login.terms_of_service')}
              </button>
              {t('login.terms_and')}
              <button className='text-btn' onClick={() => router.push('/privacy')}>
                {t('login.privacy_policy')}
              </button>
            </Text>
          </View>
        </View>
      </View>
    </View>
  );
}
