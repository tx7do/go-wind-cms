'use client';

import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';

export default function PhoneLoginPage() {
    const t = useTranslations('authentication');
    const [phone, setPhone] = useState('');
    const [verificationCode, setVerificationCode] = useState('');
    const [codeSent, setCodeSent] = useState(false);
    const [countdown, setCountdown] = useState(0);

    useEffect(() => {
        let timer: number;
        if (countdown > 0) {
            timer = window.setTimeout(() => setCountdown(countdown - 1), 1000);
        } else if (codeSent && countdown === 0) {
            setTimeout(() => setCodeSent(false));
        }
        return () => clearTimeout(timer);
    }, [countdown, codeSent]);

    const handleSendCode = () => {
        if (!phone) return;
        setCodeSent(true);
        setCountdown(60);
    };

    const handleLogin = () => {
        if (!phone || !verificationCode) return;
        console.log('登录信息:', { phone, verificationCode });
    };

    return (
        <div className="space-y-4">
            <div className="space-y-2">
                <label htmlFor="login-phone" className="block text-sm font-medium text-foreground">
                    {t('register.phone')}
                </label>
                <input
                    id="login-phone"
                    type="tel"
                    className="w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder={t('login.placeholder_phone')}
                    autoComplete="tel"
                    autoCapitalize="none"
                    autoCorrect="off"
                    spellCheck="false"
                />
            </div>

            <div className="space-y-2">
                <label htmlFor="login-code" className="block text-sm font-medium text-foreground">
                    {t('register.code')}
                </label>
                <div className="flex gap-2">
                    <input
                        id="login-code"
                        type="text"
                        className="flex-1 rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                        value={verificationCode}
                        onChange={(e) => setVerificationCode(e.target.value)}
                        placeholder={t('login.placeholder_code')}
                        maxLength={6}
                        autoComplete="one-time-code"
                        autoCapitalize="none"
                        autoCorrect="off"
                        spellCheck="false"
                        onKeyDown={(e) => e.key === 'Enter' && handleLogin()}
                    />
                    <button
                        className={`shrink-0 cursor-pointer rounded-lg border border-border bg-card px-4 py-2.5 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98] ${codeSent ? 'cursor-not-allowed opacity-60' : ''}`}
                        onClick={handleSendCode}
                        disabled={codeSent}
                    >
                        {codeSent ? `${countdown}s` : t('register.send_code')}
                    </button>
                </div>
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
