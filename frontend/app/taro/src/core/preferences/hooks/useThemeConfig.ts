import {useMemo, useSyncExternalStore} from 'react';

import {usePreferencesStore, type PreferencesState} from '../store';

/**
 * 主题配置 Hook（shadcn/ui 版本，无 antd 依赖）
 *
 * 返回 Tailwind CSS 变量风格的主题配置信息，
 * 由 globals.css 中的 .dark / .light class 自动切换。
 */
export interface ThemeConfig {
    /** 是否暗色模式 */
    isDark: boolean;
    /** 是否紧凑模式 */
    isCompact: boolean;
    /** 主色调 */
    colorPrimary: string;
    /** 圆角 */
    radius: string;
}

export const useThemeConfig = (): ThemeConfig => {
    const {theme, app} = usePreferencesStore((state: PreferencesState) => state.preferences);

    // 响应式跟踪系统暗色模式（auto 模式需要）
    const subscribeSystemDark = (callback: () => void) => {
        if (typeof window === 'undefined') return () => {};
        const mq = window.matchMedia('(prefers-color-scheme: dark)');
        mq.addEventListener('change', callback);
        return () => mq.removeEventListener('change', callback);
    };

    const getSystemIsDark = () => {
        if (typeof window === 'undefined') return false;
        return window.matchMedia('(prefers-color-scheme: dark)').matches;
    };

    const systemIsDark = useSyncExternalStore(
        subscribeSystemDark,
        getSystemIsDark,
        () => false,
    );

    return useMemo((): ThemeConfig => {
        let isDark: boolean;
        if (theme.mode === 'dark') isDark = true;
        else if (theme.mode === 'light') isDark = false;
        else isDark = systemIsDark; // auto

        return {
            isDark,
            isCompact: app.compact,
            colorPrimary: theme.colorPrimary,
            radius: theme.radius,
        };
    }, [theme.mode, theme.colorPrimary, theme.radius, app.compact, systemIsDark]);
};
