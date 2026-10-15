'use client';

import React, {useMemo} from 'react';
import {cn} from '@/lib/utils';
import ContentViewer from '@/components/content/ContentViewer';
import {Button} from '@/components/ui/button';
import Image from '@/components/ui/image';

import type {
    contentservicev1_Section,
    contentservicev1_SectionTranslation,
} from '@/api/generated/app/service/v1';
import {currentLocaleLanguageCode} from '@/i18n';

/**
 * Section 内容键契约（三端渲染器共用，改键先同步 vue/flutter 端与本文档）
 *
 * 区块语言相关内容存 translations[].content（map<string,string>），缺译回退：
 * 匹配当前语言 → 第一条翻译。语言无关样式存 config（map<string,string>）。
 *
 * 主内容键一律为 `value`（admin 页面编辑器写入的既定键），其余为可选键：
 *
 * | 类型                | value                     | 可选内容键 | config 键                                    |
 * |---------------------|---------------------------|-----------|----------------------------------------------|
 * | RICH_TEXT           | HTML                      | -         | className                                    |
 * | MARKDOWN            | Markdown                  | -         | className                                    |
 * | TITLE               | 标题文本                   | -         | level(1-6), align(left|center|right), className |
 * | IMAGE               | 图片 URL                  | -         | alt, caption, href, className                |
 * | GALLERY / CAROUSEL  | 图片 URL 数组（JSON 或换行分隔） | -      | columns(2-4), className                      |
 * | VIDEO               | 视频地址（mp4/webm 直链或站点链接） | -     | poster, className                            |
 * | BUTTON              | 按钮文字                   | -         | href, target(_blank), variant(primary|outline), className |
 * | DIVIDER             | -                         | -         | className                                    |
 * | SPACER              | -                         | -         | height(px 数值), className                   |
 * | CODE                | 代码文本                   | -         | language(shiki 语言名), className             |
 * | HTML                | 原始 HTML（渲染层消毒）     | -         | className                                    |
 * | FORM                | 表单定义 JSON             | -         | -（提交后端待 form_submission 落地，前台渲染占位） |
 * | CUSTOM              | 任意文本（按 Markdown 兜底渲染） | -     | className                                    |
 */

/** 取区块当前语言的翻译：匹配当前语言 → 回退第一条 */
function pickSectionTranslation(
    section: contentservicev1_Section,
): contentservicev1_SectionTranslation | null {
    if (!section?.translations || section.translations.length === 0) return null;
    const locale = currentLocaleLanguageCode();
    return (
        section.translations.find(
            (tr: contentservicev1_SectionTranslation) => tr.languageCode === locale,
        ) || section.translations[0]
    );
}

function getContent(section: contentservicev1_Section): string {
    return pickSectionTranslation(section)?.content?.value || '';
}

/** GALLERY/CAROUSEL 的图片列表：value 支持两种形态——JSON 数组（URL 或 {src,caption}）或换行分隔的 URL 列表 */
function getImageList(section: contentservicev1_Section): {src: string; caption?: string}[] {
    const raw = getContent(section).trim();
    if (!raw) return [];
    if (raw.startsWith('[')) {
        try {
            const parsed = JSON.parse(raw);
            if (Array.isArray(parsed)) {
                return parsed
                    .map((item: unknown) => {
                        if (typeof item === 'string') return {src: item};
                        if (item && typeof item === 'object') {
                            const obj = item as {src?: string; url?: string; caption?: string};
                            const src = obj.src || obj.url || '';
                            return src ? {src, caption: obj.caption} : null;
                        }
                        return null;
                    })
                    .filter((x): x is {src: string; caption?: string} => x !== null);
            }
        } catch {
            // JSON 解析失败时按纯 URL 文本处理，不中断整页渲染
        }
    }
    return raw
        .split(/\r?\n/)
        .map((line) => line.trim())
        .filter(Boolean)
        .map((src) => ({src}));
}

const TITLE_ALIGN_CLASS: Record<string, string> = {
    left: 'text-start',
    center: 'text-center',
    right: 'text-end',
};

