'use client';

import dynamic from 'next/dynamic';

/**
 * 客户端入口：动态加载 DevTools，ssr: false
 * 生产环境由组件内部的 NODE_ENV 判断排除
 */
const ReactQueryDevtoolsPanel = dynamic(
    () => import('./ReactQueryDevtoolsPanel'),
    {ssr: false}
);

export default function DevToolsClient() {
    if (process.env.NODE_ENV !== 'development') return null;
    return <ReactQueryDevtoolsPanel/>;
}
