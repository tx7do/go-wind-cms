'use client';

import {useMemo} from 'react';

import XIcon from '@/plugins/xicon';
import {usePreferences} from '@/core/preferences';
import {useI18n} from '@/i18n';

import styles from './ControlPanel.module.css';

export default function ControlPanel() {
    const {isDark, toggleTheme} = usePreferences();
    const {locale, changeLocale} = useI18n();

    // 语言选项
    const languageOptions = useMemo(() => [
        {key: 'zh-CN', label: '中文'},
        {key: 'en-US', label: 'English'}
    ], []);

    // 切换语言
    const handleSelectLanguage = (key: string) => {
        changeLocale(key);
    };

    return (
        <div className={styles['control-panel']}>
            <select
                className={styles['language-select']}
                value={locale}
                onChange={(e) => handleSelectLanguage(e.target.value)}
            >
                {languageOptions.map((option) => (
                    <option key={option.key} value={option.key}>
                        {option.label}
                    </option>
                ))}
            </select>
            <button
                className={`${styles['control-btn']} ${styles['theme-toggle']}`}
                onClick={toggleTheme}
                aria-label="Toggle theme"
            >
                <XIcon name={isDark ? 'carbon:sun' : 'carbon:moon'} size={18}/>
            </button>
        </div>
    );
}
