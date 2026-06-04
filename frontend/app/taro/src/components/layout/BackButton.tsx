import React from 'react';
import {View} from '@tarojs/components';
import {XIcon} from '@/plugins/xicon';

const BackButton: React.FC<{
    label: string;
    onClick: () => void;
    className?: string;
}> = ({label, onClick, className}) => {
    return (
        <View
          className={`flex items-center gap-[8rpx] px-[16rpx] py-[12rpx] min-h-touch ${className || ''}`}
          onClick={onClick}
          hoverClass='tap-active'
        >
            <XIcon name='carbon:arrow-left' size={16} className='text-textSec' />
            <View className='text-desc text-textSec'>{label}</View>
        </View>
    );
};

export default BackButton;
