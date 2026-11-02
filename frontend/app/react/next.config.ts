import type {NextConfig} from "next";
import createNextIntlPlugin from 'next-intl/plugin';

const withNextIntl = createNextIntlPlugin('./src/i18n/request.ts');

const nextConfig: NextConfig = {
    // 启用 React 严格模式
    reactStrictMode: true,

    // 静态导出模式（2026-11-02）：产物为纯静态站点（nginx 直接托管，无 node 服务）。
    //   - 全站客户端取数：详情页 id/slug 解析、SEO title 均在客户端完成；
    //   - trailingSlash: true 导出 目录/index.html 形态，nginx try_files $uri/ 直接命中；
    //   - middleware（src/proxy.ts）在 export 下不支持，已移除，locale 协商由
    //     根 page.tsx 客户端重定向兜底（裸深链可在 nginx 层按前缀 307 补）。
    // 如需恢复 SSR 模式（爬虫可抓服务端正文与 meta）：去掉 output/trailingSlash、
    // 恢复 src/proxy.ts，并把详情页改回服务端组件取数（历史实现见 git log）。
    output: 'export',
    trailingSlash: true,

    distDir: 'dist',
};

export default withNextIntl(nextConfig);
