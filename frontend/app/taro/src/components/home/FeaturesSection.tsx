import {View, Text} from '@tarojs/components';
import {useTranslations} from '@/lib/next-intl-compat';
import {XIcon} from '@/plugins/xicon';

/**
 * 首页「平台特性」区块
 * - 2列 grid 布局
 * - 图标+标题+描述
 */
export default function FeaturesSection() {
    const t = useTranslations('page.home');

    const features = [
        {icon: 'carbon:document', title: t('flexible_content_management'), description: t('content_management_desc')},
        {icon: 'carbon:cloud', title: t('multi_tenant_architecture'), description: t('multi_tenant_desc')},
        {icon: 'carbon:security', title: t('enterprise_security'), description: t('security_desc')},
        {icon: 'carbon:analytics', title: t('advanced_analytics'), description: t('analytics_desc')},
        {icon: 'carbon:api', title: t('api_integration'), description: t('api_integration_desc')},
        {icon: 'carbon:collaborate', title: t('real_time_collaboration'), description: t('real_time_collaboration_desc')},
    ];

    return (
        <View className='w-full'>
            {/* 标题行 */}
            <View className='flex items-center justify-center gap-[8rpx] mb-[24rpx]'>
                <XIcon name='carbon:rocket' size={20} className='text-primary' />
                <Text className='text-card-title font-bold text-textMain'>
                    {t('platform_features')}
                </Text>
            </View>

            {/* 2列特性网格 */}
            <View className='flex flex-wrap gap-[16rpx]'>
                {features.map((feature, index) => (
                    <View
                      key={index}
                      className='rounded bg-cardBg p-[24rpx] flex flex-col items-center text-center'
                      style={{width: 'calc(50% - 8rpx)'}}
                    >
                        <View
                          className='flex items-center justify-center w-[72rpx] h-[72rpx] rounded mb-[16rpx]'
                          style={{backgroundColor: 'rgba(22,119,255,0.08)'}}
                        >
                            <XIcon name={feature.icon || ''} size={28} className='text-primary' />
                        </View>
                        <Text className='text-desc font-bold text-textMain mb-[8rpx]'>
                            {feature.title}
                        </Text>
                        <Text className='text-tips text-textSec leading-[1.6]'>
                            {feature.description}
                        </Text>
                    </View>
                ))}
            </View>
        </View>
    );
}
