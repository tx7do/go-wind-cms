import React, {useEffect, useMemo, useRef} from 'react';
import {View} from '@tarojs/components';
import {usePreferences} from '@/core/preferences';

/**
 * 主题客户端 Provider
 *
 * 职责：
 * 1. 根据 preferences.theme.mode 切换暗色/亮色（CSS 变量方式）
 *    - H5: 在 <html> 上切换 .dark / .light class
 *    - 小程序: 在 page 上切换 .dark class（通过 CSS page.dark 选择器）
 * 2. 将 preferences.theme 中的动态色板注入为 CSS 自定义属性
 */
export default function ThemeClientProvider({children}: { children: React.ReactNode }) {
    const {theme, isDark} = usePreferences();
    const mode = theme.mode;
    const mqRef = useRef<MediaQueryList | null>(null);
    // H5 端：操作 document.documentElement
    useEffect(() => {
        if (typeof document === 'undefined') return;
        const html = document.documentElement;

        if (mqRef.current) {
            mqRef.current.onchange = null;
        }

        if (mode === 'auto') {
            if (typeof window === 'undefined') return;
            const mq = window.matchMedia('(prefers-color-scheme: dark)');
            mqRef.current = mq;

            const currentTheme = mq.matches ? 'dark' : 'light';
            html.classList.remove('dark', 'light');
            html.classList.add(currentTheme);

            mq.onchange = (e) => {
                const newTheme = e.matches ? 'dark' : 'light';
                html.classList.remove('dark', 'light');
                html.classList.add(newTheme);
            };
        } else {
            html.classList.remove('dark', 'light');
            html.classList.add(mode);
        }

        return () => {
            if (mqRef.current) {
                mqRef.current.onchange = null;
            }
        };
    }, [mode]);

    // 全平台：通过内联 CSS 变量注入色板（最可靠的跨端方案，不依赖 page 选择器在 H5 中的匹配）
    const themeStyle = useMemo(() => {
        if (isDark) {
            return {
                '--color-text-main': '#ffffffd9',
                '--color-text-sec': '#ffffffb3',
                '--color-text-third': '#ffffff8a',
                '--color-text-weak': '#ffffff59',
                '--color-page-bg': '#17171a',
                '--color-card-bg': '#232326',
                '--color-split-line': '#3a3a3c',
                '--color-bar-bg': 'rgba(35, 35, 38, 0.92)',
            } as React.CSSProperties;
        }
        return {
            '--color-text-main': '#1d2129',
            '--color-text-sec': '#4e5969',
            '--color-text-third': '#86909c',
            '--color-text-weak': '#c9cdd4',
            // 亮色背景：页面底色略深于纯白，卡片纯白以形成对比
            '--color-page-bg': '#ebedf0',
            '--color-card-bg': '#ffffff',
            // 分割线：更轻以避免与卡片边框叠加后显得粗重
            '--color-split-line': '#e8eaed',
            '--color-bar-bg': 'rgba(255, 255, 255, 0.92)',
        } as React.CSSProperties;
    }, [isDark]);

    return (
        <View style={themeStyle}>
            {children}
        </View>
    );
}
