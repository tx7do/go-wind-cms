import { useRef, useState } from 'react';
import { App, Button, Dropdown } from 'antd';
import {
  ExportOutlined,
  ImportOutlined,
  FileExcelOutlined,
  FileTextOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { PaginationQuery } from '@/core';
import ImportModal from '@/components/common/ImportModal';
import type { ImportField } from '@/utils/import';
import { fetchListPosts, useCreatePost } from '@/api/hooks/post';

/**
 * 帖子内容搬运按钮组：
 *  - 导入 Excel：复用通用 ImportModal（xlsx 逐行 create，强制 DRAFT）
 *  - 导入 JSON：接受本组件导出的全保真格式（translations 完整回填），逐帖 create（强制 DRAFT）
 *  - 导出 JSON：全量迭代（10000 条封顶），含 translations 完整内容——CSV 导出不含正文的补充
 *
 * 安全口径：导入一律落 DRAFT，不因批量导入意外发布内容。
 */

/** Excel 导入字段（label 即模板表头） */
const XLSX_FIELDS: ImportField[] = [
  { label: '标题', prop: 'title' },
  { label: '摘要', prop: 'summary' },
  { label: '正文(Markdown)', prop: 'content' },
  { label: '分类ID(逗号分隔)', prop: 'categoryIds' },
  { label: '标签ID(逗号分隔)', prop: 'tagIds' },
];

const EXPORT_FIELD_MASK =
  'id,status,sort_order,is_featured,author_name,created_at,code,editor_type,disallow_comment,thumbnail,custom_fields,publish_time,category_ids,tag_ids,translations.id,translations.post_id,translations.language_code,translations.title,translations.slug,translations.summary,translations.content';

const PAGE_SIZE = 200;
const MAX_ROWS = 10_000;

export interface PostTransferButtonsProps {
  /** 导入/导出完成后的列表刷新回调 */
  onChanged: () => void;
}

export default function PostTransferButtons({ onChanged }: PostTransferButtonsProps) {
  const { t, i18n } = useTranslation('post');
  const { message } = App.useApp();
  const createMutation = useCreatePost();

  const [xlsxOpen, setXlsxOpen] = useState(false);
  const [busy, setBusy] = useState(false);
  const jsonInputRef = useRef<HTMLInputElement>(null);

  /** 全量迭代拉取（页大小 200，10000 条封顶） */
  const fetchAllPosts = async () => {
    const all: Record<string, unknown>[] = [];
    for (let page = 1; all.length < MAX_ROWS; page++) {
      const resp = await fetchListPosts(
        new PaginationQuery({
          paging: { page, pageSize: PAGE_SIZE },
          fieldMask: EXPORT_FIELD_MASK,
          orderBy: ['-id'],
        }),
      );
      const items = (resp?.items ?? []) as Record<string, unknown>[];
      all.push(...items);
      const total = Number(resp?.total) || 0;
      if (items.length < PAGE_SIZE || all.length >= total) break;
    }
    return all;
  };

  const handleExportJson = async () => {
    setBusy(true);
    try {
      const posts = await fetchAllPosts();
      const payload = {
        version: 1,
        exportedAt: new Date().toISOString(),
        count: posts.length,
        posts,
      };
      const blob = new Blob([JSON.stringify(payload, null, 2)], {
        type: 'application/json',
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `posts-export-${new Date().toISOString().slice(0, 10)}.json`;
      a.click();
      URL.revokeObjectURL(url);
      message.success(t('transfer.exportSuccess', { count: posts.length }));
    } catch (error) {
      console.error('Export posts json failed:', error);
      message.error(t('transfer.exportFailed'));
    } finally {
      setBusy(false);
    }
  };

  /** JSON 导入：逐帖 create（强制 DRAFT），单帖失败不阻断其余，末尾汇总 */
  const handleImportJsonFile = async (file: File) => {
    setBusy(true);
    try {
      const text = await file.text();
      const parsed = JSON.parse(text) as {
        posts?: Record<string, unknown>[];
      };
      const posts = Array.isArray(parsed.posts) ? parsed.posts : [];
      if (posts.length === 0) {
        message.warning(t('transfer.noPosts'));
        return;
      }
      let ok = 0;
      const failed: number[] = [];
      for (let i = 0; i < posts.length; i++) {
        const p = posts[i];
        try {
          await createMutation.mutateAsync({
            ...(p as Record<string, unknown>),
            status: 'POST_STATUS_DRAFT',
            id: undefined,
          });
          ok += 1;
        } catch (error) {
          console.error(`Import post #${i} failed:`, error);
          failed.push(i + 1);
        }
      }
      if (failed.length > 0) {
        message.warning(t('transfer.importPartial', { ok, failed: failed.join(',') }));
      } else {
        message.success(t('transfer.importSuccess', { count: ok }));
      }
      onChanged();
    } catch (error) {
      console.error('Import posts json failed:', error);
      message.error(t('transfer.importFailed'));
    } finally {
      setBusy(false);
    }
  };

  /** Excel 行 → create 请求体：分类/标签逗号分隔解析，翻译挂当前语言，状态强制 DRAFT */
  const createFromXlsxRow = async (values: Record<string, string>) => {
    const lang = i18n.language || 'zh-CN';
    const toIds = (v: unknown): number[] | undefined => {
      const ids = String(v ?? '')
        .split(/[,，]/)
        .map((s) => Number(s.trim()))
        .filter((n) => Number.isInteger(n) && n > 0);
      return ids.length > 0 ? ids : undefined;
    };
    return createMutation.mutateAsync({
      editorType: values.editorType || 'EDITOR_TYPE_MARKDOWN',
      status: 'POST_STATUS_DRAFT',
      categoryIds: toIds(values.categoryIds),
      tagIds: toIds(values.tagIds),
      translations: [
        {
          languageCode: lang,
          title: values.title,
          summary: values.summary || '',
          content: values.content || '',
        },
      ],
    });
  };

  return (
    <>
      <Dropdown
        menu={{
          items: [
            { key: 'xlsx', icon: <FileExcelOutlined />, label: t('transfer.importXlsx') },
            { key: 'json', icon: <FileTextOutlined />, label: t('transfer.importJson') },
          ],
          onClick: ({ key }) => {
            if (key === 'xlsx') setXlsxOpen(true);
            if (key === 'json') jsonInputRef.current?.click();
          },
        }}
      >
        <Button icon={<ImportOutlined />} size="small" loading={busy}>
          {t('transfer.import')}
        </Button>
      </Dropdown>

      <Button
        icon={<ExportOutlined />}
        size="small"
        loading={busy}
        onClick={handleExportJson}
      >
        {t('transfer.exportJson')}
      </Button>

      <input
        ref={jsonInputRef}
        type="file"
        accept=".json,application/json"
        style={{ display: 'none' }}
        onChange={(e) => {
          const file = e.target.files?.[0];
          e.target.value = '';
          if (file) handleImportJsonFile(file);
        }}
      />

      <ImportModal
        open={xlsxOpen}
        fields={XLSX_FIELDS}
        createRow={createFromXlsxRow}
        onClose={() => setXlsxOpen(false)}
        onSuccess={() => onChanged()}
      />
    </>
  );
}
