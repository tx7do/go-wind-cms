import { Modal, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import {
  type mediaservicev1_MediaAsset,
  type mediaservicev1_MediaAssetDuplicateGroup,
} from '@/api/generated/admin/service/v1';
import { useListDuplicatedMediaAssets } from '@/api/hooks/media-asset';
import { formatBytes } from '@/utils';

/**
 * 重复资源弹窗：同租户内同内容哈希 >1 条的分组视图。
 * 组内按创建时间升序——最早上传者建议保留，其余可清理。
 */
const MediaAssetDuplicatesModal: React.FC<{
  open: boolean;
  onClose: () => void;
  onDeleteAsset?: (id: number) => void;
}> = ({ open, onClose, onDeleteAsset }) => {
  const { t } = useTranslation('media-asset');
  // limit 传 0 表示不拉取（弹窗未开时禁用查询）
  const { data, isLoading } = useListDuplicatedMediaAssets(open ? 50 : 0);
  const groups = data?.groups || [];

  const columns: ColumnsType<mediaservicev1_MediaAsset> = [
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      width: 140,
      render: (value) =>
        value ? dayjs(String(value)).format('YYYY-MM-DD HH:mm') : '',
    },
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: t('filename'), dataIndex: 'filename', ellipsis: true },
    {
      title: t('size'),
      dataIndex: 'size',
      width: 90,
      render: (value) => formatBytes(Number(value || 0)),
    },
    {
      title: t('duplicates.keepOrDelete'),
      width: 130,
      render: (_, record, index) =>
        index === 0 ? (
          <Tag color="green">{t('duplicates.keepFirst')}</Tag>
        ) : (
          <a
            style={{ color: 'var(--ant-color-error)' }}
            onClick={() => record.id && onDeleteAsset?.(record.id)}
          >
            {t('duplicates.deleteThis')}
          </a>
        ),
    },
  ];

  return (
    <Modal
      title={t('duplicates.title')}
      open={open}
      onCancel={onClose}
      footer={null}
      width={820}
      destroyOnHidden
    >
      {isLoading ? (
        <div style={{ textAlign: 'center', padding: 32 }}>...</div>
      ) : groups.length === 0 ? (
        <div
          style={{
            textAlign: 'center',
            padding: 32,
            color: 'var(--ant-color-text-secondary)',
          }}
        >
          {t('duplicates.empty')}
        </div>
      ) : (
        groups.map((g: mediaservicev1_MediaAssetDuplicateGroup) => (
          <div key={g.fileHash} style={{ marginBottom: 20 }}>
            <div
              style={{
                marginBottom: 8,
                fontSize: 12,
                color: 'var(--ant-color-text-secondary)',
                wordBreak: 'break-all',
              }}
            >
              <Tag color="orange">
                {t('duplicates.count', { count: Number(g.count ?? 0) })}
              </Tag>
              {g.fileHash}
            </div>
            <Table<mediaservicev1_MediaAsset>
              rowKey="id"
              size="small"
              pagination={false}
              columns={columns}
              dataSource={g.assets || []}
            />
          </div>
        ))
      )}
    </Modal>
  );
};

export default MediaAssetDuplicatesModal;
