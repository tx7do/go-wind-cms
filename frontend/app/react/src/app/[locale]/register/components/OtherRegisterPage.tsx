'use client';

import {useTranslations} from 'next-intl';

import XIcon from '@/plugins/xicon';

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
        <div className="space-y-3">
            {/* Social Buttons */}
            <button
                type="button"
                className="flex w-full cursor-pointer items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98]"
                onClick={handleButtonGoogleRegister}
            >
                <XIcon name="logos:google-icon" size={18}/>
                {t('login.social_google')}
            </button>

            <button
                type="button"
                className="flex w-full cursor-pointer items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98]"
                onClick={handleButtonGithubRegister}
            >
                <XIcon name="logos:github-icon" size={18}/>
                {t('login.social_github')}
            </button>
        </div>
    );
}
