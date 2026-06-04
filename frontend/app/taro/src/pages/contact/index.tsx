import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';
import {usePageTitle} from '@/hooks/usePageTitle';

export default function ContactPage() {
    const {t} = useTranslation();
    usePageTitle('page.title.contact');

    return (
        <View className='w-full bg-pageBg'>
            {/* 页面标题 */}
            <View className='bg-cardBg px-[24rpx] py-[32rpx] border-b-[1rpx] border-splitLine'>
                <Text className='text-title font-bold text-textMain block mb-[8rpx]'>{t('page.legal.contact.title')}</Text>
                <Text className='text-desc text-textSec'>{t('page.legal.contact.description')}</Text>
            </View>

            {/* 联系信息 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <View className='rounded bg-cardBg p-[24rpx]'>
                    <View className='py-[16rpx] border-b-[1rpx] border-splitLine'>
                        <Text className='text-body text-textMain'>{t('page.legal.contact.item_1')}</Text>
                    </View>
                    <View className='py-[16rpx] border-b-[1rpx] border-splitLine'>
                        <Text className='text-body text-textMain'>{t('page.legal.contact.item_2')}</Text>
                    </View>
                    <View className='py-[16rpx]'>
                        <Text className='text-body text-textMain'>{t('page.legal.contact.item_3')}</Text>
                    </View>
                </View>
            </View>
        </View>
    );
}
