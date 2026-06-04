import {View} from '@tarojs/components';
import React from 'react';

/**
 * 页面内容区通用容器
 * 移动端统一 px-[24rpx] 内边距
 */
const SectionContainer: React.FC<{
    children: React.ReactNode;
    className?: string;
}> = ({children, className}) => {
    return (
        <View className={`w-full px-[24rpx] py-[32rpx] ${className || ''}`}>
            {children}
        </View>
    );
};

export default SectionContainer;
