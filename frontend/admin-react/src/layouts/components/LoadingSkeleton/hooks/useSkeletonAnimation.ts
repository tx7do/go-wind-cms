import { useState, useEffect, useRef } from 'react';
import type { SkeletonAnimation } from '../Skeleton/types';

interface UseSkeletonAnimationOptions {
  /** 是否激活动画 */
  active?: boolean;

  /** 动画类型 */
  animation?: SkeletonAnimation;

  /** 延迟显示时间（避免闪烁） */
  delay?: number;
}

/**
 * 骨架屏动画控制 Hook
 * 功能：延迟显示 + 动画类名生成 + 性能优化
 */
export const useSkeletonAnimation = ({
  active = true,
  animation = 'wave',
  delay = 0,
}: UseSkeletonAnimationOptions) => {
  const [visible, setVisible] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // 延迟显示（避免快速加载时的闪烁）
  useEffect(() => {
    if (delay > 0) {
      timerRef.current = setTimeout(() => {
        setVisible(true);
      }, delay);
    } else {
      setVisible(true);
    }

    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, [delay]);

  // 生成动画类名
  const animationClass = active && visible ? `animate-${animation}` : '';

  // 清理
  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, []);

  return {
    visible: active && visible,
    animationClass,
  };
};
