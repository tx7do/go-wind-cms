import { useRef, useState } from 'react';
import type { ActionType } from '@ant-design/pro-components';
import { App, Button, Input, Modal, Select, Space, Switch, Tag } from 'antd';
import {
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  SendOutlined,
  HistoryOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ListTable from '@/components/common/ListTable';
import TableActions from '@/components/common/TableActions';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import type { webhookservicev1_Webhook as Webhook } from '@/api/generated/admin/service/v1';
import {
  WEBHOOK_EVENT_OPTIONS,
  useCreateWebhook,
  useDeleteWebhook,
  useListWebhookDeliveries,
  useListWebhooks,
  useTestWebhookPing,
  useUpdateWebhook,
} from '@/api/hooks/webhook';

/**
 * Webhook 管理页：内容事件外发订阅 CRUD + 测试投递 + 投递记录。
 * 事件为空数组表示订阅全部；投递带 X-Webhook-Signature: sha256=HMAC(body)。
 */
const WebhookManagement = () => {
  const { t } = useTranslation('webhook');
  const { message } = App.useApp();
  const actionRef = useRef<ActionType>(null);

  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<Webhook | undefined>();
  const [form, setForm] = useState<{ url: string; secret: string; events: string[]; enabled: boolean; description: string }>({
    url: '',
    secret: '',
    events: [],
    enabled: true,
    description: '',
  });

  const [deliveriesTarget, setDeliveriesTarget] = useState<Webhook | undefined>();
  const [pingLoadingId, setPingLoadingId] = useState<number | undefined>();

  const listMutation = useListWebhooks();
  const createMutation = useCreateWebhook();
  const updateMutation = useUpdateWebhook();
  const deleteMutation = useDeleteWebhook();
  const pingMutation = useTestWebhookPing();
  const deliveriesQuery = useListWebhookDeliveries(
    { webhookId: deliveriesTarget?.id ?? null },
    { enabled: !!deliveriesTarget },
  );

  const loadItems = async () => {
    const resp = await listMutation.mutateAsync();
    const items = (resp?.items ?? []) as Webhook[];
    setDataSource(items);
    return { data: items, total: Number(resp?.total) || items.length, success: true };
  };

  const openCreate = () => {
    setEditing(undefined);
    setForm({ url: '', secret: '', events: [], enabled: true, description: '' });
    setFormOpen(true);
  };

  const openEdit = (record: Webhook) => {
    setEditing(record);
    setForm({
      url: record.url ?? '',
      secret: record.secret ?? '',
      events: record.events ?? [],
      enabled: record.enabled ?? true,
      description: record.description ?? '',
    });
    setFormOpen(true);
  };

  const handleSave = async () => {
    if (!form.url.trim()) {
      message.warning(t('urlRequired'));
      return;
    }
    try {
      if (editing?.id) {
        await updateMutation.mutateAsync({
          id: editing.id,
          data: { ...form, url: form.url.trim() },
          updateMask: 'url,secret,events,enabled,description',
        });
        message.success(t('saveSuccess'));
      } else {
        await createMutation.mutateAsync({
          data: { ...form, url: form.url.trim() },
        });
        message.success(t('saveSuccess'));
      }
      setFormOpen(false);
      await loadItems();
    } catch (error) {
      console.error('Save webhook failed:', error);
    }
  };

  const handleDelete = async (record: Webhook) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
      message.success(t('deleteSuccess'));
      await loadItems();
    } catch (error) {
      console.error('Delete webhook failed:', error);
    }
  };

  const handleToggle = async (record: Webhook, checked: boolean) => {
    if (!record.id) return;
    try {
      await updateMutation.mutateAsync({
        id: record.id,
        data: { events: record.events ?? [], enabled: checked },
        updateMask: 'enabled',
      });
      await loadItems();
    } catch (error) {
      console.error('Toggle webhook failed:', error);
    }
  };

  const handlePing = async (record: Webhook) => {
    if (!record.id) return;
    setPingLoadingId(record.id);
    try {
      const delivery = await pingMutation.mutateAsync({ id: record.id });
      if (delivery.status === 'DELIVERY_STATUS_SUCCESS') {
        message.success(t('pingSuccess', { status: delivery.httpStatus ?? 0 }));
      } else {
        message.error(t('pingFailed', { error: delivery.error || '-' }));
      }
    } catch (error) {
      console.error('Ping webhook failed:', error);
    } finally {
      setPingLoadingId(undefined);
    }
  };

  const [dataSource, setDataSource] = useState<Webhook[]>([]);

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: t('url'), dataIndex: 'url', ellipsis: true },
    {
      title: t('events'),
      dataIndex: 'events',
      width: 260,
      render: (_: unknown, record: Webhook) => {
        const events = record.events ?? [];
        if (events.length === 0) return <Tag>{t('allEvents')}</Tag>;
        return (
          <Space wrap size={4}>
            {events.slice(0, 3).map((e) => (
              <Tag key={e} color="blue">{e}</Tag>
            ))}
            {events.length > 3 && <Tag>+{events.length - 3}</Tag>}
          </Space>
        );
      },
    },
    {
      title: t('enabled'),
      dataIndex: 'enabled',
      width: 90,
      render: (_: unknown, record: Webhook) => (
        <Switch
          checked={record.enabled ?? false}
          onChange={(checked) => handleToggle(record, checked)}
        />
      ),
    },
    { title: t('description'), dataIndex: 'description', ellipsis: true },
    {
      title: t('action'),
      width: 300,
      render: (_: unknown, record: Webhook) => (
        <TableActions
          record={record}
          actions={[
            {
              key: 'test',
              label: t('test'),
              icon: <SendOutlined />,
              loading: pingLoadingId === record.id,
              onClick: () => handlePing(record),
            },
            {
              key: 'deliveries',
              label: t('deliveries'),
              icon: <HistoryOutlined />,
              onClick: () => setDeliveriesTarget(record),
            },
            {
              key: 'edit',
              label: t('common:button.edit'),
              icon: <EditOutlined />,
              onClick: () => openEdit(record),
            },
            {
              key: 'delete',
              label: t('common:button.delete'),
              icon: <DeleteOutlined />,
              danger: true,
              confirm: { title: t('deleteConfirm') },
              onClick: () => handleDelete(record),
            },
          ]}
        />
      ),
    },
  ];

  return (
    <ContentContainer>
      <ListTable
        headerTitle={t('title')}
        actionRef={actionRef}
        search={false}
        loading={listMutation.isPending}
        rowKey="id"
        columns={columns}
        request={loadItems}
        dataSource={dataSource}
        toolBarRender={() => [
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            {t('create')}
          </Button>,
        ]}
        pagination={{ pageSize: 10 }}
      />

      <Modal
        open={formOpen}
        title={editing ? t('edit') : t('create')}
        onCancel={() => setFormOpen(false)}
        onOk={handleSave}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={560}
        destroyOnHidden
      >
        <Space direction="vertical" style={{ width: '100%' }} size={12}>
          <Input
            placeholder={t('urlPlaceholder')}
            value={form.url}
            onChange={(e) => setForm((f) => ({ ...f, url: e.target.value }))}
          />
          <Input.Password
            placeholder={t('secretPlaceholder')}
            value={form.secret}
            onChange={(e) => setForm((f) => ({ ...f, secret: e.target.value }))}
          />
          <Select
            mode="multiple"
            allowClear
            placeholder={t('eventsPlaceholder')}
            value={form.events}
            options={WEBHOOK_EVENT_OPTIONS as unknown as { value: string; label: string }[]}
            onChange={(v: string[]) => setForm((f) => ({ ...f, events: v }))}
            style={{ width: '100%' }}
          />
          <Input
            placeholder={t('descriptionPlaceholder')}
            value={form.description}
            onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
          />
          <Space>
            <span>{t('enabled')}</span>
            <Switch checked={form.enabled} onChange={(v) => setForm((f) => ({ ...f, enabled: v }))} />
          </Space>
        </Space>
      </Modal>

      <Modal
        open={!!deliveriesTarget}
        title={`${t('deliveries')} - ${deliveriesTarget?.url ?? ''}`}
        onCancel={() => setDeliveriesTarget(undefined)}
        footer={null}
        width={760}
        destroyOnHidden
      >
        <div style={{ maxHeight: 480, overflowY: 'auto' }}>
          {(deliveriesQuery.data?.items ?? []).map((d) => (
            <div
              key={d.id}
              style={{
                display: 'flex',
                gap: 8,
                alignItems: 'center',
                padding: '6px 0',
                borderBottom: '1px solid var(--ant-color-split)',
                fontSize: 12,
              }}
            >
              <Tag color={d.status === 'DELIVERY_STATUS_SUCCESS' ? 'success' : 'error'}>
                {d.status === 'DELIVERY_STATUS_SUCCESS' ? t('statusSuccess') : t('statusFailed')}
              </Tag>
              <Tag>{d.event}</Tag>
              <span>HTTP {d.httpStatus ?? 0}</span>
              <span style={{ color: 'var(--ant-color-text-tertiary)' }}>{d.durationMs ?? 0}ms</span>
              <span style={{ color: 'var(--ant-color-text-tertiary)', marginLeft: 'auto' }}>
                {d.createdAt ? new Date(d.createdAt).toLocaleString() : '-'}
              </span>
              {d.error && (
                <span style={{ color: 'var(--ant-color-error)', maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {d.error}
                </span>
              )}
            </div>
          ))}
          {(deliveriesQuery.data?.items ?? []).length === 0 && !deliveriesQuery.isFetching && (
            <div style={{ textAlign: 'center', padding: 24, color: 'var(--ant-color-text-tertiary)' }}>
              {t('noDeliveries')}
            </div>
          )}
        </div>
      </Modal>
    </ContentContainer>
  );
};

export default WebhookManagement;