/** YouTube / Bilibili 等站点链接转嵌入 iframe；直链视频走 <video> */
function toVideoEmbed(url: string): string | null {
    try {
        const u = new URL(url);
        if (u.hostname.includes('youtube.com')) {
            const id = u.searchParams.get('v');
            if (id) return `https://www.youtube.com/embed/${id}`;
        }
        if (u.hostname === 'youtu.be') {
            return `https://www.youtube.com/embed${u.pathname}`;
        }
        if (u.hostname.includes('bilibili.com')) {
            const id = u.searchParams.get('bvid');
            if (id) return `https://player.bilibili.com/player.html?bvid=${id}`;
        }
    } catch {
        // 非 URL 按 <video> 直链处理
    }
    return null;
}

interface SectionProps {
    section: contentservicev1_Section;
}

function RichTextSection({section}: SectionProps) {
    return <ContentViewer content={getContent(section)} type="html"/>;
}

function MarkdownSection({section}: SectionProps) {
    return <ContentViewer content={getContent(section)} type="markdown"/>;
}

function TitleSection({section}: SectionProps) {
    const config = section.config ?? {};
    const level = Math.min(Math.max(Number(config.level) || 2, 1), 6);
    const align = TITLE_ALIGN_CLASS[config.align] || '';
    const Tag = `h${level}` as keyof React.JSX.IntrinsicElements;
    const sizeClass = ['text-4xl', 'text-3xl', 'text-2xl', 'text-xl', 'text-lg', 'text-base'][level - 1];
    return (
        <Tag className={cn('font-bold text-foreground', sizeClass, align)}>
            {getContent(section)}
        </Tag>
    );
}

function ImageSection({section}: SectionProps) {
    const config = section.config ?? {};
    const src = getContent(section);
    if (!src) return null;
    const img = (
        <Image src={src} alt={config.alt || ''} className="w-full rounded-xl object-cover"/>
    );
    return (
        <figure className="overflow-hidden">
            {config.href ? (
                <a href={config.href} target={config.target || undefined} rel="noopener noreferrer">
                    {img}
                </a>
            ) : (
                img
            )}
            {config.caption && (
                <figcaption className="mt-2 text-center text-sm text-muted-foreground">
                    {config.caption}
                </figcaption>
            )}
        </figure>
    );
}

