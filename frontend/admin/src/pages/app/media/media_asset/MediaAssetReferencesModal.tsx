import { Modal, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { type mediaservicev1_MediaAssetReference } from '@/api/generated/admin/service/v1';
import { useListMediaAssetReferences } from '@/api/hooks/media-asset';

/**
 * 媒体引用关系弹窗：哪些文章/页面在使用该资源（引用图谱查询面）。
 */
const MediaAssetReferencesModal: React.FC<{
  open: boolean;
  mediaAssetId?: number;
  mediaAssetName?: string;
  onClose: () => void;
}> = ({ open, mediaAssetId, mediaAssetName, onClose }) => {
  const { t } = useTranslation('media-asset');
  const { data, isLoading } = useListMediaAssetReferences(
    open ? mediaAssetId : undefined,
  );
  const items = data?.items || [];

  const columns: ColumnsType<mediaservicev1_MediaAssetReference> = [
    {
      title: t('references.sourceType'),
      dataIndex: 'sourceType',
      width: 90,
      render: (value) => (
        <Tag color={value === 'post' ? 'blue' : 'purple'}>
          {value === 'post' ? t('references.post') : t('references.page')}
        </Tag>
      ),
    },
    { title: 'ID', dataIndex: 'sourceId', width: 70 },
    {
      title: t('references.sourceTitle'),
      dataIndex: 'sourceTitle',
      ellipsis: true,
      render: (value) =>
        value || (
          <span style={{ color: 'var(--ant-color-text-quaternary)' }}>
            {t('references.sourceMissing')}
          </span>
        ),
    },
    {
      title: t('references.field'),
      dataIndex: 'field',
      width: 110,
      render: (value) => {
        const map: Record<string, [string, string]> = {
          body: [t('references.fieldBody'), 'blue'],
          thumbnail: [t('references.fieldThumbnail'), 'cyan'],
          cover: [t('references.fieldCover'), 'geekblue'],
        };
        const [label, color] =
          map[String(value)] || [String(value ?? ''), 'default'];
        return <Tag color={color}>{label}</Tag>;
      },
    },
    {
      title: t('updatedAt'),
      dataIndex: 'updatedAt',
      width: 140,
      render: (value) =>
        value ? dayjs(String(value)).format('YYYY-MM-DD HH:mm') : '',
    },
  ];

  return (
    <Modal
      title={
        mediaAssetName
          ? `${t('references.title')} · ${mediaAssetName}`
          : t('references.title')
      }
      open={open}
      onCancel={onClose}
      footer={null}
      width={720}
      destroyOnHidden
    >
      <Table<mediaservicev1_MediaAssetReference>
        rowKey={(r) => `${r.sourceType}-${r.sourceId}-${r.field}`}
        size="small"
        loading={isLoading}
        pagination={false}
        columns={columns}
        dataSource={items}
        locale={{ emptyText: t('references.empty') }}
      />
    </Modal>
  );
};

export default MediaAssetReferencesModal;
