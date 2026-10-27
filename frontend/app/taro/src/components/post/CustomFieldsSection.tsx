import {View, Text, Image, RichText} from '@tarojs/components';
import {useState} from 'react';
import {useQuery} from '@tanstack/react-query';
import {useTranslation} from 'react-i18next';
import {apiClient} from '@/api/client';
import {useI18nRouter} from '@/i18n/helpers';
import {mediaRefId, mediaTransformUrl} from '@/utils/media';

/**
 * 自定义字段渲染（内容建模前台读侧，小程序口径）：
 *   - RICHTEXT → RichText（weapp 原生富文本组件，接受 HTML 字符串）；
 *   - IMAGE → media_asset:<id> 引用优先经即时图片变换端点渲染（只缩小不放大；
 *     gif/svg 415 或资产缺失 404 回退元数据 url，再败回原样文本），http 直链直用；
 *   - FILE → media_asset:<id> 引用经公开媒体端点换文件名展示；
 *   - RELATION → "post:12"/"page:7" 引用批量解析为公开摘要：post 目标可点击
 *     跳详情页；taro 无页面详情路由，page 目标降级为标题纯文本；
 *   - COMPONENT/DYNAMIC_ZONE → 组件实例 JSON 解析渲染键值卡（模型响应携带
 *     组件 schema，动态区按 _component 分发），解析失败回退原样字符串；
 *   - 其余原样展示。模型未绑定或加载失败静默不渲染。
 */

/** RELATION 引用值形态（与 core contentRefPattern 对齐；category/tag v1 不解析） */
const RELATION_REF_RE = /^(post|page):\d+$/;

interface FieldDef {
  name?: string;
  label?: string;
  type?: string;
  componentConfig?: {
    componentCode?: string;
    repeatable?: boolean;
    allowedComponentCodes?: string[];
  };
}

interface InstanceData {
  _component?: string;
  [key: string]: unknown;
}

