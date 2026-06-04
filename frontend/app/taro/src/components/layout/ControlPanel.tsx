import {View, Picker, Text} from '@tarojs/components';
import {useMemo} from 'react';
import XIcon from '@/plugins/xicon';
import {usePreferences} from '@/core/preferences';
import {useI18n} from '@/i18n';

export default function ControlPanel() {
    const {isDark, toggleTheme} = usePreferences();
    const {locale, changeLocale} = useI18n();

    const languageOptions = useMemo(() => [
        {key: 'zh-CN', label: '中文'},
        {key: 'en-US', label: 'English'}
    ], []);

    return (
        <View className='flex gap-[16rpx] p-[16rpx]'>
            <Picker
              mode='selector'
              range={languageOptions.map(o => o.label)}
              value={languageOptions.findIndex(o => o.key === locale)}
              onChange={(e) => {
                    const idx = e.detail.value;
                    if (languageOptions[idx]) changeLocale(languageOptions[idx].key);
                }}
            >
                <View className='flex items-center rounded border-[1rpx] border-splitLine bg-cardBg px-[24rpx] h-[64rpx] text-desc'>
                    <Text className='text-desc text-textMain'>
                        {languageOptions.find(o => o.key === locale)?.label || locale}
                    </Text>
                    <XIcon name='carbon:chevron-down' size={14} className='ml-[8rpx] text-textThird' />
                </View>
            </Picker>
            <View
              className='flex items-center justify-center w-[64rpx] h-[64rpx] rounded border-[1rpx] border-splitLine bg-cardBg'
              onClick={toggleTheme}
            >
                <XIcon name={isDark ? 'carbon:sun' : 'carbon:moon'} size={18} className='text-textSec' />
            </View>
        </View>
    );
}
