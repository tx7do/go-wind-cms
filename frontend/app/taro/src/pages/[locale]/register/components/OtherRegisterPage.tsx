'use client';

import {useTranslations} from 'next-intl';

import XIcon from '@/plugins/xicon';
import styles from '../register.scss';

export default function OtherRegisterPage() {
    const t = useTranslations('authentication');

    const handleButtonGoogleRegister = () => {
        // TODO: 实现 Google 登录逻辑
        console.log('Google 注册');
    };

    const handleButtonGithubRegister = () => {
        // TODO: 实现 GitHub 登录逻辑
        console.log('GitHub 注册');
    };

    return (
        <div className={styles['other-register']}>
            {/* Social Buttons */}
            <button
                type="button"
                className={styles['social-button']}
                onClick={handleButtonGoogleRegister}
            >
                <XIcon name="logos:google-icon" size={18}/>
                {t('login.social_google')}
            </button>

            <button
                type="button"
                className={styles['social-button']}
                onClick={handleButtonGithubRegister}
            >
                <XIcon name="logos:github-icon" size={18}/>
                {t('login.social_github')}
            </button>
        </div>
    );
}
