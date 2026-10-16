import {View, Text, Image, RichText} from '@tarojs/components';
import {useQuery} from '@tanstack/react-query';
import {apiClient} from '@/api/client';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

/**
 * 自定义字段渲染（内容建模前台读侧，小程序口径）：
 * 按文章分类链绑定的内容模型字段定义，渲染 post.customFields。
 * 值约定为字符串（admin 编辑器 Record<string,string> 口径）：
 *   - RICHTEXT → RichText（weapp 原生富文本组件，接受 HTML 字符串）；
 *   - IMAGE/FILE → media_asset:<id> 引用经公开媒体端点换 URL（H5 直接 URL 直用）；
 *   - RELATION → "post:12"/"page:7" 引用批量解析为公开摘要：post 目标可点击
 *     跳详情页；taro 无页面详情路由，page 目标降级为标题纯文本；
 *     草稿/软删目标解析缺失，回退原样字符串；
 *   - 其余（TEXT/NUMBER 等）原样展示。
 * 模型未绑定或加载失败静默不渲染（自定义字段是增强信息，不阻断正文）。
 */

/** RELATION 引用值形态（与 core contentRefPattern 对齐；category/tag v1 不解析） */
const RELATION_REF_RE = /^(post|page):\d+$/;

/** 解析 media_asset:<id> 引用；非引用形态返回 null */
function parseMediaRef(value: string): number | null {
    const m = /^media_asset:(\d+)$/.exec(value.trim());
    return m ? Number(m[1]) : null;
}

/** 按 id 换媒体资产公开 URL（值形如 media_asset:123 时使用） */
function useMediaAssetUrl(assetId: number | null, enabled: boolean) {
    return useQuery({
        queryKey: ['mediaAssetUrl', assetId],
        queryFn: () => apiClient.mediaAssetService.Get({id: assetId!}),
        enabled: enabled && !!assetId,
        retry: 0,
        staleTime: 5 * 60 * 1000,
    });
}

function FieldValue({type, value, relations}: {type: string; value: string; relations?: Map<string, {title: string; targetId?: number}>}) {
    const router = useI18nRouter();
    const isUrl = /^https?:\/\//.test(value);
    const mediaId = parseMediaRef(value);
    const needsAsset = (type === 'FIELD_TYPE_IMAGE' || type === 'FIELD_TYPE_FILE') && mediaId !== null;
    const assetQuery = useMediaAssetUrl(mediaId, needsAsset);
    const assetUrl = assetQuery.data?.url || '';

    if (type === 'FIELD_TYPE_RICHTEXT') {
        // weapp rich-text 接受 HTML 字符串；H5 端由 Taro 编译为等价节点
        return <RichText nodes={value}/>;
    }
    if (type === 'FIELD_TYPE_IMAGE') {
        const src = mediaId !== null ? assetUrl : isUrl ? value : '';
        if (!src) return <Text className='text-tips text-textSec'>{value}</Text>;
        return <Image src={src} mode='widthFix' style={{width: '100%', borderRadius: '12px'}}/>;
    }
    if (type === 'FIELD_TYPE_FILE') {
        const name = mediaId !== null ? assetQuery.data?.filename || `#${mediaId}` : value;
        return <Text className='break-all text-primary'>{name}</Text>;
    }
    if (type === 'FIELD_TYPE_RELATION') {
        const ref = value.trim();
        if (!RELATION_REF_RE.test(ref)) return <Text className='break-all'>{value}</Text>;
        const target = relations?.get(ref);
        if (!target) return <Text className='break-all'>{value}</Text>;
        if (target.targetId) {
            return (
                <Text className='break-all text-primary' onClick={() => router.push(`/post/${target.targetId}`)}>
                    {target.title || ref}
                </Text>
            );
        }
        // page 目标：taro 无页面详情路由，标题纯文本降级
        return <Text className='break-all'>{target.title || ref}</Text>;
    }
    return <Text className='break-all'>{value}</Text>;
}

export default function CustomFieldsSection({
    postId,
    customFields,
}: {
    /** 详情页的 postId 为 number|null（路由参数解析语义），null 时不渲染 */
    postId?: number | null;
    customFields?: {[key: string]: string};
}) {
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
        .filter((f) => {
            const v = customFields[f.name ?? ''];
            return v !== undefined && v !== '';
        })
        .map((f) => ({
            label: f.label || f.name || '',
            type: f.type ?? 'FIELD_TYPE_TEXT',
            name: f.name ?? '',
            value: customFields[f.name ?? ''],
        }));
    if (rows.length === 0) return null;

    const relationRefs = Array.from(new Set(rows
        .filter((r) => r.type === 'FIELD_TYPE_RELATION' && RELATION_REF_RE.test(r.value.trim()))
        .map((r) => r.value.trim())));

    const relationsQuery = useQuery({
        queryKey: ['contentReferences', relationRefs.join(',')],
        queryFn: () => apiClient.postService.ResolveContentReferences({refs: relationRefs}),
        enabled: relationRefs.length > 0,
        retry: 0,
        staleTime: 5 * 60 * 1000,
    });
    const relations = new Map<string, {title: string; targetId?: number}>(
        (relationsQuery.data?.items ?? []).map((it) => [
            it.ref ?? '',
            {title: it.title ?? '', targetId: it.targetEntityType === 'post' ? it.targetId : undefined},
        ]),
    );

    return (
        <View className='mb-6 rounded-2xl bg-card px-4 py-4' style={{border: '1px solid rgba(148,163,184,0.2)'}}>
            <Text className='mb-3 block text-base font-semibold text-textMain'>
                {modelQuery.data?.name || '自定义字段'}
            </Text>
            <View className='flex flex-col gap-3'>
                {rows.map((row) => (
                    <View key={row.name} className='flex gap-3'>
                        <Text className='w-24 shrink-0 text-xs text-textSec'>{row.label}</Text>
                        <View className='min-w-0 flex-1'>
                            <FieldValue type={row.type} value={row.value} relations={relations}/>
                        </View>
                    </View>
                ))}
            </View>
        </View>
    );
}
