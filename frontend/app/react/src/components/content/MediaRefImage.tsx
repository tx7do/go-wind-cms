'use client';

import {useState} from 'react';
import {apiClient} from '@/api/client';
import {env} from '@/config/env';

/**
 * 媒体引用图片（即时图片变换接入）：
 *   - IMAGE 值为 `media_asset:<id>` 引用时，优先经免鉴权变换端点渲染
 *     （只缩小不放大，w=1280 覆盖详情列宽 2x 余量、q 默认 80）；
 *   - gif/svg 等不可变换源（415）或资产缺失（404）回退元数据 url，再败回占位图；
 *   - http(s) 直链原样渲染（无资产 id，无法参与变换）。
 */

const MEDIA_REF_RE = /^media_asset:(\d+)$/;

const TRANSFORM_QUERY = 'w=1280&q=80';

const PLACEHOLDER_SRC = '/placeholder.svg';

/** 解析 media_asset:<id> 引用，非引用返回 null */
export function mediaRefId(value: string): string | null {
    const m = MEDIA_REF_RE.exec(value.trim());
    return m ? m[1]! : null;
}

export function mediaTransformUrl(id: string): string {
    return `${env.apiBaseUrl}/app/v1/media-assets/${id}/content?${TRANSFORM_QUERY}`;
}

export function isMediaImageValue(value: string): boolean {
    return mediaRefId(value) !== null || /^https?:\/\//.test(value);
}

export default function MediaRefImage({alt, className, value}: {
    value: string;
    alt: string;
    className?: string;
}) {
    const id = mediaRefId(value);
    // 渲染链：变换 URL → 元数据 url（回退）→ 占位图（再败终态）
    const [fallbackUrl, setFallbackUrl] = useState<string | null>(null);
    const [failed, setFailed] = useState(false);

    if (!id && !isMediaImageValue(value)) return null;

    const handleError = () => {
        if (id && fallbackUrl === null && !failed) {
            // 变换端点不可用：回退元数据 url。取元数据失败仅告警——
            // 这是 best-effort 兜底链的中间环，终态仍会落占位图。
            apiClient.mediaAssetService
                .Get({id: Number(id)})
                .then((asset) => {
                    if (asset.url) setFallbackUrl(asset.url);
                    else setFailed(true);
                })
                .catch((err: unknown) => {
                    console.error('[MediaRefImage] 元数据回退失败:', err);
                    setFailed(true);
                });
            return;
        }
        setFailed(true);
    };

    const src = failed ? PLACEHOLDER_SRC : id && !fallbackUrl ? mediaTransformUrl(id) : fallbackUrl ?? value;

    return (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={src} alt={alt} loading="lazy" className={className} onError={handleError}/>
    );
}
