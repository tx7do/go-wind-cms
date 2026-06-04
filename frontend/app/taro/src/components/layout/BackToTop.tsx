import React, {useEffect, useState} from 'react';
import {View} from '@tarojs/components';
import Taro from '@tarojs/taro';
import XIcon from '@/plugins/xicon';

const BackToTop: React.FC<{
    scrollThreshold?: number;
    onClick?: () => void;
}> = ({scrollThreshold = 500, onClick}) => {
    const [visible, setVisible] = useState(false);

    useEffect(() => {
        if (typeof window === 'undefined') return;
        const handleScroll = () => setVisible(window.scrollY > scrollThreshold);
        window.addEventListener('scroll', handleScroll);
        handleScroll();
        return () => window.removeEventListener('scroll', handleScroll);
    }, [scrollThreshold]);

    if (!visible) return null;

    return (
        <View
          className='fixed bottom-[40rpx] right-[40rpx] z-[999] flex w-[88rpx] h-[88rpx] items-center justify-center rounded-full bg-primary'
          onClick={() => {
                Taro.pageScrollTo({scrollTop: 0, duration: 300});
                onClick?.();
            }}
          hoverClass='tap-active'
        >
            <XIcon name='carbon:arrow-up' size={20} className='text-white' />
        </View>
    );
};

export default BackToTop;
