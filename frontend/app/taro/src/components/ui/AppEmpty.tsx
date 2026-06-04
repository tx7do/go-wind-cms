import React from 'react';
import {View, Text} from '@tarojs/components';
import {XIcon} from '@/plugins/xicon';
import {cn} from '@/lib/utils';

interface AppEmptyProps {
    variant?: 'default' | 'error' | 'noData';
    inContainer?: boolean;
    image?: React.ReactNode;
    description?: React.ReactNode;
    className?: string;
    style?: React.CSSProperties;
}

export const AppEmpty: React.FC<AppEmptyProps> = ({
    inContainer = false,
    image,
    description,
    className,
    style,
}) => {
    const defaultIcon = image || (
        <XIcon name='carbon:inbox' size={48} className='text-textWeak' />
    );

    return (
        <View
          className={cn(
                'flex w-full flex-col items-center justify-center gap-[24rpx] py-[96rpx] px-[24rpx]',
                inContainer && 'my-[120rpx]',
                className,
            )}
          style={style}
        >
            <View className='opacity-50'>{defaultIcon}</View>
            {description && (
                <Text className='text-desc text-textThird'>{description}</Text>
            )}
        </View>
    );
};

export default AppEmpty;
