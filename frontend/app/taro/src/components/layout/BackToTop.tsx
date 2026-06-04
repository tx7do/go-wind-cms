import React, {useEffect, useState} from 'react';
import {View} from '@tarojs/components';
import Taro from '@tarojs/taro';
import XIcon from '@/plugins/xicon';
import {cn} from '@/lib/utils';

export interface BackToTopProps {
    scrollThreshold?: number;
    onClick?: () => void;
    className?: string;
    visible?: boolean;
}

const BackToTop: React.FC<BackToTopProps> = ({
    scrollThreshold = 500,
    onClick,
    className = '',
    visible: controlledVisible
}) => {
    const [showBackToTop, setShowBackToTop] = useState(false);

    useEffect(() => {
        // Taro 使用 onPageScroll 监听滚动（在页面级别）
        // H5 模式下使用 window 滚动事件
        let scrollListener: any = null;
        if (typeof window !== 'undefined') {
            const handleScroll = () => {
                setShowBackToTop(window.scrollY > scrollThreshold);
            };
            window.addEventListener('scroll', handleScroll);
            handleScroll();
            scrollListener = handleScroll;
        }
        return () => {
            if (scrollListener) {
                window.removeEventListener('scroll', scrollListener);
            }
        };
    }, [scrollThreshold]);

    const isVisible = controlledVisible !== undefined ? controlledVisible : showBackToTop;

    const scrollToTop = () => {
        Taro.pageScrollTo({scrollTop: 0, duration: 300});
        onClick?.();
    };

    if (!isVisible) return null;

    return (
        <View
          onClick={scrollToTop}
          className={cn(
                'fixed bottom-10 right-10 z-999 flex h-14 w-14 items-center justify-center',
                'rounded-full border-none bg-primary text-white cursor-pointer',
                'text-2xl shadow-lg transition-all duration-300',
                'hover:-translate-y-1 hover:shadow-xl hover:bg-primary/80',
                'active:-translate-y-0.5',
                className,
            )}
        >
            <XIcon name='carbon:arrow-up' />
        </View>
    );
};

export default BackToTop;
