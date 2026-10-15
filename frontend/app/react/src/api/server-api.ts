import type { contentservicev1_Post, contentservicev1_Page } from '@/api/generated/app/service/v1';

/**
 * 服务端数据获取（仅可在 Server Component / generateMetadata 中调用）。
 *
 * 直连 app 服务（默认 http://localhost:6700，可用 API_INTERNAL_URL 覆盖），
 * 不经过浏览器侧的 axios 实例（其依赖 window 上下文）。
 * 公开端点仅返回已发布内容；匿名可访问，无需鉴权。
 *
 * 缓存策略：revalidate 60s —— 内容新鲜度与渲染性能的折中；
 * 帖子编辑后最迟 60 秒反映到 SEO HTML。
 */

const API_INTERNAL = (process.env.API_INTERNAL_URL ?? 'http://localhost:6700').replace(/\/$/, '');

async function serverGet<T>(path: string, revalidate = 60): Promise<T | null> {
  try {
    const res = await fetch(`${API_INTERNAL}${path}`, {
      next: { revalidate },
      headers: { Accept: 'application/json' },
    });
    if (!res.ok) return null;
    return (await res.json()) as T;
  } catch {
    // 后端不可达：返回 null，页面呈现"未找到"态而非 500
    return null;
  }
}

export interface ServerSearchParams {
  id?: string;
  slug?: string;
  locale?: string;
}

/** 服务端按 id 获取已发布帖子（SSR 详情页用） */
export function serverFetchPostById(id: number, locale?: string): Promise<contentservicev1_Post | null> {
  return serverGet<contentservicev1_Post>(`/app/v1/posts/${id}?locale=${encodeURIComponent(locale ?? '')}`);
}

/** 服务端按 id 获取已发布 CMS 页面（SSR ?id= 路由用） */
export function serverFetchPageById(id: number, locale?: string): Promise<contentservicev1_Page | null> {
  return serverGet<contentservicev1_Page>(`/app/v1/pages/${id}?locale=${encodeURIComponent(locale ?? '')}`);
}

/** 服务端按别名（code）获取已发布帖子（SEO 路由 /post/[slug] 用） */
export function serverFetchPostBySlug(slug: string, locale?: string): Promise<contentservicev1_Post | null> {
  return serverGet<contentservicev1_Post>(`/app/v1/posts/by-slug/${encodeURIComponent(slug)}?locale=${encodeURIComponent(locale ?? '')}`);
}

/** 服务端按 slug 获取已发布 CMS 页面（SSR 页面路由用） */
export function serverFetchPageBySlug(slug: string, locale?: string): Promise<contentservicev1_Page | null> {
  return serverGet<contentservicev1_Page>(`/app/v1/pages/by-slug/${encodeURIComponent(slug)}?locale=${encodeURIComponent(locale ?? '')}`);
}

/** 从帖子 DTO 提取指定语言的翻译（匹配 → 第一条回退，与客户端约定一致） */
export function pickTranslation<T extends { languageCode?: string }>(
  translations: T[] | undefined,
  locale: string,
): T | null {
  if (!translations || translations.length === 0) return null;
  return translations.find((t) => t.languageCode === locale) ?? translations[0];
}

/** 正文去 Markdown 符号后截取 description（粗略，够 meta 用） */
export function toMetaDescription(content: string, max = 160): string {
  const plain = content
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/!?\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/[#>*`~\-|]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim();
  return plain.length > max ? `${plain.slice(0, max - 1)}…` : plain;
}
