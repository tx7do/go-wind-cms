'use client';

import React, {useEffect} from 'react';
import {Spin} from 'antd';
import {usePathname} from 'next/navigation';

import {useLoading} from '@/store/core/loading/store';

import styles from './GlobalLoading.module.css';

export default function GlobalLoading() {
    const {isLoading, finish} = useLoading();
    const pathname = usePathname();

    // 超时保护 - 如果 loading 超过 5 秒，自动隐藏
    useEffect(() => {
        if (isLoading) {
            const timer = setTimeout(() => {
                console.warn('[GlobalLoading] Loading timeout (>5s), forcing finish');
                finish();
            }, 5000);
            return () => clearTimeout(timer);
        }
    }, [isLoading, finish]);

    // 监听路由变化（Next.js 页面渲染完成）
    useEffect(() => {
        if (isLoading) {
            const timer = setTimeout(() => {
                console.log('[GlobalLoading] Route changed, page rendered, hiding loading');
                finish();
            }, 300);
            return () => clearTimeout(timer);
        }
    }, [pathname, isLoading, finish]);

    if (!isLoading) return null;

    return (
        <div className={styles.loadingOverlay}>
            <div className={styles.loadingContent}>
                <Spin size="large"/>
                <p className={styles.loadingText}>加载中...</p>
            </div>
        </div>
    );
}
