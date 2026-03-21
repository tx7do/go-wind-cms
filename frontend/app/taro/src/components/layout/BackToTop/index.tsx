import React, {useEffect, useState} from 'react';
import {View} from '@tarojs/components';
import Taro from '@tarojs/taro';

import XIcon from '@/plugins/xicon';

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
          if (res && res.scrollTop !== undefined) {
            setShowBackToTop(res.scrollTop > scrollThreshold);
          }
        })
        .exec();
    };

    // 立即执行一次
    handleScroll();

    // 使用定时器定期检查滚动位置（兼容方案）
    const timer = setInterval(handleScroll, 300);

    // 尝试添加页面滚动监听
    const pages = Taro.getCurrentPages();
    const currentPage = pages[pages.length - 1];
    if (currentPage) {
      // 保存原始的 onPageScroll
      const originalOnPageScroll = (currentPage as any).onPageScroll;

      (currentPage as any).onPageScroll = function(e: { scrollTop: number }) {
        if (originalOnPageScroll) {
          originalOnPageScroll.call(this, e);
        }
        setShowBackToTop(e.scrollTop > scrollThreshold);
      };

      return () => {
        // 清理定时器
        clearInterval(timer);
        // 恢复原始 onPageScroll
        (currentPage as any).onPageScroll = originalOnPageScroll;
      };
    }

    return () => {
      clearInterval(timer);
    };
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
      <XIcon name='carbon:arrow-up' size={20} className='icon' />
    </View>
  );
};

export default BackToTop;
