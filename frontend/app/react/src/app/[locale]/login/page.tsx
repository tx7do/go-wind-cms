'use client';

import {useState} from 'react';
import {useTranslations} from 'next-intl';

import {useI18nRouter} from "@/i18n/helpers";

import ControlPanel from '@/components/layout/ControlPanel';

import AccountLoginPage from './components/AccountLoginPage';
import EmailLoginPage from './components/EmailLoginPage';
import PhoneLoginPage from './components/PhoneLoginPage';
import OtherLoginPage from './components/OtherLoginPage';

export default function LoginPage() {
    const t = useTranslations('authentication');
    const [activeTab, setActiveTab] = useState<'account' | 'email' | 'phone' | 'other'>('account');

    const router = useI18nRouter();

    const textBtn = 'text-sm text-primary transition-colors hover:text-primary/80 hover:underline cursor-pointer bg-transparent border-none';
    const tabBase = 'flex-1 cursor-pointer border-b-2 px-2 py-3 text-center text-sm font-medium transition-colors bg-transparent border-border text-muted-foreground hover:text-foreground';
    const tabActive = 'border-b-primary text-primary';

    return (
        <div className="flex min-h-screen bg-gradient-to-br from-primary/5 via-background to-background">
            {/* 顶部控制按钮 */}
            <ControlPanel/>

            {/* 左侧品牌区 */}
            <div className="hidden flex-1 flex-col justify-center px-12 max-md:hidden lg:flex">
                <div className="flex flex-col items-start gap-4">
                    <img src="/logo.png" alt={t('login.logo_alt')} className="h-16 w-auto"/>
                    <h1 className="text-3xl font-bold text-foreground">{t('login.brand_title')}</h1>
                    <p className="text-lg text-muted-foreground">{t('login.brand_subtitle')}</p>
                </div>

                <div className="mt-12 space-y-4">
                    <div className="flex items-center gap-3 text-sm text-muted-foreground">
                        <span className="text-primary">✓</span>
                        <span>{t('login.feature_projects')}</span>
                    </div>
                    <div className="flex items-center gap-3 text-sm text-muted-foreground">
                        <span className="text-primary">✓</span>
                        <span>{t('login.feature_isolation')}</span>
                    </div>
                    <div className="flex items-center gap-3 text-sm text-muted-foreground">
                        <span className="text-primary">✓</span>
                        <span>{t('login.feature_permissions')}</span>
                    </div>
                    <div className="flex items-center gap-3 text-sm text-muted-foreground">
                        <span className="text-primary">✓</span>
                        <span>{t('login.feature_analytics')}</span>
                    </div>
                </div>
            </div>

            {/* 右侧登录卡片 */}
            <div className="flex flex-1 items-center justify-center p-4 lg:flex-[0_0_480px]">
                <div className="w-full max-w-[420px] rounded-2xl border border-border bg-card p-8 shadow-lg max-md:p-6">
                    <div className="mb-6">
                        <h2 className="text-2xl font-bold text-foreground">{t('login.title')}</h2>
                        <p className="mt-2 text-sm text-muted-foreground">{t('login.login_with')}</p>
                    </div>

                    {/* Tab 切换 */}
                    <div className="mb-6 flex border-b border-border">
                        <button
                            className={`${tabBase} ${activeTab === 'account' ? tabActive : ''}`}
                            onClick={() => setActiveTab('account')}
                        >
                            {t('login.tab_account')}
                        </button>
                        <button
                            className={`${tabBase} ${activeTab === 'email' ? tabActive : ''}`}
                            onClick={() => setActiveTab('email')}
                        >
                            {t('login.tab_email')}
                        </button>
                        <button
                            className={`${tabBase} ${activeTab === 'phone' ? tabActive : ''}`}
                            onClick={() => setActiveTab('phone')}
                        >
                            {t('login.tab_phone')}
                        </button>
                        <button
                            className={`${tabBase} ${activeTab === 'other' ? tabActive : ''}`}
                            onClick={() => setActiveTab('other')}
                        >
                            {t('login.tab_other')}
                        </button>
                    </div>

                    {/* 登录表单内容 */}
                    <div className="mb-6">
                        {activeTab === 'account' && <AccountLoginPage/>}
                        {activeTab === 'email' && <EmailLoginPage/>}
                        {activeTab === 'phone' && <PhoneLoginPage/>}
                        {activeTab === 'other' && <OtherLoginPage/>}
                    </div>

                    {/* 注册链接 */}
                    <div className="mb-4 text-center text-sm text-muted-foreground">
                        <p>
                            {t('login.no_account')}
                            <button className={`${textBtn} ml-1`} onClick={() => router.push('/register')}>
                                {t('login.register_now')}
                            </button>
                        </p>
                    </div>

                    {/* 返回首页 */}
                    <div className="mb-4 text-center">
                        <button className={textBtn} onClick={() => router.push('/')}>
                            ← {t('login.back_home')}
                        </button>
                    </div>

                    {/* 服务条款 */}
                    <div className="text-center">
                        <small className="text-xs text-muted-foreground">
                            {t('login.terms_prefix')}
                            <button className={`${textBtn} mx-1`} onClick={() => router.push('/terms')}>
                                {t('login.terms_of_service')}
                            </button>
                            {t('login.terms_and')}
                            <button className={`${textBtn} ml-1`} onClick={() => router.push('/privacy')}>
                                {t('login.privacy_policy')}
                            </button>
                        </small>
                    </div>
                </div>
            </div>
        </div>
    );
}
