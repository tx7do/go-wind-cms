import { useState } from 'react';
import { App, Button, Tag } from 'antd';
import { DownloadOutlined, ReloadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ListTable from '@/components/common/ListTable';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { fetchDownloadBackup, useListBackups, type BackupArchive } from '@/api/hooks/backup';

function formatBytes(n: number): string {
  if (!n) return '-';
  const units = ['B', 'KB', 'MB', 'GB'];
  let v = n;
  let i = 0;
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024;
    i += 1;
  }
  return `${v.toFixed(v >= 100 || i === 0 ? 0 : 1)} ${units[i]}`;
}

/**
 * 备份档案管理页（只读）：backup asynq 任务写入 OSS 的历史档案列表与下载。
 * 恢复无页面入口（有意防误触，走一次性程序 + 演练文档）；保留期清理为任务侧行为。
 */
export default function BackupManagement() {
  const { t } = useTranslation('backup');
  const { message } = App.useApp();
  const listQuery = useListBackups();
  const [downloading, setDownloading] = useState<string | undefined>();

  const dataSource = listQuery.data?.items ?? [];

  const handleDownload = async (record: BackupArchive) => {
    setDownloading(record.objectKey);
    try {
      await fetchDownloadBackup(record);
    } catch (error) {
      console.error('Download backup failed:', error);
      message.error((error as Error)?.message || t('downloadFailed'));
    } finally {
      setDownloading(undefined);
    }
  };

  const columns = [
    {
      title: t('displayName'),
      dataIndex: 'displayName',
      ellipsis: true,
    },
    {
      title: t('size'),
      dataIndex: 'sizeBytes',
      width: 110,
      render: (_: unknown, record: BackupArchive) => formatBytes(Number(record.sizeBytes || 0)),
    },
    {
      title: t('lastModified'),
      dataIndex: 'lastModified',
      width: 180,
      render: (_: unknown, record: BackupArchive) =>
        record.lastModified ? new Date(record.lastModified).toLocaleString() : '-',
    },
    {
      title: t('contentType'),
      dataIndex: 'objectKey',
      width: 120,
      render: () => <Tag color="blue">application/gzip</Tag>,
    },
    {
      title: t('action'),
      width: 110,
      render: (_: unknown, record: BackupArchive) => (
        <Button
          size="small"
          icon={<DownloadOutlined />}
          loading={downloading === record.objectKey}
          onClick={() => handleDownload(record)}
        >
          {t('download')}
        </Button>
      ),
    },
  ];

  return (
    <ContentContainer>
      <ListTable<BackupArchive>
        headerTitle={t('title')}
        search={false}
        loading={listQuery.isPending}
        rowKey="objectKey"
        columns={columns}
        request={async () => {
          const items = dataSource;
          return { data: items, total: Number(listQuery.data?.total) || items.length, success: true };
        }}
        dataSource={dataSource}
        toolBarRender={() => [
          <Button key="reload" icon={<ReloadOutlined />} onClick={() => listQuery.refetch()}>
            {t('reload')}
          </Button>,
        ]}
        pagination={{ pageSize: 15 }}
      />
    </ContentContainer>
  );
}
