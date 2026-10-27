import React, {useMemo, useState} from 'react';
import {View, Text, Image, Video, Button, Input, Textarea, ScrollView} from '@tarojs/components';
import {cn} from '@/lib/utils';
import ContentViewer from '@/components/content/ContentViewer';
import {apiClient} from '@/api/client';
import {currentLocaleLanguageCode} from '@/i18n';

import type {
    contentservicev1_Section,
    contentservicev1_SectionTranslation,
} from '@/api/generated/app/service/v1';

/** FORM 区块字段定义（config.fields JSON 数组元素；与 core 侧校验契约一致） */
interface FormFieldDef {
    name: string;
    type?: string;
    required?: boolean;
}

/**
 * Section 内容键契约（三端渲染器共用，改键先同步 react/vue 端 SectionRenderer 与本文档）
 *
 * 区块语言相关内容存 translations[].content（map<string,string>），缺译回退：
 * 匹配当前语言 → 第一条翻译。语言无关样式存 config（map<string,string>）。
 * 主内容键一律为 `value`，其余可选键见 react 端 SectionRenderer 顶部契约表。
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

/** GALLERY/CAROUSEL 的图片列表：value 支持 JSON 数组或换行分隔 URL 列表 */
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

const isH5 = process.env.TARO_ENV === 'h5';

interface SectionProps {
    section: contentservicev1_Section;
}

function RichTextSection({section}: SectionProps) {
    return <ContentViewer content={getContent(section)} type='html' />;
}

function MarkdownSection({section}: SectionProps) {
    return <ContentViewer content={getContent(section)} type='markdown' />;
}

const TITLE_SIZE: Record<number, string> = {
    1: 'text-[length:44rpx]',
    2: 'text-[length:40rpx]',
    3: 'text-[length:36rpx]',
    4: 'text-[length:32rpx]',
    5: 'text-[length:30rpx]',
    6: 'text-[length:28rpx]',
};

function TitleSection({section}: SectionProps) {
    const config = section.config ?? {};
    const level = Math.min(Math.max(Number(config.level) || 2, 1), 6);
    const align = config.align === 'center' ? 'justify-center' : config.align === 'right' ? 'justify-end' : 'justify-start';
    return (
        <View className={cn('flex', align)}>
            <Text className={cn('font-bold text-textMain', TITLE_SIZE[level])}>
                {getContent(section)}
            </Text>
        </View>
    );
}

function ImageSection({section}: SectionProps) {
    const config = section.config ?? {};
    const src = getContent(section);
    if (!src) return null;
    return (
        <View className='overflow-hidden rounded-[20rpx]'>
            <Image src={src} mode='widthFix' className='w-full' lazyLoad />
            {config.caption && (
                <Text className='block text-center text-tips text-textThird mt-[10rpx]'>
                    {config.caption}
                </Text>
            )}
        </View>
    );
}