function parseInstances(value: string, expectArray: boolean): InstanceData[] | null {
  if (!value.trim()) return [];
  try {
    const parsed: unknown = JSON.parse(value);
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

/** 变换图片：变换 URL 失败回退元数据 url，再败回原样文本（媒体资产是增强信息不阻断） */
function TransformImage({assetId, fallbackText}: {assetId: number; fallbackText: string}) {
  const [fallbackUrl, setFallbackUrl] = useState<string | null>(null);
  const [failed, setFailed] = useState(false);

  if (failed) return <Text className='break-all text-sm text-textMain'>{fallbackText}</Text>;

  const handleError = () => {
    if (fallbackUrl === null) {
      // 变换端点不可用：回退元数据 url。取元数据失败仅告警——
      // 这是 best-effort 兜底链的中间环，终态落回原样文本。
      apiClient.mediaAssetService
        .Get({id: assetId})
        .then((asset) => {
          if (asset.url) setFallbackUrl(asset.url);
          else setFailed(true);
        })
        .catch((error: unknown) => {
          console.error('[CustomFieldsSection] media asset fallback failed:', error);
          setFailed(true);
        });
      return;
    }
    setFailed(true);
  };

  return (
    <Image
      src={fallbackUrl ?? mediaTransformUrl(String(assetId))}
      mode='widthFix'
      className='w-full rounded-[16rpx]'
      lazyLoad
      onError={handleError}
    />
  );
}

function FieldValue({type, value, relations}: {type: string; value: string; relations?: Map<string, {title: string; targetId?: number}>}) {
  const router = useI18nRouter();
  const isUrl = /^https?:\/\//.test(value);
  const mediaId = mediaRefId(value);
  const needsAsset = (type === 'FIELD_TYPE_IMAGE' || type === 'FIELD_TYPE_FILE') && mediaId !== null;
  const assetQuery = useQuery({
    queryKey: ['mediaAssetUrl', mediaId],
    queryFn: () => apiClient.mediaAssetService.Get({id: Number(mediaId)}),
    enabled: needsAsset,
    retry: 0,
    staleTime: 5 * 60 * 1000,
  });

  if (type === 'FIELD_TYPE_RICHTEXT') {
    // weapp rich-text 接受 HTML 字符串；H5 端由 Taro 编译为等价节点
    return <RichText nodes={value} />;
  }
  if (type === 'FIELD_TYPE_IMAGE') {
    if (mediaId !== null) return <TransformImage assetId={Number(mediaId)} fallbackText={value} />;
    if (isUrl) return <Image src={value} mode='widthFix' className='w-full rounded-[16rpx]' lazyLoad />;
    return <Text className='break-all text-sm text-textMain'>{value}</Text>;
  }
  if (type === 'FIELD_TYPE_FILE') {
    const name = mediaId !== null ? assetQuery.data?.filename || `#${mediaId}` : value;
    return <Text className='break-all text-sm text-primary'>{name}</Text>;
  }
  if (type === 'FIELD_TYPE_RELATION') {
    const ref = value.trim();
    if (!RELATION_REF_RE.test(ref)) return <Text className='break-all'>{value}</Text>;
    const target = relations?.get(ref);
    if (!target) return <Text className='break-all'>{value}</Text>;
    if (target.targetId) {
      return (
        <Text className='break-all text-primary' onClick={() => router.push(`/post/detail?id=${target.targetId}`)}>
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
    const {t} = useTranslation();
    const modelQuery = useQuery({
        queryKey: ['postContentModel', postId],
        queryFn: () => apiClient.postService.GetPostContentModel({postId: postId!}),
        enabled: !!postId,
        retry: 0,
        staleTime: 5 * 60 * 1000,
    });

    const fields = modelQuery.data?.fields ?? [];

    // 仅渲染模型定义过且文章实际携带的字段，按模型字段顺序
    // （rows 在早退前计算：relationsQuery 须无条件调用，hooks 顺序不能依赖数据）
    const rows = (!postId || !customFields)
        ? []
        : fields
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

    if (!postId || fields.length === 0 || !customFields) return null;
    if (rows.length === 0) return null;
    const relations = new Map<string, {title: string; targetId?: number}>(
        (relationsQuery.data?.items ?? []).map((it) => [
            it.ref ?? '',
            {title: it.title ?? '', targetId: it.targetEntityType === 'post' ? it.targetId : undefined},
        ]),
    );

    // 组件 schema（模型响应携带）
    const componentMap = new Map<string, {name: string; fields: FieldDef[]}>(
        (modelQuery.data?.components ?? []).map((c) => [
            c.definition?.code ?? '',
            {name: c.definition?.name ?? '', fields: (c.fields ?? []) as FieldDef[]},
        ]),
    );

    const renderComponentValue = (field: FieldDef, value: string) => {
        const cc = field.componentConfig;
        if (!cc) return <Text className='break-all'>{value}</Text>;
        const expectArray = !!cc.repeatable || field.type === 'FIELD_TYPE_DYNAMIC_ZONE';
        const instances = parseInstances(value, expectArray);
        if (instances === null) return <Text className='break-all'>{value}</Text>;
        if (instances.length === 0) return null;
        return (
            <View className='flex flex-col gap-[16rpx]'>
                {instances.map((inst, i) => {
                    const code = field.type === 'FIELD_TYPE_DYNAMIC_ZONE'
                        ? String(inst._component ?? '')
                        : cc.componentCode ?? '';
                    const schema = componentMap.get(code);
                    const pairs = (schema?.fields ?? [])
                        .map((f) => {
                            const raw = inst[f.name ?? ''];
                            const text = typeof raw === 'string' ? raw : '';
                            return {
                                label: f.label || f.name || '',
                                text,
                                image: f.type === 'FIELD_TYPE_IMAGE' && mediaRefId(text) !== null,
                                assetId: mediaRefId(text),
                            };
                        })
                        .filter((p) => p.text !== '');
                    return (
                        <View key={i} className='rounded-[16rpx] p-[20rpx] border-[length:1rpx] border-splitLine'>
                            <Text className='mb-[8rpx] block text-tips text-textSec'>{schema?.name ?? code}</Text>
                            {pairs.length === 0 ? (
                                <Text className='text-tips text-textSec'>—</Text>
                            ) : pairs.map((p, j) => (
                                <View key={j} className='py-[4rpx]'>
                                    <Text className='block text-tips text-textSec'>{p.label}</Text>
                                    {p.image && p.assetId !== null ? (
                                        <TransformImage assetId={Number(p.assetId)} fallbackText={p.text} />
                                    ) : (
                                        <Text className='break-all text-desc text-textMain'>{p.text}</Text>
                                    )}
                                </View>
                            ))}
                        </View>
                    );
                })}
            </View>
        );
    };

    return (
        <View className='au-card p-[32rpx] mb-[28rpx]'>
            <Text className='au-section-title block mb-[22rpx]'>
                {modelQuery.data?.name || t('page.post_detail.custom_fields')}
            </Text>
            <View className='flex flex-col gap-[24rpx]'>
                {rows.map((row) => {
                    const isComponent =
                        row.type === 'FIELD_TYPE_COMPONENT' || row.type === 'FIELD_TYPE_DYNAMIC_ZONE';
                    const fieldDef = fields.find((f) => f.name === row.name);
                    return (
                        <View key={row.name} className='flex gap-[24rpx]'>
                            <Text className='w-[168rpx] shrink-0 text-tips text-textSec'>{row.label}</Text>
                            <View className='min-w-0 flex-1'>
                                {isComponent && fieldDef ? (
                                    renderComponentValue(fieldDef as FieldDef, row.value)
                                ) : (
                                    <FieldValue type={row.type} value={row.value} relations={relations} />
                                )}
                            </View>
                        </View>
                    );
                })}
            </View>
        </View>
    );
}
