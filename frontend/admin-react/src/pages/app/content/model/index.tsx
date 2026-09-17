import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import {
  fetchListContentModels,
  useDeleteContentModel,
} from '@/api/hooks/content-model';
import ContentModelDrawer from './ContentModelDrawer';

/**
 * 内容模型列表页
 */
const ContentModelList = () => {
  const { t } = useTranslation('content-model');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingRow, setEditingRow] = useState<any>(null);

  const deleteMutation = useDeleteContentModel({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listContentModels'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  const columns: ProColumns<any>[] = [
    {
      title: t('name'),
      dataIndex: 'name',
    },
    {
      title: t('code'),
      dataIndex: 'code',
    },
    {
      title: t('description'),
      dataIndex: 'description',
      ellipsis: true,
      render: (_, record) => {
        const text = String(record.description ?? '');
        return text.length > 50 ? `${text.slice(0, 50)}...` : text;
      },
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      search: false,
      width: 160,
      render: (_, record) =>
        record.createdAt ? dayjs(record.createdAt).format('YYYY-MM-DD HH:mm') : '',
    },
    {
      // 搜索专用：时间范围 → created_at__gte/lte
      title: t('createdAtRange'),
      dataIndex: 'created_at_range',
      valueType: 'dateTimeRange',
      hideInTable: true,
      fieldProps: { showTime: true, presets: [
        { label: t('dateRange.today'), value: [dayjs().startOf('day'), dayjs().endOf('day')] },
        { label: t('dateRange.yesterday'), value: [dayjs().subtract(1, 'day').startOf('day'), dayjs().subtract(1, 'day').endOf('day')] },
        { label: t('dateRange.thisWeek'), value: [dayjs().startOf('week'), dayjs().endOf('week')] },
        { label: t('dateRange.lastWeek'), value: [dayjs().subtract(1, 'week').startOf('week'), dayjs().subtract(1, 'week').endOf('week')] },
        { label: t('dateRange.thisMonth'), value: [dayjs().startOf('month'), dayjs().endOf('month')] },
        { label: t('dateRange.lastMonth'), value: [dayjs().subtract(1, 'month').startOf('month'), dayjs().subtract(1, 'month').endOf('month')] },
      ] },
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
      title: t('action'),
      valueType: 'option',
      width: 100,
      fixed: 'right',
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setEditingRow(record);
            setDrawerOpen(true);
          }}
        >
          <EditOutlined />
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
        <ProTable<any>
          actionRef={actionRef}
          columns={columns}
          headerTitle={t('modelList')}
          request={async (params) => {
            try {
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

              const response = await fetchListContentModels(query);

              return {
                data: response.items || [],
                total: response.total || 0,
                success: true,
              };
            } catch (error: any) {
              message.error(error.message || t('fetchFailed'));
              return { data: [], total: 0, success: false };
            }
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
            <TableExportButton
              key="export"
              fetcher={fetchListContentModels}
              columns={columns}
              filename="content-models"
            />,
            <Button
              key="create"
              type="primary"
              icon={<PlusOutlined />}
              size="small"
              onClick={() => {
                setEditingRow(null);
                setDrawerOpen(true);
              }}
            >
              {t('create')}
            </Button>,
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
          scroll={{ y: tableScrollY }}
        />
      </div>

      <ContentModelDrawer
        open={drawerOpen}
        data={editingRow}
        onClose={() => {
          setDrawerOpen(false);
          setEditingRow(null);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />
    </>
  );
};

export default ContentModelList;
