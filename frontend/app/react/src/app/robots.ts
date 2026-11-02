import type { MetadataRoute } from 'next';

/**
 * robots.txt：全放行 + sitemap 指向。
 *
 * sitemap 本体由 app 服务提供（GET /app/v1/sitemap.xml，含已发布帖子/页面
 * 的多语言 URL 条目）。SITE_BASE_URL 为站点对外基址（与 sitemap 生成共用，
 * 未配置回退 dev 端口 15001）。
 */

const SITE_BASE_URL = (process.env.SITE_BASE_URL ?? 'http://localhost:15001').replace(/\/$/, '');

// output: 'export' 要求 route handler 显式声明 force-static（构建期生成静态文件）
export const dynamic = 'force-static';

export default function robots(): MetadataRoute.Robots {
  return {
    rules: { userAgent: '*', allow: '/' },
    sitemap: `${SITE_BASE_URL}/sitemap.xml`,
  };
}
