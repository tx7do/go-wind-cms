import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import ListTable from '@/components/common/ListTable';
import TableActions from '@/components/common/TableActions';
import { Tag, App } from 'antd';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import { useDeleteSocialAccount, fetchListSocialAccounts } from '@/api/hooks/social';
import { getPlatformColor, getPlatformLabel, platformOptions } from './constants';
import SocialAccountDrawer from './SocialAccountDrawer';

/**
 * 社交账号列表页：租户在各社交平台上的发布身份管理
 */
const SocialAccountList = () => {
  const { t } = useTranslation('social');
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingRow, setEditingRow] = useState<any>(null);

  const deleteMutation = useDeleteSocialAccount({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listSocialAccounts'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  const columns: ProColumns<any>[] = [
    {
      title: t('accountName'),
      dataIndex: 'name',
      width: 160,
    },
    {
      title: t('platform'),
      dataIndex: 'platform',
      valueType: 'select',
      width: 120,
      fieldProps: {
        options: platformOptions(t),
        showSearch: true,
      },
      render: (_, record) => (
        <Tag color={getPlatformColor(record.platform)}>
          {getPlatformLabel(t, record.platform)}
        </Tag>
      ),
    },
    {
      title: t('handle'),
      dataIndex: 'handle',
      width: 180,
      copyable: true,
    },
    {
      title: t('credentialHint'),
      dataIndex: 'credentialHint',
      search: false,
      width: 180,
      render: (_, record) => (
        <span style={{ fontFamily: 'monospace' }}>{record.credentialHint || '—'}</span>
      ),
    },
    {
      title: t('isEnabled'),
      dataIndex: 'isEnabled',
      search: false,
      width: 90,
      render: (_, record) =>
        record.isEnabled !== false ? (
          <Tag color="green">{t('enabled')}</Tag>
        ) : (
          <Tag color="default">{t('disabled')}</Tag>
        ),
    },
    {
      title: t('remark'),
      dataIndex: 'remark',
      search: false,
      ellipsis: true,
      width: 160,
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      search: false,
      width: 150,
      render: (_, record) =>
        record.createdAt ? dayjs(record.createdAt).format('YYYY-MM-DD HH:mm') : '',
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <TableActions
          record={record}
          actions={[
            {
              key: 'edit',
              label: t('common:button.edit'),
              icon: <EditOutlined />,
              onClick: () => {
                setEditingRow(record);
                setDrawerOpen(true);
              },
            },
            {
              key: 'delete',
              label: t('common:button.delete'),
              icon: <DeleteOutlined />,
              danger: true,
              confirm: {
                title: t('deleteConfirmTitle'),
                description: t('deleteAccountConfirmDesc'),
              },
              onClick: () => record.id && deleteMutation.mutate({ id: record.id }),
            },
          ]}
        />
      ),
    },
  ];

  return (
    <>
      <div ref={containerRef} className="page-container-content" style={{ padding: '0 8px', height: '100%' }}>
        <ListTable<any>
          actionRef={actionRef}
          columns={columns}
          headerTitle={t('accountList')}
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

            const response = await fetchListSocialAccounts(query);

            return {
              data: response.items || [],
              total: Number(response.total) || 0,
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
            <a
              key="create"
              onClick={() => {
                setEditingRow(null);
                setDrawerOpen(true);
              }}
            >
              + {t('createAccount')}
            </a>,
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
          scroll={{ x: 1050, y: tableScrollY }}
        />
      </div>

      <SocialAccountDrawer
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

export default SocialAccountList;
