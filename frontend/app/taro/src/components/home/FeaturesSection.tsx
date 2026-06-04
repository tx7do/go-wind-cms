import {View, Text} from '@tarojs/components';
import {useTranslations} from '@/lib/next-intl-compat';
import {XIcon} from '@/plugins/xicon';

/**
 * 首页「平台特性」区块
 * - 竖向排列，一行一个卡片
 */
const FEATURE_ICONS = [
    'carbon:document', 'carbon:cloud', 'carbon:security',
    'carbon:analytics', 'carbon:api', 'carbon:collaborate',
];

const FEATURE_BG_COLORS = [
    'rgba(22,119,255,0.08)', 'rgba(0,180,42,0.08)', 'rgba(255,125,0,0.08)',
    'rgba(114,46,209,0.08)', 'rgba(19,194,194,0.08)', 'rgba(235,47,150,0.08)',
];

const FEATURE_ICON_COLORS = ['#1677ff', '#00b42a', '#ff7d00', '#722ed1', '#13c2c2', '#eb2f96'];

export default function FeaturesSection() {
    const t = useTranslations('page.home');

    const features = [
        {title: t('flexible_content_management')},
        {title: t('multi_tenant_architecture')},
        {title: t('enterprise_security')},
        {title: t('advanced_analytics')},
        {title: t('api_integration')},
        {title: t('real_time_collaboration')},
    ];

    return (
        <View className='w-full'>
            {/* 区块标题 - 左对齐，与其他区块一致 */}
            <View className='flex items-center gap-[8rpx] mb-[32rpx]'>
                <XIcon name='carbon:rocket' size={18} className='text-primary' />
                <Text className='text-card-title font-bold text-textMain'>
                    {t('platform_features')}
                </Text>
            </View>

            {/* 竖向排列特性卡片 - 独立卡片质感 */}
            <View className='flex flex-col gap-[16rpx]'>
                {features.map((feature, index) => (
                    <View
                        key={index}
                        className='rounded-[16rpx] bg-cardBg flex flex-row items-center shadow-sm'
                        style={{
                            padding: '24rpx',
                        }}
                    >
                        {/* 彩色图标圆 */}
                        <View
                            className='flex items-center justify-center rounded-full flex-shrink-0'
                            style={{
                                width: '72rpx',
                                height: '72rpx',
                                backgroundColor: FEATURE_BG_COLORS[index],
                            }}
                        >
                            <XIcon
                                name={FEATURE_ICONS[index]}
                                size={28}
                                style={{color: FEATURE_ICON_COLORS[index]}}
                            />
                        </View>
                        {/* 标题 */}
                        <Text
                            className='text-desc font-semibold text-textMain ml-[16rpx] flex-1'
                            style={{
                                lineHeight: '1.5',
                            }}
                        >
                            {feature.title || ''}
                        </Text>
                    </View>
                ))}
            </View>
        </View>
    );
}
