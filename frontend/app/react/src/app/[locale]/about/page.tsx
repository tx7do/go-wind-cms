'use client';

import {useTranslations} from 'next-intl';
import {Button} from '@/components/ui/button';

import XIcon from '@/plugins/xicon';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用

export default function AboutPage() {
    const t = useTranslations('page.about');

    const features = [
        {
            icon: 'carbon:document-add',
            title: t('feature_content'),
            description: t('feature_content_desc'),
        },
        {
            icon: 'carbon:cloud-upload',
            title: t('feature_multi_tenant'),
            description: t('feature_multi_tenant_desc'),
        },
        {
            icon: 'carbon:security',
            title: t('feature_security'),
            description: t('feature_security_desc'),
        },
        {
            icon: 'carbon:api',
            title: t('feature_api'),
            description: t('feature_api_desc'),
        },
        {
            icon: 'carbon:collaborate',
            title: t('feature_collaboration'),
            description: t('feature_collaboration_desc'),
        },
        {
            icon: 'carbon:analytics',
            title: t('feature_analytics'),
            description: t('feature_analytics_desc'),
        },
    ];

    const teamMembers = [
        {
            name: t('team_member_1'),
            role: t('team_role_1'),
            avatar: '/logo.png',
        },
        {
            name: t('team_member_2'),
            role: t('team_role_2'),
            avatar: '/logo.png',
        },
        {
            name: t('team_member_3'),
            role: t('team_role_3'),
            avatar: '/logo.png',
        },
    ];

    return (
        <div className="w-full">
            {/* Hero Section */}
            <section className="relative w-full flex min-h-[400px] items-center justify-center overflow-hidden bg-background py-28 text-center">
                <div className="relative z-10 w-full max-w-[900px] px-8 max-md:px-4">
                    <h1 className="mb-6 text-5xl font-bold tracking-tight text-foreground max-md:text-3xl">
                        {t('title')}
                    </h1>
                    <p className="mb-4 text-2xl font-light text-primary max-md:text-xl">
                        {t('subtitle')}
                    </p>
                    <p className="w-full max-w-[700px] mx-auto text-lg text-muted-foreground max-md:text-base">
                        {t('description')}
                    </p>
                </div>
            </section>

            {/* About Section */}
            <section className="w-full bg-background py-20">
                <div className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                    <div className="grid grid-cols-1 gap-12 lg:grid-cols-2 lg:gap-16">
                        <div>
                            <h2 className="mb-6 text-3xl font-bold text-foreground">{t('about_us')}</h2>
                            <p className="mb-4 leading-relaxed text-muted-foreground">{t('about_us_desc_1')}</p>
                            <p className="mb-4 leading-relaxed text-muted-foreground">{t('about_us_desc_2')}</p>
                            <p className="leading-relaxed text-muted-foreground">{t('about_us_desc_3')}</p>
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            {[
                                {num: '10K+', label: t('stat_users')},
                                {num: '500+', label: t('stat_projects')},
                                {num: '99.9%', label: t('stat_uptime')},
                                {num: '24/7', label: t('stat_support')},
                            ].map((stat) => (
                                <div key={stat.label} className="rounded-xl border border-border bg-card p-6 text-center shadow-sm">
                                    <div className="mb-2 text-3xl font-bold text-primary">{stat.num}</div>
                                    <div className="text-sm text-muted-foreground">{stat.label}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </section>

            {/* Features Section */}
            <section className="w-full bg-muted/30 py-20">
                <div className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                    <div className="mb-12 text-center">
                        <h2 className="mb-4 text-3xl font-bold text-foreground">{t('features')}</h2>
                        <p className="text-muted-foreground">{t('features_desc')}</p>
                    </div>
                    <div className="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-6">
                        {features.map((feature) => (
                            <div key={feature.title} className="group rounded-xl border border-border bg-card p-8 text-center shadow-sm transition-all duration-300 hover:-translate-y-1 hover:border-primary/40 hover:shadow-lg">
                                <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-xl bg-primary/10 text-primary transition-all duration-300 group-hover:bg-primary group-hover:text-white group-hover:shadow-md">
                                    <XIcon name={feature.icon} size={32}/>
                                </div>
                                <h3 className="mb-2 text-lg font-bold text-foreground transition-colors group-hover:text-primary">{feature.title}</h3>
                                <p className="text-sm leading-relaxed text-muted-foreground">{feature.description}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Team Section */}
            <section className="w-full bg-background py-20">
                <div className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                    <div className="mb-12 text-center">
                        <h2 className="mb-4 text-3xl font-bold text-foreground">{t('team')}</h2>
                        <p className="text-muted-foreground">{t('team_desc')}</p>
                    </div>
                    <div className="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-6">
                        {teamMembers.map((member) => (
                            <div key={member.name} className="rounded-xl border border-border bg-card p-8 text-center shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-lg">
                                <div className="mx-auto mb-4 h-[120px] w-[120px] overflow-hidden rounded-full ring-2 ring-primary/20">
                                    {/* eslint-disable-next-line @next/next/no-img-element */}
                                    <img src={member.avatar} alt={member.name} width={120} height={120}/>
                                </div>
                                <h3 className="mb-1 text-lg font-bold text-foreground">{member.name}</h3>
                                <p className="text-sm text-muted-foreground">{member.role}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Values Section */}
            <section className="w-full bg-muted/30 py-20">
                <div className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                    <div className="mb-12 text-center">
                        <h2 className="mb-4 text-3xl font-bold text-foreground">{t('values')}</h2>
                        <p className="text-muted-foreground">{t('values_desc')}</p>
                    </div>
                    <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
                        {[
                            {title: t('value_innovation'), desc: t('value_innovation_desc')},
                            {title: t('value_reliability'), desc: t('value_reliability_desc')},
                            {title: t('value_customer'), desc: t('value_customer_desc')},
                        ].map((value) => (
                            <div key={value.title} className="rounded-xl border border-border bg-card p-8 shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-lg">
                                <h3 className="mb-3 text-xl font-bold text-primary">{value.title}</h3>
                                <p className="leading-relaxed text-muted-foreground">{value.desc}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* CTA Section */}
            <section className="w-full bg-primary/5 py-20">
                <div className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                    <div className="mx-auto max-w-[800px] text-center">
                        <h2 className="mb-4 text-3xl font-bold text-foreground">{t('cta_title')}</h2>
                        <p className="mb-8 text-lg text-muted-foreground">{t('cta_desc')}</p>
                        <div className="flex flex-wrap justify-center gap-4">
                            <Button size="lg">
                                {t('cta_explore')}
                            </Button>
                            <Button variant="outline" size="lg">
                                {t('cta_contact')}
                            </Button>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    );
}
