import { useEffect, useState } from 'react';
import { App, Button, Card, Input, Select, Space, Switch } from 'antd';
import { SaveOutlined, SendOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import {
  useCreateEmailProvider,
  useEmailProvidersList,
  useTestEmailProviderSend,
  useUpdateEmailProvider,
} from '@/api/hooks/email-provider';

/**
 * 邮件出站 SMTP 配置页（平台基础设施）：
 * 单配置卡形态（邮件链取第一个启用项）；凭据 AES-GCM 加密落库、接口永不回显。
 * notify.email worker 读库优先、SMTP_* 环境变量回退。
 */

const TLS_MODES = [
  { value: 'starttls', label: 'STARTTLS (587)' },
  { value: 'ssl', label: 'SSL/TLS (465)' },
  { value: 'none', label: 'None（仅内网调试）' },
];

interface ProviderForm {
  name: string;
  host: string;
  port: number;
  username: string;
  password: string;
  from: string;
  tlsMode: string;
  enabled: boolean;
  description: string;
}

const EMPTY: ProviderForm = {
  name: '主 SMTP',
  host: '',
  port: 587,
  username: '',
  password: '',
  from: '',
  tlsMode: 'starttls',
  enabled: true,
  description: '',
};

const EmailProviderManagement = () => {
  const { t } = useTranslation('email-provider');
  const { message } = App.useApp();

  const [form, setForm] = useState<ProviderForm>(EMPTY);
  const [existingId, setExistingId] = useState<number | undefined>();
  const [testTo, setTestTo] = useState('');
  const [saving, setSaving] = useState(false);

  const listMutation = useEmailProvidersList();
  const createMutation = useCreateEmailProvider();
  const updateMutation = useUpdateEmailProvider();
  const testSendMutation = useTestEmailProviderSend();

  useEffect(() => {
    listMutation
      .mutateAsync()
      .then((resp) => {
        const first = resp?.items?.[0];
        if (first?.id) {
          setExistingId(first.id);
          setForm({
            name: first.name ?? '',
            host: first.host ?? '',
            port: first.port ?? 587,
            username: first.username ?? '',
            password: '',
            from: first.from ?? '',
            tlsMode: first.tlsMode ?? 'starttls',
            enabled: first.enabled ?? true,
            description: first.description ?? '',
          });
        }
      })
      .catch((error) => {
        console.error('Load email providers failed:', error);
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleSave = async () => {
    if (!form.host.trim() || !form.from.trim()) {
      message.warning(t('requiredFields'));
      return;
    }
    setSaving(true);
    try {
      if (existingId) {
        await updateMutation.mutateAsync({
          id: existingId,
          data: { ...form },
          // 密码仅在填写时随 mask 提交（留空=不修改）
          updateMask: `name,host,port,username,from,tlsMode,enabled,description${form.password ? ',password' : ''}`,
        });
      } else {
        const created = await createMutation.mutateAsync({ data: { ...form } });
        setExistingId(created?.id);
      }
      message.success(t('saveSuccess'));
    } catch (error) {
      console.error('Save email provider failed:', error);
    } finally {
      setSaving(false);
    }
  };

  const handleTestSend = async () => {
    if (!existingId) {
      message.warning(t('saveBeforeTest'));
      return;
    }
    if (!testTo.trim()) {
      message.warning(t('testToRequired'));
      return;
    }
    try {
      const resp = await testSendMutation.mutateAsync({
        id: existingId,
        to: testTo.trim(),
      });
      if (resp?.success) {
        message.success(t('testSuccess'));
      } else {
        message.error(`${t('testFailed')}：${resp?.error ?? '-'}`);
      }
    } catch (error) {
      console.error('Test send failed:', error);
    }
  };

  return (
    <ContentContainer>
      <Card title={t('title')} style={{ maxWidth: 720 }}>
        <Space direction="vertical" style={{ width: '100%' }} size={12}>
          <Input
            placeholder={t('namePlaceholder')}
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
          />
          <Space.Compact style={{ width: '100%' }}>
            <Input
              style={{ width: '70%' }}
              placeholder={t('hostPlaceholder')}
              value={form.host}
              onChange={(e) => setForm((f) => ({ ...f, host: e.target.value }))}
            />
            <Input
              style={{ width: '30%' }}
              type="number"
              placeholder="587"
              value={form.port}
              onChange={(e) => setForm((f) => ({ ...f, port: Number(e.target.value) || 587 }))}
            />
          </Space.Compact>
          <Input
            placeholder={t('fromPlaceholder')}
            value={form.from}
            onChange={(e) => setForm((f) => ({ ...f, from: e.target.value }))}
          />
          <Input
            placeholder={t('usernamePlaceholder')}
            value={form.username}
            onChange={(e) => setForm((f) => ({ ...f, username: e.target.value }))}
          />
          <Input.Password
            placeholder={existingId ? t('passwordKeepPlaceholder') : t('passwordPlaceholder')}
            value={form.password}
            onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
          />
          <Select
            value={form.tlsMode}
            options={TLS_MODES}
            onChange={(v: string) => setForm((f) => ({ ...f, tlsMode: v }))}
            style={{ width: 240 }}
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
          <Space wrap>
            <Button type="primary" icon={<SaveOutlined />} loading={saving} onClick={handleSave}>
              {t('save')}
            </Button>
            <Input
              style={{ width: 240 }}
              placeholder={t('testToPlaceholder')}
              value={testTo}
              onChange={(e) => setTestTo(e.target.value)}
            />
            <Button icon={<SendOutlined />} loading={testSendMutation.isPending} onClick={handleTestSend}>
              {t('testSend')}
            </Button>
          </Space>
          <div style={{ color: 'var(--ant-color-text-tertiary)', fontSize: 12 }}>
            {t('hint')}
          </div>
        </Space>
      </Card>
    </ContentContainer>
  );
};

export default EmailProviderManagement;
