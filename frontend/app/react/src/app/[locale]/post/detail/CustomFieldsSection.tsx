'use client';

import {useQuery} from '@tanstack/react-query';
import {useLocale, useTranslations} from 'next-intl';
import {apiClient} from '@/api/client';
import MediaRefImage, {isMediaImageValue} from '@/components/content/MediaRefImage';
import type {
  contentservicev1_ContentModelComponent,
  contentservicev1_ContentReferenceSummary,
  contentservicev1_FieldDefinition,
} from '@/api/generated/app/service/v1';

/**
 * 自定义字段渲染（内容建模前台读侧）：
 * 按文章分类链绑定的内容模型字段定义，渲染 post.customFields。
 * 值约定为字符串（admin 编辑器 Record<string,string> 口径）：
 *   - RICHTEXT 按富文本注入；IMAGE/FILE 值形似 URL 时按图片/链接渲染；
 *   - RELATION 值为 "post:12"/"page:7" 引用时批量解析为公开摘要，
 *     渲染成指向目标详情页的链接（草稿/软删目标解析缺失，回退原样字符串）；
 *   - COMPONENT/DYNAMIC_ZONE 值为组件实例 JSON（模型响应携带组件 schema），
 *     解析后渲染键值卡（动态区按 _component 分发），解析失败回退原样字符串；
 *   - 其余（TEXT/NUMBER 等）原样展示。
 * 模型未绑定或加载失败静默不渲染（自定义字段是增强信息，不阻断正文）。
 */

/** RELATION 引用值形态（与 core contentRefPattern 对齐；category/tag v1 不解析） */
const RELATION_REF_RE = /^(post|page):\d+$/;

type FieldDef = Partial<contentservicev1_FieldDefinition>;
type ComponentSchemaMap = Record<string, contentservicev1_ContentModelComponent | undefined>;
type InstanceData = {_component?: string; [key: string]: unknown};

function parseInstances(value: string, expectArray: boolean): InstanceData[] | null {
  if (!value.trim()) return [];
  try {
    const parsed = JSON.parse(value) as unknown;
    if (expectArray) {
      if (!Array.isArray(parsed)) return null;
      return parsed.filter(
        (x): x is InstanceData => Boolean(x) && typeof x === 'object' && !Array.isArray(x),
      );
    }
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) return null;
    return [parsed as InstanceData];
  } catch {
    return null;
  }
}

function buildComponentMap(components?: contentservicev1_ContentModelComponent[]): ComponentSchemaMap {
  const map: ComponentSchemaMap = {};
  for (const c of components ?? []) {
    if (c.definition?.code) map[c.definition.code] = c;
  }
  return map;
}

/** 组件实例 → 标签/文本/图片对（IMAGE 支持 http 直链与 media_asset:<id> 引用） */
function instanceRows(schema: contentservicev1_ContentModelComponent | undefined, inst: InstanceData) {
  const fields = schema?.fields ?? [];
  const rows: {label: string; text: string; image: boolean}[] = [];
  for (const f of fields) {
    const key = f.name ?? '';
    const raw = inst[key];
    if (raw == null || typeof raw !== 'string' || raw === '') continue;
    const type = f.type ?? 'FIELD_TYPE_TEXT';
    const isImg = type === 'FIELD_TYPE_IMAGE' && isMediaImageValue(raw);
    rows.push({label: f.label || key, text: raw, image: isImg});
  }
  return rows;
}

