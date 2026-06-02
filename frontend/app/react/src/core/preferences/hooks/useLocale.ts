import {useMemo} from 'react';
import {useLocale as useNextLocale} from 'next-intl';
import zhCN from 'antd/locale/zh_CN';
import enUS from 'antd/locale/en_US';
import type {Locale} from 'antd/es/locale';

import {usePreferences} from './usePreferences';
import type {SupportedLanguagesType} from '../types/layout';

/**
 * Ant Design 语言包映射
 */
const antdLocales: Record<SupportedLanguagesType, Locale> = {
    'zh-CN': zhCN,
    'en-US': enUS,
};

/**
 * 语言管理 Hook（偏好设置层）
 *
 * 同时整合了三方面：
 * 1. next-intl 的 URL locale（路由层）
 * 2. preferences zustand 的 locale 偏好（持久化层）
 * 3. Ant Design 的语言包（UI 层）
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

    // Ant Design 语言包
    const antdLocale = useMemo(() => {
        return antdLocales[locale] ?? zhCN;
    }, [locale]);

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
        antdLocale,

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
