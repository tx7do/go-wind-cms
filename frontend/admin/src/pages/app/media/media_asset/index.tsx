import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import ListTable from '@/components/common/ListTable';
import { Button, Popconfirm, Tag, App, Upload } from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  UploadOutlined,
  CopyOutlined,
  LinkOutlined,
} from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { formatBytes } from '@/utils';
import {
  fetchListMediaAssets,
  useDeleteMediaAsset,
  uploadMediaAsset,
} from '@/api/hooks/media-asset';
import {
  assetTypeOptions,
  getAssetTypeColor,
  getAssetTypeLabel,
  getProcessingStatusColor,
  getProcessingStatusLabel,
  processingStatusOptions,
} from './constants';
import MediaAssetDrawer from './MediaAssetDrawer';
import MediaAssetDuplicatesModal from './MediaAssetDuplicatesModal';
import MediaAssetReferencesModal from './MediaAssetReferencesModal';

/**
 * 媒体资源列表页
 */
const MediaAssetList = () => {
  const { t } = useTranslation('media-asset');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [editingRow, setEditingRow] = useState<any>(null);
  const [duplicatesOpen, setDuplicatesOpen] = useState(false);
  const [referencesRow, setReferencesRow] = useState<{
    id?: number;
    filename?: string;
  } | null>(null);

  const deleteMutation = useDeleteMediaAsset({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listMediaAssets'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  /* 上传媒体资源（antd Upload customRequest 适配） */
  const handleUploadMediaAsset = async (options: any) => {
    const { file, onSuccess, onError } = options;
    try {
      await uploadMediaAsset(
        {
          fileDirectory: 'media',
          title: file.name,
        },
        file,
      );
      onSuccess?.({}, file);
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listMediaAssets'] });
      message.success(t('uploadSuccess'));
    } catch (error) {
      console.error('上传媒体资源失败', error);
      try {
        onError?.(error, file);
      } catch {}
      const dupHint =
        error instanceof Error ? error.message : String(error ?? '');
      if (dupHint.includes('duplicate of media asset')) {
        message.warning(`${t('uploadDuplicate')}: ${dupHint}`);
      } else {
        message.error(t('uploadFailed'));
      }
    }
  };

  const rangePresets = [
    { label: t('dateRange.today'), value: [dayjs().startOf('day'), dayjs().endOf('day')] },
    {
      label: t('dateRange.yesterday'),
      value: [dayjs().subtract(1, 'day').startOf('day'), dayjs().subtract(1, 'day').endOf('day')],
    },
    { label: t('dateRange.thisWeek'), value: [dayjs().startOf('week'), dayjs().endOf('week')] },
    {
      label: t('dateRange.lastWeek'),
      value: [dayjs().subtract(1, 'week').startOf('week'), dayjs().subtract(1, 'week').endOf('week')],
    },
    { label: t('dateRange.thisMonth'), value: [dayjs().startOf('month'), dayjs().endOf('month')] },
    {
      label: t('dateRange.lastMonth'),
      value: [dayjs().subtract(1, 'month').startOf('month'), dayjs().subtract(1, 'month').endOf('month')],
    },
  ];

  const columns: ProColumns<any>[] = [
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      search: false,
      width: 150,
      render: (_, record) =>
        record.createdAt ? dayjs(record.createdAt).format('YYYY-MM-DD HH:mm') : '',
    },
    {
      // 搜索专用：时间范围 → created_at__gte/lte
      title: t('createdAtRange'),
      dataIndex: 'created_at_range',
      valueType: 'dateTimeRange',
      hideInTable: true,
      fieldProps: { showTime: true, presets: rangePresets },
      search: {
        transform: (value: any[]) => ({
          created_at__gte: value?.[0]
            ? dayjs(value[0]).format('YYYY-MM-DD HH:mm:ss')
            : undefined,
          created_at__lte: value?.[1]
            ? dayjs(value[1]).format('YYYY-MM-DD HH:mm:ss')
            : undefined,
        }),
      },
    },
    {
      title: t('filename'),
      dataIndex: 'filename',
    },
    {
      title: t('type'),
      dataIndex: 'type',
      valueType: 'select',
      fieldProps: {
        options: assetTypeOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getAssetTypeColor(record.type)}>
          {getAssetTypeLabel(t, record.type)}
        </Tag>
      ),
    },
    {
      title: t('size'),
      dataIndex: 'size',
      search: false,
      width: 100,
      render: (_, record) => formatBytes(Number(record.size || 0)),
    },
    {
      title: t('storagePath'),
      dataIndex: 'storagePath',
      search: false,
      ellipsis: true,
    },
    {
      title: t('title'),
      dataIndex: 'title',
      search: false,
    },
    {
      title: t('caption'),
      dataIndex: 'caption',
      search: false,
      ellipsis: true,
    },
    {
      title: t('processingStatus'),
      dataIndex: 'processingStatus',
      valueType: 'select',
      fieldProps: {
        options: processingStatusOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getProcessingStatusColor(record.processingStatus)}>
          {getProcessingStatusLabel(t, record.processingStatus)}
        </Tag>
      ),
    },
    {
      title: t('referenceCount'),
      dataIndex: 'referenceCount',
      search: false,
      width: 100,
      render: (_, record) => (
        <a onClick={() => setReferencesRow(record)}>
          {Number(record.referenceCount || 0)}
        </a>
      ),
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 140,
      fixed: 'right',
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setEditingRow(record);
            setDrawerMode('edit');
            setDrawerOpen(true);
          }}
        >
          <EditOutlined />
        </a>,
        <a
          key="references"
          title={t('references.title')}
          onClick={() => setReferencesRow(record)}
        >
          <LinkOutlined />
        </a>,
        <Popconfirm
          key="delete"
          title={t('deleteConfirmTitle')}
          description={t('deleteConfirmDesc', { moduleName: t('moduleName') })}
          onConfirm={() => {
            record.id && deleteMutation.mutate({ id: record.id });
          }}
          okText={t('common:button.ok')}
          cancelText={t('common:button.cancel')}
        >
          <a style={{ color: 'var(--ant-color-error)' }}>
            <DeleteOutlined />
          </a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <>
      <div ref={containerRef} className="page-container-content" style={{ padding: '0 8px', height: '100%' }}>
        <ListTable<any>
          actionRef={actionRef}
          columns={columns}
          headerTitle={t('assetList')}
          request={async (params) => {
            const query = new PaginationQuery({
              paging: {
                page: params.current || 1,
                pageSize: params.pageSize || TABLE.DEFAULT_PAGE_SIZE,
              },
              formValues: Object.fromEntries(
                Object.entries(params).filter(
                  ([key]) => !['current', 'pageSize'].includes(key),
                ),
              ),
              orderBy: ['-created_at'],
            });

            const response = await fetchListMediaAssets(query);

            return {
              data: response.items || [],
              total: response.total || 0,
              success: true,
            };
          
}}
          rowKey="id"
          search={{
            labelWidth: 'auto',
            defaultCollapsed: false,
            span: 12,
          }}
          pagination={{
            defaultPageSize: TABLE.DEFAULT_PAGE_SIZE,
            showSizeChanger: true,
            showQuickJumper: true,
          }}
          toolBarRender={() => [
            <Button
              key="duplicates"
              icon={<CopyOutlined />}
              size="small"
              onClick={() => setDuplicatesOpen(true)}
            >
              {t('duplicates.toolbar')}
            </Button>,
            <Upload
              key="upload"
              multiple={false}
              customRequest={handleUploadMediaAsset}
              showUploadList={false}
            >
              <Button type="primary" icon={<UploadOutlined />} size="small">
                {t('upload')}
              </Button>
            </Upload>,
            <TableExportButton
              key="export"
              fetcher={fetchListMediaAssets}
              columns={columns}
              filename="media-assets"
            />,
          ]}
          options={{
            density: false,
            fullScreen: false,
            setting: false,
            reload: true,
          }}
          size="small"
          bordered
          cardBordered={false}
          scroll={{ x: 1200, y: tableScrollY }}
        />
      </div>

      <MediaAssetDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={editingRow}
        onClose={() => {
          setDrawerOpen(false);
          setEditingRow(null);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />

      <MediaAssetDuplicatesModal
        open={duplicatesOpen}
        onClose={() => setDuplicatesOpen(false)}
        onDeleteAsset={(id) => {
          deleteMutation.mutate({ id });
          setDuplicatesOpen(false);
        }}
      />

      <MediaAssetReferencesModal
        open={!!referencesRow}
        mediaAssetId={referencesRow?.id}
        mediaAssetName={referencesRow?.filename}
        onClose={() => setReferencesRow(null)}
      />
    </>
  );
};

export default MediaAssetList;
