import type {NextConfig} from "next";
import createNextIntlPlugin from 'next-intl/plugin';

const withNextIntl = createNextIntlPlugin('./src/i18n/request.ts');

const nextConfig: NextConfig = {
    // 启用 React 严格模式
    reactStrictMode: true,

    // 静态导出模式（部署到 Nginx / CDN / 对象存储）
    output: 'export',

    // 所有 URL 以 / 结尾，生成 /path/index.html 目录结构
    // Nginx try_files 可以方便地做 fallback
    trailingSlash: true,
};

export default withNextIntl(nextConfig);
