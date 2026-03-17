import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import './index.scss';

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
    const {t} = useTranslation('ui');
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
        {key: 'twitter', name: 'Twitter', icon: '🐦'},
        {key: 'facebook', name: 'Facebook', icon: '📘'},
        {key: 'instagram', name: 'Instagram', icon: '📷'},
        {key: 'telegram', name: 'Telegram', icon: '✈️'},
        {key: 'wechat', name: 'WeChat', icon: '💬'},
        {key: 'weibo', name: 'Weibo', icon: '🌐'},
        {key: 'qq', name: 'QQ', icon: '🐧'},
    ];

    const handleFooterLinkClick = (key: FooterLinkKey) => {
        router.push(footerRouteMap[key]);
    };

    const handleSocialClick = (name: string) => {
        console.log(name);
    };

    return (
        <View className='footer-outer'>
            <View className='footer-container'>
                <View className='footer-links' aria-label='Footer links'>
                    {footerLinks.map(link => (
                        <Text
                            key={link.key}
                            className='footer-link'
                            onClick={() => handleFooterLinkClick(link.key)}
                        >
                            {t(link.labelKey)}
                        </Text>
                    ))}
                </View>

                <View className='footer-meta'>
                    <Text className='copyright'>
                        {t('copyright')}
                    </Text>
                    <View className='social-list' aria-label='Social links'>
                        {socialLinks.map(social => (
                            <View
                                key={social.key}
                                className='social-btn'
                                aria-label={social.name}
                                onClick={() => handleSocialClick(social.name)}
                            >
                                <Text>{social.icon}</Text>
                            </View>
                        ))}
                    </View>
                </View>
            </View>
        </View>
    );
}
