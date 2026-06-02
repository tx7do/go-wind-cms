'use client';

import {useState} from 'react';
import {useTranslations} from 'next-intl';

export default function EmailRegisterPage() {
    const t = useTranslations('authentication');

    const [email, setEmail] = useState('');
    const [visibleEnter, setVisibleEnter] = useState(false);

    // 简单的邮箱格式验证
    const isValidEmail = () => {
        if (!email) return false;
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    };

    const handleButtonNext = () => {
        if (!isValidEmail()) {
            return;
        }
        setVisibleEnter(true);
    };

    const inputBase = 'w-full rounded-lg border border-border bg-background px-4 py-2.5 text-sm text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15';

    return (
        <div>
            {!visibleEnter ? (
                <div className="space-y-4">
                    {/* Email Input Group */}
                    <div className="space-y-2">
                        <label htmlFor="register-email-address" className="block text-sm font-medium text-foreground">
                            {t('register.email')}
                        </label>
                        <input
                            id="register-email-address"
                            type="text"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder={t('register.input_email')}
                            autoComplete="email"
                            className={`${inputBase} ${email && !isValidEmail() ? 'border-destructive focus:border-destructive focus:ring-destructive/15' : ''}`}
                        />
                        {email && !isValidEmail() && (
                            <span className="text-xs text-destructive">{t('register.invalid_email')}</span>
                        )}
                    </div>

                    {/* Next Button */}
                    <button
                        type="button"
                        className="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60"
                        disabled={!isValidEmail()}
                        onClick={handleButtonNext}
                    >
                        {t('register.next_step')}
                    </button>
                </div>
            ) : (
                <EmailRegisterEnterCodePage email={email}/>
            )}
        </div>
    );
}

function EmailRegisterEnterCodePage({email}: { email: string }) {
    const t = useTranslations('authentication');
    const [code, setCode] = useState(Array(6).fill(''));
    const [isCodeComplete, setIsCodeComplete] = useState(false);
    const [codeSent, setCodeSent] = useState(true);
    const [codeCountdown, setCodeCountdown] = useState(60);

    // 启动倒计时
    const startCountdown = () => {
        setCodeSent(true);
        setCodeCountdown(60);

        const timer = setInterval(() => {
            setCodeCountdown((prev) => {
                if (prev <= 1) {
                    clearInterval(timer);
                    setCodeSent(false);
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);
    };

    const handleInputComplete = (value: string) => {
        console.log('验证码输入完成：', value);
        setIsCodeComplete(true);
    };

    const handleButtonResend = () => {
        if (codeSent) return;
        startCountdown();
    };

    const handleButtonConfirm = () => {
        if (!isCodeComplete) return;
        console.log('确认注册');
    };

    return (
        <div className="space-y-4">
            {/* Email Hint */}
            <div className="rounded-lg bg-primary/5 p-4 text-center">
                <p className="text-sm font-medium text-foreground">{t('register.code_sent_title')}</p>
                <p className="mt-1 text-base font-semibold text-primary">{email}</p>
                <p className="mt-1 text-xs text-muted-foreground">{t('register.code_sent_subtitle')}</p>
            </div>

            {/* Verification Code Input */}
            <div>
                <input
                    type="text"
                    maxLength={6}
                    value={code.join('')}
                    onChange={(e) => {
                        const value = e.target.value.split('').slice(0, 6);
                        setCode(value);
                        if (value.length === 6) {
                            handleInputComplete(value.join(''));
                        }
                    }}
                    className="w-full rounded-lg border border-border bg-background px-4 py-3 text-center text-lg tracking-widest text-foreground transition-colors hover:border-primary focus:border-primary focus:outline-none focus:ring-[3px] focus:ring-primary/15"
                    placeholder="请输入验证码"
                />
            </div>

            {/* Resend Button */}
            <div className="text-center">
                <button
                    className="text-sm text-primary transition-colors hover:text-primary/80 hover:underline cursor-pointer bg-transparent border-none disabled:cursor-not-allowed disabled:text-muted-foreground disabled:no-underline"
                    disabled={codeSent}
                    onClick={handleButtonResend}
                >
                    {codeSent ? `${codeCountdown}s ${t('register.resend_after')}` : t('register.resend')}
                </button>
            </div>

            {/* Confirm Button */}
            <button
                type="button"
                className="w-full cursor-pointer rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground transition-all hover:bg-primary/90 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60"
                disabled={!isCodeComplete}
                onClick={handleButtonConfirm}
            >
                {t('register.confirm')}
            </button>
        </div>
    );
}
