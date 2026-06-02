'use client';

import React, {useEffect} from 'react';
import {Spinner} from '@/components/ui/spinner';
import {usePathname} from 'next/navigation';

import {useLoading} from '@/store/core/loading/store';

export default function GlobalLoading() {
    const {isLoading, finish} = useLoading();
    const pathname = usePathname();

    useEffect(() => {
        if (isLoading) {
            const timer = setTimeout(() => {
                console.warn('[GlobalLoading] Loading timeout (>5s), forcing finish');
                finish();
            }, 5000);
            return () => clearTimeout(timer);
        }
    }, [isLoading, finish]);

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
        <div className="fixed inset-0 z-[9999] flex items-center justify-center bg-white/80 backdrop-blur-sm dark:bg-black/70">
            <div className="flex flex-col items-center gap-4">
                <Spinner size="lg"/>
                <p className="mt-2 text-base font-medium text-foreground">加载中...</p>
            </div>
        </div>
    );
}
