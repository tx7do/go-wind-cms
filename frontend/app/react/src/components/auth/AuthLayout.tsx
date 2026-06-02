'use client';

import React from 'react';
import {useTranslations} from 'next-intl';
import {useI18nRouter} from '@/i18n/helpers';
import ControlPanel from '@/components/layout/ControlPanel';

export type AuthTab = 'account' | 'email' | 'phone' | 'other';

interface AuthLayoutProps {
    /** 表单标题（如"登录"或"注册"） */
    title: string;
    /** 副标题（如"使用以下方式登录"） */
    subtitle: string;
    /** 当前激活的 Tab */
    activeTab: AuthTab;
    /** Tab 切换回调 */
    onTabChange: (tab: AuthTab) => void;
    /** 表单内容 */
    children: React.ReactNode;
    /** 底部链接区域：切换登录/注册的链接 */
    switchLink?: React.ReactNode;
}

const TAB_KEYS: AuthTab[] = ['account', 'email', 'phone', 'other'];

/**
 * 认证页面通用布局（login/register 共享）
 *
 * 包含：左侧品牌区 + 右侧卡片（Tab切换 + 表单内容 + 底部链接 + 服务条款）
 */
export default function AuthLayout({
    title,
    subtitle,
    activeTab,
    onTabChange,
    children,
    switchLink,
}: AuthLayoutProps) {
    const t = useTranslations('authentication');
    const router = useI18nRouter();

    const textBtn = 'text-sm text-primary transition-colors hover:text-primary/80 hover:underline cursor-pointer bg-transparent border-none';
    const tabBase = 'flex-1 cursor-pointer border-b-2 px-2 py-3 text-center text-sm font-medium transition-colors bg-transparent border-border text-muted-foreground hover:text-foreground';
    const tabActive = 'border-b-primary text-primary';

    return (
        <div className="flex min-h-screen bg-gradient-to-br from-primary/5 via-background to-background">
            <ControlPanel/>

            {/* 左侧品牌区 */}
            <div className="hidden flex-1 flex-col justify-center px-12 max-md:hidden lg:flex">
                <div className="flex flex-col items-start gap-4">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img src="/logo.png" alt={t('login.logo_alt')} className="h-16 w-auto"/>
                    <h1 className="text-3xl font-bold text-foreground">{t('login.brand_title')}</h1>
                    <p className="text-lg text-muted-foreground">{t('login.brand_subtitle')}</p>
                </div>

                <div className="mt-12 space-y-4">
                    {['feature_projects', 'feature_isolation', 'feature_permissions', 'feature_analytics'].map(key => (
                        <div key={key} className="flex items-center gap-3 text-sm text-muted-foreground">
                            <span className="text-primary">✓</span>
                            <span>{t(`login.${key}`)}</span>
                        </div>
                    ))}
                </div>
            </div>

            {/* 右侧认证卡片 */}
            <div className="flex flex-1 items-center justify-center p-4 lg:flex-[0_0_480px]">
                <div className="w-full max-w-[420px] rounded-2xl border border-border bg-card p-8 shadow-lg max-md:p-6">
                    <div className="mb-6">
                        <h2 className="text-2xl font-bold text-foreground">{title}</h2>
                        <p className="mt-2 text-sm text-muted-foreground">{subtitle}</p>
                    </div>

                    {/* Tab 切换 */}
                    <div className="mb-6 flex border-b border-border">
                        {TAB_KEYS.map(tab => (
                            <button
                                key={tab}
                                className={`${tabBase} ${activeTab === tab ? tabActive : ''}`}
                                onClick={() => onTabChange(tab)}
                            >
                                {t(`login.tab_${tab}`)}
                            </button>
                        ))}
                    </div>

                    {/* 表单内容 */}
                    <div className="mb-6">
                        {children}
                    </div>

                    {/* 切换链接 */}
                    {switchLink && (
                        <div className="mb-4 text-center text-sm text-muted-foreground">
                            {switchLink}
                        </div>
                    )}

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
