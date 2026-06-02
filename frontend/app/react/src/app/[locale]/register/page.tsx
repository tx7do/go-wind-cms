'use client';

import {useState} from 'react';
import {useTranslations} from 'next-intl';
import {useI18nRouter} from "@/i18n/helpers";

import AuthLayout, {type AuthTab} from '@/components/auth/AuthLayout';

import AccountRegisterPage from './components/AccountRegisterPage';
import EmailRegisterPage from './components/EmailRegisterPage';
import PhoneRegisterPage from './components/PhoneRegisterPage';
import OtherRegisterPage from './components/OtherRegisterPage';

export default function RegisterPage() {
    const t = useTranslations('authentication');
    const router = useI18nRouter();
    const [activeTab, setActiveTab] = useState<AuthTab>('account');

    const textBtn = 'text-sm text-primary transition-colors hover:text-primary/80 hover:underline cursor-pointer bg-transparent border-none';

    return (
        <AuthLayout
            title={t('register.title')}
            subtitle={t('register.register_with')}
            activeTab={activeTab}
            onTabChange={setActiveTab}
            switchLink={
                <p>
                    {t('register.already_have_account')}
                    <button className={`${textBtn} ml-1`} onClick={() => router.push('/login')}>
                        {t('register.login_now')}
                    </button>
                </p>
            }
        >
            {activeTab === 'account' && <AccountRegisterPage/>}
            {activeTab === 'email' && <EmailRegisterPage/>}
            {activeTab === 'phone' && <PhoneRegisterPage/>}
            {activeTab === 'other' && <OtherRegisterPage/>}
        </AuthLayout>
    );
}
