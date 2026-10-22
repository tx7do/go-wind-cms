import { useState } from 'react';
import { List, Button, Drawer, Form, Input, App } from 'antd';
import { MailOutlined, MobileOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQueryClient } from '@tanstack/react-query';
import type {
  identityservicev1_BindContactRequest,
  identityservicev1_VerifyContactRequest,
  identityservicev1_User as User,
} from '@/api/generated/admin/service/v1';
import {
  useGetUserProfile,
  useBindContact,
  useVerifyContact,
} from '@/api/hooks/user-profile';
import { useDrawerCloseGuard } from '@/components/common/use-drawer-close-guard';

const ICON_STYLE = { fontSize: 40 };

interface AccountBindProps {
  onSwitchTab: (key: string) => void;
}

const maskEmail = (email: string) => {
  const at = email.lastIndexOf('@');
  if (at <= 0) return email;
  return `${email.slice(0, Math.min(2, at))}***${email.slice(at)}`;
};

const maskPhone = (phone: string) => {
  if (phone.length < 7) return phone;
  return `${phone.slice(0, 3)}****${phone.slice(-4)}`;
};

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

const AccountBind = ({ onSwitchTab }: AccountBindProps) => {
  const { t } = useTranslation(['profile', 'common']);
  const { message } = App.useApp();
  const queryClient = useQueryClient();

  const [bindForm] = Form.useForm();
  const [bindOpen, setBindOpen] = useState(false);
  const [codeSent, setCodeSent] = useState(false);
  const [sending, setSending] = useState(false);
  const [verifying, setVerifying] = useState(false);

  const { data: userData } = useGetUserProfile();
  const user = userData as User | undefined | null;

  const closeDrawer = () => setBindOpen(false);
  const { guardClose, handleValuesChange } = useDrawerCloseGuard(bindOpen, closeDrawer);

  const bindMutation = useBindContact();
  const verifyMutation = useVerifyContact();

  const bindList = [
    {
      key: 'email',
      title: t('accountBind.email'),
      description: user?.email
        ? t('accountBind.bound', { value: maskEmail(user.email) })
        : t('accountBind.unbound'),
      icon: <MailOutlined style={ICON_STYLE} />,
    },
    {
      key: 'phone',
      title: t('accountBind.phone'),
      description: user?.mobile
        ? t('accountBind.bound', { value: maskPhone(user.mobile) })
        : t('accountBind.unbound'),
      icon: <MobileOutlined style={ICON_STYLE} />,
    },
  ];

  const openBind = () => {
    bindForm.setFieldsValue({ email: user?.email ?? '' });
    bindForm.resetFields(['code']);
    setCodeSent(false);
    setBindOpen(true);
  };

  const handleSendCode = async () => {
    const values = await bindForm.validateFields(['email']);
    const email = values.email?.trim();
    if (!email) {
      message.error(t('accountBind.newEmailRequired'));
      return;
    }
    setSending(true);
    try {
      const req: identityservicev1_BindContactRequest = { email: { email } };
      await bindMutation.mutateAsync(req);
      setCodeSent(true);
      message.success(t('accountBind.codeSent'));
    } catch (error) {
      console.error('Failed to send bind code:', error);
      message.error(errMessage(error, t('accountBind.vcodeSendFailed')));
    } finally {
      setSending(false);
    }
  };

  const handleVerify = async () => {
    const values = await bindForm.validateFields();
    const email = values.email?.trim();
    const code = values.code?.trim();
    if (!email || !code) {
      message.error(t('accountBind.vcodeRequired'));
      return;
    }
    setVerifying(true);
    try {
      const req: identityservicev1_VerifyContactRequest = { email: { email, code } };
      await verifyMutation.mutateAsync(req);
      message.success(t('accountBind.bindSuccess'));
      setBindOpen(false);
      // 刷新用户信息以更新绑定状态展示（本页读 getMe，顶栏头像等也读 getMe）
      queryClient.invalidateQueries({ queryKey: ['getMe'] });
    } catch (error) {
      console.error('Failed to verify bind code:', error);
      message.error(errMessage(error, t('accountBind.bindFailed')));
    } finally {
      setVerifying(false);
    }
  };

  return (
    <>
      <List
        dataSource={bindList}
        renderItem={(item) => (
          <List.Item
            actions={[
              item.key === 'email' && (
                <Button type="link" size="small" onClick={openBind}>
                  {user?.email ? t('accountBind.rebind') : t('accountBind.bindNow')}
                </Button>
              ),
              <Button type="link" size="small" onClick={() => onSwitchTab('1')}>
                {t('accountBind.modify')}
              </Button>,
            ].filter(Boolean)}
          >
            <List.Item.Meta
              avatar={item.icon}
              title={item.title}
              description={item.description}
            />
          </List.Item>
        )}
      />

      <Drawer
        title={t('accountBind.bindEmailTitle')}
        open={bindOpen}
        onClose={guardClose}
        width={480}
        destroyOnHidden
        footer={
          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
            <Button onClick={guardClose}>{t('common:button.cancel')}</Button>
            <Button type="primary" loading={verifying} onClick={handleVerify}>
              {t('common:button.ok')}
            </Button>
          </div>
        }
      >
        <Form form={bindForm} layout="horizontal" labelCol={{ flex: '0 0 120px' }} wrapperCol={{ flex: 'auto' }} onValuesChange={handleValuesChange}>
          <Form.Item
            name="email"
            label={t('accountBind.newEmail')}
            rules={[{ required: true, message: t('accountBind.newEmailRequired') }]}
          >
            <Input placeholder={t('accountBind.newEmailPlaceholder')} />
          </Form.Item>
          <Form.Item
            name="code"
            label={t('accountBind.vcode')}
            rules={[{ required: true, message: t('accountBind.vcodeRequired') }]}
          >
            <div style={{ display: 'flex', gap: 8 }}>
              <Input placeholder={t('accountBind.vcodePlaceholder')} />
              <Button disabled={codeSent} loading={sending} onClick={handleSendCode}>
                {codeSent ? t('accountBind.codeSent') : t('accountBind.sendVcode')}
              </Button>
            </div>
          </Form.Item>
        </Form>
      </Drawer>
    </>
  );
};

export default AccountBind;
