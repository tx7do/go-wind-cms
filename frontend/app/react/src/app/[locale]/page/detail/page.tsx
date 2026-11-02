import {Suspense} from 'react';

import PageDetailClient from './client-page';

/**
 * CMS 页面详情（静态导出）：
 *   /[locale]/page/detail?slug=about   —— SEO 别名路由（by-slug 端点）
 *   /[locale]/page/detail?id=5         —— 兼容 id 路由
 *
 * output: 'export' 下无服务端取数（middleware/generateMetadata 对动态
 * searchParams 不可用），slug/id 解析与 SEO title 均在客户端完成
 * （client-page 拉取后 document.title 回写）。SPA meta 由布局层静态输出。
 */
// client-page 使用 useSearchParams：静态导出 prerender 必须有 Suspense 边界，
// 否则整页 CSR bailout 报 missing-suspense-with-csr-bailout。
export default function Page() {
    return (
        <Suspense fallback={null}>
            <PageDetailClient />
        </Suspense>
    );
}
