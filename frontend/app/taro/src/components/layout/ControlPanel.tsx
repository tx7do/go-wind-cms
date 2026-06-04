import {View, Picker, Text} from '@tarojs/components';

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
        <View className={cn(
            'fixed top-8 right-8 z-100 flex gap-3 rounded-2xl border border-border bg-card p-2.5 px-3.5',
            'backdrop-blur-md shadow-lg transition-all duration-300',
            'hover:-translate-y-0.5 hover:shadow-xl',
            'max-lg:top-4 max-lg:right-4 max-lg:p-1.5 max-lg:gap-1.5',
            'max-md:top-3 max-md:right-3 max-md:p-1 max-md:gap-1',
        )}
        >
            <Picker
              mode='selector'
              range={languageOptions.map(o => o.label)}
              value={languageOptions.findIndex(o => o.key === locale)}
              onChange={(e) => {
                const idx = e.detail.value;
                if (languageOptions[idx]) handleSelectLanguage(languageOptions[idx].key);
              }}
            >
                <View className={cn(
                    'flex items-center rounded-lg border border-border bg-background px-3 py-2 text-sm font-medium',
                    'text-foreground cursor-pointer transition-all duration-300',
                    'hover:border-primary hover:bg-card',
                    'max-md:px-2 max-md:py-1.5 max-md:text-xs',
                )}
                >
                    <Text>{languageOptions.find(o => o.key === locale)?.label || locale}</Text>
                    <XIcon name='carbon:chevron-down' size={14} className='ml-1 opacity-60' />
                </View>
            </Picker>
            <View
              className={cn(
                    'flex min-w-9 h-9 items-center justify-center gap-1.5 rounded-lg',
                    'bg-background border border-border text-foreground cursor-pointer',
                    'transition-all duration-300',
                    'hover:-translate-y-0.5 hover:shadow-md hover:border-primary hover:bg-card',
                    'max-md:min-w-8 max-md:h-8 max-md:text-sm',
                )}
              onClick={toggleTheme}
              aria-label='Toggle theme'
            >
                <XIcon key={isDark ? 'sun' : 'moon'} name={isDark ? 'carbon:sun' : 'carbon:moon'} size={18} className='theme-icon-animate' />
            </View>
        </View>
    );
}
