import { useMemo, useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import ListTable from '@/components/common/ListTable';
import {
  ModalForm,
  DrawerForm,
  ProFormText,
  ProFormDigit,
  ProFormSelect,
  ProFormTextArea,
  ProFormSwitch,
  ProFormDependency,
} from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App } from 'antd';
import type { ProFormItemProps } from '@ant-design/pro-components';
import { EditOutlined, DeleteOutlined, PlusOutlined, SendOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type {
  notifyservicev1_NotificationChannel as NotificationChannel,
} from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import { TABLE } from '@/config/constants';
import {
  fetchListNotificationChannels,
  useDeleteNotificationChannel,
  useSendTestEmail,
  useUpdateNotificationChannel,
  useCreateNotificationChannel,
} from '@/api/hooks/notification-channel';
import { useListChannelSchemas, type ChannelSchema } from '@/api/hooks/notification-channel-schema';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

/**
 * 渠道类型 → 文案/颜色。列必须按 record.type 渲染：一个渠道类型对应一列自己的配置
 *（SMTP 那几列对 WEBHOOK 没有意义，反之 webhookUrl 对邮件渠道为空），混着显示会读成"没配"。
 */
const CHANNEL_TYPE_LABEL_KEY: Partial<Record<NonNullable<NotificationChannel['type']>, string>> = {
  EMAIL: 'typeEmail',
  WEBHOOK: 'typeWebhook',
};

const CHANNEL_TYPE_COLOR: Partial<Record<NonNullable<NotificationChannel['type']>, string>> = {
  EMAIL: 'blue',
  WEBHOOK: 'purple',
};

// EXT_* 扩展渠道的展示名/分组来自后端 schema（settings 列那批），不在 i18n 里逐个维护。
// 判定"是不是扩展渠道"的唯一真理源 = schema 清单（hook 数据），组件内用 schemaByType。

/** settings 读回显 JSON 容错解析：坏数据返回 undefined，动态表单按全空处理。 */
function safeParseSettings(raw: string): Record<string, unknown> | undefined {
  try {
    const parsed = JSON.parse(raw);
    return typeof parsed === 'object' && parsed !== null ? (parsed as Record<string, unknown>) : undefined;
  } catch {
    // 存储态坏了按全空处理，不阻塞表单渲染
    return undefined;
  }
}

/** EXT_* 扩展渠道的 Tag 统一色。 */
const EXT_TAG_COLOR = 'geekblue';

/**
 * 签名风格 → 文案。列值就是后端 ent 列里的同一组字符串（枚举按名字配对），
 * 所以这里既不换算也不加前缀。
 */
const SIGN_STYLE_LABEL_KEY: Record<
  NonNullable<NotificationChannel['webhookSignStyle']>,
  string
> = {
  CUSTOM: 'signStyleCustom',
  NONE: 'signStyleNone',
  DINGTALK: 'signStyleDingtalk',
  FEISHU: 'signStyleFeishu',
  WECOM: 'signStyleWecom',
};

// 这一列可空，读到空值时按 CUSTOM 显示——它的实际行为就是 CUSTOM。
// EMAIL 行没有这一列，返回 undefined 由列渲染 '-'。
const signStyleKeyOf = (record: NotificationChannel): string | undefined =>
  record.type === 'WEBHOOK'
    ? SIGN_STYLE_LABEL_KEY[record.webhookSignStyle ?? 'CUSTOM']
    : undefined;

/**
 * 载荷模板可引用的占位符清单，作为词条的 `vars` 插值值传进去，而不是把花括号写进词条：
 * 一是这些名字属于接口契约（各语言都不该翻译），二是 i18next 会把词条里的 `{{x}}`
 * 当成插值语法吃掉，写进去只会渲染成一串空白。
 */
const WEBHOOK_TEMPLATE_VARS = [
  'title',
  'content',
  'event_type',
  'timestamp',
  'sign',
  'nonce',
  'recipient_user_id',
  'related_id',
  'delivered_at',
]
  .map((name) => `{{${name}}}`)
  .join(' ');

/**
 * 通知渠道管理页面（平台级配置）
 * 两类渠道：EMAIL（SMTP 账号）与 WEBHOOK（回调地址 + 可选签名密钥）。
 */
const NotificationChannelManagement = () => {
  const { t } = useTranslation('notification-channel');
  const actionRef = useRef<ActionType>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  // 扩展渠道 schema：展示名/分组/动态表单的唯一数据源
  const schemasQuery = useListChannelSchemas();
  const schemas = useMemo(() => schemasQuery.data?.items ?? [], [schemasQuery.data]);
  const schemaByType = useMemo(
    () => new Map<string, ChannelSchema>(schemas.map((s) => [String(s.type), s])),
    [schemas],
  );
  const isExtType = (t?: string) => !!t && schemaByType.has(t);
  // 类型下拉按 schema 分组（EXT 部分），EMAIL/WEBHOOK 固定在最前
  const extTypeOptions = useMemo(() => {
    const groups = new Map<string, { label: string; value: string }[]>();
    for (const s of schemas) {
      const g = String(s.group || 'EXT');
      const bucket = groups.get(g) ?? [];
      bucket.push({ label: String(s.label), value: String(s.type) });
      groups.set(g, bucket);
    }
    return [...groups.entries()].map(([group, options]) => ({ group, label: t(`extGroup.${group}`), options }));
  }, [schemas, t]);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [selected, setSelected] = useState<NotificationChannel | undefined>();
  const [testTarget, setTestTarget] = useState<NotificationChannel | undefined>();

  const refresh = () => {
    actionRef.current?.reload();
    queryClient.invalidateQueries({ queryKey: ['listNotificationChannels'] });
  };

  const deleteMutation = useDeleteNotificationChannel({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      refresh();
    },
  });

  const createMutation = useCreateNotificationChannel();
  const updateMutation = useUpdateNotificationChannel();
  const testMutation = useSendTestEmail();

  const handleDelete = async (record: NotificationChannel) => {
    if (!record.id) return;
    try {
      await deleteMutation.mutateAsync({ id: record.id });
    } catch (error) {
      // 原始错误必须留在控制台：用户可见的那句翻译不包含服务端的原因
      console.error('delete notification channel failed', error);
      message.error(errMessage(error, t('deleteFailed')));
    }
  };

  const handleSubmit = async (values: Record<string, unknown>) => {
    // 密钥与 settings 是请求级字段（不进 data）：读视图只有 hasPassword / hasWebhookSecret 布尔
    const { password, webhookSecret, settings, ...data } = values as {
      password?: string;
      webhookSecret?: string;
      settings?: Record<string, unknown>;
    } & Record<string, unknown>;
    // 扩展渠道：动态表单的 settings 对象序列化为 JSON 明文上送（sensitive 留空 = 保留已存值）
    if (isExtType(String(data.type))) {
      data.settings = settings ? JSON.stringify(settings) : '';
    } else {
      delete data.settings;
    }
    try {
      if (drawerMode === 'create') {
        await createMutation.mutateAsync({
          data: data as unknown as NotificationChannel,
          password: password || undefined,
          webhookSecret: webhookSecret || undefined,
        });
        message.success(t('createSuccess'));
      } else if (selected?.id) {
        await updateMutation.mutateAsync({
          id: selected.id,
          data: data as unknown as NotificationChannel,
          password: password || undefined,
          webhookSecret: webhookSecret || undefined,
          updateMask: isExtType(String(data.type))
            ? 'name,type,settings,enabled,remark'
            : 'name,type,smtpHost,smtpPort,smtpUsername,smtpFrom,smtpFromName,smtpTls,webhookUrl,webhookSignStyle,webhookPayloadTemplate,enabled,remark',
        });
        message.success(t('updateSuccess'));
      }
      setDrawerOpen(false);
      refresh();
      return true;
    } catch (error) {
      console.error('save notification channel failed', error);
      message.error(errMessage(error, t('saveFailed')));
      return false;
    }
  };

  const handleTestSend = async (values: { recipient: string }) => {
    if (!testTarget?.id) return false;
    try {
      const resp = await testMutation.mutateAsync({
        id: testTarget.id,
        recipient: values.recipient,
      });
      if (resp.success) {
        message.success(t('testSendSuccess'));
        setTestTarget(undefined);
      } else {
        message.error(resp.error || t('testSendFailed'));
      }
      return resp.success;
    } catch (error) {
      console.error('send test email failed', error);
      message.error(errMessage(error, t('testSendFailed')));
      return false;
    }
  };

  const columns: ProColumns<NotificationChannel>[] = [
    {
      title: t('name'),
      dataIndex: 'name',
      width: 150,
      ellipsis: true,
    },
    {
      title: t('type'),
      dataIndex: 'type',
      width: 100,
      render: (_, record) =>
        record.type ? (
          <Tag color={CHANNEL_TYPE_COLOR[record.type] ?? EXT_TAG_COLOR}>
            {schemaByType.get(record.type)?.label ?? t(CHANNEL_TYPE_LABEL_KEY[record.type] ?? record.type)}
          </Tag>
        ) : (
          '-'
        ),
    },
    {
      title: t('smtpHost'),
      dataIndex: 'smtpHost',
      width: 160,
      ellipsis: true,
      render: (_, record) => record.smtpHost || '-',
    },
    {
      title: t('smtpPort'),
      dataIndex: 'smtpPort',
      width: 90,
      render: (_, record) => record.smtpPort ?? '-',
    },
    {
      title: t('smtpFrom'),
      dataIndex: 'smtpFrom',
      width: 180,
      ellipsis: true,
      render: (_, record) => record.smtpFrom || '-',
    },
    {
      title: t('smtpTls'),
      dataIndex: 'smtpTls',
      width: 110,
      render: (_, record) => {
        const map: Record<string, string> = {
          NONE: t('tlsNone'),
          START_TLS: t('tlsStartTls'),
          SSL: t('tlsSsl'),
        };
        return map[record.smtpTls || ''] || record.smtpTls || '-';
      },
    },
    {
      title: t('webhookUrl'),
      dataIndex: 'webhookUrl',
      width: 200,
      ellipsis: true,
      render: (_, record) => record.webhookUrl || '-',
    },
    {
      title: t('webhookSignStyle'),
      dataIndex: 'webhookSignStyle',
      width: 140,
      render: (_, record) => {
        const key = signStyleKeyOf(record);
        return key ? <Tag>{t(key)}</Tag> : '-';
      },
    },
    {
      title: t('hasPassword'),
      dataIndex: 'hasPassword',
      width: 100,
      render: (_, record) =>
        record.hasPassword ? (
          <Tag color="green">{t('passwordSet')}</Tag>
        ) : (
          <Tag>{t('passwordNotSet')}</Tag>
        ),
    },
    {
      title: t('hasWebhookSecret'),
      dataIndex: 'hasWebhookSecret',
      width: 100,
      render: (_, record) =>
        record.hasWebhookSecret ? (
          <Tag color="green">{t('passwordSet')}</Tag>
        ) : (
          <Tag>{t('passwordNotSet')}</Tag>
        ),
    },
    {
      title: t('enabled'),
      dataIndex: 'enabled',
      width: 90,
      render: (_, record) =>
        record.enabled ? <Tag color="success">{t('enabledOn')}</Tag> : <Tag>{t('enabledOff')}</Tag>,
    },
    {
      title: t('remark'),
      dataIndex: 'remark',
      width: 140,
      ellipsis: true,
      hideInTable: window.innerWidth < 1600,
      render: (_, record) => record.remark || '-',
    },
    {
      title: t('actions'),
      valueType: 'option',
      width: 220,
      fixed: 'right',
      render: (_, record) => [
        <Button
          key="edit"
          type="link"
          size="small"
          icon={<EditOutlined />}
          onClick={() => {
            setDrawerMode('edit');
            setSelected(record);
            setDrawerOpen(true);
          }}
        >
          {t('edit')}
        </Button>,
        // 测试发送对 EMAIL 与扩展渠道有意义（验证配置本身）；WEBHOOK 的端到端验证在路由规则页，
        // 那里测的是"事件 → 渠道 → 台账"整条链，渠道页只能证明配置本身。
        (record.type === 'EMAIL' || isExtType(record.type)) && (
          <Button key="test" type="link" size="small" icon={<SendOutlined />} onClick={() => setTestTarget(record)}>
            {t('testSend')}
          </Button>
        ),
        <Popconfirm key="delete" title={t('deleteConfirm')} onConfirm={() => handleDelete(record)}>
          <Button danger type="link" size="small" icon={<DeleteOutlined />}>
            {t('delete')}
          </Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <div ref={containerRef} className="page-container-content">
        <ListTable<NotificationChannel>
          actionRef={actionRef}
          columns={columns}
          request={async (params) => {
            const query = new PaginationQuery({
              paging: { page: params.current || 1, pageSize: params.pageSize || 20 },
              formValues: Object.fromEntries(
                Object.entries(params).filter(([key]) => !['current', 'pageSize'].includes(key)),
              ),
            });
            const response = await fetchListNotificationChannels(query);
            return { data: response.items || [], total: Number(response.total) || 0, success: true };
          }}
          rowKey="id"
          search={false}
          pagination={{
            defaultPageSize: TABLE.DEFAULT_PAGE_SIZE,
            showSizeChanger: true,
          }}
          options={{ density: true, fullScreen: true, setting: true, reload: true }}
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
          scroll={{ y: tableScrollY, x: 1720 }}
        />
      </div>

      {/* 新建/编辑用抽屉不用弹窗：SMTP 配置字段多（7 个），
          弹窗高度挤不下；抽屉纵向空间足。 */}
      <DrawerForm
        title={drawerMode === 'create' ? t('create') : t('edit')}
        width={560}
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        drawerProps={{ destroyOnHidden: true, maskClosable: false }}
        submitTimeout={3000}
        onFinish={handleSubmit}
        initialValues={
          drawerMode === 'create'
            ? { type: 'EMAIL', smtpTls: 'START_TLS', smtpPort: 587, webhookSignStyle: 'CUSTOM', enabled: true }
            : {
                ...selected,
                password: undefined,
                webhookSecret: undefined,
                // 同列表列：该列为空时按 CUSTOM 摆进表单，否则这一项是空的、改不动
                webhookSignStyle: selected?.webhookSignStyle ?? 'CUSTOM',
                // settings 读回显是脱敏 JSON 串：parse 成对象供动态表单回填
                //（sensitive 字段值为空串=占位，留空提交即保留已存密文）
                settings: selected?.settings ? safeParseSettings(selected.settings) : undefined,
              }
        }
      >
        <ProFormText name="name" label={t('name')} rules={[{ required: true, message: t('requiredName') }]} />
        <ProFormSelect
          name="type"
          label={t('type')}
          options={[
            { label: t('typeEmail'), value: 'EMAIL' },
            { label: t('typeWebhook'), value: 'WEBHOOK' },
            // 扩展渠道按 schema 分组（IM / PUSH / SMS / MAIL / AUTO / OPS）
            ...extTypeOptions.map((g) => ({ label: g.label, title: g.label, options: g.options })),
          ]}
          disabled={drawerMode === 'edit'}
          showSearch
          rules={[{ required: true }]}
        />
        {/* 两类渠道各自的配置项分组显示：表单里同时摆几组，管理员要自己猜哪一组生效。
            编辑态 type 不可改（改类型等于换一条渠道实现），所以几组不会在同一个会话里来回切。 */}
        <ProFormDependency name={['type']}>
          {({ type }) =>
            isExtType(String(type)) ? (
              // 扩展渠道：settings 动态表单（schema 声明字段，sensitive 渲染为密码框）
              <SettingsFormFields schema={schemaByType.get(String(type))} />
            ) : type === 'WEBHOOK' ? (
              <>
                <ProFormText
                  name="webhookUrl"
                  label={t('webhookUrl')}
                  placeholder="https://example.com/hooks/notification"
                  tooltip={t('webhookUrlHint')}
                  rules={drawerMode === 'create' ? [{ required: true, message: t('requiredWebhookUrl') }] : []}
                />
                <ProFormSelect
                  name="webhookSignStyle"
                  label={t('webhookSignStyle')}
                  tooltip={t('webhookSignStyleHint')}
                  options={(['CUSTOM', 'NONE', 'DINGTALK', 'FEISHU', 'WECOM'] as const).map((value) => ({
                    value,
                    label: t(SIGN_STYLE_LABEL_KEY[value]),
                  }))}
                  rules={[{ required: true, message: t('requiredSignStyle') }]}
                />
                {/* 留空 = 用该风格的内置默认正文（内置形状由后端持有，
                    这里不复制一份，否则两处会漂）。 */}
                <ProFormTextArea
                  name="webhookPayloadTemplate"
                  label={t('webhookPayloadTemplate')}
                  tooltip={t('webhookPayloadTemplateHint', { vars: WEBHOOK_TEMPLATE_VARS })}
                  placeholder='{"msgtype":"text","text":{"content":"{{title}}"}}'
                  fieldProps={{ rows: 3, style: { fontFamily: 'var(--font-mono, monospace)' } }}
                />
                <ProFormText.Password
                  name="webhookSecret"
                  label={t('webhookSecret')}
                  placeholder={drawerMode === 'edit' ? t('webhookSecretKeepHint') : t('webhookSecretPlaceholder')}
                  tooltip={t('webhookSecretHint')}
                />
              </>
            ) : (
              <>
                <ProFormText name="smtpHost" label={t('smtpHost')} placeholder="smtp.example.com" />
                <ProFormDigit
                  name="smtpPort"
                  label={t('smtpPort')}
                  min={1}
                  max={65535}
                  fieldProps={{ precision: 0 }}
                />
                <ProFormText name="smtpUsername" label={t('smtpUsername')} />
                <ProFormText.Password
                  name="password"
                  label={t('password')}
                  placeholder={drawerMode === 'edit' ? t('passwordKeepHint') : t('passwordPlaceholder')}
                  rules={drawerMode === 'create' ? [{ required: true, message: t('requiredPassword') }] : []}
                />
                <ProFormText name="smtpFrom" label={t('smtpFrom')} placeholder="noreply@example.com" />
                <ProFormText name="smtpFromName" label={t('smtpFromName')} placeholder="GoWind" />
                <ProFormSelect
                  name="smtpTls"
                  label={t('smtpTls')}
                  options={[
                    { label: t('tlsNone'), value: 'NONE' },
                    { label: t('tlsStartTls'), value: 'START_TLS' },
                    { label: t('tlsSsl'), value: 'SSL' },
                  ]}
                />
              </>
            )
          }
        </ProFormDependency>
        <ProFormSwitch name="enabled" label={t('enabled')} />
        <ProFormTextArea name="remark" label={t('remark')} fieldProps={{ rows: 2 }} />
      </DrawerForm>

      <ModalForm<{ recipient: string }>
        title={t('testSendTitle', { name: testTarget?.name || '' })}
        width={440}
        open={!!testTarget}
        onOpenChange={(open) => {
          if (!open) setTestTarget(undefined);
        }}
        modalProps={{ destroyOnHidden: true }}
        submitTimeout={5000}
        onFinish={handleTestSend}
      >
        <ProFormText
          name="recipient"
          label={t('testRecipient')}
          rules={[
            { required: true, message: t('requiredRecipient') },
            // 邮箱格式只对 EMAIL 渠道校验；扩展渠道的 recipient 语义是运行目标
            //（chat_id/topic/设备 Key…），只做非空不做格式——格式错了由服务商报错进台账。
            ...(testTarget?.type === 'EMAIL' ? [{ type: 'email' as const, message: t('invalidEmail') }] : []),
          ]}
          placeholder={testTarget?.type === 'EMAIL' ? 'you@example.com' : undefined}
        />
      </ModalForm>
    </ContentContainer>
  );
};

export default NotificationChannelManagement;

/** schema 字段类型 → 动态表单（settings 命名空间：提交时 values.settings = {name: value}）。 */
function SettingsFormFields({ schema }: { schema: ChannelSchema | undefined }) {
  const { t } = useTranslation('notification-channel');
  if (!schema) return null;
  return (
    <>
      {schema.targetLabel && (
        <div style={{ color: 'var(--ant-color-text-tertiary, #aaa)', fontSize: 12, marginBottom: 8 }}>
          {t('ext.targetLabel', { label: schema.targetLabel })}
        </div>
      )}
      {(schema.fields ?? []).map((f) => {
        const name = ['settings', String(f.name)];
        const label = String(f.label ?? f.name);
        const rules = f.required ? [{ required: true, message: t('ext.required', { label }) }] : [];
        const common = { name, label, rules, tooltip: f.hint || undefined, extra: undefined };
        switch (f.type) {
          case 'textarea':
            return <ProFormTextArea key={f.name} {...common} fieldProps={{ rows: 3, placeholder: f.placeholder }} />;
          case 'number':
            return <ProFormDigit key={f.name} {...common} fieldProps={{ precision: 0, placeholder: f.placeholder }} />;
          case 'password':
            return (
              <ProFormText.Password
                key={f.name}
                {...(common as ProFormItemProps)}
                placeholder={f.placeholder || t('ext.secretKeepHint')}
              />
            );
          case 'select':
            return (
              <ProFormSelect
                key={f.name}
                {...common}
                options={(f.options ?? []).map((o) => ({ label: o === '' ? t('ext.optionDefault') : o, value: o }))}
                fieldProps={{ allowClear: true }}
              />
            );
          default:
            return <ProFormText key={f.name} {...common} placeholder={f.placeholder} />;
        }
      })}
    </>
  );
}
