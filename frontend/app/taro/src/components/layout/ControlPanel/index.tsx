import {useState, useMemo} from 'react';
import {View, Text, Picker} from '@tarojs/components';

import {useI18n} from '@/i18n';

import './index.scss';

// 获取初始主题状态（仅在客户端调用）
function getInitialTheme(): boolean {
  if (typeof window === 'undefined') {
    return false;
  }
  return document.documentElement.classList.contains('dark');
}

export default function ControlPanel() {
  const {changeLocale} = useI18n();

  const [isDark, setIsDark] = useState<boolean>(getInitialTheme);

  // 语言选项
  const languageOptions = useMemo(() => [
    {key: 'zh-CN', label: '中文'},
    {key: 'en-US', label: 'English'}
  ], []);

  // 切换语言
  const handleSelectLanguage = (value: string) => {
    changeLocale(value);
  };

  // 切换主题
  const handleToggleTheme = () => {
    const newIsDark = !document.documentElement.classList.contains('dark');
    setIsDark(newIsDark);
    if (newIsDark) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  };

  return (
    <View className='control-panel'>
      <Picker
        mode='selector'
        range={languageOptions}
        rangeKey='label'
        onChange={(e) => handleSelectLanguage(languageOptions[e.detail.value].key)}
      >
        <View className='language-select'>
          <Text>🌐</Text>
          <Text>{languageOptions.find(opt => opt.key === (typeof window !== 'undefined' && document.documentElement.lang) || 'zh-CN')?.label}</Text>
        </View>
      </Picker>
      <View
        className={`control-btn theme-toggle ${isDark ? 'dark' : ''}`}
        onClick={handleToggleTheme}
        aria-label='Toggle theme'
      >
        <Text>{isDark ? '☀️' : '🌙'}</Text>
      </View>
    </View>
  );
}
