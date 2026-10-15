import type {NextConfig} from "next";
import createNextIntlPlugin from 'next-intl/plugin';

const withNextIntl = createNextIntlPlugin('./src/i18n/request.ts');

const nextConfig: NextConfig = {
    // 启用 React 严格模式
    reactStrictMode: true,

    // SSR 模式（Node 服务部署：next build && next start -p $PORT）。
    // 公开内容页（帖子/CMS 页面详情）为服务端组件，爬虫可直接抓取正文与 meta；
    // 认证与交互页仍为客户端组件。静态导出模式如需恢复，加回 output: 'export'
    // 并将 server 组件取数路径改回客户端 useEffect 拉取。
    // trailingSlash: true, // SSR 模式下关闭尾斜杠强制，canonical 由 generateMetadata 输出

    distDir: 'dist',
};

export default withNextIntl(nextConfig);
