import type {NextRequest} from 'next/server';

import handleI18nRouting from '@/i18n/middleware';

// next-intl 请求处理入口：Next 16 约定文件名为 proxy.ts（旧名 middleware.ts 已弃用）。
// 2026-03 起本层曾被移除、仅靠根 page.tsx 客户端重定向兜底 `/`，裸深链
// （如 admin 预览链接 /post/detail?...）因此直接 404；恢复本层后裸路径
// 统一经 locale 协商 307 到带前缀路由。
// config/matcher 必须在本文件静态声明，Next 不允许 re-export。
export default function proxy(request: NextRequest) {
    return handleI18nRouting(request);
}

export const config = {
    // 匹配全部页面路径：localePrefix=always 下裸路径（如 admin 预览链接
    // /post/detail）无路由会 404，必须经本层 307 到带前缀路由；
    // 排除 api、_next 与带扩展名的静态资源（favicon/robots.txt 等）
    matcher: ['/', '/((?!api|_next|.*\\..*).*)'],
};
