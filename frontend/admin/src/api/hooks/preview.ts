import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type CreatePagePreviewTokenRequest,
  type CreatePostPreviewTokenRequest,
  type PreviewToken,
} from '@/api/generated/admin/service/v1';
import { apiClient } from '@/api/client';
import { supportedLocales } from '@/locales';

// ==============================
// 草稿预览：admin 签发短时效令牌，C 端公开 BFF 校验后放行草稿
// ==============================

export function useCreatePostPreviewToken(
  options?: UseMutationOptions<PreviewToken, Error, CreatePostPreviewTokenRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.previewService.CreatePostPreviewToken(req),
    ...options,
  });
}

export function useCreatePagePreviewToken(
  options?: UseMutationOptions<PreviewToken, Error, CreatePagePreviewTokenRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.previewService.CreatePagePreviewToken(req),
    ...options,
  });
}

/** C 端站点基址：预览在新标签打开 C 端详情路由（可按部署环境覆盖） */
export function cendBaseUrl(): string {
  return (import.meta.env.VITE_CEND_BASE_URL ?? 'http://localhost:15001').replace(/\/$/, '');
}

/** 预览链接的语言前缀：C 端 localePrefix=always，裸路径无路由会 404。
 *  admin 与 C 端语言码一致（zh-CN/en-US），取当前界面语言，非支持语言回退默认语言 */
export function cendLocalePrefix(language: string | undefined): string {
  const hit = supportedLocales.find((l) => l === language);
  return hit ?? supportedLocales[0];
}
