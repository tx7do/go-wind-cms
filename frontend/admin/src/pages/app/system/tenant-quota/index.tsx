import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import ListTable from '@/components/common/ListTable';
import TableActions from '@/components/common/TableActions';
import { App, Button, Tag } from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
} from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { identityservicev1_TenantQuota as TenantQuota } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import {
  fetchListTenantQuotas,
  useDeleteTenantQuota,
} from '@/api/hooks/tenant-quota';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import TenantQuotaDrawer from './components/TenantQuotaDrawer';

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

const quotaTypeColor: Record<string, string> = {
  USER_LIMIT: 'blue',
  STORAGE: 'cyan',
  API_CALL: 'geekblue',
  AI_TOKENS: 'purple',
};

/**
 * 租户配额管理页（平台管理员）：配额按租户直配，闸门与用量展示同源。
 */
export default function TenantQuotaPage() {
  const { t } = useTranslation('tenant-quota');
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [selected, setSelected] = useState<TenantQuota>();

  const deleteMutation = useDeleteTenantQuota({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listTenantQuotas'] });
    },
    onError: (error: Error) => {
      // 原始错误必须留在控制台：用户可见的那句翻译不包含服务端的原因
      console.error('delete tenant quota failed', error);
      message.error(errMessage(error, t('fetchFailed')));
    },
  });

  const columns: ProColumns<TenantQuota>[] = [
    {
      title: t('tenantId'),
      dataIndex: 'tenantId',
      width: 110,
    },
    {
      title: t('quotaType'),
      dataIndex: 'quotaType',
      width: 150,
      render: (_, record) => (
        <Tag color={quotaTypeColor[record.quotaType as string] ?? 'default'}>
          {record.quotaType}
        </Tag>
      ),
    },
    {
      title: t('quotaValue'),
      dataIndex: 'quotaValue',
      width: 180,
      render: (_, record) => Number(record.quotaValue ?? 0).toLocaleString(),
    },
    {
      title: t('updatedAt'),
      dataIndex: 'updatedAt',
      width: 170,
      valueType: 'dateTime',
      hideInSearch: true,
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
                setDrawerMode('edit');
                setSelected(record);
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
                description: t('deleteConfirmDesc'),
              },
              onClick: () => record.id && deleteMutation.mutate({ id: record.id }),
            },
          ]}
        />
      ),
    },
  ];

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <div ref={containerRef} className="page-container-content">
        <ListTable<TenantQuota>
          actionRef={actionRef}
          columns={columns}
          rowKey="id"
          search={{ labelWidth: 'auto', defaultCollapsed: false }}
          scroll={{ y: tableScrollY }}
          pagination={{
            defaultPageSize: 20,
            showSizeChanger: true,
          }}
          request={async (params) => {
            const { current, pageSize, ...rest } = params;
            const query = new PaginationQuery({
              paging: {
                page: current || 1,
                pageSize: pageSize || 20,
              },
              formValues: rest,
            });
            const res = await fetchListTenantQuotas(query);
            return {
              data: res.items || [],
              total: Number(res.total || 0),
              success: true,
            };
          }}
          toolBarRender={() => [
            <Button
              key="create"
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setDrawerMode('create');
                setSelected(undefined);
                setDrawerOpen(true);
              }}
            >
              {t('create')}
            </Button>,
          ]}
          size="middle"
          bordered
        />
      </div>
      <TenantQuotaDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={selected}
        onClose={() => setDrawerOpen(false)}
        onSuccess={() => actionRef.current?.reload()}
      />
    </ContentContainer>
  );
}
