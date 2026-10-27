import type { Metadata } from 'next';
import type { contentservicev1_Post } from '@/api/generated/app/service/v1';
import { Suspense } from 'react';
import PostDetailPage from './client-page';
import {
    serverFetchPostById,
    serverFetchPostBySlug,
    pickTranslation,
    toMetaDescription,
    type ServerSearchParams,
} from '@/api/server-api';
import {defaultLocale, locales} from '@/i18n/config';

/**
 * 帖子详情（SSR）：服务端取数 + 爬虫可见 meta + initialData 下传客户端。
 *
 * 路由：
 *   /post/detail?id={id}&lang={locale}   —— 兼容既有 query 路由
 *   /post/detail?slug={code}&lang={locale} —— SEO 别名路由（by-slug 端点）
 *
 * 服务端直连 app 服务（API_INTERNAL_URL），仅返回已发布内容；
 * 未发布/不存在 → initialPost 为 null，客户端组件渲染 not-found 态。
 */

type Props = {
    params: Promise<{ locale: string }>;
    searchParams: Promise<ServerSearchParams>;
};

/** 服务端解析 locale（next-intl requestLocale 在 meta 中不可用，直接取参数校验） */
function resolveLocale(raw: string | undefined): string {
    return raw && (locales as readonly string[]).includes(raw) ? raw : defaultLocale;
}

async function fetchPost(params: Params['params'], searchParams: ServerSearchParams): Promise<contentservicev1_Post | null> {
    const { locale: localeParam } = await params;
    const locale = resolveLocale(localeParam);
    if (searchParams.slug) {
        return serverFetchPostBySlug(searchParams.slug, locale);
    }
    const id = Number(searchParams.id);
    if (Number.isInteger(id) && id > 0) {
        return serverFetchPostById(id, locale);
    }
    return null;
}

type Params = { params: Promise<{ locale: string }> };

export async function generateMetadata({ params, searchParams }: Props): Promise<Metadata> {
    const { locale: localeParam } = await params;
    const sp = await searchParams;
    const locale = resolveLocale(localeParam);
    const post = await fetchPost(params, sp);
    if (!post) return { title: 'GoWind Content Hub' };

    const tr = pickTranslation(post.translations, locale);
    // SEO 字段优先（admin 编辑侧 seoTitle/metaDescription/og* /canonical；翻译表 seo JSON）
    const seo = tr?.seo;
    const title = seo?.seoTitle || tr?.title || 'GoWind Content Hub';
    const description = seo?.metaDescription || toMetaDescription(tr?.summary || tr?.content || '');
    const ogImage = seo?.ogImage || post.thumbnail || undefined;
    const keywords = seo?.metaKeywords || undefined;

    return {
        title,
        description,
        keywords,
        alternates: {
            canonical: seo?.canonicalUrl
                || (sp.slug ? `/post/detail?slug=${sp.slug}&lang=${locale}` : `/post/detail?id=${post.id}&lang=${locale}`),
        },
        openGraph: {
            title: seo?.ogTitle || title,
            description: seo?.ogDescription || description,
            images: ogImage ? [ogImage] : undefined,
            type: 'article',
        },
    };
}

export default async function Page({ params, searchParams }: Props) {
    const { locale: localeParam } = await params;
    const sp = await searchParams;
    const locale = resolveLocale(localeParam);

    const post =
        sp.slug
            ? await serverFetchPostBySlug(sp.slug, locale)
            : Number(sp.id) > 0
                ? await serverFetchPostById(Number(sp.id), locale)
                : null;

    return (
        <Suspense fallback={null}>
            <PostDetailPage initialPost={post} />
        </Suspense>
    );
}
