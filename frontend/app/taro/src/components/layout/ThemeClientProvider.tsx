import React, {useEffect, useRef} from 'react';
import {usePreferences} from '@/core/preferences';

/**
 * 主题客户端 Provider
 *
 * 职责：
 * 1. 根据 preferences.theme.mode 切换 .dark / .light class（暗色/亮色/跟随系统）
 * 2. 将 preferences.theme 中的动态色板注入为 CSS 自定义属性，
 *    使全站 Tailwind/shadcn 组件实时响应偏好设置变更
 */
export default function ThemeClientProvider({children}: { children: React.ReactNode }) {
    const {theme, isDark} = usePreferences();
    const mode = theme.mode;
    const mqRef = useRef<MediaQueryList | null>(null);

    // ========== 1. 切换 dark/light class ==========
    useEffect(() => {
        if (typeof document === 'undefined') return;
        const html = document.documentElement;

        // 清理之前的监听器
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

    // ========== 2. 动态注入主题色 CSS 变量 ==========
    useEffect(() => {
        if (typeof document === 'undefined') return;
        const root = document.documentElement;

        // 主色调 → --primary / --ring
        root.style.setProperty('--primary', theme.colorPrimary);

        // 根据暗色/亮色自动计算 primary-foreground（绿底按钮上的文字必须是纯白色）
        root.style.setProperty('--primary-foreground', '0 0% 100%');

        // 成功色
        root.style.setProperty('--success', theme.colorSuccess);
        root.style.setProperty('--success-foreground', '0 0% 98%');

        // 警告色
        root.style.setProperty('--warning', theme.colorWarning);
        root.style.setProperty('--warning-foreground', '0 0% 98%');

        // 错误/危险色
        root.style.setProperty('--destructive', theme.colorDestructive);
        root.style.setProperty('--destructive-foreground', '0 0% 98%');

        // 圆角
        root.style.setProperty('--radius', theme.radius);

        // ring 颜色跟随 primary
        root.style.setProperty('--ring', theme.colorPrimary);

        // accent 色保持赛博科技蓝
        root.style.setProperty('--accent', '217.2 91.2% 59.8%');
        root.style.setProperty('--accent-foreground', isDark ? '210 20% 98%' : '224 71.4% 4.1%');
    }, [
        theme.colorPrimary,
        theme.colorSuccess,
        theme.colorWarning,
        theme.colorDestructive,
        theme.radius,
        isDark,
    ]);

    return children;
}
