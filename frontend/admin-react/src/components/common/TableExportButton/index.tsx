import { useMemo, useState } from 'react';
import { App, Button, Dropdown } from 'antd';
import { DownloadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { PaginationQuery } from '@/core';
import { exportAuditLogs, type AuditExportFormat } from '@/utils/csv';

/**
 * 列表导出按钮（CSV / Excel 下拉）：内置分页聚合（默认上限 1 万行），
 * 按 ProColumns（dataIndex/title，自动剔除 hideInTable 与操作列）生成文件并触发下载。
 * fetcher 契约与审计导出一致：(query: PaginationQuery) => Promise<{ items, total }>。
 */
export interface TableExportButtonProps {
  fetcher: (query: PaginationQuery) => Promise<{ items?: any[]; total?: number }>;
  /** ProColumns[]（或 {title, dataIndex}[]），自动过滤出可导出列 */
  columns: any[];
  filename?: string;
  maxRows?: number;
}

export default function TableExportButton({
  fetcher,
  columns,
  filename = 'export',
  maxRows = 10_000,
}: TableExportButtonProps) {
  const { t } = useTranslation('common');
  const { message } = App.useApp();
  const [exporting, setExporting] = useState(false);

  const exportColumns = useMemo(
    () =>
      (columns || [])
        .filter(
          (c: any) =>
            c.dataIndex && !c.hideInTable && c.valueType !== 'option',
        )
        .map((c: any) => ({
          key: String(c.dataIndex),
          title:
            typeof c.title === 'string' ? c.title : String(c.dataIndex),
        })),
    [columns],
  );

  const handleExport = async (format: AuditExportFormat) => {
    setExporting(true);
    try {
      const count = await exportAuditLogs(
        { fetcher, filename, columns: exportColumns, maxRows },
        format,
      );
      message.success(t('export.success', { count }));
    } catch (err: any) {
      message.error(`${t('export.failed')}：${err?.message ?? ''}`);
    } finally {
      setExporting(false);
    }
  };

  return (
    <Dropdown
      menu={{
        items: [
          { key: 'csv', label: t('export.csv') },
          { key: 'xlsx', label: t('export.xlsx') },
        ],
        onClick: ({ key }) => handleExport(key as AuditExportFormat),
      }}
    >
      <Button icon={<DownloadOutlined />} loading={exporting}>
        {t('export.title')}
      </Button>
    </Dropdown>
  );
}
