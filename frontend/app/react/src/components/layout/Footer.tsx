'use client';

import React from 'react';
import {Button} from '@/components/ui/button';
import {Separator} from '@/components/ui/separator';
import {Icon} from '@iconify/react';
import {useTranslations} from 'next-intl';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

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
        <footer className="w-full bg-card flex justify-center">
            <div className="flex w-full max-w-[1200px] items-center justify-between gap-4 border-t border-border px-6 py-4 text-muted-foreground min-h-[72px] max-md:flex-col max-md:items-start max-md:gap-3 max-md:py-4 max-md:pb-5">
                <nav className="flex flex-wrap items-center gap-1" aria-label="Footer links">
                    {footerLinks.map(link => (
                        <Button
                            key={link.key}
                            variant="ghost"
                            className="px-1 text-muted-foreground hover:text-foreground transition-colors"
                            onClick={() => handleFooterLinkClick(link.key)}
                        >
                            {t(link.labelKey)}
                        </Button>
                    ))}
                </nav>

                <div className="flex shrink-0 items-center gap-3 max-md:w-full max-md:justify-between">
                    <span className="text-sm whitespace-nowrap">
                        {t('copyright')}
                    </span>
                    <Separator orientation="vertical" className="h-3.5 w-px bg-border max-md:hidden"/>
                    <div className="flex items-center gap-1" aria-label="Social links">
                        {socialLinks.map(social => (
                            <Button
                                key={social.key}
                                variant="ghost"
                                size="icon"
                                className="text-2xl text-muted-foreground hover:text-primary transition-colors duration-300 hover:-translate-y-0.5"
                                aria-label={social.name}
                                onClick={() => handleSocialClick(social.name)}
                            >
                                <Icon icon={social.icon}/>
                            </Button>
                        ))}
                    </div>
                </div>
            </div>
        </footer>
    );
}
