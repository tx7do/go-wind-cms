import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Input} from '@tarojs/components';
import type {InputProps, CommonEventFunction} from '@tarojs/components';

import './index.scss';

export default function SearchBar() {
  const {t} = useTranslation('navbar.top');
  const [searchQuery, setSearchQuery] = useState('');

  const handleInputChange: CommonEventFunction<InputProps.inputEventDetail> = (e) => {
    setSearchQuery(e.detail.value);
  };

  const handleSearch = () => {
    console.log('Searching for:', searchQuery);
    // TODO: 实现搜索逻辑
  };

  return (
    <View className='search-bar-wrapper'>
      <View className='search-bar'>
        <Text className='search-icon'>🔍</Text>
        <Input
          value={searchQuery}
          onInput={handleInputChange}
          onConfirm={handleSearch}
          placeholder={t('search_placeholder')}
          className='search-input'
        />
      </View>
    </View>
  );
}
