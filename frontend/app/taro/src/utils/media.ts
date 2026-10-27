/**
 * 媒体引用工具（即时图片变换接入）：
 * customFields 的 IMAGE 值约定为字符串——`media_asset:<id>` 引用或 http 直链。
 * 变换端点只缩小不放大；gif/svg 等不可变换源返回 415、资产缺失 404，由调用方兜底。
 */

import {env} from '@/config/env';

/** media_asset:<id> 引用值形态（与 admin 写侧约定对齐） */
export const MEDIA_REF_RE = /^media_asset:(\d+)$/;

/** 变换端点默认参数：w=1280 覆盖详情列宽 2x 余量、q=80 */
export const MEDIA_TRANSFORM_QUERY = 'w=1280&q=80';

/** 解析 media_asset:<id> 引用，非引用返回 null */
export function mediaRefId(value: string): string | null {
    return MEDIA_REF_RE.exec(value.trim())?.[1] ?? null;
}

/** 变换端点 URL（免鉴权；weapp <image> 与 API 同域，不新增白名单域名） */
export function mediaTransformUrl(id: string): string {
    return `${env.apiBaseUrl}/app/v1/media-assets/${id}/content?${MEDIA_TRANSFORM_QUERY}`;
}

/** IMAGE 值是否可按图片渲染（http 直链或 media_asset 引用） */
export function isMediaImageValue(value: string): boolean {
    return mediaRefId(value) !== null || /^https?:\/\//.test(value);
}
