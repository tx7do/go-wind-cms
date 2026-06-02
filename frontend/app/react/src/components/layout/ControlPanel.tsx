'use client';

import {useMemo} from 'react';

import XIcon from '@/plugins/xicon';
import {usePreferences} from '@/core/preferences';
import {useI18n} from '@/i18n';

import {cn} from '@/lib/utils';

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
        <div className={cn(
            'fixed top-8 right-8 z-[100] flex gap-3 rounded-2xl border border-border bg-card p-2.5 px-3.5',
            'backdrop-blur-md shadow-lg transition-all duration-300',
            'hover:-translate-y-0.5 hover:shadow-xl',
            'max-lg:top-4 max-lg:right-4 max-lg:p-1.5 max-lg:gap-1.5',
            'max-md:top-3 max-md:right-3 max-md:p-1 max-md:gap-1',
        )}>
            <select
                className={cn(
                    'rounded-lg border border-border bg-background px-3 py-2 text-sm font-medium',
                    'text-foreground cursor-pointer outline-none transition-all duration-300',
                    'hover:border-primary hover:bg-card',
                    'max-md:px-2 max-md:py-1.5 max-md:text-xs',
                )}
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
                className={cn(
                    'flex min-w-9 h-9 items-center justify-center gap-1.5 rounded-lg',
                    'bg-background border border-border text-foreground cursor-pointer',
                    'transition-all duration-300',
                    'hover:-translate-y-0.5 hover:shadow-md hover:border-primary hover:bg-card',
                    'max-md:min-w-8 max-md:h-8 max-md:text-sm',
                )}
                onClick={toggleTheme}
                aria-label="Toggle theme"
            >
                <XIcon name={isDark ? 'carbon:sun' : 'carbon:moon'} size={18}/>
            </button>
        </div>
    );
}
