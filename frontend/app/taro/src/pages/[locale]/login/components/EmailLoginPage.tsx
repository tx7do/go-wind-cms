'use client';

import {useState} from 'react';
import {useTranslations} from 'next-intl';
import styles from '../login.scss';
import {useAuthenticationStore} from "@/store/slices/authentication/hooks";

export default function EmailLoginPage() {
    const t = useTranslations('authentication');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [rememberMe, setRememberMe] = useState(false);

    const authenticationStore = useAuthenticationStore();

    const handleLogin = async () => {
        if (!email || !password) {
            return;
        }
        console.log('登录信息:', {email, password, rememberMe});

        await authenticationStore.login({
            email: email,
            password: password,
        });
    };

    const handleForgotPassword = () => {
        console.log('忘记密码');
    };

    return (
        <div className={styles['login-form']}>
            <div className={styles['form-group']}>
                <label htmlFor="login-email">{t('register.email')}</label>
                <input
                    id="login-email"
                    type="email"
                    className={styles['input-field']}
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder={t('login.placeholder_email')}
                    autoComplete="email"
                />
            </div>

            <div className={styles['form-group']}>
                <label htmlFor="login-password-email">{t('login.password')}</label>
                <input
                    id="login-password-email"
                    type="password"
                    className={styles['input-field']}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder={t('login.placeholder_password')}
                    autoComplete="current-password"
                    onKeyDown={(e) => e.key === 'Enter' && handleLogin()}
                />
            </div>

            <div className={styles['form-options']}>
                <label className={styles['checkbox-label']}>
                    <input
                        type="checkbox"
                        checked={rememberMe}
                        onChange={(e) => setRememberMe(e.target.checked)}
                    />
                    <span>{t('login.remember_me')}</span>
                </label>
                <button className={styles['text-btn']} onClick={handleForgotPassword}>
                    {t('login.forgot_password')}
                </button>
            </div>

            <button className={styles['login-button']} onClick={handleLogin}>
                {t('login.login')}
            </button>
        </div>
    );
}
