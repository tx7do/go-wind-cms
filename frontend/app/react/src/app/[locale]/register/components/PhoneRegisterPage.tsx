'use client';

import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';

export default function PhoneRegisterPage() {
    const t = useTranslations('authentication');
    
    const [phone, setPhone] = useState('');
    const [verificationCode, setVerificationCode] = useState('');
    const [codeSent, setCodeSent] = useState(false);
    const [codeCountdown, setCodeCountdown] = useState(0);

    useEffect(() => {
        let codeTimer: number | null = null;

        if (codeSent && codeCountdown > 0) {
            codeTimer = window.setInterval(() => {
                setCodeCountdown((prev) => {
                    if (prev <= 1) {
                        setCodeSent(false);
                        return 0;
                    }
                    return prev - 1;
                });
            }, 1000);
        }

        return () => {
            if (codeTimer !== null) {
                clearInterval(codeTimer);
            }
        };
    }, [codeSent, codeCountdown]);

    /**
     * 发送验证码
     */
    const handleButtonSendVerifyCode = () => {
        if (!phone) {
            return;
        }

        setCodeSent(true);
        setCodeCountdown(60);
    };

    /**
     * 登录或者注册
     */
    const handleButtonRegisterOrLogin = () => {
        if (!phone || !verificationCode) {
            return;
        }
        console.log('手机号注册/登录:', {phone, verificationCode});
    };

    const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15';

    return (
        <div className="space-y-4">
            {/* Phone Number Group */}
            <div className="space-y-2">
                <label htmlFor="register-phone-number" className="block text-sm font-medium text-foreground">
                    {t('register.phone')}
                </label>
                <input
                    id="register-phone-number"
                    type="text"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder={t('register.input_phone')}
                    autoComplete="tel"
                    className={inputBase}
                />
            </div>

            {/* Verification Code Group */}
            <div className="space-y-2">
                <label htmlFor="register-phone-code" className="block text-sm font-medium text-foreground">
                    {t('register.code')}
                </label>
                <div className="flex gap-2">
                    <input
                        id="register-phone-code"
                        type="text"
                        value={verificationCode}
                        onChange={(e) => setVerificationCode(e.target.value)}
                        placeholder={t('register.input_code')}
                        maxLength={6}
                        autoComplete="one-time-code"
                        className="flex-1 rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                    />
                    <button
                        type="button"
                        disabled={codeSent}
                        className={`shrink-0 cursor-pointer rounded-lg border border-border bg-card px-4 py-2.5 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98] ${codeSent ? 'cursor-not-allowed opacity-60' : ''}`}
                        onClick={handleButtonSendVerifyCode}
                    >
                        {codeSent ? `${codeCountdown}s` : t('register.send_code')}
                    </button>
                </div>
            </div>

            {/* Register Button */}
            <button
                type="button"
                className="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98]"
                onClick={handleButtonRegisterOrLogin}
            >
                {t('register.register_or_login')}
            </button>
        </div>
    );
}
