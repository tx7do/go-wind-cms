'use client';

import {Suspense, useState} from 'react';
import {useTranslations} from 'next-intl';
import {useI18nRouter} from "@/i18n/helpers";

import AuthLayout, {type AuthTab} from '@/components/auth/AuthLayout';

import AccountLoginPage from './components/AccountLoginPage';
import EmailLoginPage from './components/EmailLoginPage';
import PhoneLoginPage from './components/PhoneLoginPage';
import OtherLoginPage from './components/OtherLoginPage';

export default function LoginPage() {
    const t = useTranslations('authentication');
    const router = useI18nRouter();
    const [activeTab, setActiveTab] = useState<AuthTab>('account');

    const textBtn = 'text-sm text-primary transition-colors hover:text-primary/80 hover:underline cursor-pointer bg-transparent border-none';

    return (
        <AuthLayout
            title={t('login.title')}
            subtitle={t('login.login_with')}
            activeTab={activeTab}
            onTabChange={setActiveTab}
            switchLink={
                <p>
                    {t('login.no_account')}
                    <button className={`${textBtn} ms-1`} onClick={() => router.push('/register')}>
                        {t('login.register_now')}
                    </button>
                </p>
            }
        >
            {activeTab === 'account' && <AccountLoginPage/>}
            {activeTab === 'email' && <EmailLoginPage/>}
            {activeTab === 'phone' && <PhoneLoginPage/>}
            {activeTab === 'other' && (
                // OtherLoginPage 读 useSearchParams（OAuth 回调 code/state）：
                // Suspense 边界满足静态导出 prerender 要求
                <Suspense fallback={null}>
                    <OtherLoginPage/>
                </Suspense>
            )}
        </AuthLayout>
    );
}
