import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App, Empty } from 'antd';
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { fetchListPermissions, useDeletePermission } from '@/api/hooks/permission';
import { getStatusMap, getStatusOptions } from './constants';
import PermissionDrawer from './PermissionDrawer';

interface PermissionListProps {
  groupId: number | null;
}

/**
 * 权限点列表
 */
const PermissionList: React.FC<PermissionListProps> = ({ groupId }) => {
  const { t } = useTranslation('permission');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [editingPermission, setEditingPermission] = useState<any>(null);

  const statusMap = getStatusMap(t);

  // 删除 mutation
  const deleteMutation = useDeletePermission({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listPermissions'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  // 列配置
  const columns: ProColumns<any>[] = [
    {
      title: t('name'),
      dataIndex: 'name',
      width: 150,
      fixed: 'left',
    },
    {
      title: t('code'),
      dataIndex: 'code',
      width: 150,
      fixed: 'left',
    },
    {
      title: t('groupName'),
      dataIndex: 'groupName',
      width: 120,
      hideInSearch: true,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      width: 95,
      valueType: 'select',
      fieldProps: {
        options: getStatusOptions(t),
      },
      render: (_, record) => {
        const status = record.status as keyof typeof statusMap;
        const config = statusMap[status] || { text: status, color: 'default' };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 90,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setEditingPermission(record);
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
          onConfirm={() => record.id && deleteMutation.mutate({ id: record.id })}
          okText={t('common:button.ok')}
          cancelText={t('common:button.cancel')}
        >
          <a style={{ color: 'var(--ant-color-error)' }}><DeleteOutlined /></a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <>
      <div ref={containerRef} className="page-container-content" style={{ padding: '0 8px', height: '100%' }}>
        {groupId ? (
          <ProTable<any>
            actionRef={actionRef}
            columns={columns}
            headerTitle={false}
            params={{ groupId }}
            request={async (params) => {
              try {
                const query = new PaginationQuery({
                  paging: {
                    page: params.current || 1,
                    pageSize: params.pageSize || TABLE.DEFAULT_PAGE_SIZE,
                  },
                  formValues: {
                    ...Object.fromEntries(
                      Object.entries(params).filter(
                        ([key]) => !['current', 'pageSize', 'groupId'].includes(key),
                      ),
                    ),
                    group_id: groupId,
                  },
                });

                const response = await fetchListPermissions(query);

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
            }}
            pagination={{
              defaultPageSize: TABLE.DEFAULT_PAGE_SIZE,
              showSizeChanger: true,
              showQuickJumper: true,
            }}
            toolBarRender={() => [
              <TableExportButton key="export" fetcher={fetchListPermissions} columns={columns} filename="permissions" />,
              <Button
                key="create"
                type="primary"
                icon={<PlusOutlined />}
                size="small"
                onClick={() => {
                  setEditingPermission(null);
                  setDrawerMode('create');
                  setDrawerOpen(true);
                }}
              >
                {t('create')}
              </Button>,
            ]}
            options={{
              density: true,
              fullScreen: true,
              setting: true,
              reload: true,
            }}
            size="small"
            bordered
            cardBordered={false}
            scroll={{ y: tableScrollY }}
          />
        ) : (
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%' }}>
            <Empty description={t('selectGroupFirst')} />
          </div>
        )}
      </div>

      <PermissionDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={editingPermission}
        groupId={groupId!}
        onClose={() => {
          setDrawerOpen(false);
          setEditingPermission(null);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />
    </>
  );
};

export default PermissionList;
