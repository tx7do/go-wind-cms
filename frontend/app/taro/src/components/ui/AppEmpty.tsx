import React from 'react';
import {View} from '@tarojs/components';
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
        <XIcon name='carbon:inbox' size={64} className='text-muted-foreground' />
    );

    return (
        <View
          className={cn(
                'flex w-full items-center justify-center gap-4 py-12 px-5',
                inContainer && 'my-20',
                className,
            )}
          style={style}
        >
            <View className='opacity-50'>{defaultIcon}</View>
            {description && (
                <View className='text-sm text-muted-foreground'>{description}</View>
            )}
        </View>
    );
};

export default AppEmpty;