export default function CustomFieldsSection({
    postId,
    customFields,
}: {
    postId?: number | null;
    customFields?: {[key: string]: string};
}) {
    const t = useTranslations('page');
    const locale = useLocale();

    const modelQuery = useQuery({
        queryKey: ['postContentModel', postId],
        queryFn: () => apiClient.postService.GetPostContentModel({postId: postId!}),
        enabled: !!postId,
        retry: 0,
        staleTime: 5 * 60 * 1000,
    });

    const fields = modelQuery.data?.fields ?? [];

    // 仅渲染模型定义过且文章实际携带的字段，按模型字段顺序
    // （rows 在早退前计算：resolveQuery 须无条件调用，hooks 顺序不能依赖数据）
    const rows = fields
        .filter((f) => customFields?.[f.name ?? ''] !== undefined && customFields?.[f.name ?? ''] !== '')
        .map((f) => ({
            label: f.label || f.name || '',
            type: f.type ?? 'FIELD_TYPE_TEXT',
            name: f.name ?? '',
            value: customFields?.[f.name ?? ''] ?? '',
        }));

    const isUrl = (v: string) => /^https?:\/\//.test(v);

    const relationRefs = Array.from(new Set(rows
        .filter((r) => r.type === 'FIELD_TYPE_RELATION' && RELATION_REF_RE.test(r.value.trim()))
        .map((r) => r.value.trim())));

    const resolveQuery = useQuery({
        queryKey: ['contentReferences', relationRefs.join(','), locale],
        queryFn: () => apiClient.postService.ResolveContentReferences({refs: relationRefs, languageCode: locale}),
        enabled: relationRefs.length > 0,
        retry: 0,
        staleTime: 5 * 60 * 1000,
    });
    const resolved = new Map<string, contentservicev1_ContentReferenceSummary>(
        (resolveQuery.data?.items ?? []).map((it) => [it.ref ?? '', it]),
    );

    const renderRelation = (value: string) => {
        const ref = value.trim();
        if (!RELATION_REF_RE.test(ref)) return <span className="break-all">{value}</span>;
        const target = resolved.get(ref);
        if (!target?.slug) return <span className="break-all">{value}</span>;
        const href = target.targetEntityType === 'page'
            ? `/${locale}/page/detail?slug=${encodeURIComponent(target.slug)}&lang=${locale}`
            : `/${locale}/post/detail?slug=${encodeURIComponent(target.slug)}&lang=${locale}`;
        return (
            <a href={href} className="break-all text-primary hover:underline">
                {target.title || ref}
            </a>
        );
    };

    const componentMap = buildComponentMap(modelQuery.data?.components);

    const renderComponentValue = (field: FieldDef, value: string) => {
        const cc = field.componentConfig;
        if (!cc) return <span className="break-all">{value}</span>;
        const expectArray = !!cc.repeatable || field.type === 'FIELD_TYPE_DYNAMIC_ZONE';
        const instances = parseInstances(value, expectArray);
        // 解析失败（非法 JSON）：回退原样字符串
        if (instances === null) return <span className="break-all">{value}</span>;
        if (instances.length === 0) return null;

        return (
            <div className="flex flex-col gap-2">
                {instances.map((inst, i) => {
                    const code = field.type === 'FIELD_TYPE_DYNAMIC_ZONE'
                        ? String(inst._component ?? '')
                        : cc.componentCode ?? '';
                    const schema = componentMap[code];
                    const compName = schema?.definition?.name ?? code;
                    const pairs = instanceRows(schema, inst);
                    return (
                        <div key={i} className="rounded-lg border border-border p-3">
                            <div className="mb-1 text-xs font-medium text-muted-foreground">
                                {compName}
                            </div>
                            {pairs.length === 0 ? (
                                <span className="text-xs text-muted-foreground">—</span>
                            ) : (
                                pairs.map((p, j) => (
                                    <div key={j} className="flex flex-col py-0.5">
                                        <span className="text-xs text-muted-foreground">{p.label}</span>
                                        {p.image ? (
                                            <MediaRefImage value={p.text} alt={p.label}
                                                className="max-h-48 rounded-md object-cover"/>
                                        ) : (
                                            <span className="break-all text-sm">{p.text}</span>
                                        )}
                                    </div>
                                ))
                            )}
                        </div>
                    );
                })}
            </div>
        );
    };

    if (!postId || fields.length === 0 || !customFields) return null;
    if (rows.length === 0) return null;

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
                            ) : row.type === 'FIELD_TYPE_IMAGE' && isMediaImageValue(row.value) ? (
                                <MediaRefImage value={row.value} alt={row.label}
                                    className="max-h-64 rounded-lg object-cover"/>
                            ) : row.type === 'FIELD_TYPE_FILE' && isUrl(row.value) ? (
                                <a href={row.value} target="_blank" rel="noopener noreferrer"
                                    className="break-all text-primary hover:underline">
                                    {row.value}
                                </a>
                            ) : row.type === 'FIELD_TYPE_RELATION' ? (
                                renderRelation(row.value)
                            ) : row.type === 'FIELD_TYPE_COMPONENT' || row.type === 'FIELD_TYPE_DYNAMIC_ZONE' ? (
                                renderComponentValue(
                                    fields.find((f) => f.name === row.name) ?? {},
                                    row.value,
                                )
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
