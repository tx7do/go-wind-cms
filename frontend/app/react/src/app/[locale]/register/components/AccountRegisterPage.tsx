'use client';

import {useState, useMemo} from 'react';
import {useTranslations} from 'next-intl';

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

    const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15';

    return (
        <div className="space-y-4">
            {/* Username */}
            <div className="space-y-2">
                <label htmlFor="register-account-username" className="block text-sm font-medium text-foreground">
                    {t('register.username')}
                </label>
                <input
                    id="register-account-username"
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder={t('register.input_username')}
                    autoComplete="username"
                    autoCapitalize="none"
                    autoCorrect="off"
                    spellCheck="false"
                    className={`${inputBase} ${username && !isValidUsername ? 'border-destructive focus:border-destructive focus:ring-destructive/15' : ''}`}
                />
                {username && !isValidUsername && (
                    <span className="text-xs text-destructive">{t('register.invalid_username')}</span>
                )}
            </div>

            {/* Password */}
            <div className="space-y-2">
                <label htmlFor="register-account-password" className="block text-sm font-medium text-foreground">
                    {t('register.password')}
                </label>
                <input
                    id="register-account-password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder={t('register.input_password')}
                    autoComplete="new-password"
                    className={`${inputBase} ${password && !isValidPassword ? 'border-destructive focus:border-destructive focus:ring-destructive/15' : ''}`}
                />
                {password && !isValidPassword && (
                    <span className="text-xs text-destructive">{t('register.invalid_password')}</span>
                )}
            </div>

            {/* Confirm Password */}
            <div className="space-y-2">
                <label htmlFor="register-account-confirm-password" className="block text-sm font-medium text-foreground">
                    {t('register.confirm_password')}
                </label>
                <input
                    id="register-account-confirm-password"
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder={t('register.input_confirm_password')}
                    autoComplete="new-password"
                    className={`${inputBase} ${confirmPassword && !isPasswordMatch ? 'border-destructive focus:border-destructive focus:ring-destructive/15' : ''}`}
                />
                {confirmPassword && !isPasswordMatch && (
                    <span className="text-xs text-destructive">{t('register.password_not_match')}</span>
                )}
            </div>

            {/* Register Button */}
            <button
                type="button"
                className="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60"
                disabled={!isFormValid}
                onClick={handleButtonRegister}
            >
                {t('register.register')}
            </button>
        </div>
    );
}
