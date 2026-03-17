import {useTranslations} from 'next-intl';
import XIcon from '@/plugins/xicon';
import styles from '../login.scss';

export default function OtherLoginPage() {
    const t = useTranslations('authentication');

    const handleGoogleLogin = () => {
        console.log('Google 登录');
    };

    const handleGithubLogin = () => {
        console.log('GitHub 登录');
    };

    return (
        <div className={styles['other-login']}>
            <button className={styles['social-button']} onClick={handleGithubLogin}>
                <XIcon name="mdi:github" size={20}/>
                <span>{t('login.social_github')}</span>
            </button>

            <button className={styles['social-button']} onClick={handleGoogleLogin}>
                <XIcon name="mdi:google" size={20}/>
                <span>{t('login.social_google')}</span>
            </button>
        </div>
    );
}
