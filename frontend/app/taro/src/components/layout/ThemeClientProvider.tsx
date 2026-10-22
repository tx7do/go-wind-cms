import React, {useEffect, useMemo, useRef} from 'react';
import Taro from '@tarojs/taro';
import {View} from '@tarojs/components';
import {usePreferences} from '@/core/preferences';

/**
 * 主题客户端 Provider
 *
 * 1. H5：按 mode 在 <html> 上切 .dark / .light（app.css 的 html.dark token 生效）
 * 2. 小程序：<page> 挂不上 class，改在主题根 View 上内联同一份色板供子树继承，
 *    并让根 View 自己带上 .dark，app.css 的 .dark .au-* 派生规则才有后代可匹配
 * 3. 小程序：原生导航栏 / tabBar / 页面底色读不到 CSS 变量，调 API 跟主题
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

    // 色值必须与 app.css 的 :root / html.dark token 逐字一致，
    // 否则 H5 上 html 继承值与子树内联值会同时生效，出现半亮半暗
    const themeStyle = useMemo(() => {
        if (isDark) {
            return {
                '--color-primary': '#2e96ff',
                '--color-primary-strong': '#70b3ff',
                '--color-primary-tint': 'rgba(46, 150, 255, 0.16)',
                '--color-text-main': '#ffffffd9',
                '--color-text-sec': '#ffffffb3',
                '--color-text-third': '#ffffff8a',
                '--color-text-weak': '#ffffff59',
                '--color-page-bg': '#050b14',
                '--color-card-bg': '#0d1626',
                '--color-split-line': '#1e293b',
                '--color-bar-bg': 'rgba(13, 22, 38, 0.92)',
            } as React.CSSProperties;
        }
        return {
            '--color-primary': '#006be6',
            '--color-primary-strong': '#005ac2',
            '--color-primary-tint': '#e5efff',
            '--color-text-main': '#0f172a',
            '--color-text-sec': '#475569',
            '--color-text-third': '#64748b',
            '--color-text-weak': '#cbd5e1',
            '--color-page-bg': '#f4f7f9',
            '--color-card-bg': '#ffffff',
            '--color-split-line': '#e2e8f0',
            '--color-bar-bg': 'rgba(255, 255, 255, 0.92)',
        } as React.CSSProperties;
    }, [isDark]);

    // 小程序端：原生导航栏、tabBar 与页面底色由客户端绘制，读不到 CSS 变量，只能调 API 跟主题。
    // 页面底色尤其要紧：暗色下内联 token 只覆盖到根 View 自己的盒子，
    // 内容不足一屏时的下半截、以及下拉回弹区都是原生 page 底色（app.config 里的亮色基线），
    // 不调这个 API 就会从内容底下透出浅色。非 tab 页调用 setTabBarStyle 会 reject，属正常路径。
    useEffect(() => {
        if (process.env.TARO_ENV === 'h5') return;
        const page = isDark ? '#050b14' : '#f4f7f9';
        Taro.setNavigationBarColor({
            frontColor: isDark ? '#ffffff' : '#000000',
            backgroundColor: page,
        }).catch(() => {});
        Taro.setBackgroundColor({
            backgroundColor: page,
            backgroundColorTop: page,
            backgroundColorBottom: page,
        }).catch(() => {});
        Taro.setTabBarStyle({
            color: isDark ? '#8c96a6' : '#7a7e83',
            selectedColor: isDark ? '#2e96ff' : '#006be6',
            backgroundColor: isDark ? '#0d1626' : '#ffffff',
            borderStyle: 'black',
        }).catch(() => {});
    }, [isDark]);

    // 根 View 必须自己带上 dark/light class：app.css 里 .dark .au-card、.dark .au-glass、
    // .dark .au-grad-text、.dark .au-page::before 等 6 条「高级感」规则靠后代选择器生效，
    // 而小程序的 <page> 挂不上 class（H5 靠 <html>，小程序没有这一层）。
    // 只注入内联变量的话，小程序暗色下只剩换色，玻璃/极光/渐变字/骨架光全数落空。
    return (
        <View className={`app-theme-root ${isDark ? 'dark' : 'light'}`} style={themeStyle}>
            {children}
        </View>
    );
}
