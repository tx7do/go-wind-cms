'use client';

import {useState, useMemo} from 'react';
import {useTranslations} from 'next-intl';
import styles from '../register.scss';

export default function AccountRegisterPage() {
    const t = useTranslations('authentication');

    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');

    // 用户名验证（3-20个字符，只能包含字母、数字、下划线）
    const isValidUsername = useMemo(() => {
        if (!username) return false;
        const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;
        return usernameRegex.test(username);
    }, [username]);

    // 密码强度验证（至少 6 个字符）
    const isValidPassword = useMemo(() => {
        return password.length >= 6;
    }, [password]);

    // 确认密码验证
    const isPasswordMatch = useMemo(() => {
        return password === confirmPassword && confirmPassword.length > 0;
    }, [password, confirmPassword]);

    // 表单是否可提交
    const isFormValid = useMemo(() => {
        return isValidUsername && isValidPassword && isPasswordMatch;
    }, [isValidUsername, isValidPassword, isPasswordMatch]);

    const handleButtonRegister = () => {
        if (!isFormValid) {
            return;
        }

        console.log('注册信息：', {
            username: username,
            password: password,
        });
    };

    return (
        <div className={styles['register-form']}>
            {/* Username */}
            <div className={styles['form-group']}>
                <label htmlFor="register-account-username">
                    {t('register.username')}
                </label>
                <input
                    id="register-account-username"
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder={t('register.input_username')}
                    autoComplete="username"
                    className={`${styles['input-field']} ${username && !isValidUsername ? styles.error : ''}`}
                />
                {username && !isValidUsername && (
                    <span className={styles['error-hint']}>{t('register.invalid_username')}</span>
                )}
            </div>

            {/* Password */}
            <div className={styles['form-group']}>
                <label htmlFor="register-account-password">
                    {t('register.password')}
                </label>
                <input
                    id="register-account-password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder={t('register.input_password')}
                    autoComplete="new-password"
                    className={`${styles['input-field']} ${password && !isValidPassword ? styles.error : ''}`}
                />
                {password && !isValidPassword && (
                    <span className={styles['error-hint']}>{t('register.invalid_password')}</span>
                )}
            </div>

            {/* Confirm Password */}
            <div className={styles['form-group']}>
                <label htmlFor="register-account-confirm-password">
                    {t('register.confirm_password')}
                </label>
                <input
                    id="register-account-confirm-password"
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder={t('register.input_confirm_password')}
                    autoComplete="new-password"
                    className={`${styles['input-field']} ${confirmPassword && !isPasswordMatch ? styles.error : ''}`}
                />
                {confirmPassword && !isPasswordMatch && (
                    <span className={styles['error-hint']}>{t('register.password_not_match')}</span>
                )}
            </div>

            {/* Register Button */}
            <button
                type="button"
                className={styles['register-button']}
                disabled={!isFormValid}
                onClick={handleButtonRegister}
            >
                {t('register.register')}
            </button>
        </div>
    );
}
