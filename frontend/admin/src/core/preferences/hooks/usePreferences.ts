import {useMemo} from 'react';
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
    const shortcutKeys = useMemo(() => preferences.shortcutKeys, [preferences.shortcutKeys]);
    const transition = useMemo(() => preferences.transition, [preferences.transition]);
    const copyright = useMemo(() => preferences.copyright, [preferences.copyright]);
    const navigation = useMemo(() => preferences.navigation, [preferences.navigation]);
    const tabbar = useMemo(() => preferences.tabbar, [preferences.tabbar]);
    const breadcrumb = useMemo(() => preferences.breadcrumb, [preferences.breadcrumb]);
    const sidebar = useMemo(() => preferences.sidebar, [preferences.sidebar]);
    const header = useMemo(() => preferences.header, [preferences.header]);
    const footer = useMemo(() => preferences.footer, [preferences.footer]);
    const logo = useMemo(() => preferences.logo, [preferences.logo]);
    const widget = useMemo(() => preferences.widget, [preferences.widget]);

    // 计算属性
    const isDark = useMemo(() => {
        return theme.mode === 'dark' ||
            (theme.mode === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches);
    }, [theme.mode]);

    const isMobile = useMemo(() => {
        return app.isMobile;
    }, [app.isMobile]);

    // 更新方法
    const updatePreferences = (overrides: DeepPartial<Preferences>) => {
        setPreferences(overrides);
    };

    const updateTheme = (themeOverrides: DeepPartial<Preferences['theme']>) => {
        setPreferences({theme: themeOverrides});
    };

    const updateApp = (appOverrides: DeepPartial<Preferences['app']>) => {
        setPreferences({app: appOverrides});
    };

    // 主题切换
    const toggleTheme = () => {
        const newMode = theme.mode === 'dark' ? 'light' : 'dark';
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
        shortcutKeys,
        transition,
        copyright,
        navigation,
        tabbar,
        breadcrumb,
        sidebar,
        header,
        footer,
        logo,
        widget,

        // 计算属性
        isDark,
        isMobile,

        // 通用更新方法
        updatePreferences,
        resetPreferences,
        getPreference,

        // 便捷更新方法
        updateTheme,
        updateApp,
        toggleTheme,
        setThemeMode,
        setLanguage,
    };
}
