import React, {useEffect, useState} from 'react';
import {View, Text} from '@tarojs/components';
import Taro from '@tarojs/taro';

import './index.scss';

export interface BackToTopProps {
    /** 滚动超过多少像素后显示，默认 500px */
    scrollThreshold?: number;
    /** 点击回调函数 */
    onClick?: () => void;
    /** 自定义样式类名 */
    className?: string;
    /** 是否可见 */
    visible?: boolean;
}

/**
 * 回到顶部按钮组件
 * 当页面滚动超过指定阈值时显示，点击后平滑滚动到页面顶部
 */
const BackToTop: React.FC<BackToTopProps> = ({
    scrollThreshold = 500,
    onClick,
    className = '',
    visible: controlledVisible
}) => {
    const [showBackToTop, setShowBackToTop] = useState(false);

    // 监听滚动事件
    useEffect(() => {
        const handleScroll = () => {
            // Taro 中获取页面滚动位置
            Taro.createSelectorQuery()
                .selectViewport()
                .scrollOffset((res) => {
                    setShowBackToTop(res.scrollTop > scrollThreshold);
                })
                .exec();
        };

        // Taro 中添加滚动监听
        const onPageScroll = () => {
            handleScroll();
        };

        Taro.getCurrentPages().forEach((page) => {
            if (page && page.onPageScroll) {
                page.onPageScroll(onPageScroll);
            }
        });

        // 初始化时检查一次
        handleScroll();
    }, [scrollThreshold]);

    // 如果传入了 controlledVisible，使用受控模式
    const isVisible = controlledVisible !== undefined ? controlledVisible : showBackToTop;

    // 滚动到顶部
    const scrollToTop = () => {
        Taro.pageScrollTo({scrollTop: 0, duration: 300});
        onClick?.();
    };

    if (!isVisible) return null;

    return (
        <View 
            className={`back-to-top ${className}`}
            onClick={scrollToTop}
            aria-label='Back to top'
        >
            <Text className='icon'>↑</Text>
        </View>
    );
};

export default BackToTop;
