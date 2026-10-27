import type { Metadata } from 'next';
import { Suspense } from 'react';
import PageDetailClient from './client-page';
import {
    serverFetchPageById,
    serverFetchPageBySlug,
    pickTranslation,
    toMetaDescription,
    type ServerSearchParams,
} from '@/api/server-api';
import { defaultLocale, locales } from '@/i18n/config';
import type { contentservicev1_Page } from '@/api/generated/app/service/v1';

/**
 * CMS 页面详情（SSR）：服务端取数 + 爬虫可见 meta + initialPage 下传客户端。
 *
 * 路由：
 *   /[locale]/page/detail?slug=about   —— SEO 别名路由（by-slug 端点）
 *   /[locale]/page/detail?id=5         —— 兼容 id 路由
 *
 * 服务端直连 app 服务（API_INTERNAL_URL），仅返回已发布内容；
 * 未发布/不存在 → initialPage 为 null，客户端组件渲染 not-found 态。
 */

type Props = {
    params: Promise<{ locale: string }>;
    searchParams: Promise<ServerSearchParams>;
};

function resolveLocale(raw: string | undefined): string {
    return raw && (locales as readonly string[]).includes(raw) ? raw : defaultLocale;
}

async function fetchPageBySearchParams(searchParams: ServerSearchParams, locale: string): Promise<contentservicev1_Page | null> {
    if (searchParams.slug) {
        return serverFetchPageBySlug(searchParams.slug, locale);
    }
    const id = Number(searchParams.id);
    if (Number.isInteger(id) && id > 0) {
        return serverFetchPageById(id, locale);
    }
    return null;
}

export async function generateMetadata({ params, searchParams }: Props): Promise<Metadata> {
    const { locale: localeParam } = await params;
    const sp = await searchParams;
    const locale = resolveLocale(localeParam);
    const page = await fetchPageBySearchParams(sp, locale);
    if (!page) return { title: 'GoWind Content Hub' };

    const tr = pickTranslation(page.translations, locale);
    const title = tr?.title || 'GoWind Content Hub';
    // 页面正文在 sections（区块）里：拼接当前语言的区块文本生成 description
    const sectionText = (page.sections ?? [])
        .map((sec) => {
            const st = pickTranslation(sec.translations, locale);
            return Object.values(st?.content ?? {}).join(' ');
        })
        .join(' ');
    const description = toMetaDescription(`${title} ${sectionText}`.trim());
    // SEO 字段优先（admin 编辑侧 seoTitle/metaDescription/og* /canonical；翻译表 seo JSON）
    const seo = tr?.seo;
    const finalTitle = seo?.seoTitle || title;
    const finalDescription = seo?.metaDescription || description;
    const ogImage = seo?.ogImage || page.thumbnail || undefined;

    return {
        title: finalTitle,
        description: finalDescription,
        keywords: seo?.metaKeywords || undefined,
        alternates: {
            canonical: seo?.canonicalUrl
                || (sp.slug ? `/page/detail?slug=${sp.slug}&lang=${locale}` : `/page/detail?id=${page.id}&lang=${locale}`),
        },
        openGraph: {
            title: seo?.ogTitle || finalTitle,
            description: seo?.ogDescription || finalDescription,
            images: ogImage ? [ogImage] : undefined,
            type: 'website',
        },
    };
}

export default async function Page({ params, searchParams }: Props) {
    const { locale: localeParam } = await params;
    const sp = await searchParams;
    const locale = resolveLocale(localeParam);

    const page = await fetchPageBySearchParams(sp, locale);

    return (
        <Suspense fallback={null}>
            <PageDetailClient initialPage={page} />
        </Suspense>
    );
}
