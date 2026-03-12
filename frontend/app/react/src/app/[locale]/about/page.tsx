'use client';

import {useTranslations} from 'next-intl';
import {Image, Button} from 'antd';
import XIcon from '@/plugins/xicon';
import styles from './about.module.css';
import '../../globals.css'; // ✅ 导入全局 CSS，确保 CSS 变量可用

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
        <div className={styles.aboutPage}>
            {/* Hero Section */}
            <section className={styles.hero}>
                <div className={styles.heroContent}>
                    <h1 className={styles.heroTitle}>{t('title')}</h1>
                    <p className={styles.heroSubtitle}>{t('subtitle')}</p>
                    <p className={styles.heroDescription}>{t('description')}</p>
                </div>
            </section>

            {/* About Section */}
            <section className={styles.aboutSection}>
                <div className={styles.sectionContainer}>
                    <div className={styles.aboutContent}>
                        <h2>{t('about_us')}</h2>
                        <p>{t('about_us_desc_1')}</p>
                        <p>{t('about_us_desc_2')}</p>
                        <p>{t('about_us_desc_3')}</p>
                    </div>
                    <div className={styles.aboutStats}>
                        <div className={styles.statCard}>
                            <div className={styles.statNumber}>10K+</div>
                            <div className={styles.statLabel}>{t('stat_users')}</div>
                        </div>
                        <div className={styles.statCard}>
                            <div className={styles.statNumber}>500+</div>
                            <div className={styles.statLabel}>{t('stat_projects')}</div>
                        </div>
                        <div className={styles.statCard}>
                            <div className={styles.statNumber}>99.9%</div>
                            <div className={styles.statLabel}>{t('stat_uptime')}</div>
                        </div>
                        <div className={styles.statCard}>
                            <div className={styles.statNumber}>24/7</div>
                            <div className={styles.statLabel}>{t('stat_support')}</div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Features Section */}
            <section className={styles.featuresSection}>
                <div className={styles.sectionContainer}>
                    <div className={styles.sectionHeader}>
                        <h2>{t('features')}</h2>
                        <p>{t('features_desc')}</p>
                    </div>
                    <div className={styles.featuresGrid}>
                        {features.map((feature) => (
                            <div key={feature.title} className={styles.featureCard}>
                                <div className={styles.featureIcon}>
                                    <XIcon name={feature.icon} size={32} />
                                </div>
                                <h3>{feature.title}</h3>
                                <p>{feature.description}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Team Section */}
            <section className={styles.teamSection}>
                <div className={styles.sectionContainer}>
                    <div className={styles.sectionHeader}>
                        <h2>{t('team')}</h2>
                        <p>{t('team_desc')}</p>
                    </div>
                    <div className={styles.teamGrid}>
                        {teamMembers.map((member) => (
                            <div key={member.name} className={styles.teamCard}>
                                <div className={styles.teamAvatar}>
                                    <Image
                                        src={member.avatar}
                                        alt={member.name}
                                        width={120}
                                        height={120}
                                        preview={false}
                                    />
                                </div>
                                <h3>{member.name}</h3>
                                <p>{member.role}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Values Section */}
            <section className={styles.valuesSection}>
                <div className={styles.sectionContainer}>
                    <div className={styles.sectionHeader}>
                        <h2>{t('values')}</h2>
                        <p>{t('values_desc')}</p>
                    </div>
                    <div className={styles.valuesGrid}>
                        <div className={styles.valueCard}>
                            <h3>{t('value_innovation')}</h3>
                            <p>{t('value_innovation_desc')}</p>
                        </div>
                        <div className={styles.valueCard}>
                            <h3>{t('value_reliability')}</h3>
                            <p>{t('value_reliability_desc')}</p>
                        </div>
                        <div className={styles.valueCard}>
                            <h3>{t('value_customer')}</h3>
                            <p>{t('value_customer_desc')}</p>
                        </div>
                    </div>
                </div>
            </section>

            {/* CTA Section */}
            <section className={styles.ctaSection}>
                <div className={styles.sectionContainer}>
                    <h2>{t('cta_title')}</h2>
                    <p>{t('cta_desc')}</p>
                    <div className={styles.ctaButtons}>
                        <Button type="primary" size="large">
                            {t('cta_explore')}
                        </Button>
                        <Button size="large">
                            {t('cta_contact')}
                        </Button>
                    </div>
                </div>
            </section>
        </div>
    );
}