function GallerySection({section}: SectionProps) {
    const config = section.config ?? {};
    const columns = Math.min(Math.max(Number(config.columns) || 3, 1), 4);
    const images = getImageList(section);
    if (images.length === 0) return null;
    return (
        <div
            className="grid gap-4"
            style={{gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))`}}
        >
            {images.map((img, i) => (
                <figure key={i} className="overflow-hidden rounded-lg">
                    <Image src={img.src} alt={img.caption || ''} className="aspect-[4/3] w-full object-cover"/>
                    {img.caption && (
                        <figcaption className="mt-1 text-center text-xs text-muted-foreground">
                            {img.caption}
                        </figcaption>
                    )}
                </figure>
            ))}
        </div>
    );
}

function CarouselSection({section}: SectionProps) {
    const images = getImageList(section);
    if (images.length === 0) return null;
    // 轻量轮播：横向滚动吸附，避免为静态导出再引运行时轮播状态
    return (
        <div className="-mx-2 flex snap-x snap-mandatory gap-4 overflow-x-auto px-2 pb-2">
            {images.map((img, i) => (
                <figure
                    key={i}
                    className="w-[80%] shrink-0 snap-center overflow-hidden rounded-xl max-md:w-[85%]"
                >
                    <Image src={img.src} alt={img.caption || ''} className="aspect-video w-full object-cover"/>
                    {img.caption && (
                        <figcaption className="mt-1 text-center text-xs text-muted-foreground">
                            {img.caption}
                        </figcaption>
                    )}
                </figure>
            ))}
        </div>
    );
}

function VideoSection({section}: SectionProps) {
    const config = section.config ?? {};
    const src = getContent(section);
    if (!src) return null;
    const embed = toVideoEmbed(src);
    return (
        <div className="overflow-hidden rounded-xl border border-border">
            {embed ? (
                <iframe
                    src={embed}
                    title="video"
                    className="aspect-video w-full"
                    allowFullScreen
                    loading="lazy"
                />
            ) : (
                <video src={src} poster={config.poster || undefined} controls className="aspect-video w-full"/>
            )}
        </div>
    );
}

function ButtonSection({section}: SectionProps) {
    const config = section.config ?? {};
    const label = getContent(section);
    const href = config.href || config.url || config.link;
    if (!label || !href) return null;
    const external = /^https?:\/\//.test(href) || config.target === '_blank';
    return (
        <div className={cn(TITLE_ALIGN_CLASS[config.align] || '')}>
            {external ? (
                <a href={href} target="_blank" rel="noopener noreferrer">
                    <Button variant={config.variant === 'outline' ? 'outline' : 'default'}>{label}</Button>
                </a>
            ) : (
                <a href={href}>
                    <Button variant={config.variant === 'outline' ? 'outline' : 'default'}>{label}</Button>
                </a>
            )}
        </div>
    );
}

function DividerSection() {
    return <hr className="border-border"/>;
}

function SpacerSection({section}: SectionProps) {
    const height = Math.min(Math.max(Number(section.config?.height) || 48, 8), 400);
    return <div style={{height}} aria-hidden="true"/>;
}

function CodeSection({section}: SectionProps) {
    const language = section.config?.language || '';
    const code = getContent(section);
    if (!code) return null;
    // 借道 ContentViewer 的 Markdown 管线获得 Shiki 高亮与样式
    return <ContentViewer content={'```' + language + '\n' + code + '\n```'} type="markdown"/>;
}

function HtmlSection({section}: SectionProps) {
    return <ContentViewer content={getContent(section)} type="html"/>;
}

function FormSection({section}: SectionProps) {
    // 表单提交后端（form_submission）落地前的占位：可见但不误导访客输入
    const label = pickSectionTranslation(section)?.content?.title || section.name || '';
    if (!label) return null;
    return (
        <div className="rounded-xl border border-dashed border-border bg-muted/40 p-6 text-center text-sm text-muted-foreground">
            {label}
        </div>
    );
}

function CustomSection({section}: SectionProps) {
    const value = getContent(section);
    if (!value) return null;
    return <ContentViewer content={value} type="markdown"/>;
}

const SECTION_RENDERERS: Record<string, React.ComponentType<SectionProps>> = {
    SECTION_TYPE_RICH_TEXT: RichTextSection,
    SECTION_TYPE_MARKDOWN: MarkdownSection,
    SECTION_TYPE_TITLE: TitleSection,
    SECTION_TYPE_IMAGE: ImageSection,
    SECTION_TYPE_GALLERY: GallerySection,
    SECTION_TYPE_CAROUSEL: CarouselSection,
    SECTION_TYPE_VIDEO: VideoSection,
    SECTION_TYPE_BUTTON: ButtonSection,
    SECTION_TYPE_DIVIDER: DividerSection,
    SECTION_TYPE_SPACER: SpacerSection,
    SECTION_TYPE_CODE: CodeSection,
    SECTION_TYPE_HTML: HtmlSection,
    SECTION_TYPE_FORM: FormSection,
    SECTION_TYPE_CUSTOM: CustomSection,
};

export interface SectionRendererProps {
    sections?: contentservicev1_Section[];
    className?: string;
}

/**
 * CMS Page 区块渲染器：按 sortOrder 排序渲染页面全部区块。
 * 未知类型按 Markdown 兜底渲染 value，保证后台新增类型时页面不空白。
 */
export default function SectionRenderer({sections, className}: SectionRendererProps) {
    const sorted = useMemo(() => {
        return [...(sections ?? [])].sort(
            (a, b) => (a.sortOrder ?? 0) - (b.sortOrder ?? 0),
        );
    }, [sections]);

    if (sorted.length === 0) return null;

    return (
        <div className={cn('space-y-8', className)}>
            {sorted.map((section) => {
                const Renderer =
                    (section.type && SECTION_RENDERERS[section.type]) || CustomSection;
                const customClass = section.config?.className;
                return (
                    <div key={section.id ?? section.name ?? section.sortOrder} className={customClass}>
                        <Renderer section={section}/>
                    </div>
                );
            })}
        </div>
    );
}
