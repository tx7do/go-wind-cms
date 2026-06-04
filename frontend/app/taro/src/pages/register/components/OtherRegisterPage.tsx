import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import XIcon from '@/plugins/xicon';export default function OtherRegisterPage() {
  const {t} = useTranslation();

  const handleButtonGoogleRegister = () => {
    // TODO: 实现 Google 登录逻辑
    console.log('Google 注册');
  };

  const handleButtonGithubRegister = () => {
    // TODO: 实现 GitHub 登录逻辑
    console.log('GitHub 注册');
  };

  return (
    <View className='other-register'>
      {/* Google Button */}
      <View
        className='social-button'
        onClick={handleButtonGoogleRegister}
      >
        <XIcon name='logos:google-icon' size={20} />
        <Text>{t('authentication.login.social_google')}</Text>
      </View>
      
      {/* GitHub Button */}
      <View
        className='social-button'
        onClick={handleButtonGithubRegister}
      >
        <XIcon name='logos:github-icon' size={20} />
        <Text>{t('authentication.login.social_github')}</Text>
      </View>
    </View>
  );
}
