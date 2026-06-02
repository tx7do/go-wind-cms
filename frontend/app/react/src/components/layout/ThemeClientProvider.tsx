'use client';

import React, {useEffect, useRef} from 'react';
import {usePreferences} from '@/core/preferences';

export default function ThemeClientProvider({children}: { children: React.ReactNode }) {
    const {theme} = usePreferences();
    const mode = theme.mode;
    const mqRef = useRef<MediaQueryList | null>(null);

    useEffect(() => {
        const html = document.documentElement;

        // 清理之前的监听器
        if (mqRef.current) {
            mqRef.current.onchange = null;
        }

        if (mode === 'auto') {
            // 跟随系统
            const mq = window.matchMedia('(prefers-color-scheme: dark)');
            mqRef.current = mq;

            // 立即应用当前系统主题
            const currentTheme = mq.matches ? 'dark' : 'light';
            html.classList.remove('dark', 'light');
            html.classList.add(currentTheme);

            mq.onchange = (e) => {
                const newTheme = e.matches ? 'dark' : 'light';
                html.classList.remove('dark', 'light');
                html.classList.add(newTheme);
            };
        } else {
            // 直接应用指定主题
            html.classList.remove('dark', 'light');
            html.classList.add(mode);
        }

        // 清理函数
        return () => {
            if (mqRef.current) {
                mqRef.current.onchange = null;
            }
        };
    }, [mode]);

    return children;
}
