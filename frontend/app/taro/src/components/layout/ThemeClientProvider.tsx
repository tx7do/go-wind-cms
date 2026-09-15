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
    // 色值与 app.css token 及 docs/design-language.md 保持同源
    const themeStyle = useMemo(() => {
        if (isDark) {
            return {
                '--color-primary': '#2e96ff',
                '--color-primary-tint': 'rgba(46, 150, 255, 0.16)',
                '--color-text-main': '#ffffffe6',
                '--color-text-sec': '#ffffffcc',
                '--color-text-third': '#ffffffa3',
                '--color-text-weak': '#ffffff6b',
                '--color-page-bg': '#050b14',
                '--color-card-bg': '#0d1626',
                '--color-split-line': '#1e293b',
                '--color-bar-bg': 'rgba(13, 22, 38, 0.92)',
            } as React.CSSProperties;
        }
        return {
            '--color-primary': '#006be6',
            '--color-primary-tint': '#e5efff',
            '--color-text-main': '#0f172a',
            '--color-text-sec': '#475569',
            '--color-text-third': '#64748b',
            '--color-text-weak': '#cbd5e1',
            // 亮色背景：页面底色略深于纯白，卡片纯白以形成对比
            '--color-page-bg': '#f4f7f9',
            '--color-card-bg': '#ffffff',
            // 分割线：更轻以避免与卡片边框叠加后显得粗重
            '--color-split-line': '#e2e8f0',
            '--color-bar-bg': 'rgba(255, 255, 255, 0.92)',
        } as React.CSSProperties;
    }, [isDark]);

    return (
        <View style={themeStyle}>
            {children}
        </View>
    );
}
