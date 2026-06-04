import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import XIcon from '@/plugins/xicon';export default function OtherLoginPage() {
  const {t} = useTranslation();

  const handleGoogleLogin = () => {
    console.log('Google 登录');
  };

  const handleGithubLogin = () => {
    console.log('GitHub 登录');
  };

  return (
    <View className='other-login'>
      <View className='social-button' onClick={handleGoogleLogin}>
        <XIcon name='logos:google-icon' size={20} />
        <Text>{t('authentication.login.social_google')}</Text>
      </View>

      <View className='social-button' onClick={handleGithubLogin}>
        <XIcon name='logos:github-icon' size={20} />
        <Text>{t('authentication.login.social_github')}</Text>
      </View>
    </View>
  );
}
