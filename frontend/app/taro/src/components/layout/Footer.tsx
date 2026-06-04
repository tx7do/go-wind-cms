import {View, Text} from '@tarojs/components';
import {useTranslations} from '@/lib/next-intl-compat';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export default function Footer() {
    const t = useTranslations('ui');
    const router = useI18nRouter();

    const footerLinks = [
        {key: 'about_us', label: t('button.about_us'), url: '/about'},
        {key: 'contact_us', label: t('button.contact_us'), url: '/contact'},
        {key: 'non_responsibility', label: t('button.non_responsibility'), url: '/disclaimer'},
        {key: 'privacy_agreement', label: t('button.privacy_agreement'), url: '/privacy'},
        {key: 'terms_of_service', label: t('button.terms_of_service'), url: '/terms'},
    ];

    return (
        <View className='w-full bg-cardBg border-t-[1rpx] border-splitLine'>
            <View className='px-[24rpx] py-[32rpx]'>
                {/* 链接行 */}
                <View className='flex flex-wrap justify-center gap-[16rpx] mb-[24rpx]'>
                    {footerLinks.map(link => (
                        <Text
                          key={link.key}
                          className='text-desc text-textThird px-[16rpx] py-[8rpx]'
                          onClick={() => router.push(link.url)}
                        >
                            {link.label}
                        </Text>
                    ))}
                </View>
                {/* 版权 */}
                <Text className='block text-center text-tips text-textWeak'>
                    {t('copyright')}
                </Text>
            </View>
        </View>
    );
}
