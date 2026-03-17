import React from 'react';
import {View, Text, Picker} from '@tarojs/components';

import {useI18n} from '@/i18n';

import './index.scss';

/**
 * 语言切换选择器组件
 *
 * 通常用于 Header 或导航栏中，供用户切换界面语言
 */
export const LocaleSwitcher: React.FC = () => {
  const {changeLocale} = useI18n();

  const localeOptions = [
    {key: 'zh-CN', label: '简体中文'},
    {key: 'en-US', label: 'English'},
  ];

  const handleLanguageChange = (value: string) => {
    changeLocale(value);
  };

  return (
    <Picker
      mode='selector'
      range={localeOptions}
      rangeKey='label'
      onChange={(e) => handleLanguageChange(localeOptions[e.detail.value].key)}
    >
      <View className='lang-switcher' aria-label='Language switcher'>
        <Text className='lang-icon'>🌐</Text>
      </View>
    </Picker>
  );
};
