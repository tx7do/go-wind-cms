import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {usePageTitle} from '@/hooks/usePageTitle';
import {useI18nRouter} from '@/i18n/helpers';

export default function AboutPage() {
    const {t} = useTranslation();
    usePageTitle('page.title.about');
    const router = useI18nRouter();

    const features = [
        {icon: 'carbon:document-add', title: t('page.about.feature_content'), description: t('page.about.feature_content_desc')},
        {icon: 'carbon:cloud-upload', title: t('page.about.feature_multi_tenant'), description: t('page.about.feature_multi_tenant_desc')},
        {icon: 'carbon:security', title: t('page.about.feature_security'), description: t('page.about.feature_security_desc')},
        {icon: 'carbon:api', title: t('page.about.feature_api'), description: t('page.about.feature_api_desc')},
        {icon: 'carbon:collaborate', title: t('page.about.feature_collaboration'), description: t('page.about.feature_collaboration_desc')},
        {icon: 'carbon:analytics', title: t('page.about.feature_analytics'), description: t('page.about.feature_analytics_desc')},
    ];

    return (
        <View className='au-page min-h-screen w-full bg-pageBg pb-[80rpx]'>
            {/* 品牌头 */}
            <View
              className='relative overflow-hidden px-[32rpx] pt-[44rpx] pb-[64rpx]'
              style={{background: 'linear-gradient(128deg, #57a3ff 0%, #1677ff 52%, #0b55d4 100%)'}}
            >
                <View
                  className='absolute rounded-full'
                  style={{
                      width: '360rpx',
                      height: '360rpx',
                      right: '-110rpx',
                      top: '-150rpx',
                      background: 'radial-gradient(circle, rgba(255,255,255,0.26), transparent 68%)',
                  }}
                />
                <Text className='au-hero-title text-white relative block'>{t('page.about.title')}</Text>
                <Text className='text-desc text-white/85 relative block mt-[10rpx]'>{t('page.about.subtitle')}</Text>
            </View>

            <View className='px-[24rpx] -mt-[36rpx]'>
                {/* 关于我们 */}
                <View className='au-card au-card--glow p-[32rpx] au-anim'>
                    <Text className='text-card-title font-bold text-textMain block mb-[16rpx]'>{t('page.about.about_us')}</Text>
                    <Text className='text-body text-textSec leading-[1.8] block mb-[16rpx]'>{t('page.about.about_us_desc_1')}</Text>
                    <Text className='text-body text-textSec leading-[1.8] block mb-[16rpx]'>{t('page.about.about_us_desc_2')}</Text>
                    <Text className='text-body text-textSec leading-[1.8]'>{t('page.about.about_us_desc_3')}</Text>
                </View>

                {/* 功能特性 */}
                <View className='au-card p-[32rpx] au-anim au-anim-1 mt-[24rpx]'>
                    <Text className='text-card-title font-bold text-textMain block mb-[8rpx]'>{t('page.about.features')}</Text>
                    <Text className='text-desc text-textSec mb-[32rpx]'>{t('page.about.features_desc')}</Text>
                    <View className='grid grid-cols-2 gap-[24rpx]'>
                        {features.map((feature) => (
                            <View key={feature.title} className='rounded-[24rpx] bg-pageBg border-[length:1rpx] border-splitLine p-[24rpx]'>
                                <View className='flex items-center justify-center w-[72rpx] h-[72rpx] rounded-[22rpx] bg-primary-tint mb-[16rpx]'>
                                    <XIcon name={feature.icon} size={28} className='text-primary' />
                                </View>
                                <Text className='text-desc font-bold text-textMain block mb-[8rpx]'>{feature.title}</Text>
                                <Text
                                  className='text-tips text-textSec leading-[1.6]'
                                  style={{
                                      overflow: 'hidden',
                                      textOverflow: 'ellipsis',
                                      display: '-webkit-box',
                                      WebkitLineClamp: 3,
                                      WebkitBoxOrient: 'vertical',
                                  }}
                                >
                                    {feature.description}
                                </Text>
                            </View>
                        ))}
                    </View>
                </View>

                {/* 联系我们 CTA */}
                <View
                  className='au-card mt-[24rpx] p-[32rpx] flex items-center justify-between au-anim au-anim-2'
                  hoverClass='tap-active'
                  onClick={() => router.push('/contact')}
                >
                    <View>
                        <Text className='text-body font-bold text-textMain block'>{t('page.about.cta_contact')}</Text>
                        <Text className='text-tips text-textThird block mt-[6rpx]'>{t('page.contact.description')}</Text>
                    </View>
                    <XIcon name='carbon:arrow-right' size={22} className='text-primary shrink-0' />
                </View>
            </View>
        </View>
    );
}
