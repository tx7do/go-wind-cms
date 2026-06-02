import CategoryDetailPage from './client-page';

/**
 * 静态导出：返回占位 ID 让构建通过
 * Nginx try_files 将所有 /category/xxx/ 回退到 locale index.html，客户端路由接管
 */
export function generateStaticParams() {
    return [{id: '0'}];
}

export default function Page() {
    return <CategoryDetailPage />;
}
