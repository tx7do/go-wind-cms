'use client';

import {ReactQueryDevtools} from '@tanstack/react-query-devtools';
import {queryClient} from '@/core';

/**
 * React Query DevTools - 仅开发环境加载
 * 生产环境 build 时通过 next/dynamic + process.env.NODE_ENV 排除
 */
export default function ReactQueryDevtoolsPanel() {
    if (process.env.NODE_ENV !== 'development') return null;

    return (
        <ReactQueryDevtools
            client={queryClient}
            initialIsOpen={false}
            buttonPosition="bottom-right"
        />
    );
}
