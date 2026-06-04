import {View, Text} from '@tarojs/components';
import {cn} from '@/lib/utils';
import {XIcon} from '@/plugins/xicon';

interface LegalPageProps {
    icon: string;
    title: string;
    description: string;
    items: string[];
}

export default function LegalPage({icon, title, description, items}: LegalPageProps) {
    return (
        <View className='w-full bg-pageBg'>
            {/* 页面标题栏 */}
            <View className='bg-cardBg px-[24rpx] py-[32rpx] border-b-[1rpx] border-splitLine'>
                <View className='flex items-center gap-[16rpx] mb-[16rpx]'>
                    <XIcon name={icon} size={24} className='text-primary' />
                    <Text className='text-title font-bold text-textMain'>{title}</Text>
                </View>
                <Text className='text-desc text-textSec'>{description}</Text>
            </View>

            {/* 内容列表 */}
            <View className='px-[24rpx] py-[32rpx]'>
                <View className='rounded bg-cardBg p-[24rpx]'>
                    {items.map((item, index) => (
                        <View
                          key={index}
                          className='flex items-start gap-[16rpx] py-[16rpx]'
                        >
                            <Text className='shrink-0 text-desc font-bold text-primary min-w-[48rpx]'>
                                {String(index + 1).padStart(2, '0')}.
                            </Text>
                            <Text className='text-body text-textMain leading-relaxed'>{item}</Text>
                        </View>
                    ))}
                </View>
            </View>
        </View>
    );
}
