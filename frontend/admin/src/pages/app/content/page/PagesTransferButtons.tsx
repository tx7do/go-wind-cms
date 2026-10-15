import { useState } from 'react';
import { App, Button } from 'antd';
import { ExportOutlined, ImportOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { PaginationQuery } from '@/core';
import { fetchListPages, useCreatePage } from '@/api/hooks/page';

/**
 * 页面内容搬运按钮（JSON 全保真）：
 *  - 导出：全量迭代（页大小 200、1 万条封顶），含 translations + sections 完整数据
 *  - 导入：接受本组件导出格式，逐页 create（区块随 Create 整体写入），强制 DRAFT
 *
 * 页面是 Section 结构化数据，Excel 行式导入不适用（帖子才有 Excel 通道）。
 * 安全口径：导入一律落 DRAFT。
 */

const PAGE_SIZE = 200;
const MAX_ROWS = 10_000;

const EXPORT_FIELD_MASK =
  'id,status,type,editor_type,slug,author_name,disallow_comment,sort_order,thumbnail,custom_fields,content_model_id,created_at,translations.id,translations.page_id,translations.language_code,translations.title,translations.slug,translations.cover_image,translations.full_path,translations.seo,sections.id,sections.type,sections.name,sections.sort_order,sections.config,sections.translations.id,sections.translations.section_id,sections.translations.language_code,sections.translations.content';

export interface PagesTransferButtonsProps {
  /** 导入成功后的列表刷新回调 */
  onChanged: () => void;
}

export default function PagesTransferButtons({ onChanged }: PagesTransferButtonsProps) {
  const { t } = useTranslation('page');
  const { message } = App.useApp();
  const createMutation = useCreatePage();
  const [busy, setBusy] = useState(false);

  const fetchAllPages = async () => {
    const all: Record<string, unknown>[] = [];
    for (let page = 1; all.length < MAX_ROWS; page++) {
      const resp = await fetchListPages(
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
      const pages = await fetchAllPages();
      const payload = {
        version: 1,
        exportedAt: new Date().toISOString(),
        count: pages.length,
        pages,
      };
      const blob = new Blob([JSON.stringify(payload, null, 2)], {
        type: 'application/json',
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `pages-export-${new Date().toISOString().slice(0, 10)}.json`;
      a.click();
      URL.revokeObjectURL(url);
      message.success(t('transfer.exportSuccess', { count: pages.length }));
    } catch (error) {
      console.error('Export pages json failed:', error);
      message.error(t('transfer.exportFailed'));
    } finally {
      setBusy(false);
    }
  };

  const handleImportJsonFile = async (file: File) => {
    setBusy(true);
    try {
      const text = await file.text();
      const parsed = JSON.parse(text) as {
        pages?: Record<string, unknown>[];
      };
      const pages = Array.isArray(parsed.pages) ? parsed.pages : [];
      if (pages.length === 0) {
        message.warning(t('transfer.noPages'));
        return;
      }
      let ok = 0;
      const failed: number[] = [];
      for (let i = 0; i < pages.length; i++) {
        try {
          await createMutation.mutateAsync({
            ...(pages[i] as Record<string, unknown>),
            status: 'PAGE_STATUS_DRAFT',
            id: undefined,
          });
          ok += 1;
        } catch (error) {
          console.error(`Import page #${i} failed:`, error);
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
      console.error('Import pages json failed:', error);
      message.error(t('transfer.importFailed'));
    } finally {
      setBusy(false);
    }
  };

  return (
    <>
      <Button
        icon={<ImportOutlined />}
        size="small"
        loading={busy}
        onClick={() => {
          const input = document.createElement('input');
          input.type = 'file';
          input.accept = '.json,application/json';
          input.onchange = () => {
            const file = input.files?.[0];
            if (file) handleImportJsonFile(file);
          };
          input.click();
        }}
      >
        {t('transfer.importJson')}
      </Button>
      <Button
        icon={<ExportOutlined />}
        size="small"
        loading={busy}
        onClick={handleExportJson}
      >
        {t('transfer.exportJson')}
      </Button>
    </>
  );
}
