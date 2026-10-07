import {useMemo, useState, useEffect} from 'react';
import {theme as antdTheme, type ThemeConfig, type MappingAlgorithm} from 'antd';

import {darkThemeComponents, darkThemeTokens} from '../config';
import {SECONDARY_TEXT_DARK, SECONDARY_TEXT_LIGHT} from '../config/constants';
import {mixColor} from '@/utils/color';
import {usePreferencesStore} from '../store';

export const useThemeConfig = (): ThemeConfig => {
    const {theme, app} = usePreferencesStore((state) => state.preferences);

    // 响应式跟踪系统暗色模式（仅 auto 模式需要）
    const [systemIsDark, setSystemIsDark] = useState(
        () => window.matchMedia('(prefers-color-scheme: dark)').matches,
    );

    useEffect(() => {
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
        const handler = (e: MediaQueryListEvent) => setSystemIsDark(e.matches);
        mediaQuery.addEventListener('change', handler);
        return () => mediaQuery.removeEventListener('change', handler);
    }, []);

    // 计算有效主题：auto 模式下跟随系统偏好
    const effectiveIsDark = useMemo(() => {
        if (theme.mode === 'dark') return true;
        if (theme.mode === 'light') return false;
        if (theme.mode === 'auto') return systemIsDark;
        return false;
    }, [theme.mode, systemIsDark]);

    return useMemo((): ThemeConfig => {
        const algorithms: MappingAlgorithm[] = [];
        if (effectiveIsDark) algorithms.push(antdTheme.darkAlgorithm);
        if (app.compact) algorithms.push(antdTheme.compactAlgorithm);

        const tokens: ThemeConfig['token'] = {
            colorPrimary: theme.colorPrimary,
            colorSuccess: theme.colorSuccess,
            colorWarning: theme.colorWarning,
            colorError: theme.colorDestructive,
            // 与全局 ThemeProvider 同步：内层 ConfigProvider 会重算 token，
            // 不写这两行则该子树的链接/info 回到 antd 默认 #1677FF（另一支蓝）。
            colorInfo: theme.colorPrimary,
            // 与全局 ThemeProvider 同步（同一套推导，否则内层子树又回到 2.90:1 的暗色主色文字）
            colorLink: effectiveIsDark ? mixColor(theme.colorPrimary, '#ffffff', 0.7) : theme.colorPrimary,
            colorTextDescription: effectiveIsDark ? SECONDARY_TEXT_DARK : SECONDARY_TEXT_LIGHT,
            borderRadius: Number.parseInt(theme.radius) || 6,
        };
        // 与全局 ThemeProvider 保持同一份暗色色阶，避免内层 ConfigProvider 覆盖外层
        if (effectiveIsDark) Object.assign(tokens, darkThemeTokens);

        return {
            algorithm: algorithms.length > 0 ? algorithms : antdTheme.defaultAlgorithm,
            token: tokens,
            components: effectiveIsDark ? darkThemeComponents : undefined,
        };
    }, [effectiveIsDark, theme, app]);
};