function GallerySection({section}: SectionProps) {
    const config = section.config ?? {};
    const columns = Math.min(Math.max(Number(config.columns) || 2, 1), 4);
    const images = getImageList(section);
    if (images.length === 0) return null;
    return (
        <View
          className='grid gap-[16rpx]'
          style={{gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))`}}
        >
            {images.map((img, i) => (
                <View key={i} className='overflow-hidden rounded-[16rpx]'>
                    <Image
                      src={img.src}
                      mode='aspectFill'
                      className='w-full h-[220rpx]'
                      lazyLoad
                    />
                    {img.caption && (
                        <Text className='block text-center text-tips text-textThird mt-[8rpx]'>
                            {img.caption}
                        </Text>
                    )}
                </View>
            ))}
        </View>
    );
}

function CarouselSection({section}: SectionProps) {
    const images = getImageList(section);
    if (images.length === 0) return null;
    // 轻量轮播：横向滚动吸附（两端一致的 scroll-view 实现，不引第三方轮播）
    return (
        <View className='-mx-[24rpx]'>
            <ScrollView scrollX className='w-full' showScrollbar={false}>
                <View className='flex gap-[20rpx] px-[24rpx] pb-[8rpx]'>
                    {images.map((img, i) => (
                        <View key={i} className='w-[520rpx] shrink-0 overflow-hidden rounded-[20rpx]'>
                            <Image
                              src={img.src}
                              mode='aspectFill'
                              className='w-[520rpx] h-[300rpx]'
                              lazyLoad
                            />
                            {img.caption && (
                                <Text className='block text-center text-tips text-textThird mt-[8rpx]'>
                                    {img.caption}
                                </Text>
                            )}
                        </View>
                    ))}
                </View>
            </ScrollView>
        </View>
    );
}

function VideoSection({section}: SectionProps) {
    const config = section.config ?? {};
    const src = getContent(section);
    if (!src) return null;
    // 小程序不支持 iframe 嵌入；H5 上的 YouTube/Bilibili 嵌入链也与小程序不通用——
    // 统一退化为 <Video> 直链播放（mp4/webm），站点页跳系统浏览器由 LinkTap 处理
    return (
        <View className='overflow-hidden rounded-[20rpx]'>
            <Video
              src={src}
              poster={config.poster || undefined}
              controls
              className='w-full h-[400rpx]'
            />
        </View>
    );
}

function ButtonSection({section}: SectionProps) {
    const config = section.config ?? {};
    const label = getContent(section);
    const href = config.href || config.url || config.link;
    if (!label || !href) return null;
    const align = config.align === 'center' ? 'justify-center' : config.align === 'right' ? 'justify-end' : 'justify-start';
    const openLink = () => {
        // 外链跳转只在 H5 生效（小程序需业务域名白名单，且无 window）
        if (isH5 && typeof window !== 'undefined' && /^https?:\/\//.test(href)) {
            window.location.href = href;
        }
    };
    return (
        <View className={cn('flex', align)}>
            <Button
              className={cn(
                  'min-h-[80rpx] px-[48rpx] rounded-full text-body font-bold',
                  config.variant === 'outline'
                      ? 'border-[length:2rpx] border-splitLine bg-cardBg text-textMain'
                      : 'btn-hero',
              )}
              hoverClass='tap-active'
              onClick={openLink}
            >
                {label}
            </Button>
        </View>
    );
}

function DividerSection() {
    return <View className='h-[1rpx] bg-splitLine my-[8rpx]' />;
}

function SpacerSection({section}: SectionProps) {
    const height = Math.min(Math.max(Number(section.config?.height) || 48, 8), 400);
    return <View style={{height: `${height}rpx`}} />;
}

function CodeSection({section}: SectionProps) {
    const language = section.config?.language || '';
    const code = getContent(section);
    if (!code) return null;
    // 借道 ContentViewer 的 Markdown 管线获得代码块样式
    return <ContentViewer content={'```' + language + '\n' + code + '\n```'} type='markdown' />;
}

function HtmlSection({section}: SectionProps) {
    return <ContentViewer content={getContent(section)} type='html' />;
}

function FormSection({section}: SectionProps) {
    // 契约：config.fields = JSON 数组（语言无关结构）；翻译 content 的
    // `<字段名>` 键为标签（随语言变化），content.title 为表单标题。
    const [values, setValues] = useState<Record<string, string>>({});
    const [state, setState] = useState<'idle' | 'submitting' | 'done' | 'error'>('idle');

    const tr = pickSectionTranslation(section);
    const title = tr?.content?.title || section.name || '';

    const fields = useMemo<FormFieldDef[]>(() => {
        try {
            const parsed = JSON.parse(section.config?.fields || '[]');
            if (Array.isArray(parsed)) {
                return parsed.filter(
                    (f: unknown): f is FormFieldDef => !!f && typeof (f as FormFieldDef).name === 'string',
                );
            }
        } catch {
            // 字段定义非法：渲染为空表单（仅标题），不中断整页
        }
        return [];
    }, [section.config?.fields]);

    if (fields.length === 0 && !title) return null;

    const handleSubmit = async () => {
        if (!section.pageId || state === 'submitting') return;
        setState('submitting');
        try {
            await apiClient.formSubmissionService.Submit({
                pageId: section.pageId,
                sectionId: section.id ?? 0,
                data: values,
            });
            setState('done');
            setValues({});
        } catch (error) {
            console.error('Form submit failed:', error);
            setState('error');
        }
    };

    if (state === 'done') {
        return (
            <View className='au-card p-[32rpx] text-center'>
                <Text className='text-desc text-textSec'>
                    {tr?.content?.successText || title}
                </Text>
            </View>
        );
    }

    return (
        <View className='au-card p-[32rpx] flex flex-col gap-[24rpx]'>
            {title && (
                <Text className='text-body font-bold text-textMain'>{title}</Text>
            )}
            {fields.map((field) => {
                const label = tr?.content?.[field.name] || field.name;
                return (
                    <View key={field.name} className='flex flex-col gap-[10rpx]'>
                        <Text className='text-tips text-textSec'>
                            {label}
                            {field.required && <Text className='text-danger'> *</Text>}
                        </Text>
                        {field.type === 'textarea' ? (
                            <View className='au-textarea-box rounded-[16rpx] bg-pageBg border-[length:1rpx] border-splitLine px-[20rpx] py-[16rpx] box-border'>
                                <Textarea
                                  className='au-textarea'
                                  autoHeight
                                  value={values[field.name] || ''}
                                  onInput={(e) => setValues((v) => ({...v, [field.name]: e.detail.value}))}
                                />
                            </View>
                        ) : (
                            <View className='au-input'>
                                <Input
                                  className='au-input-field'
                                  value={values[field.name] || ''}
                                  onInput={(e) => setValues((v) => ({...v, [field.name]: e.detail.value}))}
                                />
                            </View>
                        )}
                    </View>
                );
            })}
            {state === 'error' && (
                <Text className='text-tips text-danger'>
                    {tr?.content?.errorText || 'Submit failed, please retry'}
                </Text>
            )}
            <Button
              className={cn('min-h-[84rpx] rounded-full text-body font-bold', 'btn-hero')}
              hoverClass='tap-active'
              disabled={state === 'submitting'}
              onClick={handleSubmit}
            >
                {state === 'submitting' ? '…' : tr?.content?.submitText || 'Submit'}
            </Button>
        </View>
    );
}

function CustomSection({section}: SectionProps) {
    const value = getContent(section);
    if (!value) return null;
    return <ContentViewer content={value} type='markdown' />;
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
        <View className={cn('flex flex-col gap-[36rpx]', className)}>
            {sorted.map((section) => {
                const Renderer =
                    (section.type && SECTION_RENDERERS[section.type]) || CustomSection;
                const customClass = section.config?.className;
                return (
                    <View key={section.id ?? section.name ?? section.sortOrder} className={customClass}>
                        <Renderer section={section} />
                    </View>
                );
            })}
        </View>
    );
}
