import {View, Text} from '@tarojs/components';
import XIcon from '@/plugins/xicon';

export interface OptionSheetOption {
    value: string;
    label: string;
    icon?: string;
}

interface OptionSheetProps {
    visible: boolean;
    title: string;
    options: OptionSheetOption[];
    /** 当前选中值，用于右侧勾选与高亮 */
    value: string;
    onSelect: (value: string) => void;
    onClose: () => void;
}

/** 一组互斥取值的底部弹层（主题/语言）；选中只回调不自行关闭，便于主题即时预览与调用方统一收尾 */
export function OptionSheet({visible, title, options, value, onSelect, onClose}: OptionSheetProps) {
    if (!visible) return null;

    return (
        <View className='fixed top-0 left-0 right-0 bottom-0 z-[100]'>
            <View className='au-sheet-mask' onClick={onClose} />
            <View className='au-sheet'>
                <View className='flex items-center justify-between px-[36rpx] h-[112rpx] border-b-[length:1rpx] border-splitLine'>
                    <Text className='text-card-title font-bold text-textMain'>{title}</Text>
                    <View
                      className='w-[60rpx] h-[60rpx] flex items-center justify-center rounded-full bg-pageBg'
                      hoverClass='tap-active'
                      onClick={onClose}
                    >
                        <XIcon name='carbon:close' size={20} className='text-textSec' />
                    </View>
                </View>

                {options.map((option, index) => {
                    const active = option.value === value;
                    return (
                        <View
                          key={option.value}
                          className={`flex items-center gap-[24rpx] px-[36rpx] min-h-touch py-[20rpx] ${index < options.length - 1 ? 'border-b-[length:1rpx] border-splitLine' : ''}`}
                          hoverClass='tap-active'
                          onClick={() => onSelect(option.value)}
                        >
                            {option.icon && (
                                <View className={`w-[64rpx] h-[64rpx] rounded-[16rpx] flex items-center justify-center ${active ? 'bg-primary-tint' : 'bg-pageBg'}`}>
                                    <XIcon name={option.icon} size={28} className={active ? 'text-primary' : 'text-textThird'} />
                                </View>
                            )}
                            <Text className={`flex-1 text-body ${active ? 'text-primary font-bold' : 'text-textMain'}`}>
                                {option.label}
                            </Text>
                            {active && <XIcon name='carbon:checkmark' size={24} className='text-primary' />}
                        </View>
                    );
                })}
            </View>
        </View>
    );
}

export default OptionSheet;
