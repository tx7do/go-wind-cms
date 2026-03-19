import {PropsWithChildren, createElement, ReactElement} from 'react';
import {useLaunch} from '@tarojs/taro';
import {Provider} from 'react-redux';

import store from './store';
import './i18n';

import '@nutui/nutui-taro/dist/styles/themes/default.scss'
import './app.scss';

function App({children}: PropsWithChildren): ReactElement {
  useLaunch(() => {
    console.log('App launched.');
  });

  return createElement(
    Provider as any,
    {store},
    children
  ) as ReactElement;
}

export default App;
