import React from 'react';
import {useTranslations} from 'next-intl';

import {XIcon} from '@/plugins/xicon';

import {cn} from '@/lib/utils';

interface Feature {
    icon?: string;
    title: string;
    description: string;
}

export default function FeaturesSection() {
    const t = useTranslations('page.home');

    // 特性列表
    const features: Feature[] = [
        {
            icon: 'carbon:document',
            title: t('flexible_content_management'),
            description: t('content_management_desc'),
        },
        {
            icon: 'carbon:cloud',
            title: t('multi_tenant_architecture'),
            description: t('multi_tenant_desc'),
        },
        {
            icon: 'carbon:security',
            title: t('enterprise_security'),
            description: t('security_desc'),
        },
        {
            icon: 'carbon:analytics',
            title: t('advanced_analytics'),
            description: t('analytics_desc'),
        },
        {
            icon: 'carbon:api',
            title: t('api_integration'),
            description: t('api_integration_desc'),
        },
        {
            icon: 'carbon:collaborate',
            title: t('real_time_collaboration'),
            description: t('real_time_collaboration_desc'),
        },
    ];

    return (
        <section className="w-full bg-background py-16">
            <div className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                <div className="mb-12 text-center">
                    <h2 className="inline-flex items-center gap-2 text-[2.25rem] font-extrabold leading-tight tracking-tight text-foreground">
                        <XIcon name="carbon:rocket" size={28} className="mr-2 text-primary"/>
                        {t('platform_features')}
                    </h2>
                </div>
                <div className="grid grid-cols-[repeat(auto-fit,minmax(300px,1fr))] gap-8">
                    {features.map((feature, index) => (
                        <div
                            key={index}
                            className={cn(
                                'group relative overflow-hidden rounded-xl border border-border bg-card p-10 text-center',
                                'shadow-sm transition-all duration-400',
                                'hover:-translate-y-2 hover:border-primary/40 hover:shadow-lg',
                            )}
                        >
                            <div className={cn(
                                'mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-xl',
                                'bg-primary/10 text-3xl text-primary',
                                'transition-all duration-300',
                                'group-hover:scale-110 group-hover:-rotate-10 group-hover:bg-primary group-hover:text-white group-hover:shadow-lg',
                            )}>
                                <XIcon name={feature.icon || ''} size={48}/>
                            </div>
                            <h3 className="mb-3 text-xl font-bold text-foreground transition-colors group-hover:text-primary">
                                {feature.title}
                            </h3>
                            <p className="text-[15px] leading-relaxed text-muted-foreground">
                                {feature.description}
                            </p>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
}
