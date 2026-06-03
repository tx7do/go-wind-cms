import {PropsWithChildren, createElement, ReactElement} from 'react';
import {useLaunch} from '@tarojs/taro';

import './i18n';
import Layout from './components/layout/Layout';
import StoreProvider from './store/StoreProvider';

import './app.scss';
import './app.css';

function App({children}: PropsWithChildren): ReactElement {
  useLaunch(() => {
    console.log('App launched.');
  });

  return createElement(
    StoreProvider,
    {},
    createElement(
      Layout,
      {},
      children
    )
  ) as ReactElement;
}

export default App;
