import React from 'react';
import ReactDOM from 'react-dom/client';
import { App as AntdApp } from 'antd';

import { QueryClientProvider } from '@tanstack/react-query';

// Devtools 仅在开发环境导入，生产构建时自动 tree-shake
// Devtools 仅在开发环境导入，生产构建时自动 tree-shake；
// 另设 VITE_QUERY_DEVTOOLS=false 可在 dev 下关闭（其悬浮球会遮挡右下角弹层按钮，自动化走查/录屏时可关）。
const ReactQueryDevtools =
  import.meta.env.DEV && import.meta.env.VITE_QUERY_DEVTOOLS !== 'false'
    ? await import('@tanstack/react-query-devtools').then((m) => m.ReactQueryDevtools)
    : () => null;

import { ThemeProvider } from '@/core/preferences/components/ThemeProvider';
import { queryClient } from './core';
import { bootstrap } from '@/bootstrap.ts';

// 样式
import 'nprogress/nprogress.css';

import './styles/global.css';
import './styles/tailwind.css';
import './styles/skeleton-force-dark.css';
import './styles/pro-layout-overrides.css';
import './styles/pro-components-dark.css';
import './styles/semantic-text.css';
import './styles/list-table.css';

import './index.css';
import App from './app.tsx';

// 执行全局初始化
bootstrap().then(() => {
  const root = ReactDOM.createRoot(document.getElementById('root')!);

  root.render(
    <React.StrictMode>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider>
          <AntdApp>
            <App />
          </AntdApp>
        </ThemeProvider>

        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </React.StrictMode>,
  );
});
