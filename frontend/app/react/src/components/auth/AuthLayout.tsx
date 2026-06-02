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
 * 布局结构：
 * - 外层 max-w-5xl 居中容器，消除宽屏中央黑洞
 * - 左侧品牌区：流光渐变标题 + 特性列表 + 径向光晕背景
 * - 右侧认证卡片：Tab 切换 + 表单 + 底部链接 + 服务条款
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

    // Tab 样式：未选中用 muted-foreground（足够可读），hover 变亮，选中用 primary 下划线
    const tabBase = 'flex-1 cursor-pointer border-b-2 px-2 py-3 text-center text-sm font-medium transition-colors bg-transparent border-border text-muted-foreground hover:text-foreground';
    const tabActive = 'border-b-primary text-primary';

    // 底部条款链接：primary 色 + 加粗 + 下划线，行动手持端点击热区
    const linkBtn = 'text-sm text-primary font-medium underline underline-offset-4 tracking-wide transition-colors hover:text-primary/80 cursor-pointer bg-transparent border-none';

    return (
        <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-primary/5 via-background to-background pb-12 max-md:pb-16">
            <ControlPanel/>

            {/* 中央居中的最大宽度容器 */}
            <div className="w-full max-w-5xl mx-auto px-6">
                <div className="grid grid-cols-1 gap-12 items-center lg:grid-cols-2 lg:gap-16">
                    {/* ─── 左侧：品牌宣传区 ─── */}
                    <div className="hidden lg:flex relative flex-col justify-center">
                        {/* 径向光晕背景：极淡的绿色弥散光 */}
                        <div
                            className="absolute inset-0 -z-10 opacity-60 pointer-events-none"
                            style={{
                                background: `radial-gradient(ellipse at 30% 40%, hsl(var(--primary) / 0.15) 0%, transparent 60%), radial-gradient(ellipse at 70% 70%, hsl(200 95% 60% / 0.06) 0%, transparent 55%)`,
                            }}
                        />

                        <div className="flex flex-col items-start gap-4">
                            {/* eslint-disable-next-line @next/next/no-img-element */}
                            <img src="/logo.png" alt={t('login.logo_alt')} className="h-16 w-auto"/>
                            <h1 className="text-4xl font-bold leading-tight">
                                {/* 流光渐变文字：白 → 翠绿 → 天蓝 */}
                                <span
                                    className="bg-gradient-to-r from-foreground via-primary to-sky-400 bg-clip-text text-transparent"
                                >
                                    {t('login.brand_title')}
                                </span>
                            </h1>
                            <p className="text-lg text-muted-foreground">{t('login.brand_subtitle')}</p>
                        </div>

                        <div className="mt-12 space-y-4">
                            {['feature_projects', 'feature_isolation', 'feature_permissions', 'feature_analytics'].map(key => (
                                <div key={key} className="flex items-center gap-3 text-sm text-muted-foreground">
                                    <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/15 text-primary text-xs font-bold">
                                        ✓
                                    </span>
                                    <span>{t(`login.${key}`)}</span>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* ─── 右侧：认证卡片 ─── */}
                    <div className="flex justify-center lg:justify-end">
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
                                <button className={linkBtn} onClick={() => router.push('/')}>
                                    ← {t('login.back_home')}
                                </button>
                            </div>

                            {/* 服务条款 — 加粗 + 下划线 + tracking-wide 优化手持端点击热区 */}
                            <div className="text-center leading-relaxed">
                                <small className="text-xs text-muted-foreground/80">
                                    {t('login.terms_prefix')}
                                    <button className={`${linkBtn} mx-1 text-xs`} onClick={() => router.push('/terms')}>
                                        {t('login.terms_of_service')}
                                    </button>
                                    {t('login.terms_and')}
                                    <button className={`${linkBtn} ml-1 text-xs`} onClick={() => router.push('/privacy')}>
                                        {t('login.privacy_policy')}
                                    </button>
                                </small>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
