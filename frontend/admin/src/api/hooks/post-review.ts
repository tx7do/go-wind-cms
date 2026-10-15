import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_Post,
  type contentservicev1_RejectPostReviewRequest,
  type contentservicev1_ReviewPostRequest,
  type contentservicev1_SubmitPostReviewRequest,
} from '@/api/generated/admin/service/v1';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';

// ==============================
// 发布审批工作流：提交审核 / 通过 / 驳回
// 权限卡点 = 接口管理把 review 端点绑定到审核角色（部署期配置）。
// ==============================

export function useSubmitPostReview(
  options?: UseMutationOptions<contentservicev1_Post, Error, contentservicev1_SubmitPostReviewRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.postService.SubmitReview(req),
    ...options,
  });
}

export function useApprovePostReview(
  options?: UseMutationOptions<contentservicev1_Post, Error, contentservicev1_ReviewPostRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.postService.ApproveReview(req),
    ...options,
  });
}

export function useRejectPostReview(
  options?: UseMutationOptions<contentservicev1_Post, Error, contentservicev1_RejectPostReviewRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.postService.RejectReview(req),
    ...options,
  });
}

/** 供操作列刷新：失效文章列表与详情缓存 */
export function invalidatePostCaches(postId?: number) {
  queryClient.invalidateQueries({ queryKey: ['listPosts'] });
  if (postId) queryClient.invalidateQueries({ queryKey: ['getPost', { id: postId }] });
}
