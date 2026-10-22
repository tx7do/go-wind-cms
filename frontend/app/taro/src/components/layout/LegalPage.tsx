import {View, Text} from '@tarojs/components';
import {XIcon} from '@/plugins/xicon';

interface LegalPageProps {
    icon: string;
    title: string;
    description: string;
    items: string[];
}

export default function LegalPage({icon, title, description, items}: LegalPageProps) {
    return (
        <View className='au-page w-full bg-pageBg min-h-screen pb-[64rpx]'>
            {/* 页面标题栏 */}
            <View className='px-[32rpx] pt-[48rpx] pb-[32rpx] au-anim'>
                <View className='w-[96rpx] h-[96rpx] rounded-[28rpx] bg-primary-tint flex items-center justify-center mb-[22rpx]'>
                    <XIcon name={icon} size={44} className='text-primary' />
                </View>
                <Text className='text-title font-bold text-textMain block mb-[10rpx]'>{title}</Text>
                <Text className='text-desc text-textSec leading-[1.7]'>{description}</Text>
            </View>

            {/* 内容列表 */}
            <View className='px-[24rpx]'>
                <View className='au-card au-card--glow p-[32rpx] au-anim au-anim-1'>
                    {items.map((item, index) => (
                        <View
                          key={index}
                          className={`flex items-start gap-[18rpx] py-[20rpx] ${index < items.length - 1 ? 'border-b-[length:1rpx] border-splitLine' : ''}`}
                        >
                            <View className='shrink-0 w-[44rpx] h-[44rpx] rounded-[12rpx] bg-primary-tint flex items-center justify-center'>
                                <Text className='text-tips font-bold text-primary'>
                                    {index + 1}
                                </Text>
                            </View>
                            <Text className='text-body text-textMain leading-[1.8] flex-1'>{item}</Text>
                        </View>
                    ))}
                </View>
            </View>
        </View>
    );
}
