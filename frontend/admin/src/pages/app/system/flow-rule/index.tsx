import { useRef, useState } from 'react';
import type { ActionType } from '@ant-design/pro-components';
import { App, Button, Input, Modal, Popconfirm, Select, Space, Switch, Tag } from 'antd';
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ListTable from '@/components/common/ListTable';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import type { flowservicev1_FlowRule as FlowRule } from '@/api/generated/admin/service/v1';
import {
  FLOW_EVENT_OPTIONS,
  FLOW_STATUS_CONDITION_OPTIONS,
  useCreateFlowRule,
  useDeleteFlowRule,
  useListFlowRules,
  useUpdateFlowRule,
} from '@/api/hooks/flow-rule';

/** v1 动作编辑模型：webhook 动作（url + 可选签名密钥） */
interface WebhookAction {
  type: 'webhook';
  url: string;
  secret: string;
}

const serializeActions = (webhooks: WebhookAction[]) =>
  JSON.stringify(webhooks.filter((w) => w.url.trim()));

const deserializeActions = (actionsJSON: string): WebhookAction[] => {
  try {
    const raw: unknown = JSON.parse(actionsJSON || '[]');
    if (!Array.isArray(raw)) return [];
    return raw
      .filter((a): a is { type?: string; url?: string; secret?: string } =>
        Boolean(a) && typeof a === 'object' && (a as { type?: string }).type === 'webhook')
      .map((a) => ({
        type: 'webhook' as const,
        url: String(a.url ?? ''),
        secret: String(a.secret ?? ''),
      }));
  } catch {
    return [];
  }
};

/**
 * 自动化 Flow 规则管理页：内容事件 → 条件（状态等值）→ 动作链。
 * v1 动作支持 webhook（HMAC 签名投递，与 webhook 订阅同款签名头）；
 * 命中次数为观测字段。
 */
