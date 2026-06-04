import {useLocale as useNextLocale} from '@/lib/next-intl-compat';

import {usePreferences} from './usePreferences';
import type {SupportedLanguagesType} from '../types/layout';

/**
 * 语言管理 Hook（偏好设置层）
 *
 * 同时整合了两方面：
 * 1. next-intl 的 URL locale（路由层）
 * 2. preferences zustand 的 locale 偏好（持久化层）
 *
 * 注意：next-intl 也有一个 useLocale，这里命名为 usePreferencesLocale 避免冲突
 */
export function usePreferencesLocale() {
    const nextLocale = useNextLocale();
    const {app, setLanguage} = usePreferences();

    // 优先使用 next-intl 的 URL locale（真实来源），fallback 到 preferences
    const locale = (nextLocale as SupportedLanguagesType) || app.locale;

    // 是否为中文
    const isZhCN = locale === 'zh-CN';

    // 是否为英文
    const isEnUS = locale === 'en-US';

    // 切换语言（仅更新偏好，不涉及路由）
    const setLocale = (newLocale: SupportedLanguagesType) => {
        setLanguage(newLocale);
    };

    // 切换到中文
    const setZhCN = () => setLocale('zh-CN');

    // 切换到英文
    const setEnUS = () => setLocale('en-US');

    // 在中英文之间切换
    const toggleLocale = () => {
        const newLocale = locale === 'zh-CN' ? 'en-US' : 'zh-CN';
        setLocale(newLocale);
    };

    // 支持的语言列表
    const supportedLocales: Array<{
        value: SupportedLanguagesType;
        label: string;
    }> = [
        {value: 'zh-CN', label: '简体中文'},
        {value: 'en-US', label: 'English'},
    ];

    return {
        // 当前语言
        locale,

        // 语言判断
        isZhCN,
        isEnUS,

        // 切换方法（偏好层）
        setLocale,
        setZhCN,
        setEnUS,
        toggleLocale,

        // 辅助
        supportedLocales,
    };
}
