import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import ListTable from '@/components/common/ListTable';
import TableActions from '@/components/common/TableActions';
import { App, Button, Form, Input, InputNumber, Modal, Select, Space, Tag } from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { identityservicev1_TenantQuota as TenantQuota } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import {
  fetchListTenantQuotas,
  useDeleteTenantQuota,
  useListQuotaPlans,
  useCreateQuotaPlan,
  useDeleteQuotaPlan,
  useApplyQuotaPlan,
} from '@/api/hooks/tenant-quota';
import { fetchListTenants } from '@/api/hooks/tenant';
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

  // ── 配额套餐模板 ──
  const [planOpen, setPlanOpen] = useState(false);
  const [applyOpen, setApplyOpen] = useState(false);
  const [applyPlan, setApplyPlan] = useState<{ id: number; name: string } | undefined>();
  const [tenantOptions, setTenantOptions] = useState<{ value: number; label: string }[]>([]);
  const [planForm] = Form.useForm();
  const [applyForm] = Form.useForm();

  const plansQuery = useListQuotaPlans();
  const refreshPlans = () => queryClient.invalidateQueries({ queryKey: ['quotaPlans'] });

  const createPlanMutation = useCreateQuotaPlan({
    onSuccess: () => {
      message.success(t('planCreated'));
      setPlanOpen(false);
      planForm.resetFields();
      refreshPlans();
    },
    onError: (error: Error) => {
      console.error('create quota plan failed', error);
      message.error(errMessage(error, t('fetchFailed')));
    },
  });

  const deletePlanMutation = useDeleteQuotaPlan({
    onSuccess: () => {
      message.success(t('planDeleted'));
      refreshPlans();
    },
    onError: (error: Error) => {
      console.error('delete quota plan failed', error);
      message.error(errMessage(error, t('fetchFailed')));
    },
  });

  const applyPlanMutation = useApplyQuotaPlan({
    onSuccess: (resp) => {
      message.success(t('planApplied', { count: Number(resp?.applied ?? 0) }));
      setApplyOpen(false);
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listTenantQuotas'] });
    },
    onError: (error: Error) => {
      console.error('apply quota plan failed', error);
      message.error(errMessage(error, t('fetchFailed')));
    },
  });

  const openApply = async (plan: { id: number; name: string }) => {
    setApplyPlan(plan);
    applyForm.resetFields();
    setApplyOpen(true);
    // 租户选项懒加载（数量小，一次拉全）
    if (tenantOptions.length === 0) {
      try {
        const res = await fetchListTenants(new PaginationQuery({
          paging: { page: 1, pageSize: 100 },
        }));
        setTenantOptions((res.items ?? []).map((x) => ({
          value: Number((x as { id?: unknown }).id),
          label: String((x as { name?: unknown }).name ?? (x as { id?: unknown }).id),
        })));
      } catch (error) {
        console.error('load tenants for plan apply failed:', error);
      }
    }
  };

  const submitPlan = async () => {
    const values = await planForm.validateFields().catch(() => null);
    if (values === null) return;
    const items: {
      quotaType: 'USER_LIMIT' | 'STORAGE' | 'API_CALL' | 'AI_TOKENS';
      limit: number;
    }[] = (['USER_LIMIT', 'STORAGE', 'API_CALL', 'AI_TOKENS'] as const)
      .map((type) => ({ quotaType: type, limit: Number(values[type] ?? 0) }))
      .filter((x) => Number.isFinite(x.limit) && x.limit > 0);
    if (items.length === 0) {
      message.warning(t('planItemRequired'));
      return;
    }
    await createPlanMutation.mutateAsync({
      data: { name: values.name, description: values.description, items },
    });
  };

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
              key="plans"
              icon={<ThunderboltOutlined />}
              onClick={() => {
                setPlanOpen(true);
                plansQuery.refetch();
              }}
            >
              {t('plans')}
            </Button>,
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

      {/* 套餐模板管理弹窗：列表 + 创建 */}
      <Modal
        title={t('plans')}
        open={planOpen}
        footer={null}
        onCancel={() => setPlanOpen(false)}
        width={620}
        destroyOnHidden
      >
        <Form form={planForm} layout="vertical" onFinish={submitPlan}>
          <Space.Compact block style={{ marginBottom: 8 }}>
            <Form.Item
              name="name"
              rules={[{ required: true, message: t('planNameRequired') }]}
              style={{ flex: 1, marginRight: 8, marginBottom: 0 }}
            >
              <Input placeholder={t('planName')} />
            </Form.Item>
            <Form.Item name="description" style={{ flex: 2, marginBottom: 0 }}>
              <Input placeholder={t('planDescription')} />
            </Form.Item>
          </Space.Compact>
          <Space.Compact block style={{ marginBottom: 8 }}>
            {(['USER_LIMIT', 'STORAGE', 'API_CALL', 'AI_TOKENS'] as const).map((type) => (
              <Form.Item
                key={type}
                name={type}
                label={type}
                style={{ flex: 1, marginRight: 6, marginBottom: 0 }}
              >
                <InputNumber min={0} style={{ width: '100%' }} placeholder="0" />
              </Form.Item>
            ))}
          </Space.Compact>
          <Button type="primary" htmlType="submit" loading={createPlanMutation.isPending}>
            {t('planSave')}
          </Button>
        </Form>

        <div style={{ marginTop: 16 }}>
          {(plansQuery.data?.items ?? []).map((p) => (
            <div
              key={p.id}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 8,
                padding: '8px 0',
                borderBottom: '1px solid var(--ant-color-border-secondary, #f0f0f0)',
              }}
            >
              <span style={{ fontWeight: 600, minWidth: 120 }}>{p.name}</span>
              <span style={{ flex: 1, color: 'var(--ant-color-text-secondary, #888)', fontSize: 12 }}>
                {(p.items ?? [])
                  .map((i) => `${i.quotaType}=${Number(i.limit ?? 0).toLocaleString()}`)
                  .join(' / ')}
              </span>
              <Button
                size="small"
                type="primary"
                ghost
                icon={<ThunderboltOutlined />}
                onClick={() => openApply({ id: Number(p.id), name: String(p.name) })}
              >
                {t('apply')}
              </Button>
              <Button
                size="small"
                danger
                icon={<DeleteOutlined />}
                onClick={() =>
                  deletePlanMutation.mutate({ id: Number(p.id) })
                }
              />
            </div>
          ))}
          {(plansQuery.data?.items ?? []).length === 0 && (
            <div style={{ color: 'var(--ant-color-text-tertiary, #aaa)' }}>{t('noPlans')}</div>
          )}
        </div>
      </Modal>

      {/* 套用确认弹窗 */}
      <Modal
        title={`${t('apply')} · ${applyPlan?.name ?? ''}`}
        open={applyOpen}
        onCancel={() => setApplyOpen(false)}
        onOk={async () => {
          const values = await applyForm.validateFields().catch(() => null);
          if (values === null || !applyPlan) return;
          await applyPlanMutation.mutateAsync({
            planId: applyPlan.id,
            tenantId: values.tenantId,
          });
        }}
        confirmLoading={applyPlanMutation.isPending}
        destroyOnHidden
      >
        <p style={{ color: 'var(--ant-color-text-secondary, #888)' }}>{t('applyHint')}</p>
        <Form form={applyForm} layout="vertical">
          <Form.Item name="tenantId" rules={[{ required: true, message: t('tenantRequired') }]}>
            <Select options={tenantOptions} placeholder={t('selectTenant')} />
          </Form.Item>
        </Form>
      </Modal>
    </ContentContainer>
  );
}
