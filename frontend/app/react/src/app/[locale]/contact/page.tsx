'use client';

import React from 'react';
import {useTranslations} from 'next-intl';

import PageHero from '@/components/layout/PageHero';
import {XIcon} from '@/plugins/xicon';
import {cn} from '@/lib/utils';

/**
 * 联系页面 — 3 栏卡片布局
 *
 * 视觉重心居中，避免左下角死角的"头重脚轻"
 * 卡片：技术支持 / 商务合作 / 服务时间
 * 邮箱支持 mailto: 唤起 + select-all 双击全选
 */
export default function ContactPage() {
    const t = useTranslations('page.legal.contact');

    const cards: Array<{
        icon: string;
        title: string;
        desc: string;
        email?: string;
        extra?: React.ReactNode;
        accent: 'primary' | 'blue' | 'slate';
    }> = [
        {
            icon: 'carbon:help',
            title: t('support_title'),
            desc: t('support_desc'),
            email: t('support_email'),
            accent: 'primary',
        },
        {
            icon: 'mdi:briefcase-outline',
            title: t('business_title'),
            desc: t('business_desc'),
            email: t('business_email'),
            accent: 'blue',
        },
        {
            icon: 'carbon:time',
            title: t('hours_title'),
            desc: t('hours_desc'),
            extra: (
                <div className="mt-3 space-y-1">
                    <p className="text-sm font-medium text-foreground">{t('hours_weekday')}</p>
                    <p className="text-xs text-muted-foreground">{t('hours_timezone')}</p>
                </div>
            ),
            accent: 'slate',
        },
    ];

    return (
        <div className="w-full">
            {/* Hero 区 */}
            <PageHero
                title={t('title')}
                description={t('description')}
                icon="carbon:email"
                size="sm"
            />

            {/* 3 栏卡片宫格 */}
            <section className="w-full py-12 max-md:py-8">
                <div className="mx-auto grid max-w-5xl grid-cols-1 gap-6 px-4 md:grid-cols-3">
                    {cards.map((card, index) => (
                        <ContactCard key={index} {...card} extra={card.extra}/>
                    ))}
                </div>
            </section>
        </div>
    );
}

/** 联系卡片 */
function ContactCard({
    icon,
    title,
    desc,
    email,
    extra,
    accent,
}: {
    icon: string;
    title: string;
    desc: string;
    email?: string;
    extra?: React.ReactNode;
    accent: 'primary' | 'blue' | 'slate';
}) {
    // 主题色映射
    const accentMap = {
        primary: {
            iconBg: 'bg-primary/10 text-primary border-primary/20',
            hoverBorder: 'hover:border-primary/40',
            emailText: 'text-primary',
        },
        blue: {
            iconBg: 'bg-blue-500/10 text-blue-500 border-blue-500/20 dark:text-blue-400',
            hoverBorder: 'hover:border-blue-500/40',
            emailText: 'text-blue-500 dark:text-blue-400',
        },
        slate: {
            iconBg: 'bg-muted text-muted-foreground border-border',
            hoverBorder: 'hover:border-border',
            emailText: 'text-foreground',
        },
    }[accent];

    return (
        <div
            className={cn(
                'group relative flex flex-col items-center rounded-2xl border border-border bg-card p-6 text-center shadow-sm',
                'transition-all duration-300 ease-out',
                'hover:-translate-y-1 hover:shadow-lg',
                accentMap.hoverBorder,
            )}
        >
            {/* 顶部图标 */}
            <div
                className={cn(
                    'mb-4 flex h-12 w-12 items-center justify-center rounded-xl border',
                    accentMap.iconBg,
                )}
            >
                <XIcon name={icon} size={22}/>
            </div>

            {/* 标题 */}
            <h3 className="text-base font-semibold text-foreground">{title}</h3>

            {/* 描述 */}
            <p className="mt-2 text-xs leading-relaxed text-muted-foreground">
                {desc}
            </p>

            {/* 邮箱 或 自定义内容 */}
            {email && (
                <a
                    href={`mailto:${email}`}
                    className={cn(
                        'mt-3 inline-block rounded-md px-3 py-1.5 text-sm font-medium select-all',
                        'transition-colors hover:underline',
                        'bg-primary/5 dark:bg-primary/10',
                        accentMap.emailText,
                    )}
                    aria-label={`Send email to ${email}`}
                >
                    {email}
                </a>
            )}

            {extra && <div className="w-full">{extra}</div>}
        </div>
    );
}
