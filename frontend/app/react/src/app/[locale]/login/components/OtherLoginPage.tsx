'use client';

import {useTranslations} from 'next-intl';
import XIcon from '@/plugins/xicon';

export default function OtherLoginPage() {
    const t = useTranslations('authentication');

    const handleGoogleLogin = () => {
        console.log('Google 登录');
    };

    const handleGithubLogin = () => {
        console.log('GitHub 登录');
    };

    return (
        <div className="space-y-3">
            <button
                className="flex w-full cursor-pointer items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98]"
                onClick={handleGithubLogin}
            >
                <XIcon name="mdi:github" size={20}/>
                <span>{t('login.social_github')}</span>
            </button>

            <button
                className="flex w-full cursor-pointer items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5 active:scale-[0.98]"
                onClick={handleGoogleLogin}
            >
                <XIcon name="mdi:google" size={20}/>
                <span>{t('login.social_google')}</span>
            </button>
        </div>
    );
}
