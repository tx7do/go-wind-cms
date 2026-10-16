import { useRef, useState } from 'react';
import type { ActionType } from '@ant-design/pro-components';
import {
  App,
  Button,
  Input,
  Modal,
  Popconfirm,
  Radio,
  Select,
  Space,
  Switch,
  Tag,
} from 'antd';
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

/** 动作编辑模型：webhook（url+secret）/ notification（title/content/收件人）/ delay（seconds） */
interface FlowActionForm {
  type: 'webhook' | 'notification' | 'delay';
  url: string;
  secret: string;
  title: string;
  content: string;
  /** 收件人用户 ID 逗号分隔输入态 */
  targetUserIdsText: string;
  /** delay 动作秒数输入态 */
  secondsText: string;
}

const emptyAction = (type: FlowActionForm['type']): FlowActionForm => ({
  type,
  url: '',
  secret: '',
  title: '',
  content: '',
  targetUserIdsText: '',
  secondsText: '',
});

const serializeActions = (actions: FlowActionForm[]) =>
  JSON.stringify(
    actions
      .map((a) => {
        if (a.type === 'webhook') {
          if (!a.url.trim()) return null;
          const out: Record<string, unknown> = {
            type: 'webhook',
            url: a.url.trim(),
          };
          if (a.secret.trim()) out.secret = a.secret.trim();
          return out;
        }
        if (a.type === 'delay') {
          const seconds = Number(a.secondsText);
          if (!Number.isInteger(seconds) || seconds < 1) return null;
          return { type: 'delay', seconds };
        }
        const ids = a.targetUserIdsText
          .split(/[,，\s]+/)
          .map((x) => Number(x))
          .filter((x) => Number.isInteger(x) && x > 0);
        if (!a.title.trim() || !a.content.trim() || ids.length === 0) return null;
        return {
          type: 'notification',
          title: a.title.trim(),
          content: a.content.trim(),
          target_user_ids: ids,
        };
      })
      .filter(Boolean),
  );

const deserializeActions = (actionsJSON: string): FlowActionForm[] => {
  try {
    const raw: unknown = JSON.parse(actionsJSON || '[]');
    if (!Array.isArray(raw)) return [];
    return raw
      .filter(
        (a): a is Record<string, unknown> =>
          Boolean(a) && typeof a === 'object',
      )
      .map((a) => {
        const type =
          a.type === 'notification'
            ? 'notification'
            : a.type === 'delay'
              ? 'delay'
              : 'webhook';
        return {
          type,
          url: String(a.url ?? ''),
          secret: String(a.secret ?? ''),
          title: String(a.title ?? ''),
          content: String(a.content ?? ''),
          targetUserIdsText: Array.isArray(a.target_user_ids)
            ? (a.target_user_ids as unknown[]).map(String).join(',')
            : '',
          secondsText:
            a.type === 'delay' && Number.isFinite(Number(a.seconds))
              ? String(a.seconds)
              : '',
        };
      });
  } catch {
    return [];
  }
};

