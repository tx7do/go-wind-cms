'use client';

import {useQuery} from '@tanstack/react-query';
import {useTranslations} from 'next-intl';
import {apiClient} from '@/api/client';

/**
 * 自定义字段渲染（内容建模前台读侧）：
 * 按文章分类链绑定的内容模型字段定义，渲染 post.customFields。
 * 值约定为字符串（admin 编辑器 Record<string,string> 口径）：
 *   - RICHTEXT 按富文本注入；IMAGE/FILE 值形似 URL 时按图片/链接渲染；
 *   - 其余（TEXT/NUMBER/RELATION 等）原样展示。
 * 模型未绑定或加载失败静默不渲染（自定义字段是增强信息，不阻断正文）。
 */
export default function CustomFieldsSection({
    postId,
    customFields,
}: {
    postId?: number | null;
    customFields?: {[key: string]: string};
}) {
    const t = useTranslations('page');

    const modelQuery = useQuery({
        queryKey: ['postContentModel', postId],
        queryFn: () => apiClient.postService.GetPostContentModel({postId: postId!}),
        enabled: !!postId,
        retry: 0,
        staleTime: 5 * 60 * 1000,
    });

    const fields = modelQuery.data?.fields ?? [];
    if (!postId || fields.length === 0 || !customFields) return null;

    // 仅渲染模型定义过且文章实际携带的字段，按模型字段顺序
    const rows = fields
        .filter((f) => customFields[f.name ?? ''] !== undefined && customFields[f.name ?? ''] !== '')
        .map((f) => ({
            label: f.label || f.name || '',
            type: f.type ?? 'FIELD_TYPE_TEXT',
            name: f.name ?? '',
            value: customFields[f.name ?? ''],
        }));
    if (rows.length === 0) return null;

    const isUrl = (v: string) => /^https?:\/\//.test(v);

    return (
        <section className="mb-8 rounded-xl border border-border bg-card p-6">
            <h2 className="mb-4 text-lg font-semibold text-foreground">
                {modelQuery.data?.name || t('post_detail.custom_fields')}
            </h2>
            <dl className="grid grid-cols-1 gap-x-8 gap-y-3 sm:grid-cols-[180px_1fr]">
                {rows.map((row) => (
                    <div key={row.name} className="contents">
                        <dt className="text-sm font-medium text-muted-foreground">{row.label}</dt>
                        <dd className="min-w-0 text-sm text-foreground">
                            {row.type === 'FIELD_TYPE_RICHTEXT' ? (
                                <div className="prose-sm max-w-none [&_a]:text-primary [&_img]:max-w-full"
                                    dangerouslySetInnerHTML={{__html: row.value}}/>
                            ) : row.type === 'FIELD_TYPE_IMAGE' && isUrl(row.value) ? (
                                <img src={row.value} alt={row.label} className="max-h-64 rounded-lg object-cover"/>
                            ) : row.type === 'FIELD_TYPE_FILE' && isUrl(row.value) ? (
                                <a href={row.value} target="_blank" rel="noopener noreferrer"
                                    className="break-all text-primary hover:underline">
                                    {row.value}
                                </a>
                            ) : (
                                <span className="break-all">{row.value}</span>
                            )}
                        </dd>
                    </div>
                ))}
            </dl>
        </section>
    );
}
