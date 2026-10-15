import {
  useMutation,
  type UseMutationOptions,
} from '@tanstack/react-query';
import {
  type contentservicev1_DeleteFormSubmissionRequest,
  type contentservicev1_FormSubmission,
} from '@/api/generated/admin/service/v1';
import { queryClient } from '@/core';
import { apiClient } from '@/api/client';
import { type PaginationQuery } from '@/core/transport/rest';

// ==============================
// 表单提交管理（页面 FORM 区块的访客提交）
// ==============================

export function useFormSubmissionsList() {
  return useMutation({
    mutationFn: (query: PaginationQuery) =>
      apiClient.formSubmissionService.List(query.toRawParams()),
  });
}

export async function fetchFormSubmission(id: number) {
  return queryClient.fetchQuery({
    queryKey: ['getFormSubmission', id],
    queryFn: () => apiClient.formSubmissionService.Get({ id }),
    retry: 0,
  });
}

export function useDeleteFormSubmission(
  options?: UseMutationOptions<Record<string, never>, Error, contentservicev1_DeleteFormSubmissionRequest>,
) {
  return useMutation({
    mutationFn: (req) => apiClient.formSubmissionService.Delete(req),
    ...options,
  });
}

/** CSV 导出：当前列表数据在前端拼装（量级小，无需服务端流式导出） */
export function submissionsToCsv(
  items: contentservicev1_FormSubmission[],
): string {
  const fieldSet = new Set<string>(['id', 'pageId', 'sectionId', 'createdAt']);
  for (const item of items) {
    for (const key of Object.keys(item.data ?? {})) {
      fieldSet.add(key);
    }
  }
  const columns = [...fieldSet];
  const escape = (v: unknown) => {
    const s = v === undefined || v === null ? '' : String(v);
    return /[",\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s;
  };
  const lines = [columns.join(',')];
  for (const item of items) {
    lines.push(
      columns
        .map((c) => {
          if (c === 'id' || c === 'pageId' || c === 'sectionId') return escape(item[c]);
          if (c === 'createdAt') return escape(item.createdAt ?? '');
          return escape(item.data?.[c] ?? '');
        })
        .join(','),
    );
  }
  // BOM 便于 Excel 识别 UTF-8
  return '\ufeff' + lines.join('\n');
}
