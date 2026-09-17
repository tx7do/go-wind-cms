import {useCallback, useMemo} from 'react';

import zhCN from 'antd/locale/zh_CN';
import enUS from 'antd/locale/en_US';
import type {Locale} from 'antd/es/locale';
import {useTranslation} from 'react-i18next';

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
 * 语言切换 Hook
 * 提供便捷的语言切换和管理功能，同时支持 Ant Design 国际化
 */
export function useLocale() {
    const {app, updateApp, preferences} = usePreferences();

    // 当前语言
    const locale = app.locale;

    // 是否为中文
    const isZhCN = locale === 'zh-CN';

    // 是否为英文
    const isEnUS = locale === 'en-US';

    // Ant Design 语言包
    const antdLocale = useMemo(() => {
        return antdLocales[locale];
    }, [locale]);

    // 切换语言
    const setLocale = useCallback((newLocale: SupportedLanguagesType) => {
        updateApp({locale: newLocale});

        // 可选：刷新页面以应用语言变更
        // window.location.reload();
    }, [updateApp]);

    // 切换到中文
    const setZhCN = useCallback(() => {
        setLocale('zh-CN');
    }, [setLocale]);

    // 切换到英文
    const setEnUS = useCallback(() => {
        setLocale('en-US');
    }, [setLocale]);

    // 切换语言（在中英文之间切换）
    const toggleLocale = useCallback(() => {
        const newLocale = locale === 'zh-CN' ? 'en-US' : 'zh-CN';
        setLocale(newLocale);
    }, [locale, setLocale]);

    const {t} = useTranslation('auth');

    // 获取语言显示名称
    const getLocaleName = useCallback((lang: SupportedLanguagesType): string => {
        const names: Record<SupportedLanguagesType, string> = {
            'zh-CN': t('langZhCN'),
            'en-US': t('langEnUS'),
        };
        return names[lang];
    }, [t]);

    // 当前语言的显示名称
    const localeName = getLocaleName(locale);

    // 获取所有支持的语言列表
    const supportedLocales: Array<{
        value: SupportedLanguagesType;
        label: string;
    }> = [
        {value: 'zh-CN', label: t('langZhCN')},
        {value: 'en-US', label: t('langEnUS')},
    ];

    return {
        // 当前语言
        locale,
        localeName,

        // Ant Design 语言包
        antdLocale,

        // 语言判断
        isZhCN,
        isEnUS,

        // 切换方法
        setLocale,
        setZhCN,
        setEnUS,
        toggleLocale,

        // 辅助方法
        getLocaleName,
        supportedLocales,

        // 完整偏好设置（供高级使用）
        preferences,
    };
}
