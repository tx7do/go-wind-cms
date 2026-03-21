import {useEffect} from 'react';
import {View, Text} from '@tarojs/components';
import {useDispatch} from 'react-redux';

import {useLoading} from '@/store/core/loading/hooks';
import {finishLoading} from '@/store/core/loading/slice';

import './index.scss';

export default function GlobalLoading() {
  const {isLoading} = useLoading();
  const dispatch = useDispatch();

  // 超时保护 - 如果 loading 超过 5 秒，自动隐藏
  useEffect(() => {
    if (isLoading) {
      const timer = setTimeout(() => {
        console.warn('[GlobalLoading] Loading timeout (>5s), forcing finish');
        dispatch(finishLoading());
      }, 5000);

      return () => clearTimeout(timer);
    }
  }, [isLoading, dispatch]);

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
