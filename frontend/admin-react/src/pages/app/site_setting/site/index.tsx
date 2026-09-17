import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { fetchListSites, useDeleteSite } from '@/api/hooks/site';
import { getStatusColor, getStatusLabel, statusOptions, getBoolLabel } from './constants';
import SiteDrawer from './SiteDrawer';

/**
 * 站点管理列表页
 */
const SiteList = () => {
  const { t } = useTranslation('site');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [editingSite, setEditingSite] = useState<any>(null);

  const deleteMutation = useDeleteSite({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listSites'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  const columns: ProColumns<any>[] = [
    {
      title: t('name'),
      dataIndex: 'name',
      fixed: 'left',
    },
    {
      title: t('slug'),
      dataIndex: 'slug',
    },
    {
      title: t('domain'),
      dataIndex: 'domain',
    },
    {
      title: t('status'),
      dataIndex: 'status',
      valueType: 'select',
      fieldProps: {
        options: statusOptions(t),
      },
      render: (_, record) => (
        <Tag color={getStatusColor(record.status)}>
          {getStatusLabel(t, record.status)}
        </Tag>
      ),
    },
    {
      title: t('defaultLocale'),
      dataIndex: 'defaultLocale',
    },
    {
      title: t('template'),
      dataIndex: 'template',
    },
    {
      title: t('theme'),
      dataIndex: 'theme',
    },
    {
      title: t('isDefault'),
      dataIndex: 'isDefault',
      render: (_, record) => getBoolLabel(t, record.isDefault),
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 90,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setEditingSite(record);
            setDrawerMode('edit');
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
          headerTitle={t('siteList')}
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

              const response = await fetchListSites(query);

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
            <TableExportButton key="export" fetcher={fetchListSites} columns={columns} filename="sites" />,
            <Button
              key="create"
              type="primary"
              icon={<PlusOutlined />}
              size="small"
              onClick={() => {
                setEditingSite(null);
                setDrawerMode('create');
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

      <SiteDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={editingSite}
        onClose={() => {
          setDrawerOpen(false);
          setEditingSite(null);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />
    </>
  );
};

export default SiteList;
