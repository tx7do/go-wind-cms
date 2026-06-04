import {View} from '@tarojs/components';
import { useEffect, useRef, useState } from 'react';

import { useLoading } from '@/store/core/loading/store';
import { cn } from '@/lib/utils';

/**
 * 顶部路由进度条
 *
 * - 路由切换时自动展示（通过 usePathname 变化检测）
 * - 也响应 zustand loading store 的编程式触发
 * - 使用 CSS transform 动画，性能友好
 * - 不再使用全屏遮罩，避免阻塞用户交互
 */
export default function NavigationProgress() {
    const { isLoading } = useLoading();
    const pathname = typeof window !== 'undefined' ? window.location.pathname : '';

    const [progress, setProgress] = useState(0);
    const [visible, setVisible] = useState(false);

    const progressRef = useRef(0);
    const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
    const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const prevPathRef = useRef<string | null>(null);

    // 清理所有定时器
    const clearTimers = () => {
        if (intervalRef.current) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
        }
        if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
            timeoutRef.current = null;
        }
    };

    // 启动进度条
    const startProgress = () => {
        clearTimers();
        setVisible(true);
        progressRef.current = 0;
        setProgress(0);

        // 模拟进度递增（在 80% 处减缓，等待实际完成）
        intervalRef.current = setInterval(() => {
            const current = progressRef.current;
            if (current < 30) {
                progressRef.current = current + Math.random() * 8 + 2;
            } else if (current < 60) {
                progressRef.current = current + Math.random() * 4 + 1;
            } else if (current < 85) {
                progressRef.current = current + Math.random() * 2;
            } else {
                // 停留在 85%，等待完成信号
                if (intervalRef.current) {
                    clearInterval(intervalRef.current);
                    intervalRef.current = null;
                }
            }
            setProgress(progressRef.current);
        }, 200);
    };

    // 完成进度条
    const finishProgress = () => {
        clearTimers();

        // 快速推进到 100%
        progressRef.current = 100;
        setProgress(100);

        // 100% 后短暂延迟再隐藏，让用户看到完成效果
        timeoutRef.current = setTimeout(() => {
            setVisible(false);
            // 隐藏后重置进度，避免下次显示时闪跳
            setTimeout(() => {
                progressRef.current = 0;
                setProgress(0);
            }, 100);
        }, 250);
    };

    // 监听路由变化 -> 自动启动进度条
    useEffect(() => {
        if (prevPathRef.current === null) {
            // 首次渲染，记录当前路径但不触发
            prevPathRef.current = pathname;
            return;
        }
        if (prevPathRef.current !== pathname) {
            prevPathRef.current = pathname;
            startProgress();
            // 路由变化后，等待页面渲染完成
            // Next.js 客户端导航通常在 300-800ms 内完成
            timeoutRef.current = setTimeout(() => {
                finishProgress();
            }, 600);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [pathname]);

    // 监听 loading store（编程式触发）
    useEffect(() => {
        if (isLoading) {
            startProgress();
        } else if (visible) {
            finishProgress();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [isLoading]);

    // 组件卸载时清理
    useEffect(() => {
        return () => clearTimers();
    }, []);

    if (!visible && progress === 0) return null;

    return (
        <View
          className='fixed left-0 right-0 top-0 z-[9999]'
          role='progressbar'
          aria-valuenow={Math.round(progress)}
          aria-valuemin={0}
          aria-valuemax={100}
        >
            <View className='h-0.75 w-full bg-transparent'>
                <View
                  className={cn(
                        'h-full bg-primary',
                        'shadow-[0_0_10px_hsl(var(--primary)/0.6),0_0_4px_hsl(var(--primary)/0.8)]',
                        'transition-[width,opacity] duration-300 ease-out',
                    )}
                  style={{
                        width: `${progress}%`,
                        opacity: visible ? 1 : 0,
                    }}
                />
            </View>
        </View>
    );
}
