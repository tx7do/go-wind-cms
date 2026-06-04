import {PropsWithChildren, createElement, ReactElement, useEffect} from 'react';
import {useLaunch} from '@tarojs/taro';

import './i18n';
import Layout from './components/layout/Layout';
import StoreProvider from './store/StoreProvider';
import './app.css';

function App({children}: PropsWithChildren): ReactElement {
  const clearAppLoading = () => {
    if (typeof document === 'undefined') return;
    document.documentElement.classList.remove('app-loading');
  };

  useLaunch(() => {
    console.log('App launched.');
    // 正常路径：Taro 启动后移除首屏 loading 标记
    clearAppLoading();
  });

  useEffect(() => {
    // 兜底路径：只要 React 完成挂载就移除 loading，避免 useLaunch 未触发时卡住
    clearAppLoading();
  }, []);

  return createElement(
    StoreProvider,
    null,
    createElement(Layout, null, children)
  ) as ReactElement;
}

export default App;