/** 动作链编辑器：主链与条件不满足分支共用（addLabel 撑三条添加按钮文案） */
const ActionEditor = ({
  actions,
  onChange,
}: {
  actions: FlowActionForm[];
  onChange: (next: FlowActionForm[]) => void;
}) => {
  const { t } = useTranslation('flow-rule');
  return (
    <>
      {actions.map((a, i) => {
        const patch = (part: Partial<FlowActionForm>) =>
          onChange(actions.map((x, j) => (j === i ? { ...x, ...part } : x)));
        return (
          <div
            key={i}
            style={{
              border: '1px solid var(--ant-color-border-secondary)',
              borderRadius: 8,
              padding: 12,
            }}
          >
            <Space
              style={{ width: '100%', justifyContent: 'space-between', marginBottom: 8 }}
            >
              <Select
                value={a.type}
                options={[
                  { value: 'webhook', label: t('actionWebhook') },
                  { value: 'notification', label: t('actionNotification') },
                  { value: 'delay', label: t('actionDelay') },
                ]}
                onChange={(v: FlowActionForm['type']) =>
                  onChange(actions.map((x, j) => (j === i ? { ...emptyAction(v) } : x)))
                }
                style={{ width: 160 }}
              />
              <Button
                icon={<DeleteOutlined />}
                disabled={actions.length <= 1}
                onClick={() => onChange(actions.filter((_, j) => j !== i))}
              />
            </Space>
            {a.type === 'webhook' ? (
              <Space direction="vertical" style={{ width: '100%' }} size={6}>
                <Input
                  placeholder={t('webhookUrlPlaceholder')}
                  value={a.url}
                  onChange={(e) => patch({ url: e.target.value })}
                />
                <Input.Password
                  placeholder={t('webhookSecretPlaceholder')}
                  value={a.secret}
                  onChange={(e) => patch({ secret: e.target.value })}
                />
              </Space>
            ) : a.type === 'delay' ? (
              <Input
                type="number"
                min={1}
                placeholder={t('delaySecondsPlaceholder')}
                value={a.secondsText}
                onChange={(e) => patch({ secondsText: e.target.value })}
              />
            ) : (
              <Space direction="vertical" style={{ width: '100%' }} size={6}>
                <Input
                  placeholder={t('notifyTitlePlaceholder')}
                  value={a.title}
                  onChange={(e) => patch({ title: e.target.value })}
                />
                <Input.TextArea
                  placeholder={t('notifyContentPlaceholder')}
                  value={a.content}
                  rows={3}
                  onChange={(e) => patch({ content: e.target.value })}
                />
                <Input
                  placeholder={t('notifyRecipientsPlaceholder')}
                  value={a.targetUserIdsText}
                  onChange={(e) => patch({ targetUserIdsText: e.target.value })}
                />
              </Space>
            )}
          </div>
        );
      })}
      <Space.Compact block>
        <Button
          block
          type="dashed"
          icon={<PlusOutlined />}
          onClick={() => onChange([...actions, emptyAction('webhook')])}
        >
          {t('addWebhookAction')}
        </Button>
        <Button
          block
          type="dashed"
          icon={<PlusOutlined />}
          onClick={() => onChange([...actions, emptyAction('notification')])}
        >
          {t('addNotificationAction')}
        </Button>
        <Button
          block
          type="dashed"
          icon={<PlusOutlined />}
          onClick={() => onChange([...actions, emptyAction('delay')])}
        >
          {t('addDelayAction')}
        </Button>
      </Space.Compact>
    </>
  );
};

/**
 * 自动化 Flow 规则管理页：内容事件/定时（cron）→ 条件（状态等值）→ 动作链。
 * 动作支持 webhook（HMAC 签名投递）、notification（站内信）与 delay
 * （延时后由后台任务接续剩余动作）；条件不满足分支（elseActions）可配置
 * 另一组动作；定时规则的触发事件/状态条件不生效。命中次数为观测字段。
 */
