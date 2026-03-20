import {useTranslation} from 'react-i18next';

import XIcon from '@/plugins/xicon';

import '../index.scss';

export default function OtherLoginPage() {
  const {t} = useTranslation();

  const handleGoogleLogin = () => {
    console.log('Google 登录');
  };

  const handleGithubLogin = () => {
    console.log('GitHub 登录');
  };

  return (
    <div className="other-login">
      <button className="social-button" onClick={handleGithubLogin}>
        <XIcon name="mdi:github" size={20}/>
        <span>{t('authentication.login.social_github')}</span>
      </button>

      <button className="social-button" onClick={handleGoogleLogin}>
        <XIcon name="mdi:google" size={20}/>
        <span>{t('authentication.login.social_google')}</span>
      </button>
    </div>
  );
}
