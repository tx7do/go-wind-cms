import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';
import {usePageTitle} from '@/hooks/usePageTitle';

/**
 * 联系我们：展示站点联系方式（i18n 文案）。
 * CMS 的表单收集走「页面 FORM 区块」（FormSubmissionService，见 page/detail 渲染器），
 * 本页不重复造表单。
 */
export default function ContactPage() {
    const {t} = useTranslation();
    usePageTitle('page.title.contact');

    const channels = [
        {icon: 'carbon:email', label: t('page.legal.contact.item_1')},
        {icon: 'carbon:chat', label: t('page.legal.contact.item_2')},
        {icon: 'carbon:time', label: t('page.legal.contact.item_3')},
    ];

    return (
        <View className='au-page min-h-screen w-full bg-pageBg pb-[80rpx]'>
            <View className='px-[24rpx] pt-[36rpx]'>
                <View className='au-anim mb-[28rpx]'>
                    <Text className='au-hero-title text-textMain block'>{t('page.legal.contact.title')}</Text>
                    <Text className='text-desc text-textSec block mt-[10rpx] leading-[1.7]'>
                        {t('page.legal.contact.description')}
                    </Text>
                </View>

                <View className='au-card au-card--glow overflow-hidden au-anim au-anim-1'>
                    {channels.map((channel, index) => (
                        <View
                          key={channel.icon}
                          className={`flex items-start gap-[20rpx] px-[28rpx] py-[28rpx] ${index < channels.length - 1 ? 'border-b-[length:1rpx] border-splitLine' : ''}`}
                        >
                            <View className='w-[60rpx] h-[60rpx] rounded-[16rpx] bg-primary-tint flex items-center justify-center shrink-0'>
                                <XIcon name={channel.icon} size={26} className='text-primary' />
                            </View>
                            <Text className='text-body text-textMain leading-[1.7] flex-1'>{channel.label}</Text>
                        </View>
                    ))}
                </View>
            </View>
        </View>
    );
}
