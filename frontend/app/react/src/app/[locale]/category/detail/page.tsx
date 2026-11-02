'use client';

import {Suspense} from 'react';

import CategoryDetailPage from './client-page';

// client-page 使用 useSearchParams：静态导出 prerender 必须有 Suspense 边界，
// 否则整页 CSR bailout 报 missing-suspense-with-csr-bailout。
export default function Page() {
    return (
        <Suspense fallback={null}>
            <CategoryDetailPage />
        </Suspense>
    );
}
