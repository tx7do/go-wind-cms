'use client';

import React, {useEffect, useState} from 'react';
import XIcon from '@/plugins/xicon';

import {cn} from '@/lib/utils';

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
            setShowBackToTop(window.scrollY > scrollThreshold);
        };

        window.addEventListener('scroll', handleScroll);
        // 初始化时检查一次
        handleScroll();

        return () => window.removeEventListener('scroll', handleScroll);
    }, [scrollThreshold]);

    // 如果传入了 controlledVisible，使用受控模式
    const isVisible = controlledVisible !== undefined ? controlledVisible : showBackToTop;

    // 滚动到顶部
    const scrollToTop = () => {
        window.scrollTo({top: 0, behavior: 'smooth'});
        onClick?.();
    };

    if (!isVisible) return null;

    return (
        <button
            onClick={scrollToTop}
            className={cn(
                'fixed bottom-10 right-10 z-[999] flex h-14 w-14 items-center justify-center',
                'rounded-full border-none bg-primary text-white cursor-pointer',
                'text-2xl shadow-lg transition-all duration-300',
                'hover:-translate-y-1 hover:shadow-xl hover:bg-primary/80',
                'active:translate-y-[-2px]',
                className,
            )}
            aria-label="Back to top"
        >
            <XIcon name="carbon:arrow-up"/>
        </button>
    );
};

export default BackToTop;
