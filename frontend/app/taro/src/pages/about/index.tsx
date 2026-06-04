import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';

export default function AboutPage() {
    const {t} = useTranslation();

    const features = [
        {icon: 'carbon:document-add', title: t('page.about.feature_content'), description: t('page.about.feature_content_desc')},
        {icon: 'carbon:cloud-upload', title: t('page.about.feature_multi_tenant'), description: t('page.about.feature_multi_tenant_desc')},
        {icon: 'carbon:security', title: t('page.about.feature_security'), description: t('page.about.feature_security_desc')},
        {icon: 'carbon:api', title: t('page.about.feature_api'), description: t('page.about.feature_api_desc')},
        {icon: 'carbon:collaborate', title: t('page.about.feature_collaboration'), description: t('page.about.feature_collaboration_desc')},
        {icon: 'carbon:analytics', title: t('page.about.feature_analytics'), description: t('page.about.feature_analytics_desc')},
    ];

    const stats = [
        {number: '10K+', label: t('page.about.stat_users')},
        {number: '500+', label: t('page.about.stat_projects')},
        {number: '99.9%', label: t('page.about.stat_uptime')},
        {number: '24/7', label: t('page.about.stat_support')},
    ];

    return (
        <View className='w-full bg-pageBg'>
            {/* 页面标题 */}
            <View className='bg-cardBg px-[24rpx] py-[32rpx] border-b-[1rpx] border-splitLine'>
                <Text className='text-title font-bold text-textMain block mb-[8rpx]'>{t('page.about.title')}</Text>
                <Text className='text-desc text-textSec block'>{t('page.about.subtitle')}</Text>
            </View>

            {/* 关于我们 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <Text className='text-card-title font-bold text-textMain block mb-[16rpx]'>{t('page.about.about_us')}</Text>
                <Text className='text-body text-textSec leading-relaxed block mb-[16rpx]'>{t('page.about.about_us_desc_1')}</Text>
                <Text className='text-body text-textSec leading-relaxed block mb-[16rpx]'>{t('page.about.about_us_desc_2')}</Text>
                <Text className='text-body text-textSec leading-relaxed'>{t('page.about.about_us_desc_3')}</Text>
            </View>

            {/* 统计数据 */}
            <View className='grid grid-cols-4 gap-[16rpx] px-[24rpx] pb-[32rpx]'>
                {stats.map((stat) => (
                    <View key={stat.label} className='rounded bg-cardBg p-[16rpx] text-center'>
                        <Text className='text-card-title font-bold text-primary block'>{stat.number}</Text>
                        <Text className='text-tips text-textThird block'>{stat.label}</Text>
                    </View>
                ))}
            </View>

            {/* 功能特性 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <Text className='text-card-title font-bold text-textMain block mb-[8rpx]'>{t('page.about.features')}</Text>
                <Text className='text-desc text-textSec block mb-[24rpx]'>{t('page.about.features_desc')}</Text>
                <View className='grid grid-cols-2 gap-[16rpx]'>
                    {features.map((feature) => (
                        <View key={feature.title} className='rounded bg-cardBg p-[24rpx]'>
                            <View className='flex items-center justify-center w-[72rpx] h-[72rpx] rounded bg-primary/10 mb-[16rpx]'>
                                <XIcon name={feature.icon} size={28} className='text-primary' />
                            </View>
                            <Text className='text-desc font-bold text-textMain block mb-[8rpx]'>{feature.title}</Text>
                            <Text className='text-tips text-textSec leading-relaxed'>{feature.description}</Text>
                        </View>
                    ))}
                </View>
            </View>
        </View>
    );
}
