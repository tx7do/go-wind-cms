'use client';

import React from 'react';
import {Button, Divider} from 'antd';
import {Icon} from '@iconify/react';
import {useTranslations} from 'next-intl';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import styles from './Footer.module.css';

type FooterLinkKey =
    | 'about_us'
    | 'contact_us'
    | 'non_responsibility'
    | 'privacy_agreement'
    | 'terms_of_service';

interface FooterLinkItem {
    key: FooterLinkKey;
    labelKey: string;
}

interface SocialItem {
    key: 'twitter' | 'facebook' | 'instagram' | 'telegram' | 'wechat' | 'weibo' | 'qq';
    name: string;
    icon: string;
}

export default function Footer() {
    const t = useTranslations('ui');
    const router = useI18nRouter();

    const footerLinks: FooterLinkItem[] = [
        {key: 'about_us', labelKey: 'button.about_us'},
        {key: 'contact_us', labelKey: 'button.contact_us'},
        {key: 'non_responsibility', labelKey: 'button.non_responsibility'},
        {key: 'privacy_agreement', labelKey: 'button.privacy_agreement'},
        {key: 'terms_of_service', labelKey: 'button.terms_of_service'},
    ];

    const footerRouteMap: Record<FooterLinkKey, string> = {
        about_us: '/about',
        contact_us: '/contact',
        non_responsibility: '/disclaimer',
        privacy_agreement: '/privacy',
        terms_of_service: '/terms',
    };

    const socialLinks: SocialItem[] = [
        {key: 'twitter', name: 'Twitter', icon: 'fa:twitter'},
        {key: 'facebook', name: 'Facebook', icon: 'fa:facebook'},
        {key: 'instagram', name: 'Instagram', icon: 'fa:instagram'},
        {key: 'telegram', name: 'Telegram', icon: 'fa:telegram'},
        {key: 'wechat', name: 'WeChat', icon: 'fa:wechat'},
        {key: 'weibo', name: 'Weibo', icon: 'fa:weibo'},
        {key: 'qq', name: 'QQ', icon: 'fa:qq'},
    ];

    const handleFooterLinkClick = (key: FooterLinkKey) => {
        router.push(footerRouteMap[key]);
    };

    const handleSocialClick = (name: string) => {
        console.log(name);
    };

    return (
        <footer className={styles.footerOuter}>
            <div className={styles.footerContainer}>
                <nav className={styles.footerLinks} aria-label="Footer links">
                    {footerLinks.map(link => (
                        <Button
                            key={link.key}
                            type="text"
                            className={styles.footerLink}
                            onClick={() => handleFooterLinkClick(link.key)}
                        >
                            {t(link.labelKey)}
                        </Button>
                    ))}
                </nav>

                <div className={styles.footerMeta}>
                    <span className={styles.copyright}>
                        {t('copyright')}
                    </span>
                    <Divider vertical className={styles.metaDivider}/>
                    <div className={styles.socialList} aria-label="Social links">
                        {socialLinks.map(social => (
                            <Button
                                key={social.key}
                                type="text"
                                className={styles.socialBtn}
                                aria-label={social.name}
                                onClick={() => handleSocialClick(social.name)}
                                icon={<Icon icon={social.icon}/>}
                            />
                        ))}
                    </div>
                </div>
            </div>
        </footer>
    );
}