const FlowRuleManagement = () => {
  const { t } = useTranslation('flow-rule');
  const { message } = App.useApp();
  const actionRef = useRef<ActionType>(null);

  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<FlowRule | undefined>();
  const [form, setForm] = useState<{
    name: string;
    triggerEvent: string;
    conditionStatus: string;
    webhooks: WebhookAction[];
    enabled: boolean;
  }>({ name: '', triggerEvent: '', conditionStatus: '', webhooks: [], enabled: true });

  const listQuery = useListFlowRules();
  const createMutation = useCreateFlowRule();
  const updateMutation = useUpdateFlowRule();
  const deleteMutation = useDeleteFlowRule();

  const dataSource = (listQuery.data?.items ?? []) as FlowRule[];

  const reload = () => listQuery.refetch();

  const openCreate = () => {
    setEditing(undefined);
    setForm({
      name: '',
      triggerEvent: 'post.published',
      conditionStatus: '',
      webhooks: [{ type: 'webhook', url: '', secret: '' }],
      enabled: true,
    });
    setFormOpen(true);
  };

  const openEdit = (record: FlowRule) => {
    setEditing(record);
    setForm({
      name: record.name ?? '',
      triggerEvent: record.triggerEvent ?? 'post.published',
      conditionStatus: record.conditionStatus ?? '',
      webhooks:
        deserializeActions(record.actions ?? '').length > 0
          ? deserializeActions(record.actions ?? '')
          : [{ type: 'webhook', url: '', secret: '' }],
      enabled: record.enabled ?? true,
    });
    setFormOpen(true);
  };

  const handleSave = async () => {
    if (!form.name.trim()) {
      message.warning(t('nameRequired'));
      return;
    }
    const webhooks = form.webhooks.filter((w) => w.url.trim());
    if (webhooks.length === 0) {
      message.warning(t('webhookRequired'));
      return;
    }
    const payload = {
      name: form.name.trim(),
      triggerEvent: form.triggerEvent,
      conditionStatus: form.conditionStatus || undefined,
      actions: serializeActions(webhooks),
      enabled: form.enabled,
    };
    try {
      if (editing?.id) {
        await updateMutation.mutateAsync({
          id: editing.id,
          data: payload,
          updateMask: 'name,triggerEvent,conditionStatus,actions,enabled',
        });
      } else {
        await createMutation.mutateAsync({ data: payload });
      }
      message.success(t('saveSuccess'));
      setFormOpen(false);
      await reload();
    } catch (error) {
      console.error('Save flow rule failed:', error);
    }
  };

  const handleDelete = async (record: FlowRule) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
      message.success(t('deleteSuccess'));
      await reload();
    } catch (error) {
      console.error('Delete flow rule failed:', error);
    }
  };

  const handleToggle = async (record: FlowRule, checked: boolean) => {
    if (!record.id) return;
    try {
      await updateMutation.mutateAsync({
        id: record.id,
        data: { enabled: checked },
        updateMask: 'enabled',
      });
      await reload();
    } catch (error) {
      console.error('Toggle flow rule failed:', error);
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: t('name'), dataIndex: 'name', ellipsis: true },
    {
      title: t('triggerEvent'),
      dataIndex: 'triggerEvent',
      width: 170,
      render: (_: unknown, record: FlowRule) => (
        <Tag color="blue">{record.triggerEvent}</Tag>
      ),
    },
    {
      title: t('conditionStatus'),
      dataIndex: 'conditionStatus',
      width: 200,
      render: (_: unknown, record: FlowRule) =>
        record.conditionStatus ? (
          <Tag color="orange">{record.conditionStatus}</Tag>
        ) : (
          <Tag>{t('noCondition')}</Tag>
        ),
    },
    {
      title: t('webhooks'),
      dataIndex: 'actions',
      ellipsis: true,
      render: (_: unknown, record: FlowRule) => {
        const webhooks = deserializeActions(record.actions ?? '');
        if (webhooks.length === 0) return '-';
        return (
          <Space wrap size={4}>
            {webhooks.map((w, i) => (
              <Tag key={i} style={{ maxWidth: 220, overflow: 'hidden', textOverflow: 'ellipsis' }}>
                {w.url}
              </Tag>
            ))}
          </Space>
        );
      },
    },
    {
      title: t('matchCount'),
      dataIndex: 'matchCount',
      width: 90,
      render: (_: unknown, record: FlowRule) => Number(record.matchCount || 0),
    },
    {
      title: t('enabled'),
      dataIndex: 'enabled',
      width: 90,
      render: (_: unknown, record: FlowRule) => (
        <Switch
          checked={record.enabled ?? false}
          onChange={(checked) => handleToggle(record, checked)}
        />
      ),
    },
    {
      title: t('action'),
      width: 150,
      render: (_: unknown, record: FlowRule) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
            {t('edit')}
          </Button>
          <Popconfirm title={t('deleteConfirm')} onConfirm={() => handleDelete(record)}>
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              {t('delete')}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <ContentContainer>
      <ListTable
        headerTitle={t('title')}
        actionRef={actionRef}
        search={false}
        loading={listQuery.isPending}
        rowKey="id"
        columns={columns}
        request={async () => {
          const items = dataSource;
          return { data: items, total: Number(listQuery.data?.total) || items.length, success: true };
        }}
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
        width={620}
        destroyOnHidden
      >
        <Space direction="vertical" style={{ width: '100%' }} size={12}>
          <Input
            placeholder={t('namePlaceholder')}
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
          />
          <Select
            placeholder={t('triggerEventPlaceholder')}
            value={form.triggerEvent}
            options={FLOW_EVENT_OPTIONS}
            onChange={(v: string) => setForm((f) => ({ ...f, triggerEvent: v }))}
            style={{ width: '100%' }}
          />
          <Select
            allowClear
            placeholder={t('conditionStatusPlaceholder')}
            value={form.conditionStatus || undefined}
            options={FLOW_STATUS_CONDITION_OPTIONS}
            onChange={(v?: string) => setForm((f) => ({ ...f, conditionStatus: v ?? '' }))}
            style={{ width: '100%' }}
          />
          {form.webhooks.map((w, i) => (
            <Space key={i} direction="vertical" style={{ width: '100%' }} size={6}>
              <Space.Compact style={{ width: '100%' }}>
                <Input
                  placeholder={t('webhookUrlPlaceholder')}
                  value={w.url}
                  onChange={(e) =>
                    setForm((f) => ({
                      ...f,
                      webhooks: f.webhooks.map((x, j) =>
                        j === i ? { ...x, url: e.target.value } : x,
                      ),
                    }))
                  }
                />
                <Button
                  icon={<DeleteOutlined />}
                  disabled={form.webhooks.length <= 1}
                  onClick={() =>
                    setForm((f) => ({
                      ...f,
                      webhooks: f.webhooks.filter((_, j) => j !== i),
                    }))
                  }
                />
              </Space.Compact>
              <Input.Password
                placeholder={t('webhookSecretPlaceholder')}
                value={w.secret}
                onChange={(e) =>
                  setForm((f) => ({
                    ...f,
                    webhooks: f.webhooks.map((x, j) =>
                      j === i ? { ...x, secret: e.target.value } : x,
                    ),
                  }))
                }
              />
            </Space>
          ))}
          <Button
            block
            type="dashed"
            icon={<PlusOutlined />}
            onClick={() =>
              setForm((f) => ({
                ...f,
                webhooks: [...f.webhooks, { type: 'webhook', url: '', secret: '' }],
              }))
            }
          >
            {t('addWebhook')}
          </Button>
          <Space>
            <span>{t('enabled')}</span>
            <Switch
              checked={form.enabled}
              onChange={(v) => setForm((f) => ({ ...f, enabled: v }))}
            />
          </Space>
        </Space>
      </Modal>
    </ContentContainer>
  );
};

export default FlowRuleManagement;
