import {useMemo, useSyncExternalStore} from 'react';
import {usePreferencesStore} from '../store';
import type {DeepPartial, Preferences} from '../types';

/**
 * 偏好设置 Hook
 * 提供便捷的偏好设置访问和更新方法
 */
export function usePreferences() {
    const {
        preferences,
        setPreferences,
        resetPreferences,
        getPreference,
    } = usePreferencesStore();

    // 便捷访问器
    const app = useMemo(() => preferences.app, [preferences.app]);
    const theme = useMemo(() => preferences.theme, [preferences.theme]);
    const content = useMemo(() => preferences.content, [preferences.content]);
    const copyright = useMemo(() => preferences.copyright, [preferences.copyright]);
    const widget = useMemo(() => preferences.widget, [preferences.widget]);
    const transition = useMemo(() => preferences.transition, [preferences.transition]);

    // 响应式跟踪系统暗色模式（auto 模式需要）
    // H5: 使用 window.matchMedia
    // 小程序: 使用 Taro.getSystemInfoSync().theme
    const subscribeSystemDark = (callback: () => void) => {
        if (typeof window !== 'undefined') {
            const mq = window.matchMedia('(prefers-color-scheme: dark)');
            mq.addEventListener('change', callback);
            return () => mq.removeEventListener('change', callback);
        }
        // 小程序端：Taro.onThemeChange
        try {
            const Taro = require('@tarojs/taro');
            Taro.onThemeChange(callback as any);
            return () => Taro.offThemeChange(callback as any);
        } catch {
            return () => {};
        }
    };
    const getSystemIsDark = () => {
        if (typeof window !== 'undefined') {
            return window.matchMedia('(prefers-color-scheme: dark)').matches;
        }
        // 小程序端
        try {
            const Taro = require('@tarojs/taro');
            const systemInfo = Taro.getSystemInfoSync();
            return systemInfo.theme === 'dark';
        } catch {
            return false;
        }
    };
    const systemIsDark = useSyncExternalStore(subscribeSystemDark, getSystemIsDark, () => false);

    // 计算属性：是否暗色模式（auto 模式下跟随系统）
    const isDark = useMemo(() => {
        if (theme.mode === 'dark') return true;
        if (theme.mode === 'light') return false;
        // auto: 跟随系统
        return systemIsDark;
    }, [theme.mode, systemIsDark]);

    const isMobile = useMemo(() => app.isMobile, [app.isMobile]);

    // 通用更新方法
    const updatePreferences = (overrides: DeepPartial<Preferences>) => {
        setPreferences(overrides);
    };

    const updateApp = (appOverrides: DeepPartial<Preferences['app']>) => {
        setPreferences({app: appOverrides});
    };

    const updateTheme = (themeOverrides: DeepPartial<Preferences['theme']>) => {
        setPreferences({theme: themeOverrides});
    };

    const updateContent = (contentOverrides: DeepPartial<Preferences['content']>) => {
        setPreferences({content: contentOverrides});
    };

    const updateWidget = (widgetOverrides: DeepPartial<Preferences['widget']>) => {
        setPreferences({widget: widgetOverrides});
    };

    // 主题切换（在 dark/light 之间切换，不影响 auto）
    const toggleTheme = () => {
        const newMode = isDark ? 'light' : 'dark';
        updateTheme({mode: newMode});
    };

    const setThemeMode = (mode: Preferences['theme']['mode']) => {
        updateTheme({mode});
    };

    // 语言切换
    const setLanguage = (locale: Preferences['app']['locale']) => {
        updateApp({locale});
    };

    return {
        // 完整偏好设置
        preferences,

        // 分组访问
        app,
        theme,
        content,
        copyright,
        widget,
        transition,

        // 计算属性
        isDark,
        isMobile,

        // 通用更新方法
        updatePreferences,
        resetPreferences,
        getPreference,

        // 便捷更新方法
        updateApp,
        updateTheme,
        updateContent,
        updateWidget,
        toggleTheme,
        setThemeMode,
        setLanguage,
    };
}
