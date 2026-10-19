import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import ListTable from '@/components/common/ListTable';
import { Button, Modal, Popconfirm, Tag, Typography, App } from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { access_keyservicev1_AccessKey as AccessKey } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import {
  fetchListAccessKeys,
  useDeleteAccessKey,
  useResetAccessKeySecret,
} from '@/api/hooks/access-key';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import AccessKeyDrawer from './components/AccessKeyDrawer';

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

/**
 * OpenAPI 访问凭证（AK/SK）管理页。
 * SK 只在创建/重置响应中明文出现一次，服务端只存 SHA-256 摘要。
 */
export default function AccessKeyPage() {
  const { t } = useTranslation('access-key');
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const actionRef = useRef<ActionType>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [selected, setSelected] = useState<AccessKey>();

  const resetMutation = useResetAccessKeySecret();
  const [resetSecretFor, setResetSecretFor] = useState<AccessKey>();
  const [resetSecretValue, setResetSecretValue] = useState<string | null>(null);

  const deleteMutation = useDeleteAccessKey({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listAccessKeys'] });
    },
    onError: (error: Error) => {
      // 原始错误必须留在控制台：用户可见的那句翻译不包含服务端的原因
      console.error('delete access key failed', error);
      message.error(errMessage(error, t('fetchFailed')));
    },
  });

  const statusMap: Record<string, { text: string; color: string }> = {
    ON: { text: t('statusMap.ON'), color: 'success' },
    OFF: { text: t('statusMap.OFF'), color: 'default' },
  };

  const columns: ProColumns<AccessKey>[] = [
    { title: t('name'), dataIndex: 'name', ellipsis: true },
    {
      title: t('accessKey'),
      dataIndex: 'accessKey',
      copyable: true,
      ellipsis: true,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      width: 90,
      render: (_, record) => {
        const cfg = statusMap[record.status as string] ?? {
          text: record.status,
          color: 'default',
        };
        return <Tag color={cfg.color}>{cfg.text}</Tag>;
      },
    },
    {
      title: t('expiresAt'),
      dataIndex: 'expiresAt',
      width: 170,
      valueType: 'dateTime',
      hideInSearch: true,
    },
    {
      title: t('lastUsedAt'),
      dataIndex: 'lastUsedAt',
      width: 170,
      valueType: 'dateTime',
      hideInSearch: true,
    },
    {
      title: t('createdAt'),
      dataIndex: 'createdAt',
      width: 170,
      valueType: 'dateTime',
      hideInSearch: true,
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 90,
      fixed: 'right',
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setDrawerMode('edit');
            setSelected(record);
            setDrawerOpen(true);
          }}
        >
          <EditOutlined />
        </a>,
        <a
          key="reset"
          title={t('resetSecret')}
          onClick={() => {
            if (!record.id) return;
            resetMutation.mutate(
              { id: record.id },
              {
                onSuccess: (resp) => {
                  message.success(t('resetSuccess'));
                  setResetSecretFor(record);
                  setResetSecretValue(resp.secret ?? null);
                },
                onError: (error: Error) => {
                  console.error('reset access key secret failed', error);
                  message.error(errMessage(error, t('resetFailed')));
                },
              },
            );
          }}
        >
          <ReloadOutlined />
        </a>,
        <Popconfirm
          key="delete"
          title={t('deleteConfirmTitle')}
          description={t('deleteConfirmDesc', { moduleName: t('moduleName') })}
          onConfirm={() =>
            record.id && deleteMutation.mutate({ id: record.id })
          }
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
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <div ref={containerRef} className="page-container-content">
        <ListTable<AccessKey>
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
            const res = await fetchListAccessKeys(query);
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
      <Modal
        title={t('secretDialogTitle')}
        open={resetSecretValue !== null}
        onOk={() => {
          setResetSecretValue(null);
          setResetSecretFor(undefined);
          actionRef.current?.reload();
        }}
        onCancel={() => {
          setResetSecretValue(null);
          setResetSecretFor(undefined);
        }}
        cancelButtonProps={{ style: { display: 'none' } }}
        width={560}
      >
        <Typography.Paragraph type="warning">
          {t('secretDialogHint')}
        </Typography.Paragraph>
        <Typography.Paragraph copyable={{ text: resetSecretValue ?? '' }}>
          <Typography.Text code style={{ wordBreak: 'break-all' }}>
            {resetSecretValue}
          </Typography.Text>
        </Typography.Paragraph>
        <Typography.Text type="secondary">
          {t('resetSecretFor')}:&nbsp;
          {resetSecretFor?.accessKey}
        </Typography.Text>
      </Modal>
      <AccessKeyDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={selected}
        onClose={() => setDrawerOpen(false)}
        onSuccess={() => actionRef.current?.reload()}
      />
    </ContentContainer>
  );
}
