'use client';

import {useState} from 'react';
import {useTranslations} from 'next-intl';
import styles from '../login.scss';
import {useAuthenticationStore} from "@/store/slices/authentication/hooks";

export default function AccountLoginPage() {
    const t = useTranslations('authentication');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [rememberMe, setRememberMe] = useState(false);

    const authenticationStore = useAuthenticationStore();

    const handleLogin = async () => {
        if (!username || !password) {
            return;
        }
        console.log('登录信息:', {username, password, rememberMe});

        await authenticationStore.login({
            username: username,
            password: password,
        });
    };

    const handleForgotPassword = () => {
        console.log('忘记密码');
    };

    return (
        <div className={styles['login-form']}>
            <div className={styles['form-group']}>
                <label htmlFor="login-username">{t('login.username')}</label>
                <input
                    id="login-username"
                    type="text"
                    className={styles['input-field']}
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder={t('register.input_username')}
                    autoComplete="username"
                />
            </div>

            <div className={styles['form-group']}>
                <label htmlFor="login-password">{t('login.password')}</label>
                <input
                    id="login-password"
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