const FlowRuleManagement = () => {
  const { t } = useTranslation('flow-rule');
  const { message } = App.useApp();
  const actionRef = useRef<ActionType>(null);

  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<FlowRule | undefined>();
  const [form, setForm] = useState<{
    name: string;
    /** 触发方式：event=内容事件 / cron=定时（triggerCron 非空） */
    triggerType: 'event' | 'cron';
    triggerEvent: string;
    triggerCron: string;
    conditionStatus: string;
    actions: FlowActionForm[];
    /** 条件不满足分支：开启后编辑 elseActions（提交为空串即清除） */
    hasElse: boolean;
    elseActions: FlowActionForm[];
    enabled: boolean;
  }>({
    name: '',
    triggerType: 'event',
    triggerEvent: '',
    triggerCron: '',
    conditionStatus: '',
    actions: [emptyAction('webhook')],
    hasElse: false,
    elseActions: [emptyAction('webhook')],
    enabled: true,
  });

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
      triggerType: 'event',
      triggerEvent: 'post.published',
      triggerCron: '',
      conditionStatus: '',
      actions: [emptyAction('webhook')],
      hasElse: false,
      elseActions: [emptyAction('webhook')],
      enabled: true,
    });
    setFormOpen(true);
  };

  const openEdit = (record: FlowRule) => {
    setEditing(record);
    const triggerCron = record.triggerCron ?? '';
    const elseActions = deserializeActions(record.elseActions ?? '');
    setForm({
      name: record.name ?? '',
      triggerType: triggerCron ? 'cron' : 'event',
      triggerEvent: record.triggerEvent ?? 'post.published',
      triggerCron,
      conditionStatus: record.conditionStatus ?? '',
      actions:
        deserializeActions(record.actions ?? '').length > 0
          ? deserializeActions(record.actions ?? '')
          : [emptyAction('webhook')],
      hasElse: elseActions.length > 0,
      elseActions: elseActions.length > 0 ? elseActions : [emptyAction('webhook')],
      enabled: record.enabled ?? true,
    });
    setFormOpen(true);
  };

  const handleSave = async () => {
    if (!form.name.trim()) {
      message.warning(t('nameRequired'));
      return;
    }
    const isCron = form.triggerType === 'cron';
    if (isCron) {
      if (!form.triggerCron.trim()) {
        message.warning(t('cronRequired'));
        return;
      }
    }
    const actions = serializeActions(form.actions);
    if (!Array.isArray(JSON.parse(actions)) || (JSON.parse(actions) as unknown[]).length === 0) {
      message.warning(t('actionRequired'));
      return;
    }
    let elseActions = '';
    if (form.hasElse) {
      elseActions = serializeActions(form.elseActions);
      if ((JSON.parse(elseActions) as unknown[]).length === 0) {
        message.warning(t('elseActionRequired'));
        return;
      }
    }
    const payload = {
      name: form.name.trim(),
      triggerEvent: form.triggerEvent,
      triggerCron: isCron ? form.triggerCron.trim() : undefined,
      conditionStatus: form.conditionStatus || undefined,
      actions,
      // 空串=清除 else 分支（update_mask 含 elseActions 时按空串落库）
      elseActions: form.hasElse ? elseActions : '',
      enabled: form.enabled,
    };
    try {
      if (editing?.id) {
        await updateMutation.mutateAsync({
          id: editing.id,
          data: payload,
          updateMask:
            'name,triggerEvent,triggerCron,conditionStatus,actions,elseActions,enabled',
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
      title: t('trigger'),
      dataIndex: 'triggerEvent',
      width: 190,
      render: (_: unknown, record: FlowRule) =>
        record.triggerCron ? (
          <Tag color="purple">{`cron: ${record.triggerCron}`}</Tag>
        ) : (
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
      title: t('actions'),
      dataIndex: 'actions',
      ellipsis: true,
      render: (_: unknown, record: FlowRule) => {
        const actions = deserializeActions(record.actions ?? '');
        if (actions.length === 0 && !record.elseActions) return '-';
        return (
          <Space wrap size={4}>
            {actions.map((a, i) => (
              <Tag
                key={i}
                color={
                  a.type === 'webhook'
                    ? 'blue'
                    : a.type === 'delay'
                      ? 'purple'
                      : 'green'
                }
                style={{ maxWidth: 240, overflow: 'hidden', textOverflow: 'ellipsis' }}
              >
                {a.type === 'webhook'
                  ? `${t('actionWebhook')}: ${a.url}`
                  : a.type === 'delay'
                    ? `${t('actionDelay')}: ${a.secondsText}s`
                    : `${t('actionNotification')}: ${a.title}`}
              </Tag>
            ))}
            {record.elseActions ? <Tag color="red">{t('elseBranchTag')}</Tag> : null}
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
          <Radio.Group
            value={form.triggerType}
            onChange={(e) =>
              setForm((f) => ({ ...f, triggerType: e.target.value }))
            }
          >
            <Radio.Button value="event">{t('triggerTypeEvent')}</Radio.Button>
            <Radio.Button value="cron">{t('triggerTypeCron')}</Radio.Button>
          </Radio.Group>
          {form.triggerType === 'event' ? (
            <Select
              placeholder={t('triggerEventPlaceholder')}
              value={form.triggerEvent}
              options={FLOW_EVENT_OPTIONS}
              onChange={(v: string) => setForm((f) => ({ ...f, triggerEvent: v }))}
              style={{ width: '100%' }}
            />
          ) : (
            <Input
              placeholder={t('triggerCronPlaceholder')}
              value={form.triggerCron}
              onChange={(e) => setForm((f) => ({ ...f, triggerCron: e.target.value }))}
            />
          )}
          {form.triggerType === 'event' && (
            <Select
              allowClear
              placeholder={t('conditionStatusPlaceholder')}
              value={form.conditionStatus || undefined}
              options={FLOW_STATUS_CONDITION_OPTIONS}
              onChange={(v?: string) => setForm((f) => ({ ...f, conditionStatus: v ?? '' }))}
              style={{ width: '100%' }}
            />
          )}

          <div>
            <div style={{ marginBottom: 8 }}>
              <Tag color="green">{t('matchedBranchTitle')}</Tag>
            </div>
            <ActionEditor
              actions={form.actions}
              onChange={(actions) => setForm((f) => ({ ...f, actions }))}
            />
          </div>

          {form.triggerType === 'event' && (
            <div>
              <Space style={{ marginBottom: 8 }}>
                <Switch
                  size="small"
                  checked={form.hasElse}
                  onChange={(v) => setForm((f) => ({ ...f, hasElse: v }))}
                />
                <Tag color={form.hasElse ? 'red' : 'default'}>{t('elseBranchTitle')}</Tag>
              </Space>
              {form.hasElse && (
                <ActionEditor
                  actions={form.elseActions}
                  onChange={(elseActions) => setForm((f) => ({ ...f, elseActions }))}
                />
              )}
            </div>
          )}

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
