'use client';

import {useState} from 'react';
import {useTranslations} from 'next-intl';
import {useAuth} from "@/api/hooks/auth";

export default function AccountLoginPage() {
    const t = useTranslations('authentication');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [rememberMe, setRememberMe] = useState(false);

    const auth = useAuth();

    const handleLogin = async () => {
        if (!username || !password) {
            return;
        }
        console.log('登录信息:', {username, password, rememberMe});

        await auth.login({
            username: username,
            password: password,
        });
    };

    const handleForgotPassword = () => {
        console.log('忘记密码');
    };

    return (
        <div className="space-y-4">
            <div className="space-y-2">
                <label htmlFor="login-username" className="block text-sm font-medium text-foreground">
                    {t('login.username')}
                </label>
                <input
                    id="login-username"
                    type="text"
                    className="w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder={t('register.input_username')}
                    autoComplete="username"
                    autoCapitalize="none"
                    autoCorrect="off"
                    spellCheck="false"
                />
            </div>

            <div className="space-y-2">
                <label htmlFor="login-password" className="block text-sm font-medium text-foreground">
                    {t('login.password')}
                </label>
                <input
                    id="login-password"
                    type="password"
                    className="w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder={t('login.placeholder_password')}
                    autoComplete="current-password"
                    onKeyDown={(e) => e.key === 'Enter' && handleLogin()}
                />
            </div>

            <div className="flex w-full items-center justify-between">
                <label className="flex cursor-pointer items-center gap-2 text-sm text-muted-foreground">
                    <input
                        type="checkbox"
                        checked={rememberMe}
                        onChange={(e) => setRememberMe(e.target.checked)}
                        className="h-4 w-4 rounded border-border accent-primary"
                    />
                    <span>{t('login.remember_me')}</span>
                </label>
                <button
                    className="text-sm text-primary transition-colors hover:text-primary/80 hover:underline cursor-pointer bg-transparent border-none"
                    onClick={handleForgotPassword}
                >
                    {t('login.forgot_password')}
                </button>
            </div>

            <button
                className="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98]"
                onClick={handleLogin}
            >
                {t('login.login')}
            </button>
        </div>
    );
}
