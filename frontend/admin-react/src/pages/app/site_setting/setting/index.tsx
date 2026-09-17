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
import { fetchListSiteSettings, useDeleteSiteSetting } from '@/api/hooks/site-setting';
import { getTypeColor, getTypeLabel, typeOptions } from './constants';
import SiteSettingDrawer from './SiteSettingDrawer';

/**
 * 站点配置列表页
 */
const SiteSettingList = () => {
  const { t } = useTranslation('site-setting');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [editingRow, setEditingRow] = useState<any>(null);

  const deleteMutation = useDeleteSiteSetting({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listSiteSettings'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  const columns: ProColumns<any>[] = [
    {
      title: t('label'),
      dataIndex: 'label',
      fixed: 'left',
    },
    {
      title: t('key'),
      dataIndex: 'key',
    },
    {
      title: t('value'),
      dataIndex: 'value',
      ellipsis: true,
    },
    {
      title: t('type'),
      dataIndex: 'type',
      valueType: 'select',
      fieldProps: {
        options: typeOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getTypeColor(record.type)}>
          {getTypeLabel(t, record.type)}
        </Tag>
      ),
    },
    {
      title: t('locale'),
      dataIndex: 'locale',
      width: 90,
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 90,
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
          headerTitle={t('settingList')}
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

              const response = await fetchListSiteSettings(query);

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
              fetcher={fetchListSiteSettings}
              columns={columns}
              filename="site-settings"
            />,
            <Button
              key="create"
              type="primary"
              icon={<PlusOutlined />}
              size="small"
              onClick={() => {
                setEditingRow(null);
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

      <SiteSettingDrawer
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
    </>
  );
};

export default SiteSettingList;
