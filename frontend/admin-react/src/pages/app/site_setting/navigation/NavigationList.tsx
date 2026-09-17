import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import {
  fetchListNavigations,
  useDeleteNavigation,
} from '@/api/hooks/navigation';
import { getBoolLabel, getLocationColor, getLocationLabel, locationOptions } from './constants';
import NavigationDrawer from './NavigationDrawer';

interface NavigationListProps {
  currentNavigationId: number | null;
  onNavigationSelect: (navigationId: number) => void;
}

/**
 * 导航列表（主表）：点击行联动右侧导航项列表
 */
const NavigationList: React.FC<NavigationListProps> = ({
  currentNavigationId,
  onNavigationSelect,
}) => {
  const { t } = useTranslation('navigation');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [editingRow, setEditingRow] = useState<any>(null);


  const deleteMutation = useDeleteNavigation({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listNavigations'] });
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
      title: t('location'),
      dataIndex: 'location',
      valueType: 'select',
      fieldProps: {
        options: locationOptions(t),
      },
      render: (_, record) => (
        <Tag color={getLocationColor(record.location)}>
          {getLocationLabel(t, record.location)}
        </Tag>
      ),
    },
    {
      title: t('locale'),
      dataIndex: 'locale',
    },
    {
      title: t('isActive'),
      dataIndex: 'isActive',
      render: (_, record) => getBoolLabel(t, record.isActive),
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      width: 150,
      render: (_, record) =>
        record.createdAt ? dayjs(record.createdAt).format('YYYY-MM-DD HH:mm') : '',
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 90,
      render: (_, record) => [
        <a
          key="edit"
          onClick={(e) => {
            e.stopPropagation();
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
          <a
            style={{ color: 'var(--ant-color-error)' }}
            onClick={(e) => e.stopPropagation()}
          >
            <DeleteOutlined />
          </a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <>
      <div
        ref={containerRef}
        className="page-container-content"
        style={{ padding: '0 4px', height: '100%', display: 'flex', flexDirection: 'column' }}
      >
        <ProTable<any>
          actionRef={actionRef}
          columns={columns}
          headerTitle={t('moduleName')}
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

              const response = await fetchListNavigations(query);

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
          }}
          toolBarRender={() => [
            <TableExportButton
              key="export"
              fetcher={fetchListNavigations}
              columns={columns}
              filename="navigations"
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
          onRow={(record) => ({
            onClick: () => {
              if (typeof record.id === 'number') {
                onNavigationSelect(record.id);
              }
            },
            style: { cursor: 'pointer' },
            className:
              record.id === currentNavigationId
                ? 'ant-table-row-selected'
                : undefined,
          })}
        />
      </div>

      <NavigationDrawer
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

export default NavigationList;
