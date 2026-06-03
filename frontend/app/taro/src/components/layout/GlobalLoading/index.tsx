import {useEffect} from 'react';
import {View, Text} from '@tarojs/components';
import {useLoadingStore} from '@/store/core/loading/store';

import './index.scss';

export default function GlobalLoading() {
  const isLoading = useLoadingStore((state) => state.isLoading);
  const finish = useLoadingStore((state) => state.finish);

  // 超时保护 - 如果 loading 超过 5 秒，自动隐藏
  useEffect(() => {
    if (isLoading) {
      const timer = setTimeout(() => {
        console.warn('[GlobalLoading] Loading timeout (>5s), forcing finish');
        finish();
      }, 5000);

      return () => clearTimeout(timer);
    }
  }, [isLoading, finish]);

  if (!isLoading) return null;

  return (
    <View className='loading-overlay'>
      <View className='loading-content'>
        <View className='spinner' />
        <Text className='loading-text'>加载中...</Text>
      </View>
    </View>
  );
}
