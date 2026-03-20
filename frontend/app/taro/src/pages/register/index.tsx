import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';

import {useI18nRouter} from '@/i18n/helpers';
import XIcon from '@/plugins/xicon';

import AccountRegisterPage from './components/AccountRegisterPage';
import EmailRegisterPage from './components/EmailRegisterPage';
import PhoneRegisterPage from './components/PhoneRegisterPage';
import OtherRegisterPage from './components/OtherRegisterPage';
import ControlPanel from '@/components/layout/ControlPanel';

import './register.scss';

export default function RegisterPage() {
  const {t} = useTranslation();
  const [activeTab, setActiveTab] = useState<'account' | 'email' | 'phone' | 'other'>('account');

  const router = useI18nRouter();

  // 登录
  const handleLoginClick = () => {
    router.push('/login');
  };

  // 返回首页
  const handleBackHome = () => {
    router.push('/');
  };

  // 服务条款
  const handleTermsClick = () => {
    router.push('/terms');
  };

  // 隐私政策
  const handlePrivacyClick = () => {
    router.push('/privacy');
  };

  return (
    <View className='register-page'>
      {/* 顶部控制按钮 */}
      <ControlPanel/>

      {/* 左侧品牌区 */}
      <View className='register-left'>
        <View className='brand'>
          <Image src='/logo.png' className='brand-logo'/>
          <Text className='brand-title'>{t('authentication.login.brand_title')}</Text>
          <Text className='brand-subtitle'>{t('authentication.login.brand_subtitle')}</Text>
        </View>

        <View className='benefits-list'>
          <View className='benefit-item'>
            <XIcon name='carbon:checkmark' size={20} className='benefit-icon' />
            <Text>{t('authentication.login.feature_projects')}</Text>
          </View>
          <View className='benefit-item'>
            <XIcon name='carbon:checkmark' size={20} className='benefit-icon' />
            <Text>{t('authentication.login.feature_isolation')}</Text>
          </View>
          <View className='benefit-item'>
            <XIcon name='carbon:checkmark' size={20} className='benefit-icon' />
            <Text>{t('authentication.login.feature_permissions')}</Text>
          </View>
          <View className='benefit-item'>
            <XIcon name='carbon:checkmark' size={20} className='benefit-icon' />
            <Text>{t('authentication.login.feature_analytics')}</Text>
          </View>
        </View>
      </View>

      {/* 右侧注册卡片 */}
      <View className='register-right'>
        <View className='register-card'>
          <View className='card-header'>
            <Text className='card-title'>{t('authentication.register.title')}</Text>
            <Text className='card-subtitle'>{t('authentication.register.register_with')}</Text>
          </View>

          {/* Tab 切换 */}
          <View className='register-tabs'>
            <button
              className={`tab ${activeTab === 'account' ? 'active' : ''}`}
              onClick={() => setActiveTab('account')}
            >
              {t('authentication.login.tab_account')}
            </button>
            <button
              className={`tab ${activeTab === 'email' ? 'active' : ''}`}
              onClick={() => setActiveTab('email')}
            >
              {t('authentication.login.tab_email')}
            </button>
            <button
              className={`tab ${activeTab === 'phone' ? 'active' : ''}`}
              onClick={() => setActiveTab('phone')}
            >
              {t('authentication.login.tab_phone')}
            </button>
            <button
              className={`tab ${activeTab === 'other' ? 'active' : ''}`}
              onClick={() => setActiveTab('other')}
            >
              {t('authentication.login.tab_other')}
            </button>
          </View>

          {/* 注册表单内容 */}
          <View className='register-content'>
            {activeTab === 'account' && <AccountRegisterPage/>}
            {activeTab === 'email' && <EmailRegisterPage/>}
            {activeTab === 'phone' && <PhoneRegisterPage/>}
            {activeTab === 'other' && <OtherRegisterPage/>}
          </View>

          {/* 登录链接 */}
          <View className='login-section'>
            <Text>
              {t('authentication.register.already_have_account')}
              <button className='text-btn' onClick={handleLoginClick}>
                {t('authentication.register.login_now')}
              </button>
            </Text>
          </View>

          {/* 返回首页 */}
          <View className='back-home'>
            <button className='text-btn' onClick={handleBackHome}>
              ← {t('authentication.login.back_home')}
            </button>
          </View>

          {/* 服务条款 */}
          <View className='terms'>
            <Text>
              {t('authentication.login.terms_prefix')}
              <button className='text-btn' onClick={handleTermsClick}>
                {t('authentication.login.terms_of_service')}
              </button>
              {t('login.terms_and')}
              <button className='text-btn' onClick={handlePrivacyClick}>
                {t('authentication.login.privacy_policy')}
              </button>
            </Text>
          </View>
        </View>
      </View>
    </View>
  );
}
